package sipparser

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

		{" Via: SIP/2.0/UDP 10.4.1.1:5070;branch=123", false, 0, ""},
		{"Via2: SIP/2.0/UDP 10.4.1.1:5070;branch=123", false, len("Via2: "), ""},
		{"Via: SIP/2/UDP 10.4.1.1:5070;branch=123", false, len("Via: SIP/2"), ""},
		{"Via: SIP/2.0UDP 10.4.1.1:5070;branch=123", false, len("Via: SIP/2.0"), ""},
		{"Via: SIP/2.0/@ 10.4.1.1:5070;branch=123", false, len("Via: SIP/2.0/"), ""},
		{"Via: SIP/2.0/UDP\r\n10.4.1.1:5070;branch=123", false, len("Via: SIP/2.0/UDP\r\n"), ""},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)

	for i, v := range testdata {
		header, _ := NewSipHeaderVia(context)
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

		if v.ok && v.encode != header.String(context) {
			t.Errorf("TestSipHeaderViaParse[%d] failed, encode = %s, wanted = %s\n", i, header.String(context), v.encode)
			continue
		}
	}

}

func BenchmarkSipHeaderViaParse(b *testing.B) {
	b.StopTimer()
	//v := []byte("Via: SIP/2.0/UDP 10.4.1.1:5070;branch=123")
	v := []byte("Via: SIP/2.0/UDP 24.15.255.101:5060")
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	header, _ := NewSipHeaderVia(context)
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

func BenchmarkSipHeaderViaEncode(b *testing.B) {
	b.StopTimer()
	//v := []byte("Via: SIP/2.0/UDP 10.4.1.1:5070;branch=123")
	v := []byte("Via: SIP/2.0/UDP 24.15.255.101:5060")
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	header, _ := NewSipHeaderVia(context)
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
