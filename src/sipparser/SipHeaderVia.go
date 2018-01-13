package sipparser

import (
	//"bytes"
	//"fmt"
	"unsafe"
)

type SipHeaderVia struct {
	protocolName    AbnfBuf
	protocolVersion AbnfBuf
	transport       AbnfBuf
	sentBy          SipHostPort
	params          SipGenericParams
}

func NewSipHeaderVia(context *ParseContext) AbnfPtr {
	addr := context.allocator.Alloc(uint32(unsafe.Sizeof(SipHeaderVia{})))
	if addr == ABNF_PTR_NIL {
		return ABNF_PTR_NIL
	}
	addr.GetSipHeaderVia(context).Init()
	return addr
}

func (this *SipHeaderVia) Init() {
	this.protocolName.Init()
	this.protocolVersion.Init()
	this.transport.Init()
	this.sentBy.Init()
	this.params.Init()
}

func (this *SipHeaderVia) AllowMulti() bool { return true }
func (this *SipHeaderVia) HasValue() bool   { return true }

/* RFC3261
 *
 * Via               =  ( "Via" / "v" ) HCOLON via-parm *(COMMA via-parm)
 * via-parm          =  sent-protocol LWS sent-by *( SEMI via-params )
 * via-params        =  via-ttl / via-maddr
 *                      / via-received / via-branch
 *                      / via-extension
 * via-ttl           =  "ttl" EQUAL ttl
 * via-maddr         =  "maddr" EQUAL host
 * via-received      =  "received" EQUAL (IPv4address / IPv6address)
 * via-branch        =  "branch" EQUAL token
 * via-extension     =  generic-param
 * sent-protocol     =  protocol-name SLASH protocol-version
 *                      SLASH transport
 * protocol-name     =  "SIP" / token
 * protocol-version  =  token
 * transport         =  "UDP" / "TCP" / "TLS" / "SCTP"
 *                      / other-transport
 * sent-by           =  host [ COLON port ]
 * ttl               =  1*3DIGIT ; 0 to 255
 *
 * RFC3581
 *
 * response-port     = "rport" [EQUAL 1*DIGIT]
 * via-params        =  via-ttl / via-maddr
 *                      / via-received / via-branch
 *                      / response-port / via-extension
 *
 */
func (this *SipHeaderVia) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	return this.ParseWithoutInit(context, src, pos)
}

func (this *SipHeaderVia) ParseWithoutInit(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	name, newPos, err := ParseHeaderName(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if !EqualNoCase(src[name.Begin:name.End], StringToByteSlice(ABNF_NAME_SIP_HDR_VIA)) &&
		!EqualNoCase(src[name.Begin:name.End], StringToByteSlice(ABNF_NAME_SIP_HDR_VIA_S)) {
		return newPos, &AbnfError{"Via parse: wrong header-name", src, newPos}
	}

	return this.ParseValueWithoutInit(context, src, newPos)
}

func (this *SipHeaderVia) ParseValue(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	return this.ParseValueWithoutInit(context, src, pos)
}

func (this *SipHeaderVia) ParseValueWithoutInit(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos
	newPos, err = this.protocolName.ParseSipToken(context, src, newPos)
	if err != nil {
		return newPos, err
	}

	newPos, err = ParseSWSMark(src, newPos, '/')
	if err != nil {
		return newPos, err
	}

	newPos, err = this.protocolVersion.ParseSipToken(context, src, newPos)
	if err != nil {
		return newPos, err
	}

	newPos, err = ParseSWSMark(src, newPos, '/')
	if err != nil {
		return newPos, err
	}

	newPos, err = this.transport.ParseSipToken(context, src, newPos)
	if err != nil {
		return newPos, err
	}

	var ok bool
	newPos, ok = ParseLWS(src, newPos)
	if !ok {
		return newPos, &AbnfError{"Via parse: wrong LWS", src, newPos}
	}

	newPos, err = this.sentBy.ParseWithoutInit(context, src, newPos)
	if err != nil {
		return newPos, err
	}

	return this.params.ParseWithoutInit(context, src, newPos, ';')
}

func (this *SipHeaderVia) Encode(context *ParseContext, buf *AbnfByteBuffer) {
	buf.WriteString(ABNF_NAME_SIP_HDR_VIA_COLON)
	this.EncodeValue(context, buf)
}

func (this *SipHeaderVia) EncodeValue(context *ParseContext, buf *AbnfByteBuffer) {
	this.protocolName.Encode(context, buf)
	buf.WriteByte('/')
	this.protocolVersion.Encode(context, buf)
	buf.WriteByte('/')
	this.transport.Encode(context, buf)
	buf.WriteByte(' ')
	this.sentBy.Encode(context, buf)
	this.params.Encode(context, buf, ';')
}

func (this *SipHeaderVia) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

func (this *SipHeaderVia) StringValue(context *ParseContext) string {
	var buf AbnfByteBuffer
	this.EncodeValue(context, &buf)
	return buf.String()
}

func ParseSipVia(context *ParseContext, src []byte, pos int) (newPos int, parsed AbnfPtr, err error) {
	addr := NewSipHeaderVia(context)
	if addr == ABNF_PTR_NIL {
		return newPos, ABNF_PTR_NIL, &AbnfError{"Via parse: out of memory for new header", src, newPos}
	}
	newPos, err = addr.GetSipHeaderVia(context).ParseValueWithoutInit(context, src, pos)
	return newPos, addr, err
}

func EncodeSipViaValue(parsed AbnfPtr, context *ParseContext, buf *AbnfByteBuffer) {
	if parsed == ABNF_PTR_NIL {
		return
	}
	parsed.GetSipHeaderVia(context).EncodeValue(context, buf)
}
