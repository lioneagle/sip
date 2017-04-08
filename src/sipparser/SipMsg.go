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
		// get boundary
		boundary, ok := contentType.GetBoundary(context)
		if !ok {
			return newPos, &AbnfError{"SipMsg parse: no boundary for multipart body", src, newPos}
		}
		return this.ParseMultiBody(context, src, newPos, boundary)
	}

	return this.ParseSingleBody(context, src, newPos)
}

func (this *SipMsg) ParseSingleBody(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	left := len(src) - pos
	bodySize := 0
	parsedPtr, ok := this.headers.GetSingleHeaderParsed(context, ABNF_NAME_SIP_HDR_CONTENT_LENGTH)
	if ok {
		bodySize = int(parsedPtr.GetSipHeaderContentLength(context).size)
		if bodySize > left {
			bodySize = left
		}
	}

	if bodySize == 0 {
		return newPos, nil
	}

	body, addr := NewSipMsgBody(context)
	if body == nil {
		return newPos, &AbnfError{"SipMsg parse: out of memory for sip-mshg-body", src, newPos}
	}
	// copy Content-* headers from sip-msg to sip-msg-body's headers
	body.headers.CopyContentHeaders(context, &this.headers)

	body.SetBody(context, src[pos:pos+bodySize])
	this.bodies.AddBody(context, addr)
	return newPos, nil
}

func (this *SipMsg) ParseMultiBody(context *ParseContext, src []byte, pos int, boundary []byte) (newPos int, err error) {
	var ok bool

	dash_boundary := append([]byte{'-', '-'}, boundary...)
	delimiter := append([]byte{'\r', '\n'}, dash_boundary...)

	newPos = pos
	pos1 := bytes.Index(src[newPos:], dash_boundary)
	if pos1 != 0 {
		return newPos, &AbnfError{"SipMsg ParseMultiBody: no first dash-bounday", src, newPos}
	}

	newPos += len(dash_boundary)

	for newPos < len(src) {
		if newPos+1 >= len(src) {

			return newPos, &AbnfError{"SipMsg ParseMultiBody: reach end without close-delimiter", src, newPos}
		}

		if src[newPos] == '-' || src[newPos+1] == '-' {
			// reach close-delimiter
			return newPos, nil
		}

		// skip transport-padding CRLF
		_, newPos, ok = FindCrlfRFC3261(src, newPos)
		if !ok {
			return newPos, &AbnfError{"SipMsg ParseMultiBody: no CRLF after dash-bounday", src, newPos}
		}

		body, addr := NewSipMsgBody(context)
		if body == nil {
			return newPos, &AbnfError{"SipMsg ParseMultiBody: out of memory for body", src, newPos}
		}

		newPos, err = body.headers.Parse(context, src, newPos)
		if err != nil {
			return newPos, &AbnfError{"SipMsg ParseMultiBody: parse headers failed", src, newPos}
		}

		begin := newPos
		end := bytes.Index(src[newPos:], delimiter)
		if end == -1 {
			return newPos, &AbnfError{"SipMsg ParseMultiBody: no delimiter", src, newPos}
		}

		end += newPos
		newPos = end + len(delimiter)

		body.body.SetByteSlice(context, src[begin:end])
		this.bodies.AddBody(context, addr)
	}

	return newPos, &AbnfError{"SipMsg ParseMultiBody: no close-delimiter", src, newPos}
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
		// create Content-Type header
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
	boundary = StringToByteSlice(ABNF_SIP_DEFAULT_BOUNDARY)
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

	var boundary []byte

	if this.bodies.Size() > 1 {
		// remove Content-* headers from sip message except Content-Length and Content-Type*/
		this.headers.RemoveContentHeaders(context)

		_, ok := this.headers.GetSingleHeader(context, "MIME-Version")
		if !ok {
			addr, _ := this.headers.GenerateAndAddSingleHeader(context, "MIME-Version", "1.0")
			if addr == nil {
				return &AbnfError{"SipMsg encode: out of memory for adding MIME-Version", nil, 0}
			}
		}

		boundary = this.FindOrCreateBoundary(context)
	}

	this.headers.Encode(context, buf)
	buf.WriteString("\r\n")

	if this.bodies.Empty() {
		return nil
	}

	bodyStart := len(buf.Bytes())

	if this.bodies.Size() == 1 {
		this.bodies.EncodeSingle(context, buf)
	} else {
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
