package sipparser

import (
	"bytes"
	//"fmt"
	"unsafe"
)

type SipHeaderFrom struct {
	addr   SipAddr
	params SipGenericParams
}

func NewSipHeaderFrom(context *ParseContext) AbnfPtr {
	addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipHeaderFrom{})))
	if addr == ABNF_PTR_NIL {
		return ABNF_PTR_NIL
	}
	addr.GetSipHeaderFrom(context).Init()
	return addr
}

func (this *SipHeaderFrom) Init() {
	this.addr.Init()
	this.params.Init()
}

func (this *SipHeaderFrom) AllowMulti() bool { return false }
func (this *SipHeaderFrom) HasValue() bool   { return true }

/* RFC3261
 *
 * From        =  ( "From" / "f" ) HCOLON from-spec
 * from-spec   =  ( name-addr / addr-spec )
 *                *( SEMI from-param )
 */
func (this *SipHeaderFrom) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	return this.ParseWithoutInit(context, src, pos)
}

func (this *SipHeaderFrom) ParseWithoutInit(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	name, newPos, err := ParseHeaderName(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if !EqualNoCase(src[name.Begin:name.End], StringToByteSlice(ABNF_NAME_SIP_HDR_FROM)) &&
		!EqualNoCase(src[name.Begin:name.End], StringToByteSlice(ABNF_NAME_SIP_HDR_FROM_S)) {
		return newPos, &AbnfError{"From parse: wrong header-name", src, newPos}
	}

	return this.ParseValueWithoutInit(context, src, newPos)
}

func (this *SipHeaderFrom) ParseValue(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	return this.ParseValueWithoutInit(context, src, pos)
}

func (this *SipHeaderFrom) ParseValueWithoutInit(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos, err = this.addr.Parse(context, src, pos)
	if err != nil {
		return newPos, err
	}

	return this.params.ParseWithoutInit(context, src, newPos, ';')
}

func (this *SipHeaderFrom) Encode(context *ParseContext, buf *bytes.Buffer) {
	buf.WriteString(ABNF_NAME_SIP_HDR_FROM_COLON)
	this.EncodeValue(context, buf)
}

func (this *SipHeaderFrom) EncodeValue(context *ParseContext, buf *bytes.Buffer) {
	this.addr.Encode(context, buf)
	this.params.Encode(context, buf, ';')
}

func (this *SipHeaderFrom) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

func ParseSipFrom(context *ParseContext, src []byte, pos int) (newPos int, parsed AbnfPtr, err error) {
	addr := NewSipHeaderFrom(context)
	if addr == ABNF_PTR_NIL {
		return newPos, ABNF_PTR_NIL, &AbnfError{"From parse: out of memory for new header", src, newPos}
	}
	newPos, err = addr.GetSipHeaderFrom(context).ParseValueWithoutInit(context, src, pos)
	return newPos, addr, err
}

func EncodeSipFromValue(parsed AbnfPtr, context *ParseContext, buf *bytes.Buffer) {
	if parsed == ABNF_PTR_NIL {
		return
	}
	parsed.GetSipHeaderFrom(context).EncodeValue(context, buf)
}
