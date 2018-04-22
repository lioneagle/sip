package sipparser

import (
	"fmt"
	"unsafe"
)

func init() {
	for i, v := range g_SipHeaderInfos {
		v.index = SipHeaderIndexType(i)
	}

	for i := 0; i < 256; i++ {
		g_tolower_table[i] = toLower(byte(i))
		g_toupper_table[i] = toUpper(byte(i))

		g_byteAsString_table[i] = fmt.Sprintf("%d", i)
	}

	fmt.Println("sizeof(bool)                 =", unsafe.Sizeof(true))
	fmt.Println("sizeof(int)                  =", unsafe.Sizeof(1))
	fmt.Println("sizeof(AbnfPtr)              =", unsafe.Sizeof(AbnfPtr(1)))
	fmt.Println("sizeof(AbnfRef)              =", unsafe.Sizeof(AbnfRef{}))
	fmt.Println("sizeof(SipHostPort)          =", unsafe.Sizeof(SipHostPort{}))
	fmt.Println("sizeof(SipUri)               =", unsafe.Sizeof(SipUri{}))
	fmt.Println("sizeof(SipAddr)              =", unsafe.Sizeof(SipAddr{}))
	fmt.Println("sizeof(SipGenericParam)      =", unsafe.Sizeof(SipGenericParam{}))
	fmt.Println("sizeof(SipHeaderFrom)        =", unsafe.Sizeof(SipHeaderFrom{}))
	fmt.Println("sizeof(SipHeaderCallId)      =", unsafe.Sizeof(SipHeaderCallId{}))
	fmt.Println("sizeof(SipHeaderCseq)        =", unsafe.Sizeof(SipHeaderCseq{}))
	fmt.Println("sizeof(SipHeaderMaxForwards) =", unsafe.Sizeof(SipHeaderMaxForwards{}))
	fmt.Println("sizeof(SipVersion)           =", unsafe.Sizeof(SipVersion{}))
	fmt.Println("sizeof(SipStartLine)         =", unsafe.Sizeof(SipStartLine{}))
}

func toLower(ch byte) byte {
	if IsUpper(ch) {
		//return ch - 'A' + 'a'
		return ch | 0x20
	}
	return ch
}

func toUpper(ch byte) byte {
	if IsLower(ch) {
		//return ch - 'a' + 'A'
		return ch & 0xDF
	}
	return ch
}
