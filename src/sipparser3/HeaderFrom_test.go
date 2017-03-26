package sipparser3

import (
	"bytes"
	"fmt"
	"testing"
)

func TestSipHeaderFromParse(t *testing.T) {

	testdata := []struct {
		src    string
		ok     bool
		newPos int
		encode string
	}{
		//{"From: sip:abc@a.com;tag=1", true, len("From: sip:abc@a.com;tag=1"), "From: <sip:abc@a.com>;tag=1"},
		//{"From: <sip:abc@a.com;user=ip>;tag=1", true, len("From: <sip:abc@a.com;user=ip>;tag=1"), "From: <sip:abc@a.com;user=ip>;tag=1"},
		{"From: abc<sip:abc@a.com;user=ip>;tag=1", true, len("From: abc<sip:abc@a.com;user=ip>;tag=1"), "From: abc<sip:abc@a.com;user=ip>;tag=1"},
	}

	context := NewParseContext()

	for i, v := range testdata {
		header := NewSipHeaderFrom()
		newPos, err := header.Parse(context, []byte(v.src), 0)

		if v.ok && err != nil {
			t.Errorf("TestSipHeaderFromParse[%d] failed, err = %s\n", i, err)
			continue
		}

		if !v.ok && err == nil {
			t.Errorf("TestSipHeaderFromParse[%d] failed, should parse failed", i)
			continue
		}

		if v.newPos != newPos {
			t.Errorf("TestSipHeaderFromParse[%d] failed, newPos = %d, wanted = %d\n", i, newPos, v.newPos)
		}

		if v.encode != header.String() {
			t.Errorf("TestSipHeaderFromParse[%d] failed, encode = %s, wanted = %s\n", i, header.String(), v.encode)
			continue
		}
	}

}

func BenchmarkSipHeaderFromParse(b *testing.B) {
	b.StopTimer()
	//v := []byte("sip:biloxi.com;transport=tcp;method=REGISTER?to=sip:bob%40biloxi.com")
	//v := []byte("sip:biloxi.com")
	//v := []byte("sip:abc@biloxi.com;transport=tcp")
	v := []byte("From: <sip:abc@biloxi.com;transport=tcp;method=REGISTER>")
	context := NewParseContext()
	header := NewSipHeaderFrom()

	b.ReportAllocs()
	b.SetBytes(2)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		header.Init()
		header.Parse(context, v, 0)
	}
	//fmt.Printf("uri = %s\n", uri.String())
	fmt.Printf("")
}

func BenchmarkSipHeaderFromEncode(b *testing.B) {
	b.StopTimer()
	//v := []byte("sip:biloxi.com;transport=tcp;method=REGISTER?to=sip:bob%40biloxi.com")
	v := []byte("From: <sip:abc@biloxi.com;transport=tcp;method=REGISTER>")
	context := NewParseContext()
	header := NewSipHeaderFrom()
	header.Parse(context, v, 0)
	b.SetBytes(2)
	b.ReportAllocs()

	buf := bytes.NewBuffer(make([]byte, 1024*1024))
	//buf := &bytes.Buffer{}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		buf.Reset()
		header.Encode(buf)
	}
}

func BenchmarkSipHeaderFromString(b *testing.B) {
	b.StopTimer()
	//v := []byte("sip:biloxi.com;transport=tcp;method=REGISTER?to=sip:bob%40biloxi.com")
	v := []byte("From: <sip:abc@biloxi.com;transport=tcp;method=REGISTER>")
	context := NewParseContext()
	header := NewSipHeaderFrom()
	header.Parse(context, v, 0)
	b.SetBytes(2)
	b.ReportAllocs()

	buf := bytes.NewBuffer(make([]byte, 1024*1024))
	//buf := &bytes.Buffer{}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		buf.Reset()
		header.String()
	}
}
