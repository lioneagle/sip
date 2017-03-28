package sipparser3

import (
	"bytes"
	"fmt"
	"testing"
)

func TestSipHeaderCallIdParse(t *testing.T) {

	testdata := []struct {
		src    string
		ok     bool
		newPos int
		encode string
	}{
		{"Call-ID: abc123@a.com", true, len("Call-ID: abc123@a.com"), "Call-ID: abc123@a.com"},
		{"Call-ID: abc123", true, len("Call-ID: abc123"), "Call-ID: abc123"},
		{"Call-ID: abc123\r\n", true, len("Call-ID: abc123"), "Call-ID: abc123"},

		{" Call-ID: abc123@", false, 0, ""},
		{"Call-ID1: abc123@", false, len("Call-ID1: "), ""},
		{"Call-ID: abc123@", false, len("Call-ID: abc123@"), ""},
		{"Call-ID: @abc", false, len("Call-ID: "), ""},
	}

	context := NewParseContext()

	for i, v := range testdata {
		header := NewSipHeaderCallId()
		newPos, err := header.Parse(context, []byte(v.src), 0)

		if v.ok && err != nil {
			t.Errorf("TestSipHeaderCallIdParse[%d] failed, err = %s\n", i, err)
			continue
		}

		if !v.ok && err == nil {
			t.Errorf("TestSipHeaderCallIdParse[%d] failed, should parse failed", i)
			continue
		}

		if v.newPos != newPos {
			t.Errorf("TestSipHeaderCallIdParse[%d] failed, newPos = %d, wanted = %d\n", i, newPos, v.newPos)
		}

		if v.ok && v.encode != header.String() {
			t.Errorf("TestSipHeaderCallIdParse[%d] failed, encode = %s, wanted = %s\n", i, header.String(), v.encode)
			continue
		}
	}

}

func BenchmarkSipHeaderCallIdParse(b *testing.B) {
	b.StopTimer()
	v := []byte("Call-ID: abc123@a.com")
	context := NewParseContext()
	header := NewSipHeaderCallId()

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

func BenchmarkSipHeaderCallIdEncode(b *testing.B) {
	b.StopTimer()
	v := []byte("Call-ID: abc123@a.com")
	context := NewParseContext()
	header := NewSipHeaderCallId()
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
