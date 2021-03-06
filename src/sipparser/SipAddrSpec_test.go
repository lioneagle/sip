package sipparser

import (
	//"fmt"
	"testing"
)

//*
func TestSipAddrSpecParse(t *testing.T) {

	testdata := []struct {
		src    string
		ok     bool
		newPos int
		encode string
	}{
		{"sip:123@abc.com;ttl=10;user=phone;a;b;c;d;e?xx=yy&x1=aa", true, len("sip:123@abc.com;ttl=10;user=phone;a;b;c;d;e?xx=yy&x1=aa"), "sip:123@abc.com;ttl=10;user=phone;a;b;c;d;e?xx=yy&x1=aa"},
		{"sips:123:tsdd@[1080::8:800:200c:417a]:5061", true, len("sips:123:tsdd@[1080::8:800:200c:417a]:5061"), "sips:123:tsdd@[1080::8:800:200c:417a]:5061"},
		{"tel:861234;phone-context=+123", true, len("tel:861234;phone-context=+123"), "tel:861234;phone-context=+123"},

		{"httpx://861234/phone-context=+123", false, len("httpx:"), ""},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 300)
	prefix := FuncName()

	for i, v := range testdata {
		context.allocator.FreeAll()
		addr := NewSipAddrSpec(context)
		addrspec := addr.GetSipAddrSpec(context)
		newPos, err := addrspec.Parse(context, []byte(v.src), 0)

		if v.ok && err != nil {
			t.Errorf("%s[%d] failed: err = %s\n", prefix, i, err)
			continue
		}

		if !v.ok && err == nil {
			t.Errorf("%s[%d] failed: should parse failed", prefix, i)
			continue
		}

		if v.newPos != newPos {
			t.Errorf("%s[%d] failed: newPos = %d, wanted = %d\n", prefix, i, newPos, v.newPos)
		}

		if v.ok && v.encode != addrspec.String(context) {
			t.Errorf("%s[%d] failed: encode = %s, wanted = %s\n", prefix, i, addrspec.String(context), v.encode)
			continue
		}
	}
}

func TestSipAddrSpecParseWithouParam(t *testing.T) {

	testdata := []struct {
		src    string
		ok     bool
		newPos int
		encode string
	}{
		{"sip:123@abc.com;ttl=10;user=phone;a;b;c;d;e?xx=yy&x1=aa", true, len("sip:123@abc.com"), "sip:123@abc.com"},
		{"sips:123:tsdd@[1080::8:800:200c:417a]:5061", true, len("sips:123:tsdd@[1080::8:800:200c:417a]:5061"), "sips:123:tsdd@[1080::8:800:200c:417a]:5061"},
		{"tel:861234;phone-context=+123", true, len("tel:861234"), "tel:861234"},

		{"httpx://861234/phone-context=+123", false, len("httpx:"), ""},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	prefix := FuncName()

	for i, v := range testdata {
		context.allocator.FreeAll()
		addr := NewSipAddrSpec(context)
		addrspec := addr.GetSipAddrSpec(context)
		newPos, err := addrspec.ParseWithoutParam(context, []byte(v.src), 0)

		if v.ok && err != nil {
			t.Errorf("%s[%d] failed: err = %s\n", prefix, i, err)
			continue
		}

		if !v.ok && err == nil {
			t.Errorf("%s[%d] failed: should parse failed", prefix, i)
			continue
		}

		if v.newPos != newPos {
			t.Errorf("%s[%d] failed: newPos = %d, wanted = %d\n", prefix, i, newPos, v.newPos)
		}

		if v.ok && v.encode != addrspec.String(context) {
			t.Errorf("%s[%d] failed: encode = %s, wanted = %s\n", prefix, i, addrspec.String(context), v.encode)
			continue
		}
	}
}

func TestSipAddrSpecEqual(t *testing.T) {

	testdata := []struct {
		uri1  string
		uri2  string
		equal bool
	}{
		// Equal
		{"sip:abc.com", "sip:abc.com", true},
		{"sip:123abc@abc.com", "sip:123abc@aBC.com", true},
		{"sip:%61lice@atlanta.com;transport=TCP", "sip:alice@AtLanTa.CoM;Transport=tcp", true},
		{"sip:carol@chicago.com", "sip:carol@chicago.com;newparam=5", true},
		{"sip:carol@chicago.com;security=on", "sip:carol@chicago.com;newparam=5", true},
		{"sip:biloxi.com;transport=tcp;method=REGISTER?to=sip:bob%40biloxi.com", "sip:biloxi.com;method=REGISTER;transport=tcp?to=sip:bob%40biloxi.com", true},
		{"sip:alice@atlanta.com?subject=project%20x&priority=urgent", "sip:alice@atlanta.com?priority=urgent&subject=project%20x", true},

		{"tel:+86123", "tel:+8.6-1(2)3", true},
		{"tel:+86123;x1", "tel:+8.6-1(2)3;x1", true},
		{"tel:+86123;X2;x1", "tel:+8.6-1(2)3;X1;x2", true},
		{"tel:861234;x1=5;y;phone-context=abc.com;zz", "tel:861234;phone-context=abc.com;x1=5;y;zz", true},

		// Not equal
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

		{"tel:+86123", "tel:8.6-1(2)3", false},
		{"tel:+86123", "tel:+18.6-1(2)3", false},
		{"tel:+86123;x1", "tel:+8.6-1(2)3", false},
		{"tel:+86123;x1=ab", "tel:+8.6-1(2)3;x1=cd", false},
		{"tel:861234;x1=5;y;phone-context=abc.com;zz", "tel:861234;phone-context=abcq.com;x1=5;y;zz", false},
	}

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	prefix := FuncName()

	for i, v := range testdata {
		context.allocator.FreeAll()
		addr1 := NewSipAddrSpec(context)
		addr2 := NewSipAddrSpec(context)
		addrspec1 := addr1.GetSipAddrSpec(context)
		addrspec2 := addr2.GetSipAddrSpec(context)

		_, err := addrspec1.Parse(context, []byte(v.uri1), 0)
		if err != nil {
			t.Errorf("%s[%d] failed: uri1 parse failed: err = %s\n", prefix, i, err.Error())
			continue
		}

		_, err = addrspec2.Parse(context, []byte(v.uri2), 0)
		if err != nil {
			t.Errorf("%s[%d] failed: uri2 parse failed: err = %s\n", prefix, i, err.Error())
			continue
		}

		if v.equal && !addrspec1.Equal(context, addrspec2) {
			t.Errorf("%s[%d] failed: should be equal, uri1 = %s, uri2 = %s\n", prefix, i, v.uri1, v.uri2)
			continue
		}

		if !v.equal && addrspec1.Equal(context, addrspec2) {
			t.Errorf("%s[%d] failed: should not be equal, uri1 = %s, uri2 = %s\n", prefix, i, v.uri1, v.uri2)
			continue
		}
	}
}

//*/

func BenchmarkSipAddrSpecParse(b *testing.B) {
	b.StopTimer()
	//v := []byte("sip:biloxi.com;transport=tcp;method=REGISTER?to=sip:bob%40biloxi.com")
	//v := []byte("sip:abc@biloxi.com;transport=tcp")
	v := []byte("sip:abc@biloxi.com;transport=tcp;method=REGISTER")
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	addr := NewSipAddrSpec(context)
	addrspec := addr.GetSipAddrSpec(context)
	remain := context.allocator.Used()
	b.ReportAllocs()
	b.SetBytes(2)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		context.allocator.ClearAllocNum()
		context.allocator.FreePart(remain)
		addrspec.Parse(context, v, 0)
	}
}
