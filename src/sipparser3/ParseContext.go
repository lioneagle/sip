package sipparser3

type ParseContext struct {
	EncodeHeaderShorName bool
}

func NewParseContext() *ParseContext {
	return &ParseContext{}
}
