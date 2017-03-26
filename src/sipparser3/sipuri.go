package sipparser3

import (
	"bytes"
	//"fmt"
	//"strings"
)

type SipUri struct {
	isSecure bool
	user     AbnfToken
	password AbnfToken
	hostport SipHostPort
	params   SipUriParams
	headers  SipUriHeaders
}

func (this *SipUri) SetSipUri()      { this.isSecure = false }
func (this *SipUri) SetSipsUri()     { this.isSecure = true }
func (this *SipUri) IsSipUri() bool  { return !this.isSecure }
func (this *SipUri) IsSipsUri() bool { return this.isSecure }

func NewSipUri() *SipUri {
	uri := &SipUri{}
	uri.params.Init()
	uri.headers.Init()
	return uri
}

func (this *SipUri) Init() {
	this.user.SetNonExist()
	this.password.SetNonExist()
	this.hostport.Init()
	this.params.Init()
	this.headers.Init()
}

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

	newPos, err = this.ParseUserinfo(context, src, newPos)
	if err != nil {
		return newPos, err
	}
	//return newPos, nil

	newPos, err = this.hostport.Parse(src, newPos)
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

	newPos, err = this.ParseUserinfo(context, src, newPos)
	if err != nil {
		return newPos, err
	}
	//return newPos, nil

	newPos, err = this.hostport.Parse(src, newPos)
	if err != nil {
		return newPos, err
	}

	return newPos, nil
}

func (this *SipUri) Encode(buf *bytes.Buffer) {
	buf.WriteString(this.Scheme())
	buf.WriteByte(':')

	if this.user.Exist() {
		buf.Write(Escape(this.user.value, IsSipUser))
		if this.password.Exist() {
			buf.WriteByte(':')
			buf.Write(Escape(this.password.value, IsSipPassword))
		}
		buf.WriteByte('@')
	}

	this.hostport.Encode(buf)

	if !this.params.Empty() {
		buf.WriteByte(';')
		this.params.Encode(buf)
	}

	if !this.headers.Empty() {
		buf.WriteByte('?')
		this.headers.Encode(buf)
	}
}

func (this *SipUri) String() string {
	return AbnfEncoderToString(this)
}

func (this *SipUri) Equal(uri URI) bool {
	return this.EqualRFC3261(uri)
}

func (this *SipUri) EqualRFC3261(uri URI) bool {
	rhs, ok := uri.(*SipUri)
	if !ok {
		return false
	}

	if (this.isSecure && !rhs.isSecure) || (!this.isSecure && rhs.isSecure) {
		return false
	}

	if !this.EqualUserinfo(rhs) {
		return false
	}

	if !this.hostport.Equal(&rhs.hostport) {
		return false
	}

	if !this.params.EqualRFC3261(&rhs.params) {
		return false
	}

	if !this.headers.EqualRFC3261(&rhs.headers) {
		return false
	}

	return true
}

func (this *SipUri) EqualUserinfo(rhs *SipUri) bool {
	if (this.user.Exist() && !rhs.user.Exist()) || (!this.user.Exist() && rhs.user.Exist()) {
		return false
	}
	ret := this.user.Equal(&rhs.user) && this.password.Equal(&rhs.password)
	return ret
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

		if EqualNoCase(scheme.value, Str2bytes("sips")) {
			this.SetSipsUri()
		} else if !EqualNoCase(scheme.value, Str2bytes("sip")) {
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

		if this.user.Empty() {
			return newPos, &AbnfError{"sip-uri parse: parse user-info failed: empty user", src, newPos}
		}

		this.user.SetExist()

		if src[newPos] == ':' {
			newPos, err = this.password.ParseEscapable(context, src, newPos+1, IsSipPassword)
			if err != nil {
				return newPos, err
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
