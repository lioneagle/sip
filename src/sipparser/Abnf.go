package sipparser

import (
	//"bytes"
	"fmt"
	"reflect"
	"runtime"
	"unsafe"
)

//type AbnfByteBuffer = bytes.Buffer
type AbnfByteBuffer = ByteBuffer

var print_mem bool = false

const (

	// single characters
	ABNF_NAME_COLON = ":"
	ABNF_NAME_SPACE = " "
)

type AbnfIsInCharset func(ch byte) bool

type AbnfEncoder interface {
	Encode(context *ParseContext, buf *AbnfByteBuffer)
}

func AbnfEncoderToString(context *ParseContext, encoder AbnfEncoder) string {
	//var buf AbnfByteBuffer
	buf := NewAbnfByteBuffer(make([]byte, 64))
	encoder.Encode(context, buf)
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
	return err.description
}

func StringToByteSlice(str string) []byte {
	strHeader := (*reflect.StringHeader)(unsafe.Pointer(&str))
	retHeader := reflect.SliceHeader{Data: strHeader.Data, Len: strHeader.Len, Cap: strHeader.Len}
	return *(*[]byte)(unsafe.Pointer(&retHeader))
}

func StringToByteSlice2(str string) *[]byte {
	strHeader := (*reflect.StringHeader)(unsafe.Pointer(&str))
	retHeader := reflect.SliceHeader{Data: strHeader.Data, Len: strHeader.Len, Cap: strHeader.Len}
	return (*[]byte)(unsafe.Pointer(&retHeader))
}

func ByteSliceToString(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}

func FuncName() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	return runtime.FuncForPC(pc).Name()
}

func FuncNameN(n int) string {
	pc, _, _, ok := runtime.Caller(n)
	if !ok {
		return ""
	}
	return runtime.FuncForPC(pc).Name()
}
