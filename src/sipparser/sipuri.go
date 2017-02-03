package sipparser

import (
	"bytes"
	//"fmt"
)

type SipUri struct {
	isSecure bool
	user     SipToken
	password SipToken
	hostport SipHostPort
	params   SipUriParams
	headers  SipUriHeaders
}

func (this *SipUri) SetSipUri()      { this.isSecure = false }
func (this *SipUri) SetSipsUri()     { this.isSecure = true }
func (this *SipUri) IsSipUri() bool  { return !this.isSecure }
func (this *SipUri) IsSipsUri() bool { return this.isSecure }

func findUserinfo(src []byte, pos int) bool {
	for newPos := pos; newPos < len(src); newPos++ {
		if src[newPos] == '@' {
			return true
		} else if src[newPos] == '>' || IsLwsChar(src[newPos]) {
			return false
		}
	}
	return false
}

func (this *SipUri) Parse(src []byte, pos int) (newPos int, err error) {
	newPos = pos

	scheme := &SipToken{}
	newPos, err = scheme.ParseEscapable(src, newPos, IsSipScheme)
	if err != nil {
		return newPos, err
	}

	if bytes.Equal(scheme.value, []byte("sips")) {
		this.SetSipsUri()
	} else if !bytes.Equal(scheme.value, []byte("sip")) {
		return newPos, &SipParseError{"parse scheme failed: not sip-uri nor sips-uri", src[newPos:]}
	}

	if src[newPos] != ':' {
		return newPos, &SipParseError{"parse scheme failed: not ':'", src[newPos:]}
	}

	newPos++

	hasUserinfo := findUserinfo(src, newPos)
	if hasUserinfo {
		newPos, err = this.user.ParseEscapable(src, newPos, IsSipUser)
		if err != nil {
			return newPos, err
		}

		if src[newPos] == ':' {
			newPos, err = this.password.ParseEscapable(src, newPos+1, IsSipPassword)
			if err != nil {
				return newPos, err
			}
		}

		if newPos >= len(src) {
			return newPos, &SipParseError{"parse user-info failed: reach end, and no '@'", src[newPos:]}
		}

		if src[newPos] != '@' {
			return newPos, &SipParseError{"parse user-info failed: no '@'", src[newPos:]}
		}

		newPos++
	}

	newPos, err = this.hostport.Parse(src, newPos)
	if err != nil {
		return newPos, err
	}

	return newPos, err

}

type SipUriParam struct {
	name  SipToken
	value SipToken
}

func (this *SipUriParam) Parse(src []byte, pos int) (newPos int, err error) {
	newPos, err = this.name.ParseEscapable(src, pos, IsSipParamChar)
	if err != nil {
		return newPos, err
	}

	if src[newPos] == '=' {
		newPos, err = this.value.ParseEscapable(src, newPos, IsSipParamChar)
		if err != nil {
			return newPos, err
		}
	}
	return newPos, nil
}

type SipUriParams struct {
	maps map[string]*SipUriParam
}

func (this *SipUriParams) Empty() bool { return len(this.maps) == 0 }
func (this *SipUriParams) GetParam(name string) (val *SipUriParam, ok bool) {
	val, ok = this.maps[name]
	return val, ok
}

func (this *SipUriParams) Parse(src []byte, pos int) (newPos int, err error) {
	for newPos < len(src) {
		if src[newPos] != ';' {
			return newPos, nil
		}

		param := &SipUriParam{}
		newPos, err = param.Parse(src, pos)
		if err != nil {
			return newPos, err
		}
		this.maps[string(param.name.value)] = param
	}

	return newPos, err
}

type SipUriHeader struct {
	name  SipToken
	value SipToken
}

func (this *SipUriHeader) Parse(src []byte, pos int) (newPos int, err error) {
	newPos, err = this.name.ParseEscapable(src, pos, IsSipHeaderChar)
	if err != nil {
		return newPos, err
	}

	if src[newPos] != '=' {
		return newPos, &SipParseError{"parse header failed: no = after hname", src[newPos:]}
	}

	newPos, err = this.value.ParseEscapable(src, newPos, IsSipHeaderChar)
	if err != nil {
		return newPos, err
	}

	return newPos, nil
}

type SipUriHeaders struct {
	maps map[string]*SipUriHeader
}

func (this *SipUriHeaders) Empty() bool { return len(this.maps) == 0 }
func (this *SipUriHeaders) GetParam(name string) (val *SipUriHeader, ok bool) {
	val, ok = this.maps[name]
	return val, ok
}

func (this *SipUriHeaders) Parse(src []byte, pos int) (newPos int, err error) {
	for newPos < len(src) {
		if src[newPos] != ';' {
			return newPos, nil
		}

		header := &SipUriHeader{}
		newPos, err = header.Parse(src, pos)
		if err != nil {
			return newPos, err
		}
		this.maps[string(header.name.value)] = header
	}

	return newPos, err
}
