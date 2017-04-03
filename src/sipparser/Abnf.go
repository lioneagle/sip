package sipparser

import (
	"bytes"
	"fmt"
	"reflect"
	"unsafe"
)

var print_mem bool = false

const (
	ABNF_UNKNOWN_URI   = int32(0)
	ABNF_SIP_URI       = int32(1)
	ABNF_SIPS_URI      = int32(2)
	ABNF_TEL_URI       = int32(3)
	ABNF_ABSOULUTE_URI = int32(4)
)

const (
	ABNF_SIP_ADDR_SPEC = int32(0)
	ABNF_SIP_NAME_ADDR = int32(1)
)

type AbnfIsInCharset func(ch byte) bool

type AbnfEncoder interface {
	Encode(context *ParseContext, buf *bytes.Buffer)
}

func AbnfEncoderToString(context *ParseContext, encoder AbnfEncoder) string {
	var buf bytes.Buffer
	encoder.Encode(context, &buf)
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

func StringToByteSlice(str string) []byte {
	strHeader := (*reflect.StringHeader)(unsafe.Pointer(&str))
	retHeader := reflect.SliceHeader{Data: strHeader.Data, Len: strHeader.Len, Cap: strHeader.Len}
	return *(*[]byte)(unsafe.Pointer(&retHeader))
}

func ByteSliceToString(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}
