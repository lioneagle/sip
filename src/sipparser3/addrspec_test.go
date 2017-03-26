package sipparser3

import (
//"testing"
)

/*
func TestSipAddrSpecParseOK(t *testing.T) {

	testdata := []struct {
		src      string
		parseOk  bool
		isSipUri bool
		isTelUri bool
	}{
		{"sip:123@abc.com;ttl=10;user=phone;a;b;c;d;e?xx=yy&x1=aa", true, true, false},
		{"sips:123:tsdd@[1080::8:800:200c:417a]:5061", true, true, false},
		{"tel:861234;phone-context=+123", true, false, true},
		{"http://861234/phone-context=+123", false, false, false},
	}

	context := NewParseContext()

	for i, v := range testdata {
		addrsepc := NewSipAddrSpec()

		_, err := addrsepc.Parse(context, []byte(v.src), 0)
		if err != nil && v.parseOk {
			t.Errorf("TestSipAddrSpecParseOK[%d] failed, %s\n", i, err.Error())
			continue
		}

		if err == nil && !v.parseOk {
			t.Errorf("TestSipAddrSpecParseOK[%d] failed, parse failed wanted\n", i)
			continue
		}

		_, isSipUri := addrsepc.IsSipUri()
		_, isTelUri := addrsepc.IsTelUri()

		if v.isSipUri && !isSipUri {
			t.Errorf("TestSipAddrSpecParseOK[%d] failed, sip-uri wanted\n", i)
			continue
		}

		if v.isTelUri && !isTelUri {
			t.Errorf("TestSipAddrSpecParseOK[%d] failed, tel-uri wanted\n", i)
			continue
		}
	}
}
//*/

/*
func BenchmarkSipAddrSpecParse(b *testing.B) {
	b.StopTimer()
	//v := []byte("sip:biloxi.com;transport=tcp;method=REGISTER?to=sip:bob%40biloxi.com")
	//v := []byte("sip:abc@biloxi.com;transport=tcp")
	v := []byte("sip:abc@biloxi.com;transport=tcp;method=REGISTER")
	context := NewParseContext()
	b.ReportAllocs()
	b.SetBytes(2)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		addr := NewSipAddrSpec()
		addr.Parse(context, v, 0)
	}
}
//*/
