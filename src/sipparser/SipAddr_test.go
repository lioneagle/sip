package sipparser

import (
	"fmt"
	"testing"
)

//*
func TestSipAddrParse(t *testing.T) {

	testdata := []struct {
		src        string
		ok         bool
		newPos     int
		isNameAddr bool
		encode     string
	}{
		{"<sip:abc@a.com>", true, len("<sip:abc@a.com>"), true, "<sip:abc@a.com>"},
		{"<sip:123@abc.com;ttl=10;user=phone;a;b;c;d;e?xx=yy&x1=aa>", true, len("<sip:123@abc.com;ttl=10;user=phone;a;b;c;d;e?xx=yy&x1=aa>"), true, "<sip:123@abc.com;ttl=10;user=phone;a;b;c;d;e?xx=yy&x1=aa>"},
		{"sips:123:tsdd@[1080::8:800:200c:417a]:5061;asd>", true, len("sips:123:tsdd@[1080::8:800:200c:417a]:5061"), false, "<sips:123:tsdd@[1080::8:800:200c:417a]:5061>"},
		{"abc def ee<tel:861234;phone-context=+123>", true, len("abc def ee<tel:861234;phone-context=+123>"), true, "abc def ee<tel:861234;phone-context=+123>"},

		{"\"", false, len("\""), false, ""},
		{"\r\n<tel:123>", false, len(""), false, ""},
		{"a b@ c<tel:123>", false, len("a b"), false, ""},
		{"<tel:", false, len("<tel:"), false, ""},
		{"<tel:123", false, len("<tel:123"), false, ""},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	prefix := FuncName()

	for i, v := range testdata {
		context.allocator.FreeAll()
		addr := NewSipAddr(context)
		sipaddr := addr.GetSipAddr(context)
		newPos, err := sipaddr.Parse(context, []byte(v.src), 0)

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

		if v.ok && v.encode != sipaddr.String(context) {
			t.Errorf("%s[%d] failed: encode = %s, wanted = %s\n", prefix, i, sipaddr.String(context), v.encode)
			continue
		}
	}
}

//*/

func BenchmarkSipAddrParse_1(b *testing.B) {
	b.StopTimer()
	v := []byte("\"string\" <sip:abc@biloxi.com;transport=tcp;method=REGISTER>")
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	ptr := NewSipUri(context)
	addr := ptr.GetSipAddr(context)
	remain := context.allocator.Used()
	b.ReportAllocs()
	b.SetBytes(2)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		context.allocator.ClearAllocNum()
		context.allocator.FreePart(remain)
		addr.Parse(context, v, 0)

	}
	//fmt.Printf("uri = %s\n", uri.String())
	fmt.Printf("")
}
