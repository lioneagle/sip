package sipparser

import (
	"bytes"
	"fmt"
	"testing"
)

func TestSipHeaderContentTypeParse(t *testing.T) {

	testdata := []struct {
		src    string
		ok     bool
		newPos int
		encode string
	}{
		{"Content-Type: application/sdp", true, len("Content-Type: application/sdp"), "Content-Type: application/sdp"},
		{"c: message/sip", true, len("c: message/sip"), "Content-Type: message/sip"},

		{" Content-Type: abc", false, 0, ""},
		{"Content-Type2: abc", false, len("Content-Type2: "), ""},
		{"Content-Type: ", false, len("Content-Type: "), ""},
		{"Content-Type: abc", false, len("Content-Type: abc"), ""},
		{"Content-Type: abc/", false, len("Content-Type: abc/"), ""},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)

	for i, v := range testdata {
		header, _ := NewSipHeaderContentType(context)
		newPos, err := header.Parse(context, []byte(v.src), 0)

		if v.ok && err != nil {
			t.Errorf("TestSipHeaderContentTypeParse[%d] failed, err = %s\n", i, err)
			continue
		}

		if !v.ok && err == nil {
			t.Errorf("TestSipHeaderContentTypeParse[%d] failed, should parse failed", i)
			continue
		}

		if v.newPos != newPos {
			t.Errorf("TestSipHeaderContentTypeParse[%d] failed, newPos = %d, wanted = %d\n", i, newPos, v.newPos)
		}

		if v.ok && v.encode != header.String(context) {
			t.Errorf("TestSipHeaderContentTypeParse[%d] failed, encode = %s, wanted = %s\n", i, header.String(context), v.encode)
			continue
		}
	}

}

func BenchmarkSipHeaderContentTypeParse(b *testing.B) {
	b.StopTimer()
	v := []byte("Content-Type: application/sdp")
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	header, _ := NewSipHeaderContentType(context)
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

func BenchmarkSipHeaderContentTypeEncode(b *testing.B) {
	b.StopTimer()
	v := []byte("Content-Type: application/sdp")
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	header, _ := NewSipHeaderContentType(context)
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
