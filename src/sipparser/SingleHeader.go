package sipparser

import (
	"bytes"
	//"fmt"
	"unsafe"
)

type SipSingleHeader struct {
	info   *SipHeaderInfo
	name   AbnfBuf
	value  AbnfBuf
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

func GenerateSingleHeader(context *ParseContext, name, value string) (*SipSingleHeader, AbnfPtr) {
	header, addr := NewSipSingleHeader(context)
	if header == nil {
		return nil, ABNF_PTR_NIL
	}

	header.SetNameByteSlice(context, StringToByteSlice(name))
	header.SetValueByteSlice(context, StringToByteSlice(value))
	return header, addr
}

func (this *SipSingleHeader) Init() {
	this.info = nil
	this.name.Init()
	this.value.Init()
	this.parsed = ABNF_PTR_NIL
}

func (this *SipSingleHeader) HasInfo() bool            { return this.info != nil }
func (this *SipSingleHeader) HasShortName() bool       { return this.HasInfo() && this.info.HasShortName() }
func (this *SipSingleHeader) ShortName() []byte        { return this.info.ShortName() }
func (this *SipSingleHeader) IsParsed() bool           { return this.parsed != ABNF_PTR_NIL }
func (this *SipSingleHeader) GetParsed() AbnfPtr       { return this.parsed }
func (this *SipSingleHeader) SetParsed(parsed AbnfPtr) { this.parsed = parsed }

func (this *SipSingleHeader) SetInfo(name string) {
	this.info, _ = GetSipHeaderInfo(name)
}

func (this *SipSingleHeader) NameHasPrefixBytes(context *ParseContext, prefix []byte) bool {
	if this.info != nil {
		return HasPrefixByteSliceNoCase(this.info.name, prefix)
	}
	return this.name.HasPrefixByteSliceNoCase(context, prefix)
}

func (this *SipSingleHeader) SetNameByteSlice(context *ParseContext, name []byte) {
	this.name.SetByteSlice(context, name)
}

func (this *SipSingleHeader) SetValueByteSlice(context *ParseContext, value []byte) {
	this.value.SetByteSlice(context, value)
}

func (this *SipSingleHeader) EqualNameByteSlice(context *ParseContext, name []byte) bool {
	if this.info != nil {
		if EqualNoCase(this.info.name, name) {
			return true
		}
		return EqualNoCase(this.info.shortName, name)
	}
	return this.name.EqualByteSliceNoCase(context, name)
}

func (this *SipSingleHeader) EqualNameString(context *ParseContext, name string) bool {
	return this.EqualNameByteSlice(context, StringToByteSlice(name))
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
	if this.IsParsed() && this.info != nil && this.info.encodeFunc != nil {
		this.info.encodeFunc(this.parsed, context, buf)
	} else if this.value.Exist() {
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
func (this *SipSingleHeaders) GetHeaderByByteSlice(context *ParseContext, name []byte) (val *SipSingleHeader, ok bool) {
	for e := this.Front(context); e != nil; e = e.Next(context) {
		v := e.Value.GetSipSingleHeader(context)
		if v.EqualNameByteSlice(context, name) {
			return v, true
		}
	}

	return nil, false
}

func (this *SipSingleHeaders) GetHeaderByString(context *ParseContext, name string) (val *SipSingleHeader, ok bool) {
	return this.GetHeaderByByteSlice(context, StringToByteSlice(name))
}

// remove Content-* headers from sip message except Content-Length and Content-Type*/
func (this *SipSingleHeaders) RemoveContentHeaders(context *ParseContext) {
	prefix := StringToByteSlice("Content-")

	for e := this.Front(context); e != nil; {
		v := e.Value.GetSipSingleHeader(context)
		if v.EqualNameString(context, ABNF_NAME_SIP_HDR_CONTENT_TYPE) ||
			v.EqualNameString(context, ABNF_NAME_SIP_HDR_CONTENT_LENGTH) {
			e = e.Next(context)
			continue
		}

		if !v.NameHasPrefixBytes(context, prefix) {
			e = e.Next(context)
			continue
		}

		e = this.Remove(context, e)
	}
}

func (this *SipSingleHeaders) GetHeaderParsedByString(context *ParseContext, name string) (parsed AbnfPtr, ok bool) {
	header, ok := this.GetHeaderByString(context, name)
	if !ok {
		return ABNF_PTR_NIL, false
	}

	if header.IsParsed() {
		return header.GetParsed(), true
	}

	if header.info == nil {
		header.SetInfo(name)
	}

	if header.info == nil || header.info.parseFunc == nil || !header.value.Exist() {
		return ABNF_PTR_NIL, false
	}

	_, parsed, err := header.info.parseFunc(context, header.value.GetAsByteSlice(context), 0)
	if err != nil || parsed == ABNF_PTR_NIL {
		return ABNF_PTR_NIL, false
	}
	header.parsed = parsed
	return parsed, true
}

func (this *SipSingleHeaders) AddHeader(context *ParseContext, header AbnfPtr) {
	this.PushBack(context, header)
}

func (this *SipSingleHeaders) GenerateAndAddHeader(context *ParseContext, name, value string) (*SipSingleHeader, AbnfPtr) {
	header, addr := GenerateSingleHeader(context, name, value)
	if header == nil {
		return nil, ABNF_PTR_NIL
	}

	this.AddHeader(context, addr)
	return header, addr
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
