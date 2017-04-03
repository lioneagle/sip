package sipparser

import (
	//"bytes"
	//"fmt"
	"testing"
)

func TestAbnfBufNew(t *testing.T) {
	context := NewParseContext()
	context.allocator = NewMemAllocator(1000)
	prefix := "TestAbnfBufNew"

	buf, addr := NewAbnfBuf(context)
	if buf == nil {
		t.Errorf("%s failed, should be ok\n", prefix)
		return
	}

	if addr != 0 {
		t.Errorf("%s failed, wrong addr = %d, wanted = 0\n", prefix, addr)
		return
	}

	if buf.addr != ABNF_PTR_NIL {
		t.Errorf("%s failed, wrong addr = %d, wanted = %s\n", prefix, buf.addr, ABNF_PTR_NIL.String())
		return
	}

	if buf.Size() != 0 {
		t.Errorf("%s failed, wrong size = %d, wanted = 0\n", prefix, buf.Size())
		return
	}

	context.allocator = NewMemAllocator(1)
	buf, _ = NewAbnfBuf(context)
	if buf != nil {
		t.Errorf("%s failed, should not be ok\n", prefix)
		return
	}

}

func TestAbnfBufSetByteSlice(t *testing.T) {
	context := NewParseContext()
	context.allocator = NewMemAllocator(1000)
	prefix := "TestAbnfBufSetByteSlice"

	buf, _ := NewAbnfBuf(context)
	if buf == nil {
		t.Errorf("%s failed, should be ok\n", prefix)
		return
	}

	buf.SetByteSlice(context, []byte(""))
	if buf.size != 0 {
		t.Errorf("%s failed, wrong buf size= %d, wanted = 0\n", prefix, buf.size)
		return
	}

	bytes1 := []byte("123")
	buf.SetByteSlice(context, bytes1)
	if buf.size != 3 {
		t.Errorf("%s failed, wrong buf size= %d, wanted = 3\n", prefix, buf.size)
		return
	}

	if !buf.EqualByteSlice(context, bytes1) {
		t.Errorf("%s failed, not equal, buf = %s, wanted = %s\n", prefix, buf.String(context), string(bytes1))
		return
	}
}

func TestAbnfBufSetString(t *testing.T) {
	context := NewParseContext()
	context.allocator = NewMemAllocator(1000)
	prefix := "TestAbnfBufSetString"

	buf, _ := NewAbnfBuf(context)
	if buf == nil {
		t.Errorf("%s failed, should be ok\n", prefix)
		return
	}

	buf.SetString(context, "abc")
	if buf.Size() != 3 {
		t.Errorf("%s failed, wrong buf size= %d, wanted = 3\n", prefix, buf.Size())
		return
	}

	if !buf.EqualString(context, "abc") {
		t.Errorf("%s failed, not equal, buf = %v, wanted = abc\n", prefix, buf.String(context))
		return
	}

	if buf.GetAsString(context) != "abc" {
		t.Errorf("%s failed, not equal, buf = %v, wanted = abc\n", prefix, buf.String(context))
		return
	}
}

func TestAbnfBufEqual(t *testing.T) {
	context := NewParseContext()
	context.allocator = NewMemAllocator(1000)
	prefix := "TestAbnfBufEqual"

	buf1, _ := NewAbnfBuf(context)
	buf2, _ := NewAbnfBuf(context)

	if !buf1.Equal(context, buf2) {
		t.Errorf("%s failed, should be equal", prefix)
		return
	}

	if buf1.EqualNoCase(context, buf2) {
		t.Errorf("%s failed, should not be equal", prefix)
		return
	}

	if buf1.EqualString(context, "abc") {
		t.Errorf("%s failed, should not be equal", prefix)
		return
	}

	if buf1.EqualStringNoCase(context, "abc") {
		t.Errorf("%s failed, should not be equal", prefix)
		return
	}

	buf1.SetString(context, "abc")
	buf2.SetString(context, "abc")

	if !buf1.Equal(context, buf2) {
		t.Errorf("%s failed, not equal, buf1 = %s, buf2 = %s\n", prefix, buf1.String(context), buf2.String(context))
		return
	}

	buf1.SetString(context, "aBc")
	if !buf1.EqualNoCase(context, buf2) {
		t.Errorf("%s failed, not equal, buf1 = %s, buf2 = %s\n", prefix, buf1.String(context), buf2.String(context))
		return
	}

	if !buf1.EqualStringNoCase(context, "abc") {
		t.Errorf("%s failed, not equal, buf1 = %s, buf2 = %s\n", prefix, buf1.String(context), buf2.String(context))
		return
	}
}

func TestAbnfBufEncode(t *testing.T) {
	context := NewParseContext()
	context.allocator = NewMemAllocator(1000)
	prefix := "TestAbnfBufEncode"

	buf, _ := NewAbnfBuf(context)
	if buf == nil {
		t.Errorf("% failed, should be ok\n", prefix)
		return
	}

	if buf.String(context) != "" {
		t.Errorf("% failed, not empty string\n", prefix)
		return
	}

	buf.SetString(context, "456")

	if buf.String(context) != "456" {
		t.Errorf("% failed, not equal, buf = %v, wanted = 456\n", prefix, buf.String(context))
		return
	}
}
