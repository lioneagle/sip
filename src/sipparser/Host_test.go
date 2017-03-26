package sipparser

import (
	//"bytes"
	//"fmt"
	"net"
	"testing"
)

func TestSipHostUnknownString(t *testing.T) {
	var host SipHost

	host.Init()

	str := host.String()

	if str != "unknown host" {
		t.Errorf("TestSipHostUnknownString failed, str = %s, wanted = %s\n", str, "unknown host")
	}
}

func TestSipHostIpv4String(t *testing.T) {
	var host SipHost
	ipv4 := net.IPv4(10, 1, 1, 1)

	host.SetIpv4(ipv4)

	str := host.String()

	if str != "10.1.1.1" {
		t.Errorf("TestSipHostIpv4String failed, str = %s, wanted = %s\n", str, "10.1.1.1")
	}
}

func TestSipHostIpv6String(t *testing.T) {
	var host SipHost
	ipv4 := net.IPv4(10, 1, 1, 1)

	host.SetIpv6(ipv4.To16())

	str := host.String()

	if str != "[10.1.1.1]" {
		t.Errorf("TestSipHostIpv6String failed, str = %s, wanted = %s\n", str, "[10.1.1.1]")
	}
}

func TestSipHostHostnameString(t *testing.T) {
	var host SipHost

	host.SetHostname([]byte("abc.com"))

	if !host.IsHostname() {
		t.Errorf("TestSipHostHostnameString failed, host is not hostname\n")
	}

	str := host.String()

	if str != "abc.com" {
		t.Errorf("TestSipHostHostnameString failed, str = %s, wanted = %s\n", str, "abc.com")
	}
}

func TestSipHostParseOk(t *testing.T) {
	var host SipHost

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

	for i, v := range testdata {
		host.Init()
		newPos, err := host.Parse([]byte(v.test), 0)

		if err != nil {
			t.Errorf("TestSipHostParseOk[%d] failed, %s\n", i, err.Error())
			continue
		}

		if newPos != v.newPos {
			t.Errorf("TestSipHostParseOk[%d] failed, newPos = %d, wanted = %d\n", i, newPos, len(v.wanted))
			continue
		}

		if !v.hosttype() {
			t.Errorf("TestSipHostParseOk[%d] failed, host is not ipv4\n", i)
			continue
		}

		str := host.String()

		if str != v.wanted {
			t.Errorf("TestSipHostParseOk[%d] failed, str = %s, wanted = %s\n", i, str, v.wanted)
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

	var host SipHost

	for i, v := range testdata {
		host.Init()
		_, err := host.Parse([]byte(v.test), 0)

		if err == nil {
			t.Errorf("TestSipHostParseNOk[%d] failed, should return err\n", i)
			continue
		}
	}
}

func TestSipHostPortParseOk(t *testing.T) {
	var host SipHostPort

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

	for i, v := range testdata {
		host.Init()
		newPos, err := host.Parse([]byte(v.test), 0)

		if err != nil {
			t.Errorf("TestSipHostPortParseOk[%d] failed, %s\n", i, err.Error())
			continue
		}

		if newPos != v.newPos {
			t.Errorf("TestSipHostPortParseOk[%d] failed, newPos = %d, wanted = %d\n", i, newPos, v.newPos)
			continue
		}

		if !v.hosttype() {
			t.Errorf("TestSipHostPortParseOk[%d] failed, host is not ipv4\n", i)
			continue
		}

		if v.hasPort != host.HasPort() {
			t.Errorf("TestSipHostPortParseOk[%d] failed, has port or not failed\n", i)
			continue
		}

		if v.hasPort {
			if v.port != host.GetPort() {
				t.Errorf("TestSipHostPortParseOk[%d] failed, has port or not failed\n", i)
				continue
			}
		}

		str := host.String()

		if str != v.wanted {
			t.Errorf("TestSipHostPortParseOk[%d] failed, str = %s, wanted = %s\n", i, str, v.wanted)
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

	var host SipHostPort

	for i, v := range testdata {
		host.Init()
		_, err := host.Parse([]byte(v.test), 0)

		if err == nil {
			t.Errorf("TestSipHostPortParseNOk[%d] failed, should return err\n", i)
			continue
		}
	}
}