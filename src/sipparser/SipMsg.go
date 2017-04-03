package sipparser

import (
	"bytes"
	//"fmt"
	"unsafe"
)

type SipMsg struct {
	startLine SipStartLine
	headers   SipHeaders
}

func NewSipMsg(context *ParseContext) (*SipMsg, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipMsg{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*SipMsg)(unsafe.Pointer(mem)).Init()
	return (*SipMsg)(unsafe.Pointer(mem)), addr
}

func (this *SipMsg) Init() {
	this.startLine.Init()
	this.headers.Init()
}

func (this *SipMsg) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	newPos = pos
	newPos, err = this.startLine.Parse(context, src, newPos)
	if err != nil {
		return newPos, err
	}

	if newPos >= len(src) {
		return newPos, &AbnfError{"SipMsg parse: no headers", src, newPos}
	}

	return this.headers.Parse(context, src, newPos)
}

func (this *SipMsg) Encode(context *ParseContext, buf *bytes.Buffer) {
	this.startLine.Encode(context, buf)
	this.headers.Encode(context, buf)
	buf.WriteByte('\r')
	buf.WriteByte('\n')

}

func (this *SipMsg) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}
