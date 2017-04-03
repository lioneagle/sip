package sipparser

import (
	"bytes"
	//"fmt"
	"unsafe"
)

type SipAddrSpec struct {
	uriType int32
	uri     AbnfPtr
}

func NewSipAddrSpec(context *ParseContext) (*SipAddrSpec, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipAddrSpec{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*SipAddrSpec)(unsafe.Pointer(mem)).Init()
	return (*SipAddrSpec)(unsafe.Pointer(mem)), addr
}

func (this *SipAddrSpec) Init() {
	this.uriType = ABNF_UNKNOWN_URI
	this.uri = ABNF_PTR_NIL
}

func (this *SipAddrSpec) Equal(context *ParseContext, rhs *SipAddrSpec) bool {
	if this.uriType != rhs.uriType {
		return false
	}

	switch this.uriType {
	case ABNF_SIP_URI:
		fallthrough
	case ABNF_SIPS_URI:
		return this.uri.GetSipUri(context).Equal(context, rhs.uri.GetSipUri(context))
	case ABNF_TEL_URI:
		return this.uri.GetTelUri(context).Equal(context, rhs.uri.GetTelUri(context))
	default:
		break
	}
	return false

}

func (this *SipAddrSpec) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()

	newPos, scheme, err := ParseUriScheme(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if scheme.EqualStringNoCase(context, "sip") {
		sipuri, addr := NewSipUri(context)
		if sipuri == nil {
			return newPos, &AbnfError{"addr-spec parse: out of memory before parse sip uri", src, newPos}
		}
		sipuri.SetSipUri()
		this.uri = addr
		this.uriType = ABNF_SIP_URI
		return sipuri.ParseAfterScheme(context, src, newPos)
	}

	if scheme.EqualStringNoCase(context, "sips") {
		sipuri, addr := NewSipUri(context)
		if sipuri == nil {
			return newPos, &AbnfError{"addr-spec parse: out of memory before parse sips uri", src, newPos}
		}
		sipuri.SetSipsUri()
		this.uri = addr
		this.uriType = ABNF_SIPS_URI
		return sipuri.ParseAfterScheme(context, src, newPos)
	}

	if scheme.EqualStringNoCase(context, "tel") {
		teluri, addr := NewTelUri(context)
		if teluri == nil {
			return newPos, &AbnfError{"addr-spec parse: out of memory before parse tel uri", src, newPos}
		}
		this.uri = addr
		this.uriType = ABNF_TEL_URI
		return teluri.ParseAfterScheme(context, src, newPos)
	}

	return newPos, &AbnfError{"addr-spec parse: unsupported uri", src, newPos}
}

func (this *SipAddrSpec) ParseWithoutParam(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()

	newPos, scheme, err := ParseUriScheme(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if scheme.EqualStringNoCase(context, "sip") {
		sipuri, addr := NewSipUri(context)
		if sipuri == nil {
			return newPos, &AbnfError{"addr-spec parse: out of memory before parse sip uri", src, newPos}
		}
		sipuri.SetSipUri()
		this.uri = addr
		this.uriType = ABNF_SIP_URI
		return sipuri.ParseAfterSchemeWithoutParam(context, src, newPos)
	}

	if scheme.EqualStringNoCase(context, "sips") {
		sipuri, addr := NewSipUri(context)
		if sipuri == nil {
			return newPos, &AbnfError{"addr-spec parse: out of memory before parse sips uri", src, newPos}
		}
		sipuri.SetSipsUri()
		this.uri = addr
		this.uriType = ABNF_SIPS_URI
		return sipuri.ParseAfterSchemeWithoutParam(context, src, newPos)
	}

	if scheme.EqualStringNoCase(context, "tel") {
		teluri, addr := NewTelUri(context)
		if teluri == nil {
			return newPos, &AbnfError{"addr-spec parse: out of memory before parse tel uri", src, newPos}
		}
		this.uri = addr
		this.uriType = ABNF_TEL_URI
		return teluri.ParseAfterSchemeWithoutParam(context, src, newPos)
	}

	return newPos, &AbnfError{"addr-spec parse: unsupported uri", src, newPos}
}

func (this *SipAddrSpec) IsSipUri(context *ParseContext) (uri *SipUri, ok bool) {
	if this.uriType == ABNF_SIP_URI {
		ok = true
		uri = this.uri.GetSipUri(context)
	}
	return uri, ok
}

func (this *SipAddrSpec) IsTelUri(context *ParseContext) (uri *TelUri, ok bool) {
	if this.uriType == ABNF_SIP_URI {
		ok = true
		uri = this.uri.GetTelUri(context)
	}
	return uri, ok
}

func (this *SipAddrSpec) Encode(context *ParseContext, buf *bytes.Buffer) {
	if this.uri != ABNF_PTR_NIL {
		switch this.uriType {
		case ABNF_SIP_URI:
			fallthrough
		case ABNF_SIPS_URI:
			this.uri.GetSipUri(context).Encode(context, buf)
		case ABNF_TEL_URI:
			this.uri.GetTelUri(context).Encode(context, buf)
		default:
			break
		}
	}

}

func (this *SipAddrSpec) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}
