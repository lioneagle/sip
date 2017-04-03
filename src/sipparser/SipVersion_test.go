package sipparser

import (
	"testing"
)

//*
func TestSipVersionParse(t *testing.T) {

	testdata := []struct {
		src    string
		ok     bool
		newPos int
		dst    string
	}{
		{"Sip/2.0", true, len("Sip/2.0"), "SIP/2.0"},
		{"Sip/22.10", true, len("Sip/22.10"), "SIP/22.10"},

		{"Si", false, 0, ""},
		{"Sip", false, 0, ""},
		{"abc/2.0", false, 0, ""},
		{"Sip/a.b", false, len("Sip/"), ""},
		{"Sip/20^b", false, len("Sip/20"), ""},
		{"Sip/20.b", false, len("Sip/20."), ""},
		{"Sip/.b", false, len("Sip/"), ""},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)

	for i, v := range testdata {
		version, _ := NewSipVersion(context)

		newPos, err := version.Parse(context, []byte(v.src), 0)
		if v.ok && err != nil {
			t.Errorf("TestSipVersionParse[%d] failed, %s\n", i, err.Error())
			continue
		}

		if !v.ok && err == nil {
			t.Errorf("TestSipVersionParse[%d] failed, should not be ok\n", i)
			continue
		}

		if newPos != v.newPos {
			t.Errorf("TestSipVersionParse[%d] failed, newPos = %d, wanted = %d\n", i, newPos, v.newPos)
			continue
		}

		if v.ok && v.dst != version.String(context) {
			t.Errorf("TestSipVersionParse[%d] failed, version = %s, wanted = %s\n", i, version.String(context), v.dst)
			continue
		}
	}
}
