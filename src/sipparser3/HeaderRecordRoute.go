package sipparser3

import (
	"bytes"
	//"fmt"
	//"strings"
)

type SipHeaderRecordRoute struct {
	addr   SipNameAddr
	params SipGenericParams
}

func NewSipHeaderRecordRoute() *SipHeaderRecordRoute {
	ret := &SipHeaderRecordRoute{}
	ret.Init()
	return ret
}

func (this *SipHeaderRecordRoute) Init() {
	this.params.Init()
}

func (this *SipHeaderRecordRoute) AllowMulti() bool { return false }
func (this *SipHeaderRecordRoute) HasValue() bool   { return true }

/* RFC3261
 *
 * Record-Route  =  "Record-Route" HCOLON rec-route *(COMMA rec-route)
 * rec-route     =  name-addr *( SEMI rr-param )
 * rr-param      =  generic-param
 */
func (this *SipHeaderRecordRoute) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	name, newPos, err := ParseHeaderName(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if !name.EqualStringNoCase("Record-Route") {
		return newPos, &AbnfError{"Record-Route parse: wrong header-name", src, newPos}
	}

	return this.ParseValue(context, src, newPos)
}

func (this *SipHeaderRecordRoute) ParseValue(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos
	newPos, err = this.addr.Parse(context, src, newPos)
	if err != nil {
		return newPos, err
	}

	return this.params.Parse(context, src, newPos, ';')
}

func (this *SipHeaderRecordRoute) Encode(buf *bytes.Buffer) {
	buf.WriteString("Record-Route: ")
	this.addr.Encode(buf)
	this.params.Encode(buf, ';')
}

func (this *SipHeaderRecordRoute) String() string {
	return AbnfEncoderToString(this)
}

func ParseSipRecordRoute(context *ParseContext, src []byte, pos int) (newPos int, parsed SipHeaderParsed, err error) {
	header := SipHeaderRecordRoute{}
	header.Init()

	newPos, err = header.ParseValue(context, src, pos)
	return newPos, &header, err
}
