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

func HasPrefixByteSliceNoCase(s1, s2 []byte) bool {
	if len(s1) < len(s2) {
		return false
	}

	if len(s2) <= 0 {
		return false
	}
	return EqualNoCase(s1[:len(s2)], s2)
}

func unescapeToByte(src []byte) byte {
	return HexToByte(src[1])<<4 | HexToByte(src[2])
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
	if pos >= len(src) || !IsDigit(src[pos]) {
		return 0, 0, pos, false
	}

	num = 0
	digit = 0
	newPos = pos

	for newPos < len(src) && IsDigit(src[newPos]) {
		digit = digit*10 + int(src[newPos]) - '0'
		newPos++
		num++
	}

	return digit, num, newPos, true

}

func ParseUriScheme(context *ParseContext, src []byte, pos int) (newPos int, scheme *AbnfBuf, err error) {
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

	scheme = &AbnfBuf{}

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

	return newPos, scheme, nil
}

func ParseHcolon(src []byte, pos int) (newPos int, err error) {
	/* RFC3261 Section 25.1, page 220
	 *
	 * HCOLON  =  *( SP / HTAB ) ":" SWS
	 */
	ref := AbnfRef{}
	newPos = ref.Parse(src, pos, IsWspChar)

	if newPos >= len(src) {
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

	newPos1, err := ParseLWS(src, newPos)
	if err == nil {
		newPos = newPos1
	}
	return newPos, nil
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

func ParseCRLF(src []byte, pos int) (newPos int, err error) {
	if ((pos + 1) >= len(src)) || (src[pos] != '\r') || (src[pos+1] != '\n') {
		return pos, &AbnfError{"CRLF parse: wrong CRLF", src, pos}
	}
	return pos + 2, nil
}

func EncodeUInt(buf *bytes.Buffer, digit uint64) {
	buf.WriteString(strconv.FormatUint(uint64(digit), 10))
}
