package sipparser3

import (
	"bytes"
	//"fmt"
)

type SipMultiHeader struct {
	info    *SipHeaderInfo
	name    AbnfToken
	headers SipSingleHeaders
}

func NewSipMultiHeader() *SipMultiHeader {
	ret := &SipMultiHeader{}
	ret.Init()
	return ret
}

func (this *SipMultiHeader) Init() {
	this.headers.Init()
}

func (this *SipMultiHeader) HasInfo() bool        { return this.info != nil }
func (this *SipMultiHeader) HasShortName() bool   { return this.HasInfo() && this.info.HasShortName() }
func (this *SipMultiHeader) ShortName() AbnfToken { return this.info.ShortName() }

func (this *SipMultiHeader) Size() int { return this.headers.Size() }

func (this *SipMultiHeader) EqualNameBytes(name []byte) bool {
	if this.info != nil {
		if EqualNoCase(this.info.name, name) {
			return true
		}
		return this.info.shortName.EqualBytesNoCase(name)
	}
	return this.name.EqualBytesNoCase(name)
}

func (this *SipMultiHeader) EqualNameString(name string) bool {
	return this.EqualNameBytes(Str2bytes(name))
}

func (this *SipMultiHeader) AddHeader(header *SipSingleHeader) *SipSingleHeader {
	return this.headers.AddHeader(header)
}

func (this *SipMultiHeader) Encode(buf *bytes.Buffer) {
	if this.info != nil {
		buf.Write(this.info.name)
	} else {
		this.name.Encode(buf)
	}
	buf.WriteString(": ")
	for i, v := range this.headers.headers {
		if i > 0 {
			buf.WriteString(", ")
		}
		v.EncodeValue(buf)
	}
}

func (this *SipMultiHeader) String() string {
	return AbnfEncoderToString(this)
}

type SipMultiHeaders struct {
	headers []SipMultiHeader
}

func NewSipMultiHeaders() *SipMultiHeaders {
	ret := &SipMultiHeaders{}
	ret.Init()
	return ret
}

func (this *SipMultiHeaders) Init() {
	this.headers = make([]SipMultiHeader, 0, 4)
}

func (this *SipMultiHeaders) Size() int   { return len(this.headers) }
func (this *SipMultiHeaders) Empty() bool { return this.Size() == 0 }
func (this *SipMultiHeaders) GetHeaderBytes(name []byte) (val *SipMultiHeader, ok bool) {
	for i, v := range this.headers {
		if v.EqualNameBytes(name) {
			return &this.headers[i], true
		}
	}
	return nil, false
}

func (this *SipMultiHeaders) GetHeaderString(name string) (val *SipMultiHeader, ok bool) {
	return this.GetHeaderBytes(Str2bytes(name))
}

func (this *SipMultiHeaders) AddHeader(header *SipMultiHeader) *SipMultiHeader {
	this.headers = append(this.headers, *header)
	return &this.headers[len(this.headers)-1]
}

func (this *SipMultiHeaders) Encode(buf *bytes.Buffer) {
	for _, v := range this.headers {
		v.Encode(buf)
		buf.WriteString("\r\n")
	}
}

func (this *SipMultiHeaders) String() string {
	return AbnfEncoderToString(this)
}
