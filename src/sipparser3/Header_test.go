package sipparser3

import (
	//"bytes"
	//"fmt"
	"testing"
)

func TestFindCrlfRFC3261(t *testing.T) {

	testdata := []struct {
		src   string
		ok    bool
		begin int
		end   int
	}{
		{"122334545\r\n", true, len("122334545"), len("122334545\r\n")},
		{"122334545\r\nadsad", true, len("122334545"), len("122334545\r\n")},
		{"122334545\n", true, len("122334545"), len("122334545\n")},
		{"122334545\nadsad", true, len("122334545"), len("122334545\n")},

		{"122334545", false, len("122334545"), len("122334545")},
		{"122334545\r", false, len("122334545\r"), len("122334545\r")},
		{"122334545\r\n ", false, len("122334545\r\n "), len("122334545\r\n ")},
		{"122334545\r\n\t", false, len("122334545\r\n\t"), len("122334545\r\n\t")},
	}

	for i, v := range testdata {
		begin, end, ok := FindCrlfRFC3261([]byte(v.src), 0)

		if v.ok && !ok {
			t.Errorf("TestFindCrlfRFC3261[%d] failed, should be ok\n", i)
			continue
		}

		if !v.ok && ok {
			t.Errorf("TestFindCrlfRFC3261[%d] failed, should be failed\n", i)
			continue
		}

		if v.begin != begin {
			t.Errorf("TestFindCrlfRFC3261[%d] failed, begin = %d, wanted = %d\n", i, begin, v.begin)
			continue
		}

		if v.end != end {
			t.Errorf("TestFindCrlfRFC3261[%d] failed, end = %d, wanted = %d\n", i, end, v.end)
			continue
		}
	}
}

//*
func TestParseSipHeaders(t *testing.T) {

	testdata := []struct {
		src    string
		ok     bool
		newPos int
		encode string
	}{
		{"X1:122334545\r\n", true, len("X1:122334545\r\n"), "X1: 122334545\r\n"},
		{"X1\t :122334545\r\n", true, len("X1\t :122334545\r\n"), "X1: 122334545\r\n"},
		{"X1 \t :\r\n\t122334545\r\n", true, len("X1 \t :\r\n\t122334545\r\n"), "X1: 122334545\r\n"},
		{"X1 \t :\r\n 122334545\r\n", true, len("X1 \t :\r\n 122334545\r\n"), "X1: 122334545\r\n"},
		{"X1:122334545\r\nX2:tt\r\n", true, len("X1:122334545\r\nX2:tt\r\n"), "X1: 122334545\r\nX2: tt\r\n"}, //*/
		{"From: <tel:12345>\r\n", true, len("From: <tel:12345>\r\n"), "From: <tel:12345>\r\n"},
		{"To: <tel:12345>\r\n", true, len("To: <tel:12345>\r\n"), "To: <tel:12345>\r\n"},
		{"f: <tel:12345>\r\n", true, len("f: <tel:12345>\r\n"), "From: <tel:12345>\r\n"},
		{"f: <tel:12345>\r\nto: <sip:456@a.com>\r\n", true, len("f: <tel:12345>\r\nto: <sip:456@a.com>\r\n"), "From: <tel:12345>\r\nto: <sip:456@a.com>\r\n"},
		{"Via: SIP/2.0/UDP 10.1.1.1:5060;branch=123\r\n", true, len("Via: SIP/2.0/UDP 10.1.1.1:5060;branch=123\r\n"), "Via: SIP/2.0/UDP 10.1.1.1:5060;branch=123\r\n"},
		{"Via: SIP/2.0/UDP 10.1.1.1:5060;branch=123\r\nVia: SIP/2.0/TCP 10.1.1.1:5060;branch=456\r\n", true, len("Via: SIP/2.0/UDP 10.1.1.1:5060;branch=123\r\nVia: SIP/2.0/TCP 10.1.1.1:5060;branch=456\r\n"), "Via: SIP/2.0/UDP 10.1.1.1:5060;branch=123, SIP/2.0/TCP 10.1.1.1:5060;branch=456\r\n"},
		{"Allow: abc, b34\r\nAllow: hhh\r\n", true, len("Allow: abc, b34\r\nAllow: hhh\r\n"), "Allow: abc, b34, hhh\r\n"},
		{"Call-ID: abc123@a.com\r\n", true, len("Call-ID: abc123@a.com\r\n"), "Call-ID: abc123@a.com\r\n"},
		{"Date: Sat, 13 Nov 2010 23:29:00 GMT\r\n", true, len("Date: Sat, 13 Nov 2010 23:29:00 GMT\r\n"), "Date: Sat, 13 Nov 2010 23:29:00 GMT\r\n"},
		{"CSeq: 1234 INVITE\r\n", true, len("CSeq: 1234 INVITE\r\n"), "CSeq: 1234 INVITE\r\n"},
		{"Content-Length: 1234\r\n", true, len("Content-Length: 1234\r\n"), "Content-Length: 1234\r\n"},
		//*/
	}

	context := NewParseContext()

	for i, v := range testdata {
		headers := NewSipHeaders()
		newPos, err := headers.Parse(context, []byte(v.src), 0)

		if v.ok && err != nil {
			t.Errorf("TestParseSipHeaders[%d] failed, err = %s\n", i, err)
			continue
		}

		if !v.ok && err == nil {
			t.Errorf("TestParseSipHeaders[%d] failed, should parse failed", i)
			continue
		}

		if v.newPos != newPos {
			t.Errorf("TestParseSipHeaders[%d] failed, newPos = %d, wanted = %d\n", i, newPos, v.newPos)
			continue
		}

		if v.encode != headers.String() {
			//t.Errorf("TestParseSipHeaders[%d] failed, encode = %v, wanted = %v\n", i, []byte(headers.String()), []byte(v.encode))
			t.Errorf("TestParseSipHeaders[%d] failed, encode = %s, wanted = %s\n", i, headers.String(), v.encode)
			continue
		}
	}

	/*
		headers := NewSipHeaders()
		headers.Parse(context, []byte("X1:122334545\r\n"), 0)
		headers.headers[0].name.SetValue([]byte("X16"))
		fmt.Println("new header =", headers.String())
		//*/
}

//*/
