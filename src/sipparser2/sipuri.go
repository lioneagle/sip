package sipparser2

import (
	//"fmt"
	"bytes"
	//"container/list"
	"strings"
)

type SipUri struct {
	isSecure bool
	user     AbnfToken
	password AbnfToken
	hostport SipHostPort
	params   SipUriParams
	headers  SipUriHeaders
}

func (this *SipUri) SetSipUri()      { this.isSecure = false }
func (this *SipUri) SetSipsUri()     { this.isSecure = true }
func (this *SipUri) IsSipUri() bool  { return !this.isSecure }
func (this *SipUri) IsSipsUri() bool { return this.isSecure }

func NewSipUri() *SipUri {
	uri := &SipUri{}
	uri.params.Init(g_allocator)
	//uri.params.Init()
	uri.headers.Init()
	return uri

}

func (this *SipUri) Scheme() string {
	if this.isSecure {
		return "sips"
	}
	return "sip"
}

func (this *SipUri) Parse(src []byte, pos int) (newPos int, err error) {
	newPos, err = this.ParseScheme(src, pos)
	if err != nil {
		return newPos, err
	}

	return this.ParseAfterScheme(src, newPos)
}

func (this *SipUri) ParseAfterScheme(src []byte, pos int) (newPos int, err error) {
	newPos = pos

	newPos, err = this.ParseUserinfo(src, newPos)
	if err != nil {
		return newPos, err
	}
	//return newPos, nil

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

func (this *SipUri) Encode(buf *bytes.Buffer) {
	buf.WriteString(this.Scheme())
	buf.WriteByte(':')

	if this.user.Exist() {
		buf.Write(Escape(this.user.value, IsSipUser))
		if this.password.Exist() {
			buf.WriteByte(':')
			buf.Write(Escape(this.password.value, IsSipPassword))
		}
		buf.WriteByte('@')
	}

	this.hostport.Encode(buf)

	if !this.params.Empty() {
		buf.WriteByte(';')
		this.params.Encode(buf)
	}

	if !this.headers.Empty() {
		buf.WriteByte('?')
		this.headers.Encode(buf)
	}
}

func (this *SipUri) String() string {
	/*var buf bytes.Buffer
	this.Encode(&buf)
	return buf.String()*/
	//*
	str := this.Scheme()
	str += ":"

	if this.user.Exist() {
		str += Bytes2str(Escape([]byte(this.user.String()), IsSipUser))
		if this.password.Exist() {
			str += ":"
			str += Bytes2str(Escape([]byte(this.password.String()), IsSipPassword))
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

	return str //*/
}

func (this *SipUri) Equal(uri URI) bool {
	rhs, ok := uri.(*SipUri)
	if !ok {
		return false
	}

	if (this.isSecure && !rhs.isSecure) || (!this.isSecure && rhs.isSecure) {
		return false
	}

	if !this.EqualUserinfo(rhs) {
		return false
	}

	if !this.hostport.Equal(&rhs.hostport) {
		return false
	}

	if !this.EqualParams(rhs) {
		return false
	}

	if !this.EqualHeaders(rhs) {
		return false
	}

	return true
}

func (this *SipUri) EqualUserinfo(rhs *SipUri) bool {
	if (this.user.Exist() && !rhs.user.Exist()) || (!this.user.Exist() && rhs.user.Exist()) {
		return false
	}
	ret := this.user.Equal(&rhs.user) && this.password.Equal(&rhs.password)
	return ret
}

func (this *SipUri) EqualParams(rhs *SipUri) bool {
	uri1 := this
	uri2 := rhs

	if uri1.params.Size() < uri2.params.Size() {
		uri1, uri2 = uri2, uri1
	}

	if !uri1.equalSpecParams(uri2) {
		return false
	}

	/*for _, v := range uri1.params.maps {
		param, ok := uri2.params.GetParam(v.name.String())
		if ok {
			if !param.value.EqualNoCase(&v.value) {
				return false
			}
		}
	}*/
	return true
}

func (this *SipUri) EqualHeaders(rhs *SipUri) bool {
	if this.headers.Size() != rhs.headers.Size() {
		return false
	}

	for _, v := range this.headers.maps {
		header, ok := rhs.headers.GetHeader(v.name.String())
		if ok {
			if !header.value.EqualNoCase(&v.value) {
				return false
			}
		}
	}
	return true
}

func (this *SipUri) equalSpecParams(rhs *SipUri) bool {
	specParams := []string{"user", "ttl", "method"}

	for _, v := range specParams {
		_, ok := this.params.GetParam(v)
		if ok {
			_, ok = rhs.params.GetParam(v)
			if !ok {
				return false
			}
		}
	}

	return true
}

func (this *SipUri) ParseScheme(src []byte, pos int) (newPos int, err error) {
	if len(src) >= 4 && (src[0] == 's') && (src[1] == 'i') && (src[2] == 'p') && (src[3] == ':') {
		this.SetSipUri()
		return 4, nil
	}

	if len(src) >= 4 && (src[0] == 's') && (src[1] == 'i') && (src[2] == 'p') && (src[3] == 's') && (src[4] == ':') {
		this.SetSipsUri()
		return 5, nil
	}

	return 0, &AbnfError{"sip-uri parse: parse scheme failed: not sip-uri nor sips-uri", src, 0}

	/*
		newPos, scheme, err := ParseUriScheme(src, pos)
		if err != nil {
			return newPos, err
		}

		if EqualNoCase(scheme.value, Str2bytes("sips")) {
			this.SetSipsUri()
		} else if !EqualNoCase(scheme.value, Str2bytes("sip")) {
			return newPos, &AbnfError{"sip-uri parse: parse scheme failed: not sip-uri nor sips-uri", src, newPos}
		} else {
			this.SetSipUri()
		}

		return newPos, nil
	*/
}

func (this *SipUri) ParseUserinfo(src []byte, pos int) (newPos int, err error) {
	newPos = pos
	hasUserinfo := findUserinfo(src, newPos)
	if hasUserinfo {
		newPos, err = this.user.ParseEscapable(src, newPos, IsSipUser)
		if err != nil {
			return newPos, err
		}

		if newPos >= len(src) {
			return newPos, &AbnfError{"sip-uri parse: parse user-info failed: reach end after user", src, newPos}
		}

		if this.user.Empty() {
			return newPos, &AbnfError{"sip-uri parse: parse user-info failed: empty user", src, newPos}
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
			return newPos, &AbnfError{"sip-uri parse: parse user-info failed: reach end, and no '@'", src, newPos}
		}

		if src[newPos] != '@' {
			return newPos, &AbnfError{"sip-uri parse: parse user-info failed: no '@'", src, newPos}
		}

		newPos++
	}

	return newPos, nil
}

func findUserinfo(src []byte, pos int) bool {
	for _, v := range src[pos:] {
		if v == '@' {
			return true
		} else if v == '>' || IsLwsChar(v) {
			break
		}
	}
	return false
}

type SipUriParam struct {
	name  AbnfToken
	value AbnfToken
}

func (this *SipUriParam) Parse(src []byte, pos int) (newPos int, err error) {
	newPos, err = this.name.ParseEscapable(src, pos, IsSipPname)
	if err != nil {
		return newPos, err
	}

	if this.name.Empty() {
		return newPos, &AbnfError{"sip-uri parse: parse sip-uri param failed: empty pname", src, newPos}
	}

	this.name.SetExist()

	if src[newPos] == '=' {
		newPos, err = this.value.ParseEscapable(src, newPos+1, IsSipPvalue)
		if err != nil {
			return newPos, err
		}

		if this.value.Empty() {
			return newPos, &AbnfError{"sip-uri parse: parse sip-uri param failed: empty pvalue", src, newPos}
		}
		this.value.SetExist()
	}
	return newPos, nil
}

func (this *SipUriParam) Encode(buf *bytes.Buffer) {
	buf.Write(Escape(this.name.value, IsSipPname))
	if this.value.Exist() {
		buf.WriteByte('=')
		buf.Write(Escape(this.value.value, IsSipPvalue))
	}
}

func (this *SipUriParam) String() string {
	str := string(Escape([]byte(this.name.String()), IsSipPname))
	if this.value.Exist() {
		str += "="
		str += string(Escape([]byte(this.value.String()), IsSipPvalue))
	}
	return str
}

type SipUriParams struct {
	/*
		orders []string
		maps   map[string]*SipUriParam*/
	//list.List
	SipList
}

func NewSipUriParams() *SipUriParams {
	return &SipUriParams{}
}

func (this *SipUriParams) Encode(buf *bytes.Buffer) {
	/*for i, v := range this.orders {
		if i > 0 {
			buf.WriteByte(';')
		}
		this.maps[v].Encode(buf)
	}*/
}

func (this *SipUriParams) String() string {
	str := ""
	/*if len(this.maps) == 0 {
		return ""
	}


	for i, v := range this.orders {
		if i > 0 {
			str += ";"
		}
		str += this.maps[v].String()
	}*/
	return str
}

//func (this *SipUriParams) Size() int   { return len(this.maps) }
func (this *SipUriParams) Size() int { return this.Len() }

//func (this *SipUriParams) Empty() bool { return len(this.maps) == 0 }
func (this *SipUriParams) Empty() bool { return this.Len() == 0 }
func (this *SipUriParams) GetParam(name string) (val *SipUriParam, ok bool) {
	//val, ok = this.maps[strings.ToLower(name)]
	for e := this.Front(); e != nil; e = e.Next() {

	}
	return val, ok
}

func (this *SipUriParams) Parse(src []byte, pos int) (newPos int, err error) {
	newPos = pos
	if newPos >= len(src) {
		return newPos, &AbnfError{"sip-uri parse: parse sip-uri param failed: reach end after ';'", src, newPos}
	}

	for newPos < len(src) {
		param := SipUriParam{}
		newPos, err = param.Parse(src, newPos)
		if err != nil {
			return newPos, err
		}
		this.PushBack(param)
		/*
			name := param.name.ToLower()
			this.orders = append(this.orders, name)
			this.maps[name] = param
		*/

		if newPos >= len(src) {
			return newPos, nil
		}

		if src[newPos] != ';' {
			return newPos, nil
		}
		newPos++
	}

	return newPos, nil
}

type SipUriHeader struct {
	name  AbnfToken
	value AbnfToken
}

func (this *SipUriHeader) Parse(src []byte, pos int) (newPos int, err error) {
	newPos, err = this.name.ParseEscapable(src, pos, IsSipHname)
	if err != nil {
		return newPos, err
	}

	if this.name.Empty() {
		return newPos, &AbnfError{"sip-uri parse: parse sip-uri header failed: empty hname", src, newPos}
	}
	this.name.SetExist()

	if newPos >= len(src) {
		return newPos, &AbnfError{"sip-uri parse: parse header failed: no = after hname", src, newPos}
	}

	if src[newPos] != '=' {
		return newPos, &AbnfError{"sip-uri parse: parse header failed: no = after hname", src, newPos}
	}

	newPos, err = this.value.ParseEscapable(src, newPos+1, IsSipHvalue)
	if err != nil {
		return newPos, err
	}

	this.value.SetExist()

	return newPos, nil
}

func (this *SipUriHeader) Encode(buf *bytes.Buffer) {
	buf.Write(Escape(this.name.value, IsSipHname))
	buf.WriteByte('=')
	if this.value.Exist() {
		buf.Write(Escape(this.value.value, IsSipHvalue))
	}
}

func (this *SipUriHeader) String() string {
	str := string(Escape([]byte(this.name.String()), IsSipHname))
	str += "="
	if this.value.Exist() {
		str += string(Escape([]byte(this.value.String()), IsSipHvalue))
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

func (this *SipUriHeaders) Size() int   { return len(this.maps) }
func (this *SipUriHeaders) Empty() bool { return len(this.maps) == 0 }
func (this *SipUriHeaders) GetHeader(name string) (val *SipUriHeader, ok bool) {
	val, ok = this.maps[strings.ToLower(name)]
	return val, ok
}

func (this *SipUriHeaders) Parse(src []byte, pos int) (newPos int, err error) {
	newPos = pos
	if newPos >= len(src) {
		return newPos, &AbnfError{"sip-uri parse: parse uri-header failed: reach end after ';'", src, newPos}
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

func (this *SipUriHeaders) Encode(buf *bytes.Buffer) {
	for i, v := range this.orders {
		if i > 0 {
			buf.WriteByte('&')
		}
		this.maps[v].Encode(buf)
	}
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
	}
	return str
}
