package sipparser3

import (
	//"bytes"
	//"fmt"
	"testing"
)

func TestSipHeaderRecordRouteParse(t *testing.T) {

	testdata := []struct {
		src    string
		ok     bool
		newPos int
		encode string
	}{
		{"Record-Route: <sip:abc@a.com>;tag=1", true, len("Record-Route: <sip:abc@a.com>;tag=1"), "Record-Route: <sip:abc@a.com>;tag=1"},
		{"Record-Route: <sip:abc@a.com;user=ip>;tag=1", true, len("Record-Route: <sip:abc@a.com;user=ip>;tag=1"), "Record-Route: <sip:abc@a.com;user=ip>;tag=1"},
		{"Record-Route: abc<sip:abc@a.com;user=ip>;tag=1", true, len("Record-Route: abc<sip:abc@a.com;user=ip>;tag=1"), "Record-Route: abc<sip:abc@a.com;user=ip>;tag=1"},
		{"Record-Route: <tel:+12358;tag=123>", true, len("Record-Route: <tel:+12358;tag=123>"), "Record-Route: <tel:+12358;tag=123>"},

		{" Record-Route: <sip:abc@a.com>;tag=1", false, 0, "0"},
		{"Record-Route1: <sip:abc@a.com>;tag=1", false, len("Record-Route1: "), ""},
		{"Record-Route: ", false, len("Record-Route: "), ""},
	}

	context := NewParseContext()

	for i, v := range testdata {
		header := NewSipHeaderRecordRoute()
		newPos, err := header.Parse(context, []byte(v.src), 0)

		if v.ok && err != nil {
			t.Errorf("TestSipHeaderRecordRouteParse[%d] failed, err = %s\n", i, err)
			continue
		}

		if !v.ok && err == nil {
			t.Errorf("TestSipHeaderRecordRouteParse[%d] failed, should parse failed", i)
			continue
		}

		if v.newPos != newPos {
			t.Errorf("TestSipHeaderRecordRouteParse[%d] failed, newPos = %d, wanted = %d\n", i, newPos, v.newPos)
		}

		if v.ok && v.encode != header.String() {
			t.Errorf("TestSipHeaderRecordRouteParse[%d] failed, encode = %s, wanted = %s\n", i, header.String(), v.encode)
			continue
		}
	}

}
