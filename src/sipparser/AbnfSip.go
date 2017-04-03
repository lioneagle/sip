package sipparser

import (
	"unsafe"
)

const (
	// basic
	ABNF_NAME_SIP_HCOLON = ABNF_NAME_COLON + ABNF_NAME_SPACE

	// uri scheme
	ABNF_NAME_URI_SCHEME_SIP  = "sip"
	ABNF_NAME_URI_SCHEME_SIPS = "sips"
	ABNF_NAME_URI_SCHEME_TEL  = "tel"

	// header names
	ABNF_NAME_SIP_HDR_FROM                        = "From"
	ABNF_NAME_SIP_HDR_FROM_S                      = "f"
	ABNF_NAME_SIP_HDR_FROM_COLON                  = ABNF_NAME_SIP_HDR_FROM + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_FROM_S_COLON                = ABNF_NAME_SIP_HDR_FROM_S + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_TO                          = "To"
	ABNF_NAME_SIP_HDR_TO_S                        = "t"
	ABNF_NAME_SIP_HDR_TO_COLON                    = ABNF_NAME_SIP_HDR_TO + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_TO_S_COLON                  = ABNF_NAME_SIP_HDR_TO_S + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_VIA                         = "Via"
	ABNF_NAME_SIP_HDR_VIA_S                       = "v"
	ABNF_NAME_SIP_HDR_VIA_COLON                   = ABNF_NAME_SIP_HDR_VIA + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_VIA_S_COLON                 = ABNF_NAME_SIP_HDR_VIA_S + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_CALL_ID                     = "Call-ID"
	ABNF_NAME_SIP_HDR_CALL_ID_S                   = "i"
	ABNF_NAME_SIP_HDR_CALL_ID_COLON               = ABNF_NAME_SIP_HDR_CALL_ID + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_CALL_ID_S_COLON             = ABNF_NAME_SIP_HDR_CALL_ID_S + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_CONTACT_ID                  = "Contact"
	ABNF_NAME_SIP_HDR_CONTACT_ID_S                = "m"
	ABNF_NAME_SIP_HDR_CONTACT_ID_COLON            = ABNF_NAME_SIP_HDR_CONTACT_ID + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_CONTACT_ID_S_COLON          = ABNF_NAME_SIP_HDR_CONTACT_ID_S + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_CONTENT_ENCODING            = "Content-Encoding"
	ABNF_NAME_SIP_HDR_CONTENT_ENCODING_S          = "e"
	ABNF_NAME_SIP_HDR_CONTENT_ENCODING_COLON      = ABNF_NAME_SIP_HDR_CONTENT_ENCODING + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_CONTENT_ENCODING_S_COLON    = ABNF_NAME_SIP_HDR_CONTENT_ENCODING_S + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_CONTENT_LENGTH              = "Content-Length"
	ABNF_NAME_SIP_HDR_CONTENT_LENGTH_S            = "l"
	ABNF_NAME_SIP_HDR_CONTENT_LENGTH_COLON        = ABNF_NAME_SIP_HDR_CONTENT_LENGTH + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_CONTENT_LENGTH_S_COLON      = ABNF_NAME_SIP_HDR_CONTENT_LENGTH_S + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_CONTENT_TYPE                = "Content-Type"
	ABNF_NAME_SIP_HDR_CONTENT_TYPE_S              = "c"
	ABNF_NAME_SIP_HDR_CONTENT_TYPE_COLON          = ABNF_NAME_SIP_HDR_CONTENT_TYPE + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_CONTENT_TYPE_S_COLON        = ABNF_NAME_SIP_HDR_CONTENT_TYPE_S + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_SUBJECTE                    = "Subject"
	ABNF_NAME_SIP_HDR_SUBJECTE_S                  = "s"
	ABNF_NAME_SIP_HDR_SUBJECTE_COLON              = ABNF_NAME_SIP_HDR_SUBJECTE + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_SUBJECTE_S_COLON            = ABNF_NAME_SIP_HDR_SUBJECTE_S + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_SUPPORTED                   = "Supported"
	ABNF_NAME_SIP_HDR_SUPPORTED_S                 = "k"
	ABNF_NAME_SIP_HDR_SUPPORTED_COLON             = ABNF_NAME_SIP_HDR_SUPPORTED + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_SUPPORTED_S_COLON           = ABNF_NAME_SIP_HDR_SUPPORTED_S + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_ALLOW_EVENTS                = "Allow-Events"
	ABNF_NAME_SIP_HDR_ALLOW_EVENTS_S              = "u"
	ABNF_NAME_SIP_HDR_ALLOW_EVENTS_COLON          = ABNF_NAME_SIP_HDR_ALLOW_EVENTS + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_ALLOW_EVENTS_S_COLON        = ABNF_NAME_SIP_HDR_ALLOW_EVENTS_S + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_EVENT                       = "Event"
	ABNF_NAME_SIP_HDR_EVENT_S                     = "o"
	ABNF_NAME_SIP_HDR_EVENT_COLON                 = ABNF_NAME_SIP_HDR_EVENT + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_EVENT_S_COLON               = ABNF_NAME_SIP_HDR_EVENT_S + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_REFER_TO                    = "Refer-To"
	ABNF_NAME_SIP_HDR_REFER_TO_S                  = "r"
	ABNF_NAME_SIP_HDR_REFER_TO_COLON              = ABNF_NAME_SIP_HDR_REFER_TO + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_REFER_TO_S_COLON            = ABNF_NAME_SIP_HDR_REFER_TO_S + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_ACCEPT_CONTACT              = "Accept-Contact"
	ABNF_NAME_SIP_HDR_ACCEPT_CONTACT_S            = "a"
	ABNF_NAME_SIP_HDR_ACCEPT_CONTACT_COLON        = ABNF_NAME_SIP_HDR_ACCEPT_CONTACT + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_ACCEPT_CONTACT_S_COLON      = ABNF_NAME_SIP_HDR_ACCEPT_CONTACT_S + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_REJECT_CONTACT              = "Reject-Contact"
	ABNF_NAME_SIP_HDR_REJECT_CONTACT_S            = "j"
	ABNF_NAME_SIP_HDR_REJECT_CONTACT_COLON        = ABNF_NAME_SIP_HDR_REJECT_CONTACT + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_REJECT_CONTACT_S_COLON      = ABNF_NAME_SIP_HDR_REJECT_CONTACT_S + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_REQUEST_DISPOSITION         = "Request-Disposition"
	ABNF_NAME_SIP_HDR_REQUEST_DISPOSITION_S       = "d"
	ABNF_NAME_SIP_HDR_REQUEST_DISPOSITION_COLON   = ABNF_NAME_SIP_HDR_REQUEST_DISPOSITION + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_REQUEST_DISPOSITION_S_COLON = ABNF_NAME_SIP_HDR_REQUEST_DISPOSITION_S + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_REFERRED_BY                 = "Referred-By"
	ABNF_NAME_SIP_HDR_REFERRED_BY_S               = "b"
	ABNF_NAME_SIP_HDR_REFERRED_BY_COLON           = ABNF_NAME_SIP_HDR_REFERRED_BY + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_REFERRED_BY_S_COLON         = ABNF_NAME_SIP_HDR_REFERRED_BY_S + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_SESSION_EXPIRES             = "Session-Expires"
	ABNF_NAME_SIP_HDR_SESSION_EXPIRES_S           = "x"
	ABNF_NAME_SIP_HDR_SESSION_EXPIRES_COLON       = ABNF_NAME_SIP_HDR_SESSION_EXPIRES + ABNF_NAME_SIP_HCOLON
	ABNF_NAME_SIP_HDR_SESSION_EXPIRES_S_COLON     = ABNF_NAME_SIP_HDR_SESSION_EXPIRES_S + ABNF_NAME_SIP_HCOLON
)

