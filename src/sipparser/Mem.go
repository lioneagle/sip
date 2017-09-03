package sipparser

import (
	"fmt"
	//"logger"
	"reflect"
	"strconv"
	"unsafe"
)

const SLICE_HEADER_LEN = int32(unsafe.Sizeof(reflect.SliceHeader{}))
const SIP_MEM_ALIGN = int32(8)
const SIP_MEM_LIGN_MASK = ^(SIP_MEM_ALIGN - 1)
const SIP_MEM_LIGN_MASK2 = (SIP_MEM_ALIGN - 1)

func RoundToAlign(x, align int32) int32 {
	return (x + align - 1) & ^(align - 1)
}

type MemAllocatorStat struct {
	allocNum      int32
	allocNumOk    int32
	freeAllNum    int32
	freePartNum   int32
	allocReqBytes int32
	allocBytes    int32
}

func (this *MemAllocatorStat) Init() {
	this.allocNum = 0
	this.allocNumOk = 0
	this.freeAllNum = 0
	this.freePartNum = 0
	this.allocReqBytes = 0

}

func (this *MemAllocatorStat) String() string {
	stat := []struct {
		name string
		num  int32
	}{
		{"alloc num", this.allocNum},
		{"alloc num ok", this.allocNumOk},
		{"free all num", this.freeAllNum},
		{"free part num", this.freePartNum},
		{"alloc request bytes", this.allocReqBytes},
	}

	str := ""
	for _, v := range stat {
		str += v.name
		str += ": "
		str += strconv.FormatUint(uint64(v.num), 10)
		str += "\n"
	}
	return str
}

type MemAllocator struct {
	mem  []byte
	used int32

	stat MemAllocatorStat
}

func NewMemAllocator(capacity int32) *MemAllocator {
	ret := MemAllocator{}
	ret.Init(capacity)
	return &ret
}

func (this *MemAllocator) Init(capacity int32) *MemAllocator {
	this.used = 0
	this.mem = make([]byte, int(capacity))
	this.stat.Init()
	return this
}

func (this *MemAllocator) Stat() *MemAllocatorStat {
	return &this.stat
}

func (this *MemAllocator) Used() int32 {
	return this.used
}

func (this *MemAllocator) ClearAllocNum() {
	this.stat.allocNum = 0
}

func (this *MemAllocator) AllocReqBytes() int32 {
	return this.stat.allocReqBytes
}

func (this *MemAllocator) AllocBytes() int32 {
	return this.stat.allocBytes
}

func (this *MemAllocator) AllocNum() int32 {
	return this.stat.allocNum
}

func (this *MemAllocator) AllocNumOk() int32 {
	return this.stat.allocNumOk
}

func (this *MemAllocator) FreeAllNum() int32 {
	return this.stat.freeAllNum
}

func (this *MemAllocator) FreePartNum() int32 {
	return this.stat.freePartNum
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

func (this *MemAllocator) Alloc(size int32) (addr AbnfPtr) {
	this.stat.allocNum++
	this.stat.allocReqBytes += size

	//logger.PrintStack()
	/*
		if print_mem {
			fmt.Printf("MEM: alloc_request = %d\n", size)
			fmt.Println(FuncNameN(2))
		} //*/
	if size <= 0 {
		return ABNF_PTR_NIL
	}

	/*if this.used&0x3 != 0 {
		fmt.Printf("mem.used = %x\n", this.used)
	}*/
	newSize := RoundToAlign(this.used+size, SIP_MEM_ALIGN)
	//newSize := (this.used + size + SIP_MEM_LIGN_MASK2) & SIP_MEM_LIGN_MASK
	//if newSize > this.Capacity() {
	if newSize > int32(cap(this.mem)) {
		//fmt.Println("newSize =", newSize)
		//fmt.Println("alloc_size =", size)
		//panic("ERROR: out of memory")
		return ABNF_PTR_NIL
	}
	used := this.used
	this.stat.allocNumOk++

	//fmt.Printf("========================== [%d] size = %d, more = %d  ==========================\n", this.stat.allocNum, size, newSize-size-used)

	//mem = (*byte)(unsafe.Pointer(&this.mem[used]))
	//addr = AbnfPtr(used)
	this.used = newSize

	/*
		if print_mem {
			fmt.Printf("MEM: used = %d\n", this.used)
		} //*/

	/*if this.used&0x3 != 0 {
		fmt.Printf("mem.used = %x\n", this.used)
	}*/

	return AbnfPtr(used)
}

func (this *MemAllocator) AllocEx(size int32) (addr AbnfPtr, alloc int32) {
	this.stat.allocNum++
	this.stat.allocReqBytes += size

	if size <= 0 {
		return ABNF_PTR_NIL, 0
	}

	used := this.used

	newSize := RoundToAlign(used+size, SIP_MEM_ALIGN)
	if newSize > int32(cap(this.mem)) {
		return ABNF_PTR_NIL, 0
	}

	//alloc = newSize - used
	this.stat.allocNumOk++
	//mem = (*byte)(unsafe.Pointer(&this.mem[used]))
	//addr = AbnfPtr(used)
	this.used = newSize

	return AbnfPtr(used), newSize - used
}

func (this *MemAllocator) FreeAll() {
	this.stat.freeAllNum++
	this.used = 0
}

func (this *MemAllocator) FreePart(remain int32) {
	this.stat.freePartNum++
	if remain >= this.used {
		return
	}
	this.used = remain
	if this.used < 0 {
		this.used = 0
	}
}

func (this *MemAllocator) String(memBegin, memEnd int) string {
	str := "-------------------------- MemAllocator show begin --------------------------\r\n"
	str += PrintByteSliceHex(memBegin, memEnd, this.mem)
	str += "-----------------------------------------------------------------------------\r\n"
	str += "MemAllocator stat:\r\n"
	str += this.stat.String()
	str += fmt.Sprintf("Used     = %d\n", this.Used())
	str += fmt.Sprintf("Left     = %d\n", this.Left())
	str += fmt.Sprintf("Capacity = %d\n", this.Capacity())
	str += "-------------------------- MemAllocator show end   --------------------------\r\n"
	return str
}

func PrintByteSliceHex(begin, end int, buf []byte) string {
	size := len(buf)
	if size == 0 {
		return ""
	}

	if begin < 0 {
		begin = 0
	}

	if end >= size {
		end = size
	}

	if begin >= end {
		return ""
	}

	size = end - begin

	lines := size / 16
	last := size % 16

	str := ""
	for i := 0; i < lines; i++ {
		str += PrintByteSliceHexOneLine(begin, buf[begin:begin+16])
		begin += 16
	}

	if last > 0 {
		str += PrintByteSliceHexOneLine(begin, buf[begin:begin+last])
	}

	return str
}

func PrintByteSliceHexOneLine(begin int, line []byte) string {
	str := fmt.Sprintf("%08xh: ", begin)
	for i := 0; i < len(line); i++ {
		str += fmt.Sprintf("%02X ", line[i])
	}
	for i := 0; i < 16-len(line); i++ {
		str += "   "
	}
	str += "; "
	for i := 0; i < len(line); i++ {
		if strconv.IsPrint(rune(line[i])) {
			if line[i] < 128 {
				str += fmt.Sprintf("%c", line[i])
			} else {
				str += "?"
			}
		} else {
			str += "."
		}
	}
	str += "\n"
	return str
}

func GetBytesDataAddr(buf *[]byte) unsafe.Pointer {
	return unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(buf)).Data)
}

//var g_allocator *MemAllocator = NewMemAllocator(1024 * 128)
