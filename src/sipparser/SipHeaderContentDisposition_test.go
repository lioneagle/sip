package sipparser

import (
	//"bytes"
	//"fmt"
	"testing"
)

func TestSipHeaderContentDispositionParse(t *testing.T) {

	testdata := []struct {
		src    string
		ok     bool
		newPos int
		encode string
	}{
		{"Content-Disposition: session", true, len("Content-Disposition: session"), "Content-Disposition: session"},
		{"Content-Disposition: early-session", true, len("Content-Disposition: early-session"), "Content-Disposition: early-session"},

		{" Content-Disposition: abc", false, 0, ""},
		{"Content-Disposition2: abc", false, len("Content-Disposition2: "), ""},
		{"Content-Disposition: ", false, len("Content-Disposition: "), ""},
		{"Content-Disposition: @", false, len("Content-Disposition: "), ""},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	prefix := FuncName()

	for i, v := range testdata {
		addr := NewSipHeaderContentDisposition(context)
		header := addr.GetSipHeaderContentDisposition(context)
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
