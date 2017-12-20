package sipparser

import (
	"bytes"
	"fmt"
	"net"
	"testing"
)

func TestSipHostUnknownString(t *testing.T) {
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)

	addr := NewSipHost(context)
	host := addr.GetSipHost(context)

	str := host.String(context)

	if str != "unknown host" {
		t.Errorf("TestSipHostUnknownString failed: str = %s, wanted = %s\n", str, "unknown host")
	}
}

func TestSipHostIpv4String(t *testing.T) {
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)

	addr := NewSipHost(context)
	host := addr.GetSipHost(context)

	//ipv4 := net.IPv4(10, 1, 1, 1)
	ipv4 := []byte{10, 1, 1, 1}

	host.SetIpv4(context, ipv4)

	str := host.String(context)

	if str != "10.1.1.1" {
		t.Errorf("TestSipHostIpv4String failed: str = %s, wanted = %s\n", str, "10.1.1.1")
	}
}

func TestSipHostIpv6String(t *testing.T) {
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)

	addr := NewSipHost(context)
	host := addr.GetSipHost(context)

	ipv4 := net.IPv4(10, 1, 1, 1)

	host.SetIpv6(context, ipv4.To16())

	str := host.String(context)

	if str != "[10.1.1.1]" {
		t.Errorf("TestSipHostIpv6String failed: str = %s, wanted = %s\n", str, "[10.1.1.1]")
	}
}

func TestSipHostHostnameString(t *testing.T) {
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)

	addr := NewSipHost(context)
	host := addr.GetSipHost(context)

	host.SetHostname(context, []byte("abc.com"))

	if !host.IsHostname() {
		t.Errorf("TestSipHostHostnameString failed: host is not hostname\n")
	}

	str := host.String(context)

	if str != "abc.com" {
		t.Errorf("TestSipHostHostnameString failed: str = %s, wanted = %s\n", str, "abc.com")
	}
}

func TestSipHostParseOk(t *testing.T) {
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)

	addr := NewSipHost(context)
	host := addr.GetSipHost(context)

	testdata := []struct {
		test     string
		wanted   string
		newPos   int
		hosttype func() bool
	}{
		{"10.43.12.14", "10.43.12.14", len("10.43.12.14"), host.IsIpv4},
		{"10.43.12.14!", "10.43.12.14", len("10.43.12.14"), host.IsIpv4},
		{"10.43.12.14ab", "10.43.12.14ab", len("10.43.12.14ab"), host.IsHostname},
		{"10.43.12.14.", "10.43.12.14.", len("10.43.12.14."), host.IsHostname},
		{"10.43.12.14-", "10.43.12.14-", len("10.43.12.14-"), host.IsHostname},
		{"10.43.ab", "10.43.ab", len("10.43.ab"), host.IsHostname},
		{"10.43.1", "10.43.1", len("10.43.1"), host.IsHostname},
		{"10.43!", "10.43", len("10.43"), host.IsHostname},
		{"[1080:0:0:0:8:800:200C:417A]ab", "[1080::8:800:200c:417a]", len("[1080:0:0:0:8:800:200C:417A]"), host.IsIpv6},

		{"ab-c.com", "ab-c.com", len("ab-c.com"), host.IsHostname},
	}

	prefix := FuncName()

	for i, v := range testdata {
		host.Init()
		newPos, err := host.Parse(context, []byte(v.test), 0)

		if err != nil {
			t.Errorf("%s[%d] failed: %s\n", prefix, i, err.Error())
			continue
		}

		if newPos != v.newPos {
			t.Errorf("%s[%d] failed: newPos = %d, wanted = %d\n", prefix, i, newPos, len(v.wanted))
			continue
		}

		if !v.hosttype() {
			t.Errorf("%s[%d] failed: host is not ipv4\n", prefix, i)
			continue
		}

		str := host.String(context)

		if str != v.wanted {
			t.Errorf("%s[%d] failed: str = %s, wanted = %s\n", prefix, i, str, v.wanted)
			continue
		}
	}

}

