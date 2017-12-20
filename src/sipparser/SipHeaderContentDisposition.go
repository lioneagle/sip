package sipparser

import (
	//"bytes"
	//"fmt"
	"unsafe"
)

type SipHeaderContentDisposition struct {
	dispType AbnfBuf
	params   SipGenericParams
}

func NewSipHeaderContentDisposition(context *ParseContext) AbnfPtr {
	addr := context.allocator.Alloc(uint32(unsafe.Sizeof(SipHeaderContentDisposition{})))
	if addr == ABNF_PTR_NIL {
		return ABNF_PTR_NIL
	}
	addr.GetSipHeaderContentDisposition(context).Init()
	return addr
}

func (this *SipHeaderContentDisposition) Init() {
	this.dispType.Init()
	this.params.Init()
}

func (this *SipHeaderContentDisposition) AllowMulti() bool { return false }
func (this *SipHeaderContentDisposition) HasValue() bool   { return true }

/* RFC3261
 *
 * Content-Disposition   =  "Content-Disposition" HCOLON
 *                          disp-type *( SEMI disp-param )
 * disp-type             =  "render" / "session" / "icon" / "alert"
 *                          / disp-extension-token
 * disp-param            =  handling-param / generic-param
 * handling-param        =  "handling" EQUAL
 *                          ( "optional" / "required"
 *                          / other-handling )
 * other-handling        =  token
 * disp-extension-token  =  token
 *
 */
func (this *SipHeaderContentDisposition) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	return this.ParseWithoutInit(context, src, newPos)
}

func (this *SipHeaderContentDisposition) ParseWithoutInit(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	name, newPos, err := ParseHeaderName(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if !EqualNoCase(src[name.Begin:name.End], StringToByteSlice("Content-Disposition")) {
		return newPos, &AbnfError{"Content-Disposition parse: wrong header-name", src, newPos}
	}

	return this.ParseValueWithoutInit(context, src, newPos)
}

func (this *SipHeaderContentDisposition) ParseValue(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	return this.ParseValueWithoutInit(context, src, pos)
}

func (this *SipHeaderContentDisposition) ParseValueWithoutInit(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos, err = this.dispType.ParseSipToken(context, src, pos)
	if err != nil {
		return newPos, err
	}

	return this.params.ParseWithoutInit(context, src, newPos, ';')
}

func (this *SipHeaderContentDisposition) Encode(context *ParseContext, buf *AbnfByteBuffer) {
	buf.WriteString("Content-Disposition: ")
	this.EncodeValue(context, buf)
}

func (this *SipHeaderContentDisposition) EncodeValue(context *ParseContext, buf *AbnfByteBuffer) {
	this.dispType.Encode(context, buf)
	this.params.Encode(context, buf, ';')
}

func (this *SipHeaderContentDisposition) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

func ParseSipContentDisposition(context *ParseContext, src []byte, pos int) (newPos int, parsed AbnfPtr, err error) {
	addr := NewSipHeaderContentDisposition(context)
	if addr == ABNF_PTR_NIL {
		return newPos, ABNF_PTR_NIL, &AbnfError{"Content-Disposition parse: out of memory for new header", src, newPos}
	}
	newPos, err = addr.GetSipHeaderContentDisposition(context).ParseValueWithoutInit(context, src, pos)
	return newPos, addr, err
}

func EncodeSipContentDispositionValue(parsed AbnfPtr, context *ParseContext, buf *AbnfByteBuffer) {
	if parsed == ABNF_PTR_NIL {
		return
	}
	parsed.GetSipHeaderContentDisposition(context).EncodeValue(context, buf)
}
