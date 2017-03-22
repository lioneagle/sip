package sipparser2

import (
	//"fmt"
	"bytes"
	"strings"
)

type TelUri struct {
	isGlobalNumber bool
	number         AbnfToken
	user           AbnfToken
	context        TelUriContext
	params         TelUriParams
}

func (this *TelUri) SetGlobalNumber()     { this.isGlobalNumber = true }
func (this *TelUri) SetLocalNumber()      { this.isGlobalNumber = false }
func (this *TelUri) IsGlobalNumber() bool { return this.isGlobalNumber }
func (this *TelUri) IsLocalNumber() bool  { return !this.isGlobalNumber }

func NewTelUri() *TelUri {
	uri := &TelUri{}
	uri.params.Init()
	return uri
}

func (this *TelUri) Scheme() string {
	return "tel"
}

func (this *TelUri) Parse(src []byte, pos int) (newPos int, err error) {

	newPos, err = this.ParseScheme(src, pos)
	if err != nil {
		return newPos, err
	}

	return this.ParseAfterScheme(src, newPos)
}

func (this *TelUri) String() string {
	str := "tel:"

	str += this.number.String()

	if !this.isGlobalNumber {
		str += this.context.String()
	}

	if !this.params.Empty() {
		str += this.params.String()
	}

	return str
}

func (this *TelUri) Encode(buf *bytes.Buffer) {
	buf.WriteString("tel:")
	this.number.Encode(buf)

	if !this.isGlobalNumber {
		this.context.Encode(buf)
	}

	if !this.params.Empty() {
		this.params.Encode(buf)
	}
}

func (this *TelUri) Equal(uri URI) bool {
	rhs, ok := uri.(*TelUri)
	if !ok {
		return false
	}

	if (this.isGlobalNumber && !rhs.isGlobalNumber) || (!this.isGlobalNumber && rhs.isGlobalNumber) {
		return false
	}

	if !this.number.Equal(&rhs.number) {
		return false
	}

	if !this.context.desc.EqualNoCase(&rhs.context.desc) {
		return false
	}

	if !this.EqualParams(rhs) {
		return false
	}

	return true
}

func (this *TelUri) EqualParams(rhs *TelUri) bool {
	if this.params.Size() != rhs.params.Size() {
		return false
	}

	for _, v := range this.params.maps {
		param, ok := rhs.params.GetParam(v.name.String())
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

func (this *TelUri) ParseScheme(src []byte, pos int) (newPos int, err error) {
	newPos, scheme, err := ParseUriScheme(src, pos)
	if err != nil {
		return newPos, err
	}

	if !EqualNoCase(scheme.value, Str2bytes("tel")) {
		return newPos, &AbnfError{"tel-uri parse: parse scheme failed: not sip-uri nor sips-uri", src, newPos}
	}

	return newPos, nil
}

func (this *TelUri) ParseAfterScheme(src []byte, pos int) (newPos int, err error) {
	newPos = pos

	newPos, err = this.ParseNumber(src, newPos)
	if err != nil {
		return newPos, err
	}

	if newPos >= len(src) {
		return newPos, nil
	}

	return this.ParseParams(src, newPos)
}

func (this *TelUri) ParseNumber(src []byte, pos int) (newPos int, err error) {
	if src[pos] == '+' {
		this.SetGlobalNumber()
		return this.ParseGlobalNumber(src, pos)
	}

	this.SetLocalNumber()
	return this.ParseLocalNumber(src, pos)
}

func (this *TelUri) ParseGlobalNumber(src []byte, pos int) (newPos int, err error) {
	newPos, err = this.number.Parse(src, pos+1, IsTelPhoneDigit)
	if err != nil {
		return newPos, err
	}

	this.number.value = src[pos:newPos]

	this.number.value = this.RemoveVisualSeperator(this.number.value)

	if this.number.Size() <= 1 {
		return newPos, &AbnfError{"tel-uri parse: parse global-number failed: empty number", src, newPos}
	}

	this.number.SetExist()

	return newPos, nil
}

func (this *TelUri) ParseLocalNumber(src []byte, pos int) (newPos int, err error) {
	newPos, err = this.number.Parse(src, pos, IsTelPhoneDigitHex)
	if err != nil {
		return newPos, err
	}

	this.number.value = this.RemoveVisualSeperator(this.number.value)

	if this.number.Empty() {
		return newPos, &AbnfError{"tel-uri parse: parse global-number failed: empty number", src, newPos}
	}

	this.number.SetExist()

	return newPos, nil
}

func (this *TelUri) RemoveVisualSeperator(number []byte) []byte {
	newNumber := make([]byte, 0)
	for _, v := range number {
		if !IsTelVisualSperator(v) {
			newNumber = append(newNumber, v)
		}
	}
	return newNumber
}

func (this *TelUri) ParseParams(src []byte, pos int) (newPos int, err error) {
	newPos = pos
	if newPos >= len(src) {
		return newPos, &AbnfError{"tel-uri parse: parse tel-uri param failed: reach end after ';'", src, newPos}
	}

	for newPos < len(src) {
		if src[newPos] != ';' {
			return newPos, nil
		}

		param := &TelUriParam{}

		newPos, err = param.Parse(src, newPos+1)
		if err != nil {
			return newPos, err
		}

		name := param.name.ToLower()
		if name == "phone-context" {
			this.context.exist = true
			this.context.isDomainName = (param.value.value[0] != '+')
			this.context.desc = param.value
			if !this.context.isDomainName {
				this.context.desc.value = this.RemoveVisualSeperator(this.context.desc.value)
			}
		} else {
			this.params.orders = append(this.params.orders, name)
			this.params.maps[name] = param
		}
	}

	return newPos, nil
}

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

func (this *TelUriParam) Parse(src []byte, pos int) (newPos int, err error) {
	newPos, err = this.name.ParseEscapable(src, pos, IsTelPname)
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
		newPos, err = this.value.ParseEscapable(src, newPos+1, IsTelPvalue)
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
