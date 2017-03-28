package sipparser3

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

	for i, v := range testdata {
		header := NewSipHeaderContentDisposition()
		newPos, err := header.Parse(context, []byte(v.src), 0)

		if v.ok && err != nil {
			t.Errorf("TestSipHeaderContentDispositionParse[%d] failed, err = %s\n", i, err)
			continue
		}

		if !v.ok && err == nil {
			t.Errorf("TestSipHeaderContentDispositionParse[%d] failed, should parse failed", i)
			continue
		}

		if v.newPos != newPos {
			t.Errorf("TestSipHeaderContentDispositionParse[%d] failed, newPos = %d, wanted = %d\n", i, newPos, v.newPos)
		}

		if v.ok && v.encode != header.String() {
			t.Errorf("TestSipHeaderContentDispositionParse[%d] failed, encode = %s, wanted = %s\n", i, header.String(), v.encode)
			continue
		}
	}

}
