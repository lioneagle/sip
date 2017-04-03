package sipparser

import (
	"bytes"
	//"fmt"
	"unsafe"
)

type SipHeaderMaxForwards struct {
	size uint32
}

func NewSipHeaderMaxForwards(context *ParseContext) (*SipHeaderMaxForwards, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipHeaderMaxForwards{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*SipHeaderMaxForwards)(unsafe.Pointer(mem)).Init()
	return (*SipHeaderMaxForwards)(unsafe.Pointer(mem)), addr
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
	name, newPos, err := ParseHeaderName(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if !EqualNoCase(src[name.Begin:name.End], StringToByteSlice("Max-Forwards")) {
		return newPos, &AbnfError{"Max-Forwards parse: wrong header-name", src, newPos}
	}

	return this.ParseValue(context, src, newPos)
}

func (this *SipHeaderMaxForwards) ParseValue(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	newPos = pos
	digit, _, newPos, ok := ParseUInt(src, newPos)
	if !ok {
		return newPos, &AbnfError{"Max-Forwards parse: wrong num", src, newPos}
	}

	this.size = uint32(digit)
	return newPos, nil
}

func (this *SipHeaderMaxForwards) Encode(context *ParseContext, buf *bytes.Buffer) {
	buf.WriteString("Max-Forwards: ")
	EncodeUInt(buf, uint64(this.size))
}

func (this *SipHeaderMaxForwards) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

func ParseSipMaxForwards(context *ParseContext, src []byte, pos int) (newPos int, parsed AbnfPtr, err error) {
	header, addr := NewSipHeaderMaxForwards(context)
	if header == nil || addr == ABNF_PTR_NIL {
		return newPos, ABNF_PTR_NIL, &AbnfError{"Max-Forwards parse: out of memory for new header", src, newPos}
	}
	newPos, err = header.ParseValue(context, src, pos)
	return newPos, addr, err
}
