package sipparser

import (
	"bytes"
	//"fmt"
	"unsafe"
)

type SipQuotedString struct {
	value AbnfBuf
}

func NewSipQuotedString(context *ParseContext) (*SipQuotedString, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipQuotedString{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*SipQuotedString)(unsafe.Pointer(mem)).Init()
	return (*SipQuotedString)(unsafe.Pointer(mem)), addr
}

func (this *SipQuotedString) Init() {
	this.value.Init()
}

func (this *SipQuotedString) Encode(context *ParseContext, buf *bytes.Buffer) {
	buf.WriteByte('"')
	this.value.Encode(context, buf)
	buf.WriteByte('"')
}

func (this *SipQuotedString) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

func (this *SipQuotedString) SetValue(context *ParseContext, value []byte) {
	this.value.SetByteSlice(context, value)
}

func (this *SipQuotedString) GetAsByteSlice(context *ParseContext) []byte {
	return this.value.GetAsByteSlice(context)
}

func (this *SipQuotedString) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	/* RFC3261 Section 25.1, page 222
	 *
	 * quoted-string  =  SWS DQUOTE *(qdtext / quoted-pair ) DQUOTE
	 * qdtext         =  LWS / %x21 / %x23-5B / %x5D-7E
	 *                 / UTF8-NONASCII
	 * quoted-pair  =  "\" (%x00-09 / %x0B-0C
	 *               / %x0E-7F)
	 */
	newPos = pos
	newPos, err = ParseSWS(src, newPos)
	if err != nil {
		return newPos, err
	}

	if src[newPos] != '"' {
		return newPos, &AbnfError{"quoted-string parse: no DQUOTE for quoted-string begin", src, newPos}
	}

	newPos++
	tokenBegin := newPos
	for (newPos < len(src)) && (src[newPos] != '"') {
		if IsLwsChar(src[newPos]) {
			newPos, err = ParseLWS(src, newPos)
			if err != nil {
				return newPos, err
			}
		} else if IsSipQuotedText(src[newPos]) {
			newPos++
		} else if src[newPos] == '\\' {
			if (newPos + 1) >= len(src) {
				return newPos, &AbnfError{"quoted-string parse: no char after \\", src, newPos}
			}
			newPos += 2
		} else {
			return newPos, &AbnfError{"quoted-string parse: not qdtext or quoted-pair", src, newPos}
		}
	}

	this.value.SetByteSlice(context, src[tokenBegin:newPos])

	if newPos >= len(src) {
		return newPos, &AbnfError{"quoted-string parse: reach end before DQUOTE", src, newPos}
	}

	if src[newPos] != '"' {
		return newPos, &AbnfError{"quoted-string parse: no DQUOTE for quoted-string end", src, newPos}
	}

	return newPos + 1, nil
}
