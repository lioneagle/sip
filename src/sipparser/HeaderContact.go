package sipparser

import (
	"bytes"
	//"fmt"
	"unsafe"
)

type SipHeaderContact struct {
	isStar bool
	addr   SipAddr
	params SipGenericParams
}

func NewSipHeaderContact(context *ParseContext) (*SipHeaderContact, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipHeaderContact{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*SipHeaderContact)(unsafe.Pointer(mem)).Init()
	return (*SipHeaderContact)(unsafe.Pointer(mem)), addr
}

func (this *SipHeaderContact) Init() {
	this.isStar = false
	this.addr.Init()
	this.params.Init()
}

func (this *SipHeaderContact) AllowMulti() bool { return false }
func (this *SipHeaderContact) HasValue() bool   { return true }

/* RFC3261
 *
 * Contact        =  ("Contact" / "m" ) HCOLON
 *                   ( STAR / (contact-param *(COMMA contact-param)))
 * contact-param  =  (name-addr / addr-spec) *(SEMI contact-params)
 * name-addr      =  [ display-name ] LAQUOT addr-spec RAQUOT
 * addr-spec      =  SIP-URI / SIPS-URI / absoluteURI
 * display-name   =  *(token LWS)/ quoted-string

 * contact-params     =  c-p-q / c-p-expires
 *                      / contact-extension
 * c-p-q              =  "q" EQUAL qvalue
 * c-p-expires        =  "expires" EQUAL delta-seconds
 * contact-extension  =  generic-param
 * delta-seconds      =  1*DIGIT
 *
 * RFC3840, page 13
 *
 * feature-param    =  enc-feature-tag [EQUAL LDQUOT (tag-value-list
 *                     / string-value ) RDQUOT]
 * enc-feature-tag  =  base-tags / other-tags
 * base-tags        =  "audio" / "automata" /
 *                     "class" / "duplex" / "data" /
 *                     "control" / "mobility" / "description" /
 *                     "events" / "priority" / "methods" /
 *                     "schemes" / "application" / "video" /
 *                     "language" / "type" / "isfocus" /
 *                     "actor" / "text" / "extensions"
 * other-tags      =  "+" ftag-name
 * ftag-name       =  ALPHA *( ALPHA / DIGIT / "!" / "'" /
 *                    "." / "-" / "%" )
 * tag-value-list  =  tag-value *("," tag-value)
 * tag-value       =  ["!"] (token-nobang / boolean / numeric)
 * token-nobang    =  1*(alphanum / "-" / "." / "%" / "*"
 *                    / "_" / "+" / "`" / "'" / "~" )
 * boolean         =  "TRUE" / "FALSE"
 * numeric         =  "#" numeric-relation number
 * numeric-relation  =  ">=" / "<=" / "=" / (number ":")
 * number          =  [ "+" / "-" ] 1*DIGIT ["." 0*DIGIT]
 * string-value    =  "<" *(qdtext-no-abkt / quoted-pair ) ">"
 * qdtext-no-abkt  =  LWS / %x21 / %x23-3B / %x3D
 *                    / %x3F-5B / %x5D-7E / UTF8-NONASCII
 *
 * draft-ietf-sip-gruu-15.txt, page 22
 *
 * contact-params  =/ temp-gruu / pub-gruu
 * temp-gruu       =  "temp-gruu" EQUAL LDQUOT *(qdtext / quoted-pair )
 *                    RDQUOT
 * pub-gruu        =  "pub-gruu" EQUAL LDQUOT *(qdtext / quoted-pair )
 *                    RDQUOT
 *
 * uri-parameter   =/ gr-param
 * gr-param        = "gr" ["=" pvalue]   ; defined in RFC3261
 *
 * draft-ietf-sip-outbound-10.txt
 *
 * c-p-reg        = "reg-id" EQUAL 1*DIGIT ; 1 to 2**31
 * c-p-instance   =  "+sip.instance" EQUAL
 *                   LDQUOT "<" instance-val ">" RDQUOT
 * instance-val   = *uric ; defined in RFC 2396
 *
 */
func (this *SipHeaderContact) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	name, newPos, err := ParseHeaderName(context, src, pos)
	if err != nil {
		return newPos, err
	}

	if !EqualNoCase(src[name.Begin:name.End], StringToByteSlice("Contact")) && !EqualNoCase(src[name.Begin:name.End], StringToByteSlice("m")) {
		return newPos, &AbnfError{"Contact parse: wrong header-name", src, newPos}
	}

	return this.ParseValue(context, src, newPos)
}

func (this *SipHeaderContact) ParseValue(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	newPos = pos
	newPos, err = ParseSWS(src, newPos)
	if err != nil {
		return newPos, err
	}

	if newPos >= len(src) {
		return newPos, &AbnfError{"Contact parse: empty", src, newPos}
	}

	if src[newPos] == '*' {
		this.isStar = true
		return ParseSWS(src, newPos+1)
	}

	newPos, err = this.addr.Parse(context, src, newPos)
	if err != nil {
		return newPos, err
	}

	return this.params.Parse(context, src, newPos, ';')
}

func (this *SipHeaderContact) Encode(context *ParseContext, buf *bytes.Buffer) {
	buf.WriteString("Contact: ")
	if this.isStar {
		buf.WriteByte('*')
	} else {
		this.addr.Encode(context, buf)
		this.params.Encode(context, buf, ';')
	}
}

func (this *SipHeaderContact) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

func ParseSipContact(context *ParseContext, src []byte, pos int) (newPos int, parsed AbnfPtr, err error) {
	header, addr := NewSipHeaderContact(context)
	if header == nil || addr == ABNF_PTR_NIL {
		return newPos, ABNF_PTR_NIL, &AbnfError{"Contact parse: out of memory for new header", src, newPos}
	}
	newPos, err = header.ParseValue(context, src, pos)
	return newPos, addr, err
}
