package sipparser3

import (
	"bytes"
	//"fmt"
)

var g_SipHeaderFullNameMaps = map[string]string{
	"f": "From",
	"t": "To",
	"v": "Via",
	"i": "Call-ID",

	"m": "Contact",
	"e": "Content-Encoding",
	"l": "Content-Length",
	"c": "Content-Type",
	"s": "Subject",
	"k": "Supported",
	"u": "Allow-Events",
	"o": "Event",
	"r": "Refer-To",
	"a": "Accept-Contact",
	"j": "Reject-Contact",
	"d": "Request-Disposition",
	"b": "Referred-By",
	"x": "Session-Expires",
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
	String() string
	Encode(buf *bytes.Buffer)
}

type SipPaseOneHeaderValue func(context *ParseContext, src []byte, pos int) (newPos int, parsed SipHeaderParsed, err error)

type SipHeaderInfo struct {
	name       []byte
	shortName  AbnfToken
	allowMulti bool
	needParse  bool
	parseFunc  SipPaseOneHeaderValue
}

func (this *SipHeaderInfo) HasShortName() bool   { return this.shortName.Exist() }
func (this *SipHeaderInfo) AllowMulti() bool     { return this.allowMulti }
func (this *SipHeaderInfo) ShortName() AbnfToken { return this.shortName }

var g_SipHeaderInfoMaps = map[string]SipHeaderInfo{
	"From":    {name: Str2bytes("From"), shortName: AbnfToken{true, Str2bytes("f")}, allowMulti: false, needParse: true, parseFunc: ParseSipFrom},
	"To":      {name: Str2bytes("To"), shortName: AbnfToken{true, Str2bytes("t")}, allowMulti: false, needParse: true, parseFunc: ParseSipTo},
	"Via":     {name: Str2bytes("Via"), shortName: AbnfToken{true, Str2bytes("v")}, allowMulti: true, needParse: true, parseFunc: ParseSipVia},
	"Call-ID": {name: Str2bytes("Call-ID"), shortName: AbnfToken{true, Str2bytes("i")}, allowMulti: false, needParse: true, parseFunc: ParseSipCallId},
	"CSeq":    {name: Str2bytes("CSeq"), allowMulti: false, needParse: true, parseFunc: ParseSipCseq},

	"Allow":            {name: Str2bytes("Allow"), allowMulti: true},
	"Contact":          {name: Str2bytes("Contact"), shortName: AbnfToken{true, Str2bytes("m")}, allowMulti: true},
	"Content-Encoding": {name: Str2bytes("Content-Encoding"), shortName: AbnfToken{true, Str2bytes("e")}, allowMulti: true},
	"Content-Length":   {name: Str2bytes("Content-Length"), shortName: AbnfToken{true, Str2bytes("l")}, allowMulti: false},
	"Content-Type":     {name: Str2bytes("Content-Type"), shortName: AbnfToken{true, Str2bytes("l")}, allowMulti: false},
	"Date":             {name: Str2bytes("Date"), allowMulti: false},

	"Subject":   {name: Str2bytes("Subject"), shortName: AbnfToken{true, Str2bytes("s")}, allowMulti: false},
	"Supported": {name: Str2bytes("Supported"), shortName: AbnfToken{true, Str2bytes("k")}, allowMulti: true},

	"Allow-Events": {name: Str2bytes("Allow-Events"), shortName: AbnfToken{true, Str2bytes("u")}, allowMulti: true},
	"Event":        {name: Str2bytes("Event"), shortName: AbnfToken{true, Str2bytes("o")}, allowMulti: false},

	"Refer-To":            {name: Str2bytes("Refer-To"), shortName: AbnfToken{true, Str2bytes("r")}, allowMulti: false},
	"Accept-Contact":      {name: Str2bytes("Accept-Contact"), shortName: AbnfToken{true, Str2bytes("a")}, allowMulti: true},
	"Reject-Contact":      {name: Str2bytes("Reject-Contact"), shortName: AbnfToken{true, Str2bytes("j")}, allowMulti: true},
	"Request-Disposition": {name: Str2bytes("Request-Disposition"), shortName: AbnfToken{true, Str2bytes("d")}, allowMulti: true},

	"Referred-By":     {name: Str2bytes("Referred-By"), shortName: AbnfToken{true, Str2bytes("b")}, allowMulti: false},
	"Session-Expires": {name: Str2bytes("Session-Expires"), shortName: AbnfToken{true, Str2bytes("xvchaogaosuc")}, allowMulti: false},
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

func (this *SipHeaders) Size() int {
	return this.singleHeaders.Size() + this.multiHeaders.Size() + this.unknownHeaders.Size()
}
func (this *SipHeaders) Empty() bool { return this.Size() == 0 }

func (this *SipHeaders) Encode(buf *bytes.Buffer) {
	this.singleHeaders.Encode(buf)
	this.multiHeaders.Encode(buf)
	this.unknownHeaders.Encode(buf)
}

func (this *SipHeaders) String() string {
	return AbnfEncoderToString(this)
}

func (this *SipHeaders) GetSingleHeader(name string) (header *SipSingleHeader, ok bool) {
	return this.singleHeaders.GetHeaderString(name)
}

func (this *SipHeaders) GetMultiHeader(name string) (header *SipMultiHeader, ok bool) {
	return this.multiHeaders.GetHeaderString(name)
}

func (this *SipHeaders) GetUnknownHeader(name string) (header *SipSingleHeader, ok bool) {
	return this.unknownHeaders.GetHeaderString(name)
}

func (this *SipHeaders) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos

	for newPos < len(src) {
		if src[newPos] == '\r' && ((newPos + 1) < len(src)) && (src[newPos+1] == '\n') {
			/* reach message-body */
			return newPos + 2, nil
		}
		var name AbnfToken

		name, newPos, err = ParseHeaderName(context, src, newPos)
		if err != nil {
			return newPos, err
		}

		fullName := GetSipHeaderFullName(name.String())

		info, ok := GetSipHeaderInfo(fullName)
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
			header := SipSingleHeader{name: name}
			header.value, newPos, err = ParseHeaderValue(context, src, newPos)
			if err != nil {
				return newPos, err
			}
			this.unknownHeaders.AddHeader(&header)
		}

	}
	return newPos, nil
}

