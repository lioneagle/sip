package sipparser

import (
	"bytes"
	"fmt"
	"testing"
	_ "unsafe"
)

//*
func TestSipMsgParse(t *testing.T) {

	testdata := []struct {
		src    string
		ok     bool
		newPos int
		encode string
	}{
		{"INVITE sip:123@a.com SIP/2.0\r\nFrom: sip:abc@a.com;tag=1\r\nAllow: a, b\r\nContent-Length: 123\r\n\r\n", true, len("INVITE sip:123@a.com SIP/2.0\r\nFrom: sip:abc@a.com;tag=1\r\nAllow: a, b\r\nContent-Length: 123\r\n\r\n"), "INVITE sip:123@a.com SIP/2.0\r\nFrom: <sip:abc@a.com>;tag=1\r\nContent-Length:          0\r\nAllow: a, b\r\n\r\n"},
		{"INVITE sip:123@a.com SIP/2.0\r\nFrom: sip:abc@a.com;tag=1\r\nAllow: a, b\r\n\r\n", true, len("INVITE sip:123@a.com SIP/2.0\r\nFrom: sip:abc@a.com;tag=1\r\nAllow: a, b\r\n\r\n"), "INVITE sip:123@a.com SIP/2.0\r\nFrom: <sip:abc@a.com>;tag=1\r\nContent-Length:          0\r\nAllow: a, b\r\n\r\n"},

		{" INVITE sip:123@a.com SIP/2.0\r\n", false, 0, ""},
		{"INVITE sip:123@a.com SIP/2.0\r\n", false, len("INVITE sip:123@a.com SIP/2.0\r\n"), ""},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 100)
	prefix := FuncName()

	for i, v := range testdata {
		addr := NewSipMsg(context)
		sipmsg := addr.GetSipMsg(context)
		newPos, err := sipmsg.Parse(context, []byte(v.src), 0)

		if v.ok && err != nil {
			t.Errorf("%s[%d] failed: parse sip-msg err = %s\n", prefix, i, err)
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

}

func TestSipMsgParseWithOneBody(t *testing.T) {

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 100)
	prefix := FuncName()

	src := "INVITE sip:6135000@24.15.255.4 SIP/2.0\r\n" +
		"Content-Length: 10\r\n" +
		"Via: SIP/2.0/UDP 24.15.255.101:5060\r\n" +
		"From: \"User ID\" <sip:6140000@24.15.255.4>;tag=dab70900252036d7134be-4ec05abe\r\n" +
		"To: <sip:6135000@24.15.255.4>\r\n" +
		"Call-ID: 0009b7da-0352000f-30a69b83-0e7b53d6@24.15.255.101\r\n" +
		"CSeq: 101 INVITE\r\n" +
		"Content-Disposition: render\r\n" +
		"Content-XYZ: abc, def\r\n" +
		"Content-XYZ: 123\r\n" +
		"Content-Encoding: gzip, tar\r\n" +
		"Contact: sip:6140000@24.15.255.101:5060\r\n" +
		"Content-Type: application/sdp\r\n" +
		"\r\n" +
		"1234567890"

	dst := "INVITE sip:6135000@24.15.255.4 SIP/2.0\r\n" +
		"Content-Length:         10\r\n" +
		"From: \"User ID\"<sip:6140000@24.15.255.4>;tag=dab70900252036d7134be-4ec05abe\r\n" +
		"To: <sip:6135000@24.15.255.4>\r\n" +
		"Call-ID: 0009b7da-0352000f-30a69b83-0e7b53d6@24.15.255.101\r\n" +
		"CSeq: 101 INVITE\r\n" +
		"Content-Disposition: render\r\n" +
		"Content-Type: application/sdp\r\n" +
		"Via: SIP/2.0/UDP 24.15.255.101:5060\r\n" +
		"Content-Encoding: gzip, tar\r\n" +
		"Contact: <sip:6140000@24.15.255.101:5060>\r\n" +
		"Content-XYZ: abc, def\r\n" +
		"Content-XYZ: 123\r\n" +
		"\r\n" +
		"1234567890"

	addr := NewSipMsg(context)
	sipmsg := addr.GetSipMsg(context)
	_, err := sipmsg.Parse(context, []byte(src), 0)

	if err != nil {
		t.Errorf("%s failed: parse sip-msg err = %s\n", prefix, err)
		return
	}

	encoded := sipmsg.String(context)

	if encoded != dst {
		t.Errorf("%s failed: \nencode = \n%s \n\nwanted = \n%s\n", prefix, encoded, dst)
		return
	}
}

func TestSipMsgParseWithMultiBody(t *testing.T) {

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 50)
	prefix := FuncName()

	boundary := "simple-boundary"
	src := "INVITE sip:6135000@24.15.255.4 SIP/2.0\r\n" +
		"Content-Length: 10\r\n" +
		"Via: SIP/2.0/UDP 24.15.255.101:5060\r\n" +
		"From: \"User ID\" <sip:6140000@24.15.255.4>;tag=dab70900252036d7134be-4ec05abe\r\n" +
		"To: <sip:6135000@24.15.255.4>\r\n" +
		"Call-ID: 0009b7da-0352000f-30a69b83-0e7b53d6@24.15.255.101\r\n" +
		"CSeq: 101 INVITE\r\n" +
		"Contact: sip:6140000@24.15.255.101:5060\r\n" +
		"Content-Type: multipart/mixed;boundary=" + boundary + "\r\n" +
		"\r\n" +
		"--" + boundary + "padding\r\n" +
		"Content-Encoding: gzip, tar\r\n" +
		"\r\n" +
		"1234567890" +
		"\r\n" +
		"--" + boundary + "padding\r\n" +
		"Content-XYZ: abc, def\r\n" +
		"\r\n" +
		"abcsdfsdfsf" +
		"\r\n" +
		"--" + boundary + "--"

	dst := "INVITE sip:6135000@24.15.255.4 SIP/2.0\r\n" +
		"Content-Length:        138\r\n" +
		"From: \"User ID\"<sip:6140000@24.15.255.4>;tag=dab70900252036d7134be-4ec05abe\r\n" +
		"To: <sip:6135000@24.15.255.4>\r\n" +
		"Call-ID: 0009b7da-0352000f-30a69b83-0e7b53d6@24.15.255.101\r\n" +
		"CSeq: 101 INVITE\r\n" +
		"Content-Type: multipart/mixed;boundary=" + boundary + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Via: SIP/2.0/UDP 24.15.255.101:5060\r\n" +
		"Contact: <sip:6140000@24.15.255.101:5060>\r\n" +
		"\r\n" +
		"--" + boundary + "\r\n" +
		"Content-Encoding: gzip, tar\r\n" +
		"\r\n" +
		"1234567890" +
		"\r\n" +
		"--" + boundary + "\r\n" +
		"Content-XYZ: abc, def\r\n" +
		"\r\n" +
		"abcsdfsdfsf" +
		"\r\n" +
		"--" + boundary + "--"

	//fmt.Println(context.allocator.String(0, 10))

	addr := NewSipMsg(context)
	sipmsg := addr.GetSipMsg(context)
	_, err := sipmsg.Parse(context, []byte(src), 0)

	if err != nil {
		t.Errorf("%s failed: parse sip-msg err = %s\n", prefix, err)
		return
	}

	//fmt.Println(context.allocator.String(0, 10))

	//fmt.Println("src len =", len(src))

	encoded := sipmsg.String(context)

	if encoded != dst {
		t.Errorf("%s failed: \nencode = \n%s \n\nwanted = \n%s\n", prefix, encoded, dst)
		return
	}
}

