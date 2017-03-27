package sipparser3

import (
	"bytes"
	"strconv"
	//"fmt"
	//"strings"
)

type SipStartLine struct {
	isRequest    bool
	method       AbnfToken
	version      SipVersion
	uri          SipUri
	statusCode   uint16
	reasonPhrase AbnfToken
}

func NewSipStartLine() *SipStartLine {
	ret := &SipStartLine{}
	ret.Init()
	return ret
}

func (this *SipStartLine) Init() {
	this.uri.Init()
}

func (this *SipStartLine) AllowMulti() bool { return false }
func (this *SipStartLine) HasValue() bool   { return true }

/* RFC3261 Section 25.1, page 222
 *
 * Request-Line   =  Method SP Request-URI SP SIP-Version CRLF
 * Status-Line    =  SIP-Version SP Status-Code SP Reason-Phrase CRLF
 * SIP-Version    =  "SIP" "/" 1*DIGIT "." 1*DIGIT
 */
func (this *SipStartLine) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	name, newPos, err := ParseHeaderName(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if !name.EqualStringNoCase("CSeq") {
		return newPos, &AbnfError{"CSeq parse: wrong header-name", src, newPos}
	}

	return this.ParseValue(context, src, newPos)
}

func (this *SipStartLine) ParseValue(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos
	digit, _, newPos, ok := ParseUInt(src, newPos)
	if !ok {
		return newPos, &AbnfError{"CSeq parse: wrong num", src, newPos}
	}

	this.id = uint32(digit)

	newPos, err = ParseLWS(src, newPos)
	if err != nil {
		return newPos, err
	}

	return this.method.Parse(context, src, newPos, IsSipToken)
}

func (this *SipStartLine) Encode(buf *bytes.Buffer) {
	buf.WriteString("CSeq: ")
	buf.WriteString(strconv.FormatUint(uint64(this.id), 10))
	buf.WriteByte(' ')
	this.method.Encode(buf)
}

func (this *SipStartLine) String() string {
	return AbnfEncoderToString(this)
}
