package sipparser3

import (
	//"fmt"
	"bytes"
	//"strings"
)

type TelUriContext struct {
	exist        bool
	isDomainName bool
	desc         AbnfToken
}

func (this *TelUriContext) Exist() bool  { return this.exist }
func (this *TelUriContext) SetExist()    { this.exist = true }
func (this *TelUriContext) SetNonExist() { this.exist = false }

func (this *TelUriContext) Encode(buf *bytes.Buffer) {
	if this.exist {
		buf.WriteString(";phone-context=")
		buf.Write(Escape(this.desc.value, IsTelPvalue))
	}
}

func (this *TelUriContext) String() string {
	return AbnfEncoderToString(this)
}

type TelUriParam struct {
	name  AbnfToken
	value AbnfToken
}

func (this *TelUriParam) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos, err = this.name.ParseEscapable(context, src, pos, IsTelPname)
	if err != nil {
		return newPos, err
	}

	if newPos >= len(src) {
		return newPos, nil
	}

	if src[newPos] == '=' {
		newPos, err = this.value.ParseEscapable(context, src, newPos+1, IsTelPvalue)
		if err != nil {
			return newPos, err
		}
	}
	return newPos, nil
}

func (this *TelUriParam) Encode(buf *bytes.Buffer) {
	buf.Write(Escape(this.name.value, IsTelPname))
	if this.value.Exist() {
		buf.WriteByte('=')
		buf.Write(Escape(this.value.value, IsTelPvalue))
	}
}

func (this *TelUriParam) String() string {
	return AbnfEncoderToString(this)
}

type TelUriParams struct {
	params []TelUriParam
}

func (this *TelUriParams) Init() {
	this.params = make([]TelUriParam, 0, 2)
}

func (this *TelUriParams) Size() int   { return len(this.params) }
func (this *TelUriParams) Empty() bool { return len(this.params) == 0 }

func (this *TelUriParams) GetParam(name string) (val *TelUriParam, ok bool) {
	for i, v := range this.params {
		if v.name.EqualStringNoCase(name) {
			return &this.params[i], true
		}
	}
	return nil, false
}

func (this *TelUriParams) Equal(rhs *TelUriParams) bool {
	if this.Size() != rhs.Size() {
		return false
	}

	for _, v := range this.params {
		param, ok := rhs.GetParam(v.name.String())
		if ok {
			if !param.value.EqualNoCase(&v.value) {
				return false
			}
		} else {
			return false
		}

	}
	return true
}

func (this *TelUriParams) Encode(buf *bytes.Buffer) {
	for _, v := range this.params {
		buf.WriteByte(';')
		v.Encode(buf)
	}

}

func (this *TelUriParams) String() string {
	return AbnfEncoderToString(this)
}
