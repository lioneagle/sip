package sipparser

import (
	"bytes"
	//"fmt"
	"unsafe"
)

type SipSingleHeader struct {
	info   *SipHeaderInfo
	name   AbnfToken
	value  AbnfToken
	parsed AbnfPtr
}

func NewSipSingleHeader(context *ParseContext) (*SipSingleHeader, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipSingleHeader{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*SipSingleHeader)(unsafe.Pointer(mem)).Init()
	return (*SipSingleHeader)(unsafe.Pointer(mem)), addr
}

func (this *SipSingleHeader) Init() {
	this.info = nil
	this.name.Init()
	this.value.Init()
	this.parsed = ABNF_PTR_NIL
}

func (this *SipSingleHeader) HasInfo() bool      { return this.info != nil }
func (this *SipSingleHeader) HasShortName() bool { return this.HasInfo() && this.info.HasShortName() }
func (this *SipSingleHeader) ShortName() []byte  { return this.info.ShortName() }

func (this *SipSingleHeader) EqualNameBytes(context *ParseContext, name []byte) bool {
	if this.info != nil {
		if EqualNoCase(this.info.name, name) {
			return true
		}
		return EqualNoCase(this.info.shortName, name)
	}
	return this.name.EqualBytesNoCase(context, name)
}

func (this *SipSingleHeader) EqualNameString(context *ParseContext, name string) bool {
	return this.EqualNameBytes(context, StringToByteSlice(name))
}

func (this *SipSingleHeader) Encode(context *ParseContext, buf *bytes.Buffer) {
	if this.info != nil {
		buf.Write(this.info.name)
	} else {
		this.name.Encode(context, buf)
	}
	buf.WriteString(": ")
	this.EncodeValue(context, buf)
}

func (this *SipSingleHeader) EncodeValue(context *ParseContext, buf *bytes.Buffer) {
	if this.value.Exist() {
		this.value.Encode(context, buf)
	}
}

func (this *SipSingleHeader) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

type SipSingleHeaders struct {
	AbnfList
}

func NewSipSingleHeaders(context *ParseContext) (*SipSingleHeaders, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipSingleHeaders{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*SipSingleHeaders)(unsafe.Pointer(mem)).Init()
	return (*SipSingleHeaders)(unsafe.Pointer(mem)), addr
}

func (this *SipSingleHeaders) Init() {
	this.AbnfList.Init()
}

func (this *SipSingleHeaders) Size() int32 { return this.Len() }
func (this *SipSingleHeaders) Empty() bool { return this.Len() == 0 }
func (this *SipSingleHeaders) GetHeaderByBytes(context *ParseContext, name []byte) (val *SipSingleHeader, ok bool) {
	for e := this.Front(context); e != nil; e = e.Next(context) {
		v := e.Value.GetSipSingleHeader(context)
		if v.EqualNameBytes(context, name) {
			return v, true
		}
	}

	return nil, false
}

func (this *SipSingleHeaders) GetHeaderByString(context *ParseContext, name string) (val *SipSingleHeader, ok bool) {
	return this.GetHeaderByBytes(context, StringToByteSlice(name))
}

func (this *SipSingleHeaders) AddHeader(context *ParseContext, header AbnfPtr) {
	this.PushBack(context, header)
}

func (this *SipSingleHeaders) Encode(context *ParseContext, buf *bytes.Buffer) {
	for e := this.Front(context); e != nil; e = e.Next(context) {
		v := e.Value.GetSipSingleHeader(context)
		v.Encode(context, buf)
		buf.WriteString("\r\n")

	}
}

func (this *SipSingleHeaders) EncodeSameValues(context *ParseContext, buf *bytes.Buffer) {
	e := this.Front(context)
	if e != nil {
		e.Value.GetSipSingleHeader(context).EncodeValue(context, buf)
		e = e.Next(context)
	}
	for ; e != nil; e = e.Next(context) {
		buf.WriteString(", ")
		e.Value.GetSipSingleHeader(context).EncodeValue(context, buf)
	}
	buf.WriteString("\r\n")
}

func (this *SipSingleHeaders) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}
