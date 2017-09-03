package sipparser

import (
	"bytes"
	//"fmt"
	"unsafe"
)

type SipHeaderCallId struct {
	id1 AbnfBuf
	id2 AbnfBuf
}

func NewSipHeaderCallId(context *ParseContext) AbnfPtr {
	addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipHeaderCallId{})))
	if addr == ABNF_PTR_NIL {
		return ABNF_PTR_NIL
	}
	addr.GetSipHeaderCallId(context).Init()

	return addr
}

func (this *SipHeaderCallId) Init() {
	this.id1.Init()
	this.id2.Init()
}

func (this *SipHeaderCallId) AllowMulti() bool { return false }
func (this *SipHeaderCallId) HasValue() bool   { return true }

/* RFC3261
 *
 * Call-ID  =  ( "Call-ID" / "i" ) HCOLON callid
 * callid   =  word [ "@" word ]
 */
func (this *SipHeaderCallId) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	name, newPos, err := ParseHeaderName(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if !EqualNoCase(src[name.Begin:name.End], StringToByteSlice(ABNF_NAME_SIP_HDR_CALL_ID)) &&
		!EqualNoCase(src[name.Begin:name.End], StringToByteSlice(ABNF_NAME_SIP_HDR_CALL_ID_S)) {
		return newPos, &AbnfError{"Call-ID parse: wrong header-name", src, newPos}
	}

	return this.ParseValue(context, src, newPos)
}

func (this *SipHeaderCallId) ParseValue(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	newPos = pos
	newPos, err = this.id1.ParseSipWord(context, src, newPos)
	if err != nil {
		return newPos, err
	}

	if newPos >= len(src) {
		return newPos, nil
	}

	if src[newPos] == '@' {
		return this.id2.ParseSipWord(context, src, newPos+1)
	}

	return newPos, nil
}

func (this *SipHeaderCallId) Encode(context *ParseContext, buf *bytes.Buffer) {
	buf.WriteString(ABNF_NAME_SIP_HDR_CALL_ID_COLON)
	this.EncodeValue(context, buf)
}

func (this *SipHeaderCallId) EncodeValue(context *ParseContext, buf *bytes.Buffer) {
	this.id1.Encode(context, buf)
	if this.id2.Exist() {
		buf.WriteByte('@')
		this.id2.Encode(context, buf)
	}
}

func (this *SipHeaderCallId) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

func ParseSipCallId(context *ParseContext, src []byte, pos int) (newPos int, parsed AbnfPtr, err error) {
	addr := NewSipHeaderCallId(context)
	if addr == ABNF_PTR_NIL {
		return newPos, ABNF_PTR_NIL, &AbnfError{"Call-ID parse: out of memory for new header", src, newPos}
	}
	newPos, err = addr.GetSipHeaderCallId(context).ParseValue(context, src, pos)
	return newPos, addr, err
}

func EncodeSipCallIdValue(parsed AbnfPtr, context *ParseContext, buf *bytes.Buffer) {
	if parsed == ABNF_PTR_NIL {
		return
	}
	parsed.GetSipHeaderCallId(context).EncodeValue(context, buf)
}
