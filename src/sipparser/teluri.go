package sipparser

import (
//"fmt"
//"strings"
)

type TelUri struct {
	isSecure bool
	user     SipToken
	password SipToken
	hostport SipHostPort
	params   SipUriParams
	headers  SipUriHeaders
}
