package sipparser

import (
	"bytes"
	//"fmt"
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

type SipHeaderInfo struct {
	name         []byte
	hasShortName bool
	shortName    []byte
	allowMulti   bool
	needParse    bool
	parseFunc    SipPaseOneHeaderValue
}

func (this *SipHeaderInfo) AllowMulti() bool   { return this.allowMulti }
func (this *SipHeaderInfo) HasShortName() bool { return this.hasShortName }
func (this *SipHeaderInfo) ShortName() []byte  { return this.shortName }

var g_SipHeaderInfoMaps = map[string]SipHeaderInfo{
	"From":    {name: []byte("From"), hasShortName: true, shortName: []byte("f"), needParse: true, parseFunc: ParseSipFrom},
	"To":      {name: []byte("To"), hasShortName: true, shortName: []byte("t"), needParse: true, parseFunc: ParseSipTo},
	"Via":     {name: []byte("Via"), hasShortName: true, shortName: []byte("v"), allowMulti: true, needParse: true, parseFunc: ParseSipVia},
	"Call-ID": {name: []byte("Call-ID"), hasShortName: true, shortName: []byte("i"), needParse: true, parseFunc: ParseSipCallId},
	"CSeq":    {name: []byte("CSeq"), needParse: true, parseFunc: ParseSipCseq},

	"Allow":               {name: []byte("Allow"), allowMulti: true},
	"Contact":             {name: []byte("Contact"), hasShortName: true, shortName: []byte("m"), allowMulti: true, needParse: true, parseFunc: ParseSipContact},
	"Content-Disposition": {name: []byte("Content-Disposition"), needParse: true, parseFunc: ParseSipContentDisposition},
	"Content-Encoding":    {name: []byte("Content-Encoding"), hasShortName: true, shortName: []byte("e"), allowMulti: true},
	"Content-Length":      {name: []byte("Content-Length"), hasShortName: true, shortName: []byte("l"), needParse: true, parseFunc: ParseSipContentLength},
	"Content-Type":        {name: []byte("Content-Type"), hasShortName: true, shortName: []byte("c"), needParse: true, parseFunc: ParseSipContentType},

	"Date": {name: []byte("Date")},

	"Max-Forwards": {name: []byte("Max-Forwards"), needParse: true, parseFunc: ParseSipMaxForwards},

	"Record-Route": {name: []byte("Record-Route"), allowMulti: true, needParse: true, parseFunc: ParseSipRecordRoute},
	"Route":        {name: []byte("Route"), allowMulti: true, needParse: true, parseFunc: ParseSipRoute},

	"Subject":   {name: []byte("Subject"), hasShortName: true, shortName: []byte("s")},
	"Supported": {name: []byte("Supported"), hasShortName: true, shortName: []byte("k"), allowMulti: true},

	"Allow-Events": {name: []byte("Allow-Events"), hasShortName: true, shortName: []byte("u")},
	"Event":        {name: []byte("Event"), hasShortName: true, shortName: []byte("o")},

	"Refer-To":            {name: []byte("Refer-To"), hasShortName: true, shortName: []byte("r")},
	"Accept-Contact":      {name: []byte("Accept-Contact"), hasShortName: true, shortName: []byte("a"), allowMulti: true},
	"Reject-Contact":      {name: []byte("Reject-Contact"), hasShortName: true, shortName: []byte("j"), allowMulti: true},
	"Request-Disposition": {name: []byte("Request-Disposition"), hasShortName: true, shortName: []byte("d"), allowMulti: true},

	"Referred-By":     {name: []byte("Referred-By"), hasShortName: true, shortName: []byte("b")},
	"Session-Expires": {name: []byte("Session-Expires"), hasShortName: true, shortName: []byte("x")},
}

func GetSipHeaderInfo(name string) (info *SipHeaderInfo, ok bool) {
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

func (this *SipHeaders) GetMultiHeader(context *ParseContext, name string) (header *SipMultiHeader, ok bool) {
	return this.multiHeaders.GetHeaderByString(context, name)
}

func (this *SipHeaders) GetUnknownHeader(context *ParseContext, name string) (header *SipSingleHeader, ok bool) {
	return this.unknownHeaders.GetHeaderByString(context, name)
}

func (this *SipHeaders) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos

	for newPos < len(src) {
		if src[newPos] == '\r' && ((newPos + 1) < len(src)) && (src[newPos+1] == '\n') {
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
			if !info.allowMulti {
				newPos, err = this.ParseSingleKnownHeader(context, src, newPos, info)
			} else {
				newPos, err = this.ParseMultiKnownHeader(context, src, newPos, info)
			}
			if err != nil {
				return newPos, nil
			}
		} else {
			/* unknown or unporecessed header */
			newPos, err = this.ParseUnprocessedHeader(context, name, src, newPos, info)
			if err != nil {
				return newPos, nil
			}
		}

	}
	return newPos, nil
}

func (this *SipHeaders) ParseUnprocessedHeader(context *ParseContext, name AbnfRef, src []byte, pos int, info *SipHeaderInfo) (newPos int, err error) {
	newPos = pos
	header, addr := NewSipSingleHeader(context)
	if header == nil {
		return newPos, &AbnfError{"SipHeaders  parse: out of memory for unknown headers", src, newPos}
	}
	header.name.SetExist()
	header.name.SetValue(context, src[name.Begin:name.End])
	header.value, newPos, err = ParseHeaderValue(context, src, newPos)
	if err != nil {
		return newPos, err
	}
	this.unknownHeaders.AddHeader(context, addr)
	return newPos, nil
}

func (this *SipHeaders) ParseSingleKnownHeader(context *ParseContext, src []byte, pos int, info *SipHeaderInfo) (newPos int, err error) {
	newPos = pos
	//*
	_, ok := this.singleHeaders.GetHeaderByBytes(context, info.name)
	if ok {
		// discard this header
		_, newPos, err = ParseHeaderValue(context, src, newPos)
		return newPos, err
	} //*/

	if info.parseFunc != nil && info.needParse {
		begin := newPos
		newPos, parsed, err := info.parseFunc(context, src, newPos)
		//newPos, parsed, err := ParseSipContentLength(context, src, newPos)
		if err != nil {
			return newPos, err
		}

		newHeader, addr := NewSipSingleHeader(context)
		if newHeader == nil {
			return newPos, &AbnfError{"SipHeaders  parse: out of memory for known single headers", src, newPos}
		}
		newHeader.info = info
		newHeader.parsed = parsed
		newHeader.name.SetExist()
		newHeader.name.SetValue(context, info.name)
		if newPos > begin {
			newHeader.value.SetExist()
			newHeader.value.SetValue(context, src[begin:newPos])
		}
		this.singleHeaders.AddHeader(context, addr)
		return ParseCRLF(src, newPos)
	} else {
		newHeader, addr := NewSipSingleHeader(context)
		if newHeader == nil {
			return newPos, &AbnfError{"SipHeaders  parse: out of memory for known unprocessed single headers", src, newPos}
		}
		newHeader.info = info
		newHeader.name.SetExist()
		newHeader.name.SetValue(context, info.name)
		newHeader.value, newPos, err = ParseHeaderValue(context, src, newPos)
		if err != nil {
			return newPos, err
		}
		this.singleHeaders.AddHeader(context, addr)
	}

	return newPos, nil
}

func (this *SipHeaders) ParseMultiKnownHeader(context *ParseContext, src []byte, pos int, info *SipHeaderInfo) (newPos int, err error) {
	var multiHeader *SipMultiHeader

	newPos = pos
	multiHeader, ok := this.multiHeaders.GetHeaderByBytes(context, info.name)
	if !ok {
		var addr AbnfPtr
		multiHeader, addr = NewSipMultiHeader(context)
		if multiHeader == nil {
			return newPos, &AbnfError{"SipHeaders  parse: out of memory for known multi headers", src, newPos}
		}
		multiHeader.info = info
		multiHeader.name.SetExist()
		multiHeader.name.SetValue(context, info.name)
		this.multiHeaders.AddHeader(context, addr)
	}

	//*
	if info.parseFunc != nil && info.needParse {
		for newPos < len(src) {

			begin := newPos
			newPos, parsed, err := info.parseFunc(context, src, newPos)
			if err != nil {
				return newPos, err
			}

			newHeader, addr := NewSipSingleHeader(context)
			if newHeader == nil {
				return newPos, &AbnfError{"SipHeaders  parse: out of memory for known multi headers", src, newPos}
			}
			newHeader.info = info
			newHeader.parsed = parsed
			newHeader.name.SetExist()
			newHeader.name.SetValue(context, info.name)
			if newPos > begin {
				newHeader.value.SetExist()
				newHeader.value.SetValue(context, src[begin:newPos])
			}
			multiHeader.AddHeader(context, addr)

			// now should be COMMA or CRLF
			newPos1, err := ParseSWSMark(src, newPos, ',')
			if err != nil {
				// should be CRLF
				return ParseCRLF(src, newPos)
			}
			newPos = newPos1

		}
	} else {
		// here one SipMultiHeader main contain some headers of same type seperated by COMMA
		newHeader, addr := NewSipSingleHeader(context)
		if newHeader == nil {
			return newPos, &AbnfError{"SipHeaders  parse: out of memory for known unprocessed multi headers", src, newPos}
		}
		newHeader.info = info
		newHeader.name.SetExist()
		newHeader.name.SetValue(context, info.name)
		newHeader.value, newPos, err = ParseHeaderValue(context, src, newPos)
		if err != nil {
			return newPos, err
		}
		multiHeader.AddHeader(context, addr)
	}

	return newPos, nil
}

func ParseHeaderName(context *ParseContext, src []byte, pos int) (name AbnfRef, newPos int, err error) {
	newPos = name.Parse(src, pos, IsSipToken)
	if name.End <= name.Begin {
		return name, newPos, &AbnfError{"SipHeaders parse: no header-name", src, newPos}
	}

	newPos, err = ParseHcolon(src, newPos)
	return name, newPos, err
}

func ParseHeaderValue(context *ParseContext, src []byte, pos int) (value AbnfToken, newPos int, err error) {
	newPos = pos
	begin, end, ok := FindCrlfRFC3261(src, newPos)
	if !ok {
		return value, end, &AbnfError{"SipHeaders parse: no CRLF for one header", src, newPos}
	}

	if begin > newPos {
		value.SetExist()
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
