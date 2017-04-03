package sipparser

import (
	"bytes"
	//"fmt"
	"net"
	"strconv"
	"unsafe"
)

const (
	HOST_TYPE_UNKNOWN = 0
	HOST_TYPE_IPV4    = 1
	HOST_TYPE_IPV6    = 2
	HOST_TYPE_NAME    = 3
)

type SipHost struct {
	id byte
	//data []byte
	data AbnfBuf
}

func NewSipHost(context *ParseContext) (*SipHost, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipHost{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*SipHost)(unsafe.Pointer(mem)).Init()
	return (*SipHost)(unsafe.Pointer(mem)), addr
}

func (this *SipHost) Init() {
	this.id = HOST_TYPE_UNKNOWN
	this.data.Init()
}

func (this *SipHost) Encode(context *ParseContext, buf *bytes.Buffer) {
	if this.id == HOST_TYPE_UNKNOWN {
		buf.WriteString("unknown host")
	} else if this.IsIpv4() {
		buf.WriteString(this.GetIpString(context))
	} else if this.IsIpv6() {
		buf.WriteByte('[')
		buf.WriteString(this.GetIpString(context))
		buf.WriteByte(']')
	} else {
		this.data.Encode(context, buf)
	}
}

func (this *SipHost) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

func (this *SipHost) IsIpv4() bool     { return this.id == HOST_TYPE_IPV4 }
func (this *SipHost) IsIpv6() bool     { return this.id == HOST_TYPE_IPV6 }
func (this *SipHost) IsHostname() bool { return this.id == HOST_TYPE_NAME }

func (this *SipHost) SetIpv4(context *ParseContext, ip net.IP) {
	this.id = HOST_TYPE_IPV4
	this.data.SetByteSlice(context, ip)
}

func (this *SipHost) SetIpv6(context *ParseContext, ip net.IP) {
	this.id = HOST_TYPE_IPV6
	this.data.SetByteSlice(context, ip)
}

func (this *SipHost) SetHostname(context *ParseContext, hostname []byte) {
	this.id = HOST_TYPE_NAME
	this.data.SetByteSlice(context, hostname)
}

func (this *SipHost) GetIp(context *ParseContext) net.IP { return this.data.GetAsByteSlice(context) }
func (this *SipHost) GetIpString(context *ParseContext) string {
	return net.IP(this.data.GetAsByteSlice(context)).String()
}

func (this *SipHost) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	if pos >= len(src) {
		return pos, &AbnfError{"host parse: parse failed", src, newPos}
	}

	newPos = pos

	if src[pos] == '[' {
		return this.parseIpv6(context, src, pos+1)
	}

	if IsAlpha(src[pos]) {
		return this.parseHostname(context, src, pos)
	}

	var ok bool

	newPos, ok = this.parseIpv4(context, src, pos)

	if !ok {
		return this.parseHostname(context, src, pos)
	}
	return newPos, nil
}

func (this *SipHost) Equal(context *ParseContext, rhs *SipHost) bool {
	if this.id != rhs.id {
		return false
	}
	if this.IsHostname() {
		return this.data.EqualNoCase(context, &rhs.data)
	}
	return this.data.Equal(context, &rhs.data)

}

func (this *SipHost) parseIpv6(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos
	for ; newPos < len(src); newPos++ {
		if src[newPos] == ']' {
			break
		}
	}

	if newPos >= len(src) {
		return newPos, &AbnfError{"host parse: no \"]\" for ipv6-reference", src, newPos}
	}

	ipv6 := net.ParseIP(ByteSliceToString(src[pos:newPos]))
	if ipv6 == nil {
		return newPos, &AbnfError{"host parse: parse ipv6 failed", src, newPos}
	}
	this.SetIpv6(context, ipv6)
	return newPos + 1, nil
}

func (this *SipHost) parseIpv4(context *ParseContext, src []byte, pos int) (newPos int, ok bool) {

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

		digit, _, newPos, ok = ParseUInt(src, newPos)

		if !ok || digit > 0xff {
			return newPos, false
		}

		ipv4[num] = byte(digit)
	}

	if newPos < len(src) && IsHostname(src[newPos]) {
		return newPos, false
	}

	this.SetIpv4(context, net.IPv4(ipv4[0], ipv4[1], ipv4[2], ipv4[3]))

	return newPos, true
}

func (this *SipHost) parseHostname(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	for newPos = pos; newPos < len(src) && IsHostname(src[newPos]); newPos++ {
	}

	if newPos <= pos {
		return newPos, &AbnfError{"host parse: null hostname", src, newPos}
	}
	this.SetHostname(context, src[pos:newPos])
	return newPos, nil

}

type SipHostPort struct {
	SipHost

	hasPort bool
	port    uint16
}

func NewSipHostPort(context *ParseContext) (*SipHostPort, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipHostPort{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*SipHostPort)(unsafe.Pointer(mem)).Init()
	return (*SipHostPort)(unsafe.Pointer(mem)), addr
}

func (this *SipHostPort) Init() {
	this.SipHost.Init()
	this.hasPort = false
	this.port = 0
}

func (this *SipHostPort) HasPort() bool   { return this.hasPort }
func (this *SipHostPort) GetPort() uint16 { return this.port }

func (this *SipHostPort) SetPort(port uint16) {
	this.hasPort = true
	this.port = port
}

func (this *SipHostPort) Encode(context *ParseContext, buf *bytes.Buffer) {
	this.SipHost.Encode(context, buf)
	if this.hasPort {
		buf.WriteByte(':')
		buf.WriteString(strconv.FormatUint(uint64(this.port), 10))
		//buf.WriteString(fmt.Sprintf(":%d", this.port))
	}
}

func (this *SipHostPort) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

func (this *SipHostPort) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	newPos, err = this.SipHost.Parse(context, src, pos)
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

	digit, _, newPos, ok = ParseUInt(src, newPos)
	if !ok {
		return newPos, &AbnfError{"hostport parse: parse port failed after \":\"", src, newPos}
	}

	if digit < 0 || digit > 0xffff {
		return newPos, &AbnfError{"hostport parse: port wrong range \":\"", src, newPos}
	}

	this.SetPort(uint16(digit))

	return newPos, nil
}

func (this *SipHostPort) Equal(context *ParseContext, rhs *SipHostPort) bool {
	if (this.hasPort && !rhs.hasPort) || (!this.hasPort && rhs.hasPort) {
		return false
	}

	return this.SipHost.Equal(context, &rhs.SipHost)
}
