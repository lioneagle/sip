package sipparser

import (
	//"bytes"
	//"fmt"
	"testing"
)

func TestSipMultiHeader(t *testing.T) {
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	prefix := FuncName()

	multiHeader, _ := NewSipMultiHeader(context)
	multiHeader.SetNameByteSlice(context, []byte("Content-ABC"))

	multiHeader.GenerateAndAddHeader(context, "Content-ABC", "asdjhd")
	if multiHeader.Size() != 1 {
		t.Errorf("%s failed: size = %s, wanted = 1\n", prefix, multiHeader.Size())
	}

	multiHeader.GenerateAndAddHeader(context, "Content-ABC", "ttgxz")
	if multiHeader.Size() != 2 {
		t.Errorf("%s failed: size = %s, wanted = 2\n", prefix, multiHeader.Size())
	}

	if multiHeader.String(context) != "Content-ABC: asdjhd, ttgxz\r\n" {
		t.Errorf("%s failed: encode = %s, wanted = %s\n", prefix, multiHeader.String(context), "Content-ABC: asdjhd, ttgxz\r\n")
	}

	if !multiHeader.NameHasPrefixByteSlice(context, []byte("Content-")) {
		t.Errorf("%s failed: wrong NameHasPrefixBytes\n", prefix)
	}

	if !multiHeader.EqualNameString(context, "Content-ABC") {
		t.Errorf("%s failed: wrong EqualNameString\n", prefix)
	}
}

func TestSipMultiHeaders(t *testing.T) {

	testdata := []struct {
		name   string
		values []string
	}{
		{"Route", []string{"<sip:123@ada.com>;ax=ads", "<tel:+1233>"}},
		{"Content-xxY", []string{"adsdfd", "hht"}},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	prefix := FuncName()

	headers, _ := NewSipMultiHeaders(context)
	for _, v := range testdata {
		for _, x := range v.values {
			headers.GenerateAndAddHeader(context, v.name, x)
		}
	}

	if headers.String(context) != "Route: <sip:123@ada.com>;ax=ads, <tel:+1233>\r\nContent-xxY: adsdfd, hht\r\n" {
		t.Errorf("%s failed: encode = %s, wanted = %s\n", prefix, headers.String(context), "Route: <sip:123@ada.com>;ax=ads, <tel:+1233>\r\nContent-xxY: adsdfd, hht\r\n")
	}

	headers.RemoveContentHeaders(context)
	if headers.String(context) != "Route: <sip:123@ada.com>;ax=ads, <tel:+1233>\r\n" {
		t.Errorf("%s failed: encode = %s, wanted = %s\n", prefix, headers.String(context), "Route: <sip:123@ada.com>;ax=ads, <tel:+1233>\r\n")
	}
}

func TestSipMultiHeadersRemoveHeaderByNameString(t *testing.T) {

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	prefix := FuncName()

	headers, _ := NewSipMultiHeaders(context)
	headers.GenerateAndAddHeader(context, "Route", "<sip:123@ada.com>;ax=ads")
	headers.GenerateAndAddHeader(context, "Route", "<tel:+1233>")
	headers.GenerateAndAddHeader(context, "Content-xxY", "adsdfd")
	headers.GenerateAndAddHeader(context, "Content-xxY", "hht")

	headers.RemoveHeaderByNameString(context, "Route")

	encoded := headers.String(context)
	dst := "Content-xxY: adsdfd, hht\r\n"
	if encoded != dst {
		t.Errorf("%s failed: encode = %s, wanted = %s\n", prefix, encoded, dst)
	}

}
