package sipparser

import (
	"bytes"
	"fmt"
	"testing"
	//"unsafe"
)

//*
func TestSipMsgParse(t *testing.T) {

	testdata := []struct {
		src    string
		ok     bool
		newPos int
		encode string
	}{
		{"INVITE sip:123@a.com SIP/2.0\r\nFrom: sip:abc@a.com;tag=1\r\nAllow: a, b\r\n\r\n", true, len("INVITE sip:123@a.com SIP/2.0\r\nFrom: sip:abc@a.com;tag=1\r\nAllow: a, b\r\n\r\n"), "INVITE sip:123@a.com SIP/2.0\r\nFrom: <sip:abc@a.com>;tag=1\r\nAllow: a, b\r\n\r\n"},

		{" INVITE sip:123@a.com SIP/2.0\r\n", false, 0, ""},
		{"INVITE sip:123@a.com SIP/2.0\r\n", false, len("INVITE sip:123@a.com SIP/2.0\r\n"), ""},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 100)
	prefix := FuncName()

	for i, v := range testdata {
		sipmsg, _ := NewSipMsg(context)
		newPos, err := sipmsg.Parse(context, []byte(v.src), 0)

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
		}

		if v.ok && v.encode != sipmsg.String(context) {
			t.Errorf("%s[%d] failed: encode = %s, wanted = %s\n", prefix, i, sipmsg.String(context), v.encode)
			continue
		}
	}

} //*/

var msg string = "INVITE sip:6135000@24.15.255.4 SIP/2.0\r\n" +
	"Content-Length: 226\r\n" +
	"Via: SIP/2.0/UDP 24.15.255.101:5060\r\n" +
	"From: \"User ID\" <sip:6140000@24.15.255.4>;tag=dab70900252036d7134be-4ec05abe\r\n" +
	"To: <sip:6135000@24.15.255.4>\r\n" +
	"Call-ID: 0009b7da-0352000f-30a69b83-0e7b53d6@24.15.255.101\r\n" +
	"CSeq: 101 INVITE\r\n" +
	//"Expires: 180\r\n" +
	//"User-Agent: Cisco-SIP-IP-Phone/2\r\n" +
	//"Accept: application/sdp\r\n" +
	"Contact: sip:6140000@24.15.255.101:5060\r\n" +
	"Content-Type: application/sdp\r\n" +
	"\r\n"

func BenchmarkSipMsgParse(b *testing.B) {
	b.StopTimer()
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	sipmsg, _ := NewSipMsg(context)
	remain := context.allocator.Used()
	msg1 := []byte(msg)
	total_headers = 0

	b.ReportAllocs()
	b.SetBytes(2)
	b.StartTimer()

	//for i := 0; i < 1; i++ {
	for i := 0; i < b.N; i++ {
		context.allocator.ClearAllocNum()
		context.allocator.FreePart(remain)
		print_mem = false
		_, err := sipmsg.Parse(context, msg1, 0)
		print_mem = false
		if err != nil {
			fmt.Println("parse sip msg failed, err =", err.Error())
			fmt.Println("msg1 = ", string(msg1))
			break
		} //*/
	}
	//fmt.Printf("msg = %s\n", sipmsg.String())
	//fmt.Printf("total_headers = %d\n", total_headers)
	//fmt.Printf("allocator.AllocNum = %d\n", context.allocator.AllocNum())
	//fmt.Printf("allocator.Used = %d\n", context.allocator.Used())
	//fmt.Printf("len(msg) = %d\n", len(msg))
	//fmt.Printf("sizeof(SipMsg) =%d\n", unsafe.Sizeof(AbnfBuf{}))

}

func BenchmarkSipMsgEncode(b *testing.B) {
	b.StopTimer()
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	sipmsg, _ := NewSipMsg(context)
	msg1 := []byte(msg)
	sipmsg.Parse(context, msg1, 0)
	remain := context.allocator.Used()
	buf := bytes.NewBuffer(make([]byte, 1024*1024))
	b.SetBytes(2)
	b.ReportAllocs()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		buf.Reset()
		context.allocator.ClearAllocNum()
		context.allocator.FreePart(remain)
		sipmsg.Encode(context, buf)
	}
	fmt.Printf("")
}
