package sipparser

import (
	//"bytes"
	//"fmt"
	"testing"
)

func TestAbnfBufNew(t *testing.T) {
	context := NewParseContext()
	context.allocator = NewMemAllocator(1000)
	prefix := FuncName()

	testdata := []struct {
		allocatorSize int32
		ok            bool
	}{
		{100, true},
		{SizeofAbnfBuf(), true},
		{SizeofAbnfBuf() - 1, false},
		{0, false},
		{1, false},
		{7, false},
	}

	for i, v := range testdata {

		context.allocator = NewMemAllocator(v.allocatorSize)
		buf, addr := NewAbnfBuf(context)

		if buf == nil && v.ok {
			t.Errorf("%s[%d]: buf should not be nil\n", prefix, i)
			continue
		}

		if buf == nil && v.ok {
			t.Errorf("%s[%d]: buf should be nil\n", prefix, i)
			continue
		}

		if addr == ABNF_PTR_NIL && v.ok {
			t.Errorf("%s[%d]: addr should not be nil\n", prefix, i)
			continue
		}

		if addr != ABNF_PTR_NIL && !v.ok {
			t.Errorf("%s[%d]: addr should be nil\n", prefix, i)
			continue
		}

		if buf == nil {
			continue
		}

		if !buf.Empty() {
			t.Errorf("%s[%d]: buf should be empty\n", prefix, i, buf.Size())
			continue
		}

		if buf.Size() != 0 {
			t.Errorf("%s[%d]: wrong size = %d, wanted = 0\n", prefix, i, buf.Size())
			continue
		}

		if buf.Exist() {
			t.Errorf("%s[%d]: should not be exist\n", prefix, i)
			continue
		}
	}
}

func TestAbnfBufSetByteSlice(t *testing.T) {
	context := NewParseContext()

	testdata := []struct {
		allocatorSize int32
		buf           string
		exist         bool
		size          int32
	}{
		{100, "123", true, 3},
		{100, "asaad", true, 5},

		{100, "", false, 0},
		{8, "123", false, 0},
		{10, "123", false, 0},
	}

	prefix := FuncName()

	var buf *AbnfBuf

	for i, v := range testdata {

		context.allocator = NewMemAllocator(v.allocatorSize)

		buf, _ = NewAbnfBuf(context)
		if buf == nil {
			t.Errorf("%s[%d]: NewAbnfBuf failed\n", prefix, i)
			continue
		}

		buf.SetByteSlice(context, []byte(v.buf))

		if !buf.Exist() && v.exist {
			t.Errorf("%s[%d]: should be exist\n", prefix, i)
			continue
		}

		if buf.Exist() && !v.exist {
			t.Errorf("%s[%d]: should not be exist\n", prefix, i)
			continue
		}

		if buf.Size() != v.size {
			t.Errorf("%s[%d]: wrong size = %d, wanted = %d\n", prefix, i, buf.Size(), v.size)
			continue
		}

		if buf.Exist() && !buf.EqualByteSlice(context, []byte(v.buf)) {
			t.Errorf("%s[%d]: should be equal\n", prefix, i)
			return
		}
	}
}

