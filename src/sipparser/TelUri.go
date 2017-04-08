package sipparser

import (
	"bytes"
	//"fmt"
	"unsafe"
)

type TelUri struct {
	isGlobalNumber bool
	number         AbnfBuf
	context        TelUriContext
	params         TelUriParams
}

func NewTelUri(context *ParseContext) (*TelUri, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(TelUri{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*TelUri)(unsafe.Pointer(mem)).Init()
	return (*TelUri)(unsafe.Pointer(mem)), addr
}

func (this *TelUri) Init() {
	this.isGlobalNumber = false
	this.number.Init()
	this.context.Init()
	this.params.Init()
}

func (this *TelUri) SetGlobalNumber()     { this.isGlobalNumber = true }
func (this *TelUri) SetLocalNumber()      { this.isGlobalNumber = false }
func (this *TelUri) IsGlobalNumber() bool { return this.isGlobalNumber }
func (this *TelUri) IsLocalNumber() bool  { return !this.isGlobalNumber }

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

func (this *TelUri) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

func (this *TelUri) Encode(context *ParseContext, buf *bytes.Buffer) {
	buf.WriteString("tel:")
	this.number.Encode(context, buf)

	if !this.isGlobalNumber {
		this.context.Encode(context, buf)
	}

	if !this.params.Empty() {
		this.params.Encode(context, buf)
	}
}

func (this *TelUri) Equal(context *ParseContext, uri URI) bool {
	rhs, ok := uri.(*TelUri)
	if !ok {
		return false
	}

	if !this.equalNumber(context, rhs) {
		return false
	}

	if !this.context.Equal(context, &rhs.context) {
		return false
	}

	if !this.params.Equal(context, &rhs.params) {
		return false
	}

	return true
}

func (this *TelUri) equalNumber(context *ParseContext, rhs *TelUri) bool {
	if (this.isGlobalNumber && !rhs.isGlobalNumber) || (!this.isGlobalNumber && rhs.isGlobalNumber) {
		return false
	}
	return this.number.Equal(context, &rhs.number)
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
	this.Init()

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

	//this.number.SetByteSlice(context, src[pos:newPos])

	this.number.SetByteSlice(context, this.RemoveVisualSeperator(context, src[pos:newPos]))

	if this.number.Size() <= 1 {
		return newPos, &AbnfError{"tel-uri parse: global-number is empty", src, newPos}
	}

	return newPos, nil
}

func (this *TelUri) ParseLocalNumber(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos, err = this.number.Parse(context, src, pos, IsTelPhoneDigitHex)
	if err != nil {
		return newPos, err
	}

	this.number.SetByteSlice(context, this.RemoveVisualSeperator(context, this.number.GetAsByteSlice(context)))

	if this.number.Empty() {
		return newPos, &AbnfError{"tel-uri parse: local-number is empty", src, newPos}
	}

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

		param, addr := NewTelUriParam(context)
		if param == nil {
			return newPos, &AbnfError{"tel-uri parse: out of memory for tel-uri param", src, newPos}
		}
		newPos, err = param.Parse(context, src, newPos+1)
		if err != nil {
			return newPos, err
		}

		if param.name.EqualStringNoCase(context, "phone-context") {
			this.context.SetExist()
			this.context.isDomainName = (param.value.GetAsByteSlice(context)[0] != '+')
			this.context.desc = param.value
			if !this.context.isDomainName {
				this.context.desc.SetByteSlice(context, this.RemoveVisualSeperator(context, this.context.desc.GetAsByteSlice(context)))
			}
		} else {
			this.params.PushBack(context, addr)
		}
	}

	return newPos, nil
}
