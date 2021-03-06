package sipparser

import (
	//"bytes"
	"fmt"
	"testing"
)

func TestSipHeaderCseqParse(t *testing.T) {

	testdata := []struct {
		src    string
		ok     bool
		newPos int
		encode string
	}{
		{"CSeq: 1234 INVITE", true, len("CSeq: 1234 INVITE"), "CSeq: 1234 INVITE"},

		{" CSeq: ", false, 0, ""},
		{"CSeq2: ", false, len("CSeq2: "), ""},
		{"CSeq: ", false, len("CSeq: "), ""},
		{"CSeq: 1234", false, len("CSeq: 1234"), ""},
		{"CSeq: 1234 \r\nINVITE", false, len("CSeq: 1234 \r\n"), ""},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	prefix := FuncName()

	for i, v := range testdata {
		addr := NewSipHeaderCseq(context)
		header := addr.GetSipHeaderCseq(context)
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

func BenchmarkSipHeaderCseqParse(b *testing.B) {
	b.StopTimer()
	v := []byte("CSeq: 101 INVITE")
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	addr := NewSipHeaderCseq(context)
	header := addr.GetSipHeaderCseq(context)
	remain := context.allocator.Used()
	b.ReportAllocs()
	b.SetBytes(2)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		context.allocator.ClearAllocNum()
		context.allocator.FreePart(remain)
		header.Parse(context, v, 0)
	}
	//fmt.Printf("uri = %s\n", uri.String())
	fmt.Printf("")
}

func BenchmarkSipHeaderCseqEncode(b *testing.B) {
	b.StopTimer()
	v := []byte("CSeq: 101 INVITE")
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	addr := NewSipHeaderCseq(context)
	header := addr.GetSipHeaderCseq(context)
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
