package sipparser

import (
	//"bytes"
	//"fmt"
	"strings"
	"testing"
)

func TestSipSingleHeader(t *testing.T) {
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	prefix := FuncName()

	header, _ := GenerateSingleHeader(context, "Content-ABC", "asdjhd")

	if header.IsParsed() {
		t.Errorf("%s failed: should not be parsed\n", prefix)
	}

	if header.HasInfo() {
		t.Errorf("%s failed: should not have info\n", prefix)
	}

	if header.GetParsed() != ABNF_PTR_NIL {
		t.Errorf("%s failed: should be ABNF_PTR_NIL\n", prefix)
	}

	if !header.NameHasPrefixBytes(context, []byte("Content-")) {
		t.Errorf("%s failed: wrong NameHasPrefixBytes\n", prefix)
	}

	if !header.EqualNameString(context, "Content-abc") {
		t.Errorf("%s failed: wrong EqualNameString\n", prefix)
	}

	if header.String(context) != "Content-ABC: asdjhd" {
		t.Errorf("%s failed: encode = %s, wanted = %s\n", prefix, header.String(context), "Content-ABC: asdjhd")
	}
}

func TestSipSingleHeaders(t *testing.T) {

	testdata := []struct {
		name  string
		value string
	}{
		{"From", "<sip:123@ada.com>;ax=ads"},
		{"Content-Length", "456"},
		{"Content-xxY", "adsdfd"},
		{"To", "<sip:123@ada.com>;ax=ads"},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	prefix := FuncName()

	headers, _ := NewSipSingleHeaders(context)

	encode1 := ""
	encode2 := ""
	for _, v := range testdata {
		headers.GenerateAndAddHeader(context, v.name, v.value)
		encode1 += v.name + ": " + v.value + "\r\n"
		if !strings.HasPrefix(v.name, "Content-") || v.name == "Content-Length" || v.name == "Content-Type" {
			encode2 += v.name + ": " + v.value + "\r\n"
		}
	}

	if headers.String(context) != encode1 {
		t.Errorf("%s failed: encode = %s, wanted = %s\n", prefix, headers.String(context), encode1)
	}

	toEncode := testdata[3].name + ": " + testdata[3].value

	to, ok := headers.GetHeaderByString(context, "To")
	if !ok {
		t.Errorf("%s failed: cannot find To header\n", prefix)
	}

	to.SetInfo("to")

	if !to.HasShortName() {
		t.Errorf("%s failed: should has short name\n", prefix)
	}

	if string(to.ShortName()) != "t" {
		t.Errorf("%s failed: wrong shor name = %s, wanted = t\n", prefix, string(to.ShortName()))
	}

	if to.String(context) != toEncode {
		t.Errorf("%s failed: toEncode = %s, wanted = %s\n", prefix, to.String(context), toEncode)
	}

	toParsed, ok := headers.GetHeaderParsedByString(context, "to")
	if !ok {
		t.Errorf("%s failed: cannot get parsed To header\n", prefix)
	}

	if toParsed.GetSipHeaderTo(context).String(context) != toEncode {
		t.Errorf("%s failed: toEncode = %s, wanted = %s\n", prefix, toParsed.GetSipHeaderTo(context).String(context), toEncode)
	}

	headers.RemoveContentHeaders(context)
	if headers.String(context) != encode2 {
		t.Errorf("%s failed: encode = %s, wanted = %s\n", prefix, headers.String(context), encode2)
	}
}

func TestSipSingleHeadersRemoveHeaderByNameString(t *testing.T) {

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	prefix := FuncName()

	headers, _ := NewSipSingleHeaders(context)
	headers.GenerateAndAddHeader(context, "Route", "<sip:123@ada.com>;ax=ads")
	headers.GenerateAndAddHeader(context, "Route", "<tel:+1233>")
	headers.GenerateAndAddHeader(context, "Content-xxY", "adsdfd")
	headers.GenerateAndAddHeader(context, "Content-xxY", "hht")

	headers.RemoveHeaderByNameString(context, "Route")

	encoded := headers.String(context)
	dst := "Content-xxY: adsdfd\r\n" +
		"Content-xxY: hht\r\n"
	if encoded != dst {
		t.Errorf("%s failed: encode = %s, wanted = %s\n", prefix, encoded, dst)
	}

}
