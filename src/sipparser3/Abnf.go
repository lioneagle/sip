package sipparser3

import (
	"bytes"
	"fmt"
)

type AbnfIsInCharset func(ch byte) bool

type AbnfEncoder interface {
	Encode(buf *bytes.Buffer)
}

func AbnfEncoderToString(encoder AbnfEncoder) string {
	var buf bytes.Buffer
	encoder.Encode(&buf)
	return buf.String()
}

type AbnfError struct {
	description string
	src         []byte
	pos         int
}

func (err *AbnfError) Error() string {
	if err.pos < len(err.src) {
		return fmt.Sprintf("%s at src[%d]: %s", err.description, err.pos, string(err.src[err.pos:]))
	}
	return fmt.Sprintf("%s at end", err.description)
}
