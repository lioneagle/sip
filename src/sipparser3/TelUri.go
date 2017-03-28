package sipparser3

import (
	"bytes"
	//"fmt"
	//"strings"
)

type TelUri struct {
	isGlobalNumber bool
	number         AbnfToken
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

func (this *TelUri) Init() {
	this.number.SetNonExist()
	this.context.SetNonExist()
	this.params.Init()
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
	return AbnfEncoderToString(this)
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
	src1 := src[pos:]
	if len(src) >= 4 && ((src1[0] | 0x20) == 't') && ((src1[1] | 0x20) == 'e') && ((src1[2] | 0x20) == 'l') && (src1[3] == ':') {
		return pos + 4, nil
	}

	return 0, &AbnfError{"tel-uri parse: parse scheme failed", src, newPos}
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

func (this *TelUri) ParseAfterSchemeWithoutParam(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	return this.ParseNumber(context, src, pos)
}

func (this *TelUri) ParseNumber(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos
	if newPos >= len(src) {
		return newPos, &AbnfError{"tel-uri parse: no number after scheme", src, newPos}
	}

	if src[pos] == '+' {
		this.SetGlobalNumber()
		return this.ParseGlobalNumber(context, src, newPos)
	}

	this.SetLocalNumber()
	return this.ParseLocalNumber(context, src, newPos)
}

func (this *TelUri) ParseGlobalNumber(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos, err = this.number.Parse(context, src, pos+1, IsTelPhoneDigit)
	if err != nil {
		return newPos, err
	}

	this.number.value = src[pos:newPos]

	this.number.value = this.RemoveVisualSeperator(context, this.number.value)

	if this.number.Size() <= 1 {
		return newPos, &AbnfError{"tel-uri parse: global-number is empty", src, newPos}
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
		return newPos, &AbnfError{"tel-uri parse: local-number is empty", src, newPos}
	}

	this.number.SetExist()

	return newPos, nil
}

func HasVisualSeperator(number []byte) bool {
	for _, v := range number {
		if IsTelVisualSperator(v) {
			return true
		}
	}
	return false
}

func (this *TelUri) RemoveVisualSeperator(context *ParseContext, number []byte) []byte {
	if !HasVisualSeperator(number) {
		return number
	}
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
		return newPos, &AbnfError{"tel-uri parse: reach end after ';'", src, newPos}
	}

	for newPos < len(src) {
		if src[newPos] != ';' {
			return newPos, nil
		}

		param := TelUriParam{}

		newPos, err = param.Parse(context, src, newPos+1)
		if err != nil {
			return newPos, err
		}

		if param.name.EqualStringNoCase("phone-context") {
			this.context.SetExist()
			this.context.isDomainName = (param.value.value[0] != '+')
			this.context.desc = param.value
			if !this.context.isDomainName {
				this.context.desc.value = this.RemoveVisualSeperator(context, this.context.desc.value)
			}
		} else {
			this.params.params = append(this.params.params, param)
		}
	}

	return newPos, nil
}
