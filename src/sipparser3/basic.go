package sipparser3

import (
	"bytes"
	//"container/list"
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

type AbnfError struct {
	description string
	src         []byte
	pos         int
}

func Str2bytes(s string) []byte {
	//h := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(&s)), Len: len(s), Cap: len(s)}
	//x := (*[2]uintptr)(unsafe.Pointer(&s))
	//h := [3]uintptr{x[0], x[1], x[1]}
	x := (*reflect.StringHeader)(unsafe.Pointer(&s))
	h := reflect.SliceHeader{Data: x.Data, Len: x.Len, Cap: x.Len}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func (err *AbnfError) Error() string {
	if err.pos < len(err.src) {
		return fmt.Sprintf("%s at src[%d]: %s", err.description, err.pos, string(err.src[err.pos:]))
	}
	return fmt.Sprintf("%s at end")
}

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

func (this *AbnfToken) Parse(context *ParseContext, src []byte, pos int, isInCharset func(ch byte) bool) (newPos int, err error) {
	begin, end, newPos, err := parseToken(context, src, pos, isInCharset)
	if err != nil {
		return newPos, err
	}

	this.value = src[begin:end]
	return newPos, nil
}

func (this *AbnfToken) ParseEscapable(context *ParseContext, src []byte, pos int, isInCharset func(ch byte) bool) (newPos int, err error) {
	begin, end, newPos, err := parseTokenEscapable(context, src, pos, isInCharset)
	if err != nil {
		return newPos, err
	}

	this.value = Unescape(context, src[begin:end])
	//this.value = src[begin:end]
	return newPos, nil
}

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

func Unescape(context *ParseContext, src []byte) (dst []byte) {
	if bytes.IndexByte(src, '%') == -1 {
		return src
	}

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

func NeedEscape(src []byte, isInCharset func(ch byte) bool) bool {
	for _, v := range src {
		if !isInCharset(v) {
			return true
		}
	}
	return false
}

func Escape(src []byte, isInCharset func(ch byte) bool) (dst []byte) {
	if !NeedEscape(src, isInCharset) {
		return src
	}

	for _, v := range src {
		if isInCharset(v) {
			dst = append(dst, v)
		} else {
			dst = append(dst, '%', ToUpperHex(v>>4), ToUpperHex(v))
		}
	}

	return dst
}

func parseToken(context *ParseContext, src []byte, pos int, isInCharset func(ch byte) bool) (tokenBegin, tokenEnd, newPos int, err error) {
	tokenBegin = pos
	for newPos = pos; newPos < len(src); newPos++ {
		if !isInCharset(src[newPos]) {
			break
		}
	}
	return tokenBegin, newPos, newPos, nil
}

func parseTokenEscapable(context *ParseContext, src []byte, pos int, isInCharset func(ch byte) bool) (tokenBegin, tokenEnd, newPos int, err error) {
	tokenBegin = pos
	for newPos = pos; newPos < len(src); newPos++ {
		if src[newPos] == '%' {
			if (newPos + 2) >= len(src) {
				return tokenBegin, newPos, newPos, &AbnfError{"token parse: parse escape token failed: reach end", src, newPos}
			}
			if !IsHex(src[newPos+1]) || !IsHex(src[newPos+2]) {
				return tokenBegin, newPos, newPos, &AbnfError{"token parse: parse escape token failed: no hex after %", src, newPos}
			}
			newPos += 2
		} else if !isInCharset(src[newPos]) {
			break
		}
	}
	return tokenBegin, newPos, newPos, nil
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

func ParseUriScheme(context *ParseContext, src []byte, pos int) (newPos int, scheme *AbnfToken, err error) {
	/* RFC3261 Section 25.1, page 223
	 *
	 * scheme         =  ALPHA *( ALPHA / DIGIT / "+" / "-" / "." )
	 */
	newPos = pos

	if newPos >= len(src) {
		return newPos, nil, &AbnfError{"uri-scheme parse: parse scheme failed: reach end", src, newPos}
	}

	if !IsAlpha(src[newPos]) {
		return newPos, nil, &AbnfError{"uri-scheme parse: parse scheme failed: fisrt char is not alpha", src, newPos}
	}

	scheme = &AbnfToken{}

	newPos, err = scheme.Parse(context, src, newPos, IsUriScheme)
	if err != nil {
		return newPos, nil, err
	}

	if newPos >= len(src) {
		return newPos, nil, &AbnfError{"uri-scheme parse: parse scheme failed: no ':' and reach end", src, newPos}
	}

	if src[newPos] != ':' {
		return newPos, nil, &AbnfError{"uri-scheme parse: parse scheme failed: no ':'", src, newPos}
	}

	newPos++
	scheme.SetExist()

	return newPos, scheme, nil
}

func ParseSWSMark(src []byte, pos int, mark byte) (newPos int, err error) {
	/* RFC3261 Section 25.1, page 220
	 *
	 * STAR    =  SWS "*" SWS ; asterisk
	 * SLASH   =  SWS "/" SWS ; slash
	 * EQUAL   =  SWS "=" SWS ; equal
	 * LPAREN  =  SWS "(" SWS ; left parenthesis
	 * RPAREN  =  SWS ")" SWS ; right parenthesis
	 * COMMA   =  SWS "," SWS ; comma
	 * SEMI    =  SWS ";" SWS ; semicolon
	 * COLON   =  SWS ":" SWS ; colon
	 */
	newPos = pos
	newPos, err = ParseSWS(src, newPos)
	if err != nil {
		return newPos, nil
	}

	if newPos >= len(src) {
		return newPos, &AbnfError{"SWSMark parse: reach end before mark", src, newPos}
	}

	if src[newPos] != mark {
		return newPos, &AbnfError{"SWSMark parse: not expected mark after SWS", src, newPos}
	}

	return ParseSWS(src, newPos)
}

func ParseSWS(src []byte, pos int) (newPos int, err error) {
	/* RFC3261 Section 25.1, page 220
	 *
	 * SWS  =  [LWS] ; sep whitespace
	 */
	newPos = pos
	if newPos >= len(src) {
		return newPos, nil
	}

	if !IsLwsChar(src[newPos]) {
		return newPos, nil
	}

	return ParseLWS(src, newPos)
}

func ParseLWS(src []byte, pos int) (newPos int, err error) {
	/* RFC3261 Section 25.1, page 220
	 *
	 * LWS  =  [*WSP CRLF] 1*WSP ; linear whitespace
	 * WSP  =  ( SP | HTAB )
	 *
	 * NOTE:
	 *
	 * 1. this defination of LWS is different from that in RFC2616 (HTTP/1.1)
	 *    RFC2616 Section 2.2, page 16:
	 *
	 *    LWS  = [CRLF] 1*( SP | HTAB )
	 *
	 * 2. WSP's defination is from RFC2234 Section 6.1, page 12
	 *
	 */
	for newPos = pos; newPos < len(src); newPos++ {
		if !IsWspChar(src[newPos]) {
			break
		}
	}

	if newPos >= len(src) {
		return newPos, nil
	}

	if (newPos+1 < len(src)) && (src[newPos] == '\r') && (src[newPos+1] == '\n') {
		newPos += 2

		if newPos >= len(src) {
			return newPos, &AbnfError{"LWS parse: no char after CRLF in LWS", src, newPos}
		}

		if !IsWspChar(src[newPos]) {
			return newPos, &AbnfError{"LWS parse: no WSP after CRLF in LWS", src, newPos}
		}

		for ; newPos < len(src); newPos++ {
			if !IsWspChar(src[newPos]) {
				break
			}
		}
	}

	return newPos, nil
}

func ParseLeftAngleQuote(src []byte, pos int) (newPos int, err error) {
	/* RFC3261 Section 25.1, page 221
	 *
	 * LAQUOT  =  SWS "<"; left angle quote
	 *
	 */
	newPos = pos

	if newPos >= len(src) {
		return newPos, &AbnfError{"LAQUOT parse: reach end at begining", src, newPos}
	}

	newPos, err = ParseSWS(src, newPos)
	if err != nil {
		return newPos, err
	}

	if newPos >= len(src) {
		return newPos, &AbnfError{"LAQUOT parse: reach end before <", src, newPos}
	}

	if src[newPos] != '<' {
		return newPos, &AbnfError{"LAQUOT parse: no <", src, newPos}
	}

	return newPos + 1, nil
}

func ParseRightAngleQuote(src []byte, pos int) (newPos int, err error) {
	/* RFC3261 Section 25.1, page 221
	 *
	 * RAQUOT  =  ">" SWS ; right angle quote
	 *
	 */
	newPos = pos
	if newPos >= len(src) {
		return newPos, &AbnfError{"RAQUOT parse: reach end at begining", src, newPos}
	}

	if src[newPos] != '>' {
		return newPos, &AbnfError{"RAQUOT parse: no >", src, newPos}
	}

	return ParseSWS(src, newPos+1)
}
