package sipparser3

type ParseContext struct {
	FuncName string
}

func NewParseContext() *ParseContext {
	return &ParseContext{}
}
