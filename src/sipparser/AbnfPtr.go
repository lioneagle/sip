package sipparser

import (
	"strconv"
	"unsafe"
)

type AbnfPtr uint32

const ABNF_PTR_NIL = AbnfPtr(0xffffffff)

func (this AbnfPtr) GetMemAddr(context *ParseContext) *byte {
	return (*byte)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetUintptr(context *ParseContext) uintptr {
	return (uintptr)(unsafe.Pointer(&context.allocator.mem[this]))
}

/*
func (this AbnfPtr) GetAbnfToken(context *ParseContext) *AbnfToken {
	return (*AbnfToken)(unsafe.Pointer(&context.allocator.mem[this]))
}*/

func (this AbnfPtr) GetAbnfBuf(context *ParseContext) *AbnfBuf {
	return (*AbnfBuf)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetAbnfListNode(context *ParseContext) *AbnfListNode {
	return (*AbnfListNode)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetAbnfList(context *ParseContext) *AbnfList {
	return (*AbnfList)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) String() string {
	if this == ABNF_PTR_NIL {
		return "nil"
	}
	return strconv.FormatUint(uint64(this), 10)
}
