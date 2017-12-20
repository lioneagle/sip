package sipparser

import (
	//"bytes"
	"fmt"
	"testing"
)

func TestSipHeaderContentTypeParse(t *testing.T) {

	testdata := []struct {
		src    string
		ok     bool
		newPos int
		encode string
	}{
		{"Content-Type: application/sdp", true, len("Content-Type: application/sdp"), "Content-Type: application/sdp"},
		{"c: message/sip", true, len("c: message/sip"), "Content-Type: message/sip"},

		{" Content-Type: abc", false, 0, ""},
		{"Content-Type2: abc", false, len("Content-Type2: "), ""},
		{"Content-Type: ", false, len("Content-Type: "), ""},
		{"Content-Type: abc", false, len("Content-Type: abc"), ""},
		{"Content-Type: abc/", false, len("Content-Type: abc/"), ""},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	prefix := FuncName()

	for i, v := range testdata {
		addr := NewSipHeaderContentType(context)
		header := addr.GetSipHeaderContentType(context)
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

func TestSipHeaderContentTypeGetBoundary(t *testing.T) {

	testdata := []struct {
		src         string
		hasBoundary bool
		boundary    string
		encode      string
	}{
		{"Content-Type: application/sdp;boundary=abc", true, "abc", "Content-Type: application/sdp;boundary=abc"},
		{"Content-Type: application/sdp;boundary=\"abc\"", true, "abc", "Content-Type: application/sdp;boundary=\"abc\""},

		{"Content-Type: application/sdp", false, "", ""},
		{"Content-Type: application/sdp;abc;yt;bound=1", false, "", ""},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	prefix := FuncName()

	for i, v := range testdata {
		addr := NewSipHeaderContentType(context)
		header := addr.GetSipHeaderContentType(context)
		header.Parse(context, []byte(v.src), 0)
		boundary, ok := header.GetBoundary(context)

		if v.hasBoundary && !ok {
			t.Errorf("%s[%d] failed: should have boundary\n", prefix, i)
			continue
		}

		if !v.hasBoundary && ok {
			t.Errorf("%s[%d] failed: should not have boundary, boundary = %s\n", prefix, i, boundary)
			continue
		}

		if !v.hasBoundary {
			continue
		}

		if string(boundary) != v.boundary {
			t.Errorf("%s[%d] failed: wrong boundary = %s, wanted = %s\n", prefix, i, string(boundary), v.boundary)
		}

		if header.String(context) != v.encode {
			t.Errorf("%s[%d] failed: wrong encode = %s, wanted = %s\n", prefix, i, header.String(context), v.encode)
		}
	}
}

func TestSipHeaderContentTypeAddBoundary(t *testing.T) {

	testdata := []struct {
		src      string
		boundary string
		encode   string
	}{
		{"Content-Type: application/sdp", "abc", "Content-Type: application/sdp;boundary=\"abc\""},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	prefix := FuncName()

	for i, v := range testdata {
		addr := NewSipHeaderContentType(context)
		header := addr.GetSipHeaderContentType(context)
		header.Parse(context, []byte(v.src), 0)
		err := header.AddBoundary(context, []byte(v.boundary))

		if err != nil {
			t.Errorf("%s[%d] failed: add boundary failed, err = %v\n", prefix, i, err)
			continue
		}

		if header.String(context) != v.encode {
			t.Errorf("%s[%d] failed: wrong encode = %s, wanted = %s\n", prefix, i, header.String(context), v.encode)
		}
	}
}

func TestSipHeaderContentTypeSetType(t *testing.T) {

	testdata := []struct {
		src      string
		mainType string
		subType  string
		encode   string
	}{
		{"Content-Type: application/sdp", "app", "tcp", "Content-Type: app/tcp"},
		{"Content-Type: application/sdp", "message", "sip", "Content-Type: message/sip"},
		{"Content-Type: Application/sdp", "application", "isup", "Content-Type: application/isup"},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	prefix := FuncName()

	for i, v := range testdata {
		addr := NewSipHeaderContentType(context)
		header := addr.GetSipHeaderContentType(context)
		header.Parse(context, []byte(v.src), 0)
		header.SetMainType(context, v.mainType)
		header.SetSubType(context, v.subType)

		if header.String(context) != v.encode {
			t.Errorf("%s[%d] failed: wrong encode = %s, wanted = %s\n", prefix, i, header.String(context), v.encode)
		}
	}
}

func BenchmarkSipHeaderContentTypeParse(b *testing.B) {
	b.StopTimer()
	v := []byte("Content-Type: application/sdp")
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	addr := NewSipHeaderContentType(context)
	header := addr.GetSipHeaderContentType(context)
	remain := context.allocator.Used()
	b.ReportAllocs()
	b.SetBytes(2)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		context.allocator.ClearAllocNum()
		context.allocator.FreePart(remain)
		header.Parse(context, v, 0)
	}
	//fmt.Printf("header = %s\n", header.String())
	fmt.Printf("")
}

func BenchmarkSipHeaderContentTypeEncode(b *testing.B) {
	b.StopTimer()
	v := []byte("Content-Type: application/sdp")
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	addr := NewSipHeaderContentType(context)
	header := addr.GetSipHeaderContentType(context)
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
}
