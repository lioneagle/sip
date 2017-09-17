package sipparser

import (
	"bytes"
	//"fmt"
	"unsafe"
)

type SipHeaderTo struct {
	addr   SipAddr
	params SipGenericParams
}

func NewSipHeaderTo(context *ParseContext) AbnfPtr {
	addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipHeaderTo{})))
	if addr == ABNF_PTR_NIL {
		return ABNF_PTR_NIL
	}
	addr.GetSipHeaderTo(context).Init()
	return addr
}

func (this *SipHeaderTo) Init() {
	this.addr.Init()
	this.params.Init()
}

func (this *SipHeaderTo) AllowMulti() bool { return false }
func (this *SipHeaderTo) HasValue() bool   { return true }

/* RFC3261
 *
 * From        =  ( "From" / "f" ) HCOLON from-spec
 * from-spec   =  ( name-addr / addr-spec )
 *                *( SEMI from-param )
 */
func (this *SipHeaderTo) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	name, newPos, err := ParseHeaderName(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if !EqualNoCase(src[name.Begin:name.End], StringToByteSlice(ABNF_NAME_SIP_HDR_TO)) &&
		!EqualNoCase(src[name.Begin:name.End], StringToByteSlice(ABNF_NAME_SIP_HDR_TO_S)) {
		return newPos, &AbnfError{"To parse: wrong header-name", src, newPos}
	}

	return this.ParseValue(context, src, newPos)
}

func (this *SipHeaderTo) ParseValue(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	//this.Init()
	newPos = pos
	newPos, err = this.addr.Parse(context, src, newPos)
	if err != nil {
		return newPos, err
	}

	return this.params.ParseWithoutInit(context, src, newPos, ';')
}

func (this *SipHeaderTo) Encode(context *ParseContext, buf *bytes.Buffer) {
	buf.WriteString(ABNF_NAME_SIP_HDR_TO_COLON)
	this.EncodeValue(context, buf)
}

func (this *SipHeaderTo) EncodeValue(context *ParseContext, buf *bytes.Buffer) {
	this.addr.Encode(context, buf)
	this.params.Encode(context, buf, ';')
}

func (this *SipHeaderTo) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

func ParseSipTo(context *ParseContext, src []byte, pos int) (newPos int, parsed AbnfPtr, err error) {
	addr := NewSipHeaderTo(context)
	if addr == ABNF_PTR_NIL {
		return newPos, ABNF_PTR_NIL, &AbnfError{"To parse: out of memory for new header", src, newPos}
	}
	newPos, err = addr.GetSipHeaderTo(context).ParseValue(context, src, pos)
	return newPos, addr, err
}

func EncodeSipToValue(parsed AbnfPtr, context *ParseContext, buf *bytes.Buffer) {
	if parsed == ABNF_PTR_NIL {
		return
	}
	parsed.GetSipHeaderTo(context).EncodeValue(context, buf)
}
