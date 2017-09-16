package sipparser

import (
	"bytes"
	"strconv"
)

func ToUpperHex(ch byte) byte {
	return "0123456789ABCDEF"[ch&0x0F]
}

func ToLowerHex(ch byte) byte {
	return "0123456789abcdef"[ch&0x0F]
}

func ToLower(ch byte) byte {
	if IsUpper(ch) {
		//return ch - 'A' + 'a'
		return ch | 0x20
	}
	return ch
}

func ToUpper(ch byte) byte {
	if IsLower(ch) {
		//return ch - 'a' + 'A'
		return ch & 0xDF
	}
	return ch
}

func HexToByte(ch byte) byte {
	if IsDigit(ch) {
		//return ch - '0'
		return ch & 0x0F
	}
	if IsLowerHexAlpha(ch) {
		return ch - 'a' + 10
	}
	return ch - 'A' + 10
}

func CompareNoCase(s1, s2 []byte) int {
	len1 := len(s1)
	if len1 != len(s2) {
		return len1 - len(s2)
	}

	for i := 0; i < len1; i++ {
		if s1[i] != s2[i] {
			ch1 := ToLower(s1[i])
			ch2 := ToLower(s2[i])
			if ch1 != ch2 {
				return int(ch1) - int(ch2)
			}
		}
	}

	return 0
}

var Count1 int = 0
var Count2 int = 0
var Count3 int = 0
var Count4 int = 0

func EqualNoCase0(s1, s2 []byte) bool {
	len1 := len(s1)
	if len1 != len(s2) {
		return false
	}

	if ToLower(s1[0]) != ToLower(s2[0]) {
		return false
	}

	for i := 1; i < len1; i++ {
		if s1[i] != s2[i] {
			if ToLower(s1[i]) != ToLower(s2[i]) {
				return false
			}
		}
	}

	return true
}

func EqualNoCase(s1, s2 []byte) bool {
	len1 := len(s1)
	if len1 != len(s2) {
		return false
	}

	if ToLower(s1[0]) != ToLower(s2[0]) {
		return false
	}

	i := 1

	for ; i < len1; i++ {
		if s1[i] != s2[i] {
			break
		}
	}

	for ; i < len1; i++ {
		if ToLower(s1[i]) != ToLower(s2[i]) {
			return false
		}
	}

	return true
}

func Unescape(src []byte) (dst []byte) {
	if bytes.IndexByte(src, '%') == -1 {
		return src
	}

	len1 := len(src)

	for i := 0; i < len1; {
		if (src[i] == '%') && ((i + 2) < len1) && IsHex(src[i+1]) && IsHex(src[i+2]) {
			dst = append(dst, unescapeToByte(src[i:]))
			i += 3
		} else {
			dst = append(dst, src[i])
			i++
		}
	}

	return dst
}

func HasPrefixByteSliceNoCase(s1, s2 []byte) bool {
	len2 := len(s2)
	if len(s1) < len2 {
		return false
	}

	if len2 <= 0 {
		return false
	}
	return EqualNoCase(s1[:len2], s2)
}

func unescapeToByte(src []byte) byte {
	return HexToByte(src[1])<<4 | HexToByte(src[2])
}

func unescapeToByte2(x1, x2 byte) byte {
	return HexToByte(x1)<<4 | HexToByte(x2)
}

func NeedEscape(src []byte, inCharset AbnfIsInCharset) bool {
	for _, v := range src {
		if !inCharset(v) {
			return true
		}
	}
	return false
}

func Escape(src []byte, inCharset AbnfIsInCharset) (dst []byte) {
	if !NeedEscape(src, inCharset) {
		return src
	}

	for _, v := range src {
		if inCharset(v) {
			dst = append(dst, v)
		} else {
			dst = append(dst, '%', ToUpperHex(v>>4), ToUpperHex(v))
		}
	}

	return dst
}

func ParseUInt(src []byte, pos int) (digit, num, newPos int, ok bool) {
	len1 := len(src)
	if pos >= len1 || !IsDigit(src[pos]) {
		return 0, 0, pos, false
	}

	num = 0
	digit = 0
	newPos = pos

	for newPos < len1 && IsDigit(src[newPos]) {
		digit = digit*10 + int(src[newPos]) - '0'
		newPos++
		num++
	}

	return digit, num, newPos, true

}

