package sipparser

import (
	"testing"
)

func TestAbnfRefParse(t *testing.T) {

	testdata := []struct {
		name        string
		isInCharset func(ch byte) bool
		src         string
		err         bool
		begin       int32
		end         int32
		newPos      int
	}{

		{"IsDigit", IsDigit, "01234abc", false, 0, 5, 5},
		{"IsDigit", IsDigit, "56789=bc", false, 0, 5, 5},
		{"IsDigit", IsDigit, "ad6789abc", true, 0, 0, 0},
	}
	prefix := FuncName()

	for i, v := range testdata {
		ref := AbnfRef{}
		newPos := ref.Parse([]byte(v.src), 0, v.isInCharset)

		if ref.Begin != v.begin {
			t.Errorf("%s[%d]: begin = %d, wanted = %d\n", prefix, i, v.name, ref.Begin, v.begin)
			continue
		}

		if ref.End != v.end {
			t.Errorf("%s[%d]: end = %d, wanted = %d\n", prefix, i, v.name, ref.End, v.end)
			continue
		}

		if newPos != v.newPos {
			t.Errorf("%s[%d]: newPos = %d, wanted = %d\n", prefix, i, v.name, newPos, v.newPos)
			continue
		}
	}
}

func TestAbnfRefParseEscapable(t *testing.T) {

	testdata := []struct {
		name        string
		isInCharset func(ch byte) bool
		src         string
		ok          bool
		begin       int32
		end         int32
		newPos      int
		escapeNum   int
	}{
		{"IsDigit", IsDigit, "01234abc", true, 0, 5, 5, 0},
		{"IsDigit", IsDigit, "56789=bc", true, 0, 5, 5, 0},
		{"IsDigit", IsDigit, "%301234abc", true, 0, 7, 7, 1},
		{"IsDigit", IsDigit, "%30%311234abc", true, 0, 10, 10, 2},
		{"IsDigit", IsDigit, "%30%31123%3a", true, 0, 12, 12, 3},
		{"IsDigit", IsDigit, "ad6789abc", true, 0, 0, 0, 0},

		{"IsDigit", IsDigit, "%3c%31123%", false, 0, 10, 9, 2},
		{"IsDigit", IsDigit, "%30%31123%F", false, 0, 10, 9, 2},
		{"IsDigit", IsDigit, "%3x%31123%F", false, 0, 0, 0, 2},
	}

	prefix := FuncName()

	for i, v := range testdata {

		ref := AbnfRef{}
		escapeNum, newPos, err := ref.ParseEscapable([]byte(v.src), 0, v.isInCharset)

		if err != nil && v.ok {
			t.Errorf("%s[%d]: should be ok\n", prefix, i, v.name)
			continue
		}

		if err == nil && !v.ok {
			t.Errorf("%s[%d]: should not be ok\n", prefix, i, v.name)
			continue
		}

		if newPos != v.newPos {
			t.Errorf("%s[%d]: newPos = %d, wanted = %d\n", prefix, i, v.name, newPos, v.newPos)
			continue
		}

		if err != nil {
			continue
		}

		if ref.Begin != v.begin {
			t.Errorf("%s[%d]: begin = %d, wanted = %d\n", prefix, i, v.name, ref.Begin, v.begin)
			continue
		}

		if ref.End != v.end {
			t.Errorf("%s[%d]: end = %d, wanted = %d\n", prefix, i, v.name, ref.End, v.end)
			continue
		}

		if escapeNum != v.escapeNum {
			t.Errorf("%s[%d]: escapeNum = %d, wanted = %d\n", prefix, i, v.name, escapeNum, v.escapeNum)
			continue
		}
	}
}
