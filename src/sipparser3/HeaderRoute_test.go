package sipparser3

import (
	"bytes"
	"fmt"
	"testing"
)

func TestSipHeaderRouteParse(t *testing.T) {

	testdata := []struct {
		src    string
		ok     bool
		newPos int
		encode string
	}{
		{"Route: <sip:abc@a.com>;tag=1", true, len("Route: <sip:abc@a.com>;tag=1"), "Route: <sip:abc@a.com>;tag=1"},
		{"Route: <sip:abc@a.com;user=ip>;tag=1", true, len("Route: <sip:abc@a.com;user=ip>;tag=1"), "Route: <sip:abc@a.com;user=ip>;tag=1"},
		{"Route: abc<sip:abc@a.com;user=ip>;tag=1", true, len("Route: abc<sip:abc@a.com;user=ip>;tag=1"), "Route: abc<sip:abc@a.com;user=ip>;tag=1"},
		{"Route: <tel:+12358;tag=123>", true, len("Route: <tel:+12358;tag=123>"), "Route: <tel:+12358;tag=123>"},

		{" Route: <sip:abc@a.com>;tag=1", false, 0, "0"},
		{"Route1: <sip:abc@a.com>;tag=1", false, len("Route1: "), ""},
		{"Route: ", false, len("Route: "), ""},
	}

	context := NewParseContext()

	for i, v := range testdata {
		header := NewSipHeaderRoute()
		newPos, err := header.Parse(context, []byte(v.src), 0)

		if v.ok && err != nil {
			t.Errorf("TestSipHeaderRouteParse[%d] failed, err = %s\n", i, err)
			continue
		}

		if !v.ok && err == nil {
			t.Errorf("TestSipHeaderRouteParse[%d] failed, should parse failed", i)
			continue
		}

		if v.newPos != newPos {
			t.Errorf("TestSipHeaderRouteParse[%d] failed, newPos = %d, wanted = %d\n", i, newPos, v.newPos)
		}

		if v.ok && v.encode != header.String() {
			t.Errorf("TestSipHeaderRouteParse[%d] failed, encode = %s, wanted = %s\n", i, header.String(), v.encode)
			continue
		}
	}

}

func BenchmarkSipHeaderRouteParse(b *testing.B) {
	b.StopTimer()
	//v := []byte("sip:biloxi.com;transport=tcp;method=REGISTER?to=sip:bob%40biloxi.com")
	//v := []byte("sip:biloxi.com")
	//v := []byte("sip:abc@biloxi.com;transport=tcp")
	v := []byte("Route: <sip:abc@biloxi.com;transport=tcp;method=REGISTER>")
	context := NewParseContext()
	header := NewSipHeaderRoute()

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

func BenchmarkSipHeaderRouteEncode(b *testing.B) {
	b.StopTimer()
	//v := []byte("sip:biloxi.com;transport=tcp;method=REGISTER?to=sip:bob%40biloxi.com")
	v := []byte("Route: <sip:abc@biloxi.com;transport=tcp;method=REGISTER>")
	context := NewParseContext()
	header := NewSipHeaderRoute()
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
