package sipparser3

import (
	"bytes"
	//"fmt"
	//"strings"
)

type SipHeaderContentLength struct {
	size uint32
}

func NewSipHeaderContentLength() *SipHeaderContentLength {
	ret := &SipHeaderContentLength{}
	ret.Init()
	return ret
}

func (this *SipHeaderContentLength) Init() {
}

func (this *SipHeaderContentLength) AllowMulti() bool { return false }
func (this *SipHeaderContentLength) HasValue() bool   { return true }

/* RFC3261
 *
 * CSeq  =  "CSeq" HCOLON 1*DIGIT LWS Method
 *
 */
func (this *SipHeaderContentLength) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	name, newPos, err := ParseHeaderName(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if !name.EqualStringNoCase("Content-Length") && !name.EqualStringNoCase("l") {
		return newPos, &AbnfError{"Content-Length parse: wrong header-name", src, newPos}
	}

	return this.ParseValue(context, src, newPos)
}

func (this *SipHeaderContentLength) ParseValue(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos
	digit, _, newPos, ok := ParseUInt(src, newPos)
	if !ok {
		return newPos, &AbnfError{"Content-Length parse: wrong num", src, newPos}
	}

	this.size = uint32(digit)
	return newPos, nil
}

func (this *SipHeaderContentLength) Encode(buf *bytes.Buffer) {
	buf.WriteString("Content-Length: ")
	EncodeUInt(buf, uint64(this.size))
}

func (this *SipHeaderContentLength) String() string {
	return AbnfEncoderToString(this)
}

func ParseSipContentLength(context *ParseContext, src []byte, pos int) (newPos int, parsed SipHeaderParsed, err error) {
	header := SipHeaderContentLength{}
	header.Init()

	newPos, err = header.ParseValue(context, src, pos)
	return newPos, &header, err
}
