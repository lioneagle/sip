package sipparser2

import (
	"bytes"
	//"fmt"
	//"strings"
)

type SipUriParam struct {
	name  AbnfToken
	value AbnfToken
}

func (this *SipUriParam) Parse(src []byte, pos int) (newPos int, err error) {
	newPos, err = this.name.ParseEscapable(src, pos, IsSipPname)
	if err != nil {
		return newPos, err
	}

	if this.name.Empty() {
		return newPos, &AbnfError{"sip-uri parse: parse sip-uri param failed: empty pname", src, newPos}
	}

	this.name.SetExist()

	if src[newPos] == '=' {
		newPos, err = this.value.ParseEscapable(src, newPos+1, IsSipPvalue)
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
	SipList
}

func NewSipUriParams() *SipUriParams {
	return &SipUriParams{}
}

func (this *SipUriParams) Size() int   { return this.Len() }
func (this *SipUriParams) Empty() bool { return this.Len() == 0 }
func (this *SipUriParams) GetParam(name string) (val *SipUriParam, ok bool) {
	for e := this.Front(); e != nil; e = e.Next() {
		param := e.Value.(*SipUriParam)

		if param.name.EqualStringNoCase(name) {
			return param, true
		}
	}
	return nil, false
}

func (this *SipUriParams) Parse(src []byte, pos int) (newPos int, err error) {
	newPos = pos
	if newPos >= len(src) {
		return newPos, &AbnfError{"sip-uri parse: parse sip-uri param failed: reach end after ';'", src, newPos}
	}

	for newPos < len(src) {
		param := &SipUriParam{}
		newPos, err = param.Parse(src, newPos)
		if err != nil {
			return newPos, err
		}
		this.PushBack(param)

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

	for e := params1.Front(); e != nil; e = e.Next() {
		param1 := e.Value.(*SipUriParam)
		param2, ok := params2.GetParam(param1.name.String())
		if ok {
			if !param2.value.EqualNoCase(&param1.value) {
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
			return param1.value.EqualNoCase(&param2.value)
		}
	}

	return true
}

func (this *SipUriParams) Encode(buf *bytes.Buffer) {
	for e := this.Front(); e != nil; e = e.Next() {
		if e != this.Front() {
			buf.WriteByte(';')
		}

		param := e.Value.(*SipUriParam)
		param.Encode(buf)
	}
}

func (this *SipUriParams) String() string {
	str := ""
	for e := this.Front(); e != nil; e = e.Next() {
		if e != this.Front() {
			str += ";"
		}

		param := e.Value.(*SipUriParam)
		str += param.String()
	}

	return str
}
