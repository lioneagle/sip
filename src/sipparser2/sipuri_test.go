package sipparser2

import (
	"bytes"
	"fmt"
	"testing"
)

//*
func TestSipUriParseOK(t *testing.T) {

	testdata := []struct {
		src       string
		user      string
		password  string
		hostport  string
		isSipsUri bool
		//params   []string
		//headers  []string
	}{
		{"sip:123@abc.com;ttl=10;user=phone;a;b;c;d;e?xx=yy&x1=aa", "123", "", "abc.com", false},
		{"Sip:123:tsdd@[1080::8:800:200c:417a]:5061", "123", "tsdd", "[1080::8:800:200c:417a]:5061", false},
		{"sip:123:@10.43.12.14", "123", "", "10.43.12.14", false},
		{"sip:%23123%31:@10.43.12.14", "#1231", "", "10.43.12.14", false},

		{"sips:123@abc.com;ttl=10;user=phone;a;b;c;d;e?xx=yy&x1=aa", "123", "", "abc.com", true},
		{"sips:123:tsdd@[1080::8:800:200c:417a]:5061", "123", "tsdd", "[1080::8:800:200c:417a]:5061", true},
		{"sips:123:@10.43.12.14", "123", "", "10.43.12.14", true},
		{"sips:%23123%31:@10.43.12.14", "#1231", "", "10.43.12.14", true},
	}

	for i, v := range testdata {
		uri := NewSipUri()

		newPos, err := uri.Parse([]byte(v.src), 0)
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

		if newPos != len(v.src) {
			t.Errorf("TestSipUriUserinfoParseOK[%d] failed, newPos = %d, wanted = %d\n", i, newPos, len(v.src))
			continue
		}

		if uri.user.String() != v.user {
			t.Errorf("TestSipUriUserinfoParseOK[%d] failed, user wrong, user = %s, wanted = %s", i, string(uri.user.value), v.user)
			continue
		}

		if uri.password.String() != v.password {
			t.Errorf("TestSipUriUserinfoParseOK[%d] failed, password wrong, password = %s, wanted = %s", i, string(uri.password.value), v.password)
			continue
		}

		if uri.hostport.String() != v.hostport {
			t.Errorf("TestSipUriUserinfoParseOK[%d] failed, host wrong, host = %s, wanted = %s", i, uri.hostport, v.hostport)
			continue
		}
	}
}

