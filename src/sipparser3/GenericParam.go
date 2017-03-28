package sipparser3

import (
	"bytes"
	//"fmt"
	//"strings"
)

const (
	SIP_GENERIC_VALUE_TYPE_NOT_EXIST     = 0
	SIP_GENERIC_VALUE_TYPE_TOKEN         = 1
	SIP_GENERIC_VALUE_TYPE_QUOTED_STRING = 2
	//SIP_GENERIC_VALUE_TYPE_HOST          = 3
)

type SipGenericParam struct {
	name      AbnfToken
	valueType int
	value     interface{}
	parsed    interface{}
}

func NewSipGenericParam() *SipGenericParam {
	return &SipGenericParam{}
}

/*
 * generic-param  =  token [ EQUAL gen-value ]
 * gen-value      =  token / host / quoted-string
 *
 */
func (this *SipGenericParam) Parse(context *ParseContext, src []byte, pos int) (newPos int, err error) {
	newPos, err = this.name.ParseEscapable(context, src, pos, IsSipPname)
	if err != nil {
		return newPos, err
	}

	if newPos >= len(src) {
		return newPos, nil
	}

	newPos, err = ParseSWSMark(src, newPos, '=')
	if err != nil {
		return newPos, err
	}

	if newPos >= len(src) {
		return newPos, &AbnfError{"generic-param parse: parse value failed: empty gen-value", src, newPos}
	}

	/* @@TODO: 目前解gen-value时，暂不考虑解析出host，因为一般没有必要解析出来，以后再考虑添加这个功能 */

	if IsSipToken(src[newPos]) {
		var token AbnfToken
		newPos, err = token.Parse(context, src, newPos, IsSipToken)
		if err != nil {
			return newPos, err
		}
		this.valueType = SIP_GENERIC_VALUE_TYPE_TOKEN
		this.value = token
	} else if (src[newPos] == '"') || IsLwsChar(src[newPos]) {
		var quotedString SipQuotedString
		newPos, err = quotedString.Parse(context, src, newPos)
		if err != nil {
			return newPos, err
		}
		this.valueType = SIP_GENERIC_VALUE_TYPE_QUOTED_STRING
		this.value = quotedString
	} else {
		return newPos, &AbnfError{"generic-param parse: parse value failed: not token nor quoted-string", src, newPos}
	}

	return newPos, nil
}

func (this *SipGenericParam) Encode(buf *bytes.Buffer) {
	buf.Write(Escape(this.name.value, IsSipPname))
	if this.valueType == SIP_GENERIC_VALUE_TYPE_TOKEN {
		buf.WriteByte('=')
		token, ok := this.value.(AbnfToken)
		if ok {
			token.Encode(buf)
		}
	} else if this.valueType == SIP_GENERIC_VALUE_TYPE_QUOTED_STRING {
		buf.WriteByte('=')
		quotedString, ok := this.value.(SipQuotedString)
		if ok {
			quotedString.Encode(buf)
		}
	}
}

func (this *SipGenericParam) String() string {
	return AbnfEncoderToString(this)
}

type SipGenericParams struct {
	params []SipGenericParam
}

func NewSipGenericParams() *SipGenericParams {
	ret := &SipGenericParams{}
	ret.Init()
	return ret
}

func (this *SipGenericParams) Init() {
	this.params = make([]SipGenericParam, 0, 2)
}

func (this *SipGenericParams) Size() int   { return len(this.params) }
func (this *SipGenericParams) Empty() bool { return len(this.params) == 0 }
func (this *SipGenericParams) GetParam(name string) (val *SipGenericParam, ok bool) {
	for i, v := range this.params {
		if v.name.EqualStringNoCase(name) {
			return &this.params[i], true
		}
	}
	return nil, false
}

/* RFC3261
 *
 * generic-param-list   =  *( SWS seperator SWS generic-param )
 *
 * seperator 通常为分号
 *
 */
func (this *SipGenericParams) Parse(context *ParseContext, src []byte, pos int, seperator byte) (newPos int, err error) {
	newPos = pos
	if newPos >= len(src) {
		return newPos, nil
	}

	for newPos < len(src) {

		newPos, err = ParseSWSMark(src, newPos, seperator)
		if err != nil {
			return newPos, nil
		}

		param := SipGenericParam{}
		newPos, err = param.Parse(context, src, newPos)
		if err != nil {
			return newPos, err
		}
		this.params = append(this.params, param)
	}

	return newPos, nil
}

func (this *SipGenericParams) Encode(buf *bytes.Buffer, seperator byte) {
	for _, v := range this.params {
		buf.WriteByte(seperator)
		v.Encode(buf)
	}
}

func (this *SipGenericParams) String(seperator byte) string {
	var buf bytes.Buffer
	this.Encode(&buf, seperator)
	return buf.String()
}
