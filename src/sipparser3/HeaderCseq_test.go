package sipparser3

import (
	"bytes"
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

	for i, v := range testdata {
		header := NewSipHeaderCseq()
		newPos, err := header.Parse(context, []byte(v.src), 0)

		if v.ok && err != nil {
			t.Errorf("TestSipHeaderCseqParse[%d] failed, err = %s\n", i, err)
			continue
		}

		if !v.ok && err == nil {
			t.Errorf("TestSipHeaderCseqParse[%d] failed, should parse failed", i)
			continue
		}

		if v.newPos != newPos {
			t.Errorf("TestSipHeaderCseqParse[%d] failed, newPos = %d, wanted = %d\n", i, newPos, v.newPos)
		}

		if v.ok && v.encode != header.String() {
			t.Errorf("TestSipHeaderCseqParse[%d] failed, encode = %s, wanted = %s\n", i, header.String(), v.encode)
			continue
		}
	}

}

func BenchmarkSipHeaderCseqParse(b *testing.B) {
	b.StopTimer()
	v := []byte("CSeq: 1234 INVITE")
	context := NewParseContext()
	header := NewSipHeaderCseq()

	b.ReportAllocs()
	b.SetBytes(2)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		header.Init()
		header.Parse(context, v, 0)
	}
	//fmt.Printf("uri = %s\n", uri.String())
	fmt.Printf("")
}

func BenchmarkSipHeaderCseqEncode(b *testing.B) {
	b.StopTimer()
	v := []byte("CSeq: 1234 INVITE")
	context := NewParseContext()
	header := NewSipHeaderCseq()
	header.Parse(context, v, 0)
	b.SetBytes(2)
	b.ReportAllocs()

	buf := bytes.NewBuffer(make([]byte, 1024*1024))
	//buf := &bytes.Buffer{}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		buf.Reset()
		header.Encode(buf)
	}
}
