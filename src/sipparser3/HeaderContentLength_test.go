package sipparser3

import (
	//"bytes"
	//"fmt"
	"testing"
)

func TestSipHeaderContentLengthParse(t *testing.T) {

	testdata := []struct {
		src    string
		ok     bool
		newPos int
		encode string
	}{
		{"Content-Length: 1234", true, len("Content-Length: 1234"), "Content-Length: 1234"},
		{"l: 1234", true, len("l: 1234"), "Content-Length: 1234"},

		{" Content-Lengt: 1234", false, 0, ""},
		{"Content-Lengt: 1234", false, len("Content-Lengt: "), ""},
		{"Content-Length: ", false, len("Content-Length: "), ""},
		{"Content-Length: a123", false, len("Content-Length: "), ""},
	}

	context := NewParseContext()

	for i, v := range testdata {
		header := NewSipHeaderContentLength()
		newPos, err := header.Parse(context, []byte(v.src), 0)

		if v.ok && err != nil {
			t.Errorf("TestSipHeaderContentLengthParse[%d] failed, err = %s\n", i, err)
			continue
		}

		if !v.ok && err == nil {
			t.Errorf("TestSipHeaderContentLengthParse[%d] failed, should parse failed", i)
			continue
		}

		if v.newPos != newPos {
			t.Errorf("TestSipHeaderContentLengthParse[%d] failed, newPos = %d, wanted = %d\n", i, newPos, v.newPos)
		}

		if v.ok && v.encode != header.String() {
			t.Errorf("TestSipHeaderContentLengthParse[%d] failed, encode = %s, wanted = %s\n", i, header.String(), v.encode)
			continue
		}
	}

}