func TestSipMsgParseWithMultiBody2(t *testing.T) {

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 50)
	prefix := FuncName()

	boundary := "simple-boundary"
	src := "INVITE sip:6135000@24.15.255.4 SIP/2.0\r\n" +
		"Content-Length: 10\r\n" +
		"Via: SIP/2.0/UDP 24.15.255.101:5060\r\n" +
		"From: \"User ID\" <sip:6140000@24.15.255.4>;tag=dab70900252036d7134be-4ec05abe\r\n" +
		"To: <sip:6135000@24.15.255.4>\r\n" +
		"Call-ID: 0009b7da-0352000f-30a69b83-0e7b53d6@24.15.255.101\r\n" +
		"CSeq: 101 INVITE\r\n" +
		"Contact: sip:6140000@24.15.255.101:5060\r\n" +
		"Content-Type: multipart/mixed;boundary=" + boundary + "\r\n" +
		"\r\n" +
		"--" + boundary + "padding\r\n" +
		"Content-Encoding: gzip, tar\r\n" +
		"\r\n" +
		"1234567890" +
		"\r\n" +
		"--" + boundary + "padding\r\n" +
		"Content-XYZ: abc, def\r\n" +
		"\r\n" +
		"abcsdfsdfsf" +
		"\r\n" +
		"--" + boundary + "--"

	dst := "INVITE sip:6135000@24.15.255.4 SIP/2.0\r\n" +
		"Content-Length:        186\r\n" +
		"From: \"User ID\"<sip:6140000@24.15.255.4>;tag=dab70900252036d7134be-4ec05abe\r\n" +
		"To: <sip:6135000@24.15.255.4>\r\n" +
		"Call-ID: 0009b7da-0352000f-30a69b83-0e7b53d6@24.15.255.101\r\n" +
		"CSeq: 101 INVITE\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: multipart/mixed;boundary=\"" + ABNF_SIP_DEFAULT_BOUNDARY + "\"\r\n" +
		"Via: SIP/2.0/UDP 24.15.255.101:5060\r\n" +
		"Contact: <sip:6140000@24.15.255.101:5060>\r\n" +
		"\r\n" +
		"--" + ABNF_SIP_DEFAULT_BOUNDARY + "\r\n" +
		"Content-Encoding: gzip, tar\r\n" +
		"\r\n" +
		"1234567890" +
		"\r\n" +
		"--" + ABNF_SIP_DEFAULT_BOUNDARY + "\r\n" +
		"Content-XYZ: abc, def\r\n" +
		"\r\n" +
		"abcsdfsdfsf" +
		"\r\n" +
		"--" + ABNF_SIP_DEFAULT_BOUNDARY + "--"

	//fmt.Println(context.allocator.String(0, 10))

	addr := NewSipMsg(context)
	sipmsg := addr.GetSipMsg(context)
	_, err := sipmsg.Parse(context, []byte(src), 0)

	if err != nil {
		t.Errorf("%s failed: parse sip-msg err = %s\n", prefix, err)
		return
	}

	sipmsg.headers.singleHeaders.RemoveHeaderByNameString(context, "content-type")

	encoded := sipmsg.String(context)

	if encoded != dst {
		t.Errorf("%s failed: \nencode = \n%s \n\nwanted = \n%s\n", prefix, encoded, dst)
		return
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
	addr := NewSipMsg(context)
	sipmsg := addr.GetSipMsg(context)
	remain := context.allocator.Used()
	//remainAllocReqBytes := context.allocator.AllocReqBytes()
	//fmt.Printf("allocator.Used = %d\n", context.allocator.Used())
	//context.ParseSipHeaderAsRaw = true
	msg1 := []byte(msg)
	//msg2 := make([]byte, len(msg1))
	total_headers = 0

	b.ReportAllocs()
	b.SetBytes(2)
	b.StartTimer()

	//for i := 0; i < 1; i++ {
	for i := 0; i < b.N; i++ {
		//copy(msg2, msg1[:len(msg1)])
		context.allocator.ClearAllocNum()
		context.allocator.FreePart(remain)
		print_mem = true
		_, err := sipmsg.Parse(context, msg1, 0)
		print_mem = false
		if err != nil {
			fmt.Println("parse sip msg failed, err =", err.Error())
			fmt.Println("msg1 = ", string(msg1))
			break
		} //*/
	}
	//fmt.Printf("msg = %s\n", sipmsg.String())

	/*
		fmt.Printf("total_headers = %d\n", total_headers)
		fmt.Printf("allocator.AllocNum = %d\n", context.allocator.AllocNum())
		fmt.Printf("allocator.Used = %d\n", context.allocator.Used()-remain)
		fmt.Printf("allocator.AllocReqBytes = %d\n", context.allocator.AllocReqBytes()-remainAllocReqBytes)
		fmt.Printf("len(msg) = %d\n", len(msg))
		fmt.Printf("sizeof(SipMsg) =%d\n", unsafe.Sizeof(AbnfBuf{}))
		fmt.Printf("sizeof(SipAddrSpec) =%d\n", unsafe.Sizeof(SipAddrSpec{}))
		fmt.Printf("sizeof(SipAddr) =%d\n", unsafe.Sizeof(SipAddr{}))
		fmt.Printf("sizeof(SipNameAddr) =%d\n", unsafe.Sizeof(SipNameAddr{}))
		fmt.Printf("sizeof(SipUri) =%d\n", unsafe.Sizeof(SipUri{}))
		fmt.Printf("sizeof(AbnfBuf) =%d\n", unsafe.Sizeof(AbnfBuf{}))
		fmt.Printf("sizeof(SipHostPort) =%d\n", unsafe.Sizeof(SipHostPort{}))
		fmt.Printf("sizeof(SipUriParams) =%d\n", unsafe.Sizeof(SipUriParams{}))
		fmt.Printf("sizeof(SipUriHeaders) =%d\n", unsafe.Sizeof(SipUriHeaders{}))

		//*/

}

func BenchmarkSipMsgEncode(b *testing.B) {
	b.StopTimer()
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	addr := NewSipMsg(context)
	sipmsg := addr.GetSipMsg(context)
	msg1 := []byte(msg)
	sipmsg.Parse(context, msg1, 0)
	remain := context.allocator.Used()
	buf := bytes.NewBuffer(make([]byte, 1024*1024))
	//fmt.Println("BenchmarkSipMsgEncode: bodies.Size() =", sipmsg.bodies.Size())

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
