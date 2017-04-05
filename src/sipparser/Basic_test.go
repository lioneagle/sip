package sipparser

import (
	"bytes"
	"testing"
)

func TestToLowerHex(t *testing.T) {
	src := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	wanted := []byte("0123456789abcdef")

	for i, v := range src {
		u := ToLowerHex(v)
		if u != wanted[i] {
			t.Errorf("TestToLowerHex[%d] failed, ret = %c, wanted = %c\n", i, u, wanted[i])
		}
	}
}

func TestToUpper(t *testing.T) {
	src := []byte(";[]abcdefghigklmnopqrstuvwxyz012-+")
	wanted := []byte(";[]ABCDEFGHIGKLMNOPQRSTUVWXYZ012-+")

	for i, v := range src {
		u := ToUpper(v)
		if u != wanted[i] {
			t.Errorf("TestToUpper[%d] failed, ret = %c, wanted = %c\n", i, u, wanted[i])
		}
	}
}

func TestCompareNoCase(t *testing.T) {

	testdata := []struct {
		s1  string
		s2  string
		ret int
	}{
		{"aBcf", "abck", -1},
		{"aBcf", "abcf", 0},
		{"aBcf", "abcc", 1},
		{"aBcdf", "abcf", 1},
		{"aBcdf", "abcaaf", -1},
		{"089+=abcdefghigklmnopqrstuvwxyz123", "089+=ABCDEFGHIGKLMNOPQRSTUVWXYZ123", 0},
	}

	for i, v := range testdata {
		u := CompareNoCase([]byte(v.s1), []byte(v.s2))
		if u < 0 {
			u = -1
		} else if u > 0 {
			u = 1
		}
		if u != v.ret {
			t.Errorf("TestCompareNoCase[%d] failed, ret = %d, wanted = %d\n", i, u, v.ret)
		}
	}
}

func TestEqualNoCase(t *testing.T) {

	wanted := []struct {
		s1  string
		s2  string
		ret bool
	}{
		{"aBcf", "abck", false},
		{"aBcf", "abcf", true},
		{"aBcf", "abcc", false},
		{"aBcdf", "abcf", false},
		{"aBcdf", "abcaaf", false},
		{"089+=abcdefghigklmnopqrstuvwxyz123", "089+=ABCDEFGHIGKLMNOPQRSTUVWXYZ123", true},
	}

	for i, v := range wanted {
		u := EqualNoCase([]byte(v.s1), []byte(v.s2))
		if u != v.ret {
			t.Errorf("TestCompareNoCase[%d] failed, ret = %d, wanted = %d\n", i, u, v.ret)
		}
	}
}

/*
func TestUnescape(t *testing.T) {

	wanted := []struct {
		escaped   string
		unescaped string
	}{

		{"a%42c", "aBc"},
		{"a%3B", "a;"},
		{"a%3b%42", "a;B"},
		{"ac%3", "ac%3"},
		{"ac%P3", "ac%P3"},
		{"ac%", "ac%"},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)

	for i, v := range wanted {
		u := Unescape(context, []byte(v.escaped))
		if bytes.Compare(u, []byte(v.unescaped)) != 0 {
			t.Errorf("TestUnescape[%d] failed, ret = %s, wanted = %s\n", i, string(u), v.unescaped)
		}
	}
}*/

func TestEscape(t *testing.T) {

	wanted := []struct {
		name        string
		isInCharset func(ch byte) bool
	}{
		{"IsDigit", IsDigit},
		{"IsAlpha", IsAlpha},
		{"IsLower", IsLower},
		{"IsUpper", IsUpper},
		{"IsAlphanum", IsAlphanum},
		{"IsLowerHexAlpha", IsLowerHexAlpha},
		{"IsUpperHexAlpha", IsUpperHexAlpha},
		{"IsLowerHex", IsLowerHex},
		{"IsUpperHex", IsUpperHex},
		{"IsHex", IsHex},
		{"IsCrlfChar", IsCrlfChar},
		{"IsWspChar", IsWspChar},
		{"IsLwsChar", IsLwsChar},
		{"IsUtf8Char", IsUtf8Char},

		{"IsUriUnreserved", IsUriUnreserved},
		{"IsUriReserved", IsUriReserved},
		{"IsUriUric", IsUriUric},
		{"IsUriUricNoSlash", IsUriUricNoSlash},
		{"IsUriPchar", IsUriPchar},
		{"IsUriScheme", IsUriScheme},
		{"IsUriRegName", IsUriRegName},

		{"IsSipToken", IsSipToken},
		{"IsSipSeparators", IsSipSeparators},
		{"IsSipWord", IsSipWord},
		{"IsSipQuotedPair", IsSipQuotedPair},
		{"IsSipQuotedText", IsSipQuotedText},
		{"IsSipComment", IsSipComment},
		{"IsSipUser", IsSipUser},
		{"IsSipPassword", IsSipPassword},
		{"IsSipPname", IsSipPname},
		{"IsSipPvalue", IsSipPvalue},
		{"IsSipHname", IsSipHname},
		{"IsSipHvalue", IsSipHvalue},
		{"IsSipReasonPhrase", IsSipReasonPhrase},
	}

	chars := makeFullCharset()

	for i, v := range wanted {
		u := Escape(chars, v.isInCharset)
		if bytes.Compare(Unescape(u), chars) != 0 {
			t.Errorf("TestEscape[%d]: %s failed\n", i, v.name)
		}
	}
}

func makeFullCharset() (ret []byte) {
	for i := 0; i < 256; i++ {
		ret = append(ret, byte(i))
	}
	return ret
}
