package sipparser3

import (
	//"bytes"
	//"fmt"
	"testing"
)

func TestSipHeaderToParse(t *testing.T) {

	testdata := []struct {
		src    string
		ok     bool
		newPos int
		encode string
	}{
		{"To: sip:abc@a.com;tag=1", true, len("To: sip:abc@a.com;tag=1"), "To: <sip:abc@a.com>;tag=1"},
		{"To: <sip:abc@a.com;user=ip>;tag=1", true, len("To: <sip:abc@a.com;user=ip>;tag=1"), "To: <sip:abc@a.com;user=ip>;tag=1"},
		{"To: abc<sip:abc@a.com;user=ip>;tag=1", true, len("To: abc<sip:abc@a.com;user=ip>;tag=1"), "To: abc<sip:abc@a.com;user=ip>;tag=1"},
		{"To: tel:+12358;tag=123", true, len("To: tel:+12358;tag=123"), "To: <tel:+12358>;tag=123"},

		{" To: <sip:abc@a.com>;tag=1", false, 0, "0"},
		{"To1: <sip:abc@a.com>;tag=1", false, len("To1: "), ""},
		{"To: ", false, len("To: "), ""},
	}

	context := NewParseContext()

	for i, v := range testdata {
		header := NewSipHeaderTo()
		newPos, err := header.Parse(context, []byte(v.src), 0)

		if v.ok && err != nil {
			t.Errorf("TestSipHeaderToParse[%d] failed, err = %s\n", i, err)
			continue
		}

		if !v.ok && err == nil {
			t.Errorf("TestSipHeaderToParse[%d] failed, should parse failed", i)
			continue
		}

		if v.newPos != newPos {
			t.Errorf("TestSipHeaderToParse[%d] failed, newPos = %d, wanted = %d\n", i, newPos, v.newPos)
		}

		if v.ok && v.encode != header.String() {
			t.Errorf("TestSipHeaderToParse[%d] failed, encode = %s, wanted = %s\n", i, header.String(), v.encode)
			continue
		}
	}

}
