package sipparser

import (
	"bytes"
	//"fmt"
	"unsafe"
)

type SipMultiHeader struct {
	info    *SipHeaderInfo
	name    AbnfToken
	headers SipSingleHeaders
}

func NewSipMultiHeader(context *ParseContext) (*SipMultiHeader, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipMultiHeader{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*SipMultiHeader)(unsafe.Pointer(mem)).Init()
	return (*SipMultiHeader)(unsafe.Pointer(mem)), addr
}

func (this *SipMultiHeader) Init() {
	this.info = nil
	this.name.Init()
	this.headers.Init()
}

func (this *SipMultiHeader) HasInfo() bool      { return this.info != nil }
func (this *SipMultiHeader) HasShortName() bool { return this.HasInfo() && this.info.HasShortName() }
func (this *SipMultiHeader) ShortName() []byte  { return this.info.ShortName() }

func (this *SipMultiHeader) Size() int32 { return this.headers.Size() }

func (this *SipMultiHeader) EqualNameBytes(context *ParseContext, name []byte) bool {
	if this.info != nil {
		if EqualNoCase(this.info.name, name) {
			return true
		}
		return EqualNoCase(this.info.shortName, name)
	}
	return this.name.EqualBytesNoCase(context, name)
}

func (this *SipMultiHeader) EqualNameString(context *ParseContext, name string) bool {
	return this.EqualNameBytes(context, StringToByteSlice(name))
}

func (this *SipMultiHeader) AddHeader(context *ParseContext, header AbnfPtr) {
	this.headers.AddHeader(context, header)
}

func (this *SipMultiHeader) Encode(context *ParseContext, buf *bytes.Buffer) {
	if this.info != nil {
		buf.Write(this.info.name)
	} else {
		this.name.Encode(context, buf)
	}
	buf.WriteString(": ")
	this.headers.EncodeSameValues(context, buf)
}

func (this *SipMultiHeader) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

type SipMultiHeaders struct {
	AbnfList
}

func NewSipMultiHeaders(context *ParseContext) (*SipMultiHeaders, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipMultiHeaders{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*SipMultiHeaders)(unsafe.Pointer(mem)).Init()
	return (*SipMultiHeaders)(unsafe.Pointer(mem)), addr
}

func (this *SipMultiHeaders) Init() {
	this.AbnfList.Init()
}

func (this *SipMultiHeaders) Size() int32 { return this.Len() }
func (this *SipMultiHeaders) Empty() bool { return this.Len() == 0 }
func (this *SipMultiHeaders) GetHeaderByBytes(context *ParseContext, name []byte) (val *SipMultiHeader, ok bool) {
	for e := this.Front(context); e != nil; e = e.Next(context) {
		v := e.Value.GetSipMultiHeader(context)
		if v.EqualNameBytes(context, name) {
			return v, true
		}
	}
	return nil, false
}

func (this *SipMultiHeaders) GetHeaderByString(context *ParseContext, name string) (val *SipMultiHeader, ok bool) {
	return this.GetHeaderByBytes(context, StringToByteSlice(name))
}

func (this *SipMultiHeaders) AddHeader(context *ParseContext, header AbnfPtr) {
	this.PushBack(context, header)
}

func (this *SipMultiHeaders) Encode(context *ParseContext, buf *bytes.Buffer) {
	for e := this.Front(context); e != nil; e = e.Next(context) {
		v := e.Value.GetSipMultiHeader(context)
		v.Encode(context, buf)

	}
}

func (this *SipMultiHeaders) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}
