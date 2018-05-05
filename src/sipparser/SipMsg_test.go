package sipparser

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const (
	SIP_MSG_BUF_SEPERATOR = "========"
)

type SipMsgBuf struct {
	Name string
	Buf  []byte
}

func NewSipMsgBuf() *SipMsgBuf {
	return &SipMsgBuf{}
}

func (this *SipMsgBuf) Read(src []byte, pos int) (newPos int, ret int) {
	len1 := len(src)
	newPos = pos
	p1 := bytes.Index(src[newPos:], []byte(SIP_MSG_BUF_SEPERATOR))
	if p1 == -1 {
		return newPos, -1
	}
	newPos += p1

	p1 = bytes.IndexByte(src[newPos:], '\n')
	if p1 == -1 {
		return newPos, -1
	}
	newPos += p1 + 1

	for ; newPos < len1; newPos++ {
		if !IsWspChar(src[newPos]) {
			break
		}
	}

	if newPos >= len1 {
		return newPos, -1
	}

	if !bytes.Equal(src[newPos:newPos+4], []byte("name")) {
		fmt.Printf("ERROR: not name after seperator at %d\n", newPos)
		return newPos, 1
	}

	newPos += 4

	for ; newPos < len1; newPos++ {
		if !IsWspChar(src[newPos]) {
			break
		}
	}

	if newPos >= len1 {
		fmt.Printf("ERROR: reach end after name at %d\n", newPos)
		return newPos, 2
	}

	if src[newPos] != '=' {
		fmt.Errorf("ERROR: not '=' after name at %d\n", newPos)
		return newPos, 3
	}
	newPos++

	for ; newPos < len1; newPos++ {
		if !IsWspChar(src[newPos]) {
			break
		}
	}

	if newPos >= len1 {
		fmt.Printf("ERROR: reach end after '=' at %d\n", newPos)
		return newPos, 4
	}

	if src[newPos] != '"' {
		fmt.Printf("ERROR: not '\"' after '=' at %d\n", newPos)
		return newPos, 5
	}

	newPos++

	nameBegin := newPos

	p1 = bytes.IndexByte(src[newPos:], '"')
	if p1 == -1 {
		fmt.Printf("ERROR: not '\"' after name-value at %d\n", newPos)
		return newPos, 6
	}
	newPos += p1

	this.Name = string(src[nameBegin:newPos])

	newPos++

	for ; newPos < len1; newPos++ {
		if !IsLwsChar(src[newPos]) {
			break
		}
	}

	if newPos >= len1 {
		fmt.Printf("ERROR: reach end after name-value at %d\n", newPos)
		return newPos, 7
	}

	bufBegin := newPos

	p1 = bytes.Index(src[newPos:], []byte(SIP_MSG_BUF_SEPERATOR))
	if p1 == -1 {
		fmt.Printf("ERROR: reach end after msg-value at %d\n", newPos)
		return newPos, 8
	}

	newPos += p1

	this.Buf = src[bufBegin : newPos-2]

	return newPos, 0
}

type SipMsgBufs struct {
	Size  int
	Data  map[string]*SipMsgBuf
	Names []string
}

func NewSipMsgBufs() *SipMsgBufs {
	return &SipMsgBufs{Data: make(map[string]*SipMsgBuf)}
}

func (this *SipMsgBufs) ReadFromFile(filename string) bool {
	if this.Size > 0 {
		return true
	}

	src, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Errorf("ERROR: read file %s failed, err =\n", filename, err.Error())
		return false
	}

	pos := 0
	len1 := len(src)

	var ret int

	for {
		buf := NewSipMsgBuf()
		pos, ret = buf.Read(src, pos)
		if ret == -1 {
			return true
		} else if ret != 0 {
			fmt.Println("ERROR: ret =", ret)
			return false
		}

		this.Data[buf.Name] = buf
		this.Names = append(this.Names, buf.Name)
		this.Size++

		if pos >= len1 {
			break
		}
	}

	return true
}

