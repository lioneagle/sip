package sipparser

import (
	"bytes"
	//"fmt"
	"unsafe"
)

type SipAddr struct {
	addrType int32
	addr     SipNameAddr
}

func NewSipAddr(context *ParseContext) (*SipAddr, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipAddr{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*SipAddr)(unsafe.Pointer(mem)).Init()
	return (*SipAddr)(unsafe.Pointer(mem)), addr
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
		return this.addr.Parse(context, src, newPos)
	}

	_, _, err = ParseUriScheme(context, src, newPos)
	if err == nil {
		this.addrType = ABNF_SIP_ADDR_SPEC
		return this.addr.addrsepc.ParseWithoutParam(context, src, newPos)
	}

	this.addrType = ABNF_SIP_NAME_ADDR
	return this.addr.Parse(context, src, newPos)
}

func (this *SipAddr) Encode(context *ParseContext, buf *bytes.Buffer) {
	this.addr.Encode(context, buf)
}

func (this *SipAddr) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}
