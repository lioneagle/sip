package sipparser

import (
	"bytes"
	//"fmt"
	"unsafe"
)

type SipMsgBody struct {
	body    AbnfBuf
	headers SipHeaders
}

func NewSipMsgBody(context *ParseContext) (*SipMsgBody, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipMsgBody{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*SipMsgBody)(unsafe.Pointer(mem)).Init()
	return (*SipMsgBody)(unsafe.Pointer(mem)), addr
}

func (this *SipMsgBody) Init() {
	this.body.Init()
	this.headers.Init()
}

func (this *SipMsgBody) Encode(context *ParseContext, buf *bytes.Buffer) {
	headerPtr, ok := this.headers.GetSingleHeaderParsed(context, ABNF_NAME_SIP_HDR_CONTENT_LENGTH)
	if ok {
		headerPtr.GetSipHeaderContentLength(context).SetValue(this.body.Size())
	}
	this.headers.Encode(context, buf)
	buf.WriteString("\r\n")
	this.body.Encode(context, buf)
}

func (this *SipMsgBody) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

type SipMsgBodies struct {
	AbnfList
}

func NewSipMsgBodies(context *ParseContext) (*SipMsgBodies, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipMsgBodies{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*SipMsgBodies)(unsafe.Pointer(mem)).Init()
	return (*SipMsgBodies)(unsafe.Pointer(mem)), addr
}

func (this *SipMsgBodies) Init() {
	this.AbnfList.Init()
}

func (this *SipMsgBodies) Size() int32 { return this.Len() }
func (this *SipMsgBodies) Empty() bool { return this.Len() == 0 }

func (this *SipMsgBodies) AddBody(context *ParseContext, body AbnfPtr) {
	this.PushBack(context, body)
}

func (this *SipMsgBodies) EncodeSingle(context *ParseContext, buf *bytes.Buffer) {
	e := this.Front(context)
	if e.Value == ABNF_PTR_NIL {
		return
	}

	body := e.Value.GetSipMsgBody(context)
	body.body.Encode(context, buf)
}

/* RFC2046
 * boundary := 0*69<bchars> bcharsnospace
 *
 * bchars := bcharsnospace / " "
 *
 * bcharsnospace := DIGIT / ALPHA / "'" / "(" / ")" /
 *                  "+" / "_" / "," / "-" / "." /
 *                  "/" / ":" / "=" / "?"
 *
 * body-part := <"message" as defined in RFC 822, with all
 *               header fields optional, not starting with the
 *               specified dash-boundary, and with the
 *               delimiter not occurring anywhere in the
 *               body part.  Note that the semantics of a
 *               part differ from the semantics of a message,
 *               as described in the text.>
 *
 * close-delimiter := delimiter "--"
 *
 * dash-boundary := "--" boundary
 *                  ; boundary taken from the value of
 *                  ; boundary parameter of the
 *                  ; Content-Type field.
 *
 * delimiter := CRLF dash-boundary
 *
 * discard-text := *(*text CRLF)
 *                 ; May be ignored or discarded.
 *
 * encapsulation := delimiter transport-padding
 *                  CRLF body-part
 *
 * epilogue := discard-text
 *
 * multipart-body := [preamble CRLF]
 *                   dash-boundary transport-padding CRLF
 *                   body-part *encapsulation
 *                   close-delimiter transport-padding
 *                   [CRLF epilogue]
 *
 * preamble := discard-text
 *
 * transport-padding := *LWSP-char
 *                      ; Composers MUST NOT generate
 *                      ; non-zero length transport
 *                      ; padding, but receivers MUST
 *                      ; be able to handle padding
 *                      ; added by message transports.
 */
func (this *SipMsgBodies) EncodeMulti(context *ParseContext, buf *bytes.Buffer, boundary []byte) {
	for e := this.Front(context); e != nil; e = e.Next(context) {
		v := e.Value.GetSipMsgBody(context)
		// dash-boundary
		buf.WriteString("--")
		buf.Write(boundary)
		buf.WriteString("\r\n")

		// body-part
		v.Encode(context, buf)

		// CRLF
		buf.WriteString("\r\n")
	}
	//
	buf.WriteString("--")
	buf.Write(boundary)
	buf.WriteString("--")
}

func (this *SipMsgBodies) StringSingle(context *ParseContext) string {
	var buf bytes.Buffer
	this.EncodeSingle(context, &buf)
	return buf.String()
}

func (this *SipMsgBodies) StringMulti(context *ParseContext, boundary []byte) string {
	var buf bytes.Buffer
	this.EncodeMulti(context, &buf, boundary)
	return buf.String()
}
