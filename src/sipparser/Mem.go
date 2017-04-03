package sipparser

import (
	"fmt"
	"reflect"
	"unsafe"
)

const SLICE_HEADER_LEN = int32(unsafe.Sizeof(reflect.SliceHeader{}))
const SIP_MEM_ALIGN = int32(4)
const SIP_MEM_LIGN_MASK = ^(SIP_MEM_ALIGN - 1)
const SIP_MEM_LIGN_MASK2 = (SIP_MEM_ALIGN - 1)

func RoundToAlign(x, align int32) int32 {
	return (x + align - 1) & ^(align - 1)
}

type MemAllocator struct {
	mem           []byte
	used          int32
	allocNum      int32
	allocSliceNum int32
	freeAllNum    int32
}

func NewMemAllocator(capacity int32) *MemAllocator {
	ret := MemAllocator{}
	ret.Init(capacity)
	return &ret
}

func (this *MemAllocator) Init(capacity int32) *MemAllocator {
	this.used = 0
	this.mem = make([]byte, int(capacity))
	return this
}

func (this *MemAllocator) Used() int32 {
	return this.used
}

func (this *MemAllocator) ClearAllocNum() {
	this.allocNum = 0
}

func (this *MemAllocator) AllocNum() int32 {
	return this.allocNum
}

func (this *MemAllocator) FreeAllNum() int32 {
	return this.freeAllNum
}

func (this *MemAllocator) Capacity() int32 {
	return int32(cap(this.mem))
}

func (this *MemAllocator) Left() int32 {
	return int32(cap(this.mem)) - this.used
}

//func (this *MemAllocator) GetMem(addr uint32) (mem *byte, size uint32) {
func (this *MemAllocator) GetMem(addr int32) *byte {
	if addr >= this.Capacity() {
		panic("ERROR: out of memory range")
		//return nil
	}
	return (*byte)(unsafe.Pointer(&this.mem[addr]))
}

func (this *MemAllocator) Alloc(size int32) (mem *byte, addr AbnfPtr) {
	if print_mem {
		fmt.Printf("MEM: alloc_request = %d\n", size)
	}
	newSize := this.used + size + (this.used+size)%SIP_MEM_ALIGN
	//newSize := (this.used + size + SIP_MEM_LIGN_MASK2) & SIP_MEM_LIGN_MASK
	if newSize >= this.Left() {
		//panic("ERROR: out of memory")
		return nil, ABNF_PTR_NIL
	}
	this.allocNum++

	mem = (*byte)(unsafe.Pointer(&this.mem[this.used]))
	addr = AbnfPtr(this.used)
	this.used = newSize
	return mem, addr
}

func (this *MemAllocator) AllocByteSlice(size int32) []byte {
	mem, _ := this.Alloc(size + SLICE_HEADER_LEN)
	this.allocSliceNum++

	data := (*reflect.SliceHeader)(unsafe.Pointer(mem))
	data.Data = (uintptr)(unsafe.Pointer(&this.mem[this.used-size]))
	data.Len = 0
	data.Cap = int(size)

	return *(*[]byte)(unsafe.Pointer(data))
}

func (this *MemAllocator) FreeAll() {
	this.used = 0
	this.freeAllNum++
}

func (this *MemAllocator) FreePart(remain int32) {
	this.used = remain
	this.freeAllNum++
}

func GetBytesDataAddr(buf *[]byte) unsafe.Pointer {
	return unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(buf)).Data)
}

var g_allocator *MemAllocator = NewMemAllocator(1024 * 128)
