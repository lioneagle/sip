package sipparser

import (
	//"bytes"
	_ "fmt"
	//"reflect"
	//"unsafe"
)

type AbnfRef struct {
	Begin int32
	End   int32
}

func (this *AbnfRef) Len() int32 {
	return this.End - this.Begin
}

func (this *AbnfRef) Parse(src []byte, pos int, inCharset AbnfIsInCharset) (end int) {
	this.Begin = int32(pos)
	len1 := len(src)

	for end = pos; end < len1; end++ {
		if !inCharset(src[end]) {
			break
		}
	}

	this.End = int32(end)
	return end
}

func (this *AbnfRef) ParseSipToken(src []byte, pos int) (end int) {
	this.Begin = int32(pos)
	len1 := len(src)

	if pos >= len1 || !IsSipToken(src[pos]) {
		this.End = int32(pos)
		return pos
	}

	for end = pos + 1; end < len1; end++ {
		if !IsSipToken(src[end]) {
			this.End = int32(end)
			return end
		}
	}

	this.End = int32(end)
	return end
}

func (this *AbnfRef) ParseSipWord(src []byte, pos int) (end int) {
	this.Begin = int32(pos)
	len1 := len(src)

	if pos >= len1 || !IsSipWord(src[pos]) {
		this.End = int32(pos)
		return pos
	}

	for end = pos + 1; end < len1; end++ {
		if !IsSipWord(src[end]) {
			break
		}
	}

	this.End = int32(end)
	return end
}

func (this *AbnfRef) ParseUriScheme(src []byte, pos int) (end int) {
	this.Begin = int32(pos)
	len1 := len(src)

	if pos >= len1 || !IsUriScheme(src[pos]) {
		this.End = int32(pos)
		return pos
	}

	for end = pos + 1; end < len1; end++ {
		if !IsUriScheme(src[end]) {
			break
		}
	}

	this.End = int32(end)
	return end
}

func (this *AbnfRef) ParseWspChar(src []byte, pos int) (end int) {
	this.Begin = int32(pos)
	len1 := len(src)

	if pos >= len1 || !IsWspChar(src[pos]) {
		this.End = int32(pos)
		return pos
	}

	for end = pos; end < len1; end++ {
		if !IsWspChar(src[end]) {
			break
		}
	}

	this.End = int32(end)
	return end
}

func (this *AbnfRef) ParseDigit(src []byte, pos int) (end int) {
	this.Begin = int32(pos)
	len1 := len(src)

	if pos >= len1 || !IsDigit(src[pos]) {
		this.End = int32(pos)
		return pos
	}

	for end = pos; end < len1; end++ {
		if !IsDigit(src[end]) {
			break
		}
	}

	this.End = int32(end)
	return end
}

func (this *AbnfRef) ParseEscapable0(src []byte, pos int, inCharset AbnfIsInCharset) (escapeNum, newPos int, err error) {
	this.Begin = int32(pos)
	len1 := len(src)

	for newPos = pos; newPos < len1; newPos++ {
		if src[newPos] == '%' {
			if (newPos + 2) >= len1 {
				return escapeNum, newPos, &AbnfError{"AbnfRef ParseEscapable: reach end after %", src, newPos}
			}
			if !IsHex(src[newPos+1]) || !IsHex(src[newPos+2]) {
				return escapeNum, newPos, &AbnfError{"AbnfRef ParseEscapable: no hex after %", src, newPos}
			}
			escapeNum++
			newPos += 2
		} else if !inCharset(src[newPos]) {
			break
		}
	}
	this.End = int32(newPos)
	return escapeNum, newPos, nil
}

func (this *AbnfRef) ParseEscapable(src []byte, pos int, inCharset AbnfIsInCharset) (escapeNum, newPos int, err error) {
	this.Begin = int32(pos)
	len1 := len(src)

	for newPos = pos; newPos < len1; newPos++ {
		if src[newPos] == '%' {
			break
		} else if !inCharset(src[newPos]) {
			this.End = int32(newPos)
			return escapeNum, newPos, nil
		}
	}

	for ; newPos < len1; newPos++ {
		if src[newPos] == '%' {
			if (newPos + 2) >= len1 {
				return escapeNum, newPos, &AbnfError{"AbnfRef ParseEscapable: reach end after %", src, newPos}
			}
			if !IsHex(src[newPos+1]) || !IsHex(src[newPos+2]) {
				return escapeNum, newPos, &AbnfError{"AbnfRef ParseEscapable: no hex after %", src, newPos}
			}
			escapeNum++
			newPos += 2
		} else if !inCharset(src[newPos]) {
			this.End = int32(newPos)
			return escapeNum, newPos, nil
		}
	}
	this.End = int32(newPos)
	return escapeNum, newPos, nil
}

func (this *AbnfRef) ParseEscapableSipUser0(src []byte, pos int) (escapeNum, newPos int, err error) {
	this.Begin = int32(pos)
	len1 := len(src)

	for newPos = pos; newPos < len1; newPos++ {
		if src[newPos] == '%' {
			if (newPos + 2) >= len1 {
				return escapeNum, newPos, &AbnfError{"AbnfRef ParseEscapableSipUser: reach end after %", src, newPos}
			}
			if !IsHex(src[newPos+1]) || !IsHex(src[newPos+2]) {
				return escapeNum, newPos, &AbnfError{"AbnfRef ParseEscapableSipUser: no hex after %", src, newPos}
			}
			escapeNum++
			newPos += 2
		} else if !IsSipUser(src[newPos]) {
			break
		}
	}
	this.End = int32(newPos)
	return escapeNum, newPos, nil
}

func (this *AbnfRef) ParseEscapableSipUser(src []byte, pos int) (escapeNum, newPos int, err error) {
	this.Begin = int32(pos)
	len1 := len(src)

	for newPos = pos; newPos < len1; newPos++ {
		if src[newPos] == '%' {
			break
		} else if !IsSipUser(src[newPos]) {
			this.End = int32(newPos)
			return escapeNum, newPos, nil
		}
	}

	for ; newPos < len1; newPos++ {
		if src[newPos] == '%' {
			if (newPos + 2) >= len1 {
				return escapeNum, newPos, &AbnfError{"AbnfRef ParseEscapableSipUser: reach end after %", src, newPos}
			}
			if !IsHex(src[newPos+1]) || !IsHex(src[newPos+2]) {
				return escapeNum, newPos, &AbnfError{"AbnfRef ParseEscapableSipUser: no hex after %", src, newPos}
			}
			escapeNum++
			newPos += 2
		} else if !IsSipUser(src[newPos]) {
			this.End = int32(newPos)
			return escapeNum, newPos, nil
		}
	}
	this.End = int32(newPos)
	return escapeNum, newPos, nil
}
