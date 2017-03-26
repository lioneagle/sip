package sipparser3

import (
	"bytes"
	//"fmt"
	//"strings"
)

type SipHeaderFrom struct {
	addr   SipAddr
	params SipGenericParams
}

func NewSipHeaderFrom() *SipHeaderFrom {
	ret := &SipHeaderFrom{}
	ret.Init()
	return ret
}

func (this *SipHeaderFrom) Init() {
	this.params.Init()
}

func (this *SipHeaderFrom) AllowMulti() bool { return false }
func (this *SipHeaderFrom) HasValue() bool   { return true }

/* RFC3261
 *
 * From        =  ( "From" / "f" ) HCOLON from-spec
 * from-spec   =  ( name-addr / addr-spec )
 *                *( SEMI from-param )
 */
func (this *SipHeaderFrom) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	name, newPos, err := ParseHeaderName(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if !name.EqualStringNoCase("from") && !name.EqualStringNoCase("f") {
		return newPos, &AbnfError{"From parse: wrong header-name", src, newPos}
	}

	return this.ParseValue(context, src, newPos)
}

func (this *SipHeaderFrom) ParseValue(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos
	newPos, err = this.addr.Parse(context, src, newPos)
	if err != nil {
		return newPos, err
	}

	return this.params.Parse(context, src, newPos, ';')
}

func (this *SipHeaderFrom) Encode(buf *bytes.Buffer) {
	buf.WriteString("From: ")
	this.addr.Encode(buf)
	this.params.Encode(buf, ';')
}

func (this *SipHeaderFrom) String() string {
	return AbnfEncoderToString(this)
}

func ParseSipFrom(context *ParseContext, src []byte, pos int) (newPos int, parsed SipHeaderParsed, err error) {
	header := SipHeaderFrom{}
	header.Init()

	newPos, err = header.ParseValue(context, src, pos)
	return newPos, &header, err
}
