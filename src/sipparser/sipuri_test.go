package sipparser

import (
	"fmt"
	"testing"
)

func TestSipUriParseOK(t *testing.T) {

	type test_param struct {
		name  string
		value string
	}

	testdata := []struct {
		test     string
		user     string
		password string
		hostport string
		//params   []string
		//headers  []string
	}{
		{"sip:123@abc.com;user=phone", "123", "", "abc.com"},
		{"sip:123:tsdd@[1080::8:800:200c:417a]:5061", "123", "tsdd", "[1080::8:800:200c:417a]:5061"},
		{"sip:123:@10.43.12.14", "123", "", "10.43.12.14"},
	}

	for i, v := range testdata {
		uri := NewSipUri()

		newPos, err := uri.Parse([]byte(v.test), 0)
		if err != nil {
			t.Errorf("TestSipUriUserinfoParseOK[%d] failed, %s\n", i, err.Error())
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
	}
}

func TestSipUriUserinfoParseNOK(t *testing.T) {

	testdata := []struct {
		test   string
		newPos int
	}{
		{"sip:@abc.com", len("sip:")},
		{"sip::asas@abc.com", len("sip:")},
		{"sip:#123@abc.com", len("sip:")},
		{"sip:123:2#@abc.com", len("sip:123:2")},
	}

	for i, v := range testdata {
		uri := &SipUri{}

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
