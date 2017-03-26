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

type SipHeader interface {
	GetName() []byte
	GetValue() *AbnfToken
	GetParsed() SipHeaderParsed
}

type SipHeaderList struct {
	headers []*SipHeader
}
