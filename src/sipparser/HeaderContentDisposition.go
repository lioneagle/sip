package sipparser

import (
	"bytes"
	//"fmt"
	"unsafe"
)

type SipHeaderContentDisposition struct {
	dispType AbnfToken
	params   SipGenericParams
}

func NewSipHeaderContentDisposition(context *ParseContext) (*SipHeaderContentDisposition, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipHeaderContentDisposition{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*SipHeaderContentDisposition)(unsafe.Pointer(mem)).Init()
	return (*SipHeaderContentDisposition)(unsafe.Pointer(mem)), addr
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
	name, newPos, err := ParseHeaderName(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if !EqualNoCase(src[name.Begin:name.End], StringToByteSlice("Content-Disposition")) {
		return newPos, &AbnfError{"Content-Disposition parse: wrong header-name", src, newPos}
	}

	return this.ParseValue(context, src, newPos)
}

func (this *SipHeaderContentDisposition) ParseValue(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	newPos, err = this.dispType.Parse(context, src, pos, IsSipToken)
	if err != nil {
		return newPos, err
	}

	return this.params.Parse(context, src, newPos, ';')
}

func (this *SipHeaderContentDisposition) Encode(context *ParseContext, buf *bytes.Buffer) {
	buf.WriteString("Content-Disposition: ")
	this.dispType.Encode(context, buf)
	this.params.Encode(context, buf, ';')
}

func (this *SipHeaderContentDisposition) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

func ParseSipContentDisposition(context *ParseContext, src []byte, pos int) (newPos int, parsed AbnfPtr, err error) {
	header, addr := NewSipHeaderContentDisposition(context)
	if header == nil || addr == ABNF_PTR_NIL {
		return newPos, ABNF_PTR_NIL, &AbnfError{"Content-Disposition parse: out of memory for new header", src, newPos}
	}
	newPos, err = header.ParseValue(context, src, pos)
	return newPos, addr, err
}
