package sipparser3

import (
	"bytes"
	//"fmt"
	//"strings"
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

	if newPos >= len(src) {
		return newPos, nil
	}

	if src[newPos] == '=' {
		newPos, err = this.value.ParseEscapable(context, src, newPos+1, IsSipPvalue)
		if err != nil {
			return newPos, err
		}
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
	return AbnfEncoderToString(this)
}

type SipUriParams struct {
	params []SipUriParam
}

func NewSipUriParams() *SipUriParams {
	ret := &SipUriParams{}
	ret.Init()
	return ret
}

func (this *SipUriParams) Init() {
	//this.params = make([]SipUriParam, 0, 2)
	if len(this.params) != 0 {
		this.params = make([]SipUriParam, 0, 2)
	}
}

func (this *SipUriParams) Size() int   { return len(this.params) }
func (this *SipUriParams) Empty() bool { return len(this.params) == 0 }
func (this *SipUriParams) GetParam(name string) (val *SipUriParam, ok bool) {
	for i, v := range this.params {
		if v.name.EqualStringNoCase(name) {
			return &this.params[i], true
		}
	}
	return nil, false
}

func (this *SipUriParams) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos
	if newPos >= len(src) {
		return newPos, &AbnfError{"sip-uri parse: parse sip-uri param failed: reach end after ';'", src, newPos}
	}

	for newPos < len(src) {
		param := SipUriParam{}
		newPos, err = param.Parse(context, src, newPos)
		if err != nil {
			return newPos, err
		}
		this.params = append(this.params, param)

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

	for _, v := range params1.params {
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
			return param1.value.EqualNoCase(&param2.value)
		}
	}

	return true
}

func (this *SipUriParams) Encode(buf *bytes.Buffer) {
	for i, v := range this.params {
		if i > 0 {
			buf.WriteByte(';')
		}
		v.Encode(buf)
	}
}

func (this *SipUriParams) String() string {
	return AbnfEncoderToString(this)
}
