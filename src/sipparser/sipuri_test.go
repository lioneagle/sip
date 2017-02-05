package sipparser

import (
	"fmt"
	"testing"
)

func TestSipUriParseOK(t *testing.T) {

	testdata := []struct {
		test      string
		user      string
		password  string
		hostport  string
		isSipsUri bool
		//params   []string
		//headers  []string
	}{
		{"sip:123@abc.com;ttl=10;user=phone;a;b;c;d;e?xx=yy&x1=aa", "123", "", "abc.com", false},
		{"sip:123:tsdd@[1080::8:800:200c:417a]:5061", "123", "tsdd", "[1080::8:800:200c:417a]:5061", false},
		{"sip:123:@10.43.12.14", "123", "", "10.43.12.14", false},
		{"sip:%23123%31:@10.43.12.14", "#1231", "", "10.43.12.14", false},

		{"sips:123@abc.com;ttl=10;user=phone;a;b;c;d;e?xx=yy&x1=aa", "123", "", "abc.com", true},
		{"sips:123:tsdd@[1080::8:800:200c:417a]:5061", "123", "tsdd", "[1080::8:800:200c:417a]:5061", true},
		{"sips:123:@10.43.12.14", "123", "", "10.43.12.14", true},
		{"sips:%23123%31:@10.43.12.14", "#1231", "", "10.43.12.14", true},
	}

	for i, v := range testdata {
		uri := NewSipUri()

		newPos, err := uri.Parse([]byte(v.test), 0)
		if err != nil {
			t.Errorf("TestSipUriUserinfoParseOK[%d] failed, %s\n", i, err.Error())
			continue
		}

		if v.isSipsUri && !uri.IsSipsUri() {
			t.Errorf("TestSipUriUserinfoParseOK[%d] failed, sips-uri wanted\n", i)
			continue
		}

		if !v.isSipsUri && !uri.IsSipUri() {
			t.Errorf("TestSipUriUserinfoParseOK[%d] failed, sip-uri wanted\n", i)
			continue
		}

		if newPos != len(v.test) {
			t.Errorf("TestSipUriUserinfoParseOK[%d] failed, newPos = %d, wanted = %d\n", i, newPos, len(v.test))
			continue
		}

		if v.user != string(uri.user.value) {
			t.Errorf("TestSipUriUserinfoParseOK[%d] failed, user wrong, user = %s, wanted = %s", i, string(uri.user.value), v.user)
			continue
		}

		if v.password != string(uri.password.value) {
			t.Errorf("TestSipUriUserinfoParseOK[%d] failed, password wrong, password = %s, wanted = %s", i, string(uri.password.value), v.password)
			continue
		}

		if v.hostport != uri.hostport.String() {
			t.Errorf("TestSipUriUserinfoParseOK[%d] failed, host wrong, host = %s, wanted = %s", i, uri.hostport, v.hostport)
			continue
		}

		//fmt.Printf("uri encode = %s\n", uri)
	}
}

func TestSipUriParamsParseOK(t *testing.T) {

	uri := NewSipUri()
	src := "sip:123@abc.com;ttl=10;user=phone;a;b;c;d;e?xx=yy&x1=aa"

	_, err := uri.Parse([]byte(src), 0)
	if err != nil {
		t.Errorf("TestSipUriParamsParseOK failed, err = %s\n", err.Error())
		return
	}

	testdata := []struct {
		name     string
		value    string
		hasValue bool
	}{
		{"Ttl", "10", true},
		{"UseR", "phone", true},
		{"A", "", false},
		{"b", "", false},
		{"c", "", false},
		{"D", "", false},
		{"E", "", false},
	}

	for i, v := range testdata {
		param, ok := uri.params.GetParam(v.name)
		if !ok {
			t.Errorf("TestSipUriParamsParseOK[%d] failed, cannot get ttl param\n", i)
			continue
		}

		if param.value.Exist() && !v.hasValue {
			t.Errorf("TestSipUriParamsParseOK[%d] failed, should have no pvalue\n", i)
			continue
		}

		if !param.value.Exist() && v.hasValue {
			t.Errorf("TestSipUriParamsParseOK[%d] failed, should have pvalue\n", i)
			continue
		}

		if param.value.Exist() && param.value.String() != v.value {
			t.Errorf("TestSipUriParamsParseOK[%d] failed, pvalue = %s, wanted = %s\n", i, param.value.String(), v.value)
			continue
		}

	}
}

