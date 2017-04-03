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

func NewSipHeaderFrom(context *ParseContext) (*SipHeaderFrom, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipHeaderFrom{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*SipHeaderFrom)(unsafe.Pointer(mem)).Init()
	return (*SipHeaderFrom)(unsafe.Pointer(mem)), addr
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
	name, newPos, err := ParseHeaderName(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if !EqualNoCase(src[name.Begin:name.End], StringToByteSlice(ABNF_NAME_SIP_HDR_FROM)) &&
		!EqualNoCase(src[name.Begin:name.End], StringToByteSlice(ABNF_NAME_SIP_HDR_FROM_S)) {
		return newPos, &AbnfError{"From parse: wrong header-name", src, newPos}
	}

	return this.ParseValue(context, src, newPos)
}

func (this *SipHeaderFrom) ParseValue(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	newPos = pos
	newPos, err = this.addr.Parse(context, src, newPos)
	if err != nil {
		return newPos, err
	}

	return this.params.Parse(context, src, newPos, ';')
}

func (this *SipHeaderFrom) Encode(context *ParseContext, buf *bytes.Buffer) {
	buf.WriteString(ABNF_NAME_SIP_HDR_FROM_COLON)
	this.addr.Encode(context, buf)
	this.params.Encode(context, buf, ';')
}

func (this *SipHeaderFrom) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

func ParseSipFrom(context *ParseContext, src []byte, pos int) (newPos int, parsed AbnfPtr, err error) {
	header, addr := NewSipHeaderFrom(context)
	if header == nil || addr == ABNF_PTR_NIL {
		return newPos, ABNF_PTR_NIL, &AbnfError{"From parse: out of memory for new header", src, newPos}
	}
	newPos, err = header.ParseValue(context, src, pos)
	return newPos, addr, err
}
