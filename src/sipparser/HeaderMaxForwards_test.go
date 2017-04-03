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

	for i, v := range testdata {
		header, _ := NewSipHeaderMaxForwards(context)
		newPos, err := header.Parse(context, []byte(v.src), 0)

		if v.ok && err != nil {
			t.Errorf("TestSipHeaderMaxForwardsParse[%d] failed, err = %s\n", i, err)
			continue
		}

		if !v.ok && err == nil {
			t.Errorf("TestSipHeaderMaxForwardsParse[%d] failed, should parse failed", i)
			continue
		}

		if v.newPos != newPos {
			t.Errorf("TestSipHeaderMaxForwardsParse[%d] failed, newPos = %d, wanted = %d\n", i, newPos, v.newPos)
		}

		if v.ok && v.encode != header.String(context) {
			t.Errorf("TestSipHeaderMaxForwardsParse[%d] failed, encode = %s, wanted = %s\n", i, header.String(context), v.encode)
			continue
		}
	}

}