func TestAbnfBufSetByteSliceWithUnescape(t *testing.T) {
	context := NewParseContext()

	testdata := []struct {
		allocatorSize int32
		buf           string
		escapeNum     int
		dst           string
		exist         bool
		size          int32
	}{
		{100, "123", 0, "123", true, 3},
		{100, "123%32", 1, "1232", true, 4},
		{100, "%31asa%33ad", 2, "1asa3ad", true, 7},

		{100, "", 0, "", false, 0},
		{8, "123", 0, "123", false, 0},
		{10, "123", 0, "123", false, 0},
	}

	prefix := FuncName()

	var buf *AbnfBuf

	for i, v := range testdata {

		context.allocator = NewMemAllocator(v.allocatorSize)

		buf, _ = NewAbnfBuf(context)
		if buf == nil {
			t.Errorf("%s[%d]: NewAbnfBuf failed\n", prefix, i)
			continue
		}

		buf.SetByteSliceWithUnescape(context, []byte(v.buf), v.escapeNum)

		if !buf.Exist() && v.exist {
			t.Errorf("%s[%d]: should be exist\n", prefix, i)
			continue
		}

		if buf.Exist() && !v.exist {
			t.Errorf("%s[%d]: should not be exist\n", prefix, i)
			continue
		}

		if buf.Size() != v.size {
			t.Errorf("%s[%d]: wrong size = %d, wanted = %d\n", prefix, i, buf.Size(), v.size)
			continue
		}

		if buf.Exist() && !buf.EqualByteSlice(context, []byte(v.dst)) {
			t.Errorf("%s[%d]: should be equal\n", prefix, i)
			return
		}
	}
}

func TestAbnfBufSetString(t *testing.T) {
	context := NewParseContext()
	context.allocator = NewMemAllocator(1000)
	prefix := FuncName()

	buf, _ := NewAbnfBuf(context)
	if buf == nil {
		t.Errorf("%s failed: should be ok\n", prefix)
		return
	}

	buf.SetString(context, "abc")
	if buf.Size() != 3 {
		t.Errorf("%s failed: wrong buf size= %d, wanted = 3\n", prefix, buf.Size())
		return
	}

	if !buf.EqualString(context, "abc") {
		t.Errorf("%s failed: not equal, buf = %v, wanted = abc\n", prefix, buf.String(context))
		return
	}

	if buf.GetAsString(context) != "abc" {
		t.Errorf("%s failed: not equal, buf = %v, wanted = abc\n", prefix, buf.String(context))
		return
	}
}

func TestAbnfBufEqual(t *testing.T) {
	context := NewParseContext()
	context.allocator = NewMemAllocator(1000)
	prefix := FuncName()

	buf1, _ := NewAbnfBuf(context)
	buf2, _ := NewAbnfBuf(context)

	if !buf1.Equal(context, buf2) {
		t.Errorf("%s failed: should be equal", prefix)
		return
	}

	if buf1.EqualNoCase(context, buf2) {
		t.Errorf("%s failed: should not be equal", prefix)
		return
	}

	if buf1.EqualString(context, "abc") {
		t.Errorf("%s failed: should not be equal", prefix)
		return
	}

	if buf1.EqualStringNoCase(context, "abc") {
		t.Errorf("%s failed: should not be equal", prefix)
		return
	}

	buf1.SetString(context, "abc")
	buf2.SetString(context, "abc")

	if !buf1.Equal(context, buf2) {
		t.Errorf("%s failed: not equal, buf1 = %s, buf2 = %s\n", prefix, buf1.String(context), buf2.String(context))
		return
	}

	buf1.SetString(context, "aBc")
	if !buf1.EqualNoCase(context, buf2) {
		t.Errorf("%s failed: not equal, buf1 = %s, buf2 = %s\n", prefix, buf1.String(context), buf2.String(context))
		return
	}

	if !buf1.EqualStringNoCase(context, "abc") {
		t.Errorf("%s failed: not equal, buf1 = %s, buf2 = %s\n", prefix, buf1.String(context), buf2.String(context))
		return
	}
}

func TestAbnfBufEncode(t *testing.T) {
	context := NewParseContext()
	context.allocator = NewMemAllocator(1000)
	prefix := FuncName()

	buf, _ := NewAbnfBuf(context)
	if buf == nil {
		t.Errorf("%s failed: should be ok\n", prefix)
		return
	}

	if buf.String(context) != "" {
		t.Errorf("%s failed: not empty string\n", prefix)
		return
	}

	buf.SetString(context, "456")

	if buf.String(context) != "456" {
		t.Errorf("%s failed: not equal, buf = %v, wanted = 456\n", prefix, buf.String(context))
		return
	}
}
