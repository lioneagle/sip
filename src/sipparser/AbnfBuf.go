package sipparser

import (
	"bytes"
	//"fmt"
	"reflect"
	"unsafe"
)

type AbnfBuf struct {
	// using highest bit of size as exist flag
	addr  AbnfPtr
	size  uint32
	alloc uint32
}

const (
	ABNF_BUF_EXIST_BIT  = uint32(0x80000000)
	ABNF_BUF_EXIST_MASK = uint32(0x7fffffff)
)

func SizeofAbnfBuf() int32 {
	return int32(unsafe.Sizeof(AbnfBuf{}))
}

func NewAbnfBuf(context *ParseContext) (*AbnfBuf, AbnfPtr) {
	mem, addr := context.allocator.Alloc(SizeofAbnfBuf())
	if mem == nil {
		return nil, ABNF_PTR_NIL
	}

	(*AbnfBuf)(unsafe.Pointer(mem)).Init()
	//addr.GetAbnfBuf(context).Init()
	return (*AbnfBuf)(unsafe.Pointer(mem)), addr
}

func (this *AbnfBuf) Init() {
	this.addr = ABNF_PTR_NIL
	this.size = 0
	this.alloc = 0
}

func (this *AbnfBuf) Empty() bool {
	return this.Size() == 0
}

func (this *AbnfBuf) SetEmpty() {
	this.size = ABNF_BUF_EXIST_BIT
}

func (this *AbnfBuf) Size() int32 {
	return int32(this.size & ABNF_BUF_EXIST_MASK)
}

func (this *AbnfBuf) Exist() bool {
	return (this.size & ABNF_BUF_EXIST_BIT) != 0
}

func (this *AbnfBuf) SetExist() {
	this.size |= ABNF_BUF_EXIST_BIT
}

func (this *AbnfBuf) SetNonExist() {
	this.size = 0
}

func (this *AbnfBuf) setSize(size int32) {
	this.size = uint32(size) | ABNF_BUF_EXIST_BIT
}

func (this *AbnfBuf) SetValue(context *ParseContext, value []byte) {
	this.SetByteSlice(context, value)
}

func (this *AbnfBuf) allocMem(context *ParseContext, size int32) bool {
	if size == 0 {
		this.SetNonExist()
		return true
	}

	//if this.Size() < size {
	if int32(this.alloc) < size {
		addr, alloc := context.allocator.AllocEx(size)
		if addr == ABNF_PTR_NIL {
			// keep unchanged
			return false
		}
		this.addr = addr
		this.alloc = uint32(alloc)
	}

	this.setSize(size)
	return true
}

func (this *AbnfBuf) SetByteSlice(context *ParseContext, buf []byte) {
	size := int32(len(buf))
	if size == 0 {
		this.SetNonExist()
		return
	}

	/*
		if this.allocMem(context, int32(len(buf))) {
			//this.SetExist()
			copy(this.GetAsByteSlice(context), buf)
		}

		// remain unchanged
		return
		//*/

	//*
	if int32(this.alloc) < size {
		addr, alloc := context.allocator.AllocEx(size)
		if addr == ABNF_PTR_NIL {
			// keep unchanged
			return
		}
		this.addr = addr
		this.alloc = uint32(alloc)
	}
	this.setSize(size)

	copy(this.GetAsByteSlice(context), buf)
	//*/

}

func (this *AbnfBuf) SetByteSliceWithUnescape(context *ParseContext, buf []byte, escapeNum int) {
	if escapeNum <= 0 {
		this.SetByteSlice(context, buf)
		return
	}
	if !this.allocMem(context, int32(len(buf)-2*escapeNum)) {
		this.SetByteSlice(context, buf)
		return
	}
	this.SetExist()
	dst := this.GetAsByteSlice(context)
	var j int32 = 0
	for i := 0; i < len(buf); {
		if (buf[i] == '%') && ((i + 2) < len(buf)) && IsHex(buf[i+1]) && IsHex(buf[i+2]) {
			dst[j] = unescapeToByte(buf[i:])
			i += 3
		} else {
			dst[j] = buf[i]
			i++
		}
		j++
	}
}

func (this *AbnfBuf) SetString(context *ParseContext, str string) {
	this.SetByteSlice(context, StringToByteSlice(str))
}

