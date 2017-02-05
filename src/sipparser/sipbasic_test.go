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

	wanted := []struct {
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

	for i, v := range wanted {
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

	for i, v := range wanted {
		u := Unescape([]byte(v.escaped))
		if bytes.Compare(u, []byte(v.unescaped)) != 0 {
			t.Errorf("TestUnescape[%d] failed, ret = %s, wanted = %s\n", i, string(u), v.unescaped)
		}
	}
}

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
		{"IsSipQuotedString", IsSipQuotedString},
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

func TestParseToken(t *testing.T) {

	wanted := []struct {
		name        string
		isInCharset func(ch byte) bool
		src         string
		err         bool
		begin       int
		end         int
		newPos      int
	}{

		{"IsDigit", IsDigit, "01234abc", false, 0, 5, 5},
		{"IsDigit", IsDigit, "56789=bc", false, 0, 5, 5},
		{"IsDigit", IsDigit, "ad6789abc", true, 0, 0, 0},
	}

	for i, v := range wanted {
		begin, end, newPos, err := parseToken([]byte(v.src), 0, v.isInCharset)
		if err != nil {
			if !v.err {
				t.Errorf("TestParseToken[%d], %s: parse failed, want success\n", i, v.name)
			}
		} else {
			if begin != v.begin {
				t.Errorf("TestParseToken[%d], %s: begin = %d, want = %d\n", i, v.name, begin, v.begin)
			} else if end != v.end {
				t.Errorf("TestParseToken[%d], %s: end = %d, want = %d\n", i, v.name, end, v.end)
			} else if newPos != v.newPos {
				t.Errorf("TestParseToken[%d], %s: newPos = %d, want = %d\n", i, v.name, newPos, v.newPos)
			}
		}
	}
}

func TestParseTokenEscapable(t *testing.T) {

	wanted := []struct {
		name        string
		isInCharset func(ch byte) bool
		src         string
		err         bool
		begin       int
		end         int
		newPos      int
	}{

		{"IsDigit", IsDigit, "01234abc", false, 0, 5, 5},
		{"IsDigit", IsDigit, "56789=bc", false, 0, 5, 5},
		{"IsDigit", IsDigit, "ad6789abc", true, 0, 0, 0},
		{"IsDigit", IsDigit, "%301234abc", false, 0, 7, 7},
		{"IsDigit", IsDigit, "%30%311234abc", false, 0, 10, 10},
		{"IsDigit", IsDigit, "%30%31123%3a", false, 0, 12, 12},
		{"IsDigit", IsDigit, "%3c%31123%", true, 0, 10, 10},
		{"IsDigit", IsDigit, "%30%31123%F", true, 0, 10, 10},
		{"IsDigit", IsDigit, "%3x%31123%F", true, 0, 0, 0},
	}

	for i, v := range wanted {
		begin, end, newPos, err := parseTokenEscapable([]byte(v.src), 0, v.isInCharset)
		if err != nil {
			if !v.err {
				t.Errorf("TestParseTokenEscapable[%d], %s: parse failed, want success\n", i, v.name)
			}
		} else {
			if begin != v.begin {
				t.Errorf("TestParseTokenEscapable[%d], %s: begin = %d, want = %d\n", i, v.name, begin, v.begin)
			} else if end != v.end {
				t.Errorf("TestParseTokenEscapable[%d], %s: end = %d, want = %d\n", i, v.name, end, v.end)
			} else if newPos != v.newPos {
				t.Errorf("TestParseTokenEscapable[%d], %s: newPos = %d, want = %d\n", i, v.name, newPos, v.newPos)
			}
		}
	}
}
