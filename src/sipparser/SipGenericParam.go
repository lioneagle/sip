package sipparser

import (
	"bytes"
	//"fmt"
	"unsafe"
)

const (
	SIP_GENERIC_VALUE_TYPE_NOT_EXIST     = 0
	SIP_GENERIC_VALUE_TYPE_TOKEN         = 1
	SIP_GENERIC_VALUE_TYPE_QUOTED_STRING = 2
	SIP_GENERIC_VALUE_TYPE_IPV6          = 3
)

type SipGenericParam struct {
	name      AbnfBuf
	valueType int32
	value     AbnfPtr
	//parsed    AbnfPtr
}

func NewSipGenericParam(context *ParseContext) AbnfPtr {
	addr := context.allocator.Alloc(uint32(unsafe.Sizeof(SipGenericParam{})))
	if addr == ABNF_PTR_NIL {
		return ABNF_PTR_NIL
	}

	addr.GetSipGenericParam(context).Init()
	return addr
}

func (this *SipGenericParam) Init() {
	this.name.Init()
	this.valueType = SIP_GENERIC_VALUE_TYPE_NOT_EXIST
	this.value = ABNF_PTR_NIL
	//this.parsed = ABNF_PTR_NIL
}

/*
 * generic-param  =  token [ EQUAL gen-value ]
 * gen-value      =  token / host / quoted-string
 *
 */
func (this *SipGenericParam) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	return this.ParseWithoutInit(context, src, pos)
}

func (this *SipGenericParam) ParseWithoutInit(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos, err = this.name.ParseEscapable(context, src, pos, IsSipPname)
	if err != nil {
		return newPos, err
	}

	if newPos >= len(src) {
		return newPos, nil
	}

	if src[newPos] != '=' && !IsLwsChar(src[newPos]) {
		return newPos, nil
	}

	var matchMark bool

	newPos, matchMark, err = ParseSWSMarkCanOmmit(src, newPos, '=')
	if err != nil {
		return newPos, err
	}

	if !matchMark {
		return newPos, nil
	}

	return this.ParseValue(context, src, newPos)
}

func (this *SipGenericParam) ParseValue(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos
	if newPos >= len(src) {
		return newPos, &AbnfError{"generic-param ParseValue: empty gen-value", src, newPos}
	}

	/* @@TODO: 目前解gen-value时，暂不考虑解析出host，因为一般没有必要解析出来，以后再考虑添加这个功能 */
	if IsSipToken(src[newPos]) {
		return this.parseValueToken(context, src, newPos)
	} else if (src[newPos] == '"') || IsLwsChar(src[newPos]) {
		return this.parseValueQuotedString(context, src, newPos)
	} else if src[newPos] == '[' {
		return this.parseValueIpv6(context, src, newPos)
	}

	return newPos, &AbnfError{"generic-param ParseValue: not token nor quoted-string", src, newPos}
}

func (this *SipGenericParam) parseValueToken(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos
	addr := NewAbnfBuf(context)
	if addr == ABNF_PTR_NIL {
		return newPos, &AbnfError{"generic-param  ParseValue: out of memory for token value", src, newPos}
	}
	newPos, err = addr.GetAbnfBuf(context).ParseSipToken(context, src, newPos)
	if err != nil {
		return newPos, err
	}
	this.valueType = SIP_GENERIC_VALUE_TYPE_TOKEN
	this.value = addr
	return newPos, nil
}

func (this *SipGenericParam) parseValueQuotedString(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos
	addr := NewSipQuotedString(context)
	if addr == ABNF_PTR_NIL {
		return newPos, &AbnfError{"generic-param  ParseValue: out of memory for quoted-string value", src, newPos}
	}
	newPos, err = addr.GetSipQuotedString(context).Parse(context, src, newPos)
	if err != nil {
		return newPos, err
	}
	this.valueType = SIP_GENERIC_VALUE_TYPE_QUOTED_STRING
	this.value = addr
	return newPos, nil
}

func (this *SipGenericParam) parseValueIpv6(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos
	p1 := bytes.IndexByte(src[newPos:], ']')
	if p1 == -1 {
		return newPos, &AbnfError{"no ']' for ipv6 for gen-value", src, newPos}
	}

	addr := NewAbnfBuf(context)
	if addr == ABNF_PTR_NIL {
		return newPos, &AbnfError{"generic-param  ParseValue: out of memory for ipv6 value", src, newPos}
	}

	buf := addr.GetAbnfBuf(context)
	buf.SetByteSlice(context, src[newPos:newPos+p1+1])

	newPos += p1 + 1
	this.valueType = SIP_GENERIC_VALUE_TYPE_IPV6
	this.value = addr
	return newPos, nil
}

func (this *SipGenericParam) SetNameAsString(context *ParseContext, name string) {
	this.name.SetString(context, name)
}

func (this *SipGenericParam) SetValueToken(context *ParseContext, value []byte) {
	this.valueType = SIP_GENERIC_VALUE_TYPE_TOKEN
	addr := NewAbnfBuf(context)
	if addr == ABNF_PTR_NIL {
		return
	}

	addr.GetAbnfBuf(context).SetValue(context, value)
	this.value = addr
}

