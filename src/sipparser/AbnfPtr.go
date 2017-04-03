package sipparser

import (
	"strconv"
	"unsafe"
)

type AbnfPtr int32

const ABNF_PTR_NIL = AbnfPtr(-1)

func (this AbnfPtr) GetMemAddr(context *ParseContext) *byte {
	return (*byte)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetUintptr(context *ParseContext) uintptr {
	return (uintptr)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetAbnfToken(context *ParseContext) *AbnfToken {
	return (*AbnfToken)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetAbnfListNode(context *ParseContext) *AbnfListNode {
	return (*AbnfListNode)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) String() string {
	if this == ABNF_PTR_NIL {
		return "nil"
	}
	return strconv.FormatUint(uint64(this), 10)
}
