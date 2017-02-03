package sipparser

import (
	"fmt"
	"testing"
)

func TestSipUriUserinfoParseOK(t *testing.T) {

	testdata := []struct {
		test     string
		newPos   int
		user     string
		password string
	}{
		{"sip:123@abc.com", len("sip:123@abc.com"), "123", ""},
		{"sip:123:tsdd@abc.com", len("sip:123:tsdd@abc.com"), "123", "tsdd"},
	}

	for i, v := range testdata {
		uri := &SipUri{}

		newPos, err := uri.Parse([]byte(v.test), 0)
		if err != nil {
			t.Errorf("TestSipUriUserinfoParseOK[%d] failed, %s\n", i, err.Error())
			continue
		}

		if newPos != v.newPos {
			t.Errorf("TestSipUriUserinfoParseOK[%d] failed, newPos = %d, wanted = %d\n", i, newPos, v.newPos)
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
