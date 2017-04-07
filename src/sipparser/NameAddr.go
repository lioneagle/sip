package sipparser

import (
	"bytes"
	//"fmt"
	"unsafe"
)

type SipDisplayName struct {
	isQuotedString bool
	value          AbnfPtr
}

func NewSipDisplayName(context *ParseContext) (*SipDisplayName, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipDisplayName{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*SipDisplayName)(unsafe.Pointer(mem)).Init()
	return (*SipDisplayName)(unsafe.Pointer(mem)), addr
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
		this.isQuotedString = true
		quotedString, addr := NewSipQuotedString(context)
		if quotedString == nil {
			return newPos, &AbnfError{"DisplayName parse: out of memory after first\"", src, newPos}
		}
		newPos, err = quotedString.Parse(context, src, newPos)
		if err != nil {
			return newPos, err
		}
		this.value = addr

	} else {
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
			newPos = ref.Parse(src, newPos, IsSipToken)
			newPos, err = ParseLWS(src, newPos)
			if err != nil {
				return newPos, err
			}
		}
		name, addr := NewAbnfBuf(context)
		if name == nil {
			return newPos, &AbnfError{"DisplayName parse: out of memory after tokens", src, newPos}
		}
		name.SetValue(context, src[nameBegin:newPos])
		this.value = addr
	}

	return newPos, nil
}

func (this *SipDisplayName) Encode(context *ParseContext, buf *bytes.Buffer) {
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
	addrsepc    SipAddrSpec
}

func NewSipNameAddr(context *ParseContext) (*SipNameAddr, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipNameAddr{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*SipNameAddr)(unsafe.Pointer(mem)).Init()
	return (*SipNameAddr)(unsafe.Pointer(mem)), addr
}

func (this *SipNameAddr) Init() {
	this.displayname.Init()
	this.addrsepc.Init()
}

func (this *SipNameAddr) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	/* RFC3261 Section 25.1, page 222
	 *
	 * name-addr =  [ display-name ] LAQUOT addr-spec RAQUOT
	 * RAQUOT  =  ">" SWS ; right angle quote
	 * LAQUOT  =  SWS "<"; left angle quote
	 */
	newPos = pos
	this.Init()

	newPos, err = ParseSWS(src, newPos)
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

	newPos, err = this.addrsepc.Parse(context, src, newPos)
	if err != nil {
		return newPos, err
	}

	return ParseRightAngleQuote(src, newPos)
}

func (this *SipNameAddr) Encode(context *ParseContext, buf *bytes.Buffer) {
	this.displayname.Encode(context, buf)
	buf.WriteByte('<')
	this.addrsepc.Encode(context, buf)
	buf.WriteByte('>')
}

func (this *SipNameAddr) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}
