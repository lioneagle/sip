package sipparser3

import (
	//"bytes"
	//"fmt"
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

	for i, v := range testdata {
		header := NewSipHeaderContentType()
		newPos, err := header.Parse(context, []byte(v.src), 0)

		if v.ok && err != nil {
			t.Errorf("TestSipHeaderContentTypeParse[%d] failed, err = %s\n", i, err)
			continue
		}

		if !v.ok && err == nil {
			t.Errorf("TestSipHeaderContentTypeParse[%d] failed, should parse failed", i)
			continue
		}

		if v.newPos != newPos {
			t.Errorf("TestSipHeaderContentTypeParse[%d] failed, newPos = %d, wanted = %d\n", i, newPos, v.newPos)
		}

		if v.ok && v.encode != header.String() {
			t.Errorf("TestSipHeaderContentTypeParse[%d] failed, encode = %s, wanted = %s\n", i, header.String(), v.encode)
			continue
		}
	}

}
