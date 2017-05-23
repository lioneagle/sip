package sipparser

import (
//"bytes"
//"fmt"
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
	//num := len(src)

	//*
	for end = pos; end < len(src); end++ {
		if !inCharset(src[end]) {
			break
		}
	}

	this.End = int32(end)
	return end
	//*/

	/*
		var v byte
		for end, v = range src[pos:] {
			if !inCharset(v) {
				this.End = int32(pos + end)
				return pos + end
			}
		}

		this.End = int32(len(src))
		return len(src)
		//*/
}

func (this *AbnfRef) ParseEscapable(src []byte, pos int, inCharset AbnfIsInCharset) (escapeNum, newPos int, err error) {
	this.Begin = int32(pos)

	for newPos = pos; newPos < len(src); newPos++ {
		if src[newPos] == '%' {
			if (newPos + 2) >= len(src) {
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
