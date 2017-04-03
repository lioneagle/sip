package sipparser

import (
	"bytes"
	//"fmt"
	"strings"
	"unsafe"
)

type AbnfToken struct {
	exist bool
	//value []byte
	value AbnfBuf
}

func NewAbnfToken(context *ParseContext) (*AbnfToken, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(AbnfToken{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*AbnfToken)(unsafe.Pointer(mem)).Init()
	return (*AbnfToken)(unsafe.Pointer(mem)), addr
}

func (this *AbnfToken) Init() {
	this.value.Init()
	this.exist = false
}

func (this *AbnfToken) Exist() bool  { return this.exist }
func (this *AbnfToken) Size() int32  { return this.value.Size() }
func (this *AbnfToken) Empty() bool  { return this.value.Size() == 0 }
func (this *AbnfToken) SetExist()    { this.exist = true }
func (this *AbnfToken) SetNonExist() { this.exist = false }

func (this *AbnfToken) SetValue(context *ParseContext, value []byte) {
	this.value.SetByteSlice(context, value)
}

func (this *AbnfToken) ToLower(context *ParseContext) string {
	return strings.ToLower(this.value.GetAsString(context))
}
func (this *AbnfToken) ToUpper(context *ParseContext) string {
	return strings.ToUpper(this.value.GetAsString(context))
}

func (this *AbnfToken) Equal(context *ParseContext, rhs *AbnfToken) bool {
	if this.exist != rhs.exist {
		return false
	}
	if !this.exist {
		return true
	}
	return this.value.Equal(context, &rhs.value)
}

func (this *AbnfToken) EqualNoCase(context *ParseContext, rhs *AbnfToken) bool {
	if this.exist != rhs.exist {
		return false
	}
	if !this.exist {
		return true
	}
	return this.value.EqualNoCase(context, &rhs.value)
}

func (this *AbnfToken) EqualBytes(context *ParseContext, rhs []byte) bool {
	if !this.exist {
		return false
	}
	return this.value.EqualByteSlice(context, rhs)
}

func (this *AbnfToken) EqualBytesNoCase(context *ParseContext, rhs []byte) bool {
	if !this.exist {
		return false
	}
	return this.value.EqualByteSliceNoCase(context, rhs)
}

func (this *AbnfToken) EqualString(context *ParseContext, str string) bool {
	if !this.exist {
		return false
	}
	return this.value.EqualString(context, str)
}

func (this *AbnfToken) EqualStringNoCase(context *ParseContext, str string) bool {
	if !this.exist {
		return false
	}

	return this.value.EqualStringNoCase(context, str)
}

func (this *AbnfToken) GetAsByteSlice(context *ParseContext) []byte {
	if !this.exist {
		return nil
	}
	return this.value.GetAsByteSlice(context)
}

func (this *AbnfToken) Parse(context *ParseContext, src []byte, pos int, inCharset AbnfIsInCharset) (newPos int, err error) {
	begin, end, newPos, err := parseToken(src, pos, inCharset)
	if err != nil {
		return newPos, err
	}
	if begin >= end {
		return newPos, &AbnfError{"AbnfToken parse: value is empty", src, newPos}
	}
	this.SetExist()
	this.value.SetByteSlice(context, src[begin:end])
	return newPos, nil
}

func (this *AbnfToken) ParseEscapable(context *ParseContext, src []byte, pos int, inCharset AbnfIsInCharset) (newPos int, err error) {
	begin, end, newPos, err := parseTokenEscapable(src, pos, inCharset)
	if err != nil {
		return newPos, err
	}

	if begin >= end {
		return newPos, &AbnfError{"AbnfToken ParseEscapable: value is empty", src, newPos}
	}

	this.SetExist()
	this.value.SetByteSlice(context, Unescape(context, src[begin:end]))
	return newPos, nil
}

func (this *AbnfToken) Encode(context *ParseContext, buf *bytes.Buffer) {
	if this.exist {
		this.value.Encode(context, buf)
	}
}

func (this *AbnfToken) String(context *ParseContext) string {
	if this.exist {
		return this.value.GetAsString(context)
	}
	return ""
}
