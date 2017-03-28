package sipparser3

import (
	"bytes"
)

type URI interface {
	Init()
	Scheme() string
	Parse(context *ParseContext, src []byte, pos int) (newPos int, err error)
	String() string
	Encode(buf *bytes.Buffer)
	Equal(rhs URI) bool
}
