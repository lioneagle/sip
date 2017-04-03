package sipparser

import (
	"bytes"
	//"fmt"
	"unsafe"
)

type SipUri struct {
	isSecure bool
	user     AbnfToken
	password AbnfToken
	hostport SipHostPort
	params   SipUriParams
	headers  SipUriHeaders
}

func NewSipUri(context *ParseContext) (*SipUri, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipUri{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}
	(*SipUri)(unsafe.Pointer(mem)).isSecure = false
	(*SipUri)(unsafe.Pointer(mem)).Init()
	return (*SipUri)(unsafe.Pointer(mem)), addr
}

func (this *SipUri) Init() {
	this.user.Init()
	this.password.Init()
	this.hostport.Init()
	this.params.Init()
	this.headers.Init()
}

func (this *SipUri) SetSipUri()      { this.isSecure = false }
func (this *SipUri) SetSipsUri()     { this.isSecure = true }
func (this *SipUri) IsSipUri() bool  { return !this.isSecure }
func (this *SipUri) IsSipsUri() bool { return this.isSecure }

func (this *SipUri) Scheme() string {
	if this.isSecure {
		return "sips"
	}
	return "sip"
}

func (this *SipUri) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {

	newPos, err = this.ParseScheme(context, src, pos)
	if err != nil {
		return newPos, err
	}

	return this.ParseAfterScheme(context, src, newPos)
}

func (this *SipUri) ParseAfterScheme(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos
	this.Init()

	newPos, err = this.ParseUserinfo(context, src, newPos)
	if err != nil {
		return newPos, err
	}
	//return newPos, nil

	newPos, err = this.hostport.Parse(context, src, newPos)
	if err != nil {
		return newPos, err
	}

	if newPos >= len(src) {
		return newPos, nil
	}

	if src[newPos] == ';' {
		newPos, err = this.params.Parse(context, src, newPos+1)
		if err != nil {
			return newPos, err
		}
	}

	if newPos >= len(src) {
		return newPos, nil
	}

	if src[newPos] == '?' {
		newPos, err = this.headers.Parse(context, src, newPos+1)
		if err != nil {
			return newPos, err
		}
	}

	return newPos, err
}

func (this *SipUri) ParseAfterSchemeWithoutParam(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos
	this.Init()

	newPos, err = this.ParseUserinfo(context, src, newPos)
	if err != nil {
		return newPos, err
	}
	//return newPos, nil

	newPos, err = this.hostport.Parse(context, src, newPos)
	if err != nil {
		return newPos, err
	}

	return newPos, nil
}

func (this *SipUri) Encode(context *ParseContext, buf *bytes.Buffer) {
	buf.WriteString(this.Scheme())
	buf.WriteByte(':')

	if this.user.Exist() {
		buf.Write(Escape(this.user.GetAsByteSlice(context), IsSipUser))
		if this.password.Exist() {
			buf.WriteByte(':')
			buf.Write(Escape(this.password.GetAsByteSlice(context), IsSipPassword))
		}
		buf.WriteByte('@')
	}

	this.hostport.Encode(context, buf)

	if !this.params.Empty() {
		buf.WriteByte(';')
		this.params.Encode(context, buf)
	}

	if !this.headers.Empty() {
		buf.WriteByte('?')
		this.headers.Encode(context, buf)
	}
}

func (this *SipUri) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

func (this *SipUri) Equal(context *ParseContext, uri URI) bool {
	return this.EqualRFC3261(context, uri)
}

func (this *SipUri) EqualRFC3261(context *ParseContext, uri URI) bool {
	rhs, ok := uri.(*SipUri)
	if !ok {
		return false
	}

	if (this.isSecure && !rhs.isSecure) || (!this.isSecure && rhs.isSecure) {
		return false
	}

	if !this.EqualUserinfo(context, rhs) {
		return false
	}

	if !this.hostport.Equal(context, &rhs.hostport) {
		return false
	}

	if !this.params.EqualRFC3261(context, &rhs.params) {
		return false
	}

	if !this.headers.EqualRFC3261(context, &rhs.headers) {
		return false
	}

	return true
}

func (this *SipUri) EqualUserinfo(context *ParseContext, rhs *SipUri) bool {
	if (this.user.Exist() && !rhs.user.Exist()) || (!this.user.Exist() && rhs.user.Exist()) {
		return false
	}

	if (this.password.Exist() && !rhs.password.Exist()) || (!this.password.Exist() && rhs.password.Exist()) {
		return false
	}

	if this.user.Exist() && !this.user.Equal(context, &rhs.user) {
		return false
	}

	if this.password.Exist() && !this.password.Equal(context, &rhs.user) {
		return false
	}

	return true
}

func (this *SipUri) ParseScheme(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	src1 := src[pos:]
	if len(src) >= 4 && ((src1[0] | 0x20) == 's') && ((src1[1] | 0x20) == 'i') && ((src1[2] | 0x20) == 'p') && ((src1[3] | 0x20) == ':') {
		this.SetSipUri()
		return pos + 4, nil
	}

	if len(src) >= 5 && ((src1[0] | 0x20) == 's') && ((src1[1] | 0x20) == 'i') && ((src1[2] | 0x20) == 'p') &&
		((src1[3] | 0x20) == 's') && ((src1[4] | 0x20) == ':') {
		this.SetSipsUri()
		return pos + 5, nil
	}

	return 0, &AbnfError{"sip-uri parse: parse scheme failed: not sip-uri nor sips-uri", src, newPos}

	/*
		newPos, scheme, err := ParseUriScheme(src, pos)
		if err != nil {
			return newPos, err
		}

		if EqualNoCase(scheme.value, StringToByteSlice("sips")) {
			this.SetSipsUri()
		} else if !EqualNoCase(scheme.value, StringToByteSlice("sip")) {
			return newPos, &AbnfError{"sip-uri parse: parse scheme failed: not sip-uri nor sips-uri", src, newPos}
		} else {
			this.SetSipUri()
		}

		return newPos, nil
	*/
}

func (this *SipUri) ParseUserinfo(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos
	hasUserinfo := findUserinfo(src, newPos)
	if hasUserinfo {
		newPos, err = this.user.ParseEscapable(context, src, newPos, IsSipUser)
		if err != nil {
			return newPos, err
		}

		if newPos >= len(src) {
			return newPos, &AbnfError{"sip-uri parse: parse user-info failed: reach end after user", src, newPos}
		}

		if src[newPos] == ':' {
			newPos++
			if newPos >= len(src) {
				return newPos, &AbnfError{"sip-uri parse: parse user-info failed: reach end after password :", src, newPos}
			}

			if IsSipPassword(src[newPos]) {
				newPos, err = this.password.ParseEscapable(context, src, newPos, IsSipPassword)
				if err != nil {
					return newPos, err
				}
			}
			this.password.SetExist()
		}

		if newPos >= len(src) {
			return newPos, &AbnfError{"sip-uri parse: parse user-info failed: reach end, and no '@'", src, newPos}
		}

		if src[newPos] != '@' {
			return newPos, &AbnfError{"sip-uri parse: parse user-info failed: no '@'", src, newPos}
		}

		newPos++
	}

	return newPos, nil
}

func findUserinfo(src []byte, pos int) bool {
	for _, v := range src[pos:] {
		if v == '@' {
			return true
		} else if v == '>' || IsLwsChar(v) {
			break
		}
	}
	return false
}
