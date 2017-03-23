package sipparser2

type MemAllocator struct {
	mem  []byte
	used int
}

func (this *MemAllocator) Init(capacity int) *MemAllocator {
	this.used = 0
	this.mem = make([]byte, capacity)
	return this
}

func (this *MemAllocator) Alloc(size int) []byte {
	if (this.used + size) >= len(this.mem) {
		panic("ERROR: out of memory")
	}
	old := this.used
	this.used += size
	return this.mem[old : this.used+size]
}

func (this *MemAllocator) FreeAll() {
	this.used = 0
}

var g_allocator *MemAllocator = new(MemAllocator).Init(1024 * 128)