func (this *AbnfBuf) GetAsByteSlice(context *ParseContext) []byte {
	if this.addr == ABNF_PTR_NIL {
		return nil
	}
	size := int(this.Size())
	header := reflect.SliceHeader{Data: this.addr.GetUintptr(context), Len: size, Cap: size}
	return *(*[]byte)(unsafe.Pointer(&header))
}

func (this *AbnfBuf) GetAsString(context *ParseContext) string {
	if this.addr == ABNF_PTR_NIL {
		return ""
	}
	header := reflect.StringHeader{Data: this.addr.GetUintptr(context), Len: int(this.Size())}
	return *(*string)(unsafe.Pointer(&header))
}

//*
func (this *AbnfBuf) ParseEnableEmpty(context *ParseContext, src []byte, pos int, inCharset AbnfIsInCharset) (newPos int) {
	ref := AbnfRef{}
	newPos = ref.Parse(src, pos, inCharset)
	if ref.Begin < ref.End {
		this.SetByteSlice(context, src[ref.Begin:ref.End])
	} else {
		this.SetEmpty()
	}
	return newPos
}

func (this *AbnfBuf) Parse(context *ParseContext, src []byte, pos int, inCharset AbnfIsInCharset) (newPos int, err error) {
	ref := AbnfRef{}
	newPos = ref.Parse(src, pos, inCharset)

	if ref.Begin >= ref.End {
		return newPos, &AbnfError{"AbnfBuf Parse: value is empty", src, newPos}
	}
	this.SetByteSlice(context, src[ref.Begin:ref.End])
	return newPos, nil
}

func (this *AbnfBuf) ParseSipToken(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	ref := AbnfRef{}
	newPos = ref.ParseSipToken(src, pos)

	if ref.Begin >= ref.End {
		return newPos, &AbnfError{"AbnfBuf ParseSipToken: value is empty", src, newPos}
	}
	this.SetByteSlice(context, src[ref.Begin:ref.End])
	return newPos, nil
}

func (this *AbnfBuf) ParseSipWord(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	ref := AbnfRef{}
	newPos = ref.ParseSipWord(src, pos)

	if ref.Begin >= ref.End {
		return newPos, &AbnfError{"AbnfBuf ParseSipWord: value is empty", src, newPos}
	}
	this.SetByteSlice(context, src[ref.Begin:ref.End])
	return newPos, nil
}

func (this *AbnfBuf) ParseUriScheme(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	ref := AbnfRef{}
	newPos = ref.ParseUriScheme(src, pos)

	if ref.Begin >= ref.End {
		return newPos, &AbnfError{"AbnfBuf ParseUriScheme: value is empty", src, newPos}
	}
	this.SetByteSlice(context, src[ref.Begin:ref.End])
	return newPos, nil
}

func (this *AbnfBuf) ParseWspChar(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	ref := AbnfRef{}
	newPos = ref.ParseWspChar(src, pos)

	if ref.Begin >= ref.End {
		return newPos, &AbnfError{"AbnfBuf ParseWspChar: value is empty", src, newPos}
	}
	this.SetByteSlice(context, src[ref.Begin:ref.End])
	return newPos, nil
}

func (this *AbnfBuf) ParseDigit(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	ref := AbnfRef{}
	newPos = ref.ParseDigit(src, pos)

	if ref.Begin >= ref.End {
		return newPos, &AbnfError{"AbnfBuf ParseDigit: value is empty", src, newPos}
	}
	this.SetByteSlice(context, src[ref.Begin:ref.End])
	return newPos, nil
}

func (this *AbnfBuf) ParseEscapableEnableEmpty(context *ParseContext, src []byte, pos int, inCharset AbnfIsInCharset) (newPos int, err error) {
	ref := AbnfRef{}
	escapeNum, newPos, err := ref.ParseEscapable(src, pos, inCharset)
	if err != nil {
		return newPos, err
	}

	if ref.Begin >= ref.End {
		this.SetEmpty()
		return newPos, nil
	}

	if escapeNum == 0 {
		this.SetByteSlice(context, src[ref.Begin:ref.End])

	} else {
		this.SetByteSliceWithUnescape(context, src[ref.Begin:ref.End], escapeNum)
	}
	return newPos, nil
}

