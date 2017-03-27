package sipparser3

import (
	"bytes"
	"fmt"
	"testing"
)

func TestSipStartLineParse(t *testing.T) {

	testdata := []struct {
		src       string
		ok        bool
		isRequest bool
		newPos    int
		encode    string
	}{
		{"INVITE sip:123@a.com SIP/2.0\r\n", true, true, len("INVITE sip:123@a.com SIP/2.0\r\n"), "INVITE sip:123@a.com SIP/2.0\r\n"},
		{"SIP/2.0 200 OK\r\n", true, false, len("SIP/2.0 200 OK\r\n"), "SIP/2.0 200 OK\r\n"},
		{"SIP/2.0 200 OK xx\r\n", true, false, len("SIP/2.0 200 OK xx\r\n"), "SIP/2.0 200 OK xx\r\n"},
		{"SIP/2.0 200 \r\n", true, false, len("SIP/2.0 200 \r\n"), "SIP/2.0 200 \r\n"},

		{" INVITE sip:123@a.com SIP/2.0\r\n", false, true, 0, ""},
		{"INVITE", false, true, len("INVITE"), ""},
		{"INVITE@ sip:123@a.com SIP/2.0\r\n", false, true, len("INVITE"), ""},
		{"INVITE sip: SIP/2.0\r\n", false, true, len("INVITE sip:"), ""},
		{"INVITE sip:123@a.com", false, true, len("INVITE sip:123@a.com"), ""},
		{"INVITE sip:123@a.com@ SIP/2.0\r\n", false, true, len("INVITE sip:123@a.com"), ""},
		{"INVITE sip:123@a.com pSIP/2.0\r\n", false, true, len("INVITE sip:123@a.com "), ""},
		{"INVITE sip:123@a.com SIP/2.0", false, true, len("INVITE sip:123@a.com SIP/2.0"), ""},
		{"INVITE sip:123@a.com SIP/2.0\n", false, true, len("INVITE sip:123@a.com SIP/2.0"), ""},
		{"INVITE sip:123@a.com SIP/2.0\r", false, true, len("INVITE sip:123@a.com SIP/2.0"), ""},
		{"INVITE sip:123@a.com SIP/2.0\rt", false, true, len("INVITE sip:123@a.com SIP/2.0"), ""},
		{"INVITE sip:123@a.com SIP/2.0t\n", false, true, len("INVITE sip:123@a.com SIP/2.0"), ""},

		{"pSIP/2.0 200 OK\r\n", false, true, len("pSIP"), ""},
		{"SIP/2.0", false, false, len("SIP/2.0"), ""},
		{"SIP/2.0&", false, false, len("SIP/2.0"), ""},
		{"SIP/2.0 ", false, false, len("SIP/2.0 "), ""},
		{"SIP/2.0 a", false, false, len("SIP/2.0 "), ""},
		{"SIP/2.0 123", false, false, len("SIP/2.0 123"), ""},
		{"SIP/2.0 12a", false, false, len("SIP/2.0 12"), ""},
		{"SIP/2.0 123 ", false, false, len("SIP/2.0 123 "), ""},
		{"SIP/2.0 123 X", false, false, len("SIP/2.0 123 X"), ""},
		{"SIP/2.0 123 X\r", false, false, len("SIP/2.0 123 X"), ""},
		{"SIP/2.0 123 X\n", false, false, len("SIP/2.0 123 X"), ""},
	}

	context := NewParseContext()

	for i, v := range testdata {
		startLine := NewSipStartLine()
		newPos, err := startLine.Parse(context, []byte(v.src), 0)

		if v.ok && err != nil {
			t.Errorf("TestSipStartLineParse[%d] failed, err = %s\n", i, err)
			continue
		}

		if !v.ok && err == nil {
			t.Errorf("TestSipStartLineParse[%d] failed, should parse failed", i)
			continue
		}

		if v.isRequest && !startLine.IsRequest() {
			t.Errorf("TestSipStartLineParse[%d] failed, should be Request-Line", i)
			continue
		}

		if !v.isRequest && startLine.IsRequest() {
			t.Errorf("TestSipStartLineParse[%d] failed, should be Status-Line", i)
			continue
		}

		if v.newPos != newPos {
			t.Errorf("TestSipStartLineParse[%d] failed, newPos = %d, wanted = %d\n", i, newPos, v.newPos)
		}

		if v.ok && v.encode != startLine.String() {
			t.Errorf("TestSipStartLineParse[%d] failed, encode = %s, wanted = %s\n", i, startLine.String(), v.encode)
			continue
		}
	}

}

func BenchmarkSipRequestLineParse(b *testing.B) {
	b.StopTimer()
	v := []byte("INVITE sip:123@a.com SIP/2.0\r\n")
	context := NewParseContext()
	startLine := NewSipStartLine()

	b.ReportAllocs()
	b.SetBytes(2)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		startLine.Init()
		startLine.Parse(context, v, 0)
	}
	//fmt.Printf("uri = %s\n", uri.String())
	fmt.Printf("")
}

func BenchmarkSipRequestLineEncode(b *testing.B) {
	b.StopTimer()
	v := []byte("INVITE sip:123@a.com SIP/2.0\r\n")
	context := NewParseContext()
	startLine := NewSipStartLine()
	startLine.Parse(context, v, 0)
	b.SetBytes(2)
	b.ReportAllocs()

	buf := bytes.NewBuffer(make([]byte, 1024*1024))
	//buf := &bytes.Buffer{}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		buf.Reset()
		startLine.Encode(buf)
	}
}

func BenchmarkSipStatusLineParse(b *testing.B) {
	b.StopTimer()
	v := []byte("SIP/2.0 200 OK\r\n")
	context := NewParseContext()
	startLine := NewSipStartLine()

	b.ReportAllocs()
	b.SetBytes(2)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		startLine.Init()
		startLine.Parse(context, v, 0)
	}
	//fmt.Printf("uri = %s\n", uri.String())
	fmt.Printf("")
}

func BenchmarkSipStatusLineEncode(b *testing.B) {
	b.StopTimer()
	v := []byte("SIP/2.0 200 OK\r\n")
	context := NewParseContext()
	startLine := NewSipStartLine()
	startLine.Parse(context, v, 0)
	b.SetBytes(2)
	b.ReportAllocs()

	buf := bytes.NewBuffer(make([]byte, 1024*1024))
	//buf := &bytes.Buffer{}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		buf.Reset()
		startLine.Encode(buf)
	}
}
