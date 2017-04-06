package sipparser

import (
	//"bytes"
	//"fmt"
	//"strings"
	"testing"
)

func TestSipMsgBody(t *testing.T) {
	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	prefix := FuncName()

	body, _ := NewSipMsgBody(context)
	body.headers.GenerateAndAddSingleHeader(context, "Content-Length", "123")
	body.headers.GenerateAndAddSingleHeader(context, "Content-Type", "application/sdp")
	body.headers.GenerateAndAddSingleHeader(context, "Content-Disposition", "render")

	body.headers.GenerateAndAddMultiHeader(context, "Content-Encoding", "gzip")
	body.headers.GenerateAndAddMultiHeader(context, "Content-Encoding", "tar")
	body.headers.GenerateAndAddMultiHeader(context, "Content-Language", "fr")
	body.headers.GenerateAndAddMultiHeader(context, "Content-Language", "en")

	body.headers.GenerateAndAddUnknownHeader(context, "Content-ABC", "xx1")
	body.headers.GenerateAndAddUnknownHeader(context, "Content-ABC", "xx2")

	body.body.SetString(context, "v=0\r\ns=call\r\n")

	encode := "Content-Length:         13\r\n"
	encode += "Content-Type: application/sdp\r\n"
	encode += "Content-Disposition: render\r\n"
	encode += "Content-Encoding: gzip, tar\r\n"
	encode += "Content-Language: fr, en\r\n"
	encode += "Content-ABC: xx1\r\n"
	encode += "Content-ABC: xx2\r\n"
	encode += "\r\n"
	encode += "v=0\r\ns=call\r\n"

	if body.String(context) != encode {
		t.Errorf("%s failed: encode = \n%s, wanted = \n%s\n", prefix, body.String(context), encode)
	}

}

func TestSipMsgBodies(t *testing.T) {

	context := NewParseContext()
	context.allocator = NewMemAllocator(1024 * 30)
	prefix := FuncName()

	bodies, _ := NewSipMsgBodies(context)

	body, addr := NewSipMsgBody(context)
	body.headers.GenerateAndAddSingleHeader(context, "Content-Length", "123")
	body.headers.GenerateAndAddSingleHeader(context, "Content-Type", "application/sdp")
	body.headers.GenerateAndAddSingleHeader(context, "Content-Disposition", "render")

	body.headers.GenerateAndAddMultiHeader(context, "Content-Encoding", "gzip")
	body.headers.GenerateAndAddMultiHeader(context, "Content-Encoding", "tar")

	body.headers.GenerateAndAddUnknownHeader(context, "Content-ABC", "xx1")
	body.headers.GenerateAndAddUnknownHeader(context, "Content-ABC", "xx2")

	body.body.SetString(context, "v=0\r\ns=call\r\n")
	bodies.AddBody(context, addr)

	body, addr = NewSipMsgBody(context)
	body.headers.GenerateAndAddSingleHeader(context, "Content-Length", "123")
	body.headers.GenerateAndAddSingleHeader(context, "Content-Type", "application/text")

	body.headers.GenerateAndAddMultiHeader(context, "Content-Language", "fr")
	body.headers.GenerateAndAddMultiHeader(context, "Content-Language", "en")

	body.headers.GenerateAndAddUnknownHeader(context, "Content-XYZ", "xx1")
	body.headers.GenerateAndAddUnknownHeader(context, "Content-XYZ", "xx2")

	body.body.SetString(context, "v=0\r\ns=call-2\r\n")
	bodies.AddBody(context, addr)

	boundary := "simple-boundary"
	encode := "--" + boundary + "\r\n"
	encode += "Content-Length:         13\r\n"
	encode += "Content-Type: application/sdp\r\n"
	encode += "Content-Disposition: render\r\n"
	encode += "Content-Encoding: gzip, tar\r\n"
	encode += "Content-ABC: xx1\r\n"
	encode += "Content-ABC: xx2\r\n"
	encode += "\r\n"
	encode += "v=0\r\ns=call\r\n"
	encode += "\r\n"
	encode += "--" + boundary + "\r\n"

	encode += "Content-Length:         15\r\n"
	encode += "Content-Type: application/text\r\n"
	encode += "Content-Language: fr, en\r\n"
	encode += "Content-XYZ: xx1\r\n"
	encode += "Content-XYZ: xx2\r\n"
	encode += "\r\n"
	encode += "v=0\r\ns=call-2\r\n"
	encode += "\r\n"
	encode += "--" + boundary + "--"

	if bodies.StringMulti(context, []byte(boundary)) != encode {
		t.Errorf("%s failed: encode = \n%s \nwanted = \n%s\n", prefix, bodies.StringMulti(context, []byte(boundary)), encode)
	}

	encode = "v=0\r\ns=call\r\n"

	if bodies.StringSingle(context) != encode {
		t.Errorf("%s failed: encode = \n%s \nwanted = \n%s\n", prefix, bodies.StringSingle(context), encode)
	}

}
