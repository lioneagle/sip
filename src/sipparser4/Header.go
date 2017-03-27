package sipparser4

import (
	"bytes"
)

type SipHeaderParsed interface {
	IsHeaderList() bool
	Parse(src []byte, pos int) (newPos int, err error)
	String() string
	Encode(buf *bytes.Buffer)
}

type SipHeader struct {
}

type SipHeaderList struct {
	headers []SipHeader
}
