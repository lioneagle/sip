package sipparser

import (
	"bytes"
	//"fmt"
	//"strings"
	"unsafe"
)

type SipUriHeader struct {
	name  AbnfToken
	value AbnfToken
}

func NewSipUriHeader(context *ParseContext) (*SipUriHeader, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipUriHeader{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*SipUriHeader)(unsafe.Pointer(mem)).Init()
	return (*SipUriHeader)(unsafe.Pointer(mem)), addr
}

func (this *SipUriHeader) Init() {
	this.name.Init()
	this.value.Init()
}

func (this *SipUriHeader) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos, err = this.name.ParseEscapable(context, src, pos, IsSipHname)
	if err != nil {
		return newPos, err
	}

	if newPos >= len(src) {
		return newPos, &AbnfError{"sip-uri parse: parse header failed: no = after hname", src, newPos}
	}

	if src[newPos] != '=' {
		return newPos, &AbnfError{"sip-uri parse: parse header failed: no = after hname", src, newPos}
	}

	newPos++

	this.value.SetExist()

	if newPos >= len(src) {
		return newPos, nil
	}

	if !IsSipHvalue(src[newPos]) {
		return newPos, nil
	}

	return this.value.ParseEscapable(context, src, newPos, IsSipHvalue)
}

func (this *SipUriHeader) Encode(context *ParseContext, buf *bytes.Buffer) {
	buf.Write(Escape(this.name.GetAsByteSlice(context), IsSipHname))
	buf.WriteByte('=')
	if this.value.Exist() {
		buf.Write(Escape(this.value.GetAsByteSlice(context), IsSipHvalue))
	}
}

func (this *SipUriHeader) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

type SipUriHeaders struct {
	AbnfList
}

func NewSipUriHeaders(context *ParseContext) (*SipUriHeaders, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipUriHeaders{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*SipUriHeaders)(unsafe.Pointer(mem)).Init()
	return (*SipUriHeaders)(unsafe.Pointer(mem)), addr
}

func (this *SipUriHeaders) Init() {
	this.AbnfList.Init()
}

func (this *SipUriHeaders) Size() int32 { return this.Len() }
func (this *SipUriHeaders) Empty() bool { return this.Len() == 0 }
func (this *SipUriHeaders) GetHeader(context *ParseContext, name string) (val *SipUriHeader, ok bool) {
	for e := this.Front(context); e != nil; e = e.Next(context) {
		v := e.Value.GetSipUriHeader(context)
		if v.name.EqualStringNoCase(context, name) {
			return v, true
		}
	}
	return nil, false
}

func (this *SipUriHeaders) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	newPos = pos
	if newPos >= len(src) {
		return newPos, &AbnfError{"sip-uri parse: parse uri-header failed: reach end after ';'", src, newPos}
	}

	for newPos < len(src) {
		header, addr := NewSipUriHeader(context)
		if header == nil {
			return newPos, &AbnfError{"sip-uri  parse: out of memory for sip uri headers", src, newPos}
		}
		newPos, err = header.Parse(context, src, newPos)
		if err != nil {
			return newPos, err
		}
		this.PushBack(context, addr)

		if newPos >= len(src) {
			return newPos, nil
		}

		if src[newPos] != '&' {
			return newPos, nil
		}
		newPos++
	}

	return newPos, err
}

func (this *SipUriHeaders) EqualRFC3261(context *ParseContext, rhs *SipUriHeaders) bool {
	if this.Size() != rhs.Size() {
		return false
	}

	for e := this.Front(context); e != nil; e = e.Next(context) {
		v := e.Value.GetSipUriHeader(context)
		header, ok := rhs.GetHeader(context, v.name.String(context))
		if ok {
			if !header.value.EqualNoCase(context, &v.value) {
				return false
			}
		}
	}

	return true
}

func (this *SipUriHeaders) Encode(context *ParseContext, buf *bytes.Buffer) {
	e := this.Front(context)
	if e != nil {
		e.Value.GetSipUriHeader(context).Encode(context, buf)
	}
	e = e.Next(context)

	for ; e != nil; e = e.Next(context) {
		buf.WriteByte('&')
		e.Value.GetSipUriHeader(context).Encode(context, buf)
	}
}

func (this *SipUriHeaders) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}
