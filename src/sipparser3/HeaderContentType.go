package sipparser3

import (
	"bytes"
	//"fmt"
	//"strings"
)

type SipHeaderContentType struct {
	mainType AbnfToken
	subType  AbnfToken
	params   SipGenericParams
}

func NewSipHeaderContentType() *SipHeaderContentType {
	ret := &SipHeaderContentType{}
	ret.Init()
	return ret
}

func (this *SipHeaderContentType) Init() {
	this.params.Init()
}

func (this *SipHeaderContentType) AllowMulti() bool { return false }
func (this *SipHeaderContentType) HasValue() bool   { return true }

/* RFC3261
 *
 * Content-Type     =  ( "Content-Type" / "c" ) HCOLON media-type
 * media-type       =  m-type SLASH m-subtype *(SEMI m-parameter)
 * m-type           =  discrete-type / composite-type
 * discrete-type    =  "text" / "image" / "audio" / "video"
 *                     / "application" / extension-token
 * composite-type   =  "message" / "multipart" / extension-token
 * extension-token  =  ietf-token / x-token
 * ietf-token       =  token
 * x-token          =  "x-" token
 * m-subtype        =  extension-token / iana-token
 * iana-token       =  token
 * m-parameter      =  m-attribute EQUAL m-value
 * m-attribute      =  token
 * m-value          =  token / quoted-string
 */
func (this *SipHeaderContentType) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	name, newPos, err := ParseHeaderName(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if !name.EqualStringNoCase("Content-Type") && !name.EqualStringNoCase("c") {
		return newPos, &AbnfError{"Content-Type parse: wrong header-name", src, newPos}
	}

	return this.ParseValue(context, src, newPos)
}

func (this *SipHeaderContentType) ParseValue(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos, err = this.mainType.Parse(context, src, pos, IsSipToken)
	if err != nil {
		return newPos, err
	}

	newPos, err = ParseSWSMark(src, newPos, '/')
	if err != nil {
		return newPos, err
	}

	newPos, err = this.subType.Parse(context, src, newPos, IsSipToken)
	if err != nil {
		return newPos, err
	}

	return this.params.Parse(context, src, newPos, ';')
}

func (this *SipHeaderContentType) Encode(buf *bytes.Buffer) {
	buf.WriteString("Content-Type: ")
	this.mainType.Encode(buf)
	buf.WriteByte('/')
	this.subType.Encode(buf)
	this.params.Encode(buf, ';')
}

func (this *SipHeaderContentType) String() string {
	return AbnfEncoderToString(this)
}

func ParseSipContentType(context *ParseContext, src []byte, pos int) (newPos int, parsed SipHeaderParsed, err error) {
	header := SipHeaderContentType{}
	header.Init()

	newPos, err = header.ParseValue(context, src, pos)
	return newPos, &header, err
}
