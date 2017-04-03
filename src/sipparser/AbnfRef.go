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

func (this *AbnfRef) Parse(context *ParseContext, src []byte, pos int, inCharset AbnfIsInCharset) (newPos int, err error) {
	begin, end, newPos, err := parseToken(src, pos, inCharset)
	if err != nil {
		return newPos, err
	}
	if begin >= end {
		return newPos, &AbnfError{"AbnfRef parse: value is empty", src, newPos}
	}
	this.Begin = int32(begin)
	this.End = int32(end)
	return newPos, nil
}
