package sipparser

import (
	_ "bytes"
)

type URI interface {
	Init()
	Scheme() string
	Parse(context *ParseContext, src []byte, pos int) (newPos int, err error)
	String(context *ParseContext) string
	Encode(context *ParseContext, buf *AbnfByteBuffer)
	Equal(context *ParseContext, rhs URI) bool
}
