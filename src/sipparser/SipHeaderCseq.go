package sipparser

import (
	//"bytes"
	//"fmt"
	"unsafe"
)

type SipHeaderCseq struct {
	id     uint32
	method AbnfBuf
}

func NewSipHeaderCseq(context *ParseContext) AbnfPtr {
	addr := context.allocator.Alloc(uint32(unsafe.Sizeof(SipHeaderCseq{})))
	if addr == ABNF_PTR_NIL {
		return ABNF_PTR_NIL
	}
	addr.GetSipHeaderCseq(context).Init()
	return addr
}

func (this *SipHeaderCseq) Init() {
	this.id = 0
	this.method.Init()
}

func (this *SipHeaderCseq) AllowMulti() bool { return false }
func (this *SipHeaderCseq) HasValue() bool   { return true }

/* RFC3261
 *
 * CSeq  =  "CSeq" HCOLON 1*DIGIT LWS Method
 *
 */
func (this *SipHeaderCseq) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	return this.ParseWithoutInit(context, src, newPos)
}

func (this *SipHeaderCseq) ParseWithoutInit(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	name, newPos, err := ParseHeaderName(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if !EqualNoCase(src[name.Begin:name.End], StringToByteSlice("CSeq")) {
		return newPos, &AbnfError{"CSeq parse: wrong header-name", src, newPos}
	}

	return this.ParseValueWithoutInit(context, src, newPos)
}

func (this *SipHeaderCseq) ParseValue(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	return this.ParseValueWithoutInit(context, src, pos)
}

func (this *SipHeaderCseq) ParseValueWithoutInit(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	digit, _, newPos, ok := ParseUInt(src, pos)
	if !ok {
		return newPos, &AbnfError{"CSeq parse: wrong num", src, newPos}
	}

	this.id = uint32(digit)

	newPos, ok = ParseLWS(src, newPos)
	if !ok {
		return newPos, &AbnfError{"CSeq parse: wong LWS", src, newPos}
	}

	return this.method.ParseSipToken(context, src, newPos)
}

func (this *SipHeaderCseq) Encode(context *ParseContext, buf *AbnfByteBuffer) {
	buf.WriteString("CSeq: ")
	this.EncodeValue(context, buf)
}

func (this *SipHeaderCseq) EncodeValue(context *ParseContext, buf *AbnfByteBuffer) {
	EncodeUInt(buf, uint64(this.id))
	buf.WriteByte(' ')
	this.method.Encode(context, buf)
}

func (this *SipHeaderCseq) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

func ParseSipCseq(context *ParseContext, src []byte, pos int) (newPos int, parsed AbnfPtr, err error) {
	addr := NewSipHeaderCseq(context)
	if addr == ABNF_PTR_NIL {
		return newPos, ABNF_PTR_NIL, &AbnfError{"CSeq parse: out of memory for new header", src, newPos}
	}
	newPos, err = addr.GetSipHeaderCseq(context).ParseValueWithoutInit(context, src, pos)
	return newPos, addr, err
}

func EncodeSipCseqValue(parsed AbnfPtr, context *ParseContext, buf *AbnfByteBuffer) {
	if parsed == ABNF_PTR_NIL {
		return
	}
	parsed.GetSipHeaderCseq(context).EncodeValue(context, buf)
}
