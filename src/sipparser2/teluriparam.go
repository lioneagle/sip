package sipparser2

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
	buf.WriteString(";phone-context=")
	buf.Write(Escape(this.desc.value, IsTelPvalue))
}

func (this *TelUriContext) String() string {
	str := ";phone-context="
        str += Bytes2str(Escape(this.desc.value, IsTelPvalue))
	return str
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

	if this.name.Empty() {
		return newPos, &AbnfError{"tel-uri parse: parse tel-uri param failed: empty pname", src, newPos}
	}

	this.name.SetExist()

	if newPos >= len(src) {
		return newPos, nil
	}

	if src[newPos] == '=' {
		newPos, err = this.value.ParseEscapable(context, src, newPos+1, IsTelPvalue)
		if err != nil {
			return newPos, err
		}

		if this.value.Empty() {
			return newPos, &AbnfError{"tel-uri parse: parse tel-uri param failed: empty pvalue", src, newPos}
		}
		this.value.SetExist()
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
	str := Bytes2str(Escape(this.name.value, IsTelPname))
	if this.value.Exist() {
		str += "="
		str += Bytes2str(Escape(this.value.value, IsTelPvalue))
	}
	return str
}

type TelUriParams struct {
	AbnfList
}

func (this *TelUriParams) Size() int   { return this.Len() }
func (this *TelUriParams) Empty() bool { return this.Len() == 0 }

func (this *TelUriParams) GetParam(name string) (val *TelUriParam, ok bool) {
	for e := this.Front(); e != nil; e = e.Next() {
		param := e.Value.(*TelUriParam)

		if param.name.EqualStringNoCase(name) {
			return param, true
		}
	}
	return nil, false
}

func (this *TelUriParams) Equal(rhs *TelUriParams) bool {
	if this.Size() != rhs.Size() {
		return false
	}

	for e := this.Front(); e != nil; e = e.Next() {
		param1 := e.Value.(*TelUriParam)
		param2, ok := rhs.GetParam(param1.name.String())
		if ok {
			if !param2.value.EqualNoCase(&param1.value) {
				return false
			}
		} else {
			return false
		}

	}
	return true
}

func (this *TelUriParams) Encode(buf *bytes.Buffer) {
	for e := this.Front(); e != nil; e = e.Next() {
		buf.WriteByte(';')
		param := e.Value.(*TelUriParam)
		param.Encode(buf)
	}

}

func (this *TelUriParams) String() string {
	str := ""
	for e := this.Front(); e != nil; e = e.Next() {
		str += ";"
		param := e.Value.(*TelUriParam)
		str += param.String()
	}

	return str
}
