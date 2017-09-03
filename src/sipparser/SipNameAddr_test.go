package sipparser

import (
	"testing"
)

//*
func TestSipNameAddrParse(t *testing.T) {

	testdata := []struct {
		src    string
		ok     bool
		newPos int
		encode string
	}{
		{"<sip:abc@a.com>", true, len("<sip:abc@a.com>"), "<sip:abc@a.com>"},
		{"<sip:123@abc.com;ttl=10;user=phone;a;b;c;d;e?xx=yy&x1=aa>", true, len("<sip:123@abc.com;ttl=10;user=phone;a;b;c;d;e?xx=yy&x1=aa>"), "<sip:123@abc.com;ttl=10;user=phone;a;b;c;d;e?xx=yy&x1=aa>"},
		{"\"abc\"<sips:123:tsdd@[1080::8:800:200c:417a]:5061>", true, len("\"abc\"<sips:123:tsdd@[1080::8:800:200c:417a]:5061>"), "\"abc\"<sips:123:tsdd@[1080::8:800:200c:417a]:5061>"},
		{"abc def ee<tel:861234;phone-context=+123>", true, len("abc def ee<tel:861234;phone-context=+123>"), "abc def ee<tel:861234;phone-context=+123>"},

		{"\"", false, len("\""), ""},
		{"\r\n<tel:123>", false, len(""), ""},
		{"a b@ c<tel:123>", false, len("a b"), ""},
		{"<tel:", false, len("<tel:"), ""},
		{"<tel:123", false, len("<tel:123"), ""},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	prefix := FuncName()

	for i, v := range testdata {
		addr := NewSipNameAddr(context)
		nameaddr := addr.GetSipNameAddr(context)
		newPos, err := nameaddr.Parse(context, []byte(v.src), 0)

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

		if v.ok && v.encode != nameaddr.String(context) {
			t.Errorf("%s[%d] failed: encode = %s, wanted = %s\n", prefix, i, nameaddr.String(context), v.encode)
			continue
		}
	}
}

//*/

func BenchmarkSipNameAddrSpecParse(b *testing.B) {
	b.StopTimer()
	//v := []byte("sip:biloxi.com;transport=tcp;method=REGISTER?to=sip:bob%40biloxi.com")
	//v := []byte("sip:abc@biloxi.com;transport=tcp")
	v := []byte("<sip:abc@biloxi.com;transport=tcp;method=REGISTER>")
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	addr := NewSipNameAddr(context)
	nameaddr := addr.GetSipNameAddr(context)
	remain := context.allocator.Used()
	b.ReportAllocs()
	b.SetBytes(2)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		context.allocator.ClearAllocNum()
		context.allocator.FreePart(remain)
		nameaddr.Parse(context, v, 0)
	}
}
