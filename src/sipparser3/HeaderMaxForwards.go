package sipparser3

import (
	"bytes"
	//"fmt"
	//"strings"
)

type SipHeaderMaxForwards struct {
	size uint32
}

func NewSipHeaderMaxForwards() *SipHeaderMaxForwards {
	ret := &SipHeaderMaxForwards{}
	ret.Init()
	return ret
}

func (this *SipHeaderMaxForwards) Init() {
}

func (this *SipHeaderMaxForwards) AllowMulti() bool { return false }
func (this *SipHeaderMaxForwards) HasValue() bool   { return true }

/* RFC3261
 *
 * Max-Forwards  =  "Max-Forwards" HCOLON 1*DIGIT
 */
func (this *SipHeaderMaxForwards) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	name, newPos, err := ParseHeaderName(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if !name.EqualStringNoCase("Max-Forwards") {
		return newPos, &AbnfError{"Max-Forwards parse: wrong header-name", src, newPos}
	}

	return this.ParseValue(context, src, newPos)
}

func (this *SipHeaderMaxForwards) ParseValue(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos
	digit, _, newPos, ok := ParseUInt(src, newPos)
	if !ok {
		return newPos, &AbnfError{"Max-Forwards parse: wrong num", src, newPos}
	}

	this.size = uint32(digit)
	return newPos, nil
}

func (this *SipHeaderMaxForwards) Encode(buf *bytes.Buffer) {
	buf.WriteString("Max-Forwards: ")
	EncodeUInt(buf, uint64(this.size))
}

func (this *SipHeaderMaxForwards) String() string {
	return AbnfEncoderToString(this)
}

func ParseSipMaxForwards(context *ParseContext, src []byte, pos int) (newPos int, parsed SipHeaderParsed, err error) {
	header := SipHeaderMaxForwards{}
	header.Init()

	newPos, err = header.ParseValue(context, src, pos)
	return newPos, &header, err
}
