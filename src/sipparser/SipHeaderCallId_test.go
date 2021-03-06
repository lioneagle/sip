package sipparser

import (
	//"bytes"
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
	context.allocator = NewMemAllocator(1024 * 30)
	prefix := FuncName()

	for i, v := range testdata {
		addr := NewSipHeaderCallId(context)
		header := addr.GetSipHeaderCallId(context)
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

func BenchmarkSipHeaderCallIdParse(b *testing.B) {
	b.StopTimer()
	//v := []byte("Call-ID: abc123@a.com")
	//v := []byte("Call-ID: 0009b7da-0352000f-30a69b83-0e7b53d6@24.15.255.101")
	v := []byte("Call-ID: 0009b7da-0352000f-30a69b83-0e7b53d6@24.15.255.101")
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	addr := NewSipHeaderCallId(context)
	header := addr.GetSipHeaderCallId(context)
	remain := context.allocator.Used()
	b.ReportAllocs()
	b.SetBytes(2)
	b.StartTimer()

	var i int
	for i = 0; i < b.N; i++ {
		context.allocator.ClearAllocNum()
		context.allocator.FreePart(remain)
		header.Parse(context, v, 0)
	}
	//fmt.Printf("header = %s\n", header.String(context))
	//fmt.Printf("allocator.AllocNum = %d, i= %d\n", context.allocator.AllocNum(), i)
	//fmt.Printf("allocator.Used = %d, i= %d\n", context.allocator.Used(), i)
	//fmt.Printf("")
}

func BenchmarkSipHeaderCallIdEncode(b *testing.B) {
	b.StopTimer()
	//v := []byte("Call-ID: abc123@a.com")
	v := []byte("Call-ID: 0009b7da-0352000f-30a69b83-0e7b53d6@24.15.255.101")
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	addr := NewSipHeaderCallId(context)
	header := addr.GetSipHeaderCallId(context)
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
	//fmt.Println("header =", header.String(context))
}
