package sipparser

import (
	"bytes"
	//"fmt"
	"unsafe"
)

type SipHeaderCseq struct {
	id     uint32
	method AbnfBuf
}

func NewSipHeaderCseq(context *ParseContext) (*SipHeaderCseq, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipHeaderCseq{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*SipHeaderCseq)(unsafe.Pointer(mem)).Init()
	return (*SipHeaderCseq)(unsafe.Pointer(mem)), addr
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
	name, newPos, err := ParseHeaderName(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if !EqualNoCase(src[name.Begin:name.End], StringToByteSlice("CSeq")) {
		return newPos, &AbnfError{"CSeq parse: wrong header-name", src, newPos}
	}

	return this.ParseValue(context, src, newPos)
}

func (this *SipHeaderCseq) ParseValue(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	newPos = pos
	digit, _, newPos, ok := ParseUInt(src, newPos)
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

func (this *SipHeaderCseq) Encode(context *ParseContext, buf *bytes.Buffer) {
	buf.WriteString("CSeq: ")
	this.EncodeValue(context, buf)
}

func (this *SipHeaderCseq) EncodeValue(context *ParseContext, buf *bytes.Buffer) {
	EncodeUInt(buf, uint64(this.id))
	buf.WriteByte(' ')
	this.method.Encode(context, buf)
}

func (this *SipHeaderCseq) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

func ParseSipCseq(context *ParseContext, src []byte, pos int) (newPos int, parsed AbnfPtr, err error) {
	header, addr := NewSipHeaderCseq(context)
	if header == nil || addr == ABNF_PTR_NIL {
		return newPos, ABNF_PTR_NIL, &AbnfError{"CSeq parse: out of memory for new header", src, newPos}
	}
	newPos, err = header.ParseValue(context, src, pos)
	return newPos, addr, err
}

func EncodeSipCseqValue(parsed AbnfPtr, context *ParseContext, buf *bytes.Buffer) {
	if parsed == ABNF_PTR_NIL {
		return
	}
	parsed.GetSipHeaderCseq(context).EncodeValue(context, buf)
}
