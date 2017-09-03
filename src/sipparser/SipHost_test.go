package sipparser

import (
	//"bytes"
	//"fmt"
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

	ipv4 := net.IPv4(10, 1, 1, 1)

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