func (this *SipGenericParam) SetValueQuotedString(context *ParseContext, value []byte) {
	this.valueType = SIP_GENERIC_VALUE_TYPE_QUOTED_STRING
	addr := NewSipQuotedString(context)
	if addr == ABNF_PTR_NIL {
		return
	}
	addr.GetSipQuotedString(context).SetValue(context, value)
	this.value = addr
}

func (this *SipGenericParam) GetValueAsByteSlice(context *ParseContext) ([]byte, bool) {
	if this.valueType != SIP_GENERIC_VALUE_TYPE_NOT_EXIST && this.value != ABNF_PTR_NIL {
		if this.valueType == SIP_GENERIC_VALUE_TYPE_TOKEN {
			return this.value.GetAbnfBuf(context).GetAsByteSlice(context), true
		} else if this.valueType == SIP_GENERIC_VALUE_TYPE_QUOTED_STRING {
			return this.value.GetSipQuotedString(context).GetAsByteSlice(context), true
		}
	}
	return nil, false
}

func (this *SipGenericParam) Encode(context *ParseContext, buf *AbnfByteBuffer) {
	//buf.Write(Escape(this.name.GetAsByteSlice(context), IsSipPname))
	//buf.Write(SipPnameEscape(this.name.GetAsByteSlice(context)))
	WriteSipPnameEscape(buf, this.name.GetAsByteSlice(context))
	if this.valueType != SIP_GENERIC_VALUE_TYPE_NOT_EXIST && this.value != ABNF_PTR_NIL {
		if this.valueType == SIP_GENERIC_VALUE_TYPE_TOKEN || this.valueType == SIP_GENERIC_VALUE_TYPE_IPV6 {
			buf.WriteByte('=')
			this.value.GetAbnfBuf(context).Encode(context, buf)
		} else if this.valueType == SIP_GENERIC_VALUE_TYPE_QUOTED_STRING {
			buf.WriteByte('=')
			this.value.GetSipQuotedString(context).Encode(context, buf)
		}
	}
}

func (this *SipGenericParam) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

type SipGenericParams struct {
	AbnfList
}

func NewSipGenericParams(context *ParseContext) AbnfPtr {
	addr := context.allocator.Alloc(uint32(unsafe.Sizeof(SipGenericParams{})))
	if addr == ABNF_PTR_NIL {
		return ABNF_PTR_NIL
	}
	addr.GetSipGenericParams(context).Init()
	return addr
}

func (this *SipGenericParams) Init() {
	this.AbnfList.Init()
}

func (this *SipGenericParams) Size() int32 { return this.Len() }
func (this *SipGenericParams) Empty() bool { return this.Len() == 0 }

func (this *SipGenericParams) GetParam(context *ParseContext, name string) (val *SipGenericParam, ok bool) {
	for e := this.Front(context); e != nil; e = e.Next(context) {
		v := e.Value.GetSipGenericParam(context)
		if v.name.EqualStringNoCase(context, name) {
			return v, true
		}
	}
	return nil, false
}

/* RFC3261
 *
 * generic-param-list   =  *( SWS seperator SWS generic-param )
 *
 * seperator 通常为分号
 *
 */
func (this *SipGenericParams) Parse(context *ParseContext, src []byte, pos int, seperator byte) (newPos int, err error) {
	this.Init()
	return this.ParseWithoutInit(context, src, pos, seperator)
}

func (this *SipGenericParams) ParseWithoutInit(context *ParseContext, src []byte, pos int, seperator byte) (newPos int, err error) {
	newPos = pos
	if newPos >= len(src) {
		return newPos, nil
	}

	for newPos < len(src) {

		if IsOnlyCRLF(src, newPos) {
			return newPos, nil
		}

		if src[newPos] != seperator && !IsLwsChar(src[newPos]) {
			return newPos, nil
		}

		/*
			newPos, err = ParseSWSMark(src, newPos, seperator)
			if err != nil {
				return newPos, nil
			} //*/

		//*
		var macthMark bool
		var newPos1 int
		newPos1, macthMark, err = ParseSWSMarkCanOmmit(src, newPos, seperator)
		if err != nil {
			return newPos, err
		}

		if !macthMark {
			return newPos, nil
		}

		newPos = newPos1 //*/

		addr := NewSipGenericParam(context)
		if addr == ABNF_PTR_NIL {
			return newPos, &AbnfError{"generic-param  parse: out of memory for name", src, newPos}
		}
		newPos, err = addr.GetSipGenericParam(context).ParseWithoutInit(context, src, newPos)
		if err != nil {
			return newPos, err
		}
		this.PushBack(context, addr)
	}

	return newPos, nil
}

func (this *SipGenericParams) Encode(context *ParseContext, buf *AbnfByteBuffer, seperator byte) {
	for e := this.Front(context); e != nil; e = e.Next(context) {
		buf.WriteByte(seperator)
		e.Value.GetSipGenericParam(context).Encode(context, buf)
	}
}

func (this *SipGenericParams) String(context *ParseContext, seperator byte) string {
	var buf AbnfByteBuffer
	this.Encode(context, &buf, seperator)
	return buf.String()
}
