package sipparser

import (
	//"bytes"
	//"fmt"
	"unsafe"
)

type SipHeaderRoute struct {
	addr   SipNameAddr
	params SipGenericParams
}

func NewSipHeaderRoute(context *ParseContext) AbnfPtr {
	addr := context.allocator.Alloc(uint32(unsafe.Sizeof(SipHeaderRoute{})))
	if addr == ABNF_PTR_NIL {
		return ABNF_PTR_NIL
	}
	addr.GetSipHeaderRoute(context).Init()
	return addr
}

func (this *SipHeaderRoute) Init() {
	this.addr.Init()
	this.params.Init()
}

func (this *SipHeaderRoute) AllowMulti() bool { return false }
func (this *SipHeaderRoute) HasValue() bool   { return true }

/* RFC3261
 *
 * Route        =  "Route" HCOLON route-param *(COMMA route-param)
 * route-param  =  name-addr *( SEMI rr-param )
 * rr-param     =  generic-param
 */
func (this *SipHeaderRoute) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	return this.ParseWithoutInit(context, src, pos)
}

func (this *SipHeaderRoute) ParseWithoutInit(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	name, newPos, err := ParseHeaderName(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if !EqualNoCase(src[name.Begin:name.End], StringToByteSlice("Route")) {
		return newPos, &AbnfError{"Route parse: wrong header-name", src, newPos}
	}

	return this.ParseValueWithoutInit(context, src, newPos)
}

func (this *SipHeaderRoute) ParseValue(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	return this.ParseValueWithoutInit(context, src, newPos)
}

func (this *SipHeaderRoute) ParseValueWithoutInit(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos
	newPos, err = this.addr.ParseWithoutInit(context, src, newPos)
	if err != nil {
		return newPos, err
	}

	return this.params.ParseWithoutInit(context, src, newPos, ';')
}

func (this *SipHeaderRoute) Encode(context *ParseContext, buf *AbnfByteBuffer) {
	buf.WriteString("Route: ")
	this.EncodeValue(context, buf)
}

func (this *SipHeaderRoute) EncodeValue(context *ParseContext, buf *AbnfByteBuffer) {
	this.addr.Encode(context, buf)
	this.params.Encode(context, buf, ';')
}

func (this *SipHeaderRoute) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

func ParseSipRoute(context *ParseContext, src []byte, pos int) (newPos int, parsed AbnfPtr, err error) {
	addr := NewSipHeaderRoute(context)
	if addr == ABNF_PTR_NIL {
		return newPos, ABNF_PTR_NIL, &AbnfError{"Route parse: out of memory for new header", src, newPos}
	}
	newPos, err = addr.GetSipHeaderRoute(context).ParseValueWithoutInit(context, src, pos)
	return newPos, addr, err
}

func EncodeSipRouteValue(parsed AbnfPtr, context *ParseContext, buf *AbnfByteBuffer) {
	if parsed == ABNF_PTR_NIL {
		return
	}
	parsed.GetSipHeaderRoute(context).EncodeValue(context, buf)
}
