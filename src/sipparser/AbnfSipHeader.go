package sipparser

import (
	_ "fmt"
)

var g_SipHeaderAssoValues = []byte{
	45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
	45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
	45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
	45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
	45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
	45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
	45, 45, 45, 45, 45, 3, 43, 0, 33, 10,
	30, 45, 45, 4, 38, 37, 5, 27, 45, 2,
	45, 45, 1, 11, 28, 6, 26, 45, 36, 45,
	45, 45, 45, 45, 45, 45, 45, 3, 43, 0,
	33, 10, 30, 45, 45, 4, 38, 37, 5, 27,
	45, 2, 45, 45, 1, 11, 28, 6, 26, 45,
	36, 45, 45, 45, 45, 45, 45, 45, 45, 45,
	45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
	45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
	45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
	45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
	45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
	45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
	45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
	45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
	45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
	45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
	45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
	45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
	45, 45, 45, 45, 45, 45,
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

func GetSipHeaderIndex3(name []byte) SipHeaderIndexType {
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
	34: {[]byte("d"), ABNF_SIP_HDR_REQUEST_DISPOSITION},
	35: {[]byte("From"), ABNF_SIP_HDR_FROM},
	36: {[]byte("Session-Expires"), ABNF_SIP_HDR_SESSION_EXPIRES},
	37: {[]byte("x"), ABNF_SIP_HDR_SESSION_EXPIRES},
	38: {[]byte("k"), ABNF_SIP_HDR_SUPPORTED},
	39: {[]byte("j"), ABNF_SIP_HDR_REJECT_CONTACT},
	40: {[]byte("Date"), ABNF_SIP_HDR_DATE},
	41: {[]byte("Event"), ABNF_SIP_HDR_EVENT},
	42: {[]byte("Max-Forwards"), ABNF_SIP_HDR_MAX_FORWARDS},
	43: {[]byte("MIME-Version"), ABNF_SIP_HDR_MIME_VERSION},
	44: {[]byte("b"), ABNF_SIP_HDR_REFERRED_BY},
}

func GetSipHeaderIndex(src []byte) SipHeaderIndexType {
	pos := 0
	len1 := len(src)

	if pos >= len1 {
		return ABNF_SIP_HDR_UNKNOWN
	}

	switch src[pos] | 0x20 {
	case 'a':
		pos++
		if pos >= len1 {
			return ABNF_SIP_HDR_ACCEPT_CONTACT
		}
		switch src[pos] | 0x20 {
		case 'c':
			pos++
			if (pos + 11) >= len1 {
				return ABNF_SIP_HDR_UNKNOWN
			}
			if ((src[pos] | 0x20) == 'c') &&
				((src[pos+1] | 0x20) == 'e') &&
				((src[pos+2] | 0x20) == 'p') &&
				((src[pos+3] | 0x20) == 't') &&
				(src[pos+4] == '-') &&
				((src[pos+5] | 0x20) == 'c') &&
				((src[pos+6] | 0x20) == 'o') &&
				((src[pos+7] | 0x20) == 'n') &&
				((src[pos+8] | 0x20) == 't') &&
				((src[pos+9] | 0x20) == 'a') &&
				((src[pos+10] | 0x20) == 'c') &&
				((src[pos+11] | 0x20) == 't') {
				pos += 12
				if pos >= len1 {
					return ABNF_SIP_HDR_ACCEPT_CONTACT
				}
			}
			return ABNF_SIP_HDR_UNKNOWN
		case 'l':
			pos++
			if (pos < len1) && ((src[pos] | 0x20) == 'l') {
				pos++
				if (pos < len1) && ((src[pos] | 0x20) == 'o') {
					pos++
					if (pos < len1) && ((src[pos] | 0x20) == 'w') {
						pos++
						if pos >= len1 {
							return ABNF_SIP_HDR_ALLOW
						}
						if (pos + 6) >= len1 {
							return ABNF_SIP_HDR_UNKNOWN
						}
						if (src[pos] == '-') &&
							((src[pos+1] | 0x20) == 'e') &&
							((src[pos+2] | 0x20) == 'v') &&
							((src[pos+3] | 0x20) == 'e') &&
							((src[pos+4] | 0x20) == 'n') &&
							((src[pos+5] | 0x20) == 't') &&
							((src[pos+6] | 0x20) == 's') {
							pos += 7
							if pos >= len1 {
								return ABNF_SIP_HDR_ALLOW_EVENTS
							}
						}
					}
				}
			}
			return ABNF_SIP_HDR_UNKNOWN
		}
		return ABNF_SIP_HDR_UNKNOWN
	case 'b':
		pos++
		if pos >= len1 {
			return ABNF_SIP_HDR_REFERRED_BY
		}
		return ABNF_SIP_HDR_UNKNOWN
	case 'c':
		pos++
		if pos >= len1 {
			return ABNF_SIP_HDR_CONTENT_TYPE
		}
		switch src[pos] | 0x20 {
		case 'a':
			pos++
			if (pos + 4) >= len1 {
				return ABNF_SIP_HDR_UNKNOWN
			}
			if ((src[pos] | 0x20) == 'l') &&
				((src[pos+1] | 0x20) == 'l') &&
				(src[pos+2] == '-') &&
				((src[pos+3] | 0x20) == 'i') &&
				((src[pos+4] | 0x20) == 'd') {
				pos += 5
				if pos >= len1 {
					return ABNF_SIP_HDR_CALL_ID
				}
			}
			return ABNF_SIP_HDR_UNKNOWN
		case 'o':
			pos++
			if (pos < len1) && ((src[pos] | 0x20) == 'n') {
				pos++
				if (pos < len1) && ((src[pos] | 0x20) == 't') {
					pos++
					switch src[pos] | 0x20 {
					case 'a':
						pos++
						if (pos + 1) >= len1 {
							return ABNF_SIP_HDR_UNKNOWN
						}
						if ((src[pos] | 0x20) == 'c') &&
							((src[pos+1] | 0x20) == 't') {
							pos += 2
							if pos >= len1 {
								return ABNF_SIP_HDR_CONTACT
							}
						}
						return ABNF_SIP_HDR_UNKNOWN
					case 'e':
						pos++
						if (pos < len1) && ((src[pos] | 0x20) == 'n') {
							pos++
							if (pos < len1) && ((src[pos] | 0x20) == 't') {
								pos++
								if (pos < len1) && (src[pos] == '-') {
									pos++
									switch src[pos] | 0x20 {
									case 'd':
										pos++
										if (pos + 9) >= len1 {
											return ABNF_SIP_HDR_UNKNOWN
										}
										if ((src[pos] | 0x20) == 'i') &&
											((src[pos+1] | 0x20) == 's') &&
											((src[pos+2] | 0x20) == 'p') &&
											((src[pos+3] | 0x20) == 'o') &&
											((src[pos+4] | 0x20) == 's') &&
											((src[pos+5] | 0x20) == 'i') &&
											((src[pos+6] | 0x20) == 't') &&
											((src[pos+7] | 0x20) == 'i') &&
											((src[pos+8] | 0x20) == 'o') &&
											((src[pos+9] | 0x20) == 'n') {
											pos += 10
											if pos >= len1 {
												return ABNF_SIP_HDR_CONTENT_DISPOSITION
											}
										}
										return ABNF_SIP_HDR_UNKNOWN
									case 'e':
										pos++
										if (pos + 6) >= len1 {
											return ABNF_SIP_HDR_UNKNOWN
										}
										if ((src[pos] | 0x20) == 'n') &&
											((src[pos+1] | 0x20) == 'c') &&
											((src[pos+2] | 0x20) == 'o') &&
											((src[pos+3] | 0x20) == 'd') &&
											((src[pos+4] | 0x20) == 'i') &&
											((src[pos+5] | 0x20) == 'n') &&
											((src[pos+6] | 0x20) == 'g') {
											pos += 7
											if pos >= len1 {
												return ABNF_SIP_HDR_CONTENT_ENCODING
											}
										}
										return ABNF_SIP_HDR_UNKNOWN
									case 'l':
										pos++
										if (pos + 4) >= len1 {
											return ABNF_SIP_HDR_UNKNOWN
										}
										if ((src[pos] | 0x20) == 'e') &&
											((src[pos+1] | 0x20) == 'n') &&
											((src[pos+2] | 0x20) == 'g') &&
											((src[pos+3] | 0x20) == 't') &&
											((src[pos+4] | 0x20) == 'h') {
											pos += 5
											if pos >= len1 {
												return ABNF_SIP_HDR_CONTENT_LENGTH
											}
										}
										return ABNF_SIP_HDR_UNKNOWN
									case 't':
										pos++
										if (pos + 2) >= len1 {
											return ABNF_SIP_HDR_UNKNOWN
										}
										if ((src[pos] | 0x20) == 'y') &&
											((src[pos+1] | 0x20) == 'p') &&
											((src[pos+2] | 0x20) == 'e') {
											pos += 3
											if pos >= len1 {
												return ABNF_SIP_HDR_CONTENT_TYPE
											}
										}
										return ABNF_SIP_HDR_UNKNOWN
									}
								}
							}
						}
						return ABNF_SIP_HDR_UNKNOWN
					}
				}
			}
			return ABNF_SIP_HDR_UNKNOWN
		case 's':
			pos++
			if (pos + 1) >= len1 {
				return ABNF_SIP_HDR_UNKNOWN
			}
			if ((src[pos] | 0x20) == 'e') &&
				((src[pos+1] | 0x20) == 'q') {
				pos += 2
				if pos >= len1 {
					return ABNF_SIP_HDR_CSEQ
				}
			}
			return ABNF_SIP_HDR_UNKNOWN
		}
		return ABNF_SIP_HDR_UNKNOWN
	case 'd':
		pos++
		if (pos + 2) >= len1 {
			return ABNF_SIP_HDR_UNKNOWN
		}
		if ((src[pos] | 0x20) == 'a') &&
			((src[pos+1] | 0x20) == 't') &&
			((src[pos+2] | 0x20) == 'e') {
			pos += 3
			if pos >= len1 {
				return ABNF_SIP_HDR_DATE
			}
		}
		return ABNF_SIP_HDR_UNKNOWN
	case 'e':
		pos++
		if pos >= len1 {
			return ABNF_SIP_HDR_CONTENT_ENCODING
		}
		if (pos + 3) >= len1 {
			return ABNF_SIP_HDR_UNKNOWN
		}
		if ((src[pos] | 0x20) == 'v') &&
			((src[pos+1] | 0x20) == 'e') &&
			((src[pos+2] | 0x20) == 'n') &&
			((src[pos+3] | 0x20) == 't') {
			pos += 4
			if pos >= len1 {
				return ABNF_SIP_HDR_EVENT
			}
		}
		return ABNF_SIP_HDR_UNKNOWN
	case 'f':
		pos++
		if pos >= len1 {
			return ABNF_SIP_HDR_FROM
		}
		if (pos + 2) >= len1 {
			return ABNF_SIP_HDR_UNKNOWN
		}
		if ((src[pos] | 0x20) == 'r') &&
			((src[pos+1] | 0x20) == 'o') &&
			((src[pos+2] | 0x20) == 'm') {
			pos += 3
			if pos >= len1 {
				return ABNF_SIP_HDR_FROM
			}
		}
		return ABNF_SIP_HDR_UNKNOWN
	case 'i':
		pos++
		if pos >= len1 {
			return ABNF_SIP_HDR_CALL_ID
		}
		return ABNF_SIP_HDR_UNKNOWN
	case 'j':
		pos++
		if pos >= len1 {
			return ABNF_SIP_HDR_REJECT_CONTACT
		}
		return ABNF_SIP_HDR_UNKNOWN
	case 'k':
		pos++
		if pos >= len1 {
			return ABNF_SIP_HDR_SUPPORTED
		}
		return ABNF_SIP_HDR_UNKNOWN
	case 'l':
		pos++
		if pos >= len1 {
			return ABNF_SIP_HDR_CONTENT_LENGTH
		}
		return ABNF_SIP_HDR_UNKNOWN
	case 'm':
		pos++
		if pos >= len1 {
			return ABNF_SIP_HDR_CONTACT
		}
		switch src[pos] | 0x20 {
		case 'a':
			pos++
			if (pos + 9) >= len1 {
				return ABNF_SIP_HDR_UNKNOWN
			}
			if ((src[pos] | 0x20) == 'x') &&
				(src[pos+1] == '-') &&
				((src[pos+2] | 0x20) == 'f') &&
				((src[pos+3] | 0x20) == 'o') &&
				((src[pos+4] | 0x20) == 'r') &&
				((src[pos+5] | 0x20) == 'w') &&
				((src[pos+6] | 0x20) == 'a') &&
				((src[pos+7] | 0x20) == 'r') &&
				((src[pos+8] | 0x20) == 'd') &&
				((src[pos+9] | 0x20) == 's') {
				pos += 10
				if pos >= len1 {
					return ABNF_SIP_HDR_MAX_FORWARDS
				}
			}
			return ABNF_SIP_HDR_UNKNOWN
		case 'i':
			pos++
			if (pos + 9) >= len1 {
				return ABNF_SIP_HDR_UNKNOWN
			}
			if ((src[pos] | 0x20) == 'm') &&
				((src[pos+1] | 0x20) == 'e') &&
				(src[pos+2] == '-') &&
				((src[pos+3] | 0x20) == 'v') &&
				((src[pos+4] | 0x20) == 'e') &&
				((src[pos+5] | 0x20) == 'r') &&
				((src[pos+6] | 0x20) == 's') &&
				((src[pos+7] | 0x20) == 'i') &&
				((src[pos+8] | 0x20) == 'o') &&
				((src[pos+9] | 0x20) == 'n') {
				pos += 10
				if pos >= len1 {
					return ABNF_SIP_HDR_MIME_VERSION
				}
			}
			return ABNF_SIP_HDR_UNKNOWN
		}
		return ABNF_SIP_HDR_UNKNOWN
	case 'o':
		pos++
		if pos >= len1 {
			return ABNF_SIP_HDR_EVENT
		}
		return ABNF_SIP_HDR_UNKNOWN
	case 'r':
		pos++
		if pos >= len1 {
			return ABNF_SIP_HDR_REFER_TO
		}
		switch src[pos] | 0x20 {
		case 'e':
			pos++
			switch src[pos] | 0x20 {
			case 'c':
				pos++
				if (pos + 8) >= len1 {
					return ABNF_SIP_HDR_UNKNOWN
				}
				if ((src[pos] | 0x20) == 'o') &&
					((src[pos+1] | 0x20) == 'r') &&
					((src[pos+2] | 0x20) == 'd') &&
					(src[pos+3] == '-') &&
					((src[pos+4] | 0x20) == 'r') &&
					((src[pos+5] | 0x20) == 'o') &&
					((src[pos+6] | 0x20) == 'u') &&
					((src[pos+7] | 0x20) == 't') &&
					((src[pos+8] | 0x20) == 'e') {
					pos += 9
					if pos >= len1 {
						return ABNF_SIP_HDR_RECORD_ROUTE
					}
				}
				return ABNF_SIP_HDR_UNKNOWN
			case 'f':
				pos++
				if (pos < len1) && ((src[pos] | 0x20) == 'e') {
					pos++
					if (pos < len1) && ((src[pos] | 0x20) == 'r') {
						pos++
						switch src[pos] | 0x20 {
						case '-':
							pos++
							if (pos + 1) >= len1 {
								return ABNF_SIP_HDR_UNKNOWN
							}
							if ((src[pos] | 0x20) == 't') &&
								((src[pos+1] | 0x20) == 'o') {
								pos += 2
								if pos >= len1 {
									return ABNF_SIP_HDR_REFER_TO
								}
							}
							return ABNF_SIP_HDR_UNKNOWN
						case 'r':
							pos++
							if (pos + 4) >= len1 {
								return ABNF_SIP_HDR_UNKNOWN
							}
							if ((src[pos] | 0x20) == 'e') &&
								((src[pos+1] | 0x20) == 'd') &&
								(src[pos+2] == '-') &&
								((src[pos+3] | 0x20) == 'b') &&
								((src[pos+4] | 0x20) == 'y') {
								pos += 5
								if pos >= len1 {
									return ABNF_SIP_HDR_REFERRED_BY
								}
							}
							return ABNF_SIP_HDR_UNKNOWN
						}
					}
				}
				return ABNF_SIP_HDR_UNKNOWN
			case 'j':
				pos++
				if (pos + 10) >= len1 {
					return ABNF_SIP_HDR_UNKNOWN
				}
				if ((src[pos] | 0x20) == 'e') &&
					((src[pos+1] | 0x20) == 'c') &&
					((src[pos+2] | 0x20) == 't') &&
					(src[pos+3] == '-') &&
					((src[pos+4] | 0x20) == 'c') &&
					((src[pos+5] | 0x20) == 'o') &&
					((src[pos+6] | 0x20) == 'n') &&
					((src[pos+7] | 0x20) == 't') &&
					((src[pos+8] | 0x20) == 'a') &&
					((src[pos+9] | 0x20) == 'c') &&
					((src[pos+10] | 0x20) == 't') {
					pos += 11
					if pos >= len1 {
						return ABNF_SIP_HDR_REJECT_CONTACT
					}
				}
				return ABNF_SIP_HDR_UNKNOWN
			case 'q':
				pos++
				if (pos + 15) >= len1 {
					return ABNF_SIP_HDR_UNKNOWN
				}
				if ((src[pos] | 0x20) == 'u') &&
					((src[pos+1] | 0x20) == 'e') &&
					((src[pos+2] | 0x20) == 's') &&
					((src[pos+3] | 0x20) == 't') &&
					(src[pos+4] == '-') &&
					((src[pos+5] | 0x20) == 'd') &&
					((src[pos+6] | 0x20) == 'i') &&
					((src[pos+7] | 0x20) == 's') &&
					((src[pos+8] | 0x20) == 'p') &&
					((src[pos+9] | 0x20) == 'o') &&
					((src[pos+10] | 0x20) == 's') &&
					((src[pos+11] | 0x20) == 'i') &&
					((src[pos+12] | 0x20) == 't') &&
					((src[pos+13] | 0x20) == 'i') &&
					((src[pos+14] | 0x20) == 'o') &&
					((src[pos+15] | 0x20) == 'n') {
					pos += 16
					if pos >= len1 {
						return ABNF_SIP_HDR_REQUEST_DISPOSITION
					}
				}
				return ABNF_SIP_HDR_UNKNOWN
			}
			return ABNF_SIP_HDR_UNKNOWN
		case 'o':
			pos++
			if (pos + 2) >= len1 {
				return ABNF_SIP_HDR_UNKNOWN
			}
			if ((src[pos] | 0x20) == 'u') &&
				((src[pos+1] | 0x20) == 't') &&
				((src[pos+2] | 0x20) == 'e') {
				pos += 3
				if pos >= len1 {
					return ABNF_SIP_HDR_ROUTE
				}
			}
			return ABNF_SIP_HDR_UNKNOWN
		}
		return ABNF_SIP_HDR_UNKNOWN
	case 's':
		pos++
		if pos >= len1 {
			return ABNF_SIP_HDR_SUBJECT
		}
		switch src[pos] | 0x20 {
		case 'e':
			pos++
			if (pos + 12) >= len1 {
				return ABNF_SIP_HDR_UNKNOWN
			}
			if ((src[pos] | 0x20) == 's') &&
				((src[pos+1] | 0x20) == 's') &&
				((src[pos+2] | 0x20) == 'i') &&
				((src[pos+3] | 0x20) == 'o') &&
				((src[pos+4] | 0x20) == 'n') &&
				(src[pos+5] == '-') &&
				((src[pos+6] | 0x20) == 'e') &&
				((src[pos+7] | 0x20) == 'x') &&
				((src[pos+8] | 0x20) == 'p') &&
				((src[pos+9] | 0x20) == 'i') &&
				((src[pos+10] | 0x20) == 'r') &&
				((src[pos+11] | 0x20) == 'e') &&
				((src[pos+12] | 0x20) == 's') {
				pos += 13
				if pos >= len1 {
					return ABNF_SIP_HDR_SESSION_EXPIRES
				}
			}
			return ABNF_SIP_HDR_UNKNOWN
		case 'u':
			pos++
			switch src[pos] | 0x20 {
			case 'b':
				pos++
				if (pos + 3) >= len1 {
					return ABNF_SIP_HDR_UNKNOWN
				}
				if ((src[pos] | 0x20) == 'j') &&
					((src[pos+1] | 0x20) == 'e') &&
					((src[pos+2] | 0x20) == 'c') &&
					((src[pos+3] | 0x20) == 't') {
					pos += 4
					if pos >= len1 {
						return ABNF_SIP_HDR_SUBJECT
					}
				}
				return ABNF_SIP_HDR_UNKNOWN
			case 'p':
				pos++
				if (pos + 5) >= len1 {
					return ABNF_SIP_HDR_UNKNOWN
				}
				if ((src[pos] | 0x20) == 'p') &&
					((src[pos+1] | 0x20) == 'o') &&
					((src[pos+2] | 0x20) == 'r') &&
					((src[pos+3] | 0x20) == 't') &&
					((src[pos+4] | 0x20) == 'e') &&
					((src[pos+5] | 0x20) == 'd') {
					pos += 6
					if pos >= len1 {
						return ABNF_SIP_HDR_SUPPORTED
					}
				}
				return ABNF_SIP_HDR_UNKNOWN
			}
			return ABNF_SIP_HDR_UNKNOWN
		}
		return ABNF_SIP_HDR_UNKNOWN
	case 't':
		pos++
		if pos >= len1 {
			return ABNF_SIP_HDR_TO
		}
		if (pos < len1) && ((src[pos] | 0x20) == 'o') {
			pos++
			if pos >= len1 {
				return ABNF_SIP_HDR_TO
			}
		}
		return ABNF_SIP_HDR_UNKNOWN
	case 'u':
		pos++
		if pos >= len1 {
			return ABNF_SIP_HDR_ALLOW_EVENTS
		}
		return ABNF_SIP_HDR_UNKNOWN
	case 'v':
		pos++
		if pos >= len1 {
			return ABNF_SIP_HDR_VIA
		}
		if (pos + 1) >= len1 {
			return ABNF_SIP_HDR_UNKNOWN
		}
		if ((src[pos] | 0x20) == 'i') &&
			((src[pos+1] | 0x20) == 'a') {
			pos += 2
			if pos >= len1 {
				return ABNF_SIP_HDR_VIA
			}
		}
		return ABNF_SIP_HDR_UNKNOWN
	case 'x':
		pos++
		if pos >= len1 {
			return ABNF_SIP_HDR_SESSION_EXPIRES
		}
		return ABNF_SIP_HDR_UNKNOWN
	}

	return ABNF_SIP_HDR_UNKNOWN
}

func GetSipHeaderIndex4(src []byte) SipHeaderIndexType {
	pos := 0
	len1 := len(src)

	if pos >= len1 {
		return ABNF_SIP_HDR_UNKNOWN
	}

	switch src[pos] | 0x20 {
	case 'a':
		pos++
		if pos >= len1 {
			return ABNF_SIP_HDR_ACCEPT_CONTACT
		}
		switch src[pos] | 0x20 {
		case 'c':
			pos++
			if (pos < len1) && ((src[pos] | 0x20) == 'c') {
				pos++
				if (pos < len1) && ((src[pos] | 0x20) == 'e') {
					pos++
					if (pos < len1) && ((src[pos] | 0x20) == 'p') {
						pos++
						if (pos < len1) && ((src[pos] | 0x20) == 't') {
							pos++
							if (pos < len1) && (src[pos] == '-') {
								pos++
								if (pos < len1) && ((src[pos] | 0x20) == 'c') {
									pos++
									if (pos < len1) && ((src[pos] | 0x20) == 'o') {
										pos++
										if (pos < len1) && ((src[pos] | 0x20) == 'n') {
											pos++
											if (pos < len1) && ((src[pos] | 0x20) == 't') {
												pos++
												if (pos < len1) && ((src[pos] | 0x20) == 'a') {
													pos++
													if (pos < len1) && ((src[pos] | 0x20) == 'c') {
														pos++
														if (pos < len1) && ((src[pos] | 0x20) == 't') {
															pos++
															if pos >= len1 {
																return ABNF_SIP_HDR_ACCEPT_CONTACT
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
			return ABNF_SIP_HDR_UNKNOWN
		case 'l':
			pos++
			if (pos < len1) && ((src[pos] | 0x20) == 'l') {
				pos++
				if (pos < len1) && ((src[pos] | 0x20) == 'o') {
					pos++
					if (pos < len1) && ((src[pos] | 0x20) == 'w') {
						pos++
						if pos >= len1 {
							return ABNF_SIP_HDR_ALLOW
						}
						if (pos < len1) && (src[pos] == '-') {
							pos++
							if (pos < len1) && ((src[pos] | 0x20) == 'e') {
								pos++
								if (pos < len1) && ((src[pos] | 0x20) == 'v') {
									pos++
									if (pos < len1) && ((src[pos] | 0x20) == 'e') {
										pos++
										if (pos < len1) && ((src[pos] | 0x20) == 'n') {
											pos++
											if (pos < len1) && ((src[pos] | 0x20) == 't') {
												pos++
												if (pos < len1) && ((src[pos] | 0x20) == 's') {
													pos++
													if pos >= len1 {
														return ABNF_SIP_HDR_ALLOW_EVENTS
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
			return ABNF_SIP_HDR_UNKNOWN
		}
		return ABNF_SIP_HDR_UNKNOWN
	case 'b':
		pos++
		if pos >= len1 {
			return ABNF_SIP_HDR_REFERRED_BY
		}
		return ABNF_SIP_HDR_UNKNOWN
	case 'c':
		pos++
		if pos >= len1 {
			return ABNF_SIP_HDR_CONTENT_TYPE
		}
		switch src[pos] | 0x20 {
		case 'a':
			pos++
			if (pos < len1) && ((src[pos] | 0x20) == 'l') {
				pos++
				if (pos < len1) && ((src[pos] | 0x20) == 'l') {
					pos++
					if (pos < len1) && (src[pos] == '-') {
						pos++
						if (pos < len1) && ((src[pos] | 0x20) == 'i') {
							pos++
							if (pos < len1) && ((src[pos] | 0x20) == 'd') {
								pos++
								if pos >= len1 {
									return ABNF_SIP_HDR_CALL_ID
								}
							}
						}
					}
				}
			}
			return ABNF_SIP_HDR_UNKNOWN
		case 'o':
			pos++
			if (pos < len1) && ((src[pos] | 0x20) == 'n') {
				pos++
				if (pos < len1) && ((src[pos] | 0x20) == 't') {
					pos++
					switch src[pos] | 0x20 {
					case 'a':
						pos++
						if (pos < len1) && ((src[pos] | 0x20) == 'c') {
							pos++
							if (pos < len1) && ((src[pos] | 0x20) == 't') {
								pos++
								if pos >= len1 {
									return ABNF_SIP_HDR_CONTACT
								}
							}
						}
						return ABNF_SIP_HDR_UNKNOWN
					case 'e':
						pos++
						if (pos < len1) && ((src[pos] | 0x20) == 'n') {
							pos++
							if (pos < len1) && ((src[pos] | 0x20) == 't') {
								pos++
								if (pos < len1) && (src[pos] == '-') {
									pos++
									switch src[pos] | 0x20 {
									case 'd':
										pos++
										if (pos < len1) && ((src[pos] | 0x20) == 'i') {
											pos++
											if (pos < len1) && ((src[pos] | 0x20) == 's') {
												pos++
												if (pos < len1) && ((src[pos] | 0x20) == 'p') {
													pos++
													if (pos < len1) && ((src[pos] | 0x20) == 'o') {
														pos++
														if (pos < len1) && ((src[pos] | 0x20) == 's') {
															pos++
															if (pos < len1) && ((src[pos] | 0x20) == 'i') {
																pos++
																if (pos < len1) && ((src[pos] | 0x20) == 't') {
																	pos++
																	if (pos < len1) && ((src[pos] | 0x20) == 'i') {
																		pos++
																		if (pos < len1) && ((src[pos] | 0x20) == 'o') {
																			pos++
																			if (pos < len1) && ((src[pos] | 0x20) == 'n') {
																				pos++
																				if pos >= len1 {
																					return ABNF_SIP_HDR_CONTENT_DISPOSITION
																				}
																			}
																		}
																	}
																}
															}
														}
													}
												}
											}
										}
										return ABNF_SIP_HDR_UNKNOWN
									case 'e':
										pos++
										if (pos < len1) && ((src[pos] | 0x20) == 'n') {
											pos++
											if (pos < len1) && ((src[pos] | 0x20) == 'c') {
												pos++
												if (pos < len1) && ((src[pos] | 0x20) == 'o') {
													pos++
													if (pos < len1) && ((src[pos] | 0x20) == 'd') {
														pos++
														if (pos < len1) && ((src[pos] | 0x20) == 'i') {
															pos++
															if (pos < len1) && ((src[pos] | 0x20) == 'n') {
																pos++
																if (pos < len1) && ((src[pos] | 0x20) == 'g') {
																	pos++
																	if pos >= len1 {
																		return ABNF_SIP_HDR_CONTENT_ENCODING
																	}
																}
															}
														}
													}
												}
											}
										}
										return ABNF_SIP_HDR_UNKNOWN
									case 'l':
										pos++
										if (pos < len1) && ((src[pos] | 0x20) == 'e') {
											pos++
											if (pos < len1) && ((src[pos] | 0x20) == 'n') {
												pos++
												if (pos < len1) && ((src[pos] | 0x20) == 'g') {
													pos++
													if (pos < len1) && ((src[pos] | 0x20) == 't') {
														pos++
														if (pos < len1) && ((src[pos] | 0x20) == 'h') {
															pos++
															if pos >= len1 {
																return ABNF_SIP_HDR_CONTENT_LENGTH
															}
														}
													}
												}
											}
										}
										return ABNF_SIP_HDR_UNKNOWN
									case 't':
										pos++
										if (pos < len1) && ((src[pos] | 0x20) == 'y') {
											pos++
											if (pos < len1) && ((src[pos] | 0x20) == 'p') {
												pos++
												if (pos < len1) && ((src[pos] | 0x20) == 'e') {
													pos++
													if pos >= len1 {
														return ABNF_SIP_HDR_CONTENT_TYPE
													}
												}
											}
										}
										return ABNF_SIP_HDR_UNKNOWN
									}
								}
							}
						}
						return ABNF_SIP_HDR_UNKNOWN
					}
				}
			}
			return ABNF_SIP_HDR_UNKNOWN
		case 's':
			pos++
			if (pos < len1) && ((src[pos] | 0x20) == 'e') {
				pos++
				if (pos < len1) && ((src[pos] | 0x20) == 'q') {
					pos++
					if pos >= len1 {
						return ABNF_SIP_HDR_CSEQ
					}
				}
			}
			return ABNF_SIP_HDR_UNKNOWN
		}
		return ABNF_SIP_HDR_UNKNOWN
	case 'd':
		pos++
		if (pos < len1) && ((src[pos] | 0x20) == 'a') {
			pos++
			if (pos < len1) && ((src[pos] | 0x20) == 't') {
				pos++
				if (pos < len1) && ((src[pos] | 0x20) == 'e') {
					pos++
					if pos >= len1 {
						return ABNF_SIP_HDR_DATE
					}
				}
			}
		}
		return ABNF_SIP_HDR_UNKNOWN
	case 'e':
		pos++
		if pos >= len1 {
			return ABNF_SIP_HDR_CONTENT_ENCODING
		}
		if (pos < len1) && ((src[pos] | 0x20) == 'v') {
			pos++
			if (pos < len1) && ((src[pos] | 0x20) == 'e') {
				pos++
				if (pos < len1) && ((src[pos] | 0x20) == 'n') {
					pos++
					if (pos < len1) && ((src[pos] | 0x20) == 't') {
						pos++
						if pos >= len1 {
							return ABNF_SIP_HDR_EVENT
						}
					}
				}
			}
		}
		return ABNF_SIP_HDR_UNKNOWN
	case 'f':
		pos++
		if pos >= len1 {
			return ABNF_SIP_HDR_FROM
		}
		if (pos < len1) && ((src[pos] | 0x20) == 'r') {
			pos++
			if (pos < len1) && ((src[pos] | 0x20) == 'o') {
				pos++
				if (pos < len1) && ((src[pos] | 0x20) == 'm') {
					pos++
					if pos >= len1 {
						return ABNF_SIP_HDR_FROM
					}
				}
			}
		}
		return ABNF_SIP_HDR_UNKNOWN
	case 'i':
		pos++
		if pos >= len1 {
			return ABNF_SIP_HDR_CALL_ID
		}
		return ABNF_SIP_HDR_UNKNOWN
	case 'j':
		pos++
		if pos >= len1 {
			return ABNF_SIP_HDR_REJECT_CONTACT
		}
		return ABNF_SIP_HDR_UNKNOWN
	case 'k':
		pos++
		if pos >= len1 {
			return ABNF_SIP_HDR_SUPPORTED
		}
		return ABNF_SIP_HDR_UNKNOWN
	case 'l':
		pos++
		if pos >= len1 {
			return ABNF_SIP_HDR_CONTENT_LENGTH
		}
		return ABNF_SIP_HDR_UNKNOWN
	case 'm':
		pos++
		if pos >= len1 {
			return ABNF_SIP_HDR_CONTACT
		}
		switch src[pos] | 0x20 {
		case 'a':
			pos++
			if (pos < len1) && ((src[pos] | 0x20) == 'x') {
				pos++
				if (pos < len1) && (src[pos] == '-') {
					pos++
					if (pos < len1) && ((src[pos] | 0x20) == 'f') {
						pos++
						if (pos < len1) && ((src[pos] | 0x20) == 'o') {
							pos++
							if (pos < len1) && ((src[pos] | 0x20) == 'r') {
								pos++
								if (pos < len1) && ((src[pos] | 0x20) == 'w') {
									pos++
									if (pos < len1) && ((src[pos] | 0x20) == 'a') {
										pos++
										if (pos < len1) && ((src[pos] | 0x20) == 'r') {
											pos++
											if (pos < len1) && ((src[pos] | 0x20) == 'd') {
												pos++
												if (pos < len1) && ((src[pos] | 0x20) == 's') {
													pos++
													if pos >= len1 {
														return ABNF_SIP_HDR_MAX_FORWARDS
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
			return ABNF_SIP_HDR_UNKNOWN
		case 'i':
			pos++
			if (pos < len1) && ((src[pos] | 0x20) == 'm') {
				pos++
				if (pos < len1) && ((src[pos] | 0x20) == 'e') {
					pos++
					if (pos < len1) && (src[pos] == '-') {
						pos++
						if (pos < len1) && ((src[pos] | 0x20) == 'v') {
							pos++
							if (pos < len1) && ((src[pos] | 0x20) == 'e') {
								pos++
								if (pos < len1) && ((src[pos] | 0x20) == 'r') {
									pos++
									if (pos < len1) && ((src[pos] | 0x20) == 's') {
										pos++
										if (pos < len1) && ((src[pos] | 0x20) == 'i') {
											pos++
											if (pos < len1) && ((src[pos] | 0x20) == 'o') {
												pos++
												if (pos < len1) && ((src[pos] | 0x20) == 'n') {
													pos++
													if pos >= len1 {
														return ABNF_SIP_HDR_MIME_VERSION
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
			return ABNF_SIP_HDR_UNKNOWN
		}
		return ABNF_SIP_HDR_UNKNOWN
	case 'o':
		pos++
		if pos >= len1 {
			return ABNF_SIP_HDR_EVENT
		}
		return ABNF_SIP_HDR_UNKNOWN
	case 'r':
		pos++
		if pos >= len1 {
			return ABNF_SIP_HDR_REFER_TO
		}
		switch src[pos] | 0x20 {
		case 'e':
			pos++
			switch src[pos] | 0x20 {
			case 'c':
				pos++
				if (pos < len1) && ((src[pos] | 0x20) == 'o') {
					pos++
					if (pos < len1) && ((src[pos] | 0x20) == 'r') {
						pos++
						if (pos < len1) && ((src[pos] | 0x20) == 'd') {
							pos++
							if (pos < len1) && (src[pos] == '-') {
								pos++
								if (pos < len1) && ((src[pos] | 0x20) == 'r') {
									pos++
									if (pos < len1) && ((src[pos] | 0x20) == 'o') {
										pos++
										if (pos < len1) && ((src[pos] | 0x20) == 'u') {
											pos++
											if (pos < len1) && ((src[pos] | 0x20) == 't') {
												pos++
												if (pos < len1) && ((src[pos] | 0x20) == 'e') {
													pos++
													if pos >= len1 {
														return ABNF_SIP_HDR_RECORD_ROUTE
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
				return ABNF_SIP_HDR_UNKNOWN
			case 'f':
				pos++
				if (pos < len1) && ((src[pos] | 0x20) == 'e') {
					pos++
					if (pos < len1) && ((src[pos] | 0x20) == 'r') {
						pos++
						switch src[pos] | 0x20 {
						case '-':
							pos++
							if (pos < len1) && ((src[pos] | 0x20) == 't') {
								pos++
								if (pos < len1) && ((src[pos] | 0x20) == 'o') {
									pos++
									if pos >= len1 {
										return ABNF_SIP_HDR_REFER_TO
									}
								}
							}
							return ABNF_SIP_HDR_UNKNOWN
						case 'r':
							pos++
							if (pos < len1) && ((src[pos] | 0x20) == 'e') {
								pos++
								if (pos < len1) && ((src[pos] | 0x20) == 'd') {
									pos++
									if (pos < len1) && (src[pos] == '-') {
										pos++
										if (pos < len1) && ((src[pos] | 0x20) == 'b') {
											pos++
											if (pos < len1) && ((src[pos] | 0x20) == 'y') {
												pos++
												if pos >= len1 {
													return ABNF_SIP_HDR_REFERRED_BY
												}
											}
										}
									}
								}
							}
							return ABNF_SIP_HDR_UNKNOWN
						}
					}
				}
				return ABNF_SIP_HDR_UNKNOWN
			case 'j':
				pos++
				if (pos < len1) && ((src[pos] | 0x20) == 'e') {
					pos++
					if (pos < len1) && ((src[pos] | 0x20) == 'c') {
						pos++
						if (pos < len1) && ((src[pos] | 0x20) == 't') {
							pos++
							if (pos < len1) && (src[pos] == '-') {
								pos++
								if (pos < len1) && ((src[pos] | 0x20) == 'c') {
									pos++
									if (pos < len1) && ((src[pos] | 0x20) == 'o') {
										pos++
										if (pos < len1) && ((src[pos] | 0x20) == 'n') {
											pos++
											if (pos < len1) && ((src[pos] | 0x20) == 't') {
												pos++
												if (pos < len1) && ((src[pos] | 0x20) == 'a') {
													pos++
													if (pos < len1) && ((src[pos] | 0x20) == 'c') {
														pos++
														if (pos < len1) && ((src[pos] | 0x20) == 't') {
															pos++
															if pos >= len1 {
																return ABNF_SIP_HDR_REJECT_CONTACT
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
				return ABNF_SIP_HDR_UNKNOWN
			case 'q':
				pos++
				if (pos < len1) && ((src[pos] | 0x20) == 'u') {
					pos++
					if (pos < len1) && ((src[pos] | 0x20) == 'e') {
						pos++
						if (pos < len1) && ((src[pos] | 0x20) == 's') {
							pos++
							if (pos < len1) && ((src[pos] | 0x20) == 't') {
								pos++
								if (pos < len1) && (src[pos] == '-') {
									pos++
									if (pos < len1) && ((src[pos] | 0x20) == 'd') {
										pos++
										if (pos < len1) && ((src[pos] | 0x20) == 'i') {
											pos++
											if (pos < len1) && ((src[pos] | 0x20) == 's') {
												pos++
												if (pos < len1) && ((src[pos] | 0x20) == 'p') {
													pos++
													if (pos < len1) && ((src[pos] | 0x20) == 'o') {
														pos++
														if (pos < len1) && ((src[pos] | 0x20) == 's') {
															pos++
															if (pos < len1) && ((src[pos] | 0x20) == 'i') {
																pos++
																if (pos < len1) && ((src[pos] | 0x20) == 't') {
																	pos++
																	if (pos < len1) && ((src[pos] | 0x20) == 'i') {
																		pos++
																		if (pos < len1) && ((src[pos] | 0x20) == 'o') {
																			pos++
																			if (pos < len1) && ((src[pos] | 0x20) == 'n') {
																				pos++
																				if pos >= len1 {
																					return ABNF_SIP_HDR_REQUEST_DISPOSITION
																				}
																			}
																		}
																	}
																}
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
				return ABNF_SIP_HDR_UNKNOWN
			}
			return ABNF_SIP_HDR_UNKNOWN
		case 'o':
			pos++
			if (pos < len1) && ((src[pos] | 0x20) == 'u') {
				pos++
				if (pos < len1) && ((src[pos] | 0x20) == 't') {
					pos++
					if (pos < len1) && ((src[pos] | 0x20) == 'e') {
						pos++
						if pos >= len1 {
							return ABNF_SIP_HDR_ROUTE
						}
					}
				}
			}
			return ABNF_SIP_HDR_UNKNOWN
		}
		return ABNF_SIP_HDR_UNKNOWN
	case 's':
		pos++
		if pos >= len1 {
			return ABNF_SIP_HDR_SUBJECT
		}
		switch src[pos] | 0x20 {
		case 'e':
			pos++
			if (pos < len1) && ((src[pos] | 0x20) == 's') {
				pos++
				if (pos < len1) && ((src[pos] | 0x20) == 's') {
					pos++
					if (pos < len1) && ((src[pos] | 0x20) == 'i') {
						pos++
						if (pos < len1) && ((src[pos] | 0x20) == 'o') {
							pos++
							if (pos < len1) && ((src[pos] | 0x20) == 'n') {
								pos++
								if (pos < len1) && (src[pos] == '-') {
									pos++
									if (pos < len1) && ((src[pos] | 0x20) == 'e') {
										pos++
										if (pos < len1) && ((src[pos] | 0x20) == 'x') {
											pos++
											if (pos < len1) && ((src[pos] | 0x20) == 'p') {
												pos++
												if (pos < len1) && ((src[pos] | 0x20) == 'i') {
													pos++
													if (pos < len1) && ((src[pos] | 0x20) == 'r') {
														pos++
														if (pos < len1) && ((src[pos] | 0x20) == 'e') {
															pos++
															if (pos < len1) && ((src[pos] | 0x20) == 's') {
																pos++
																if pos >= len1 {
																	return ABNF_SIP_HDR_SESSION_EXPIRES
																}
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
			return ABNF_SIP_HDR_UNKNOWN
		case 'u':
			pos++
			switch src[pos] | 0x20 {
			case 'b':
				pos++
				if (pos < len1) && ((src[pos] | 0x20) == 'j') {
					pos++
					if (pos < len1) && ((src[pos] | 0x20) == 'e') {
						pos++
						if (pos < len1) && ((src[pos] | 0x20) == 'c') {
							pos++
							if (pos < len1) && ((src[pos] | 0x20) == 't') {
								pos++
								if pos >= len1 {
									return ABNF_SIP_HDR_SUBJECT
								}
							}
						}
					}
				}
				return ABNF_SIP_HDR_UNKNOWN
			case 'p':
				pos++
				if (pos < len1) && ((src[pos] | 0x20) == 'p') {
					pos++
					if (pos < len1) && ((src[pos] | 0x20) == 'o') {
						pos++
						if (pos < len1) && ((src[pos] | 0x20) == 'r') {
							pos++
							if (pos < len1) && ((src[pos] | 0x20) == 't') {
								pos++
								if (pos < len1) && ((src[pos] | 0x20) == 'e') {
									pos++
									if (pos < len1) && ((src[pos] | 0x20) == 'd') {
										pos++
										if pos >= len1 {
											return ABNF_SIP_HDR_SUPPORTED
										}
									}
								}
							}
						}
					}
				}
				return ABNF_SIP_HDR_UNKNOWN
			}
			return ABNF_SIP_HDR_UNKNOWN
		}
		return ABNF_SIP_HDR_UNKNOWN
	case 't':
		pos++
		if pos >= len1 {
			return ABNF_SIP_HDR_TO
		}
		if (pos < len1) && ((src[pos] | 0x20) == 'o') {
			pos++
			if pos >= len1 {
				return ABNF_SIP_HDR_TO
			}
		}
		return ABNF_SIP_HDR_UNKNOWN
	case 'u':
		pos++
		if pos >= len1 {
			return ABNF_SIP_HDR_ALLOW_EVENTS
		}
		return ABNF_SIP_HDR_UNKNOWN
	case 'v':
		pos++
		if pos >= len1 {
			return ABNF_SIP_HDR_VIA
		}
		if (pos < len1) && ((src[pos] | 0x20) == 'i') {
			pos++
			if (pos < len1) && ((src[pos] | 0x20) == 'a') {
				pos++
				if pos >= len1 {
					return ABNF_SIP_HDR_VIA
				}
			}
		}
		return ABNF_SIP_HDR_UNKNOWN
	case 'x':
		pos++
		if pos >= len1 {
			return ABNF_SIP_HDR_SESSION_EXPIRES
		}
		return ABNF_SIP_HDR_UNKNOWN
	}

	return ABNF_SIP_HDR_UNKNOWN
}
