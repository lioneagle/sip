package sipparser2

import (
	"testing"
)

//*
func TestSipNameAddrParseOK(t *testing.T) {

	testdata := []struct {
		src string
	}{
		{"<sip:123@abc.com;ttl=10;user=phone;a;b;c;d;e?xx=yy&x1=aa>"},
		{"\"abc\"<sips:123:tsdd@[1080::8:800:200c:417a]:5061>"},
		{"abc def ee<tel:861234;phone-context=+123>"},
	}

	context := NewParseContext()

	for i, v := range testdata {
		nameaddr := NewSipNameAddr()

		_, err := nameaddr.Parse(context, []byte(v.src), 0)
		if err != nil {
			t.Errorf("TestSipNameAddrParseOK[%d] failed, %s\n", i, err.Error())
			continue
		}
	}
}

//*/
