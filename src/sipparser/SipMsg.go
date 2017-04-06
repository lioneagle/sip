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

func (this *SipMsg) FindOrCreateBoundary(context *ParseContext) (boundary []byte) {
	var parsedContentType *SipHeaderContentType

	parsedPtr, ok := this.headers.GetSingleHeaderParsed(context, ABNF_NAME_SIP_HDR_CONTENT_TYPE)
	if ok {
		parsedContentType = parsedPtr.GetSipHeaderContentType(context)
		boundary, ok = parsedContentType.GetBoundary(context)
		if ok {
			return boundary
		}

	} else {
		// create Conten-Type
		contentType, addr := this.headers.CreateSingleHeader(context, ABNF_NAME_SIP_HDR_CONTENT_TYPE)
		if contentType == nil {
			return nil
		}
		parsedContentType, parsedPtr = NewSipHeaderContentType(context)
		if parsedContentType == nil {
			return nil
		}
		parsedContentType.SetMainType(context, "multipart")
		parsedContentType.SetSubType(context, "mixed")
		contentType.parsed = parsedPtr
		this.headers.singleHeaders.PushBack(context, addr)
	}

	/* create boundary */
	boundary = StringToByteSlice("sip-unique-boundary-aasdasdewfd")
	parsedContentType.AddBoundary(context, boundary)

	return boundary
}

func (this *SipMsg) Encode(context *ParseContext, buf *bytes.Buffer) error {
	// create Content-Length header if not exist
	err := this.headers.CreateContentLength(context, 0)
	if err != nil {
		return err
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
		_, ok := this.headers.GetSingleHeader(context, "MIME-Version")
		if !ok {
			this.headers.GenerateAndAddSingleHeader(context, "MIME-Version", "1.0")
		}
		// remove Content-* headers from sip message except Content-Length and Content-Type*/
		this.headers.RemoveContentHeaders(context)

		boundary := this.FindOrCreateBoundary(context)
		this.bodies.EncodeMulti(context, buf, boundary)
	}

	// modify Content-Length size or create Content-Length
	bodySize := StringToByteSlice(strconv.FormatUint(uint64(len(buf.Bytes())-bodyStart), 10))
	contentLength, ok := this.headers.GetSingleHeaderParsed(context, ABNF_NAME_SIP_HDR_CONTENT_LENGTH)
	if !ok {
		return &AbnfError{"SipMsg encode: no Content-Length after encoding msg-body", nil, 0}
	}
	encodeEnd := int(contentLength.GetSipHeaderContentLength(context).encodeEnd)
	copy(buf.Bytes()[encodeEnd-len(bodySize):encodeEnd], bodySize)

	return nil
}

func (this *SipMsg) String(context *ParseContext) string {
	var buf bytes.Buffer
	this.Encode(context, &buf)
	return buf.String()
}
