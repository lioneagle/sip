package sipparser

type ParseContext struct {
	EncodeHeaderShorName bool
	allocator            *MemAllocator
	parseSrc             []byte
	parsePos             uint32
}

func NewParseContext() *ParseContext {
	return &ParseContext{}
}

func (this *ParseContext) SetAllocator(allocator *MemAllocator) {
	this.allocator = allocator
}

func (this *ParseContext) ClearAllocNum() {
	this.allocator.ClearAllocNum()
}

func (this *ParseContext) FreePart(remain int32) {
	this.allocator.FreePart(remain)
}

func (this *ParseContext) Used() int32 {
	return this.allocator.Used()
}