func TestSipHostParseNOk(t *testing.T) {
	testdata := []struct {
		test string
	}{
		{""},
		{"!10.43.12.14"},
		{"[12"},
		{"[12!]"},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)

	for i, v := range testdata {
		addr := NewSipHost(context)
		host := addr.GetSipHost(context)
		_, err := host.Parse(context, []byte(v.test), 0)

		if err == nil {
			t.Errorf("TestSipHostParseNOk[%d] failed: should return err\n", i)
			continue
		}
	}
}

func TestSipHostPortParseOk(t *testing.T) {
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)

	addr := NewSipHostPort(context)
	host := addr.GetSipHostPort(context)

	testdata := []struct {
		test     string
		wanted   string
		newPos   int
		hosttype func() bool
		hasPort  bool
		port     uint16
	}{
		{"10.43.12.14", "10.43.12.14", len("10.43.12.14"), host.IsIpv4, false, 0},
		{"10.43.12.14!", "10.43.12.14", len("10.43.12.14"), host.IsIpv4, false, 0},
		{"10.43.12.14:5000", "10.43.12.14:5000", len("10.43.12.14:5000"), host.IsIpv4, true, 5000},
		{"10.43.1:65535", "10.43.1:65535", len("10.43.1:65535"), host.IsHostname, true, 65535},
		{"[1080:0:0:0:8:800:200C:417A]:0ab", "[1080::8:800:200c:417a]:0", len("[1080:0:0:0:8:800:200C:417A]:0"), host.IsIpv6, true, 0},
		{"[1080:0:0:0:8:800:200C:417A]", "[1080::8:800:200c:417a]", len("[1080:0:0:0:8:800:200C:417A]"), host.IsIpv6, false, 0},
		{"ab-c.com:123", "ab-c.com:123", len("ab-c.com:123"), host.IsHostname, true, 123},
	}

	prefix := FuncName()

	for i, v := range testdata {
		host.Init()
		newPos, err := host.Parse(context, []byte(v.test), 0)

		if err != nil {
			t.Errorf("%s[%d] failed: %s\n", prefix, i, err.Error())
			continue
		}

		if newPos != v.newPos {
			t.Errorf("%s[%d] failed: newPos = %d, wanted = %d\n", prefix, i, newPos, v.newPos)
			continue
		}

		if !v.hosttype() {
			t.Errorf("%s[%d] failed: host is not ipv4\n", prefix, i)
			continue
		}

		if v.hasPort != host.HasPort() {
			t.Errorf("%s[%d] failed: has port or not failed\n", prefix, i)
			continue
		}

		if v.hasPort {
			if v.port != host.GetPort() {
				t.Errorf("%s[%d] failed: has port or not failed\n", prefix, i)
				continue
			}
		}

		str := host.String(context)

		if str != v.wanted {
			t.Errorf("%s[%d] failed: str = %s, wanted = %s\n", prefix, i, str, v.wanted)
			continue
		}
	}

}

func TestSipHostPortParseNOk(t *testing.T) {
	testdata := []struct {
		test string
	}{
		{""},
		{"!10.43.12.14"},
		{"[12"},
		{"[12!]"},
		{"abd:"},
		{"abc:123456"},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)

	addr := NewSipHostPort(context)
	host := addr.GetSipHostPort(context)

	for i, v := range testdata {
		host.Init()
		_, err := host.Parse(context, []byte(v.test), 0)

		if err == nil {
			t.Errorf("TestSipHostPortParseNOk[%d] failed: should return err\n", i)
			continue
		}
	}
}

func TestWriteByteAsString(t *testing.T) {
	testdata := []struct {
		val byte
		ret string
	}{
		{0, "0"},
		{1, "1"},
		{9, "9"},
		{99, "99"},
		{109, "109"},
		{255, "255"},
	}

	//buf := bytes.NewBuffer(make([]byte, 1024*64))
	buf := &AbnfByteBuffer{}
	for i, v := range testdata {
		buf.Reset()
		WriteByteAsString(buf, v.val)
		if buf.String() != v.ret {
			t.Errorf("TestWriteByteAsString[%d] failed: ret = %v, wanted = %v\n", i, buf.String(), v.ret)
			continue
		}
	}
}

