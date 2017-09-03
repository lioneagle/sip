package sipparser

import (
	//"bytes"
	//"fmt"
	"testing"
)

func TestSipHeaderMaxForwardsParse(t *testing.T) {

	testdata := []struct {
		src    string
		ok     bool
		newPos int
		encode string
	}{
		{"Max-Forwards: 1234", true, len("Max-Forwards: 1234"), "Max-Forwards: 1234"},

		{" Max-Forwards: 1234", false, 0, ""},
		{"Max-Forwards2: 1234", false, len("Max-Forwards2: "), ""},
		{"Max-Forwards: ", false, len("Max-Forwards: "), ""},
		{"Max-Forwards: a123", false, len("Max-Forwards: "), ""},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	prefix := FuncName()

	for i, v := range testdata {
		addr := NewSipHeaderMaxForwards(context)
		header := addr.GetSipHeaderMaxForwards(context)
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