package sipparser

import (
	//"bytes"
	//"fmt"
	//"strings"
	"unsafe"
)

type SipUriParam struct {
	name  AbnfBuf
	value AbnfBuf
}

func NewSipUriParam(context *ParseContext) AbnfPtr {
	addr := context.allocator.Alloc(uint32(unsafe.Sizeof(SipUriParam{})))
	if addr == ABNF_PTR_NIL {
		return ABNF_PTR_NIL
	}
	addr.GetSipUriParam(context).Init()
	return addr
}

func (this *SipUriParam) Init() {
	this.name.Init()
	this.value.Init()
}

func (this *SipUriParam) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	this.Init()
	newPos, err = this.name.ParseEscapable(context, src, pos, IsSipPname)
	if err != nil {
		return newPos, err
	}

	if newPos >= len(src) {
		return newPos, nil
	}

	if src[newPos] == '=' {
		newPos, err = this.value.ParseEscapable(context, src, newPos+1, IsSipPvalue)
		if err != nil {
			return newPos, err
		}
	}
	return newPos, nil
}

func (this *SipUriParam) Encode(context *ParseContext, buf *AbnfByteBuffer) {
	buf.Write(Escape(this.name.GetAsByteSlice(context), IsSipPname))
	if this.value.Exist() {
		buf.WriteByte('=')
		buf.Write(Escape(this.value.GetAsByteSlice(context), IsSipPvalue))
	}
}

func (this *SipUriParam) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

type SipUriParams struct {
	AbnfList
}

func NewSipUriParams(context *ParseContext) AbnfPtr {
	addr := context.allocator.Alloc(uint32(unsafe.Sizeof(SipUriParams{})))
	if addr == ABNF_PTR_NIL {
		return ABNF_PTR_NIL
	}
	addr.GetSipUriParams(context).Init()
	return addr
}

func (this *SipUriParams) Init() {
	this.AbnfList.Init()
}

func (this *SipUriParams) Size() int32 { return this.Len() }
func (this *SipUriParams) Empty() bool { return this.Len() == 0 }
func (this *SipUriParams) GetParam(context *ParseContext, name string) (val *SipUriParam, ok bool) {
	for e := this.Front(context); e != nil; e = e.Next(context) {
		v := e.Value.GetSipUriParam(context)
		if v.name.EqualStringNoCase(context, name) {
			return v, true
		}
	}
	return nil, false
}

func (this *SipUriParams) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos = pos
	if newPos >= len(src) {
		return newPos, &AbnfError{"sip-uri parse: parse sip-uri param failed: reach end after ';'", src, newPos}
	}

	for newPos < len(src) {
		addr := NewSipUriParam(context)
		if addr == ABNF_PTR_NIL {
			return newPos, &AbnfError{"sip-uri  parse: out of memory for sip uri params", src, newPos}
		}
		newPos, err = addr.GetSipUriParam(context).Parse(context, src, newPos)
		if err != nil {
			return newPos, err
		}
		this.PushBack(context, addr)

		if newPos >= len(src) {
			return newPos, nil
		}

		if src[newPos] != ';' {
			return newPos, nil
		}
		newPos++
	}

	return newPos, nil
}

func (this *SipUriParams) EqualRFC3261(context *ParseContext, rhs *SipUriParams) bool {
	params1 := this
	params2 := rhs

	if params1.Size() < params2.Size() {
		params1, params2 = params2, params1
	}

	if !params1.equalSpecParamsRFC3261(context, params2) {
		return false
	}

	for e := params1.Front(context); e != nil; e = e.Next(context) {
		v := e.Value.GetSipUriParam(context)
		param, ok := params2.GetParam(context, v.name.String(context))
		if ok {
			if !param.value.EqualNoCase(context, &v.value) {
				return false
			}
		}
	}

	return true
}

func (this *SipUriParams) equalSpecParamsRFC3261(context *ParseContext, rhs *SipUriParams) bool {
	specParams := []string{"user", "ttl", "method"}

	for _, v := range specParams {
		param1, ok := this.GetParam(context, v)
		if ok {
			param2, ok := rhs.GetParam(context, v)
			if !ok {
				return false
			}
			return param1.value.EqualNoCase(context, &param2.value)
		}
	}

	return true
}

func (this *SipUriParams) Encode(context *ParseContext, buf *AbnfByteBuffer) {
	e := this.Front(context)
	if e != nil {
		e.Value.GetSipUriParam(context).Encode(context, buf)
	}
	e = e.Next(context)

	for ; e != nil; e = e.Next(context) {
		buf.WriteByte(';')
		e.Value.GetSipUriParam(context).Encode(context, buf)
	}
}

func (this *SipUriParams) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}
