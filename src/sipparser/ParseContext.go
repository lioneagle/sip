package sipparser

type ParseContext struct {
	EncodeHeaderShorName bool
	allocator            *MemAllocator
}

func NewParseContext() *ParseContext {
	return &ParseContext{}
}
