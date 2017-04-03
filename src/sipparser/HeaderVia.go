package sipparser

import (
	"bytes"
	//"fmt"
	"unsafe"
)

type SipHeaderVia struct {
	version   SipVersion
	transport AbnfToken
	sentBy    SipHostPort
	params    SipGenericParams
}

func NewSipHeaderVia(context *ParseContext) (*SipHeaderVia, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipHeaderVia{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*SipHeaderVia)(unsafe.Pointer(mem)).Init()
	return (*SipHeaderVia)(unsafe.Pointer(mem)), addr
}

func (this *SipHeaderVia) Init() {
	this.version.Init()
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
	name, newPos, err := ParseHeaderName(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if !EqualNoCase(src[name.Begin:name.End], StringToByteSlice("Via")) && !EqualNoCase(src[name.Begin:name.End], StringToByteSlice("v")) {
		return newPos, &AbnfError{"Via parse: wrong header-name", src, newPos}
	}

	return this.ParseValue(context, src, newPos)
}

func (this *SipHeaderVia) ParseValue(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	newPos = pos
	newPos, err = this.version.Parse(context, src, newPos)
	if err != nil {
		return newPos, err
	}

	if src[newPos] != '/' {
		return newPos, &AbnfError{"Via parse: no slash after protocol-version", src, newPos}
	}

	newPos++
	newPos, err = this.transport.Parse(context, src, newPos, IsSipToken)
	if err != nil {
		return newPos, err
	}

	newPos, err = ParseLWS(src, newPos)
	if err != nil {
		return newPos, err
	}

	newPos, err = this.sentBy.Parse(context, src, newPos)
	if err != nil {
		return newPos, err
	}

	return this.params.Parse(context, src, newPos, ';')
}

func (this *SipHeaderVia) Encode(context *ParseContext, buf *bytes.Buffer) {
	buf.WriteString("Via: ")
	this.EncodeValue(context, buf)
}

func (this *SipHeaderVia) EncodeValue(context *ParseContext, buf *bytes.Buffer) {
	this.version.Encode(context, buf)
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
	var buf bytes.Buffer
	this.EncodeValue(context, &buf)
	return buf.String()
}

func ParseSipVia(context *ParseContext, src []byte, pos int) (newPos int, parsed AbnfPtr, err error) {
	//fmt.Println("enter via")
	header, addr := NewSipHeaderVia(context)
	if header == nil || addr == ABNF_PTR_NIL {
		return newPos, ABNF_PTR_NIL, &AbnfError{"Via parse: out of memory for new header", src, newPos}
	}
	newPos, err = header.ParseValue(context, src, pos)
	return newPos, addr, err
}
