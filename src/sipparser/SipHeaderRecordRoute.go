package sipparser

import (
	"bytes"
	//"fmt"
	"unsafe"
)

type SipHeaderRecordRoute struct {
	addr   SipNameAddr
	params SipGenericParams
}

func NewSipHeaderRecordRoute(context *ParseContext) AbnfPtr {
	addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipHeaderRecordRoute{})))
	if addr == ABNF_PTR_NIL {
		return ABNF_PTR_NIL
	}
	addr.GetSipHeaderRecordRoute(context).Init()
	return addr
}

func (this *SipHeaderRecordRoute) Init() {
	this.addr.Init()
	this.params.Init()
}

func (this *SipHeaderRecordRoute) AllowMulti() bool { return false }
func (this *SipHeaderRecordRoute) HasValue() bool   { return true }

/* RFC3261
 *
 * Record-Route  =  "Record-Route" HCOLON rec-route *(COMMA rec-route)
 * rec-route     =  name-addr *( SEMI rr-param )
 * rr-param      =  generic-param
 */
func (this *SipHeaderRecordRoute) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	name, newPos, err := ParseHeaderName(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if !EqualNoCase(src[name.Begin:name.End], StringToByteSlice("Record-Route")) {
		return newPos, &AbnfError{"Record-Route parse: wrong header-name", src, newPos}
	}

	return this.ParseValue(context, src, newPos)
}

func (this *SipHeaderRecordRoute) ParseValue(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	newPos = pos
	newPos, err = this.addr.Parse(context, src, newPos)
	if err != nil {
		return newPos, err
	}

	return this.params.Parse(context, src, newPos, ';')
}

func (this *SipHeaderRecordRoute) Encode(context *ParseContext, buf *bytes.Buffer) {
	buf.WriteString("Record-Route: ")
	this.EncodeValue(context, buf)
}

func (this *SipHeaderRecordRoute) EncodeValue(context *ParseContext, buf *bytes.Buffer) {
	this.addr.Encode(context, buf)
	this.params.Encode(context, buf, ';')
}

func (this *SipHeaderRecordRoute) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

func ParseSipRecordRoute(context *ParseContext, src []byte, pos int) (newPos int, parsed AbnfPtr, err error) {
	addr := NewSipHeaderRecordRoute(context)
	if addr == ABNF_PTR_NIL {
		return newPos, ABNF_PTR_NIL, &AbnfError{"Record-Route parse: out of memory for new header", src, newPos}
	}
	newPos, err = addr.GetSipHeaderRecordRoute(context).ParseValue(context, src, pos)
	return newPos, addr, err
}

func EncodeSipRecordRouteValue(parsed AbnfPtr, context *ParseContext, buf *bytes.Buffer) {
	if parsed == ABNF_PTR_NIL {
		return
	}
	parsed.GetSipHeaderRecordRoute(context).EncodeValue(context, buf)
}
