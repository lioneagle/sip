package sipparser4

import (
	//"fmt"
	"bytes"
	"strings"
)

type TelUriContext struct {
	exist        bool
	isDomainName bool
	desc         AbnfToken
}

func (this *TelUriContext) Encode(buf *bytes.Buffer) {
	buf.WriteString(";phone-context=")
	this.desc.Encode(buf)
}

func (this *TelUriContext) String() string {
	str := ";phone-context="
	str += this.desc.String()
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
	orders []string
	maps   map[string]*TelUriParam
}

func (this *TelUriParams) Init() {
	this.orders = make([]string, 0)
	this.maps = make(map[string]*TelUriParam)
}

func (this *TelUriParams) Size() int   { return len(this.maps) }
func (this *TelUriParams) Empty() bool { return len(this.maps) == 0 }

func (this *TelUriParams) GetParam(name string) (val *TelUriParam, ok bool) {
	val, ok = this.maps[strings.ToLower(name)]
	return val, ok
}

func (this *TelUriParams) Encode(buf *bytes.Buffer) {
	for _, v := range this.orders {
		buf.WriteByte(';')
		this.maps[v].Encode(buf)
	}
}

func (this *TelUriParams) String() string {
	if len(this.maps) == 0 {
		return ""
	}

	str := ""
	for _, v := range this.orders {
		str += ";"
		str += this.maps[v].String()
	}
	return str
}

func (this *TelUriParams) Equal(rhs *TelUriParams) bool {
	if this.Size() != rhs.Size() {
		return false
	}

	for _, v := range this.maps {
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