func (this *SipHeaders) ParseSingleKnownHeader(context *ParseContext, src []byte, pos int, info *SipHeaderInfo) (newPos int, err error) {
	newPos = pos
	_, ok := this.singleHeaders.GetHeaderBytes(info.name)
	if ok {
		/* discard this header */
		return newPos, nil
	}

	if info.parseFunc != nil && info.needParse {
		begin := newPos
		newPos1, parsed, err := info.parseFunc(context, src, newPos)
		newPos = newPos1
		if err != nil {
			return newPos, err
		}

		newHeader := SipSingleHeader{name: AbnfToken{true, info.name}, info: info, parsed: parsed}
		newHeader.Init()
		if newPos > begin {
			newHeader.value.SetExist()
			newHeader.value.SetValue(src[begin:newPos])
		}
		this.singleHeaders.AddHeader(&newHeader)
		return ParseCRLF(src, newPos)
	} else {
		newHeader := SipSingleHeader{name: AbnfToken{true, info.name}, info: info}
		newHeader.Init()
		newHeader.value, newPos, err = ParseHeaderValue(context, src, newPos)
		if err != nil {
			return newPos, err
		}
		this.singleHeaders.AddHeader(&newHeader)
	}

	return newPos, nil
}

func (this *SipHeaders) ParseMultiKnownHeader(context *ParseContext, src []byte, pos int, info *SipHeaderInfo) (newPos int, err error) {
	var multiHeader *SipMultiHeader

	newPos = pos
	multiHeader, ok := this.multiHeaders.GetHeaderBytes(info.name)
	if !ok {
		multiHeader = &SipMultiHeader{name: AbnfToken{true, info.name}, info: info}
		multiHeader.Init()
		multiHeader = this.multiHeaders.AddHeader(multiHeader)
	}

	if info.parseFunc != nil && info.needParse {
		for newPos < len(src) {
			begin := newPos
			newPos1, parsed, err := info.parseFunc(context, src, newPos)
			newPos = newPos1
			if err != nil {
				return newPos, err
			}

			newHeader := SipSingleHeader{name: AbnfToken{true, info.name}, info: info, parsed: parsed}
			newHeader.Init()
			if newPos > begin {
				newHeader.value.SetExist()
				newHeader.value.SetValue(src[begin:newPos])
			}
			multiHeader.AddHeader(&newHeader)

			/* now should be COMMA or CRLF */
			newPos1, err = ParseSWSMark(src, newPos, ',')
			if err != nil {
				/* should be CRLF */
				return ParseCRLF(src, newPos)
			}
			newPos = newPos1

		}
	} else {
		/* here one SipMultiHeader main contain some headers of same type seperated by COMMA */
		newHeader := SipSingleHeader{name: AbnfToken{true, info.name}, info: info}
		newHeader.value, newPos, err = ParseHeaderValue(context, src, newPos)
		if err != nil {
			return newPos, err
		}
		multiHeader.AddHeader(&newHeader)
	}

	return newPos, nil
}

func ParseHeaderName(context *ParseContext, src []byte, pos int) (name AbnfToken, newPos int, err error) {
	newPos, err = name.Parse(context, src, pos, IsSipToken)
	if err != nil {
		return name, newPos, err
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
		value.SetValue(src[newPos:begin])
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
