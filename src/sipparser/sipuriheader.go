package sipparser

import (
	//"fmt"
	"bytes"
	"strings"
)

type SipUriHeader struct {
	name  AbnfToken
	value AbnfToken
}

func (this *SipUriHeader) Parse(src []byte, pos int) (newPos int, err error) {
	newPos, err = this.name.ParseEscapable(src, pos, IsSipHname)
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

	newPos, err = this.value.ParseEscapable(src, newPos+1, IsSipHvalue)
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
	orders []string
	maps   map[string]*SipUriHeader
}

func NewSipUriHeaders() *SipUriHeaders {
	return &SipUriHeaders{orders: make([]string, 0), maps: make(map[string]*SipUriHeader)}
}

func (this *SipUriHeaders) Init() {
	this.orders = make([]string, 0)
	this.maps = make(map[string]*SipUriHeader)
}

func (this *SipUriHeaders) Size() int   { return len(this.maps) }
func (this *SipUriHeaders) Empty() bool { return len(this.maps) == 0 }
func (this *SipUriHeaders) GetHeader(name string) (val *SipUriHeader, ok bool) {
	val, ok = this.maps[strings.ToLower(name)]
	return val, ok
}

func (this *SipUriHeaders) Parse(src []byte, pos int) (newPos int, err error) {
	newPos = pos
	if newPos >= len(src) {
		return newPos, &AbnfError{"sip-uri parse: parse uri-header failed: reach end after ';'", src, newPos}
	}

	for newPos < len(src) {
		header := &SipUriHeader{}
		newPos, err = header.Parse(src, newPos)
		if err != nil {
			return newPos, err
		}
		name := header.name.ToLower()
		this.orders = append(this.orders, name)
		this.maps[name] = header

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

func (this *SipUriHeaders) Encode(buf *bytes.Buffer) {
	for i, v := range this.orders {
		if i > 0 {
			buf.WriteByte('&')
		}
		this.maps[v].Encode(buf)
	}
}

func (this *SipUriHeaders) String() string {
	if len(this.maps) == 0 {
		return ""
	}

	str := ""
	for i, v := range this.orders {
		if i > 0 {
			str += "&"
		}
		str += this.maps[v].String()
	}
	return str
}

func (this *SipUriHeaders) EqualRFC3261(rhs *SipUriHeaders) bool {
	if this.Size() != rhs.Size() {
		return false
	}

	for _, v := range this.maps {
		header, ok := rhs.GetHeader(v.name.String())
		if ok {
			if !header.value.EqualNoCase(&v.value) {
				return false
			}
		}
	}
	return true
}
