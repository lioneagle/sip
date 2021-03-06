package sipparser

import (
	//"bytes"
	"fmt"
	"testing"
)

func TestSipHeaderContactParse(t *testing.T) {

	testdata := []struct {
		src    string
		ok     bool
		newPos int
		encode string
	}{
		{"Contact: *", true, len("Contact: *"), "Contact: *"},
		{"Contact: sip:abc@a.com;tag=1", true, len("Contact: sip:abc@a.com;tag=1"), "Contact: <sip:abc@a.com>;tag=1"},
		{"Contact: <sip:abc@a.com;user=ip>;tag=1", true, len("Contact: <sip:abc@a.com;user=ip>;tag=1"), "Contact: <sip:abc@a.com;user=ip>;tag=1"},
		{"Contact: abc<sip:abc@a.com;user=ip>;tag=1", true, len("Contact: abc<sip:abc@a.com;user=ip>;tag=1"), "Contact: abc<sip:abc@a.com;user=ip>;tag=1"},
		{"Contact: tel:+12358;tag=123", true, len("Contact: tel:+12358;tag=123"), "Contact: <tel:+12358>;tag=123"},

		{" Contact: <sip:abc@a.com>;tag=1", false, 0, "0"},
		{"Contact1: <sip:abc@a.com>;tag=1", false, len("Contact1: "), ""},
		{"Contact: ", false, len("Contact: "), ""},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	prefix := FuncName()

	for i, v := range testdata {
		addr := NewSipHeaderContact(context)
		header := addr.GetSipHeaderContact(context)
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

func BenchmarkSipHeaderContactParse(b *testing.B) {
	b.StopTimer()
	v := []byte("Contact: sip:6140000@24.15.255.101:5060")
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	addr := NewSipHeaderContact(context)
	header := addr.GetSipHeaderContact(context)
	remain := context.allocator.Used()
	b.ReportAllocs()
	b.SetBytes(2)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		context.allocator.ClearAllocNum()
		context.allocator.FreePart(remain)
		header.Parse(context, v, 0)
	}
	//b.StartTimer()
	//fmt.Printf("allocator.AllocNum = %d\n", context.allocator.AllocNum())
	//fmt.Printf("allocator.Used = %d\n", context.allocator.Used()-remain)
	//fmt.Printf("len(header) = %d\n", len(v))
	//fmt.Printf("")
}

func BenchmarkSipHeaderContactEncode(b *testing.B) {
	b.StopTimer()
	v := []byte("Contact: sip:6140000@24.15.255.101:5060")
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 300)
	addr := NewSipHeaderContact(context)
	header := addr.GetSipHeaderContact(context)
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
	fmt.Printf("")
}