func (this *SipMsgBufs) GetFilteredData(filter string) (ret []*SipMsgBuf) {
	if len(filter) == 0 {
		filter = "."
	}

	if filter != "." {
		for _, v := range this.Names {
			_, ok := ByteSliceIndexNoCase([]byte(v), 0, []byte(filter))
			if ok {
				ret = append(ret, this.Data[v])
			}
		}
	} else {
		for _, v := range this.Names {
			ret = append(ret, this.Data[v])
		}
	}

	return ret
}

var g_sip_msgs *SipMsgBufs = NewSipMsgBufs()

func ReadSipMsgBufs() *SipMsgBufs {
	filename := filepath.FromSlash(os.Args[len(os.Args)-1] + "/src/testdata/sip_msg.txt")
	g_sip_msgs.ReadFromFile(filename)
	//fmt.Println("g_sip_msgs.Size", g_sip_msgs.Size)
	return g_sip_msgs
}

/*func TestSipMsgBufsReadFromFile(t *testing.T) {
	bufs := ReadSipMsgBufs()
	p := bufs.Data["sip_flow_reg_message_200"]

	fmt.Printf("bufs[\"sip_flow_reg_message_200\"] = \n%s", string(p.Buf))
}*/

//*
func TestSipMsgParse1(t *testing.T) {

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
		{"INVITE sip:123@a.com SIP/2.0\r\nVia:", false, len("INVITE sip:123@a.com SIP/2.0\r\nVia:"), ""},
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

var msg2 string = "INVITE sip:6135000@24.15.255.4 SIP/2.0\r\n" +
	"Content-Length : 226\r\n" +
	"Via: SIP/2.0/UDP 24.15.255.101:5060;branch=072c09e5.0\r\n" +
	"From: \"User ID\" <sip:6140000@24.15.255.4;user=phone>;tag=dab70900252036d7134be-4ec05abe\r\n" +
	"To: <sip:6135000@24.15.255.4;user=phone>\r\n" +
	"Call-ID: 0009b7da-0352000f-30a69b83-0e7b53d6@24.15.255.101\r\n" +
	"CSeq: 101 INVITE\r\n" +
	"Expires: 180\r\n" +
	"User-Agent: Cisco-SIP-IP-Phone/2\r\n" +
	"Accept: application/sdp\r\n" +
	"Contact: sip:6140000@24.15.255.101:5060\r\n" +
	"Content-Type: application/sdp\r\n" +
	"\r\n" +
	"v=0\r\n" +
	"o=CiscoSystemsSIP-IPPhone-UserAgent 12748 16809 IN IP4 24.15.255.101\r\n" +
	"s=SIP Call\r\n" +
	"c=IN IP4 24.15.255.101\r\n" +
	"t=0 0\r\nm=audio 26640 RTP/AVP 0 8 18 101\r\n" +
	"a=rtpmap:0 pcmu/8000\r\n" +
	"a=rtpmap:101 telephone-event/8000\r\n" +
	"a=fmtp:101 0-11\r\n"

var msg3 string = "INVITE sip:6135000@24.15.255.4 SIP/2.0\r\n" +
	"Content-Length : 351\r\n" +
	"Via: SIP/2.0/UDP 24.15.255.101:5060;branch=072c09e5.0\r\n" +
	"From: \"User ID\" <sip:6140000@24.15.255.4;user=phone>;tag=dab70900252036d7134be-4ec05abe\r\n" +
	"To: <sip:6135000@24.15.255.4;user=phone>\r\n" +
	"Call-ID: 0009b7da-0352000f-30a69b83-0e7b53d6@24.15.255.101\r\n" +
	"CSeq: 101 INVITE\r\n" +
	"Expires: 180\r\n" +
	"User-Agent: Cisco-SIP-IP-Phone/2\r\n" +
	"Accept: application/sdp\r\n" +
	"Contact: sip:6140000@24.15.255.101:5060\r\n" +
	"Content-Type: multipart/mixed;boundary=abcd\r\n" +
	//"Content-Type: application/sdp;boundary=\"abcd\"\r\n"
	//"Content-Type: application/sdp;boundary=abcd\r\n"
	//"Content-Type: application/sdp\r\n"
	"\r\n" +
	"--abcd\r\n" +
	"Content-Type: application/sdp\r\n" +
	"\r\n" +
	"v=0\r\n" +
	"o=CiscoSystemsSIP-IPPhone-UserAgent 12748 16809 IN IP4 24.15.255.101\r\n" +
	"s=SIP Call\r\n" +
	"c=IN IP4 24.15.255.101\r\n" +
	"t=0 0\r\nm=audio 26640 RTP/AVP 0 8 18 101\r\n" +
	"a=rtpmap:0 pcmu/8000\r\n" +
	"a=rtpmap:101 telephone-event/8000\r\n" +
	"a=fmtp:101 0-11\r\n" +
	"--abcd\r\n" +
	"Content-Type: application/ISUP;version=nxv3;base=etsi121\r\n" +
	"\r\n" +
	"123456\r\n" +
	"--abcd--"

var msg4 string = "INVITE sip:6135000@24.15.255.4 SIP/2.0\r\n" +
	"Content-Length \t: \r\n\t 351\r\n" +
	"Via \t: \r\n\t SIP \r\n\t / \r\n\t 2.0 \r\n\t / \r\n\t UDP \r\n\t 24.15.255.101:5060 \r\n\t ; \r\n\t ttl \r\n\t = \r\n\t 100 \r\n\t ; \r\n\t branch \r\n\t = \r\n\t 072c09e5.0 \r\n\t , \r\n\t SIP \r\n\t / \r\n\t 2.0 \r\n\t / \r\n\t UDP \r\n\t 24.15.255.101:5060 \r\n\t ; \r\n\t ttl \r\n\t = \r\n\t 45 \r\n\t ; \r\n\t branch \r\n\t = \r\n\t 072c09e5.0\r\n" +
	"Via \t: \r\n\t SIP \r\n\t / \r\n\t 2.0 \r\n\t / \r\n\t UDP \r\n\t 24.15.255.101:5060 \r\n\t ; \r\n\t branch \r\n\t = \r\n\t 072c09e5.0\r\n" +
	"Via \t: \r\n\t SIP \r\n\t / \r\n\t 2.0 \r\n\t / \r\n\t UDP \r\n\t 24.15.255.101:5060 \r\n\t ; \r\n\t branch \r\n\t = \r\n\t 072c09e5.0\r\n" +
	"Via \t: \r\n\t SIP \r\n\t / \r\n\t 2.0 \r\n\t / \r\n\t UDP \r\n\t 24.15.255.101:5060 \r\n\t ; \r\n\t branch \r\n\t = \r\n\t 072c09e5.0\r\n" +
	"From \t: \r\n\t \"User ID\" \r\n\t <sip:6140000@24.15.255.4;user=phone> \r\n\t \r\n\t ; \r\n\t tag \r\n\t = \r\n\t dab70900252036d7134be-4ec05abe\r\n" +
	"To \t: \r\n\t \r\n\t <sip:6135000@24.15.255.4;user=phone> \r\n\t \r\n" +
	"Call-ID \t: \r\n\t 0009b7da-0352000f-30a69b83-0e7b53d6@24.15.255.101\r\n" +
	"CSeq \t: \r\n\t 101 \r\n\t INVITE\r\n" +
	"Expires \t: \r\n\t 180\r\n" +
	"User-Agent \t: \r\n\t Cisco-SIP-IP-Phone/2\r\n" +
	"Accept \t: \r\n\t application \r\n\t / \r\n\t sdp \r\n\t ; \r\n\t q \r\n\t = \r\n\t 0.1 \r\n\t , \r\n\t application \r\n\t / \r\n\t isup \r\n\t ; \r\n\t q \r\n\t = \r\n\t 0.3\r\n" +
	"Accept \t: \r\n\t application \r\n\t / \r\n\t sdp \r\n\t ; \r\n\t q \r\n\t = \r\n\t 0.4\r\n" +
	"Contact \t: \r\n\t \r\n\t <sip:6140000@24.15.255.101:5060> \r\n\t \r\n\t ; \r\n\t expires \r\n\t = \r\n\t 1000 \r\n\t , \r\n\t \r\n\t <sip:6140000@24.15.255.101:5060> \r\n\t \r\n\t ; \r\n\t expires \r\n\t = \r\n\t 2000\r\n" +
	"Contact \t: \r\n\t sip:6140000@24.15.255.101:5060 \r\n\t ; \r\n\t expires \r\n\t = \r\n\t 3000\r\n" +
	"Content-Type \t: \r\n\t multipart \r\n\t / \r\n\t mixed \r\n\t ; \r\n\t boundary \r\n\t = \r\n\t abcd\r\n" +
	"\r\n" +
	"--abcd\r\n" +
	"Content-Type \t: \r\n\t application \r\n\t / \r\n\t sdp\r\n" +
	"\r\n" +
	"v=0\r\n" +
	"o=CiscoSystemsSIP-IPPhone-UserAgent 12748 16809 IN IP4 24.15.255.101\r\n" +
	"s=SIP Call\r\n" +
	"c=IN IP4 24.15.255.101\r\n" +
	"t=0 0\r\nm=audio 26640 RTP/AVP 0 8 18 101\r\n" +
	"a=rtpmap:0 pcmu/8000\r\n" +
	"a=rtpmap:101 telephone-event/8000\r\n" +
	"a=fmtp:101 0-11\r\n" +
	"--abcd\r\n" +
	"Content-Type \t: \r\n\t application \r\n\t / \r\n\t ISUP \r\n\t ; \r\n\t version \r\n\t = \r\n\t nxv3 \r\n\t ; \r\n\t base \r\n\t = \r\n\t etsi121\r\n" +
	"\r\n" +
	"123456\r\n" +
	"--abcd--"

func BenchmarkSipMsgParse1(b *testing.B) {
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
		fmt.Printf("len(msg1) = %d\n", len(msg1))
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

func BenchmarkSipMsgParse1_Raw(b *testing.B) {
	b.StopTimer()
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	addr := NewSipMsg(context)
	sipmsg := addr.GetSipMsg(context)
	remain := context.allocator.Used()
	//remainAllocReqBytes := context.allocator.AllocReqBytes()
	//fmt.Printf("allocator.Used = %d\n", context.allocator.Used())
	context.ParseSipHeaderAsRaw = true
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
		fmt.Printf("len(msg1) = %d\n", len(msg1))
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

func BenchmarkSipMsgParse2(b *testing.B) {
	b.StopTimer()
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	addr := NewSipMsg(context)
	sipmsg := addr.GetSipMsg(context)
	remain := context.allocator.Used()
	//remainAllocReqBytes := context.allocator.AllocReqBytes()
	//fmt.Printf("allocator.Used = %d\n", context.allocator.Used())
	//context.ParseSipHeaderAsRaw = true
	msg1 := []byte(msg2)
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
		fmt.Printf("len(msg1) = %d\n", len(msg1))
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

func BenchmarkSipMsgParse3(b *testing.B) {
	b.StopTimer()
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	addr := NewSipMsg(context)
	sipmsg := addr.GetSipMsg(context)
	remain := context.allocator.Used()
	//remainAllocReqBytes := context.allocator.AllocReqBytes()
	//fmt.Printf("allocator.Used = %d\n", context.allocator.Used())
	//context.ParseSipHeaderAsRaw = true
	msg1 := []byte(msg3)
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
	//fmt.Printf("msg = %s\n", sipmsg.String(context))

	/*
		fmt.Printf("total_headers = %d\n", total_headers)
		fmt.Printf("allocator.AllocNum = %d\n", context.allocator.AllocNum())
		fmt.Printf("allocator.Used = %d\n", context.allocator.Used()-remain)
		fmt.Printf("allocator.AllocReqBytes = %d\n", context.allocator.AllocReqBytes()-remainAllocReqBytes)
		fmt.Printf("len(msg1) = %d\n", len(msg1))
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

func BenchmarkSipMsgParse4(b *testing.B) {
	b.StopTimer()
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 300)
	addr := NewSipMsg(context)
	sipmsg := addr.GetSipMsg(context)
	remain := context.allocator.Used()
	//remainAllocReqBytes := context.allocator.AllocReqBytes()
	//fmt.Printf("allocator.Used = %d\n", context.allocator.Used())
	//context.ParseSipHeaderAsRaw = true
	msg1 := []byte(msg4)
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
	//fmt.Printf("msg = %s\n", sipmsg.String(context))

	/*
		fmt.Printf("total_headers = %d\n", total_headers)
		fmt.Printf("allocator.AllocNum = %d\n", context.allocator.AllocNum())
		fmt.Printf("allocator.Used = %d\n", context.allocator.Used()-remain)
		fmt.Printf("allocator.AllocReqBytes = %d\n", context.allocator.AllocReqBytes()-remainAllocReqBytes)
		fmt.Printf("len(msg1) = %d\n", len(msg1))
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

func BenchmarkSipMsgRawParse(b *testing.B) {
	b.StopTimer()
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	addr := NewSipMsg(context)
	sipmsg := addr.GetSipMsg(context)
	remain := context.allocator.Used()
	context.ParseSipHeaderAsRaw = true
	msg1 := []byte(msg)
	total_headers = 0

	b.ReportAllocs()
	b.SetBytes(2)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
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
}

func BenchmarkSipMsgRawScan1(b *testing.B) {
	b.StopTimer()
	msg1 := []byte(msg)
	b.ReportAllocs()
	b.SetBytes(2)
	b.StartTimer()

	//newPos := 0
	var err error

	for i := 0; i < b.N; i++ {
		//newPos, err = SipMsgRawScan1(msg1, 0)
		_, err = SipMsgRawScan1(msg1, 0)
		if err != nil {
			fmt.Println("SipMsgRawScan1 failed, err =", err.Error())
			fmt.Println("msg1 = ", string(msg1))
			break
		} //*/
	}

	//fmt.Println("newPos =", newPos)
}

func BenchmarkSipMsgRawScan2(b *testing.B) {
	b.StopTimer()
	msg1 := []byte(msg)
	b.ReportAllocs()
	b.SetBytes(2)
	b.StartTimer()

	//newPos := 0
	var err error

	for i := 0; i < b.N; i++ {
		//newPos, err = SipMsgRawScan2(msg1, 0)
		_, err = SipMsgRawScan2(msg1, 0)
		if err != nil {
			fmt.Println("SipMsgRawScan2 failed, err =", err.Error())
			fmt.Println("msg1 = ", string(msg1))
			break
		} //*/
	}

	//fmt.Println("newPos =", newPos)
}

func BenchmarkSipMsgRawScan3(b *testing.B) {
	b.StopTimer()
	msg1 := []byte(msg)
	b.ReportAllocs()
	b.SetBytes(2)
	b.StartTimer()

	//newPos := 0
	var err error

	for i := 0; i < b.N; i++ {
		//newPos, err = SipMsgRawScan3(msg1, 0)
		_, err = SipMsgRawScan3(msg1, 0)
		if err != nil {
			fmt.Println("SipMsgRawScan3 failed, err =", err.Error())
			fmt.Println("msg1 = ", string(msg1))
			break
		} //*/
	}

	//fmt.Println("newPos =", newPos)
}

func BenchmarkSipMsgRawScan4(b *testing.B) {
	b.StopTimer()
	msg1 := []byte(msg)
	b.ReportAllocs()
	b.SetBytes(2)
	b.StartTimer()

	//newPos := 0
	var err error

	for i := 0; i < b.N; i++ {
		//newPos, err = SipMsgRawScan4(msg1, 0)
		_, err = SipMsgRawScan4(msg1, 0)
		if err != nil {
			fmt.Println("SipMsgRawScan4 failed, err =", err.Error())
			fmt.Println("msg1 = ", string(msg1))
			break
		} //*/
	}

	//fmt.Println("newPos =", newPos)
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
	buf := &AbnfByteBuffer{}

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

func BenchmarkSipMsgRawEncode(b *testing.B) {
	b.StopTimer()
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	context.ParseSipHeaderAsRaw = true
	addr := NewSipMsg(context)
	sipmsg := addr.GetSipMsg(context)
	msg1 := []byte(msg)
	sipmsg.Parse(context, msg1, 0)
	remain := context.allocator.Used()
	buf := &AbnfByteBuffer{}

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

func BenchmarkSipMsgsRawScan(b *testing.B) {
	bufs := ReadSipMsgBufs()

	testdata := bufs.GetFilteredData(".")

	for _, v := range testdata {
		v := v

		b.Run(v.Name, func(b *testing.B) {
			//b.Parallel()

			b.StopTimer()

			msg := v.Buf
			context := NewParseContext()
			context.allocator = NewMemAllocator(1024 * 10)

			b.StartTimer()

			var err error
			for i := 0; i < b.N; i++ {
				_, err = SipMsgRawScan1(msg, 0)
				if err != nil {
					fmt.Println("SipMsgRawScan4 failed, err =", err.Error())
					fmt.Println("msg1 = ", string(msg))
					break
				}
			}
		})
	}
}

func BenchmarkSipMsgsRawParse(b *testing.B) {
	bufs := ReadSipMsgBufs()

	testdata := bufs.GetFilteredData(".")

	for _, v := range testdata {
		v := v

		b.Run(v.Name, func(b *testing.B) {
			//b.Parallel()

			b.StopTimer()
			context := NewParseContext()
			context.allocator = NewMemAllocator(1024 * 30)
			addr := NewSipMsg(context)
			sipmsg := addr.GetSipMsg(context)
			remain := context.allocator.Used()
			context.ParseSipHeaderAsRaw = true
			msg := v.Buf

			b.ReportAllocs()
			b.SetBytes(2)
			b.StartTimer()

			for i := 0; i < b.N; i++ {
				context.allocator.ClearAllocNum()
				context.allocator.FreePart(remain)
				print_mem = true
				_, err := sipmsg.Parse(context, msg, 0)
				print_mem = false
				if err != nil {
					fmt.Println("parse sip msg failed, err =", err.Error())
					fmt.Println("msg = ", string(msg))
					break
				} //*/
			}
		})
	}
}

func BenchmarkSipMsgsParse(b *testing.B) {
	bufs := ReadSipMsgBufs()

	testdata := bufs.GetFilteredData(".")

	for _, v := range testdata {
		v := v

		b.Run(v.Name, func(b *testing.B) {
			//b.Parallel()

			b.StopTimer()
			context := NewParseContext()
			context.allocator = NewMemAllocator(1024 * 30)
			addr := NewSipMsg(context)
			sipmsg := addr.GetSipMsg(context)
			remain := context.allocator.Used()
			context.ParseSipHeaderAsRaw = false
			msg := v.Buf

			b.ReportAllocs()
			b.SetBytes(2)
			b.StartTimer()

			for i := 0; i < b.N; i++ {
				context.allocator.ClearAllocNum()
				context.allocator.FreePart(remain)
				print_mem = true
				_, err := sipmsg.Parse(context, msg, 0)
				print_mem = false
				if err != nil {
					fmt.Println("parse sip msg failed, err =", err.Error())
					fmt.Println("msg = ", string(msg))
					break
				} //*/
			}
		})
	}
}
