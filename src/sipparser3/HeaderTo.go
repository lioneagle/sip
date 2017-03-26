package sipparser3

import (
	"bytes"
	//"fmt"
	//"strings"
)

type SipHeaderTo struct {
	addr   SipAddr
	params SipGenericParams
}

func NewSipHeaderTo() *SipHeaderTo {
	ret := &SipHeaderTo{}
	ret.Init()
	return ret
}

func (this *SipHeaderTo) Init() {
	this.params.Init()
}

func (this *SipHeaderTo) AllowMulti() bool { return false }
func (this *SipHeaderTo) HasValue() bool   { return true }

/* RFC3261
 *
 * From        =  ( "From" / "f" ) HCOLON from-spec
 * from-spec   =  ( name-addr / addr-spec )
 *                *( SEMI from-param )
 */
func (this *SipHeaderTo) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	name, newPos, err := ParseHeaderName(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if !name.EqualStringNoCase("To") && !name.EqualStringNoCase("f") {
		return newPos, &AbnfError{"To parse: wrong header-name", src, newPos}
	}

	return this.ParseValue(context, src, newPos)
}

func (this *SipHeaderTo) ParseValue(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos
	newPos, err = this.addr.Parse(context, src, newPos)
	if err != nil {
		return newPos, err
	}

	return this.params.Parse(context, src, newPos, ';')
}

func (this *SipHeaderTo) Encode(buf *bytes.Buffer) {
	buf.WriteString("To: ")
	this.addr.Encode(buf)
	this.params.Encode(buf, ';')
}

func (this *SipHeaderTo) String() string {
	return AbnfEncoderToString(this)
}

func ParseSipTo(context *ParseContext, src []byte, pos int) (newPos int, parsed SipHeaderParsed, err error) {
	header := SipHeaderTo{}
	header.Init()

	newPos, err = header.Parse(context, src, pos)
	return newPos, &header, err
}
