package sipparser

import (
	//"fmt"
	"bytes"
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
	/*var buf bytes.Buffer
	this.Encode(&buf)
	return buf.String()*/
	//*
	str := this.Scheme()
	str += ":"

	if this.user.Exist() {
		str += string(Escape([]byte(this.user.String()), IsSipUser))
		if this.password.Exist() {
			str += ":"
			str += string(Escape([]byte(this.password.String()), IsSipPassword))
		}
		str += "@"
	}

	str += this.hostport.String()

	if !this.params.Empty() {
		str += ";"
		str += this.params.String()
	}

	if !this.headers.Empty() {
		str += "?"
		str += this.headers.String()
	}

	return str //*/
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
	newPos, scheme, err := ParseUriScheme(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if EqualNoCase(scheme.value, []byte("sips")) {
		this.SetSipsUri()
	} else if !EqualNoCase(scheme.value, []byte("sip")) {
		return newPos, &AbnfError{"sip-uri parse: parse scheme failed: not sip-uri nor sips-uri", src, newPos}
	} else {
		this.SetSipUri()
	}

	return newPos, nil
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
	for newPos := pos; newPos < len(src); newPos++ {
		if src[newPos] == '@' {
			return true
		} else if src[newPos] == '>' || IsLwsChar(src[newPos]) {
			return false
		}
	}
	return false
}
