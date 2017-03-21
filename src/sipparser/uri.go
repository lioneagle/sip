package sipparser

import (
	"bytes"
)

type URI interface {
	Scheme() string
	Parse(src []byte, pos int) (newPos int, err error)
	String() string
	Encode(buf *bytes.Buffer)
	Equal(rhs URI) bool
}
