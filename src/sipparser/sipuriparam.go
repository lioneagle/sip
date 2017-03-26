package sipparser

import (
	"bytes"
	//"fmt"
	"strings"
)

type SipUriParam struct {
	name  AbnfToken
	value AbnfToken
}

func (this *SipUriParam) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos, err = this.name.ParseEscapable(context, src, pos, IsSipPname)
	if err != nil {
		return newPos, err
	}

	if this.name.Empty() {
		return newPos, &AbnfError{"sip-uri parse: parse sip-uri param failed: empty pname", src, newPos}
	}

	this.name.SetExist()

	if newPos >= len(src) {
		return newPos, nil
	}

	if src[newPos] == '=' {
		newPos, err = this.value.ParseEscapable(context, src, newPos+1, IsSipPvalue)
		if err != nil {
			return newPos, err
		}

		if this.value.Empty() {
			return newPos, &AbnfError{"sip-uri parse: parse sip-uri param failed: empty pvalue", src, newPos}
		}
		this.value.SetExist()
	}
	return newPos, nil
}

func (this *SipUriParam) Encode(buf *bytes.Buffer) {
	buf.Write(Escape(this.name.value, IsSipPname))
	if this.value.Exist() {
		buf.WriteByte('=')
		buf.Write(Escape(this.value.value, IsSipPvalue))
	}
}

func (this *SipUriParam) String() string {
	str := string(Escape([]byte(this.name.String()), IsSipPname))
	if this.value.Exist() {
		str += "="
		str += string(Escape([]byte(this.value.String()), IsSipPvalue))
	}
	return str
}

type SipUriParams struct {
	orders []string
	maps   map[string]*SipUriParam
}

func NewSipUriParams() *SipUriParams {
	return &SipUriParams{orders: make([]string, 0), maps: make(map[string]*SipUriParam)}
}

func (this *SipUriParams) Init() {
	this.orders = make([]string, 0)
	this.maps = make(map[string]*SipUriParam)
}

func (this *SipUriParams) Encode(buf *bytes.Buffer) {
	for i, v := range this.orders {
		if i > 0 {
			buf.WriteByte(';')
		}
		this.maps[v].Encode(buf)
	}
}

func (this *SipUriParams) String() string {
	if len(this.maps) == 0 {
		return ""
	}

	str := ""
	for i, v := range this.orders {
		if i > 0 {
			str += ";"
		}
		str += this.maps[v].String()
	}
	return str
}

func (this *SipUriParams) Size() int   { return len(this.maps) }
func (this *SipUriParams) Empty() bool { return len(this.maps) == 0 }
func (this *SipUriParams) GetParam(name string) (val *SipUriParam, ok bool) {
	val, ok = this.maps[strings.ToLower(name)]
	return val, ok
}

func (this *SipUriParams) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos
	if newPos >= len(src) {
		return newPos, &AbnfError{"sip-uri parse: parse sip-uri param failed: reach end after ';'", src, newPos}
	}

	for newPos < len(src) {
		param := &SipUriParam{}
		newPos, err = param.Parse(context, src, newPos)
		if err != nil {
			return newPos, err
		}
		name := param.name.ToLower()
		this.orders = append(this.orders, name)
		this.maps[name] = param

		if newPos >= len(src) {
			return newPos, nil
		}

		if src[newPos] != ';' {
			return newPos, nil
		}
		newPos++
	}

	return newPos, nil
}

func (this *SipUriParams) EqualRFC3261(rhs *SipUriParams) bool {
	params1 := this
	params2 := rhs

	if params1.Size() < params2.Size() {
		params1, params2 = params2, params1
	}

	if !params1.equalSpecParamsRFC3261(params2) {
		return false
	}

	for _, v := range params1.maps {
		param, ok := params2.GetParam(v.name.String())
		if ok {
			if !param.value.EqualNoCase(&v.value) {
				return false
			}
		}
	}
	return true
}

func (this *SipUriParams) equalSpecParamsRFC3261(rhs *SipUriParams) bool {
	specParams := []string{"user", "ttl", "method"}

	for _, v := range specParams {
		param1, ok := this.GetParam(v)
		if ok {
			param2, ok := rhs.GetParam(v)
			if !ok {
				return false
			}
			ret := param1.value.EqualNoCase(&param2.value)
			return ret

		}
	}

	return true
}
