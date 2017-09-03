package sipparser

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
	context.allocator = NewMemAllocator(1024 * 30)
	prefix := FuncName()

	for i, v := range testdata {
		header, _ := NewSipHeaderRoute(context)
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

func BenchmarkSipHeaderRouteParse(b *testing.B) {
	b.StopTimer()
	//v := []byte("sip:biloxi.com;transport=tcp;method=REGISTER?to=sip:bob%40biloxi.com")
	//v := []byte("sip:biloxi.com")
	//v := []byte("sip:abc@biloxi.com;transport=tcp")
	v := []byte("Route: <sip:abc@biloxi.com;transport=tcp;method=REGISTER>")
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	header, _ := NewSipHeaderRoute(context)
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

func BenchmarkSipHeaderRouteEncode(b *testing.B) {
	b.StopTimer()
	//v := []byte("sip:biloxi.com;transport=tcp;method=REGISTER?to=sip:bob%40biloxi.com")
	v := []byte("Route: <sip:abc@biloxi.com;transport=tcp;method=REGISTER>")
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	header, _ := NewSipHeaderRoute(context)
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
