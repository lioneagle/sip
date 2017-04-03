package sipparser

import (
	"bytes"
)

type URI interface {
	Init()
	Scheme() string
	Parse(context *ParseContext, src []byte, pos int) (newPos int, err error)
	String(context *ParseContext) string
	Encode(context *ParseContext, buf *bytes.Buffer)
	Equal(context *ParseContext, rhs URI) bool
}
