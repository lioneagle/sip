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
	//SIP_GENERIC_VALUE_TYPE_HOST          = 3
)

type SipGenericParam struct {
	name      AbnfBuf
	valueType int32
	value     AbnfPtr
	//parsed    AbnfPtr
}

func NewSipGenericParam(context *ParseContext) (*SipGenericParam, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipGenericParam{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*SipGenericParam)(unsafe.Pointer(mem)).Init()
	return (*SipGenericParam)(unsafe.Pointer(mem)), addr
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
	newPos, err = this.name.ParseEscapable(context, src, pos, IsSipPname)
	if err != nil {
		return newPos, err
	}

	if newPos >= len(src) {
		return newPos, nil
	}

	newPos, err = ParseSWSMark(src, newPos, '=')
	if err != nil {
		return newPos, err
	}

	if newPos >= len(src) {
		return newPos, &AbnfError{"generic-param parse: parse value failed: empty gen-value", src, newPos}
	}

	/* @@TODO: 目前解gen-value时，暂不考虑解析出host，因为一般没有必要解析出来，以后再考虑添加这个功能 */

	if IsSipToken(src[newPos]) {
		token, addr := NewAbnfBuf(context)
		if token == nil {
			return newPos, &AbnfError{"generic-param  parse: out of memory for token value", src, newPos}
		}
		newPos, err = token.Parse(context, src, newPos, IsSipToken)
		if err != nil {
			return newPos, err
		}
		this.valueType = SIP_GENERIC_VALUE_TYPE_TOKEN
		this.value = addr
	} else if (src[newPos] == '"') || IsLwsChar(src[newPos]) {
		quotedString, addr := NewSipQuotedString(context)
		if quotedString == nil {
			return newPos, &AbnfError{"generic-param  parse: out of memory for quoted-string value", src, newPos}
		}
		newPos, err = quotedString.Parse(context, src, newPos)
		if err != nil {
			return newPos, err
		}
		this.valueType = SIP_GENERIC_VALUE_TYPE_QUOTED_STRING
		this.value = addr
	} else {
		return newPos, &AbnfError{"generic-param parse: parse value failed: not token nor quoted-string", src, newPos}
	}

	return newPos, nil
}

func (this *SipGenericParam) SetNameAsString(context *ParseContext, name string) {
	this.name.SetString(context, name)
}

func (this *SipGenericParam) SetValueQuotedString(context *ParseContext, value []byte) {
	this.valueType = SIP_GENERIC_VALUE_TYPE_QUOTED_STRING

	quotedString, addr := NewSipQuotedString(context)
	if quotedString == nil {
		return
	}

	quotedString.SetValue(context, value)
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

func (this *SipGenericParam) Encode(context *ParseContext, buf *bytes.Buffer) {
	buf.Write(Escape(this.name.GetAsByteSlice(context), IsSipPname))
	if this.valueType != SIP_GENERIC_VALUE_TYPE_NOT_EXIST && this.value != ABNF_PTR_NIL {
		if this.valueType == SIP_GENERIC_VALUE_TYPE_TOKEN {
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

func NewSipGenericParams(context *ParseContext) (*SipGenericParams, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(SipGenericParams{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*SipGenericParams)(unsafe.Pointer(mem)).Init()
	return (*SipGenericParams)(unsafe.Pointer(mem)), addr
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
	newPos = pos
	if newPos >= len(src) {
		return newPos, nil
	}

	for newPos < len(src) {

		newPos, err = ParseSWSMark(src, newPos, seperator)
		if err != nil {
			return newPos, nil
		}

		param, addr := NewSipGenericParam(context)
		if param == nil {
			return newPos, &AbnfError{"generic-param  parse: out of memory for name", src, newPos}
		}
		newPos, err = param.Parse(context, src, newPos)
		if err != nil {
			return newPos, err
		}
		this.PushBack(context, addr)
	}

	return newPos, nil
}

func (this *SipGenericParams) Encode(context *ParseContext, buf *bytes.Buffer, seperator byte) {
	for e := this.Front(context); e != nil; e = e.Next(context) {
		buf.WriteByte(seperator)
		e.Value.GetSipGenericParam(context).Encode(context, buf)
	}
}

func (this *SipGenericParams) String(context *ParseContext, seperator byte) string {
	var buf bytes.Buffer
	this.Encode(context, &buf, seperator)
	return buf.String()
}
