package sipparser

import (
	"bytes"
	//"fmt"
	"unsafe"
)

type SipHeaderRoute struct {
	addr   SipNameAddr
	params SipGenericParams
}

func NewSipHeaderRoute(context *ParseContext) (*SipHeaderRoute, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipHeaderRoute{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*SipHeaderRoute)(unsafe.Pointer(mem)).Init()
	return (*SipHeaderRoute)(unsafe.Pointer(mem)), addr
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
	name, newPos, err := ParseHeaderName(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if !EqualNoCase(src[name.Begin:name.End], StringToByteSlice("Route")) {
		return newPos, &AbnfError{"Route parse: wrong header-name", src, newPos}
	}

	return this.ParseValue(context, src, newPos)
}

func (this *SipHeaderRoute) ParseValue(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	newPos = pos
	newPos, err = this.addr.Parse(context, src, newPos)
	if err != nil {
		return newPos, err
	}

	return this.params.Parse(context, src, newPos, ';')
}

func (this *SipHeaderRoute) Encode(context *ParseContext, buf *bytes.Buffer) {
	buf.WriteString("Route: ")
	this.EncodeValue(context, buf)
}

func (this *SipHeaderRoute) EncodeValue(context *ParseContext, buf *bytes.Buffer) {
	this.addr.Encode(context, buf)
	this.params.Encode(context, buf, ';')
}

func (this *SipHeaderRoute) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

func ParseSipRoute(context *ParseContext, src []byte, pos int) (newPos int, parsed AbnfPtr, err error) {
	header, addr := NewSipHeaderRoute(context)
	if header == nil || addr == ABNF_PTR_NIL {
		return newPos, ABNF_PTR_NIL, &AbnfError{"Route parse: out of memory for new header", src, newPos}
	}
	newPos, err = header.ParseValue(context, src, pos)
	return newPos, addr, err
}

func EncodeSipRouteValue(parsed AbnfPtr, context *ParseContext, buf *bytes.Buffer) {
	if parsed == ABNF_PTR_NIL {
		return
	}
	parsed.GetSipHeaderRoute(context).EncodeValue(context, buf)
}
