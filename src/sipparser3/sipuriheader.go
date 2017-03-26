package sipparser3

import (
	"bytes"
	//"fmt"
	//"strings"
)

type SipUriHeader struct {
	name  AbnfToken
	value AbnfToken
}

func (this *SipUriHeader) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos, err = this.name.ParseEscapable(context, src, pos, IsSipHname)
	if err != nil {
		return newPos, err
	}

	if this.name.Empty() {
		return newPos, &AbnfError{"sip-uri parse: parse sip-uri header failed: empty hname", src, newPos}
	}
	this.name.SetExist()

	if newPos >= len(src) {
		return newPos, &AbnfError{"sip-uri parse: parse header failed: no = after hname", src, newPos}
	}

	if src[newPos] != '=' {
		return newPos, &AbnfError{"sip-uri parse: parse header failed: no = after hname", src, newPos}
	}

	newPos, err = this.value.ParseEscapable(context, src, newPos+1, IsSipHvalue)
	if err != nil {
		return newPos, err
	}

	this.value.SetExist()

	return newPos, nil
}

func (this *SipUriHeader) Encode(buf *bytes.Buffer) {
	buf.Write(Escape(this.name.value, IsSipHname))
	buf.WriteByte('=')
	if this.value.Exist() {
		buf.Write(Escape(this.value.value, IsSipHvalue))
	}
}

func (this *SipUriHeader) String() string {
	str := string(Escape([]byte(this.name.String()), IsSipHname))
	str += "="
	if this.value.Exist() {
		str += string(Escape([]byte(this.value.String()), IsSipHvalue))
	}
	return str
}

type SipUriHeaders struct {
	headers []SipUriHeader
}

func NewSipUriHeaders() *SipUriHeaders {
	ret := &SipUriHeaders{}
	ret.Init()
	return ret
}

func (this *SipUriHeaders) Init() {
	this.headers = make([]SipUriHeader, 0, 2)
}

func (this *SipUriHeaders) Size() int   { return len(this.headers) }
func (this *SipUriHeaders) Empty() bool { return len(this.headers) == 0 }
func (this *SipUriHeaders) GetHeader(name string) (val *SipUriHeader, ok bool) {
	for _, v := range this.headers {
		if v.name.EqualStringNoCase(name) {
			return &v, true
		}
	}
	return nil, false
}

func (this *SipUriHeaders) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos
	if newPos >= len(src) {
		return newPos, &AbnfError{"sip-uri parse: parse uri-header failed: reach end after ';'", src, newPos}
	}

	for newPos < len(src) {
		header := SipUriHeader{}
		newPos, err = header.Parse(context, src, newPos)
		if err != nil {
			return newPos, err
		}
		this.headers = append(this.headers, header)

		if newPos >= len(src) {
			return newPos, nil
		}

		if src[newPos] != '&' {
			return newPos, nil
		}
		newPos++
	}

	return newPos, err
}

func (this *SipUriHeaders) EqualRFC3261(rhs *SipUriHeaders) bool {
	if this.Size() != rhs.Size() {
		return false
	}

	for _, v := range this.headers {
		header, ok := rhs.GetHeader(v.name.String())
		if ok {
			if !header.value.EqualNoCase(&v.value) {
				return false
			}
		}
	}
	return true
}

func (this *SipUriHeaders) Encode(buf *bytes.Buffer) {
	for i, v := range this.headers {
		if i > 0 {
			buf.WriteByte('&')
		}
		v.Encode(buf)
	}

}

func (this *SipUriHeaders) String() string {
	str := ""
	for i, v := range this.headers {
		if i > 0 {
			str += "&"
		}
		str += v.String()
	}
	return str
}
