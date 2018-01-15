package sipparser

import (
	//"bytes"
	//"fmt"
	"unsafe"
)

type SipDisplayName struct {
	isQuotedString bool
	value          AbnfPtr
}

func NewSipDisplayName(context *ParseContext) AbnfPtr {
	addr := context.allocator.Alloc(uint32(unsafe.Sizeof(SipDisplayName{})))
	if addr == ABNF_PTR_NIL {
		return ABNF_PTR_NIL
	}
	addr.GetSipDisplayName(context).Init()
	return addr
}

func (this *SipDisplayName) Init() {
	this.isQuotedString = false
	this.value = ABNF_PTR_NIL
}

func (this *SipDisplayName) Exist() bool {
	return this.value != ABNF_PTR_NIL
}

func (this *SipDisplayName) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	/* RFC3261 Section 25.1, page 222
	 *
	 * display-name   =  *(token LWS)/ quoted-string
	 */
	newPos = pos
	if newPos >= len(src) {
		return newPos, nil
	}

	newPos, err = ParseSWS(src, newPos)
	if err != nil {
		return newPos, err
	}

	if newPos < len(src) && src[newPos] == '"' {
		return this.parseQuotedString(context, src, newPos)
	}

	return this.parseTokens(context, src, newPos)
}

func (this *SipDisplayName) parseQuotedString(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos
	this.isQuotedString = true
	addr := NewSipQuotedString(context)
	if addr == ABNF_PTR_NIL {
		return newPos, &AbnfError{"DisplayName parse: out of memory after first\"", src, newPos}
	}
	newPos, err = addr.GetSipQuotedString(context).Parse(context, src, newPos)
	if err != nil {
		return newPos, err
	}
	this.value = addr
	return newPos, nil
}

func (this *SipDisplayName) parseTokens(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos
	this.isQuotedString = false
	newPos = pos
	if !IsSipToken(src[newPos]) {
		return newPos, &AbnfError{"DisplayName parse: no token or quoted-string", src, newPos}
	}

	nameBegin := newPos
	for newPos < len(src) {
		if !IsSipToken(src[newPos]) {
			break
		}

		ref := AbnfRef{}
		newPos = ref.ParseSipToken(src, newPos)
		var ok bool
		newPos, ok = ParseLWS(src, newPos)
		if !ok {
			return newPos, &AbnfError{"DisplayName parse: wrong LWS", src, newPos}
		}
	}
	addr := NewAbnfBuf(context)
	if addr == ABNF_PTR_NIL {
		return newPos, &AbnfError{"DisplayName parse: out of memory after tokens", src, newPos}
	}
	addr.GetAbnfBuf(context).SetValue(context, src[nameBegin:newPos])
	this.value = addr
	return newPos, nil
}

func (this *SipDisplayName) Encode(context *ParseContext, buf *AbnfByteBuffer) {
	if this.value != ABNF_PTR_NIL {
		if this.isQuotedString {
			this.value.GetSipQuotedString(context).Encode(context, buf)
		} else {
			this.value.GetAbnfBuf(context).Encode(context, buf)
		}
	}
}

func (this *SipDisplayName) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

type SipNameAddr struct {
	displayname SipDisplayName
	addrspec    SipAddrSpec
}

func NewSipNameAddr(context *ParseContext) AbnfPtr {
	addr := context.allocator.Alloc(uint32(unsafe.Sizeof(SipNameAddr{})))
	if addr == ABNF_PTR_NIL {
		return ABNF_PTR_NIL
	}
	addr.GetSipAddr(context).Init()
	return addr
}

func (this *SipNameAddr) Init() {
	this.displayname.Init()
	this.addrspec.Init()
}

func (this *SipNameAddr) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	/* RFC3261 Section 25.1, page 222
	 *
	 * name-addr =  [ display-name ] LAQUOT addr-spec RAQUOT
	 * RAQUOT  =  ">" SWS ; right angle quote
	 * LAQUOT  =  SWS "<"; left angle quote
	 */
	this.Init()
	return this.ParseWithoutInit(context, src, pos)
}

func (this *SipNameAddr) ParseWithoutInit(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	/* RFC3261 Section 25.1, page 222
	 *
	 * name-addr =  [ display-name ] LAQUOT addr-spec RAQUOT
	 * RAQUOT  =  ">" SWS ; right angle quote
	 * LAQUOT  =  SWS "<"; left angle quote
	 */
	newPos, err = ParseSWS(src, pos)
	if err != nil {
		return newPos, err
	}

	if newPos >= len(src) {
		return newPos, &AbnfError{"SipNameAddr parse: no value", src, newPos}
	}

	if src[newPos] != '<' {
		newPos, err = this.displayname.Parse(context, src, newPos)
		if err != nil {
			return newPos, err
		}
	}

	newPos, err = ParseLeftAngleQuote(src, newPos)
	if err != nil {
		return newPos, err
	}

	newPos, err = this.addrspec.ParseWithoutInit(context, src, newPos)
	if err != nil {
		return newPos, err
	}

	return ParseRightAngleQuote(src, newPos)
}

func (this *SipNameAddr) Encode(context *ParseContext, buf *AbnfByteBuffer) {
	this.displayname.Encode(context, buf)
	buf.WriteByte('<')
	this.addrspec.Encode(context, buf)
	buf.WriteByte('>')
}

func (this *SipNameAddr) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}
