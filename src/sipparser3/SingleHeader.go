package sipparser3

import (
	"bytes"
	//"fmt"
)

type SipSingleHeader struct {
	info   *SipHeaderInfo
	name   AbnfToken
	value  AbnfToken
	parsed SipHeaderParsed
}

func NewSipSingleHeader() *SipSingleHeader {
	ret := &SipSingleHeader{}
	ret.Init()
	return ret
}

func (this *SipSingleHeader) Init() {}

func (this *SipSingleHeader) HasInfo() bool        { return this.info != nil }
func (this *SipSingleHeader) HasShortName() bool   { return this.HasInfo() && this.info.HasShortName() }
func (this *SipSingleHeader) ShortName() AbnfToken { return this.info.ShortName() }

func (this *SipSingleHeader) EqualNameBytes(name []byte) bool {
	if this.info != nil {
		if EqualNoCase(this.info.name, name) {
			return true
		}
		return this.info.shortName.EqualBytesNoCase(name)
	}
	return this.name.EqualBytesNoCase(name)
}

func (this *SipSingleHeader) EqualNameString(name string) bool {
	return this.EqualNameBytes(Str2bytes(name))
}

func (this *SipSingleHeader) Encode(buf *bytes.Buffer) {
	if this.info != nil {
		buf.Write(this.info.name)
	} else {
		this.name.Encode(buf)
	}
	buf.WriteString(": ")
	this.EncodeValue(buf)
}

func (this *SipSingleHeader) EncodeValue(buf *bytes.Buffer) {
	if this.value.Exist() {
		this.value.Encode(buf)
	}
}

func (this *SipSingleHeader) String() string {
	return AbnfEncoderToString(this)
}

type SipSingleHeaders struct {
	headers []SipSingleHeader
}

func NewSipSingleHeaders() *SipSingleHeaders {
	ret := &SipSingleHeaders{}
	ret.Init()
	return ret
}

func (this *SipSingleHeaders) Init() {
	this.headers = make([]SipSingleHeader, 0, 5)
}

func (this *SipSingleHeaders) Size() int   { return len(this.headers) }
func (this *SipSingleHeaders) Empty() bool { return this.Size() == 0 }
func (this *SipSingleHeaders) GetHeaderBytes(name []byte) (val *SipSingleHeader, ok bool) {
	for i, v := range this.headers {
		if v.EqualNameBytes(name) {
			return &this.headers[i], true
		}
	}
	return nil, false
}

func (this *SipSingleHeaders) GetHeaderString(name string) (val *SipSingleHeader, ok bool) {
	return this.GetHeaderBytes(Str2bytes(name))
}

func (this *SipSingleHeaders) AddHeader(header *SipSingleHeader) *SipSingleHeader {
	this.headers = append(this.headers, *header)
	return &this.headers[len(this.headers)-1]
}

func (this *SipSingleHeaders) Encode(buf *bytes.Buffer) {
	for _, v := range this.headers {
		v.Encode(buf)
		buf.WriteString("\r\n")
	}
}

func (this *SipSingleHeaders) String() string {
	return AbnfEncoderToString(this)
}
