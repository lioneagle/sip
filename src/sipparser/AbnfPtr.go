package sipparser

import (
	"strconv"
	"unsafe"
)

type AbnfPtr int32

const ABNF_PTR_NIL = AbnfPtr(-1)

func (this AbnfPtr) GetMemAddr(context *ParseContext) *byte {
	return (*byte)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetUintptr(context *ParseContext) uintptr {
	return (uintptr)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetAbnfToken(context *ParseContext) *AbnfToken {
	return (*AbnfToken)(unsafe.Pointer(&context.allocator.mem[this]))
}

func (this AbnfPtr) GetAbnfListNode(context *ParseContext) *AbnfListNode {
	return (*AbnfListNode)(unsafe.Pointer(&context.allocator.mem[this]))
}

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

func (this AbnfPtr) String() string {
	if this == ABNF_PTR_NIL {
		return "nil"
	}
	return strconv.FormatUint(uint64(this), 10)
}
