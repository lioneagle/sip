package sipparser

import (
	//"bytes"
	//"fmt"
	"unsafe"
)

type SipHeaderMaxForwards struct {
	size uint32
}

func NewSipHeaderMaxForwards(context *ParseContext) AbnfPtr {
	addr := context.allocator.Alloc(uint32(unsafe.Sizeof(SipHeaderMaxForwards{})))
	if addr == ABNF_PTR_NIL {
		return ABNF_PTR_NIL
	}
	addr.GetSipHeaderMaxForwards(context).Init()
	return addr
}

func (this *SipHeaderMaxForwards) Init() {
	this.size = 0
}

func (this *SipHeaderMaxForwards) AllowMulti() bool { return false }
func (this *SipHeaderMaxForwards) HasValue() bool   { return true }

/* RFC3261
 *
 * Max-Forwards  =  "Max-Forwards" HCOLON 1*DIGIT
 */
func (this *SipHeaderMaxForwards) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	return this.ParseWithoutInit(context, src, pos)
}

func (this *SipHeaderMaxForwards) ParseWithoutInit(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	name, newPos, err := ParseHeaderName(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if !EqualNoCase(src[name.Begin:name.End], StringToByteSlice("Max-Forwards")) {
		return newPos, &AbnfError{"Max-Forwards parse: wrong header-name", src, newPos}
	}

	return this.ParseValueWithoutInit(context, src, newPos)
}

func (this *SipHeaderMaxForwards) ParseValue(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	return this.ParseValueWithoutInit(context, src, pos)
}

func (this *SipHeaderMaxForwards) ParseValueWithoutInit(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	digit, _, newPos, ok := ParseUInt(src, pos)
	if !ok {
		return newPos, &AbnfError{"Max-Forwards parse: wrong num", src, newPos}
	}

	this.size = uint32(digit)
	return newPos, nil
}

func (this *SipHeaderMaxForwards) Encode(context *ParseContext, buf *AbnfByteBuffer) {
	buf.WriteString("Max-Forwards: ")
	this.EncodeValue(context, buf)
}

func (this *SipHeaderMaxForwards) EncodeValue(context *ParseContext, buf *AbnfByteBuffer) {
	EncodeUInt(buf, uint64(this.size))
}

func (this *SipHeaderMaxForwards) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

func ParseSipMaxForwards(context *ParseContext, src []byte, pos int) (newPos int, parsed AbnfPtr, err error) {
	addr := NewSipHeaderMaxForwards(context)
	if addr == ABNF_PTR_NIL {
		return newPos, ABNF_PTR_NIL, &AbnfError{"Max-Forwards parse: out of memory for new header", src, newPos}
	}
	newPos, err = addr.GetSipHeaderMaxForwards(context).ParseValueWithoutInit(context, src, pos)
	return newPos, addr, err
}

func EncodeSipMaxForwardsValue(parsed AbnfPtr, context *ParseContext, buf *AbnfByteBuffer) {
	if parsed == ABNF_PTR_NIL {
		return
	}
	parsed.GetSipHeaderMaxForwards(context).EncodeValue(context, buf)
}
