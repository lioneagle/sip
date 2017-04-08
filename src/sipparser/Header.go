package sipparser

import (
	"bytes"
	//"fmt"
	"strings"
)

var total_headers int

var g_SipHeaderFullNameMaps = map[string]string{
	ABNF_NAME_SIP_HDR_FROM_S:    ABNF_NAME_SIP_HDR_FROM,
	ABNF_NAME_SIP_HDR_TO_S:      ABNF_NAME_SIP_HDR_TO,
	ABNF_NAME_SIP_HDR_VIA_S:     ABNF_NAME_SIP_HDR_VIA,
	ABNF_NAME_SIP_HDR_CALL_ID_S: ABNF_NAME_SIP_HDR_CALL_ID,

	ABNF_NAME_SIP_HDR_CONTACT_ID_S:          ABNF_NAME_SIP_HDR_CONTACT_ID,
	ABNF_NAME_SIP_HDR_CONTENT_ENCODING_S:    ABNF_NAME_SIP_HDR_CONTENT_ENCODING,
	ABNF_NAME_SIP_HDR_CONTENT_LENGTH_S:      ABNF_NAME_SIP_HDR_CONTENT_LENGTH,
	ABNF_NAME_SIP_HDR_CONTENT_TYPE_S:        ABNF_NAME_SIP_HDR_CONTENT_TYPE,
	ABNF_NAME_SIP_HDR_SUBJECTE_S:            ABNF_NAME_SIP_HDR_SUBJECTE,
	ABNF_NAME_SIP_HDR_SUPPORTED_S:           ABNF_NAME_SIP_HDR_SUPPORTED,
	ABNF_NAME_SIP_HDR_ALLOW_EVENTS_S:        ABNF_NAME_SIP_HDR_ALLOW_EVENTS,
	ABNF_NAME_SIP_HDR_EVENT_S:               ABNF_NAME_SIP_HDR_EVENT,
	ABNF_NAME_SIP_HDR_REFER_TO_S:            ABNF_NAME_SIP_HDR_REFER_TO,
	ABNF_NAME_SIP_HDR_ACCEPT_CONTACT_S:      ABNF_NAME_SIP_HDR_ACCEPT_CONTACT,
	ABNF_NAME_SIP_HDR_REJECT_CONTACT_S:      ABNF_NAME_SIP_HDR_REJECT_CONTACT,
	ABNF_NAME_SIP_HDR_REQUEST_DISPOSITION_S: ABNF_NAME_SIP_HDR_REQUEST_DISPOSITION,
	ABNF_NAME_SIP_HDR_REFERRED_BY_S:         ABNF_NAME_SIP_HDR_REFERRED_BY,
	ABNF_NAME_SIP_HDR_SESSION_EXPIRES_S:     ABNF_NAME_SIP_HDR_SESSION_EXPIRES,
}

func GetSipHeaderFullName(name string) (fullName string) {
	name = strings.ToLower(name)
	fullName, ok := g_SipHeaderFullNameMaps[name]
	if !ok {
		return name
	}
	return fullName
}

type SipHeaderParsed interface {
	HasValue() bool
	String(context *ParseContext) string
	Encode(context *ParseContext, buf *bytes.Buffer)
}

type SipPaseOneHeaderValue func(context *ParseContext, src []byte, pos int) (newPos int, parsed AbnfPtr, err error)
type SipEncodeOneHeaderValue func(parsed AbnfPtr, context *ParseContext, buf *bytes.Buffer)

type SipHeaderInfo struct {
	name         []byte
	hasShortName bool
	shortName    []byte
	allowMulti   bool
	needParse    bool
	parseFunc    SipPaseOneHeaderValue
	encodeFunc   SipEncodeOneHeaderValue
}

func (this *SipHeaderInfo) AllowMulti() bool   { return this.allowMulti }
func (this *SipHeaderInfo) HasShortName() bool { return this.hasShortName }
func (this *SipHeaderInfo) ShortName() []byte  { return this.shortName }

