package sipparser

import (
	"bytes"
	"fmt"
	"testing"
)

func TestSipHeaderToParse(t *testing.T) {

	testdata := []struct {
		src    string
		ok     bool
		newPos int
		encode string
	}{
		{"To: sip:abc@a.com;tag=1", true, len("To: sip:abc@a.com;tag=1"), "To: <sip:abc@a.com>;tag=1"},
		{"To: <sip:abc@a.com;user=ip>;tag=1", true, len("To: <sip:abc@a.com;user=ip>;tag=1"), "To: <sip:abc@a.com;user=ip>;tag=1"},
		{"To: abc<sip:abc@a.com;user=ip>;tag=1", true, len("To: abc<sip:abc@a.com;user=ip>;tag=1"), "To: abc<sip:abc@a.com;user=ip>;tag=1"},
		{"To: tel:+12358;tag=123", true, len("To: tel:+12358;tag=123"), "To: <tel:+12358>;tag=123"},

		{" To: <sip:abc@a.com>;tag=1", false, 0, "0"},
		{"To1: <sip:abc@a.com>;tag=1", false, len("To1: "), ""},
		{"To: ", false, len("To: "), ""},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	prefix := FuncName()

	for i, v := range testdata {
		addr := NewSipHeaderTo(context)
		header := addr.GetSipHeaderTo(context)
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

func BenchmarkSipHeaderToParse(b *testing.B) {
	b.StopTimer()
	v := []byte("To: <sip:6135000@24.15.255.4>")
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	addr := NewSipHeaderTo(context)
	header := addr.GetSipHeaderTo(context)
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

func BenchmarkSipHeaderToEncode(b *testing.B) {
	b.StopTimer()
	v := []byte("To: <sip:6135000@24.15.255.4>")
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	addr := NewSipHeaderTo(context)
	header := addr.GetSipHeaderTo(context)
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
		header.Encode(context, buf)
	}
}
