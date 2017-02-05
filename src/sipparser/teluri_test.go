package sipparser

import (
	//"fmt"
	"testing"
)

func TestTelUriParseOK(t *testing.T) {
	testdata := []struct {
		src                      string
		isGlobalNumber           bool
		number                   string
		phoneContext             string
		phoneContextIsDomainName bool
	}{
		{"tel:+861234", true, "+861234", "", false},
		{"tel:861234;phone-context=+123", false, "861234", "+123", false},
		{"tel:861234;phone-context=+123", false, "861234", "+123", false},
		{"tel:861234;phone-context=a.com", false, "861234", "a.com", true},
	}

	for i, v := range testdata {
		uri := NewTelUri()

		newPos, err := uri.Parse([]byte(v.src), 0)
		if err != nil {
			t.Errorf("TestSipUriUserinfoParseOK[%d] failed, %s\n", i, err.Error())
			continue
		}

		if v.isGlobalNumber && !uri.IsGlobalNumber() {
			t.Errorf("TestSipUriUserinfoParseOK[%d] failed, global-number wanted\n", i)
			continue
		}

		if !v.isGlobalNumber && !uri.IsLocalNumber() {
			t.Errorf("TestSipUriUserinfoParseOK[%d] failed, local-number wanted\n", i)
			continue
		}

		if newPos != len(v.src) {
			t.Errorf("TestSipUriUserinfoParseOK[%d] failed, newPos = %d, wanted = %d\n", i, newPos, len(v.src))
			continue
		}

		if uri.number.String() != v.number {
			t.Errorf("TestSipUriUserinfoParseOK[%d] failed, number = %s, wanted = %s\n", i, uri.number.String(), v.number)
			continue
		}

		if uri.context.desc.String() != v.phoneContext {
			t.Errorf("TestSipUriUserinfoParseOK[%d] failed, user wrong, user = %s, wanted = %s", i, uri.context.desc.String(), v.phoneContext)
			continue
		}

		if v.phoneContextIsDomainName && !uri.context.isDomainName {
			t.Errorf("TestSipUriUserinfoParseOK[%d] failed, phone-context should be domain-name", i)
			continue
		}

		if !v.phoneContextIsDomainName && uri.context.isDomainName {
			t.Errorf("TestSipUriUserinfoParseOK[%d] failed, phone-context should be global-number", i)
			continue
		}
	}
}

func TestTelUriParamsParseOK(t *testing.T) {
	uri := NewTelUri()
	src := "tel:+86123;ttl=10;user%32=phone%31;a;b;c;d;e"

	_, err := uri.Parse([]byte(src), 0)
	if err != nil {
		t.Errorf("TestTelUriParamsParseOK failed, err = %s\n", err.Error())
		return
	}

	testdata := []struct {
		name     string
		value    string
		hasValue bool
	}{
		{"Ttl", "10", true},
		{"UseR2", "phone1", true},
		{"A", "", false},
		{"b", "", false},
		{"c", "", false},
		{"D", "", false},
		{"E", "", false},
	}

	for i, v := range testdata {
		param, ok := uri.params.GetParam(v.name)
		if !ok {
			t.Errorf("TestTelUriParamsParseOK[%d] failed, cannot get ttl param\n", i)
			continue
		}

		if param.value.Exist() && !v.hasValue {
			t.Errorf("TestTelUriParamsParseOK[%d] failed, should have no pvalue\n", i)
			continue
		}

		if !param.value.Exist() && v.hasValue {
			t.Errorf("TestTelUriParamsParseOK[%d] failed, should have pvalue\n", i)
			continue
		}

		if param.value.Exist() && param.value.String() != v.value {
			t.Errorf("TestTelUriParamsParseOK[%d] failed, pvalue = %s, wanted = %s\n", i, param.value.String(), v.value)
			continue
		}

	}
}

func TestTelUriEncode(t *testing.T) {

	testdata := []struct {
		src string
		dst string
	}{
		{"tel:+861234;phonex=+123", "tel:+861234;phonex=+123"},
		{"tel:861234;phone-context=+123", "tel:861234;phone-context=+123"},
		{"tel:861234;x1=5;y;phone-context=+123;zz", "tel:861234;phone-context=+123;x1=5;y;zz"},
	}

	for i, v := range testdata {
		uri := NewTelUri()

		_, err := uri.Parse([]byte(v.src), 0)
		if err != nil {
			t.Errorf("TestTelUriEncode[%d] failed, parse failed, err = %s\n", i, err.Error())
			continue
		}

		str := uri.String()

		if str != v.dst {
			t.Errorf("TestTelUriEncode[%d] failed, uri = %s, wanted = %s\n", i, str, v.dst)
			continue
		}
	}
}