func TestSipUriParamsParseOK(t *testing.T) {

	uri := NewSipUri()
	src := "sip:123@abc.com;ttl=10;user%32=phone%31;a;b;c;d;e?xx=yy&x1=aa"

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
			t.Errorf("TestSipUriParamsParseOK[%d] failed, cannot get %s param\n", i, v.name)
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
		src    string
		newPos int
	}{
		//{"sipx:@abc.com", len("sipx:")},
		{"sipx:@abc.com", 0},
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

		newPos, err := uri.Parse([]byte(v.src), 0)
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
		{"sip:123@10.43.12.14;method=INVITE", "sip:123@10.43.12.14;method=INVITE"},
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

func TestSipUriEqual(t *testing.T) {
	testdata := []struct {
		uri1  string
		uri2  string
		equal bool
	}{
		{"sip:abc.com", "sip:abc.com", true},
		{"sip:123abc@abc.com", "sip:123abc@aBC.com", true},
		{"sip:%61lice@atlanta.com;transport=TCP", "sip:alice@AtLanTa.CoM;Transport=tcp", true},
		{"sip:carol@chicago.com", "sip:carol@chicago.com;newparam=5", true},
		{"sip:carol@chicago.com;security=on", "sip:carol@chicago.com;newparam=5", true},
		{"sip:biloxi.com;transport=tcp;method=REGISTER?to=sip:bob%40biloxi.com", "sip:biloxi.com;method=REGISTER;transport=tcp?to=sip:bob%40biloxi.com", true},
		{"sip:alice@atlanta.com?subject=project%20x&priority=urgent", "sip:alice@atlanta.com?priority=urgent&subject=project%20x", true},

		{"sip:123abc@abc.com", "sips:123abc@abc.com", false},                                  // different scheme
		{"SIP:ALICE@AtLanTa.CoM;Transport=udp", "sip:alice@AtLanTa.CoM;Transport=UDP", false}, //different usernames
		{"SIP:ALICE@AtLanTa.CoM;Transport=udp", "sip:AtLanTa.CoM;Transport=UDP", false},       //different usernames
		{"sip:bob@biloxi.com", "sip:bob@biloxi.com:5060", false},                              //can resolve to different ports
		{"sip:bob@biloxi.com", "sip:bob@biloxi.com:6000;transport=tcp", false},                //can resolve to different port and transports
		{"sip:abc.com;user=phone", "sip:abc.com;user=ip", false},                              //different param
		{"sip:abc.com;user=phone", "sip:abc.com", false},                                      //different param
		{"sip:abc.com;ttl=11", "sip:abc.com;ttl=10", false},                                   //different param
		{"sip:abc.com", "sip:abc.com;ttl=10", false},                                          //different param
		{"sip:abc.com", "sip:abc.com;method=INVITE", false},                                   //different param
		{"sip:carol@chicago.com", "sip:carol@chicago.com?Subject=next%20meeting", false},      //different header component
		{"sip:bob@phone21.boxesbybob.com", "sip:bob@192.0.2.4", false},                        //even though that's what phone21.boxesbybob.com resolves to
	}

	for i, v := range testdata {
		uri1 := NewSipUri()
		uri2 := NewSipUri()

		_, err := uri1.Parse([]byte(v.uri1), 0)
		if err != nil {
			t.Errorf("TestSipUriEqual[%d] failed, uri1 parse failed, err = %s\n", i, err.Error())
			continue
		}

		_, err = uri2.Parse([]byte(v.uri2), 0)
		if err != nil {
			t.Errorf("TestSipUriEqual[%d] failed, uri2 parse failed, err = %s\n", i, err.Error())
			continue
		}

		if v.equal && !uri1.Equal(uri2) {
			t.Errorf("TestSipUriEqual[%d] failed, should be equal, uri1 = %s, uri2 = %s\n", i, v.uri1, v.uri2)
			continue
		}

		if !v.equal && uri1.Equal(uri2) {
			t.Errorf("TestSipUriEqual[%d] failed, should not be equal, uri1 = %s, uri2 = %s\n", i, v.uri1, v.uri2)
			continue
		}
	}
}

//*/

/*
func BenchmarkStrParse(b *testing.B) {
	b.StopTimer()
	//v := []byte("sip:biloxi.com;transport=tcp;method=REGISTER?to=sip:bob%40biloxi.com")
	v := []byte("sip:biloxi.com")
	b.SetBytes(2)
	b.ReportAllocs()

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for k := 0; k < len(v); k++ {
			if IsWspChar(v[k]) {
				//fmt.Printf("k=%d, len(v)=%d\n", k, len(v))
				break
			}
		}
	}
	fmt.Printf("")
}//*/

func BenchmarkSipUriParse(b *testing.B) {
	b.StopTimer()
	//v := []byte("sip:biloxi.com;transport=tcp;method=REGISTER?to=sip:bob%40biloxi.com")
	//v := []byte("sip:biloxi.com")
	//v := []byte("sip:abc@biloxi.com;transport=tcp")
	v := []byte("sip:abc@biloxi.com;transport=tcp;method=REGISTER")
	uri := NewSipUri()

	b.ReportAllocs()
	b.SetBytes(2)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		uri.Init()
		uri.Parse(v, 0)
	}
	//fmt.Printf("uri = %s\n", uri.String())
	fmt.Printf("")
}

/*
func BenchmarkSipUriString(b *testing.B) {
	b.StopTimer()
	v := "sip:biloxi.com;transport=tcp;method=REGISTER?to=sip:bob%40biloxi.com"
	uri := NewSipUri()
	uri.Parse([]byte(v), 0)
	b.ReportAllocs()
	b.SetBytes(2)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		uri.String()
	}
}//*/
//*
func BenchmarkSipUriEncode(b *testing.B) {
	b.StopTimer()
	//v := []byte("sip:biloxi.com;transport=tcp;method=REGISTER?to=sip:bob%40biloxi.com")
	v := []byte("sip:abc@biloxi.com;transport=tcp;method=REGISTER")
	uri := NewSipUri()
	uri.Parse(v, 0)
	b.SetBytes(2)
	b.ReportAllocs()

	buf := bytes.NewBuffer(make([]byte, 1024*1024))
	//buf := &bytes.Buffer{}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		buf.Reset()
		uri.Encode(buf)
	}
}

//*/
