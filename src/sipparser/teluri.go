package sipparser

import (
	//"fmt"
	"bytes"
	//"strings"
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

func (this *TelUri) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {

	newPos, err = this.ParseScheme(context, src, pos)
	if err != nil {
		return newPos, err
	}

	return this.ParseAfterScheme(context, src, newPos)
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

	if !this.params.Equal(&rhs.params) {
		return false
	}

	return true
}


func (this *TelUri) ParseScheme(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos, scheme, err := ParseUriScheme(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if !EqualNoCase(scheme.value, []byte("tel")) {
		return newPos, &AbnfError{"tel-uri parse: parse scheme failed: not sip-uri nor sips-uri", src, newPos}
	}

	return newPos, nil
}

func (this *TelUri) ParseAfterScheme(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos

	newPos, err = this.ParseNumber(context, src, newPos)
	if err != nil {
		return newPos, err
	}

	if newPos >= len(src) {
		return newPos, nil
	}

	return this.ParseParams(context, src, newPos)
}

func (this *TelUri) ParseNumber(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	if src[pos] == '+' {
		this.SetGlobalNumber()
		return this.ParseGlobalNumber(context, src, pos)
	}

	this.SetLocalNumber()
	return this.ParseLocalNumber(context, src, pos)
}

func (this *TelUri) ParseGlobalNumber(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos, err = this.number.Parse(context, src, pos+1, IsTelPhoneDigit)
	if err != nil {
		return newPos, err
	}

	this.number.value = src[pos:newPos]

	this.number.value = this.RemoveVisualSeperator(context, this.number.value)

	if this.number.Size() <= 1 {
		return newPos, &AbnfError{"tel-uri parse: parse global-number failed: empty number", src, newPos}
	}

	this.number.SetExist()

	return newPos, nil
}

func (this *TelUri) ParseLocalNumber(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos, err = this.number.Parse(context, src, pos, IsTelPhoneDigitHex)
	if err != nil {
		return newPos, err
	}

	this.number.value = this.RemoveVisualSeperator(context, this.number.value)

	if this.number.Empty() {
		return newPos, &AbnfError{"tel-uri parse: parse global-number failed: empty number", src, newPos}
	}

	this.number.SetExist()

	return newPos, nil
}

func (this *TelUri) RemoveVisualSeperator(context *ParseContext, number []byte) []byte {
	newNumber := make([]byte, 0)
	for _, v := range number {
		if !IsTelVisualSperator(v) {
			newNumber = append(newNumber, v)
		}
	}
	return newNumber
}

func (this *TelUri) ParseParams(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos
	if newPos >= len(src) {
		return newPos, &AbnfError{"tel-uri parse: parse tel-uri param failed: reach end after ';'", src, newPos}
	}

	for newPos < len(src) {
		if src[newPos] != ';' {
			return newPos, nil
		}

		param := &TelUriParam{}

		newPos, err = param.Parse(context, src, newPos+1)
		if err != nil {
			return newPos, err
		}

		name := param.name.ToLower()
		if name == "phone-context" {
			this.context.exist = true
			this.context.isDomainName = (param.value.value[0] != '+')
			this.context.desc = param.value
			if !this.context.isDomainName {
				this.context.desc.value = this.RemoveVisualSeperator(context, this.context.desc.value)
			}
		} else {
			this.params.orders = append(this.params.orders, name)
			this.params.maps[name] = param
		}
	}

	return newPos, nil
}
