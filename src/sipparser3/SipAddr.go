package sipparser3

import (
	"bytes"
	//"fmt"
	//"strings"
)

type SipAddr struct {
	addr SipNameAddr
}

func NewSipAddr() *SipAddr {
	ret := &SipAddr{}
	ret.Init()
	return ret
}

func (this *SipAddr) Init() {
	this.addr.Init()
}

/* RFC3261
 *
 * sip-addr       =  ( name-addr / addr-spec )
 * name-addr      =  [ display-name ] LAQUOT addr-spec RAQUOT
 * addr-spec      =  SIP-URI / SIPS-URI / absoluteURI
 * display-name   =  *(token LWS)/ quoted-string
 *
 */
func (this *SipAddr) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos
	if newPos >= len(src) {
		return newPos, &AbnfError{"SipAddr parse: empty", src, newPos}
	}
	this.Init()

	newPos, err = ParseSWS(src, newPos)
	if err != nil {
		return newPos, err
	}

	if (src[newPos] == '<') || (src[newPos] == '"') {
		return this.addr.Parse(context, src, newPos)
	}

	_, _, err = ParseUriScheme(context, src, newPos)
	if err == nil {
		return this.addr.addrsepc.ParseWithoutParam(context, src, newPos)
	}

	return this.addr.Parse(context, src, newPos)
}

func (this *SipAddr) Encode(buf *bytes.Buffer) {
	this.addr.Encode(buf)
}

func (this *SipAddr) String() string {
	return AbnfEncoderToString(this)
}
