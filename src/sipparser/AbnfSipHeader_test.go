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

func TestAbnfSipHeader_GetSipHeaderIndex2(t *testing.T) {
	for i, v := range g_SipHeaderInfos {
		v := v
		j := i
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()

			test.EXPECT_EQ(t, GetSipHeaderIndex2(v.name), SipHeaderIndexType(j), "header name = %s", v.name)
			test.EXPECT_EQ(t, GetSipHeaderIndex2(bytes.ToLower(v.name)), SipHeaderIndexType(j), "header name = %s", bytes.ToLower(v.name))
			test.EXPECT_EQ(t, GetSipHeaderIndex2(bytes.ToUpper(v.name)), SipHeaderIndexType(j), "header name = %s", bytes.ToUpper(v.name))
		})
	}
}

func TestAbnfSipHeader_GetSipHeaderIndex3(t *testing.T) {
	for i, v := range g_SipHeaderInfos {
		v := v
		j := i
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()

			test.EXPECT_EQ(t, GetSipHeaderIndex3(v.name), SipHeaderIndexType(j), "header name = %s", v.name)
			test.EXPECT_EQ(t, GetSipHeaderIndex3(bytes.ToLower(v.name)), SipHeaderIndexType(j), "header name = %s", bytes.ToLower(v.name))
			test.EXPECT_EQ(t, GetSipHeaderIndex3(bytes.ToUpper(v.name)), SipHeaderIndexType(j), "header name = %s", bytes.ToUpper(v.name))
		})
	}
}

func BenchmarkGetSipHeaderIndex(b *testing.B) {
	b.StopTimer()
	b.SetBytes(2)
	b.ReportAllocs()

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for _, v := range g_SipHeaderInfos {
			GetSipHeaderIndex(v.name)
		}
	}
}

func BenchmarkGetSipHeaderIndex2(b *testing.B) {
	b.StopTimer()
	b.SetBytes(2)
	b.ReportAllocs()

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for _, v := range g_SipHeaderInfos {
			GetSipHeaderIndex2(v.name)
		}
	}
}

func BenchmarkGetSipHeaderIndex3(b *testing.B) {
	b.StopTimer()
	b.SetBytes(2)
	b.ReportAllocs()

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for _, v := range g_SipHeaderInfos {
			GetSipHeaderIndex3(v.name)
		}
	}
}
