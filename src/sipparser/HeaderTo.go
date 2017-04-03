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

func NewSipHeaderTo(context *ParseContext) (*SipHeaderTo, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipHeaderTo{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*SipHeaderTo)(unsafe.Pointer(mem)).Init()
	return (*SipHeaderTo)(unsafe.Pointer(mem)), addr
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

	if !EqualNoCase(src[name.Begin:name.End], StringToByteSlice("To")) && !EqualNoCase(src[name.Begin:name.End], StringToByteSlice("f")) {
		return newPos, &AbnfError{"To parse: wrong header-name", src, newPos}
	}

	return this.ParseValue(context, src, newPos)
}

func (this *SipHeaderTo) ParseValue(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	newPos = pos
	newPos, err = this.addr.Parse(context, src, newPos)
	if err != nil {
		return newPos, err
	}

	return this.params.Parse(context, src, newPos, ';')
}

func (this *SipHeaderTo) Encode(context *ParseContext, buf *bytes.Buffer) {
	buf.WriteString("To: ")
	this.addr.Encode(context, buf)
	this.params.Encode(context, buf, ';')
}

func (this *SipHeaderTo) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

func ParseSipTo(context *ParseContext, src []byte, pos int) (newPos int, parsed AbnfPtr, err error) {
	header, addr := NewSipHeaderTo(context)
	if header == nil || addr == ABNF_PTR_NIL {
		return newPos, ABNF_PTR_NIL, &AbnfError{"To parse: out of memory for new header", src, newPos}
	}
	newPos, err = header.ParseValue(context, src, pos)
	return newPos, addr, err
}
