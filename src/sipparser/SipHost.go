package sipparser

import (
	"bytes"
	//"fmt"
	"net"
	"unsafe"
)

const (
	HOST_TYPE_UNKNOWN = 0
	HOST_TYPE_IPV4    = 1
	HOST_TYPE_IPV6    = 2
	HOST_TYPE_NAME    = 3
)

type SipHost struct {
	id   byte
	ip   [16]byte
	data AbnfBuf
}

func NewSipHost(context *ParseContext) AbnfPtr {
	addr := context.allocator.Alloc(uint32(unsafe.Sizeof(SipHost{})))
	if addr == ABNF_PTR_NIL {
		return ABNF_PTR_NIL
	}
	addr.GetSipHost(context).Init()
	return addr
}

func (this *SipHost) Init() {
	this.id = HOST_TYPE_UNKNOWN
	this.data.Init()
}

func (this *SipHost) Encode(context *ParseContext, buf *AbnfByteBuffer) {
	if this.id == HOST_TYPE_UNKNOWN {
		buf.WriteString("unknown host")
	} else if this.IsIpv4() {
		//buf.WriteString(this.GetIpString(context))
		this.WriteIpv4AsString(buf)
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

func (this *SipHost) SetIpv4(context *ParseContext, ip []byte) {
	this.id = HOST_TYPE_IPV4
	//this.data.SetByteSlice(context, ip)
	this.ip[0] = ip[0]
	this.ip[1] = ip[1]
	this.ip[2] = ip[2]
	this.ip[3] = ip[3]
}

func (this *SipHost) SetIpv6(context *ParseContext, ip net.IP) {
	this.id = HOST_TYPE_IPV6
	//this.data.SetByteSlice(context, ip)
	this.ip[0] = ip[0]
	this.ip[1] = ip[1]
	this.ip[2] = ip[2]
	this.ip[3] = ip[3]
	this.ip[4] = ip[4]
	this.ip[5] = ip[5]
	this.ip[6] = ip[6]
	this.ip[7] = ip[7]
	this.ip[8] = ip[8]
	this.ip[9] = ip[9]
	this.ip[10] = ip[10]
	this.ip[11] = ip[11]
	this.ip[12] = ip[12]
	this.ip[13] = ip[13]
	this.ip[14] = ip[14]
	this.ip[15] = ip[15]
}

func (this *SipHost) SetHostname(context *ParseContext, hostname []byte) {
	this.id = HOST_TYPE_NAME
	this.data.SetByteSlice(context, hostname)
}

//func (this *SipHost) GetIp(context *ParseContext) net.IP { return this.data.GetAsByteSlice(context) }
func (this *SipHost) GetIp(context *ParseContext) net.IP { return this.ip[:] }

/*func (this *SipHost) GetIpString(context *ParseContext) string {
	return net.IP(this.data.GetAsByteSlice(context)).String()
}*/
func (this *SipHost) GetIpString(context *ParseContext) string {
	if this.id == HOST_TYPE_IPV4 {
		return net.IP(this.ip[0:4]).String()
	}
	return net.IP(this.ip[0:]).String()
}

func (this *SipHost) WriteIpv4AsString(buf *AbnfByteBuffer) {
	WriteByteAsString(buf, this.ip[0])
	buf.WriteByte('.')
	WriteByteAsString(buf, this.ip[1])
	buf.WriteByte('.')
	WriteByteAsString(buf, this.ip[2])
	buf.WriteByte('.')
	WriteByteAsString(buf, this.ip[3])
}

func WriteByteAsString(buf *AbnfByteBuffer, v byte) {
	buf.WriteString(g_byteAsString_table[v])
	/*
		if v == 0 {
			buf.WriteByte('0')
			return
		}

		x2 := v / 10
		x1 := v - x2*10
		x3 := x2 / 10
		x2 -= x3 * 10

		if x3 != 0 {
			buf.WriteByte('0' + x3)
			buf.WriteByte('0' + x2)
		} else if x2 != 0 {
			buf.WriteByte('0' + x2)
		}

		buf.WriteByte('0' + x1) //*/

}

func (this *SipHost) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	return this.ParseWithoutInit(context, src, pos)
}

func (this *SipHost) ParseWithoutInit(context *ParseContext, src []byte, pos int) (newPos int, err error) {
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
	if this.id == HOST_TYPE_IPV4 {
		return bytes.Equal(this.ip[0:4], rhs.ip[0:4])
	}
	return bytes.Equal(this.ip[0:], rhs.ip[0:])
	//return this.data.Equal(context, &rhs.data)

}

func (this *SipHost) parseIpv6(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos
	len1 := len(src)
	for ; newPos < len1; newPos++ {
		if src[newPos] == ']' {
			break
		}
	}

	if newPos >= len1 {
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
	len1 := len(src)
	newPos = pos

	if pos >= len1 {
		return newPos, false
	}

	var digit int

	digit, _, newPos, ok = ParseUInt(src, newPos)

	if !ok || digit > 0xff {
		return newPos, false
	}

	this.ip[0] = byte(digit)

	for num := 1; num < net.IPv4len; num++ {
		if newPos >= len1 {
			return newPos, false
		}

		if src[newPos] != '.' {
			return newPos, false
		}
		newPos++

		digit, _, newPos, ok = ParseUInt(src, newPos)

		if !ok || digit > 0xff {
			return newPos, false
		}

		this.ip[num] = byte(digit)
	}

	if newPos < len1 && IsHostname(src[newPos]) {
		return newPos, false
	}

	this.id = HOST_TYPE_IPV4

	return newPos, true
}

func (this *SipHost) parseHostname(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	len1 := len(src)
	for newPos = pos; newPos < len1 && IsHostname(src[newPos]); newPos++ {
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

func NewSipHostPort(context *ParseContext) AbnfPtr {
	addr := context.allocator.Alloc(uint32(unsafe.Sizeof(SipHostPort{})))
	if addr == ABNF_PTR_NIL {
		return ABNF_PTR_NIL
	}
	addr.GetSipHostPort(context).Init()
	return addr
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

func (this *SipHostPort) Encode(context *ParseContext, buf *AbnfByteBuffer) {
	this.SipHost.Encode(context, buf)
	if this.hasPort {
		buf.WriteByte(':')
		EncodeUInt(buf, uint64(this.port))
		//buf.WriteString(strconv.FormatUint(uint64(this.port), 10))
		//buf.WriteString(fmt.Sprintf(":%d", this.port))
	}
}

func (this *SipHostPort) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

func (this *SipHostPort) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	return this.ParseWithoutInit(context, src, pos)
}

func (this *SipHostPort) ParseWithoutInit(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos, err = this.SipHost.ParseWithoutInit(context, src, pos)
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