var g_SipHeaderInfoMaps = map[string]SipHeaderInfo{
	"from":    {name: []byte("From"), hasShortName: true, shortName: []byte("f"), needParse: true, parseFunc: ParseSipFrom, encodeFunc: EncodeSipFromValue},
	"to":      {name: []byte("To"), hasShortName: true, shortName: []byte("t"), needParse: true, parseFunc: ParseSipTo, encodeFunc: EncodeSipToValue},
	"via":     {name: []byte("Via"), hasShortName: true, shortName: []byte("v"), allowMulti: true, needParse: true, parseFunc: ParseSipVia, encodeFunc: EncodeSipViaValue},
	"call-id": {name: []byte("Call-ID"), hasShortName: true, shortName: []byte("i"), needParse: true, parseFunc: ParseSipCallId, encodeFunc: EncodeSipCallIdValue},
	"cseq":    {name: []byte("CSeq"), needParse: true, parseFunc: ParseSipCseq, encodeFunc: EncodeSipCseqValue},

	"allow":               {name: []byte("Allow"), allowMulti: true},
	"contact":             {name: []byte("Contact"), hasShortName: true, shortName: []byte("m"), allowMulti: true, needParse: true, parseFunc: ParseSipContact, encodeFunc: EncodeSipContactValue},
	"content-disposition": {name: []byte("Content-Disposition"), needParse: true, parseFunc: ParseSipContentDisposition, encodeFunc: EncodeSipContentDispositionValue},
	"content-encoding":    {name: []byte("Content-Encoding"), hasShortName: true, shortName: []byte("e"), allowMulti: true},
	"content-length":      {name: []byte("Content-Length"), hasShortName: true, shortName: []byte("l"), needParse: true, parseFunc: ParseSipContentLength, encodeFunc: EncodeSipContentLengthValue},
	"content-type":        {name: []byte("Content-Type"), hasShortName: true, shortName: []byte("c"), needParse: true, parseFunc: ParseSipContentType, encodeFunc: EncodeSipContentTypeValue},

	"date": {name: []byte("Date")},

	"max-forwards": {name: []byte("Max-Forwards"), needParse: true, parseFunc: ParseSipMaxForwards, encodeFunc: EncodeSipMaxForwardsValue},

	"record-route": {name: []byte("Record-Route"), allowMulti: true, needParse: true, parseFunc: ParseSipRecordRoute, encodeFunc: EncodeSipRecordRouteValue},
	"route":        {name: []byte("Route"), allowMulti: true, needParse: true, parseFunc: ParseSipRoute, encodeFunc: EncodeSipRouteValue},

	"subject":   {name: []byte("Subject"), hasShortName: true, shortName: []byte("s")},
	"supported": {name: []byte("Supported"), hasShortName: true, shortName: []byte("k"), allowMulti: true},

	"allow-events": {name: []byte("Allow-Events"), hasShortName: true, shortName: []byte("u")},
	"event":        {name: []byte("Event"), hasShortName: true, shortName: []byte("o")},

	"refer-to":            {name: []byte("Refer-To"), hasShortName: true, shortName: []byte("r")},
	"accept-contact":      {name: []byte("Accept-Contact"), hasShortName: true, shortName: []byte("a"), allowMulti: true},
	"reject-contact":      {name: []byte("Reject-Contact"), hasShortName: true, shortName: []byte("j"), allowMulti: true},
	"request-disposition": {name: []byte("Request-Disposition"), hasShortName: true, shortName: []byte("d"), allowMulti: true},

	"referred-by":     {name: []byte("Referred-By"), hasShortName: true, shortName: []byte("b")},
	"session-expires": {name: []byte("Session-Expires"), hasShortName: true, shortName: []byte("x")},

	"mime-version": {name: []byte("MIME-Version")},
}

func GetSipHeaderInfo(name string) (info *SipHeaderInfo, ok bool) {
	name = strings.ToLower(name)
	info1, ok := g_SipHeaderInfoMaps[name]
	if !ok {
		return nil, false
	}
	return &info1, true
}

/*
type SipHeaderType interface {
	AllowMulti() bool
	//GetHeader() SipHeader
	//Parse(context, src []byte, pos int) (newPos int, err error)
	String() string
	Encode(buf *bytes.Buffer)
	HasInfo() bool
	Info() *SipHeaderInfo
}
*/

type SipHeaders struct {
	singleHeaders  SipSingleHeaders
	multiHeaders   SipMultiHeaders
	unknownHeaders SipSingleHeaders
}

func NewSipHeaders() *SipHeaders {
	ret := &SipHeaders{}
	ret.Init()
	return ret
}

func (this *SipHeaders) Init() {
	this.singleHeaders.Init()
	this.multiHeaders.Init()
	this.unknownHeaders.Init()
}

func (this *SipHeaders) Size() int32 {
	return this.singleHeaders.Size() + this.multiHeaders.Size() + this.unknownHeaders.Size()
}
func (this *SipHeaders) Empty() bool { return this.Size() == 0 }