func ParseUriScheme(context *ParseContext, src []byte, pos int, scheme *AbnfBuf) (newPos int, err error) {
	/* RFC3261 Section 25.1, page 223
	 *
	 * scheme         =  ALPHA *( ALPHA / DIGIT / "+" / "-" / "." )
	 */
	newPos = pos

	if newPos >= len(src) {
		return newPos, &AbnfError{"uri-scheme parse: parse scheme failed: reach end", src, newPos}
	}

	if !IsAlpha(src[newPos]) {
		return newPos, &AbnfError{"uri-scheme parse: parse scheme failed: fisrt char is not alpha", src, newPos}
	}

	newPos, err = scheme.ParseUriScheme(context, src, newPos)
	if err != nil {
		return newPos, err
	}

	if newPos >= len(src) {
		return newPos, &AbnfError{"uri-scheme parse: parse scheme failed: no ':' and reach end", src, newPos}
	}

	if src[newPos] != ':' {
		return newPos, &AbnfError{"uri-scheme parse: parse scheme failed: no ':'", src, newPos}
	}

	newPos++

	return newPos, nil
}

func ParseHcolon(src []byte, pos int) (newPos int, err error) {
	/* RFC3261 Section 25.1, page 220
	 *
	 * HCOLON  =  *( SP / HTAB ) ":" SWS
	 */
	//ref := AbnfRef{}
	//newPos = ref.ParseWspChar(src, pos)
	//newPos = (&AbnfRef{}).ParseWspChar(src, pos)

	len1 := len(src)

	for newPos = pos; newPos < len1; newPos++ {
		if !IsWspChar(src[newPos]) {
			break
		}
	}

	if newPos >= len1 {
		return newPos, &AbnfError{"HCOLON parse: reach end before ':'", src, newPos}
	}

	if src[newPos] != ':' {
		return newPos, &AbnfError{"HCOLON parse: no ':' after *( SP / HTAB )", src, newPos}
	}

	return ParseSWS(src, newPos+1)
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

	return ParseSWS(src, newPos+1)
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

	newPos1, ok := ParseLWS(src, newPos)
	if ok {
		newPos = newPos1
	}
	return newPos, nil
}

func ParseLWS(src []byte, pos int) (newPos int, ok bool) {
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
	//newPos = eatWsp(src, pos)
	for newPos = pos; newPos < len(src); newPos++ {
		if !IsWspChar(src[newPos]) {
			break
		}
	}

	if (newPos + 1) >= len(src) {
		return newPos, true
	}

	//if IsCRLF2(src, newPos) {
	if (src[newPos] == '\r') && (src[newPos+1] == '\n') {
		newPos += 2

		if newPos >= len(src) {
			//return newPos, &AbnfError{"LWS parse: no char after CRLF in LWS", src, newPos}
			return newPos, false
		}

		if !IsWspChar(src[newPos]) {
			//return newPos, &AbnfError{"LWS parse: no WSP after CRLF in LWS", src, newPos}
			return newPos, false
		}

		//newPos = eatWsp(src, newPos)
		for ; newPos < len(src); newPos++ {
			if !IsWspChar(src[newPos]) {
				break
			}
		}
	}

	return newPos, true
}

func eatWsp(src []byte, pos int) (newPos int) {
	for newPos = pos; newPos < len(src); newPos++ {
		if !IsWspChar(src[newPos]) {
			break
		}
	}
	return newPos
}

func ParseLeftAngleQuote(src []byte, pos int) (newPos int, err error) {
	/* RFC3261 Section 25.1, page 221
	 *
	 * LAQUOT  =  SWS "<"; left angle quote
	 *
	 */
	newPos = pos
	len1 := len(src)

	if newPos >= len1 {
		return newPos, &AbnfError{"LAQUOT parse: reach end at begining", src, newPos}
	}

	newPos, err = ParseSWS(src, newPos)
	if err != nil {
		return newPos, err
	}

	if newPos >= len1 {
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

func IsCRLF(src []byte, pos int) bool {
	return ((pos + 1) < len(src)) && (src[pos] == '\r') && (src[pos+1] == '\n')
}

func IsCRLF2(src []byte, pos int) bool {
	return (src[pos] == '\r') && (src[pos+1] == '\n')
}

func IsOnlyCRLF(src []byte, pos int) bool {
	if (pos + 2) < len(src) {
		return (src[pos] == '\r') && (src[pos+1] == '\n') && !IsWspChar(src[pos+2])
	}
	return (pos+2) == len(src) && (src[pos] == '\r') && (src[pos+1] == '\n')
}

func ParseCRLF(src []byte, pos int) (newPos int, err error) {
	if !IsCRLF(src, pos) {
		return pos, &AbnfError{"CRLF parse: wrong CRLF", src, pos}
	}
	return pos + 2, nil
}

func EncodeUInt(buf *bytes.Buffer, digit uint64) {
	buf.WriteString(strconv.FormatUint(uint64(digit), 10))
}
