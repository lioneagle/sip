package sipparser

import (
	//"bytes"
	"fmt"
	"testing"
)

func TestSipHeaderContentDispositionParse(t *testing.T) {

	testdata := []struct {
		src    string
		ok     bool
		newPos int
		encode string
	}{
		{"Content-Disposition: session", true, len("Content-Disposition: session"), "Content-Disposition: session"},
		{"Content-Disposition: early-session", true, len("Content-Disposition: early-session"), "Content-Disposition: early-session"},

		{" Content-Disposition: abc", false, 0, ""},
		{"Content-Disposition2: abc", false, len("Content-Disposition2: "), ""},
		{"Content-Disposition: ", false, len("Content-Disposition: "), ""},
		{"Content-Disposition: @", false, len("Content-Disposition: "), ""},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	prefix := FuncName()

	for i, v := range testdata {
		addr := NewSipHeaderContentDisposition(context)
		header := addr.GetSipHeaderContentDisposition(context)
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

		if v.ok && v.encode != header.String(context) {
			t.Errorf("%s[%d] failed: encode = %s, wanted = %s\n", prefix, i, header.String(context), v.encode)
			continue
		}
	}

}

func BenchmarkSipHeaderContentDispositionParse(b *testing.B) {
	b.StopTimer()
	v := []byte("Content-Disposition: early-session;handling=optional")
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	addr := NewSipHeaderContentDisposition(context)
	header := addr.GetSipHeaderContentDisposition(context)
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

func BenchmarkSipHeaderContentDispositionEncode(b *testing.B) {
	b.StopTimer()
	v := []byte("Content-Disposition: early-session;handling=optional")
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	addr := NewSipHeaderContentDisposition(context)
	header := addr.GetSipHeaderContentDisposition(context)
	header.Parse(context, v, 0)
	remain := context.allocator.Used()
	//buf := bytes.NewBuffer(make([]byte, 1024*1024))
	buf := &AbnfByteBuffer{}
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