func (this *SipHeaders) Encode(context *ParseContext, buf *bytes.Buffer) {
	this.singleHeaders.Encode(context, buf)
	this.multiHeaders.Encode(context, buf)
	this.unknownHeaders.Encode(context, buf)
}

func (this *SipHeaders) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

func (this *SipHeaders) GetSingleHeader(context *ParseContext, name string) (header *SipSingleHeader, ok bool) {
	return this.singleHeaders.GetHeaderByString(context, name)
}

func (this *SipHeaders) GetSingleHeaderParsed(context *ParseContext, name string) (parsed AbnfPtr, ok bool) {
	return this.singleHeaders.GetHeaderParsedByString(context, name)
}

func (this *SipHeaders) GetMultiHeader(context *ParseContext, name string) (header *SipMultiHeader, ok bool) {
	return this.multiHeaders.GetHeaderByString(context, name)
}

func (this *SipHeaders) GetUnknownHeader(context *ParseContext, name string) (header *SipSingleHeader, ok bool) {
	return this.unknownHeaders.GetHeaderByString(context, name)
}

func (this *SipHeaders) GenerateAndAddSingleHeader(context *ParseContext, name, value string) (*SipSingleHeader, AbnfPtr) {
	return this.singleHeaders.GenerateAndAddHeader(context, name, value)
}

func (this *SipHeaders) GenerateAndAddMultiHeader(context *ParseContext, name, value string) (*SipSingleHeader, AbnfPtr) {
	return this.multiHeaders.GenerateAndAddHeader(context, name, value)
}

func (this *SipHeaders) GenerateAndAddUnknownHeader(context *ParseContext, name, value string) (*SipSingleHeader, AbnfPtr) {
	return this.unknownHeaders.GenerateAndAddHeader(context, name, value)
}

func (this *SipHeaders) CreateSingleHeader(context *ParseContext, name string) (*SipSingleHeader, AbnfPtr) {
	header, addr := NewSipSingleHeader(context)
	if header == nil {
		return nil, ABNF_PTR_NIL
	}
	info, _ := GetSipHeaderInfo(name)
	header.info = info
	return header, addr
}

// remove Content-* headers from sip message except Content-Length and Content-Type*/
func (this *SipHeaders) RemoveContentHeaders(context *ParseContext) {
	this.singleHeaders.RemoveContentHeaders(context)
	this.multiHeaders.RemoveContentHeaders(context)
	this.unknownHeaders.RemoveContentHeaders(context)
}

func (this *SipHeaders) CopyContentHeaders(context *ParseContext, rhs *SipHeaders) {
	this.singleHeaders.CopyContentHeaders(context, &rhs.singleHeaders)
	this.multiHeaders.CopyContentHeaders(context, &rhs.multiHeaders)
	this.unknownHeaders.CopyContentHeaders(context, &rhs.unknownHeaders)
}

func (this *SipHeaders) CreateContentLength(context *ParseContext, size uint32) error {
	headerPtr, ok := this.GetSingleHeaderParsed(context, ABNF_NAME_SIP_HDR_CONTENT_LENGTH)
	if ok {
		headerPtr.GetSipHeaderContentLength(context).size = size
		return nil
	}

	contentLength, addr := this.CreateSingleHeader(context, ABNF_NAME_SIP_HDR_CONTENT_LENGTH)
	if contentLength == nil {
		return &AbnfError{"SipHeaders: out of memory for creating Content-Length", nil, 0}
	}

	parsedContentLength, parsedPtr := NewSipHeaderContentLength(context)
	if parsedContentLength == nil {
		return &AbnfError{"SipHeaders  encode: out of memory for creating parsed Content-Length", nil, 0}
	}
	parsedContentLength.size = size
	contentLength.parsed = parsedPtr
	this.singleHeaders.AddHeader(context, addr)
	return nil
}

func (this *SipHeaders) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos

	for newPos < len(src) {
		if IsCRLF(src, newPos) {
			/* reach message-body */
			return newPos + 2, nil
		}
		var name AbnfRef

		name, newPos, err = ParseHeaderName(context, src, newPos)
		if err != nil {
			return newPos, err
		}

		total_headers++

		fullName := GetSipHeaderFullName(ByteSliceToString(src[name.Begin:name.End]))

		info, ok := GetSipHeaderInfo(fullName)
		//ok = false
		if ok {
			newPos, err = this.parseKnownHeader(context, src, newPos, info)
		} else {
			/* unknown header */
			newPos, err = this.parseUnknownHeader(context, name, src, newPos, info)
		}

		if err != nil {
			return newPos, nil
		}

	}
	return newPos, nil
}

