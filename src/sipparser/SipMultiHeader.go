package sipparser

import (
	//"fmt"
	"unsafe"
)

type SipMultiHeader struct {
	info    *SipHeaderInfo
	name    AbnfBuf
	headers SipSingleHeaders
}

func NewSipMultiHeader(context *ParseContext) AbnfPtr {
	addr := context.allocator.Alloc(uint32(unsafe.Sizeof(SipMultiHeader{})))
	if addr == ABNF_PTR_NIL {
		return ABNF_PTR_NIL
	}
	addr.GetSipMultiHeader(context).Init()
	return addr
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
	this.info, _ = GetSipHeaderInfo(StringToByteSlice(name))
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

func (this *SipMultiHeader) GenerateAndAddHeader(context *ParseContext, name, value string) AbnfPtr {
	addr := GenerateSingleHeader(context, name, value)
	if addr == ABNF_PTR_NIL {
		return ABNF_PTR_NIL
	}

	this.AddHeader(context, addr)
	return addr
}

func (this *SipMultiHeader) Parse(context *ParseContext, src []byte, pos int, info *SipHeaderInfo) (newPos int, err error) {
	this.info = info

	if info.parseFunc != nil && info.needParse && !context.ParseSipHeaderAsRaw {
		return this.parseHeader(context, src, pos, info)

	}

	/* 此时可能存在多个用逗号隔开的同类型头部 */
	return this.parseAsRawHeader(context, src, pos, info)
}

func (this *SipMultiHeader) parseHeader(context *ParseContext, src []byte, pos int, info *SipHeaderInfo) (newPos int, err error) {
	if pos >= len(src) {
		return pos, &AbnfError{"SipMultiHeader  parseHeader: no header", src, newPos}
	}
	newPos = pos
	for newPos < len(src) {
		var addr AbnfPtr

		addr, newPos, err = parseOneSingleHeader(context, src, newPos, info)
		if err != nil {
			return newPos, err
		}

		this.AddHeader(context, addr)

		// now should be COMMA or CRLF
		if IsOnlyCRLF(src, newPos) {
			return newPos + 2, nil
		}
		/*
			var macthMark bool
			var newPos1 int

			newPos1, macthMark, err = ParseSWSMarkCanOmmit(src, newPos, ',')
			if err != nil {
				return newPos, err
			}

			if !macthMark {
				return newPos, nil
			}

			newPos = newPos1
			//*/

		//*
		var newPos1 int
		var macthMark bool

		newPos1, macthMark, err = ParseSWSMarkCanOmmit(src, newPos, ',')
		if err != nil {
			return newPos + 2, nil
		}

		if !macthMark {
			return newPos, nil
		}

		newPos = newPos1 //*/
	}
	return newPos, nil
}

func (this *SipMultiHeader) parseAsRawHeader(context *ParseContext, src []byte, pos int, info *SipHeaderInfo) (newPos int, err error) {
	addr, newPos, err := parseRawHeader(context, AbnfRef{0, int32(len(info.name))}, src, pos, info)
	if err != nil {
		return newPos, err
	}
	this.AddHeader(context, addr)

	return newPos, nil
}

func (this *SipMultiHeader) Encode(context *ParseContext, buf *AbnfByteBuffer) {
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

func NewSipMultiHeaders(context *ParseContext) AbnfPtr {
	addr := context.allocator.Alloc(uint32(unsafe.Sizeof(SipMultiHeaders{})))
	if addr == ABNF_PTR_NIL {
		return ABNF_PTR_NIL
	}
	addr.GetSipMultiHeaders(context).Init()
	return addr
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

func (this *SipMultiHeaders) GetHeaderByIndex(context *ParseContext, headerIndex SipHeaderIndexType) (val *SipMultiHeader, ok bool) {
	for e := this.Front(context); e != nil; e = e.Next(context) {
		v := e.Value.GetSipMultiHeader(context)
		if v.info != nil && v.info.index == headerIndex {
			return v, true
		}
	}
	return nil, false
}

// remove Content-* headers from sip message except Content-Length and Content-Type*/
func (this *SipMultiHeaders) RemoveContentHeaders(context *ParseContext) {
	prefix := StringToByteSlice("Content-")

	for e := this.Front(context); e != nil; {
		v := e.Value.GetSipMultiHeader(context)
		//if v.EqualNameString(context, ABNF_NAME_SIP_HDR_CONTENT_TYPE) ||
		//	v.EqualNameString(context, ABNF_NAME_SIP_HDR_CONTENT_LENGTH) {
		if v.info != nil && (v.info.index == ABNF_SIP_HDR_CONTENT_TYPE || v.info.index == ABNF_SIP_HDR_CONTENT_LENGTH) {
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
		if v.info != nil && (v.info.index == ABNF_SIP_HDR_CONTENT_TYPE || v.info.index == ABNF_SIP_HDR_CONTENT_LENGTH) {
			rhs.AddHeader(context, e.Value)
			continue
		}

		if v.NameHasPrefixBytes(context, prefix) {
			rhs.AddHeader(context, e.Value)
		}
	}
}

func (this *SipMultiHeaders) AddHeader(context *ParseContext, header AbnfPtr) {
	this.PushBack(context, header)
}

func (this *SipMultiHeaders) GenerateAndAddHeader(context *ParseContext, name, value string) AbnfPtr {
	multiHeader, ok := this.GetHeaderByByteSlice(context, StringToByteSlice(name))
	if multiHeader == nil || !ok {
		var addr AbnfPtr
		addr = NewSipMultiHeader(context)
		if addr == ABNF_PTR_NIL {
			return ABNF_PTR_NIL
		}
		multiHeader = addr.GetSipMultiHeader(context)
		multiHeader.SetNameByteSlice(context, StringToByteSlice(name))
		this.AddHeader(context, addr)
	}
	return multiHeader.GenerateAndAddHeader(context, name, value)
}

func (this *SipMultiHeaders) Encode(context *ParseContext, buf *AbnfByteBuffer) {
	for e := this.Front(context); e != nil; e = e.Next(context) {
		v := e.Value.GetSipMultiHeader(context)
		v.Encode(context, buf)

	}
}

func (this *SipMultiHeaders) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}
