package sipparser

import (
	//"fmt"
	"bytes"
	"net"
	"strconv"
)

const (
	HOST_TYPE_UNKNOWN = 0
	HOST_TYPE_IPV4    = 1
	HOST_TYPE_IPV6    = 2
	HOST_TYPE_NAME    = 3
)

type SipHost struct {
	id   byte
	data []byte
}

func (this *SipHost) Init() {
	this.id = HOST_TYPE_UNKNOWN
}

func (this *SipHost) String() string {
	if this.id == HOST_TYPE_UNKNOWN {
		return "unknown host"
	}
	if this.IsIpv4() {
		return this.GetIpString()
	}

	if this.IsIpv6() {
		return "[" + this.GetIpString() + "]"
	}
	return string(this.data)
}

func (this *SipHost) IsIpv4() bool     { return this.id == HOST_TYPE_IPV4 }
func (this *SipHost) IsIpv6() bool     { return this.id == HOST_TYPE_IPV6 }
func (this *SipHost) IsHostname() bool { return this.id == HOST_TYPE_NAME }

func (this *SipHost) SetIpv4(ip net.IP) {
	this.id = HOST_TYPE_IPV4
	this.data = ip
}

func (this *SipHost) SetIpv6(ip net.IP) {
	this.id = HOST_TYPE_IPV6
	this.data = ip
}

func (this *SipHost) SetHostname(hostname []byte) {
	this.id = HOST_TYPE_NAME
	this.data = hostname
}

func (this *SipHost) GetIp() net.IP       { return this.data }
func (this *SipHost) GetIpString() string { return net.IP(this.data).String() }

func (this *SipHost) Parse(src []byte, pos int) (newPos int, err error) {
	if pos >= len(src) {
		return pos, &AbnfError{"parse failed", src, newPos}
	}

	newPos = pos

	if src[pos] == '[' {
		return this.parseIpv6(src, pos+1)
	}

	if IsAlpha(src[pos]) {
		return this.parseHostname(src, pos)
	}

	var ok bool

	newPos, ok = this.parseIpv4(src, pos)

	if !ok {
		return this.parseHostname(src, pos)
	}
	return newPos, nil
}

func (this *SipHost) Equal(rhs *SipHost) bool {
	if this.id != rhs.id {
		return false
	}
	if this.IsHostname() {
		return EqualNoCase(this.data, rhs.data)
	}
	return bytes.Equal(this.data, rhs.data)

}

func (this *SipHost) parseIpv6(src []byte, pos int) (newPos int, err error) {
	newPos = pos
	for ; newPos < len(src); newPos++ {
		if src[newPos] == ']' {
			break
		}
	}

	if newPos >= len(src) {
		return newPos, &AbnfError{"no \"]\" for ipv6-reference", src, newPos}
	}

	ipv6 := net.ParseIP(string(src[pos:newPos]))
	if ipv6 == nil {
		return newPos, &AbnfError{"parse ipv6 failed", src, newPos}
	}
	this.SetIpv6(ipv6)
	return newPos + 1, nil
}

func (this *SipHost) parseIpv4(src []byte, pos int) (newPos int, ok bool) {

	var ipv4 [net.IPv4len]byte

	newPos = pos

	for num := 0; num < net.IPv4len; num++ {
		if newPos >= len(src) {
			return newPos, false
		}

		if num > 0 {
			if src[newPos] != '.' {
				return newPos, false
			}
			newPos++
		}

		var digit int
		var ok bool

		digit, newPos, ok = ParseUInt(src, newPos)

		if !ok || digit > 0xff {
			return newPos, false
		}

		ipv4[num] = byte(digit)
	}

	if newPos < len(src) && IsHostname(src[newPos]) {
		return newPos, false
	}

	this.SetIpv4(net.IPv4(ipv4[0], ipv4[1], ipv4[2], ipv4[3]))

	return newPos, true
}

func (this *SipHost) parseHostname(src []byte, pos int) (newPos int, err error) {
	for newPos = pos; newPos < len(src) && IsHostname(src[newPos]); newPos++ {
	}

	if newPos <= pos {
		return newPos, &AbnfError{"null hostname", src, newPos}
	}
	this.SetHostname(src[pos:newPos])
	return newPos, nil

}

type SipHostPort struct {
	SipHost

	hasPort bool
	port    uint16
}

func (this *SipHostPort) Init() {
	this.SipHost.Init()
	this.hasPort = false
}

func (this *SipHostPort) HasPort() bool   { return this.hasPort }
func (this *SipHostPort) GetPort() uint16 { return this.port }

func (this *SipHostPort) SetPort(port uint16) {
	this.hasPort = true
	this.port = port
}

func (this *SipHostPort) String() string {
	str := this.SipHost.String()
	if this.hasPort {
		str += ":"
		str += strconv.FormatUint(uint64(this.port), 10)
	}

	return str
}

func (this *SipHostPort) Parse(src []byte, pos int) (newPos int, err error) {
	newPos, err = this.SipHost.Parse(src, pos)
	if err != nil {
		return newPos, err
	}

	if newPos >= len(src) {
		return newPos, nil
	}

	if src[newPos] != ':' {
		return newPos, nil
	}

	newPos++

	var digit int
	var ok bool

	digit, newPos, ok = ParseUInt(src, newPos)
	if !ok {
		return newPos, &AbnfError{"parse port failed after \":\"", src, newPos}
	}

	if digit < 0 || digit > 0xffff {
		return newPos, &AbnfError{"port wrong range \":\"", src, newPos}
	}

	this.SetPort(uint16(digit))

	return newPos, nil
}

func (this *SipHostPort) Equal(rhs *SipHostPort) bool {
	if (this.hasPort && !rhs.hasPort) || (!this.hasPort && rhs.hasPort) {
		return false
	}

	return this.SipHost.Equal(&rhs.SipHost)
}
