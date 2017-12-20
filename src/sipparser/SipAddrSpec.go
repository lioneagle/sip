package sipparser

import (
	//"bytes"
	//"fmt"
	"unsafe"
)

type SipAddrSpec struct {
	uriType int32
	uri     AbnfPtr
}

func NewSipAddrSpec(context *ParseContext) AbnfPtr {
	addr := context.allocator.Alloc(uint32(unsafe.Sizeof(SipAddrSpec{})))
	if addr == ABNF_PTR_NIL {
		return ABNF_PTR_NIL
	}

	addr.GetSipAddrSpec(context).Init()
	return addr
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
	return this.ParseWithoutInit(context, src, pos)
}

func (this *SipAddrSpec) ParseWithoutInit(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	var scheme AbnfBuf

	newPos, err = ParseUriScheme(context, src, pos, &scheme)
	if err != nil {
		return newPos, err
	}

	if scheme.EqualStringNoCase(context, ABNF_NAME_URI_SCHEME_SIP) {
		addr := NewSipUri(context)
		if addr == ABNF_PTR_NIL {
			return newPos, &AbnfError{"addr-spec parse: out of memory before parse sip uri", src, newPos}
		}
		sipuri := addr.GetSipUri(context)
		sipuri.SetSipUri()
		this.uri = addr
		this.uriType = ABNF_SIP_URI
		return sipuri.ParseAfterSchemeWithoutInit(context, src, newPos)
	}

	if scheme.EqualStringNoCase(context, ABNF_NAME_URI_SCHEME_SIPS) {
		addr := NewSipUri(context)
		if addr == ABNF_PTR_NIL {
			return newPos, &AbnfError{"addr-spec parse: out of memory before parse sips uri", src, newPos}
		}
		sipuri := addr.GetSipUri(context)
		sipuri.SetSipsUri()
		this.uri = addr
		this.uriType = ABNF_SIPS_URI
		return sipuri.ParseAfterSchemeWithoutInit(context, src, newPos)
	}

	if scheme.EqualStringNoCase(context, ABNF_NAME_URI_SCHEME_TEL) {
		addr := NewTelUri(context)
		if addr == ABNF_PTR_NIL {
			return newPos, &AbnfError{"addr-spec parse: out of memory before parse tel uri", src, newPos}
		}
		teluri := addr.GetTelUri(context)
		this.uri = addr
		this.uriType = ABNF_TEL_URI
		return teluri.ParseAfterSchemeWithoutInit(context, src, newPos)
	}

	return newPos, &AbnfError{"addr-spec parse: unsupported uri", src, newPos}
}

func (this *SipAddrSpec) ParseWithoutParam(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	return this.ParseWithoutParamNorInit(context, src, pos)
}

func (this *SipAddrSpec) ParseWithoutParamNorInit(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	var scheme AbnfBuf

	newPos, err = ParseUriScheme(context, src, pos, &scheme)
	if err != nil {
		return newPos, err
	}

	if scheme.EqualStringNoCase(context, "sip") {
		addr := NewSipUri(context)
		if addr == ABNF_PTR_NIL {
			return newPos, &AbnfError{"addr-spec parse: out of memory before parse sip uri", src, newPos}
		}
		sipuri := addr.GetSipUri(context)
		sipuri.SetSipUri()
		this.uri = addr
		this.uriType = ABNF_SIP_URI
		return sipuri.ParseAfterSchemeWithoutParam(context, src, newPos)
	}

	if scheme.EqualStringNoCase(context, "sips") {
		addr := NewSipUri(context)
		if addr == ABNF_PTR_NIL {
			return newPos, &AbnfError{"addr-spec parse: out of memory before parse sips uri", src, newPos}
		}
		sipuri := addr.GetSipUri(context)
		sipuri.SetSipsUri()
		this.uri = addr
		this.uriType = ABNF_SIPS_URI
		return sipuri.ParseAfterSchemeWithoutParam(context, src, newPos)
	}

	if scheme.EqualStringNoCase(context, "tel") {
		addr := NewTelUri(context)
		if addr == ABNF_PTR_NIL {
			return newPos, &AbnfError{"addr-spec parse: out of memory before parse tel uri", src, newPos}
		}
		teluri := addr.GetTelUri(context)
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

func (this *SipAddrSpec) Encode(context *ParseContext, buf *AbnfByteBuffer) {
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
