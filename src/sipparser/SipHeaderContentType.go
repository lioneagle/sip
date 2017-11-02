package sipparser

import (
	"bytes"
	//"fmt"
	"unsafe"
)

type SipHeaderContentType struct {
	mainType AbnfBuf
	subType  AbnfBuf
	params   SipGenericParams
}

func NewSipHeaderContentType(context *ParseContext) AbnfPtr {
	addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipHeaderContentType{})))
	if addr == ABNF_PTR_NIL {
		return ABNF_PTR_NIL
	}
	addr.GetSipHeaderContentType(context).Init()
	return addr
}

func (this *SipHeaderContentType) Init() {
	this.mainType.Init()
	this.subType.Init()
	this.params.Init()
}

func (this *SipHeaderContentType) AllowMulti() bool { return false }
func (this *SipHeaderContentType) HasValue() bool   { return true }

/* RFC3261
 *
 * Content-Type     =  ( "Content-Type" / "c" ) HCOLON media-type
 * media-type       =  m-type SLASH m-subtype *(SEMI m-parameter)
 * m-type           =  discrete-type / composite-type
 * discrete-type    =  "text" / "image" / "audio" / "video"
 *                     / "application" / extension-token
 * composite-type   =  "message" / "multipart" / extension-token
 * extension-token  =  ietf-token / x-token
 * ietf-token       =  token
 * x-token          =  "x-" token
 * m-subtype        =  extension-token / iana-token
 * iana-token       =  token
 * m-parameter      =  m-attribute EQUAL m-value
 * m-attribute      =  token
 * m-value          =  token / quoted-string
 */
func (this *SipHeaderContentType) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	return this.ParseWithoutInit(context, src, pos)
}

func (this *SipHeaderContentType) ParseWithoutInit(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	name, newPos, err := ParseHeaderName(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if !EqualNoCase(src[name.Begin:name.End], StringToByteSlice(ABNF_NAME_SIP_HDR_CONTENT_TYPE)) &&
		!EqualNoCase(src[name.Begin:name.End], StringToByteSlice(ABNF_NAME_SIP_HDR_CONTENT_TYPE_S)) {
		return newPos, &AbnfError{"Content-Type parse: wrong header-name", src, newPos}
	}

	return this.ParseValueWithoutInit(context, src, newPos)
}

func (this *SipHeaderContentType) ParseValue(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	return this.ParseValueWithoutInit(context, src, pos)
}

func (this *SipHeaderContentType) ParseValueWithoutInit(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos, err = this.mainType.ParseSipToken(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if newPos >= len(src) {
		return newPos, &AbnfError{"Content-Type parse: no slash", src, newPos}
	}

	if src[newPos] != '/' && !IsLwsChar(src[newPos]) {
		return newPos, &AbnfError{"Content-Type parse: no slash", src, newPos}
	}

	newPos, err = ParseSWSMark(src, newPos, '/')
	if err != nil {
		return newPos, err
	}

	newPos, err = this.subType.ParseSipToken(context, src, newPos)
	if err != nil {
		return newPos, err
	}

	return this.params.ParseWithoutInit(context, src, newPos, ';')
}

func (this *SipHeaderContentType) Encode(context *ParseContext, buf *bytes.Buffer) {
	buf.WriteString(ABNF_NAME_SIP_HDR_CONTENT_TYPE_COLON)
	this.EncodeValue(context, buf)
}

func (this *SipHeaderContentType) EncodeValue(context *ParseContext, buf *bytes.Buffer) {
	this.mainType.Encode(context, buf)
	buf.WriteByte('/')
	this.subType.Encode(context, buf)
	this.params.Encode(context, buf, ';')
}

func (this *SipHeaderContentType) GetBoundary(context *ParseContext) (boundary []byte, ok bool) {
	param, ok := this.params.GetParam(context, "boundary")
	if !ok {
		return nil, false
	}

	boundary, ok = param.GetValueAsByteSlice(context)
	if !ok {
		return nil, false
	}
	return boundary, ok
}

func (this *SipHeaderContentType) SetMainType(context *ParseContext, mainType string) {
	this.mainType.SetString(context, mainType)
}

func (this *SipHeaderContentType) SetSubType(context *ParseContext, subType string) {
	this.subType.SetString(context, subType)
}

func (this *SipHeaderContentType) AddBoundary(context *ParseContext, boundary []byte) error {
	addr := NewSipGenericParam(context)
	if addr == ABNF_PTR_NIL {
		return &AbnfError{"Content-Type parse: out of memory for adding boundary", nil, 0}
	}

	addr.GetSipGenericParam(context).SetNameAsString(context, "boundary")
	addr.GetSipGenericParam(context).SetValueQuotedString(context, boundary)
	this.params.PushBack(context, addr)
	return nil
}

func (this *SipHeaderContentType) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

func ParseSipContentType(context *ParseContext, src []byte, pos int) (newPos int, parsed AbnfPtr, err error) {
	addr := NewSipHeaderContentType(context)
	if addr == ABNF_PTR_NIL {
		return newPos, ABNF_PTR_NIL, &AbnfError{"Content-Type parse: out of memory for new header", nil, 0}
	}
	newPos, err = addr.GetSipHeaderContentType(context).ParseValueWithoutInit(context, src, pos)
	return newPos, addr, err
}

func EncodeSipContentTypeValue(parsed AbnfPtr, context *ParseContext, buf *bytes.Buffer) {
	if parsed == ABNF_PTR_NIL {
		return
	}
	parsed.GetSipHeaderContentType(context).EncodeValue(context, buf)
}
