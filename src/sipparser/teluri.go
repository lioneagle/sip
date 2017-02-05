package sipparser

import (
	//"fmt"
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

func (this *TelUri) ParseScheme(src []byte, pos int) (newPos int, err error) {
	newPos, scheme, err := ParseUriScheme(src, pos)
	if err != nil {
		return newPos, err
	}

	if !EqualNoCase(scheme.value, []byte("tel")) {
		return newPos, &AbnfError{"parse scheme failed: not sip-uri nor sips-uri", src, newPos}
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

	if this.number.Empty() {
		return newPos, &AbnfError{"parse global-number failed: empty number", src, newPos}
	}

	this.number.value = src[pos:newPos]

	/*@@TODO: remove visual-seperator */

	return newPos, nil
}

func (this *TelUri) ParseLocalNumber(src []byte, pos int) (newPos int, err error) {
	newPos, err = this.number.Parse(src, pos, IsTelPhoneDigitHex)
	if err != nil {
		return newPos, err
	}

	/*@@TODO: remove visual-seperator */

	if this.number.Empty() {
		return newPos, &AbnfError{"parse global-number failed: empty number", src, newPos}
	}

	return newPos, nil
}

func (this *TelUri) ParseParams(src []byte, pos int) (newPos int, err error) {
	newPos = pos
	if newPos >= len(src) {
		return newPos, &AbnfError{"parse tel-uri param failed: reach end after ';'", src, newPos}
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
		return newPos, &AbnfError{"parse tel-uri param failed: empty pname", src, newPos}
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
			return newPos, &AbnfError{"parse tel-uri param failed: empty pvalue", src, newPos}
		}
		this.value.SetExist()
	}
	return newPos, nil
}

func (this *TelUriParam) String() string {
	str := string(Escape([]byte(this.name.String()), IsTelPname))
	if this.value.Exist() {
		str += "="
		str += string(Escape([]byte(this.value.String()), IsTelPvalue))
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

func (this *TelUriParams) String() string {
	if len(this.maps) == 0 {
		return ""
	}

	str := ""
	for i, v := range this.orders {
		str += ";"
		str += this.maps[v].String()
		i++
	}
	return str
}
