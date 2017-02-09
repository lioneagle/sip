package sipparser

import (
	"fmt"
)

type SipQuotedString struct {
	value []byte
}

func (this *SipQuotedString) String() string        { return fmt.Sprintf("\"%s\"", string(this.value)) }
func (this *SipQuotedString) SetValue(value []byte) { this.value = value }

func (this *SipQuotedString) Parse(src []byte, pos int) (newPos int, err error) {
	/* RFC3261 Section 25.1, page 222
	 *
	 * quoted-string  =  SWS DQUOTE *(qdtext / quoted-pair ) DQUOTE
	 * qdtext         =  LWS / %x21 / %x23-5B / %x5D-7E
	 *                 / UTF8-NONASCII
	 * quoted-pair  =  "\" (%x00-09 / %x0B-0C
	 *               / %x0E-7F)
	 */

	return newPos, nil
}
