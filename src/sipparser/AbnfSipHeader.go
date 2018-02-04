package sipparser

import (
	_ "fmt"
)

var g_SipHeaderAssoValues = []byte{
	44, 44, 44, 44, 44, 44, 44, 44, 44, 44,
	44, 44, 44, 44, 44, 44, 44, 44, 44, 44,
	44, 44, 44, 44, 44, 44, 44, 44, 44, 44,
	44, 44, 44, 44, 44, 44, 44, 44, 44, 44,
	44, 44, 44, 44, 44, 44, 44, 44, 44, 44,
	44, 44, 44, 44, 44, 44, 44, 44, 44, 44,
	44, 44, 44, 44, 44, 3, 39, 0, 32, 10,
	30, 44, 44, 4, 37, 36, 5, 27, 44, 2,
	44, 44, 1, 11, 28, 6, 26, 44, 33, 44,
	44, 44, 44, 44, 44, 44, 44, 3, 39, 0,
	32, 10, 30, 44, 44, 4, 37, 36, 5, 27,
	44, 2, 44, 44, 1, 11, 28, 6, 26, 44,
	33, 44, 44, 44, 44, 44, 44, 44, 44, 44,
	44, 44, 44, 44, 44, 44, 44, 44, 44, 44,
	44, 44, 44, 44, 44, 44, 44, 44, 44, 44,
	44, 44, 44, 44, 44, 44, 44, 44, 44, 44,
	44, 44, 44, 44, 44, 44, 44, 44, 44, 44,
	44, 44, 44, 44, 44, 44, 44, 44, 44, 44,
	44, 44, 44, 44, 44, 44, 44, 44, 44, 44,
	44, 44, 44, 44, 44, 44, 44, 44, 44, 44,
	44, 44, 44, 44, 44, 44, 44, 44, 44, 44,
	44, 44, 44, 44, 44, 44, 44, 44, 44, 44,
	44, 44, 44, 44, 44, 44, 44, 44, 44, 44,
	44, 44, 44, 44, 44, 44, 44, 44, 44, 44,
	44, 44, 44, 44, 44, 44, 44, 44, 44, 44,
	44, 44, 44, 44, 44, 44,
}

func SipHeaderHash(name []byte) int {
	ret := len(name)
	if ret == 1 {
		ret += int(g_SipHeaderAssoValues[name[0]])
	} else {
		ret += int(g_SipHeaderAssoValues[name[0]] + g_SipHeaderAssoValues[name[1]])
	}
	return ret
}

func GetSipHeaderIndex(name []byte) SipHeaderIndexType {
	len1 := len(name)
	if MIN_KNOWN_SIP_HEADER_LENGTH <= len1 && len1 <= MAX_KNOWN_SIP_HEADER_LENGTH {
		key := SipHeaderHash(name)
		if 0 <= key && key <= MAX_KNOWN_SIP_HEADER_HASH_VALUE {
			if EqualNoCase(name, g_SipHeaderGperf[key].name) {
				return g_SipHeaderGperf[key].index
			}

		}
	}
	return ABNF_SIP_HDR_UNKNOWN
}

const (
	MIN_KNOWN_SIP_HEADER_LENGTH     = 1
	MAX_KNOWN_SIP_HEADER_LENGTH     = 19
	MIN_KNOWN_SIP_HEADER_HASH_VALUE = 1
	MAX_KNOWN_SIP_HEADER_HASH_VALUE = 44
)

var g_SipHeaderGperf = []struct {
	name  []byte
	index SipHeaderIndexType
}{
	1:  {[]byte("c"), ABNF_SIP_HDR_CONTENT_TYPE},
	2:  {[]byte("r"), ABNF_SIP_HDR_REFER_TO},
	3:  {[]byte("o"), ABNF_SIP_HDR_EVENT},
	4:  {[]byte("a"), ABNF_SIP_HDR_ACCEPT_CONTACT},
	5:  {[]byte("i"), ABNF_SIP_HDR_CALL_ID},
	6:  {[]byte("l"), ABNF_SIP_HDR_CONTENT_LENGTH},
	7:  {[]byte("u"), ABNF_SIP_HDR_ALLOW_EVENTS},
	8:  {[]byte("Route"), ABNF_SIP_HDR_ROUTE},
	9:  {[]byte("Contact"), ABNF_SIP_HDR_CONTACT},
	10: {[]byte("Call-ID"), ABNF_SIP_HDR_CALL_ID},
	11: {[]byte("e"), ABNF_SIP_HDR_CONTENT_ENCODING},
	12: {[]byte("s"), ABNF_SIP_HDR_SUBJECT},
	13: {[]byte("Allow"), ABNF_SIP_HDR_ALLOW},
	14: {[]byte("Content-Type"), ABNF_SIP_HDR_CONTENT_TYPE},
	15: {[]byte("CSeq"), ABNF_SIP_HDR_CSEQ},
	16: {[]byte("Content-Length"), ABNF_SIP_HDR_CONTENT_LENGTH},
	17: {[]byte("Accept-Contact"), ABNF_SIP_HDR_ACCEPT_CONTACT},
	18: {[]byte("Content-Encoding"), ABNF_SIP_HDR_CONTENT_ENCODING},
	19: {[]byte("Refer-To"), ABNF_SIP_HDR_REFER_TO},
	20: {[]byte("Allow-Events"), ABNF_SIP_HDR_ALLOW_EVENTS},
	21: {[]byte("Content-Disposition"), ABNF_SIP_HDR_CONTENT_DISPOSITION},
	22: {[]byte("Referred-By"), ABNF_SIP_HDR_REFERRED_BY},
	23: {[]byte("Record-Route"), ABNF_SIP_HDR_RECORD_ROUTE},
	24: {[]byte("Subject"), ABNF_SIP_HDR_SUBJECT},
	25: {[]byte("Reject-Contact"), ABNF_SIP_HDR_REJECT_CONTACT},
	26: {[]byte("Supported"), ABNF_SIP_HDR_SUPPORTED},
	27: {[]byte("v"), ABNF_SIP_HDR_VIA},
	28: {[]byte("m"), ABNF_SIP_HDR_CONTACT},
	29: {[]byte("t"), ABNF_SIP_HDR_TO},
	30: {[]byte("Request-Disposition"), ABNF_SIP_HDR_REQUEST_DISPOSITION},
	31: {[]byte("f"), ABNF_SIP_HDR_FROM},
	32: {[]byte("To"), ABNF_SIP_HDR_TO},
	33: {[]byte("Via"), ABNF_SIP_HDR_VIA},
	34: {[]byte("x"), ABNF_SIP_HDR_SESSION_EXPIRES},
	35: {[]byte("From"), ABNF_SIP_HDR_FROM},
	36: {[]byte("Session-Expires"), ABNF_SIP_HDR_SESSION_EXPIRES},
	37: {[]byte("k"), ABNF_SIP_HDR_SUPPORTED},
	38: {[]byte("j"), ABNF_SIP_HDR_REJECT_CONTACT},
	39: {[]byte("Date"), ABNF_SIP_HDR_DATE},
	40: {[]byte("b"), ABNF_SIP_HDR_REFERRED_BY},
	41: {[]byte("Event"), ABNF_SIP_HDR_EVENT},
	42: {[]byte("Max-Forwards"), ABNF_SIP_HDR_MAX_FORWARDS},
	43: {[]byte("MIME-Version"), ABNF_SIP_HDR_MIME_VERSION},
}