func TestSipHostWriteIpv4AsString(t *testing.T) {
	testdata := []struct {
		ipv4 []byte
		ret  string
	}{
		{[]byte{255, 254, 253, 252}, "255.254.253.252"},
		{[]byte{255, 0, 0, 252}, "255.0.0.252"},
		{[]byte{0, 0, 0, 252}, "0.0.0.252"},
	}

	//buf := bytes.NewBuffer(make([]byte, 1024*64))
	buf := &AbnfByteBuffer{}
	host := SipHost{id: HOST_TYPE_IPV4}
	for i, v := range testdata {
		buf.Reset()
		host.ip[0] = v.ipv4[0]
		host.ip[1] = v.ipv4[1]
		host.ip[2] = v.ipv4[2]
		host.ip[3] = v.ipv4[3]
		host.WriteIpv4AsString(buf)
		if buf.String() != v.ret {
			t.Errorf("TestSipHostWriteIpv4AsString[%d] failed: ret = %v, wanted = %v\n", i, buf.String(), v.ret)
			continue
		}
	}
}

func BenchmarkWriteByteAsString1(b *testing.B) {
	b.StopTimer()
	//buf := bytes.NewBuffer(make([]byte, 1024*64))
	buf := &AbnfByteBuffer{}
	b.SetBytes(2)
	b.ReportAllocs()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		buf.Reset()
		WriteByteAsString(buf, 255)
	}
}

func BenchmarkWriteByteAsString2(b *testing.B) {
	b.StopTimer()
	//buf := bytes.NewBuffer(make([]byte, 1024*64))
	buf := &AbnfByteBuffer{}
	b.SetBytes(2)
	b.ReportAllocs()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		buf.Reset()
		WriteByteAsString(buf, 0)
	}
}

func BenchmarkWriteByteUseFmt1(b *testing.B) {
	b.StopTimer()
	//buf := bytes.NewBuffer(make([]byte, 1024*64))
	buf := &AbnfByteBuffer{}
	b.SetBytes(2)
	b.ReportAllocs()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		buf.Reset()
		buf.WriteString(fmt.Sprintf("%d", 255))
	}
}

func BenchmarkWriteByteUseFmt2(b *testing.B) {
	b.StopTimer()
	//buf := bytes.NewBuffer(make([]byte, 1024*64))
	buf := &AbnfByteBuffer{}
	b.SetBytes(2)
	b.ReportAllocs()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		buf.Reset()
		buf.WriteString(fmt.Sprintf("%d", 0))
	}
}

func BenchmarkWriteIpv4String(b *testing.B) {
	b.StopTimer()
	ip := []byte{255, 255, 255, 255}
	buf := bytes.NewBuffer(make([]byte, 1024*64))
	b.SetBytes(2)
	b.ReportAllocs()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		buf.Reset()
		buf.WriteString(net.IP(ip).String())
	}
}

func BenchmarkWriteIpv4UseFmt(b *testing.B) {
	b.StopTimer()
	ip := []byte{255, 255, 255, 255}
	//buf := bytes.NewBuffer(make([]byte, 1024*64))
	buf := &AbnfByteBuffer{}
	b.SetBytes(2)
	b.ReportAllocs()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		buf.Reset()
		buf.WriteString(fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3]))
	}
}

func BenchmarkWriteIpv4AsString(b *testing.B) {
	b.StopTimer()
	host := SipHost{id: HOST_TYPE_IPV4}
	host.ip[0] = 255
	host.ip[1] = 255
	host.ip[2] = 255
	host.ip[3] = 255
	//buf := bytes.NewBuffer(make([]byte, 1024*64))
	buf := &AbnfByteBuffer{}
	b.SetBytes(2)
	b.ReportAllocs()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		buf.Reset()
		host.WriteIpv4AsString(buf)
	}
}

func BenchmarkSipHostParseIpv4(b *testing.B) {
	b.StopTimer()
	host := SipHost{}
	src := []byte("255.255.255.255")
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	b.SetBytes(2)
	b.ReportAllocs()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		host.parseIpv4(context, src, 0)
	}
}
