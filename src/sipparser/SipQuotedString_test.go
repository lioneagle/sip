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

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	prefix := FuncName()

	for i, v := range testdata {
		addr := NewSipQuotedString(context)
		quotedString := addr.GetSipQuotedString(context)

		_, err := quotedString.Parse(context, []byte(v.src), 0)
		if err != nil && v.parseOk {
			t.Errorf("%s[%d] failed: %s\n", prefix, i, err.Error())
			continue
		}

		if err == nil && !v.parseOk {
			t.Errorf("%s[%d] failed: parse failed wanted\n", prefix, i)
			continue
		}

		if quotedString.String(context) != v.wanted {
			t.Errorf("%s[%d] failed: value = %s, wanted = %s\n", prefix, i, quotedString.String(context), v.wanted)
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
		{"\r\n\"abc\\00\"", 0},
		{"\r\n \"abc\\", len("\r\n\"abc\\")},
		{"\r\n \"abc\r\n\\", len("\r\n \"abc\r\n")},
		{"\r\n \"abc", len("\r\n \"abc")},
		{"\r\n \"abc€", len("\r\n \"abc") + 1},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	prefix := FuncName()

	for i, v := range testdata {
		addr := NewSipQuotedString(context)
		quotedString := addr.GetSipQuotedString(context)

		newPos, err := quotedString.Parse(context, []byte(v.src), 0)

		if err == nil {
			t.Errorf("%s[%d] failed: err should not be nil", prefix, i)
			continue
		}

		if newPos != v.newPos {
			fmt.Print(err)
			t.Errorf("%s[%d] failed: newPos = %d, wanted = %d\n", prefix, i, newPos, v.newPos)
			continue
		}
	}
}
