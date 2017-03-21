package sipparser

import (
	"fmt"
	"testing"
)

func TestSipQuotedStringParseOK(t *testing.T) {

	testdata := []struct {
		src     string
		parseOk bool
		wanted  string
	}{
		{"\"abc\"", true, "\"abc\""},
		{"\"abc\\00\"", true, "\"abc\\00\""},
		{" \t\r\n \"abc\\00\\\"\"", true, "\"abc\\00\\\"\""},
		{" \t\r\n\t\"abc\\0b\"", true, "\"abc\\0b\""},
	}

	for i, v := range testdata {
		quotedString := NewSipQuotedString()

		_, err := quotedString.Parse([]byte(v.src), 0)
		if err != nil && v.parseOk {
			t.Errorf("TestSipQuotedStringParseOK[%d] failed, %s\n", i, err.Error())
			continue
		}

		if err == nil && !v.parseOk {
			t.Errorf("TestSipQuotedStringParseOK[%d] failed, parse failed wanted\n", i)
			continue
		}

		if quotedString.String() != v.wanted {
			t.Errorf("TestSipQuotedStringParseOK[%d] failed, value = %s, wanted = %s\n", i, quotedString.String(), v.wanted)
			continue
		}

	}
}

func TestSipQuotedStringParseNOK(t *testing.T) {

	testdata := []struct {
		src    string
		newPos int
	}{
		{"abc\"", 0},
		{"\r\n\"abc\\00\"", len("\r\n")},
		{"\r\n \"abc\\", len("\r\n\"abc\\")},
		{"\r\n \"abc\r\n\\", len("\r\n \"abc\r\n")},
		{"\r\n \"abc", len("\r\n \"abc")},
		{"\r\n \"abc€", len("\r\n \"abc") + 1},
	}

	for i, v := range testdata {
		quotedString := NewSipQuotedString()

		newPos, err := quotedString.Parse([]byte(v.src), 0)
		if err == nil {
			t.Errorf("TestSipQuotedStringParseNOK[%d] failed", i)
			continue
		}

		if newPos != v.newPos {
			fmt.Print(err)
			t.Errorf("TestSipQuotedStringParseNOK[%d] failed, newPos = %d, wanted = %d\n", i, newPos, v.newPos)
			continue
		}
	}
}
