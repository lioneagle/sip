package sipparser3

import (
	"bytes"
	"fmt"
	"testing"
)

func TestSipHeaderViaParse(t *testing.T) {

	testdata := []struct {
		src    string
		ok     bool
		newPos int
		encode string
	}{
		{"Via: SIP/2.0/UDP 10.4.1.1:5070;branch=123", true, len("Via: SIP/2.0/UDP 10.4.1.1:5070;branch=123"), "Via: SIP/2.0/UDP 10.4.1.1:5070;branch=123"},
	}

	context := NewParseContext()

	for i, v := range testdata {
		header := NewSipHeaderVia()
		newPos, err := header.Parse(context, []byte(v.src), 0)

		if v.ok && err != nil {
			t.Errorf("TestSipHeaderViaParse[%d] failed, err = %s\n", i, err)
			continue
		}

		if !v.ok && err == nil {
			t.Errorf("TestSipHeaderViaParse[%d] failed, should parse failed", i)
			continue
		}

		if v.newPos != newPos {
			t.Errorf("TestSipHeaderViaParse[%d] failed, newPos = %d, wanted = %d\n", i, newPos, v.newPos)
		}

		if v.ok && v.encode != header.String() {
			t.Errorf("TestSipHeaderViaParse[%d] failed, encode = %s, wanted = %s\n", i, header.String(), v.encode)
			continue
		}
	}

}

func BenchmarkSipHeaderViaParse(b *testing.B) {
	b.StopTimer()
	v := []byte("Via: SIP/2.0/UDP 10.4.1.1:5070;branch=123")
	context := NewParseContext()
	header := NewSipHeaderVia()

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

func BenchmarkSipHeaderViaEncode(b *testing.B) {
	b.StopTimer()
	v := []byte("Via: SIP/2.0/UDP 10.4.1.1:5070;branch=123")
	context := NewParseContext()
	header := NewSipHeaderVia()
	header.Parse(context, v, 0)
	b.SetBytes(2)
	b.ReportAllocs()

	buf := bytes.NewBuffer(make([]byte, 1024*1024))

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		buf.Reset()
		header.Encode(buf)
	}
}
