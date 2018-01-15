package sipparser

import (
	//"bytes"
	//"fmt"
	"unsafe"
)

type SipAddr struct {
	addrType int32
	addr     SipNameAddr
}

func NewSipAddr(context *ParseContext) AbnfPtr {
	addr := context.allocator.Alloc(uint32(unsafe.Sizeof(SipAddr{})))
	if addr == ABNF_PTR_NIL {
		return ABNF_PTR_NIL
	}

	addr.GetSipAddr(context).Init()
	return addr
}

func (this *SipAddr) Init() {
	this.addrType = ABNF_SIP_ADDR_SPEC
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
	this.Init()
	return this.ParseWithoutInit(context, src, pos)
}

func (this *SipAddr) ParseWithoutInit(context *ParseContext, src []byte, pos int) (newPos int, err error) {
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
		this.addrType = ABNF_SIP_NAME_ADDR
		return this.addr.ParseWithoutInit(context, src, newPos)
	}

	var scheme AbnfBuf

	_, err = ParseUriScheme(context, src, newPos, &scheme)
	if err == nil {
		this.addrType = ABNF_SIP_ADDR_SPEC
		return this.addr.addrspec.ParseWithoutParamNorInit(context, src, newPos)
	}

	this.addrType = ABNF_SIP_NAME_ADDR
	return this.addr.ParseWithoutInit(context, src, newPos)
}

func (this *SipAddr) Encode(context *ParseContext, buf *AbnfByteBuffer) {
	this.addr.Encode(context, buf)
}

func (this *SipAddr) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}
