package sipparser3

import (
	"bytes"
	//"fmt"
	"strings"
)

type AbnfToken struct {
	exist bool
	value []byte
}

func (this *AbnfToken) String() string {
	if this.exist {
		return Bytes2str(this.value)
	}
	return ""
}

func (this *AbnfToken) Encode(buf *bytes.Buffer) {
	if this.exist {
		buf.Write(this.value)
	}
}

func (this *AbnfToken) Exist() bool  { return this.exist }
func (this *AbnfToken) Size() int    { return len(this.value) }
func (this *AbnfToken) Empty() bool  { return len(this.value) == 0 }
func (this *AbnfToken) SetExist()    { this.exist = true }
func (this *AbnfToken) SetNonExist() { this.exist = false }

func (this *AbnfToken) SetValue(value []byte) { this.value = value }
func (this *AbnfToken) ToLower() string       { return strings.ToLower(string(this.value)) }
func (this *AbnfToken) ToUpper() string       { return strings.ToUpper(string(this.value)) }

func (this *AbnfToken) Equal(rhs *AbnfToken) bool {
	if (this.exist && !rhs.exist) || (!this.exist && rhs.exist) {
		return false
	}
	return bytes.Equal(this.value, rhs.value)
}

func (this *AbnfToken) EqualNoCase(rhs *AbnfToken) bool {
	if (this.exist && !rhs.exist) || (!this.exist && rhs.exist) {
		return false
	}
	return EqualNoCase(this.value, rhs.value)
}

func (this *AbnfToken) EqualBytes(rhs []byte) bool {
	if !this.exist {
		return false
	}
	return bytes.Equal(this.value, rhs)
}

func (this *AbnfToken) EqualBytesNoCase(rhs []byte) bool {
	if !this.exist {
		return false
	}
	return EqualNoCase(this.value, rhs)
}

func (this *AbnfToken) EqualString(str string) bool {
	if !this.exist {
		return false
	}
	return bytes.Equal(this.value, Str2bytes(str))
}

func (this *AbnfToken) EqualStringNoCase(str string) bool {
	if !this.exist {
		return false
	}

	return EqualNoCase(this.value, Str2bytes(str))
}

func (this *AbnfToken) Parse(context *ParseContext, src []byte, pos int, inCharset AbnfIsInCharset) (newPos int, err error) {
	begin, end, newPos, err := parseToken(src, pos, inCharset)
	if err != nil {
		return newPos, err
	}
	if begin >= end {
		return newPos, &AbnfError{"AbnfToken parse: value is empty", src, newPos}
	}
	this.SetExist()

	this.value = src[begin:end]
	return newPos, nil
}

func (this *AbnfToken) ParseEscapable(context *ParseContext, src []byte, pos int, inCharset AbnfIsInCharset) (newPos int, err error) {
	begin, end, newPos, err := parseTokenEscapable(src, pos, inCharset)
	if err != nil {
		return newPos, err
	}

	if begin >= end {
		return newPos, &AbnfError{"AbnfToken ParseEscapable: value is empty", src, newPos}
	}

	this.SetExist()

	this.value = Unescape(context, src[begin:end])
	//this.value = src[begin:end]
	return newPos, nil
}
