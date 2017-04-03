package sipparser

import (
	//"bytes"
	//"fmt"
	"testing"
)

func TestGenericParamParse(t *testing.T) {

	testdata := []struct {
		src    string
		ok     bool
		newPos int
		encode string
	}{
		{"a=b", true, len("a=b"), "a=b"},
		{"a\r\n\t=\r\n\tb", true, len("a\r\n\t=\r\n\tb"), "a=b"},
		{"a\r\n =\r\n b", true, len("a\r\n =\r\n b"), "a=b"},
		{"a=\"b\"", true, len("a=\"b\""), "a=\"b\""},
		{"a=\r\n\t\"b\"", true, len("a=\r\n\t\"b\""), "a=\"b\""},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)

	for i, v := range testdata {
		param, _ := NewSipGenericParam(context)
		newPos, err := param.Parse(context, []byte(v.src), 0)

		if v.ok && err != nil {
			t.Errorf("TestGenericParamParse[%d] failed, err = %s\n", i, err)
			continue
		}

		if !v.ok && err == nil {
			t.Errorf("TestGenericParamParse[%d] failed, should parse failed", i)
			continue
		}

		if v.newPos != newPos {
			t.Errorf("TestGenericParamParse[%d] failed, newPos = %d, wanted = %d\n", i, newPos, v.newPos)
		}

		if v.encode != param.String(context) {
			t.Errorf("TestGenericParamParse[%d] failed, encode = %s, wanted = %s\n", i, param.String(context), v.encode)
			continue
		}
	}

}

func TestGenericParamsParse(t *testing.T) {

	testdata := []struct {
		src    string
		ok     bool
		newPos int
		encode string
	}{
		{";a=b", true, len(";a=b"), ";a=b"},
		{";a=b;c", true, len(";a=b;c"), ";a=b;c"},
		{";a=b;c=d", true, len(";a=b;c=d"), ";a=b;c=d"},
		{";a=b\r\n\t;c=d", true, len(";a=b\r\n\t;c=d"), ";a=b;c=d"},
		{";a=b\r\n\t; c\r\n = d", true, len(";a=b\r\n\t; c\r\n = d"), ";a=b;c=d"},

		{";a=", false, len(";a="), ""},
		{";@=", false, len(";"), ""},
		{";a\r\n=", false, len(";a"), ""},
		{";a=\"ac", false, len(";a=\"ac"), ""},
		{";a=@", false, len(";a="), ""},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)

	for i, v := range testdata {
		params, _ := NewSipGenericParams(context)
		newPos, err := params.Parse(context, []byte(v.src), 0, ';')

		if v.ok && err != nil {
			t.Errorf("TestGenericParamsParse[%d] failed, err = %s\n", i, err)
			continue
		}

		if !v.ok && err == nil {
			t.Errorf("TestGenericParamsParse[%d] failed, should parse failed", i)
			continue
		}

		if v.newPos != newPos {
			t.Errorf("TestGenericParamsParse[%d] failed, newPos = %d, wanted = %d\n", i, newPos, v.newPos)
		}

		if v.ok && v.encode != params.String(context, ';') {
			t.Errorf("TestGenericParamsParse[%d] failed, encode = %s, wanted = %s\n", i, params.String(context, ';'), v.encode)
			continue
		}
	}

}
