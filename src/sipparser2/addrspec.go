package sipparser2

import (
	//"fmt"
	//"strings"
	"bytes"
)

type SipAddrSpec struct {
	uri URI
}

func NewSipAddrSpec() *SipAddrSpec {
	return &SipAddrSpec{}
}

func (this *SipAddrSpec) Encode(buf *bytes.Buffer) {
	this.uri.Encode(buf)
}

func (this *SipAddrSpec) String() string {
	return this.uri.String()
}

func (this *SipAddrSpec) Equal(rhs *SipAddrSpec) bool {
	return this.uri.Equal(rhs.uri)
}

func (this *SipAddrSpec) Parse(src []byte, pos int) (newPos int, err error) {
	newPos, scheme, err := ParseUriScheme(src, pos)
	if err != nil {
		return newPos, err
	}

	if scheme.EqualStringNoCase("sip") || scheme.EqualStringNoCase("sips") {
		sipuri := NewSipUri()
		this.uri = sipuri
		return sipuri.Parse(src, pos)
	}

	if scheme.EqualStringNoCase("tel") {
		teluri := NewTelUri()
		this.uri = teluri
		return teluri.Parse(src, pos)
	}

	return newPos, &AbnfError{"addr-spec parse: unsupported uri", src, newPos}
}

func (this *SipAddrSpec) IsSipUri() (sipuri *SipUri, ok bool) {
	sipuri, ok = this.uri.(*SipUri)
	return sipuri, ok
}

func (this *SipAddrSpec) IsTelUri() (teluri *TelUri, ok bool) {
	teluri, ok = this.uri.(*TelUri)
	return teluri, ok
}