func (this *AbnfBuf) ParseEscapable(context *ParseContext, src []byte, pos int, inCharset AbnfIsInCharset) (newPos int, err error) {
	ref := AbnfRef{}
	escapeNum, newPos, err := ref.ParseEscapable(src, pos, inCharset)
	if err != nil {
		return newPos, err
	}

	if ref.Begin >= ref.End {
		return newPos, &AbnfError{"AbnfBuf ParseEscapable: value is empty", src, newPos}
	}

	if escapeNum == 0 {
		this.SetByteSlice(context, src[ref.Begin:ref.End])

	} else {
		this.SetByteSliceWithUnescape(context, src[ref.Begin:ref.End], escapeNum)
	}
	return newPos, nil
}

func (this *AbnfBuf) ParseEscapableSipUser(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	ref := AbnfRef{}
	escapeNum, newPos, err := ref.ParseEscapableSipUser(src, pos)
	if err != nil {
		return newPos, err
	}

	if ref.Begin >= ref.End {
		return newPos, &AbnfError{"AbnfBuf ParseEscapableSipUser: value is empty", src, newPos}
	}

	if escapeNum == 0 {
		this.SetByteSlice(context, src[ref.Begin:ref.End])

	} else {
		this.SetByteSliceWithUnescape(context, src[ref.Begin:ref.End], escapeNum)
	}
	return newPos, nil
}

func (this *AbnfBuf) simpleNotEqual(rhs *AbnfBuf) bool {
	return (this.size != rhs.size) || (this.addr == ABNF_PTR_NIL) || (rhs.addr == ABNF_PTR_NIL)
}

func (this *AbnfBuf) simpleHaveNoPrefix(prefix []byte) bool {
	if !this.Exist() || this.Empty() || this.addr == ABNF_PTR_NIL || len(prefix) == 0 {
		return true
	}

	return this.Size() < int32(len(prefix))
}

func (this *AbnfBuf) HasPrefixByteSlice(context *ParseContext, prefix []byte) bool {
	if this.simpleHaveNoPrefix(prefix) {
		return false
	}

	return bytes.HasPrefix(this.GetAsByteSlice(context), prefix)
}

func (this *AbnfBuf) HasPrefixByteSliceNoCase(context *ParseContext, prefix []byte) bool {
	if this.simpleHaveNoPrefix(prefix) {
		return false
	}

	return EqualNoCase(this.GetAsByteSlice(context)[:len(prefix)], prefix)
}

func (this *AbnfBuf) Equal(context *ParseContext, rhs *AbnfBuf) bool {
	if this.addr == rhs.addr {
		return true
	}
	if this.simpleNotEqual(rhs) {
		return false
	}

	if this.Empty() {
		return true
	}
	return bytes.Equal(this.GetAsByteSlice(context), rhs.GetAsByteSlice(context))
}

func (this *AbnfBuf) EqualNoCase(context *ParseContext, rhs *AbnfBuf) bool {
	if this.addr == rhs.addr {
		return true
	}

	if this.simpleNotEqual(rhs) {
		return false
	}

	if this.Empty() {
		return true
	}

	return EqualNoCase(this.GetAsByteSlice(context), rhs.GetAsByteSlice(context))
}

func (this *AbnfBuf) EqualByteSlice(context *ParseContext, rhs []byte) bool {
	if this.addr == ABNF_PTR_NIL || this.Size() != int32(len(rhs)) {
		return false
	}
	return bytes.Equal(this.GetAsByteSlice(context), rhs)
}

func (this *AbnfBuf) EqualByteSliceNoCase(context *ParseContext, rhs []byte) bool {
	if this.addr == ABNF_PTR_NIL || this.Size() != int32(len(rhs)) {
		return false
	}
	return EqualNoCase(this.GetAsByteSlice(context), rhs)
}

func (this *AbnfBuf) EqualString(context *ParseContext, str string) bool {
	return this.EqualByteSlice(context, StringToByteSlice(str))
}

func (this *AbnfBuf) EqualStringNoCase(context *ParseContext, str string) bool {
	return this.EqualByteSliceNoCase(context, StringToByteSlice(str))
}

func (this *AbnfBuf) Encode(context *ParseContext, buf *bytes.Buffer) {
	if this.addr != ABNF_PTR_NIL && this.Size() != 0 {
		buf.Write(this.GetAsByteSlice(context))
	}
}

func (this *AbnfBuf) String(context *ParseContext) string {
	if this.addr != ABNF_PTR_NIL && this.Size() != 0 {
		return AbnfEncoderToString(context, this)
	}
	return ""
}
