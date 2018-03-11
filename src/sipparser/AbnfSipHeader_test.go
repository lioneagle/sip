package sipparser

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/lioneagle/goutil/src/test"
)

func TestAbnfSipHeader_GetSipHeaderIndex(t *testing.T) {
	for i, v := range g_SipHeaderInfos {
		v := v
		j := i
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()

			test.EXPECT_EQ(t, GetSipHeaderIndex(v.name), SipHeaderIndexType(j), "header name = %s", v.name)
			test.EXPECT_EQ(t, GetSipHeaderIndex(bytes.ToLower(v.name)), SipHeaderIndexType(j), "header name = %s", bytes.ToLower(v.name))
			test.EXPECT_EQ(t, GetSipHeaderIndex(bytes.ToUpper(v.name)), SipHeaderIndexType(j), "header name = %s", bytes.ToUpper(v.name))
		})
	}
}
