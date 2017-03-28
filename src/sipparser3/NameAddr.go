package sipparser3

import (
	//"fmt"
	//"strings"
	"bytes"
)

type SipDisplayName struct {
	isQuotedString bool
	name           AbnfToken
	quotedstring   SipQuotedString
}

func NewSipDisplayName() *SipDisplayName {
	return &SipDisplayName{}
}

func (this *SipDisplayName) Encode(buf *bytes.Buffer) {
	if this.isQuotedString {
		this.quotedstring.Encode(buf)
	} else {
		this.name.Encode(buf)
	}
}

func (this *SipDisplayName) String() string {
	return AbnfEncoderToString(this)
}

func (this *SipDisplayName) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	/* RFC3261 Section 25.1, page 222
	 *
	 * display-name   =  *(token LWS)/ quoted-string
	 */
	newPos = pos
	if newPos >= len(src) {
		return newPos, nil
	}

	newPos, err = ParseSWS(src, newPos)
	if err != nil {
		return newPos, err
	}

	if newPos < len(src) && src[newPos] == '"' {
		this.isQuotedString = true
		newPos, err = this.quotedstring.Parse(context, src, newPos)
		if err != nil {
			return newPos, err
		}
	} else {
		this.isQuotedString = false
		newPos = pos
		if !IsSipToken(src[newPos]) {
			return newPos, &AbnfError{"DisplayName parse: no token or quoted-string", src, newPos}
		}

		nameBegin := newPos
		for newPos < len(src) {
			if !IsSipToken(src[newPos]) {
				break
			}

			token := AbnfToken{}
			newPos, err = token.Parse(context, src, newPos, IsSipToken)
			if err != nil {
				break
			}

			newPos, err = ParseLWS(src, newPos)
			if err != nil {
				return newPos, err
			}
		}
		this.name.SetExist()
		this.name.SetValue(src[nameBegin:newPos])
	}

	return newPos, nil
}

type SipNameAddr struct {
	displayname SipDisplayName
	addrsepc    SipAddrSpec
}

func NewSipNameAddr() *SipNameAddr {
	return &SipNameAddr{}
}

func (this *SipNameAddr) Encode(buf *bytes.Buffer) {
	this.displayname.Encode(buf)
	buf.WriteByte('<')
	this.addrsepc.Encode(buf)
	buf.WriteByte('>')
}

func (this *SipNameAddr) String() string {
	return AbnfEncoderToString(this)
}

func (this *SipNameAddr) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	/* RFC3261 Section 25.1, page 222
	 *
	 * name-addr =  [ display-name ] LAQUOT addr-spec RAQUOT
	 * RAQUOT  =  ">" SWS ; right angle quote
	 * LAQUOT  =  SWS "<"; left angle quote
	 */
	newPos = pos

	newPos, err = ParseSWS(src, newPos)
	if err != nil {
		return newPos, err
	}

	if newPos >= len(src) {
		return newPos, &AbnfError{"SipNameAddr parse: no value", src, newPos}
	}

	if src[newPos] != '<' {
		newPos, err = this.displayname.Parse(context, src, newPos)
		if err != nil {
			return newPos, err
		}
	}

	newPos, err = ParseLeftAngleQuote(src, newPos)
	if err != nil {
		return newPos, err
	}

	newPos, err = this.addrsepc.Parse(context, src, newPos)
	if err != nil {
		return newPos, err
	}

	newPos, err = ParseRightAngleQuote(src, newPos)
	if err != nil {
		return newPos, err
	}

	return newPos, nil
}
