package sipparser

import (
	"bytes"
	"fmt"
	"testing"
)

func TestSipHeaderFromParse(t *testing.T) {

	testdata := []struct {
		src    string
		ok     bool
		newPos int
		encode string
	}{
		{"From: sip:abc@a.com;tag=1", true, len("From: sip:abc@a.com;tag=1"), "From: <sip:abc@a.com>;tag=1"},
		{"From: <sip:abc@a.com;user=ip>;tag=1", true, len("From: <sip:abc@a.com;user=ip>;tag=1"), "From: <sip:abc@a.com;user=ip>;tag=1"},
		{"From: abc<sip:abc@a.com;user=ip>;tag=1", true, len("From: abc<sip:abc@a.com;user=ip>;tag=1"), "From: abc<sip:abc@a.com;user=ip>;tag=1"},
		{"From: tel:+12358;tag=123", true, len("From: tel:+12358;tag=123"), "From: <tel:+12358>;tag=123"},

		{" From: <sip:abc@a.com>;tag=1", false, 0, "0"},
		{"From1: <sip:abc@a.com>;tag=1", false, len("From1: "), ""},
		{"From: ", false, len("From: "), ""},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	prefix := FuncName()

	for i, v := range testdata {
		addr := NewSipHeaderFrom(context)
		header := addr.GetSipHeaderFrom(context)
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

func BenchmarkSipHeaderFromParse(b *testing.B) {
	b.StopTimer()
	//v := []byte("sip:biloxi.com;transport=tcp;method=REGISTER?to=sip:bob%40biloxi.com")
	//v := []byte("sip:biloxi.com")
	//v := []byte("sip:abc@biloxi.com;transport=tcp")
	//v := []byte("From: <sip:abc@biloxi.com;transport=tcp;method=REGISTER>")
	v := []byte("From: \"User ID\" <sip:6140000@24.15.255.4>;tag=dab70900252036d7134be-4ec05abe")
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	addr := NewSipHeaderFrom(context)
	header := addr.GetSipHeaderFrom(context)
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

func BenchmarkSipHeaderFromEncode(b *testing.B) {
	b.StopTimer()
	//v := []byte("sip:biloxi.com;transport=tcp;method=REGISTER?to=sip:bob%40biloxi.com")
	//v := []byte("From: <sip:abc@biloxi.com;transport=tcp;method=REGISTER>")
	v := []byte("From: \"User ID\" <sip:6140000@24.15.255.4>;tag=dab70900252036d7134be-4ec05abe")
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	addr := NewSipHeaderFrom(context)
	header := addr.GetSipHeaderFrom(context)
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

//*
func BenchmarkSipHeaderFromString(b *testing.B) {
	b.StopTimer()
	//v := []byte("sip:biloxi.com;transport=tcp;method=REGISTER?to=sip:bob%40biloxi.com")
	v := []byte("From: <sip:abc@biloxi.com;transport=tcp;method=REGISTER>")
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	addr := NewSipHeaderFrom(context)
	header := addr.GetSipHeaderFrom(context)
	header.Parse(context, v, 0)
	remain := context.allocator.Used()
	b.SetBytes(2)
	b.ReportAllocs()
	buf := bytes.NewBuffer(make([]byte, 1024*1024))
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		buf.Reset()
		context.allocator.ClearAllocNum()
		context.allocator.FreePart(remain)
		header.String(context)
	}
} //*/
