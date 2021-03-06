package sipparser

import (
	//"bytes"
	//"fmt"
	//"strings"
	"unsafe"
)

type SipVersion struct {
	major AbnfBuf
	minor AbnfBuf
}

func NewSipVersion(context *ParseContext) AbnfPtr {
	addr := context.allocator.Alloc(uint32(unsafe.Sizeof(SipVersion{})))
	if addr == ABNF_PTR_NIL {
		return ABNF_PTR_NIL
	}
	addr.GetSipVersion(context).Init()
	return addr
}

func (this *SipVersion) Init() {
	this.major.Init()
	this.minor.Init()
}

/*
 * SIP-Version    =  "SIP" "/" 1*DIGIT "." 1*DIGIT
 */
func (this *SipVersion) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos
	if (newPos + 4) >= len(src) {
		return newPos, &AbnfError{"SipVersion parse: len not enough", src, newPos}
	}

	if !EqualNoCase(src[newPos:newPos+4], []byte{'s', 'i', 'p', '/'}) {
		return newPos, &AbnfError{"SipVersion parse: wrong name", src, newPos}
	}

	newPos += 4

	newPos, err = this.major.ParseDigit(context, src, newPos)
	if err != nil {
		return newPos, &AbnfError{"SipVersion parse: parse major version failed", src, newPos}
	}

	if src[newPos] != '.' {
		return newPos, &AbnfError{"SipVersion parse: no '.' after major version", src, newPos}
	}

	newPos++

	newPos, err = this.minor.ParseDigit(context, src, newPos)
	if err != nil {
		return newPos, &AbnfError{"SipVersion parse: parse minor version failed", src, newPos}
	}

	return newPos, nil
}

func (this *SipVersion) Encode(context *ParseContext, buf *AbnfByteBuffer) {
	buf.WriteString("SIP/")
	this.major.Encode(context, buf)
	buf.WriteByte('.')
	this.minor.Encode(context, buf)
}

func (this *SipVersion) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}
