package sipparser

import (
//"fmt"
//"strings"
)

type SipDisplayName struct {
	isQuotedString bool
	name           AbnfToken
	quotedstring   SipQuotedString
}

func NewSipDisplayName() *SipDisplayName {
	return &SipDisplayName{}
}

type SipNameAddr struct {
	displayname SipDisplayName
	addrsepc    SipAddrSpec
}

func NewSipNameAddr() *SipNameAddr {
	return &SipNameAddr{}
}
