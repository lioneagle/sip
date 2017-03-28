package sipparser3

import (
	"bytes"
	//"fmt"
	//"strings"
)

type SipHeaderContentDisposition struct {
	dispType AbnfToken
	params   SipGenericParams
}

func NewSipHeaderContentDisposition() *SipHeaderContentDisposition {
	ret := &SipHeaderContentDisposition{}
	ret.Init()
	return ret
}

func (this *SipHeaderContentDisposition) Init() {
	this.dispType.SetNonExist()
	this.params.Init()
}

func (this *SipHeaderContentDisposition) AllowMulti() bool { return false }
func (this *SipHeaderContentDisposition) HasValue() bool   { return true }

/* RFC3261
 *
 * Content-Disposition   =  "Content-Disposition" HCOLON
 *                          disp-type *( SEMI disp-param )
 * disp-type             =  "render" / "session" / "icon" / "alert"
 *                          / disp-extension-token
 * disp-param            =  handling-param / generic-param
 * handling-param        =  "handling" EQUAL
 *                          ( "optional" / "required"
 *                          / other-handling )
 * other-handling        =  token
 * disp-extension-token  =  token
 *
 */
func (this *SipHeaderContentDisposition) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	name, newPos, err := ParseHeaderName(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if !name.EqualStringNoCase("Content-Disposition") {
		return newPos, &AbnfError{"Content-Disposition parse: wrong header-name", src, newPos}
	}

	return this.ParseValue(context, src, newPos)
}

func (this *SipHeaderContentDisposition) ParseValue(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	newPos, err = this.dispType.Parse(context, src, pos, IsSipToken)
	if err != nil {
		return newPos, err
	}

	return this.params.Parse(context, src, newPos, ';')
}

func (this *SipHeaderContentDisposition) Encode(buf *bytes.Buffer) {
	buf.WriteString("Content-Disposition: ")
	this.dispType.Encode(buf)
	this.params.Encode(buf, ';')
}

func (this *SipHeaderContentDisposition) String() string {
	return AbnfEncoderToString(this)
}

func ParseSipContentDisposition(context *ParseContext, src []byte, pos int) (newPos int, parsed SipHeaderParsed, err error) {
	header := SipHeaderContentDisposition{}
	newPos, err = header.ParseValue(context, src, pos)
	return newPos, &header, err
}
