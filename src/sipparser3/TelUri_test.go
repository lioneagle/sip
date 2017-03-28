package sipparser3

import (
	"bytes"
	"fmt"
	"testing"
)

//*
func TestTelUriParseOK(t *testing.T) {
	testdata := []struct {
		src                      string
		isGlobalNumber           bool
		number                   string
		phoneContext             string
		phoneContextIsDomainName bool
	}{
		{"tel:+861234", true, "+861234", "", false},
		{"tel:+86-12.(34)", true, "+861234", "", false},
		{"tel:861234;phone-context=+123", false, "861234", "+123", false},
		{"tel:861234;phone-context=+123", false, "861234", "+123", false},
		{"tel:861234;phone-context=a.com", false, "861234", "a.com", true},
		{"tel:86-1.2(34);phone-context=a.com", false, "861234", "a.com", true},
	}

	context := NewParseContext()

	for i, v := range testdata {
		uri := NewTelUri()

		newPos, err := uri.Parse(context, []byte(v.src), 0)
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

	context := NewParseContext()

	_, err := uri.Parse(context, []byte(src), 0)
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

func TestTelUriParseNOK(t *testing.T) {

	testdata := []struct {
		src    string
		newPos int
	}{
		{"tel1:+86123", 0},
		{"tel:+", len("tel:+")},
		{"tel:", len("tel:")},
		{"tel:.-()", len("tel:.-()")},
		{"tel:zz", len("tel:")},
		{"tel:123;", len("tel:123;")},
	}

	context := NewParseContext()

	for i, v := range testdata {
		uri := NewTelUri()

		newPos, err := uri.Parse(context, []byte(v.src), 0)
		if err == nil {
			t.Errorf("TestTelUriParseNOK[%d] failed", i)
			continue
		}

		if newPos != v.newPos {
			t.Errorf("TestTelUriParseNOK[%d] failed, newPos = %d, wanted = %d\n", i, newPos, v.newPos)
			continue
		}
	}
}

func TestTelUriEncode(t *testing.T) {

	testdata := []struct {
		src string
		dst string
	}{
		{"tel:+861234", "tel:+861234"},
		{"tel:+861234;phonex=+123", "tel:+861234;phonex=+123"},
		{"tel:861234;phone-context=+123", "tel:861234;phone-context=+123"},
		{"tel:861234;x1=5;y;phone-context=+1-2.3(56);zz", "tel:861234;phone-context=+12356;x1=5;y;zz"},
		{"tel:861234;x1=5;y;phone-context=abc.com;zz", "tel:861234;phone-context=abc.com;x1=5;y;zz"},
	}

	context := NewParseContext()

	for i, v := range testdata {
		uri := NewTelUri()

		_, err := uri.Parse(context, []byte(v.src), 0)
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

func TestTelUriEqual(t *testing.T) {
	testdata := []struct {
		uri1  string
		uri2  string
		equal bool
	}{
		{"tel:+86123", "tel:+8.6-1(2)3", true},
		{"tel:+86123;x1", "tel:+8.6-1(2)3;x1", true},
		{"tel:+86123;X2;x1", "tel:+8.6-1(2)3;X1;x2", true},
		{"tel:861234;x1=5;y;phone-context=abc.com;zz", "tel:861234;phone-context=abc.com;x1=5;y;zz", true},

		{"tel:+86123", "tel:8.6-1(2)3", false},
		{"tel:+86123", "tel:+18.6-1(2)3", false},
		{"tel:+86123;x1", "tel:+8.6-1(2)3", false},
		{"tel:+86123;x1=ab", "tel:+8.6-1(2)3;x1=cd", false},
		{"tel:861234;x1=5;y;phone-context=abc.com;zz", "tel:861234;phone-context=abcq.com;x1=5;y;zz", false},
	}

	context := NewParseContext()

	for i, v := range testdata {
		uri1 := NewTelUri()
		uri2 := NewTelUri()

		_, err := uri1.Parse(context, []byte(v.uri1), 0)
		if err != nil {
			t.Errorf("TestTelUriEqual[%d] failed, uri1 parse failed, err = %s\n", i, err.Error())
			continue
		}

		_, err = uri2.Parse(context, []byte(v.uri2), 0)
		if err != nil {
			t.Errorf("TestTelUriEqual[%d] failed, uri2 parse failed, err = %s\n", i, err.Error())
			continue
		}

		if v.equal && !uri1.Equal(uri2) {
			t.Errorf("TestTelUriEqual[%d] failed, should be equal, uri1 = %s, uri2 = %s\n", i, v.uri1, v.uri2)
			continue
		}

		if !v.equal && uri1.Equal(uri2) {
			t.Errorf("TestTelUriEqual[%d] failed, should not be equal, uri1 = %s, uri2 = %s\n", i, v.uri1, v.uri2)
			continue
		}
	}
} //*/

func BenchmarkTelUriParse(b *testing.B) {
	b.StopTimer()
	//v := []byte("tel:861234;x1=5;y;phone-context=abc.com;zz")
	//v := []byte("tel:861234;x1=5;phone-context=abc.com;zz")
	v := []byte("tel:861234;x1=5;phone-context=abc.com")
	context := NewParseContext()
	uri := NewTelUri()

	b.ReportAllocs()
	b.SetBytes(2)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		uri.Parse(context, v, 0)
	}
	//fmt.Printf("uri = %s\n", uri.String())
	fmt.Printf("")
}

/*
func BenchmarkTelUriString(b *testing.B) {
	b.StopTimer()
	v := "tel:861234;x1=5;y;phone-context=abc.com;zz"
	context := NewParseContext()
	uri := NewTelUri()
	uri.Parse(context, []byte(v), 0)
	b.ReportAllocs()
	b.SetBytes(2)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		uri.String()
	}
}//*/

func BenchmarkTelUriEncode(b *testing.B) {
	b.StopTimer()
	//v := []byte("tel:861234;x1=5;y;phone-context=abc.com;zz")
	v := []byte("tel:861234;x1=5;y;phone-context=abc.com")
	//v := []byte("tel:861234")
	context := NewParseContext()
	uri := NewTelUri()
	uri.Parse(context, v, 0)
	b.ReportAllocs()
	b.SetBytes(2)

	buf := bytes.NewBuffer(make([]byte, 1024*1024))
	//buf := &bytes.Buffer{}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		buf.Reset()
		uri.Encode(buf)
	}
}

//*/
