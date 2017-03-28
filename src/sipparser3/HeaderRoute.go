package sipparser3

import (
	"bytes"
	//"fmt"
	//"strings"
)

type SipHeaderRoute struct {
	addr   SipNameAddr
	params SipGenericParams
}

func NewSipHeaderRoute() *SipHeaderRoute {
	ret := &SipHeaderRoute{}
	ret.Init()
	return ret
}

func (this *SipHeaderRoute) Init() {
	this.addr.Init()
	this.params.Init()
}

func (this *SipHeaderRoute) AllowMulti() bool { return false }
func (this *SipHeaderRoute) HasValue() bool   { return true }

/* RFC3261
 *
 * Route        =  "Route" HCOLON route-param *(COMMA route-param)
 * route-param  =  name-addr *( SEMI rr-param )
 * rr-param     =  generic-param
 */
func (this *SipHeaderRoute) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	name, newPos, err := ParseHeaderName(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if !name.EqualStringNoCase("Route") {
		return newPos, &AbnfError{"Route parse: wrong header-name", src, newPos}
	}

	return this.ParseValue(context, src, newPos)
}

func (this *SipHeaderRoute) ParseValue(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	newPos = pos
	newPos, err = this.addr.Parse(context, src, newPos)
	if err != nil {
		return newPos, err
	}

	return this.params.Parse(context, src, newPos, ';')
}

func (this *SipHeaderRoute) Encode(buf *bytes.Buffer) {
	buf.WriteString("Route: ")
	this.addr.Encode(buf)
	this.params.Encode(buf, ';')
}

func (this *SipHeaderRoute) String() string {
	return AbnfEncoderToString(this)
}

func ParseSipRoute(context *ParseContext, src []byte, pos int) (newPos int, parsed SipHeaderParsed, err error) {
	header := SipHeaderRoute{}
	newPos, err = header.ParseValue(context, src, pos)
	return newPos, &header, err
}
