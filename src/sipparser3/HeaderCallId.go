package sipparser3

import (
	"bytes"
	//"fmt"
	//"strings"
)

type SipHeaderCallId struct {
	id1 AbnfToken
	id2 AbnfToken
}

func NewSipHeaderCallId() *SipHeaderCallId {
	ret := &SipHeaderCallId{}
	ret.Init()
	return ret
}

func (this *SipHeaderCallId) Init() {
}

func (this *SipHeaderCallId) AllowMulti() bool { return false }
func (this *SipHeaderCallId) HasValue() bool   { return true }

/* RFC3261
 *
 * Call-ID  =  ( "Call-ID" / "i" ) HCOLON callid
 * callid   =  word [ "@" word ]
 */
func (this *SipHeaderCallId) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	name, newPos, err := ParseHeaderName(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if !name.EqualStringNoCase("Call-ID") && !name.EqualStringNoCase("i") {
		return newPos, &AbnfError{"Call-ID parse: wrong header-name", src, newPos}
	}

	return this.ParseValue(context, src, newPos)
}

func (this *SipHeaderCallId) ParseValue(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos
	newPos, err = this.id1.Parse(context, src, newPos, IsSipWord)
	if err != nil {
		return newPos, err
	}

	if newPos >= len(src) {
		return newPos, nil
	}

	if src[newPos] == '@' {
		return this.id2.Parse(context, src, newPos+1, IsSipWord)
	}

	return newPos, nil
}

func (this *SipHeaderCallId) Encode(buf *bytes.Buffer) {
	buf.WriteString("Call-ID: ")
	this.id1.Encode(buf)
	if this.id2.Exist() {
		buf.WriteByte('@')
		this.id2.Encode(buf)

	}
}

func (this *SipHeaderCallId) String() string {
	return AbnfEncoderToString(this)
}

func ParseSipCallId(context *ParseContext, src []byte, pos int) (newPos int, parsed SipHeaderParsed, err error) {
	header := SipHeaderCallId{}
	header.Init()

	newPos, err = header.ParseValue(context, src, pos)
	return newPos, &header, err
}
