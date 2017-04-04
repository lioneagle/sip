package sipparser

import (
	"bytes"
	//"fmt"
	"unsafe"
)

type TelUriContext struct {
	exist        bool
	isDomainName bool
	desc         AbnfBuf
}

func NewTelUriContext(context *ParseContext) (*TelUriContext, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(TelUriContext{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*TelUriContext)(unsafe.Pointer(mem)).Init()
	return (*TelUriContext)(unsafe.Pointer(mem)), addr
}

func (this *TelUriContext) Init() {
	this.exist = false
	this.isDomainName = false
	this.desc.Init()
}

func (this *TelUriContext) Equal(context *ParseContext, rhs *TelUriContext) bool {
	if (this.exist && !rhs.exist) || (!this.exist && rhs.exist) {
		return false
	}

	if this.exist {
		return this.desc.EqualNoCase(context, &rhs.desc)
	}

	return true
}

func (this *TelUriContext) Exist() bool  { return this.exist }
func (this *TelUriContext) SetExist()    { this.exist = true }
func (this *TelUriContext) SetNonExist() { this.exist = false }

func (this *TelUriContext) Encode(context *ParseContext, buf *bytes.Buffer) {
	if this.exist {
		buf.WriteString(";phone-context=")
		buf.Write(Escape(this.desc.GetAsByteSlice(context), IsTelPvalue))
	}
}

func (this *TelUriContext) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

type TelUriParam struct {
	name  AbnfBuf
	value AbnfBuf
}

func NewTelUriParam(context *ParseContext) (*TelUriParam, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(TelUriParam{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*TelUriParam)(unsafe.Pointer(mem)).Init()
	return (*TelUriParam)(unsafe.Pointer(mem)), addr
}

func (this *TelUriParam) Init() {
	this.name.Init()
	this.value.Init()
}

func (this *TelUriParam) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos, err = this.name.ParseEscapable(context, src, pos, IsTelPname)
	if err != nil {
		return newPos, err
	}

	if newPos >= len(src) {
		return newPos, nil
	}

	if src[newPos] == '=' {
		newPos, err = this.value.ParseEscapable(context, src, newPos+1, IsTelPvalue)
		if err != nil {
			return newPos, err
		}
	}
	return newPos, nil
}

func (this *TelUriParam) Encode(context *ParseContext, buf *bytes.Buffer) {
	buf.Write(Escape(this.name.GetAsByteSlice(context), IsTelPname))
	if this.value.Exist() {
		buf.WriteByte('=')
		buf.Write(Escape(this.value.GetAsByteSlice(context), IsTelPvalue))
	}
}

func (this *TelUriParam) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}

type TelUriParams struct {
	AbnfList
}

func NewTelUriParams(context *ParseContext) (*TelUriParams, AbnfPtr) {
	mem, addr := context.allocator.Alloc(int32(unsafe.Sizeof(TelUriParams{})))
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*TelUriParams)(unsafe.Pointer(mem)).Init()
	return (*TelUriParams)(unsafe.Pointer(mem)), addr
}

func (this *TelUriParams) Init() {
	this.AbnfList.Init()
}

func (this *TelUriParams) Size() int32 { return this.Len() }
func (this *TelUriParams) Empty() bool { return this.Len() == 0 }

func (this *TelUriParams) GetParam(context *ParseContext, name string) (val *TelUriParam, ok bool) {
	for e := this.Front(context); e != nil; e = e.Next(context) {
		v := e.Value.GetTelUriParam(context)
		if v.name.EqualStringNoCase(context, name) {
			return v, true
		}
	}

	return nil, false
}

func (this *TelUriParams) Equal(context *ParseContext, rhs *TelUriParams) bool {
	if this.Size() != rhs.Size() {
		return false
	}

	/* RFC3966
	 * Parameters are compared according to 'pname', regardless of the
	 * order they appeared in the URI.  If one URI has a parameter name
	 * not found in the other, the two URIs are not equal.
	 */
	for e := this.Front(context); e != nil; e = e.Next(context) {
		v := e.Value.GetTelUriParam(context)
		param, ok := rhs.GetParam(context, v.name.String(context))
		if !ok {
			return false
		}

		if !param.value.EqualNoCase(context, &v.value) {
			return false
		}
	}

	return true
}

func (this *TelUriParams) Encode(context *ParseContext, buf *bytes.Buffer) {
	for e := this.Front(context); e != nil; e = e.Next(context) {
		buf.WriteByte(';')
		e.Value.GetTelUriParam(context).Encode(context, buf)
	}
}

func (this *TelUriParams) String(context *ParseContext) string {
	return AbnfEncoderToString(context, this)
}
