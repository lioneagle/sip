package sipparser

import (
	"bytes"
	"fmt"
	"testing"
)

func TestSipHeaderContentLengthParse(t *testing.T) {

	testdata := []struct {
		src       string
		ok        bool
		newPos    int
		encode    string
		encodeEnd int
	}{
		{"Content-Length: 1234", true, len("Content-Length: 1234"), "Content-Length:       1234", len("Content-Length:       1234")},
		{"l: 1234", true, len("l: 1234"), "Content-Length:       1234", len("Content-Length:       1234")},

		{" Content-Lengt: 1234", false, 0, "", 0},
		{"Content-Lengt: 1234", false, len("Content-Lengt: "), "", 0},
		{"Content-Length: ", false, len("Content-Length: "), "", 0},
		{"Content-Length: a123", false, len("Content-Length: "), "", 0},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	prefix := FuncName()

	for i, v := range testdata {
		header, _ := NewSipHeaderContentLength(context)
		newPos, err := header.Parse(context, []byte(v.src), 0)

		if v.ok && err != nil {
			t.Errorf("%s[%d] failed: err = %s\n", prefix, i, err)
			continue
		}

		if !v.ok && err == nil {
			t.Errorf("%s[%d] failed: should parse failed", prefix, i)
			continue
		}

		if v.newPos != newPos {
			t.Errorf("%s[%d] failed: newPos = %d, wanted = %d\n", prefix, i, newPos, v.newPos)
		}

		if !v.ok {
			continue
		}

		if v.encode != header.String(context) {
			t.Errorf("%s[%d] failed: encode = %s, wanted = %s\n", prefix, i, header.String(context), v.encode)
			continue
		}

		if v.encodeEnd != int(header.encodeEnd) {
			t.Errorf("%s[%d] failed: encodeEnd = %s, wanted = %s\n", prefix, i, header.encodeEnd, v.encodeEnd)
			continue
		}
	}

}

func BenchmarkSipHeaderContentLengthParse(b *testing.B) {
	b.StopTimer()
	v := []byte("Content-Length: 226")
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	header, _ := NewSipHeaderContentLength(context)
	remain := context.allocator.Used()
	b.ReportAllocs()
	b.SetBytes(2)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		context.allocator.ClearAllocNum()
		context.allocator.FreePart(remain)
		header.Parse(context, v, 0)
	}
	//fmt.Printf("header = %s\n", header.String())
	fmt.Printf("")
}

func BenchmarkSipHeaderContentLengthEncode(b *testing.B) {
	b.StopTimer()
	v := []byte("Content-Length: 226")
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	header, _ := NewSipHeaderContentLength(context)
	header.Parse(context, v, 0)
	remain := context.allocator.Used()
	buf := bytes.NewBuffer(make([]byte, 1024*1024))
	b.SetBytes(2)
	b.ReportAllocs()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		buf.Reset()
		context.allocator.ClearAllocNum()
		context.allocator.FreePart(remain)
		header.Encode(context, buf)
	}
}
