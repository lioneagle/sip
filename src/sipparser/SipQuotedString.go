package sipparser

import (
	"bytes"
	//"fmt"
	"unsafe"
)

type SipQuotedString struct {
	value AbnfBuf
}

func NewSipQuotedString(context *ParseContext) AbnfPtr {
	addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipQuotedString{})))
	if addr == ABNF_PTR_NIL {
		return ABNF_PTR_NIL
	}
	addr.GetSipQuotedString(context).Init()
	return addr
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

	newPos, err = this.parseValue(context, src, newPos+1)
	if err != nil {
		return newPos, err
	}

	if newPos >= len(src) {
		return newPos, &AbnfError{"quoted-string parse: reach end before DQUOTE", src, newPos}
	}

	if src[newPos] != '"' {
		return newPos, &AbnfError{"quoted-string parse: no DQUOTE for quoted-string end", src, newPos}
	}

	return newPos + 1, nil
}

func (this *SipQuotedString) parseValue(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	tokenBegin := pos
	newPos = pos
	for (newPos < len(src)) && (src[newPos] != '"') {
		if IsLwsChar(src[newPos]) {
			var ok bool
			newPos, ok = ParseLWS(src, newPos)
			if !ok {
				return newPos, &AbnfError{"quoted-string parse: wrong LWS", src, newPos}
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
	return newPos, nil
}
