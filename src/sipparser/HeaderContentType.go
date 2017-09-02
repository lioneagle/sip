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

func NewSipHeaderContentType(context *ParseContext) (*SipHeaderContentType, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipHeaderContentType{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*SipHeaderContentType)(unsafe.Pointer(mem)).Init()
	return (*SipHeaderContentType)(unsafe.Pointer(mem)), addr
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
	name, newPos, err := ParseHeaderName(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if !EqualNoCase(src[name.Begin:name.End], StringToByteSlice(ABNF_NAME_SIP_HDR_CONTENT_TYPE)) &&
		!EqualNoCase(src[name.Begin:name.End], StringToByteSlice(ABNF_NAME_SIP_HDR_CONTENT_TYPE_S)) {
		return newPos, &AbnfError{"Content-Type parse: wrong header-name", src, newPos}
	}

	return this.ParseValue(context, src, newPos)
}

func (this *SipHeaderContentType) ParseValue(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	newPos, err = this.mainType.ParseSipToken(context, src, pos)
	if err != nil {
		return newPos, err
	}

	newPos, err = ParseSWSMark(src, newPos, '/')
	if err != nil {
		return newPos, err
	}

	newPos, err = this.subType.ParseSipToken(context, src, newPos)
	if err != nil {
		return newPos, err
	}

	return this.params.Parse(context, src, newPos, ';')
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
	param, addr := NewSipGenericParam(context)
	if param == nil {
		return &AbnfError{"Content-Type parse: out of memory for adding boundary", nil, 0}
	}
	param.SetNameAsString(context, "boundary")
	param.SetValueQuotedString(context, boundary)
	this.params.PushBack(context, addr)
	return nil
}

func (this *SipHeaderContentType) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

func ParseSipContentType(context *ParseContext, src []byte, pos int) (newPos int, parsed AbnfPtr, err error) {
	header, addr := NewSipHeaderContentType(context)
	if header == nil || addr == ABNF_PTR_NIL {
		return newPos, ABNF_PTR_NIL, &AbnfError{"Content-Type parse: out of memory for new header", nil, 0}
	}
	newPos, err = header.ParseValue(context, src, pos)
	return newPos, addr, err
}

func EncodeSipContentTypeValue(parsed AbnfPtr, context *ParseContext, buf *bytes.Buffer) {
	if parsed == ABNF_PTR_NIL {
		return
	}
	parsed.GetSipHeaderContentType(context).EncodeValue(context, buf)
}
