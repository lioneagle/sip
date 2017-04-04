package sipparser

import (
	//"bytes"
	//"fmt"
	"testing"
)

//*
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

	prefix := FuncName()

	for i, v := range testdata {
		begin, end, ok := FindCrlfRFC3261([]byte(v.src), 0)

		if v.ok && !ok {
			t.Errorf("%s[%d] failed: should be ok\n", prefix, i)
			continue
		}

		if !v.ok && ok {
			t.Errorf("%s[%d] failed: should be failed\n", prefix, i)
			continue
		}

		if v.begin != begin {
			t.Errorf("%s[%d] failed: begin = %d, wanted = %d\n", prefix, i, begin, v.begin)
			continue
		}

		if v.end != end {
			t.Errorf("%s[%d] failed: end = %d, wanted = %d\n", prefix, i, end, v.end)
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
		{"X1:122334545\r\nX2:tt\r\n", true, len("X1:122334545\r\nX2:tt\r\n"), "X1: 122334545\r\nX2: tt\r\n"},
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
		{"Content-Type: application/isup\r\n", true, len("Content-Type: application/isup\r\n"), "Content-Type: application/isup\r\n"},
		{"Max-Forwards: 1234\r\n", true, len("Max-Forwards: 1234\r\n"), "Max-Forwards: 1234\r\n"},
		{"Route: <tel:12345>\r\nRoute: <sip:456@a.com>\r\n", true, len("Route: <tel:12345>\r\nRoute: <sip:456@a.com>\r\n"), "Route: <tel:12345>, <sip:456@a.com>\r\n"},
		{"Record-Route: <tel:12345>\r\nRecord-Route: <sip:456@a.com>\r\n", true, len("Record-Route: <tel:12345>\r\nRecord-Route: <sip:456@a.com>\r\n"), "Record-Route: <tel:12345>, <sip:456@a.com>\r\n"},
		{"Content-Disposition: early-session\r\n", true, len("Content-Disposition: early-session\r\n"), "Content-Disposition: early-session\r\n"},
		{"Contact: <tel:12345>\r\nContact: sip:456@a.com\r\n", true, len("Contact: <tel:12345>\r\nContact: sip:456@a.com\r\n"), "Contact: <tel:12345>, sip:456@a.com\r\n"},

		{":122334545\r\n", false, 0, ""},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	prefix := FuncName()

	for i, v := range testdata {
		headers := NewSipHeaders()
		newPos, err := headers.Parse(context, []byte(v.src), 0)

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
			continue
		}

		if v.encode != headers.String(context) {
			//t.Errorf("%s[%d] failed: encode = %v, wanted = %v\n", prefix, i, []byte(headers.String(context)), []byte(v.encode))
			t.Errorf("%s[%d] failed: encode = %s, wanted = %s\n", prefix, i, headers.String(context), v.encode)
			continue
		}
	}

}

//*/
