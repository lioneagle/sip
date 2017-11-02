package sipparser

type ParseContext struct {
	EncodeHeaderShorName bool
	allocator            *MemAllocator
	parseSrc             []byte
	parsePos             uint32
	ParseSipHeaderAsRaw  bool
	//ParseSipKeyHeader    bool
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

func (this *ParseContext) GetAllocNum() uint32 {
	return this.allocator.AllocNum()
}

func (this *ParseContext) FreePart(remain uint32) {
	this.allocator.FreePart(remain)
}

func (this *ParseContext) Used() uint32 {
	return this.allocator.Used()
}
