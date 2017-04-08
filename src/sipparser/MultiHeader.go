package sipparser

import (
	"bytes"
	//"fmt"
	"unsafe"
)

type SipMultiHeader struct {
	info    *SipHeaderInfo
	name    AbnfBuf
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

func (this *SipMultiHeader) SetInfo(name string) {
	this.info, _ = GetSipHeaderInfo(name)
}

func (this *SipMultiHeader) Size() int32 { return this.headers.Size() }

func (this *SipMultiHeader) NameHasPrefixByteSlice(context *ParseContext, prefix []byte) bool {
	if this.info != nil {
		return HasPrefixByteSliceNoCase(this.info.name, prefix)
	}
	return this.name.HasPrefixByteSliceNoCase(context, prefix)
}

func (this *SipMultiHeader) SetNameByteSlice(context *ParseContext, name []byte) {
	this.name.SetByteSlice(context, name)
}

func (this *SipMultiHeader) EqualNameByteSlice(context *ParseContext, name []byte) bool {
	if this.info != nil {
		if EqualNoCase(this.info.name, name) {
			return true
		}
		return EqualNoCase(this.info.shortName, name)
	}
	return this.name.EqualByteSliceNoCase(context, name)
}

func (this *SipMultiHeader) EqualNameString(context *ParseContext, name string) bool {
	return this.EqualNameByteSlice(context, StringToByteSlice(name))
}

func (this *SipMultiHeader) AddHeader(context *ParseContext, header AbnfPtr) {
	this.headers.AddHeader(context, header)
}

func (this *SipMultiHeader) GenerateAndAddHeader(context *ParseContext, name, value string) (*SipSingleHeader, AbnfPtr) {
	header, addr := GenerateSingleHeader(context, name, value)
	if header == nil {
		return nil, ABNF_PTR_NIL
	}

	this.AddHeader(context, addr)
	return header, addr
}

func (this *SipMultiHeader) Parse(context *ParseContext, src []byte, pos int, info *SipHeaderInfo) (newPos int, err error) {
	if info.parseFunc != nil && info.needParse {
		return this.parseParsableHeader(context, src, pos, info)

	}
	return this.parseUnparsableHeader(context, src, pos, info)
}

func (this *SipMultiHeader) parseParsableHeader(context *ParseContext, src []byte, pos int, info *SipHeaderInfo) (newPos int, err error) {
	newPos = pos
	for newPos < len(src) {
		addr, newPos, err := parseOneParsableSingleHeader(context, src, newPos, info)
		if err != nil {
			return newPos, err
		}
		this.AddHeader(context, addr)

		// now should be COMMA or CRLF
		newPos1, err := ParseSWSMark(src, newPos, ',')
		if err != nil {
			// should be CRLF
			return ParseCRLF(src, newPos)
		}
		newPos = newPos1

	}
	return newPos, nil
}

func (this *SipMultiHeader) parseUnparsableHeader(context *ParseContext, src []byte, pos int, info *SipHeaderInfo) (newPos int, err error) {
	addr, newPos, err := parseOneUnparsableSingleHeader(context, AbnfRef{0, int32(len(info.name))}, src, pos, info)
	if err != nil {
		return newPos, err
	}
	this.AddHeader(context, addr)

	return newPos, nil
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
func (this *SipMultiHeaders) GetHeaderByByteSlice(context *ParseContext, name []byte) (val *SipMultiHeader, ok bool) {
	for e := this.Front(context); e != nil; e = e.Next(context) {
		v := e.Value.GetSipMultiHeader(context)
		if v.EqualNameByteSlice(context, name) {
			return v, true
		}
	}
	return nil, false
}

func (this *SipMultiHeaders) RemoveHeaderByNameString(context *ParseContext, name string) {
	name1 := StringToByteSlice(name)
	for e := this.Front(context); e != nil; {
		v := e.Value.GetSipSingleHeader(context)
		if v.EqualNameByteSlice(context, name1) {
			e = this.Remove(context, e)
			continue
		}
		e = e.Next(context)
	}
}

func (this *SipMultiHeaders) GetHeaderByString(context *ParseContext, name string) (val *SipMultiHeader, ok bool) {
	return this.GetHeaderByByteSlice(context, StringToByteSlice(name))
}

// remove Content-* headers from sip message except Content-Length and Content-Type*/
func (this *SipMultiHeaders) RemoveContentHeaders(context *ParseContext) {
	prefix := StringToByteSlice("Content-")

	for e := this.Front(context); e != nil; {
		v := e.Value.GetSipMultiHeader(context)
		if v.EqualNameString(context, ABNF_NAME_SIP_HDR_CONTENT_TYPE) ||
			v.EqualNameString(context, ABNF_NAME_SIP_HDR_CONTENT_LENGTH) {
			e = e.Next(context)
			continue
		}

		if !v.NameHasPrefixByteSlice(context, prefix) {
			e = e.Next(context)
			continue
		}

		e = this.Remove(context, e)
	}
}

func (this *SipMultiHeaders) CopyContentHeaders(context *ParseContext, rhs *SipMultiHeaders) {
	prefix := StringToByteSlice("Content-")

	for e := this.Front(context); e != nil; e = e.Next(context) {
		v := e.Value.GetSipSingleHeader(context)
		if v.NameHasPrefixBytes(context, prefix) {
			rhs.AddHeader(context, e.Value)
		}
	}
}

func (this *SipMultiHeaders) AddHeader(context *ParseContext, header AbnfPtr) {
	this.PushBack(context, header)
}

func (this *SipMultiHeaders) GenerateAndAddHeader(context *ParseContext, name, value string) (*SipSingleHeader, AbnfPtr) {
	multiHeader, ok := this.GetHeaderByByteSlice(context, StringToByteSlice(name))
	if multiHeader == nil || !ok {
		var addr AbnfPtr
		multiHeader, addr = NewSipMultiHeader(context)
		if multiHeader == nil {
			return nil, ABNF_PTR_NIL
		}
		multiHeader.SetNameByteSlice(context, StringToByteSlice(name))
		this.AddHeader(context, addr)
	}
	return multiHeader.GenerateAndAddHeader(context, name, value)
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
