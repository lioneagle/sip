package sipparser

import (
	//"bytes"
	//"fmt"
	"testing"
)

func TestAbnfBufNew(t *testing.T) {
	testdata := []struct {
		allocatorSize uint32
		ok            bool
	}{
		{100, true},
		{SizeofAbnfBuf() + 4, true},
		{SizeofAbnfBuf() - 1, false},
		{0, false},
		{1, false},
		{7, false},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1000)
	prefix := FuncName()

	for i, v := range testdata {

		context.allocator = NewMemAllocator(v.allocatorSize)
		addr := NewAbnfBuf(context)

		if addr == ABNF_PTR_NIL && v.ok {
			t.Errorf("%s[%d]: addr should not be nil\n", prefix, i)
			continue
		}

		if addr != ABNF_PTR_NIL && !v.ok {
			t.Errorf("%s[%d]: addr should be nil\n", prefix, i)
			continue
		}

		if addr == ABNF_PTR_NIL {
			continue
		}

		buf := addr.GetAbnfBuf(context)

		if !buf.Empty() {
			t.Errorf("%s[%d]: buf should be empty (buf = %v)\n", prefix, i, buf)
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
	testdata := []struct {
		allocatorSize uint32
		buf           string
		exist         bool
		size          uint32
	}{
		{100, "123", true, 3},
		{100, "asaad", true, 5},

		{100, "", false, 0},
		{16, "123", false, 0},
		{18, "123", false, 0},
	}

	context := NewParseContext()
	prefix := FuncName()

	for i, v := range testdata {

		context.allocator = NewMemAllocator(v.allocatorSize)

		addr := NewAbnfBuf(context)
		if addr == ABNF_PTR_NIL {
			t.Errorf("%s[%d]: NewAbnfBuf failed\n", prefix, i)
			continue
		}

		buf := addr.GetAbnfBuf(context)

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
	testdata := []struct {
		allocatorSize uint32
		buf           string
		escapeNum     int
		dst           string
		exist         bool
		size          uint32
	}{
		{100, "123", 0, "123", true, 3},
		{100, "123%32", 1, "1232", true, 4},
		{100, "%31asa%33ad", 2, "1asa3ad", true, 7},

		{100, "", 0, "", false, 0},
		{16, "123", 0, "123", false, 0},
		{18, "123", 0, "123", false, 0},
	}

	context := NewParseContext()
	prefix := FuncName()

	for i, v := range testdata {

		context.allocator = NewMemAllocator(v.allocatorSize)

		addr := NewAbnfBuf(context)
		if addr == ABNF_PTR_NIL {
			t.Errorf("%s[%d]: NewAbnfBuf failed\n", prefix, i)
			continue
		}

		buf := addr.GetAbnfBuf(context)

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

func TestAbnfBufHasPrefixByteSlice(t *testing.T) {
	testdata := []struct {
		buf       string
		prefix    string
		hasPrefix bool
	}{
		{"abc-", "abc-", true},

		{"", "", false},
		{"abc-", "", false},
		{"", "abc-", false},
		{"abc-", "abc-def", false},
		{"abc-as", "aBc-", false},
		{"xa", "abc-", false},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(10000)
	prefix := FuncName()

	for i, v := range testdata {
		context.allocator.FreeAll()
		addr := NewAbnfBuf(context)
		if addr == ABNF_PTR_NIL {
			t.Errorf("%s[%d]: NewAbnfBuf failed\n", prefix, i)
			continue
		}

		buf := addr.GetAbnfBuf(context)

		buf.SetString(context, v.buf)

		hasPrefix := buf.HasPrefixByteSlice(context, []byte(v.prefix))

		if !hasPrefix && v.hasPrefix {
			t.Errorf("%s[%d]: should have prefix\n", prefix, i)
			continue
		}

		if hasPrefix && !v.hasPrefix {
			t.Errorf("%s[%d]: should not have prefix\n", prefix, i)
			continue
		}
	}
}

func TestAbnfBufHasPrefixByteSliceNoCase(t *testing.T) {
	testdata := []struct {
		buf       string
		prefix    string
		hasPrefix bool
	}{
		{"abc-", "aBC-", true},

		{"", "", false},
		{"aBc-", "", false},
		{"", "abc-", false},
		{"abc-", "aBc-def", false},
		{"xa", "Abc-", false},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(10000)
	prefix := FuncName()

	for i, v := range testdata {
		context.allocator.FreeAll()
		addr := NewAbnfBuf(context)
		if addr == ABNF_PTR_NIL {
			t.Errorf("%s[%d]: NewAbnfBuf failed\n", prefix, i)
			continue
		}

		buf := addr.GetAbnfBuf(context)

		buf.SetString(context, v.buf)

		hasPrefix := buf.HasPrefixByteSliceNoCase(context, []byte(v.prefix))

		if !hasPrefix && v.hasPrefix {
			t.Errorf("%s[%d]: should have prefix\n", prefix, i)
			continue
		}

		if hasPrefix && !v.hasPrefix {
			t.Errorf("%s[%d]: should not have prefix\n", prefix, i)
			continue
		}
	}
}

func TestAbnfBufParseEnableEmpty(t *testing.T) {
	testdata := []struct {
		isInCharset func(ch byte) bool
		src         string
		isEmpty     bool
		newPos      int
		dst         string
	}{
		{IsDigit, "", true, 0, ""},
		{IsDigit, "ad6789abc", true, 0, ""},
		{IsDigit, "01234abc", false, len("01234"), "01234"},
		{IsDigit, "5678=bc", false, len("5678"), "5678"},
	}

	prefix := FuncName()
	context := NewParseContext()
	context.allocator = NewMemAllocator(1000)

	for i, v := range testdata {
		context.allocator.FreeAll()
		addr := NewAbnfBuf(context)
		if addr == ABNF_PTR_NIL {
			t.Errorf("%s[%d]: NewAbnfBuf failed\n", prefix, i)
			continue
		}

		buf := addr.GetAbnfBuf(context)

		newPos := buf.ParseEnableEmpty(context, []byte(v.src), 0, v.isInCharset)

		if !buf.Exist() {
			t.Errorf("%s[%d]: should be exist\n", prefix, i)
			continue
		}

		if buf.Empty() && !v.isEmpty {
			t.Errorf("%s[%d]: should not be empty\n", prefix, i)
			continue
		}

		if !buf.Empty() && v.isEmpty {
			t.Errorf("%s[%d]: should be empty\n", prefix, i)
			continue
		}

		if newPos != v.newPos {
			t.Errorf("%s[%d]: newPos = %d, wanted = %d\n", prefix, i, newPos, v.newPos)
			continue
		}

		if !v.isEmpty && !buf.EqualString(context, v.dst) {
			t.Errorf("%s[%d]: wrong value = %s, wanted = %s\n", prefix, i, buf.String(context), v.dst)
			continue
		}
	}
}

func TestAbnfBufParse(t *testing.T) {
	testdata := []struct {
		isInCharset func(ch byte) bool
		src         string
		ok          bool
		newPos      int
		dst         string
	}{

		{IsDigit, "01234abc", true, len("01234"), "01234"},
		{IsDigit, "5678=bc", true, len("5678"), "5678"},

		{IsDigit, "", false, 0, ""},
		{IsDigit, "ad6789abc", false, 0, ""},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1000)
	prefix := FuncName()

	for i, v := range testdata {
		context.allocator.FreeAll()
		addr := NewAbnfBuf(context)
		if addr == ABNF_PTR_NIL {
			t.Errorf("%s[%d]: NewAbnfBuf failed\n", prefix, i)
			continue
		}

		buf := addr.GetAbnfBuf(context)

		newPos, err := buf.Parse(context, []byte(v.src), 0, v.isInCharset)

		if err != nil && v.ok {
			t.Errorf("%s[%d]: should be ok\n", prefix, i)
			continue
		}

		if err == nil && !v.ok {
			t.Errorf("%s[%d]: should not be ok\n", prefix, i)
			continue
		}

		if newPos != v.newPos {
			t.Errorf("%s[%d]: newPos = %d, wanted = %d\n", prefix, i, newPos, v.newPos)
			continue
		}

		if v.ok && buf.Empty() {
			t.Errorf("%s[%d]: should not be empty\n", prefix, i)
			continue
		}

		if v.ok && !buf.EqualString(context, v.dst) {
			t.Errorf("%s[%d]: wrong value = %s, wanted = %s\n", prefix, i, buf.String(context), v.dst)
			continue
		}
	}
}

func TestAbnfBufParseEscapableEnableEmpty(t *testing.T) {
	testdata := []struct {
		isInCharset func(ch byte) bool
		src         string
		ok          bool
		isEmpty     bool
		newPos      int
		dst         string
	}{
		{IsDigit, "", true, true, 0, ""},
		{IsDigit, "ad6789abc", true, true, 0, ""},
		{IsDigit, "01234abc", true, false, len("01234"), "01234"},
		{IsDigit, "01%3234abc", true, false, len("01%3234"), "01234"},
		{IsDigit, "56%378=bc", true, false, len("56%378"), "5678"},

		{IsDigit, "56%3x8=bc", false, false, len("56"), ""},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1000)
	prefix := FuncName()

	for i, v := range testdata {
		context.allocator.FreeAll()
		addr := NewAbnfBuf(context)
		if addr == ABNF_PTR_NIL {
			t.Errorf("%s[%d]: NewAbnfBuf failed\n", prefix, i)
			continue
		}

		buf := addr.GetAbnfBuf(context)

		newPos, err := buf.ParseEscapableEnableEmpty(context, []byte(v.src), 0, v.isInCharset)
		if err != nil && v.ok {
			t.Errorf("%s[%d]: should be ok err = %s\n", prefix, i, err.Error())
			continue
		}

		if err == nil && !v.ok {
			t.Errorf("%s[%d]: should not be ok \n", prefix, i)
			continue
		}

		if newPos != v.newPos {
			t.Errorf("%s[%d]: newPos = %d, wanted = %d\n", prefix, i, newPos, v.newPos)
			continue
		}

		if !v.ok {
			continue
		}

		if !buf.Exist() && v.ok {
			t.Errorf("%s[%d]: should be exist\n", prefix, i)
			continue
		}

		if buf.Empty() && !v.isEmpty {
			t.Errorf("%s[%d]: should not be empty\n", prefix, i)
			continue
		}

		if !buf.Empty() && v.isEmpty {
			t.Errorf("%s[%d]: should be empty\n", prefix, i)
			continue
		}

		if !v.isEmpty && !buf.EqualString(context, v.dst) {
			t.Errorf("%s[%d]: wrong value = %s, wanted = %s\n", prefix, i, buf.String(context), v.dst)
			continue
		}
	}
}

func TestAbnfBufParseEscapable(t *testing.T) {
	testdata := []struct {
		isInCharset func(ch byte) bool
		src         string
		ok          bool
		newPos      int
		dst         string
	}{

		{IsDigit, "01234abc", true, len("01234"), "01234"},
		{IsDigit, "5678=bc", true, len("5678"), "5678"},
		{IsDigit, "012%334abc", true, len("012%334"), "01234"},
		{IsDigit, "567%38=bc", true, len("567%38"), "5678"},

		{IsDigit, "", false, 0, ""},
		{IsDigit, "ad6789abc", false, 0, ""},
		{IsDigit, "012%3x4abc", false, len("012"), ""},
		{IsDigit, "01223%", false, len("01223"), ""},
		{IsDigit, "01223%4", false, len("01223"), ""},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1000)
	prefix := FuncName()

	for i, v := range testdata {
		context.allocator.FreeAll()
		addr := NewAbnfBuf(context)
		if addr == ABNF_PTR_NIL {
			t.Errorf("%s[%d]: NewAbnfBuf failed\n", prefix, i)
			continue
		}

		buf := addr.GetAbnfBuf(context)

		newPos, err := buf.ParseEscapable(context, []byte(v.src), 0, v.isInCharset)

		if err != nil && v.ok {
			t.Errorf("%s[%d]: should be ok\n", prefix, i)
			continue
		}

		if err == nil && !v.ok {
			t.Errorf("%s[%d]: should not be ok\n", prefix, i)
			continue
		}

		if newPos != v.newPos {
			t.Errorf("%s[%d]: newPos = %d, wanted = %d\n", prefix, i, newPos, v.newPos)
			continue
		}

		if v.ok && buf.Empty() {
			t.Errorf("%s[%d]: should not be empty\n", prefix, i)
			continue
		}

		if v.ok && !buf.EqualString(context, v.dst) {
			t.Errorf("%s[%d]: wrong value = %s, wanted = %s\n", prefix, i, buf.String(context), v.dst)
			continue
		}
	}
}

func TestAbnfBufSetString(t *testing.T) {
	context := NewParseContext()
	context.allocator = NewMemAllocator(1000)
	prefix := FuncName()

	addr := NewAbnfBuf(context)
	if addr == ABNF_PTR_NIL {
		t.Errorf("%s failed: should be ok\n", prefix)
		return
	}

	buf := addr.GetAbnfBuf(context)

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

	addr1 := NewAbnfBuf(context)
	addr2 := NewAbnfBuf(context)

	buf1 := addr1.GetAbnfBuf(context)
	buf2 := addr2.GetAbnfBuf(context)

	if !buf1.Equal(context, buf2) {
		t.Errorf("%s failed: should be equal", prefix)
		return
	}

	if !buf1.EqualNoCase(context, buf2) {
		t.Errorf("%s failed: should be equal", prefix)
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

	buf1.SetEmpty()
	buf2.SetEmpty()

	if !buf1.Equal(context, buf2) {
		t.Errorf("%s failed: should be equal", prefix)
		return
	}

	if !buf1.EqualNoCase(context, buf2) {
		t.Errorf("%s failed: should be equal", prefix)
		return
	}

	buf1.SetNonExist()
	buf2.SetNonExist()

	if !buf1.Equal(context, buf2) {
		t.Errorf("%s failed: should be equal", prefix)
		return
	}

	if !buf1.EqualNoCase(context, buf2) {
		t.Errorf("%s failed: should be equal", prefix)
		return
	}
}

func TestAbnfBufEncode(t *testing.T) {
	context := NewParseContext()
	context.allocator = NewMemAllocator(1000)
	prefix := FuncName()

	addr := NewAbnfBuf(context)
	if addr == ABNF_PTR_NIL {
		t.Errorf("%s failed: should be ok\n", prefix)
		return
	}

	buf := addr.GetAbnfBuf(context)

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

func BenchmarkSetByteSlice(b *testing.B) {
	b.StopTimer()
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	remain := context.allocator.Used()

	b.ReportAllocs()
	b.SetBytes(2)
	b.StartTimer()

	data := []byte("0123456789012345678901234567890123456789")
	buf := &AbnfBuf{}

	for i := 0; i < b.N; i++ {
		context.allocator.ClearAllocNum()
		context.allocator.FreePart(remain)
		buf.SetByteSlice(context, data)
	}
}

func BenchmarkSetByteSlice2(b *testing.B) {
	b.StopTimer()
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	remain := context.allocator.Used()

	b.ReportAllocs()
	b.SetBytes(2)
	b.StartTimer()

	data := []byte("0123456789012345678901234567890123456789")
	buf := &AbnfBuf{}

	for i := 0; i < b.N; i++ {
		context.allocator.ClearAllocNum()
		context.allocator.FreePart(remain)
		buf.SetByteSlice2(context, &data)
	}
}
