package sipparser

import (
	"bytes"
	//"fmt"
	"reflect"
	"unsafe"
)

type AbnfBuf struct {
	//AbnfPtr
	addr AbnfPtr
	size int32
}

func NewAbnfBuf(context *ParseContext) (*AbnfBuf, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(AbnfBuf{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*AbnfBuf)(unsafe.Pointer(mem)).Init()
	return (*AbnfBuf)(unsafe.Pointer(mem)), addr
}

func (this *AbnfBuf) Init() {
	this.addr = ABNF_PTR_NIL
	this.size = 0
}

func (this *AbnfBuf) SetByteSlice(context *ParseContext, buf []byte) {
	size := int32(len(buf))
	if size == 0 {
		this.size = 0
		return
	}

	if this.size < size {
		mem, addr := context.allocator.Alloc(size)
		if mem == nil {
			return
		}
		this.addr = addr

	}

	this.size = size
	copy(this.GetAsByteSlice(context), buf)
	return
}

func (this *AbnfBuf) SetString(context *ParseContext, str string) {
	this.SetByteSlice(context, StringToByteSlice(str))
}

func (this *AbnfBuf) GetAsByteSlice(context *ParseContext) []byte {
	if this.addr == ABNF_PTR_NIL {
		return nil
	}
	header := reflect.SliceHeader{Data: this.addr.GetUintptr(context), Len: int(this.size), Cap: int(this.size)}
	return *(*[]byte)(unsafe.Pointer(&header))
}

func (this *AbnfBuf) GetAsString(context *ParseContext) string {
	if this.addr == ABNF_PTR_NIL {
		return ""
	}
	header := reflect.StringHeader{Data: this.addr.GetUintptr(context), Len: int(this.size)}
	return *(*string)(unsafe.Pointer(&header))
}

func (this *AbnfBuf) Size() int32 {
	return this.size
}

func (this *AbnfBuf) Equal(context *ParseContext, rhs *AbnfBuf) bool {
	if this.addr == rhs.addr {
		return true
	}
	return bytes.Equal(this.GetAsByteSlice(context), rhs.GetAsByteSlice(context))
}

func (this *AbnfBuf) EqualNoCase(context *ParseContext, rhs *AbnfBuf) bool {
	if this.addr == ABNF_PTR_NIL || rhs.addr == ABNF_PTR_NIL {
		return false
	}
	return EqualNoCase(this.GetAsByteSlice(context), rhs.GetAsByteSlice(context))
}

func (this *AbnfBuf) EqualByteSlice(context *ParseContext, rhs []byte) bool {
	if this.addr == ABNF_PTR_NIL {
		return false
	}
	return bytes.Equal(this.GetAsByteSlice(context), rhs)
}

func (this *AbnfBuf) EqualByteSliceNoCase(context *ParseContext, rhs []byte) bool {
	if this.addr == ABNF_PTR_NIL {
		return false
	}
	return EqualNoCase(this.GetAsByteSlice(context), rhs)
}

func (this *AbnfBuf) EqualString(context *ParseContext, str string) bool {
	return this.EqualByteSlice(context, StringToByteSlice(str))
}

func (this *AbnfBuf) EqualStringNoCase(context *ParseContext, str string) bool {
	return this.EqualByteSliceNoCase(context, StringToByteSlice(str))
}

func (this *AbnfBuf) Encode(context *ParseContext, buf *bytes.Buffer) {
	if this.addr != ABNF_PTR_NIL && this.size != 0 {
		buf.Write(this.GetAsByteSlice(context))
	}
}

func (this *AbnfBuf) String(context *ParseContext) string {
	if this.addr != ABNF_PTR_NIL && this.size != 0 {
		return AbnfEncoderToString(context, this)
	}
	return ""
}