func (this *SipHeaders) parseKnownHeader(context *ParseContext, src []byte, pos int, info *SipHeaderInfo) (newPos int, err error) {
	if !info.allowMulti {
		return this.parseSingleKnownHeader(context, src, pos, info)
	}
	return this.parseMultiKnownHeader(context, src, pos, info)
}

func (this *SipHeaders) parseUnknownHeader(context *ParseContext, name AbnfRef, src []byte, pos int, info *SipHeaderInfo) (newPos int, err error) {
	addr, newPos, err := parseOneUnparsableSingleHeader(context, name, src, pos, info)
	if err != nil {
		return newPos, err
	}
	this.unknownHeaders.AddHeader(context, addr)
	return newPos, nil
}

func (this *SipHeaders) parseSingleKnownHeader(context *ParseContext, src []byte, pos int, info *SipHeaderInfo) (newPos int, err error) {
	var addr AbnfPtr

	newPos = pos
	_, ok := this.singleHeaders.GetHeaderByByteSlice(context, info.name)
	if ok {
		// discard this header
		_, newPos, err = ParseHeaderValue(context, src, newPos)
		return newPos, err
	}

	if info.parseFunc != nil && info.needParse {
		addr, newPos, err = parseOneParsableSingleHeader(context, src, newPos, info)
		if err != nil {
			return newPos, err
		}
		this.singleHeaders.AddHeader(context, addr)
		return ParseCRLF(src, newPos)
	} else {
		addr, newPos, err = parseOneUnparsableSingleHeader(context, AbnfRef{0, int32(len(info.name))}, src, newPos, info)
		if err != nil {
			return newPos, err
		}
		this.singleHeaders.AddHeader(context, addr)
	}

	return newPos, nil
}

func (this *SipHeaders) parseMultiKnownHeader(context *ParseContext, src []byte, pos int, info *SipHeaderInfo) (newPos int, err error) {
	var multiHeader *SipMultiHeader
	//var addr AbnfPtr

	newPos = pos
	multiHeader, ok := this.multiHeaders.GetHeaderByByteSlice(context, info.name)
	if !ok {
		var addr AbnfPtr
		multiHeader, addr = NewSipMultiHeader(context)
		if multiHeader == nil {
			return newPos, &AbnfError{"SipHeaders  parse: out of memory for known multi headers", src, newPos}
		}
		multiHeader.info = info
		multiHeader.SetNameByteSlice(context, info.name)
		this.multiHeaders.AddHeader(context, addr)
	}

	return multiHeader.Parse(context, src, newPos, info)
}

func ParseHeaderName(context *ParseContext, src []byte, pos int) (name AbnfRef, newPos int, err error) {
	newPos = name.Parse(src, pos, IsSipToken)
	if name.End <= name.Begin {
		return name, newPos, &AbnfError{"SipHeaders parse: no header-name", src, newPos}
	}

	newPos, err = ParseHcolon(src, newPos)
	return name, newPos, err
}

func ParseHeaderValue(context *ParseContext, src []byte, pos int) (value AbnfBuf, newPos int, err error) {
	newPos = pos
	begin, end, ok := FindCrlfRFC3261(src, newPos)
	if !ok {
		return value, end, &AbnfError{"SipHeaders parse: no CRLF for one header", src, newPos}
	}

	if begin > newPos {
		value.SetValue(context, src[newPos:begin])
	}

	return value, end, err
}

func FindCrlfRFC3261(src []byte, pos int) (begin, end int, ok bool) {
	/* state diagram
	 *                                                              other char/found
	 *       |----------|    CR    |-------|    LF    |---------|---------------------->end
	 *  |--->| ST_START | -------> | ST_CR |--------->| ST_CRLF |                        ^
	 *  |    |----------|          |-------|          |---------|                        |
	 *  |                               |                  |        other char/not found |
	 *  |                               |------------------+-----------------------------|
	 *  |            WSP                                   |
	 *  |--------------------------------------------------|
	 *
	 *  it is an error if any character except 'LF' is after 'CR' in this routine.
	 *  'CR' or 'LF' is not equal to 'CRLF' in this routine
	 */
	end = pos
	for end < len(src) {
		for ; (end < len(src)) && (src[end] != '\n'); end++ {
		}
		if end >= len(src) {
			/* no CRLF" */
			return end, end, false
		}
		end++

		if end >= len(src) {
			break
		}

		if !IsWspChar(src[end]) {
			break
		}
	}

	if ((pos + 2) < end) && (src[end-2] == '\r') {
		begin = end - 2
	} else {
		begin = end - 1
	}

	return begin, end, true
}
