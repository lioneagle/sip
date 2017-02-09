package sipparser

import (
	//"container/list"
	"bytes"
	"fmt"
	"strings"
)

type AbnfError struct {
	description string
	src         []byte
	pos         int
}

func (err *AbnfError) Error() string {
	if err.pos < len(err.src) {
		return fmt.Sprintf("%s at %s", err.description, string(err.src[err.pos:]))
	}
	return fmt.Sprintf("%s at end")
}

type AbnfToken struct {
	exist bool
	value []byte
}

func (this *AbnfToken) String() string { return string(this.value) }
func (this *AbnfToken) Exist() bool    { return this.exist }
func (this *AbnfToken) Size() int      { return len(this.value) }
func (this *AbnfToken) Empty() bool    { return len(this.value) == 0 }
func (this *AbnfToken) SetExist()      { this.exist = true }
func (this *AbnfToken) SetNonExist()   { this.exist = false }

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

func (this *AbnfToken) EqualString(str string) bool {
	if !this.exist {
		return false
	}
	return bytes.Equal(this.value, []byte(str))
}

func (this *AbnfToken) EqualStringNoCase(str string) bool {
	if !this.exist {
		return false
	}

	return EqualNoCase(this.value, []byte(str))
}

func (this *AbnfToken) Parse(src []byte, pos int, isInCharset func(ch byte) bool) (newPos int, err error) {
	begin, end, newPos, err := parseToken(src, pos, isInCharset)
	if err != nil {
		return newPos, err
	}

	this.value = src[begin:end]
	return newPos, nil
}

func (this *AbnfToken) ParseEscapable(src []byte, pos int, isInCharset func(ch byte) bool) (newPos int, err error) {
	begin, end, newPos, err := parseTokenEscapable(src, pos, isInCharset)
	if err != nil {
		return newPos, err
	}

	this.value = Unescape(src[begin:end])
	return newPos, nil
}

/*
type SipList struct {
	list.List
}

func (this *SipList) RemoveAll() {
	var n *list.Element
	for e := this.Front(); e != nil; e = n {
		n = e.Next()
		this.Remove(e)
	}
}*/

func ToUpperHex(ch byte) byte {
	return "0123456789ABCDEF"[ch&0x0F]
}

func ToLowerHex(ch byte) byte {
	return "0123456789abcdef"[ch&0x0F]
}

func ToLower(ch byte) byte {
	if IsUpper(ch) {
		return ch - 'A' + 'a'
	}
	return ch
}

func ToUpper(ch byte) byte {
	if IsLower(ch) {
		return ch - 'a' + 'A'
	}
	return ch
}

func HexToByte(ch byte) byte {
	if IsDigit(ch) {
		return ch - '0'
	}
	if IsLowerHexAlpha(ch) {
		return ch - 'a' + 10
	}
	return ch - 'A' + 10
}

func CompareNoCase(s1, s2 []byte) int {
	if len(s1) != len(s2) {
		return len(s1) - len(s2)
	}

	for i, v := range s1 {
		if v != s2[i] {
			ch1 := ToLower(v)
			ch2 := ToLower(s2[i])
			if ch1 != ch2 {
				return int(ch1) - int(ch2)
			}
		}
	}

	return 0
}

func EqualNoCase(s1, s2 []byte) bool {
	return CompareNoCase(s1, s2) == 0
}

func Unescape(src []byte) (dst []byte) {
	for i := 0; i < len(src); {
		if (src[i] == '%') && ((i + 2) < len(src)) && IsHex(src[i+1]) && IsHex(src[i+2]) {
			dst = append(dst, unescapeToByte(src[i:]))
			i += 3
		} else {
			dst = append(dst, src[i])
			i++
		}
	}

	return dst
}

func unescapeToByte(src []byte) byte {
	return HexToByte(src[1])<<4 | HexToByte(src[2])
}

func Escape(src []byte, isInCharset func(ch byte) bool) (dst []byte) {
	for _, v := range src {
		if isInCharset(v) {
			dst = append(dst, v)
		} else {
			dst = append(dst, '%', ToUpperHex(v>>4), ToUpperHex(v))
		}
	}

	return dst
}

func parseToken(src []byte, pos int, isInCharset func(ch byte) bool) (tokenBegin, tokenEnd, newPos int, err error) {
	tokenBegin = pos
	for newPos = pos; newPos < len(src); newPos++ {
		if !isInCharset(src[newPos]) {
			break
		}
	}
	tokenEnd = newPos
	return tokenBegin, tokenEnd, newPos, nil
}

func parseTokenEscapable(src []byte, pos int, isInCharset func(ch byte) bool) (tokenBegin, tokenEnd, newPos int, err error) {
	tokenBegin = pos
	for newPos = pos; newPos < len(src); newPos++ {
		if src[newPos] == '%' {
			if (newPos + 2) >= len(src) {
				return tokenBegin, newPos, newPos, &AbnfError{"parse escape token failed: reach end", src, newPos}
			}
			if !IsHex(src[newPos+1]) || !IsHex(src[newPos+2]) {
				return tokenBegin, newPos, newPos, &AbnfError{"parse escape token failed: no hex after %", src, newPos}
			}
			newPos += 2
		} else if !isInCharset(src[newPos]) {
			break
		}
	}
	tokenEnd = newPos
	return tokenBegin, tokenEnd, newPos, nil
}

func ParseUInt(src []byte, pos int) (digit, newPos int, ok bool) {
	if pos >= len(src) || !IsDigit(src[pos]) {
		return 0, pos, false
	}

	num := 0
	digit = 0
	newPos = pos

	for newPos < len(src) && IsDigit(src[newPos]) {
		digit = digit*10 + int(src[newPos]) - '0'
		newPos++
		num++
	}

	return digit, newPos, true

}

func ParseUriScheme(src []byte, pos int) (newPos int, scheme *AbnfToken, err error) {
	newPos = pos

	if newPos >= len(src) {
		return newPos, nil, &AbnfError{"parse scheme failed: reach end", src, newPos}
	}

	if !IsAlpha(src[newPos]) {
		return newPos, nil, &AbnfError{"parse scheme failed: fisrt char is not alpha", src, newPos}
	}

	scheme = &AbnfToken{}

	newPos, err = scheme.Parse(src, newPos, IsUriScheme)
	if err != nil {
		return newPos, nil, err
	}

	if newPos >= len(src) {
		return newPos, nil, &AbnfError{"parse scheme failed: no ':' and reach end", src, newPos}
	}

	if src[newPos] != ':' {
		return newPos, nil, &AbnfError{"parse scheme failed: no ':'", src, newPos}
	}

	newPos++
	scheme.SetExist()

	return newPos, scheme, nil
}
