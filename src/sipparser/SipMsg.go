package sipparser

import (
	"bytes"
	//"fmt"
	"strconv"
	"unsafe"
)

type SipMsg struct {
	startLine SipStartLine
	headers   SipHeaders
	bodies    SipMsgBodies
}

func NewSipMsg(context *ParseContext) (*SipMsg, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipMsg{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*SipMsg)(unsafe.Pointer(mem)).Init()
	return (*SipMsg)(unsafe.Pointer(mem)), addr
}

func (this *SipMsg) Init() {
	this.startLine.Init()
	this.headers.Init()
}

func (this *SipMsg) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	newPos = pos
	newPos, err = this.startLine.Parse(context, src, newPos)
	if err != nil {
		return newPos, err
	}

	if newPos >= len(src) {
		return newPos, &AbnfError{"SipMsg parse: no headers or CRLF", src, newPos}
	}

	newPos, err = this.headers.Parse(context, src, newPos)
	if err != nil {
		return newPos, err
	}

	return this.ParseMsgBody(context, src, newPos)
}

func (this *SipMsg) ParseMsgBody(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos
	headerPtr, ok := this.headers.GetSingleHeaderParsed(context, ABNF_NAME_SIP_HDR_CONTENT_TYPE)
	if !ok {
		// no Content-Type means no msg-body
		return newPos, nil
	}

	contentType := headerPtr.GetSipHeaderContentType(context)
	if contentType.mainType.EqualStringNoCase(context, "multipart") {
		// mime bodies
		return this.ParseMultiBody(context, src, newPos)
	}

	return this.ParseSingleBody(context, src, newPos)
}

func (this *SipMsg) ParseSingleBody(context *ParseContext, src []byte, pos int) (newPos int, err error) {

	return newPos, nil
}

func (this *SipMsg) ParseMultiBody(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	return newPos, nil
}

func (this *SipMsg) Encode(context *ParseContext, buf *bytes.Buffer) error {
	// @@TODO: calculate length of msg-body and create or modify Content-Length header
	var contentLength *SipSingHeader
	headerPtr, ok := this.headers.GetSingleHeaderParsed(context, ABNF_NAME_SIP_HDR_CONTENT_LENGTH)
	if !ok {
		contentLength, addr := NewSipSingleHeader(context)
		if contentLength == nil {
			return &AbnfError{"SipMsg  encode: out of memory for creating Content-Length", nil, 0}
		}
		info, _ := GetSipHeaderInfo(ABNF_NAME_SIP_HDR_CONTENT_LENGTH)
		contentLength.info = info

	}

	this.startLine.Encode(context, buf)
	this.headers.Encode(context, buf)
	buf.WriteString("\r\n")

	bodyStart := len(buf.Bytes())
	if this.bodies.Empty() {
		return nil
	}

	if this.bodies.Size() == 1 {
		this.bodies.EncodeSingle(context, buf)
	} else {
		// remove Content-* headers from sip message */
		// @@TODO: get boundary from msg or create one
		boundary := []byte("asassdada")
		this.bodies.EncodeMulti(context, buf, boundary)
	}

	// modify Content-Length size or create Content-Length
	bodySize := strconv.FormatUint(uint64(len(buf.Bytes())-bodyStart), 10)

}

func (this *SipMsg) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}
