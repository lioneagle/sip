package sipparser3

import (
	"bytes"
	//"fmt"
	//"strings"
)

type SipStartLine struct {
	isRequest    bool
	method       AbnfToken
	version      SipVersion
	addrspec     SipAddrSpec
	statusCode   uint16
	reasonPhrase AbnfToken
}

func NewSipStartLine() *SipStartLine {
	ret := &SipStartLine{}
	ret.Init()
	return ret
}

func (this *SipStartLine) Init() {
	this.isRequest = false
	this.method.SetNonExist()
	this.version.Init()
	this.addrspec.Init()
	this.reasonPhrase.SetNonExist()
}

func (this *SipStartLine) IsRequest() bool  { return this.isRequest }
func (this *SipStartLine) AllowMulti() bool { return false }
func (this *SipStartLine) HasValue() bool   { return true }

/* RFC3261 Section 25.1, page 222
 *
 * Request-Line   =  Method SP Request-URI SP SIP-Version CRLF
 * Status-Line    =  SIP-Version SP Status-Code SP Reason-Phrase CRLF
 * SIP-Version    =  "SIP" "/" 1*DIGIT "." 1*DIGIT
 * Reason-Phrase  =  *(reserved / unreserved / escaped
 *                   / UTF8-NONASCII / UTF8-CONT / SP / HTAB)
 * Request-URI    =  SIP-URI / SIPS-URI / absoluteURI
 */
func (this *SipStartLine) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	newPos, err = this.version.Parse(context, src, pos)
	if err == nil {
		this.isRequest = false
		newPos, err = this.ParseStatusLineAfterSipVersion(context, src, newPos)
	} else {
		this.isRequest = true
		newPos, err = this.ParseRequestLine(context, src, pos)
	}

	if err != nil {
		return newPos, err
	}

	if newPos >= len(src) {
		return newPos, &AbnfError{"StartLine parse: reach end before CRLF", src, newPos}
	}

	return ParseCRLF(src, newPos)
}

func (this *SipStartLine) ParseStatusLineAfterSipVersion(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos
	if newPos >= len(src) {
		return newPos, &AbnfError{"StatusLine parse: reach end after SIP-Version", src, newPos}
	}

	if src[newPos] != ' ' {
		return newPos, &AbnfError{"StatusLine parse: no SP after SIP-Version", src, newPos}
	}

	digit, _, newPos, ok := ParseUInt(src, newPos+1)
	if !ok {
		return newPos, &AbnfError{"StatusLine parse: wrong Status-Code", src, newPos}
	}

	this.statusCode = uint16(digit)

	if newPos >= len(src) {
		return newPos, &AbnfError{"StatusLine parse: reach end after Status-Code", src, newPos}
	}

	if src[newPos] != ' ' {
		return newPos, &AbnfError{"StatusLine parse: no SP after Status-Code", src, newPos}
	}

	newPos++

	if newPos >= len(src) {
		return newPos, nil
	}

	if !IsSipReasonPhrase(src[newPos]) {
		return newPos, nil
	}

	return this.reasonPhrase.ParseEscapable(context, src, newPos, IsSipReasonPhrase)
}

func (this *SipStartLine) ParseRequestLine(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos
	newPos, err = this.method.Parse(context, src, newPos, IsSipToken)
	if err != nil {
		return newPos, &AbnfError{"RequestLine parse: wrong METHOD", src, newPos}
	}

	if newPos >= len(src) {
		return newPos, &AbnfError{"RequestLine parse: reach end after METHOD", src, newPos}
	}

	if src[newPos] != ' ' {
		return newPos, &AbnfError{"RequestLine parse: no SP after METHOD", src, newPos}
	}

	newPos, err = this.addrspec.Parse(context, src, newPos+1)
	if err != nil {
		//return newPos, &AbnfError{"RequestLine parse: wrong Request-URI", src, newPos}
		return newPos, err
	}

	if newPos >= len(src) {
		return newPos, &AbnfError{"RequestLine parse: reach end after Request-URI", src, newPos}
	}

	if src[newPos] != ' ' {
		return newPos, &AbnfError{"RequestLine parse: no SP after Request-URI", src, newPos}
	}

	newPos, err = this.version.Parse(context, src, newPos+1)
	if err != nil {
		//return newPos, &AbnfError{"RequestLine parse: wrong SIP-Version", src, newPos}
		return newPos, err
	}

	return newPos, nil
}

func (this *SipStartLine) Encode(buf *bytes.Buffer) {
	if this.isRequest {
		this.EncodeRequestLine(buf)
	} else {
		this.EncodeStatusLine(buf)
	}
	buf.WriteByte('\r')
	buf.WriteByte('\n')
}

func (this *SipStartLine) EncodeRequestLine(buf *bytes.Buffer) {
	this.method.Encode(buf)
	buf.WriteByte(' ')
	this.addrspec.Encode(buf)
	buf.WriteByte(' ')
	this.version.Encode(buf)
}

func (this *SipStartLine) EncodeStatusLine(buf *bytes.Buffer) {
	this.version.Encode(buf)
	buf.WriteByte(' ')
	EncodeUInt(buf, uint64(this.statusCode))
	buf.WriteByte(' ')
	this.reasonPhrase.Encode(buf)
}

func (this *SipStartLine) String() string {
	return AbnfEncoderToString(this)
}