func TestSipUriHeadersParseOK(t *testing.T) {

	uri := NewSipUri()
	src := "sip:123@abc.com;ttl=10;user=phone;a;b;c;d;e?xx=yy&x1=aa"

	_, err := uri.Parse([]byte(src), 0)
	if err != nil {
		t.Errorf("TestSipUriHeadersParseOK failed, err = %s\n", err.Error())
		return
	}

	testdata := []struct {
		name     string
		value    string
		hasValue bool
	}{
		{"XX", "yy", true},
		{"x1", "aa", true},
	}

	for i, v := range testdata {
		header, ok := uri.headers.GetHeader(v.name)
		if !ok {
			t.Errorf("TestSipUriHeadersParseOK[%d] failed, cannot get ttl param\n", i)
			continue
		}

		if header.value.Exist() && !v.hasValue {
			t.Errorf("TestSipUriHeadersParseOK[%d] failed, should have no pvalue\n", i)
			continue
		}

		if !header.value.Exist() && v.hasValue {
			t.Errorf("TestSipUriHeadersParseOK[%d] failed, should have pvalue\n", i)
			continue
		}

		if header.value.Exist() && header.value.String() != v.value {
			t.Errorf("TestSipUriHeadersParseOK[%d] failed, pvalue = %s, wanted = %s\n", i, header.value.String(), v.value)
			continue
		}
	}
}

func TestSipUriUserinfoParseNOK(t *testing.T) {

	testdata := []struct {
		test   string
		newPos int
	}{
		{"sipx:@abc.com", len("sipx")},
		{"Sip:@abc.com", len("Sip")},
		{"sip:@abc.com", len("sip:")},
		{"sip::asas@abc.com", len("sip:")},
		{"sip:#123@abc.com", len("sip:")},
		{"sip:123:2#@abc.com", len("sip:123:2")},
		{"sip:123:2@abc.com;;", len("sip:123:2@abc.com;")},
		{"sip:123:2@abc.com;a=;", len("sip:123:2@abc.com;a=")},
		{"sip:123:2@abc.com;ttl=10?q", len("sip:123:2@abc.com;ttl=10?q")},
		{"sip:123:2@abc.com;a=b?", len("sip:123:2@abc.com;a=b?")},
	}

	for i, v := range testdata {
		uri := NewSipUri()

		newPos, err := uri.Parse([]byte(v.test), 0)
		if err == nil {
			t.Errorf("TestSipUriUserinfoParseNOK[%d] failed", i)
			continue
		}

		if newPos != v.newPos {
			t.Errorf("TestSipUriUserinfoParseNOK[%d] failed, newPos = %d, wanted = %d\n", i, newPos, v.newPos)
			continue
		}
	}
}

func TestSipUriEncode(t *testing.T) {

	testdata := []struct {
		src string
		dst string
	}{
		{"sip:123@abc.com;ttl=10;user=phone;a;b;c;d;e?xx=yy&x1=aa", "sip:123@abc.com;ttl=10;user=phone;a;b;c;d;e?xx=yy&x1=aa"},
		{"sip:123:tsdd@[1080::8:800:200c:417a]:5061", "sip:123:tsdd@[1080::8:800:200c:417a]:5061"},
		{"sip:123:@10.43.12.14", "sip:123:@10.43.12.14"},
		{"sip:%23123%31:@10.43.12.14", "sip:%231231:@10.43.12.14"},
	}

	for i, v := range testdata {
		uri := NewSipUri()

		_, err := uri.Parse([]byte(v.src), 0)
		if err != nil {
			t.Errorf("TestSipUriEncode[%d] failed, parse failed, err = %s\n", i, err.Error())
			continue
		}

		str := uri.String()

		if str != v.dst {
			t.Errorf("TestSipUriUserinfoParseOK[%d] failed, uri = %s, wanted = %s\n", i, str, v.dst)
			continue
		}
	}
}

func BenchmarkSipUri(b *testing.B) {

	for i := 0; i < 10000; i++ {
		fmt.Sprintf("hello")
	}

}

func BenchmarkSipUri2(b *testing.B) {

	for i := 0; i < 1000000; i++ {
		fmt.Sprintf("hello")
	}

}