const (
	ABNF_UNKNOWN_URI   = int32(0)
	ABNF_SIP_URI       = int32(1)
	ABNF_SIPS_URI      = int32(2)
	ABNF_TEL_URI       = int32(3)
	ABNF_ABSOULUTE_URI = int32(4)
)

const (
	ABNF_SIP_ADDR_SPEC = int32(0)
	ABNF_SIP_NAME_ADDR = int32(1)
)

func (this AbnfPtr) GetSipSingleHeader(context *ParseContext) *SipSingleHeader {
	return (*SipSingleHeader)(unsafe.Pointer(&context.allocator.mem[this]))
}
func (this AbnfPtr) GetSipSingleHeaders(context *ParseContext) *SipSingleHeaders {
	return (*SipSingleHeaders)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetSipMultiHeader(context *ParseContext) *SipMultiHeader {
	return (*SipMultiHeader)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetSipMultiHeaders(context *ParseContext) *SipMultiHeaders {
	return (*SipMultiHeaders)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetSipHost(context *ParseContext) *SipHost {
	return (*SipHost)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetSipHostPort(context *ParseContext) *SipHostPort {
	return (*SipHostPort)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetSipUriParam(context *ParseContext) *SipUriParam {
	return (*SipUriParam)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetSipUriParams(context *ParseContext) *SipUriParams {
	return (*SipUriParams)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetSipUriHeader(context *ParseContext) *SipUriHeader {
	return (*SipUriHeader)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetSipUriHeaders(context *ParseContext) *SipUriHeaders {
	return (*SipUriHeaders)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetSipUri(context *ParseContext) *SipUri {
	return (*SipUri)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetTelUriContext(context *ParseContext) *TelUriContext {
	return (*TelUriContext)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetTelUriParam(context *ParseContext) *TelUriParam {
	return (*TelUriParam)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetTelUriParams(context *ParseContext) *TelUriParams {
	return (*TelUriParams)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetTelUri(context *ParseContext) *TelUri {
	return (*TelUri)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetSipAddrSpec(context *ParseContext) *SipAddrSpec {
	return (*SipAddrSpec)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetSipDisplayName(context *ParseContext) *SipDisplayName {
	return (*SipDisplayName)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetSipQuotedString(context *ParseContext) *SipQuotedString {
	return (*SipQuotedString)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetSipNameAddr(context *ParseContext) *SipNameAddr {
	return (*SipNameAddr)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetSipAddr(context *ParseContext) *SipAddr {
	return (*SipAddr)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetSipGenericParam(context *ParseContext) *SipGenericParam {
	return (*SipGenericParam)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetSipGenericParams(context *ParseContext) *SipGenericParams {
	return (*SipGenericParams)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetSipVersion(context *ParseContext) *SipVersion {
	return (*SipVersion)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetSipHeaderVia(context *ParseContext) *SipHeaderVia {
	return (*SipHeaderVia)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetSipHeaderTo(context *ParseContext) *SipHeaderTo {
	return (*SipHeaderTo)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetSipHeaderFrom(context *ParseContext) *SipHeaderFrom {
	return (*SipHeaderFrom)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetSipHeaderCseq(context *ParseContext) *SipHeaderCseq {
	return (*SipHeaderCseq)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetSipHeaderContentType(context *ParseContext) *SipHeaderContentType {
	return (*SipHeaderContentType)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetSipHeaderContentLength(context *ParseContext) *SipHeaderContentLength {
	return (*SipHeaderContentLength)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetSipHeaderContentDisposition(context *ParseContext) *SipHeaderContentDisposition {
	return (*SipHeaderContentDisposition)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetSipHeaderContact(context *ParseContext) *SipHeaderContact {
	return (*SipHeaderContact)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetSipHeaderCallId(context *ParseContext) *SipHeaderCallId {
	return (*SipHeaderCallId)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetSipHeaderRecordRoute(context *ParseContext) *SipHeaderRecordRoute {
	return (*SipHeaderRecordRoute)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetSipHeaderRoute(context *ParseContext) *SipHeaderRoute {
	return (*SipHeaderRoute)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetSipHeaderMaxForwards(context *ParseContext) *SipHeaderMaxForwards {
	return (*SipHeaderMaxForwards)(unsafe.Pointer(&context.allocator.mem[this]))
}

///////////////////////////////////////////////

func (this AbnfPtr) GetSipMsg(context *ParseContext) *SipMsg {
	return (*SipMsg)(unsafe.Pointer(&context.allocator.mem[this]))
}
