package sipparser3

import (
	"bytes"
	"strconv"
	//"fmt"
	//"strings"
)

type SipHeaderCseq struct {
	id     uint32
	method AbnfToken
}

func NewSipHeaderCseq() *SipHeaderCseq {
	ret := &SipHeaderCseq{}
	ret.Init()
	return ret
}

func (this *SipHeaderCseq) Init() {
}

func (this *SipHeaderCseq) AllowMulti() bool { return false }
func (this *SipHeaderCseq) HasValue() bool   { return true }

/* RFC3261
 *
 * CSeq  =  "CSeq" HCOLON 1*DIGIT LWS Method
 *
 */
func (this *SipHeaderCseq) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	name, newPos, err := ParseHeaderName(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if !name.EqualStringNoCase("CSeq") {
		return newPos, &AbnfError{"CSeq parse: wrong header-name", src, newPos}
	}

	return this.ParseValue(context, src, newPos)
}

func (this *SipHeaderCseq) ParseValue(context *ParseContext, src []byte, pos int) (newPos int, err error) {
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

func (this *SipHeaderCseq) Encode(buf *bytes.Buffer) {
	buf.WriteString("CSeq: ")
	buf.WriteString(strconv.FormatUint(uint64(this.id), 10))
	buf.WriteByte(' ')
	this.method.Encode(buf)
}

func (this *SipHeaderCseq) String() string {
	return AbnfEncoderToString(this)
}

func ParseSipCseq(context *ParseContext, src []byte, pos int) (newPos int, parsed SipHeaderParsed, err error) {
	header := SipHeaderCseq{}
	header.Init()

	newPos, err = header.ParseValue(context, src, pos)
	return newPos, &header, err
}
