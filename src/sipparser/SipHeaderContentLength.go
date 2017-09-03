package sipparser

import (
	"bytes"
	"fmt"
	"unsafe"
)

type SipHeaderContentLength struct {
	size      uint32
	encodeEnd uint32 // record end position when encoding for modify length of sip msg
}

func NewSipHeaderContentLength(context *ParseContext) AbnfPtr {
	addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipHeaderContentLength{})))
	if addr == ABNF_PTR_NIL {
		return ABNF_PTR_NIL
	}
	addr.GetSipHeaderContentLength(context).Init()
	return addr
}

func (this *SipHeaderContentLength) Init() {
	this.size = 0
	this.encodeEnd = 0
}

func (this *SipHeaderContentLength) AllowMulti() bool    { return false }
func (this *SipHeaderContentLength) HasValue() bool      { return true }
func (this *SipHeaderContentLength) SetValue(size int32) { this.size = uint32(size) }

/* RFC3261
 *
 * Content-Length  =  ( "Content-Length" / "l" ) HCOLON 1*DIGIT
 */
func (this *SipHeaderContentLength) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	name, newPos, err := ParseHeaderName(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if !EqualNoCase(src[name.Begin:name.End], StringToByteSlice(ABNF_NAME_SIP_HDR_CONTENT_LENGTH)) &&
		!EqualNoCase(src[name.Begin:name.End], StringToByteSlice(ABNF_NAME_SIP_HDR_CONTENT_LENGTH_S)) {
		return newPos, &AbnfError{"Content-Length parse: wrong header-name", src, newPos}
	}

	return this.ParseValue(context, src, newPos)
}

func (this *SipHeaderContentLength) ParseValue(context *ParseContext, src []byte, pos int) (newPos int, err error) {

	//fmt.Println("enter Content-Length ParseValue")
	this.Init()
	newPos = pos
	digit, _, newPos, ok := ParseUInt(src, newPos)
	if !ok {
		return newPos, &AbnfError{"Content-Length parse: wrong num", src, newPos}
	}

	this.size = uint32(digit)
	return newPos, nil
}

func (this *SipHeaderContentLength) Encode(context *ParseContext, buf *bytes.Buffer) {
	buf.WriteString(ABNF_NAME_SIP_HDR_CONTENT_LENGTH_COLON)
	this.EncodeValue(context, buf)
}

func (this *SipHeaderContentLength) EncodeValue(context *ParseContext, buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf(ABNF_SIP_CONTENT_LENGTH_PRINT_FMT, this.size))
	this.encodeEnd = uint32(len(buf.Bytes()))
}

func (this *SipHeaderContentLength) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

func ParseSipContentLength(context *ParseContext, src []byte, pos int) (newPos int, parsed AbnfPtr, err error) {
	addr := NewSipHeaderContentLength(context)
	if addr == ABNF_PTR_NIL {
		return newPos, ABNF_PTR_NIL, &AbnfError{"Content-Length parse: out of memory for new header", src, newPos}
	}
	newPos, err = addr.GetSipHeaderContentLength(context).ParseValue(context, src, pos)
	return newPos, addr, err
}

func EncodeSipContentLengthValue(parsed AbnfPtr, context *ParseContext, buf *bytes.Buffer) {
	if parsed == ABNF_PTR_NIL {
		return
	}
	parsed.GetSipHeaderContentLength(context).EncodeValue(context, buf)
}