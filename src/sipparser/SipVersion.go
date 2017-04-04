package sipparser

import (
	"bytes"
	//"fmt"
	//"strings"
	"unsafe"
)

type SipVersion struct {
	major AbnfBuf
	minor AbnfBuf
}

func NewSipVersion(context *ParseContext) (*SipVersion, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipVersion{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*SipVersion)(unsafe.Pointer(mem)).Init()
	return (*SipVersion)(unsafe.Pointer(mem)), addr
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

	newPos, err = this.major.Parse(context, src, newPos, IsDigit)
	if err != nil {
		return newPos, &AbnfError{"SipVersion parse: parse major version failed", src, newPos}
	}

	if src[newPos] != '.' {
		return newPos, &AbnfError{"SipVersion parse: no '.' after major version", src, newPos}
	}

	newPos++

	newPos, err = this.minor.Parse(context, src, newPos, IsDigit)
	if err != nil {
		return newPos, &AbnfError{"SipVersion parse: parse minor version failed", src, newPos}
	}

	return newPos, nil
}

func (this *SipVersion) Encode(context *ParseContext, buf *bytes.Buffer) {
	buf.WriteString("SIP/")
	this.major.Encode(context, buf)
	buf.WriteByte('.')
	this.minor.Encode(context, buf)
}

func (this *SipVersion) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}
