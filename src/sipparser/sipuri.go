package sipparser

import (
	"bytes"
	//"fmt"
	"strings"
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

func NewSipUri() *SipUri {
	//return &SipUri{params: NewSipUriParams(), headers: NewSipUriHeaders()}
	uri := &SipUri{}
	uri.params.Init()
	uri.headers.Init()
	return uri

}

func (this *SipUri) Parse(src []byte, pos int) (newPos int, err error) {

	newPos, err = this.ParseScheme(src, pos)
	if err != nil {
		return newPos, err
	}

	newPos, err = this.ParseUserinfo(src, newPos)
	if err != nil {
		return newPos, err
	}

	newPos, err = this.hostport.Parse(src, newPos)
	if err != nil {
		return newPos, err
	}

	if newPos >= len(src) {
		return newPos, nil
	}

	if src[newPos] == ';' {
		newPos, err = this.params.Parse(src, newPos+1)
		if err != nil {
			return newPos, err
		}
	}

	if newPos >= len(src) {
		return newPos, nil
	}

	if src[newPos] == '?' {
		newPos, err = this.headers.Parse(src, newPos+1)
		if err != nil {
			return newPos, err
		}
	}

	return newPos, err

}

func (this *SipUri) String() string {
	str := ""
	if this.isSecure {
		str = "sips:"
	} else {
		str = "sip:"
	}

	if this.user.Exist() {
		str += string(Escape([]byte(this.user.String()), IsSipUser))
		if this.password.Exist() {
			str += ":"
			str += string(Escape([]byte(this.password.String()), IsSipPassword))
		}
		str += "@"
	}

	str += this.hostport.String()

	if !this.params.Empty() {
		str += ";"
		str += this.params.String()
	}

	if !this.headers.Empty() {
		str += "?"
		str += this.headers.String()
	}

	return str
}

func (this *SipUri) ParseScheme(src []byte, pos int) (newPos int, err error) {
	newPos = pos

	if newPos >= len(src) {
		return newPos, &SipParseError{"parse scheme failed: reach end", src[newPos:]}
	}

	if !IsAlpha(src[newPos]) {
		return newPos, &SipParseError{"parse scheme failed: fisrt char is not alpha", src[newPos:]}
	}

	scheme := &SipToken{}

	newPos, err = scheme.Parse(src, newPos, IsSipScheme)
	if err != nil {
		return newPos, err
	}

	if bytes.Equal(scheme.value, []byte("sips")) {
		this.SetSipsUri()
	} else if !bytes.Equal(scheme.value, []byte("sip")) {
		return newPos, &SipParseError{"parse scheme failed: not sip-uri nor sips-uri", src[newPos:]}
	}

	if newPos >= len(src) {
		return newPos, &SipParseError{"parse scheme failed: no ':' and reach end", src[newPos:]}
	}

	if src[newPos] != ':' {
		return newPos, &SipParseError{"parse scheme failed: no ':'", src[newPos:]}
	}

	newPos++

	return newPos, nil
}

func (this *SipUri) ParseUserinfo(src []byte, pos int) (newPos int, err error) {
	newPos = pos
	hasUserinfo := findUserinfo(src, newPos)
	if hasUserinfo {
		newPos, err = this.user.ParseEscapable(src, newPos, IsSipUser)
		if err != nil {
			return newPos, err
		}

		if this.user.Empty() {
			return newPos, &SipParseError{"parse user-info failed: empty user", src[newPos:]}
		}

		if newPos >= len(src) {
			return newPos, &SipParseError{"parse user-info failed: reach end after user", src[newPos:]}
		}
		this.user.SetExist()

		if src[newPos] == ':' {
			newPos, err = this.password.ParseEscapable(src, newPos+1, IsSipPassword)
			if err != nil {
				return newPos, err
			}
			this.password.SetExist()
		}

		if newPos >= len(src) {
			return newPos, &SipParseError{"parse user-info failed: reach end, and no '@'", src[newPos:]}
		}

		if src[newPos] != '@' {
			return newPos, &SipParseError{"parse user-info failed: no '@'", src[newPos:]}
		}

		newPos++
	}

	return newPos, nil
}

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

type SipUriParam struct {
	name  SipToken
	value SipToken
}

func (this *SipUriParam) Parse(src []byte, pos int) (newPos int, err error) {
	newPos, err = this.name.ParseEscapable(src, pos, IsSipParamChar)
	if err != nil {
		return newPos, err
	}

	if this.name.Empty() {
		return newPos, &SipParseError{"parse sip-uri param failed: empty pname", src[newPos:]}
	}

	this.name.SetExist()

	if src[newPos] == '=' {
		newPos, err = this.value.ParseEscapable(src, newPos+1, IsSipParamChar)
		if err != nil {
			return newPos, err
		}

		if this.value.Empty() {
			return newPos, &SipParseError{"parse sip-uri param failed: empty pvalue", src[newPos:]}
		}
		this.value.SetExist()
	}
	return newPos, nil
}

func (this *SipUriParam) String() string {
	str := this.name.String()
	if this.value.Exist() {
		str += "="
		str += string(Escape([]byte(this.value.String()), IsSipParamChar))
	}
	return str
}

type SipUriParams struct {
	orders []string
	maps   map[string]*SipUriParam
}

func NewSipUriParams() *SipUriParams {
	return &SipUriParams{orders: make([]string, 0), maps: make(map[string]*SipUriParam)}
}

func (this *SipUriParams) Init() {
	this.orders = make([]string, 0)
	this.maps = make(map[string]*SipUriParam)
}

func (this *SipUriParams) String() string {
	if len(this.maps) == 0 {
		return ""
	}

	str := ""
	for i, v := range this.orders {
		if i > 0 {
			str += ";"
		}
		str += this.maps[v].String()
		i++
	}
	return str
}

func (this *SipUriParams) Empty() bool { return len(this.maps) == 0 }
func (this *SipUriParams) GetParam(name string) (val *SipUriParam, ok bool) {
	val, ok = this.maps[strings.ToLower(name)]
	return val, ok
}

func (this *SipUriParams) Parse(src []byte, pos int) (newPos int, err error) {
	newPos = pos
	if newPos >= len(src) {
		return newPos, &SipParseError{"parse uri-param failed: reach end after ';'", src[newPos:]}
	}

	for newPos < len(src) {
		param := &SipUriParam{}
		newPos, err = param.Parse(src, newPos)
		if err != nil {
			return newPos, err
		}
		name := param.name.ToLower()
		this.orders = append(this.orders, name)
		this.maps[name] = param

		if newPos >= len(src) {
			return newPos, nil
		}

		if src[newPos] != ';' {
			return newPos, nil
		}
		newPos++
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
	if this.name.Empty() {
		return newPos, &SipParseError{"parse sip-uri header failed: empty hname", src[newPos:]}
	}
	this.name.SetExist()

	if src[newPos] != '=' {
		return newPos, &SipParseError{"parse header failed: no = after hname", src[newPos:]}
	}

	newPos, err = this.value.ParseEscapable(src, newPos+1, IsSipHeaderChar)
	if err != nil {
		return newPos, err
	}

	this.value.SetExist()

	return newPos, nil
}

func (this *SipUriHeader) String() string {
	str := this.name.String()
	str += "="
	if this.value.Exist() {
		str += string(Escape([]byte(this.value.String()), IsSipHeaderChar))
	}
	return str
}

type SipUriHeaders struct {
	orders []string
	maps   map[string]*SipUriHeader
}

func NewSipUriHeaders() *SipUriHeaders {
	return &SipUriHeaders{orders: make([]string, 0), maps: make(map[string]*SipUriHeader)}
}

func (this *SipUriHeaders) Init() {
	this.orders = make([]string, 0)
	this.maps = make(map[string]*SipUriHeader)
}

func (this *SipUriHeaders) Empty() bool { return len(this.maps) == 0 }
func (this *SipUriHeaders) GetParam(name string) (val *SipUriHeader, ok bool) {
	val, ok = this.maps[strings.ToLower(name)]
	return val, ok
}

func (this *SipUriHeaders) Parse(src []byte, pos int) (newPos int, err error) {
	newPos = pos
	if newPos >= len(src) {
		return newPos, &SipParseError{"parse uri-header failed: reach end after ';'", src[newPos:]}
	}

	for newPos < len(src) {
		header := &SipUriHeader{}
		newPos, err = header.Parse(src, newPos)
		if err != nil {
			return newPos, err
		}
		name := header.name.ToLower()
		this.orders = append(this.orders, name)
		this.maps[name] = header

		if newPos >= len(src) {
			return newPos, nil
		}

		if src[newPos] != '&' {
			return newPos, nil
		}
		newPos++
	}

	return newPos, err
}

func (this *SipUriHeaders) String() string {
	if len(this.maps) == 0 {
		return ""
	}

	str := ""
	for i, v := range this.orders {
		if i > 0 {
			str += "&"
		}
		str += this.maps[v].String()
		i++
	}
	return str
}
