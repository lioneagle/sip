package sipparser3

import (
	"bytes"
	//"fmt"
)

type SipHeaderInfo struct {
	name       []byte
	shortName  AbnfToken
	allowMulti bool
}

var g_SipHeaderShortNameToFullNameMaps = map[string]string{
	"f": "From",
	"t": "To",
	"v": "Via",
	"i": "Call-ID",
}

var g_SipHeaderInfoMaps = map[string]SipHeaderInfo{
	"From":    {name: Str2bytes("From"), shortName: AbnfToken{true, Str2bytes("f")}, allowMulti: false},
	"To":      {name: Str2bytes("To"), shortName: AbnfToken{true, Str2bytes("t")}, allowMulti: false},
	"Via":     {name: Str2bytes("Via"), shortName: AbnfToken{true, Str2bytes("v")}, allowMulti: true},
	"Call-ID": {name: Str2bytes("Call-ID"), shortName: AbnfToken{true, Str2bytes("i")}, allowMulti: false},
	"CSeq":    {name: Str2bytes("Call-ID"), shortName: AbnfToken{false, nil}, allowMulti: false},
}

func GetSipHeaderInfo(name string) (info *SipHeaderInfo, ok bool) {
	info1, ok := g_SipHeaderInfoMaps[name]
	if !ok {
		return nil, false
	}
	return &info1, ok
}

type SipHeaderParsed interface {
	IsHeaderList() bool
	Parse(context, src []byte, pos int) (newPos int, err error)
	String() string
	Encode(buf *bytes.Buffer)
}

type SipHeader struct {
	info   *SipHeaderInfo
	name   AbnfToken
	value  AbnfToken
	parsed SipHeaderParsed
}

func (this *SipHeader) Encode(buf *bytes.Buffer) {
	this.name.Encode(buf)
	buf.WriteByte(':')
	if this.value.Exist() {
		this.value.Encode(buf)
	}
}

func (this *SipHeader) String() string {
	var buf bytes.Buffer
	this.Encode(&buf)
	return buf.String()
}

type SipHeaders struct {
	headers []SipHeader
}

func NewSipHeaders() *SipHeaders {
	ret := &SipHeaders{}
	ret.Init()
	return ret
}

func (this *SipHeaders) Init() {
	this.headers = make([]SipHeader, 0, 5)
}

func (this *SipHeaders) Size() int   { return len(this.headers) }
func (this *SipHeaders) Empty() bool { return len(this.headers) == 0 }

func (this *SipHeaders) Encode(buf *bytes.Buffer) {
	for _, v := range this.headers {
		v.Encode(buf)
		buf.WriteByte('\r')
		buf.WriteByte('\n')
	}
}

func (this *SipHeaders) String() string {
	var buf bytes.Buffer
	this.Encode(&buf)
	return buf.String()
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
		name.SetExist()

		info, ok := GetSipHeaderInfo(Bytes2str(name.value))
		if ok {
			if !info.allowMulti {
				this.ParseSingleKnownHeader(context, src, newPos, info)
			}

		} else {

			/* unknown header */
			header := SipHeader{name: name}
			header.value, newPos, err = ParseHeaderValue(context, src, newPos)
			if err != nil {
				return newPos, err
			}
			this.headers = append(this.headers, header)
		}

	}
	return newPos, nil
}

func (this *SipHeaders) ParseSingleKnownHeader(context *ParseContext, src []byte, pos int, info *SipHeaderInfo) (name AbnfToken, newPos int, err error) {
	return name, newPos, nil
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
