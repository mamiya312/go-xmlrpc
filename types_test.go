package xmlrpc_test

import (
	"strings"
	"testing"
	"time"

	"github.com/mamiya312/go-xmlrpc"
)

type Hoge struct {
	Message string `xmlrpc:"message"`
}

func TestNullableTypes(t *testing.T) {

	testcases := []struct {
		Value  xmlrpc.Marshaler
		Expect string
	}{
		{Value: &xmlrpc.NullBool{Value: true, Valid: true}, Expect: "<value><boolean>1</boolean></value>"},
		{Value: &xmlrpc.NullBool{Value: false, Valid: true}, Expect: "<value><boolean>0</boolean></value>"},
		{Value: &xmlrpc.NullBool{Value: true, Valid: false}, Expect: "<value><nil/></value>"},
		{Value: &xmlrpc.NullBool{Value: false, Valid: false}, Expect: "<value><nil/></value>"},
		{Value: &xmlrpc.NullInt{Value: 1, Valid: true}, Expect: "<value><int>1</int></value>"},
		{Value: &xmlrpc.NullInt{Value: 1, Valid: false}, Expect: "<value><nil/></value>"},
		{Value: &xmlrpc.NullInt32{Value: 1, Valid: true}, Expect: "<value><int>1</int></value>"},
		{Value: &xmlrpc.NullInt32{Value: 1, Valid: false}, Expect: "<value><nil/></value>"},
		{Value: &xmlrpc.NullInt64{Value: 1, Valid: true}, Expect: "<value><i8>1</i8></value>"},
		{Value: &xmlrpc.NullInt64{Value: 1, Valid: false}, Expect: "<value><nil/></value>"},
		{Value: &xmlrpc.NullFloat32{Value: -12.53, Valid: true}, Expect: "<value><double>-12.53</double></value>"},
		{Value: &xmlrpc.NullFloat32{Value: -12.53, Valid: false}, Expect: "<value><nil/></value>"},
		{Value: &xmlrpc.NullFloat64{Value: -12.53, Valid: true}, Expect: "<value><double>-12.53</double></value>"},
		{Value: &xmlrpc.NullFloat64{Value: -12.53, Valid: false}, Expect: "<value><nil/></value>"},
		{Value: &xmlrpc.NullString{Value: "hoge", Valid: true}, Expect: "<value><string>hoge</string></value>"},
		{Value: &xmlrpc.NullString{Value: "hoge", Valid: false}, Expect: "<value><nil/></value>"},
		{Value: &xmlrpc.NullTime{Value: time.Date(1998, 7, 17, 14, 8, 55, 0, time.UTC).Local(), Valid: true}, Expect: "<value><dateTime.iso8601>19980717T14:08:55Z</dateTime.iso8601></value>"},
		{Value: &xmlrpc.NullTime{Value: time.Date(1998, 7, 17, 14, 8, 55, 0, time.UTC).Local(), Valid: false}, Expect: "<value><nil/></value>"},
		{Value: &xmlrpc.NullBytes{Value: []byte("you can't read this!"), Valid: true}, Expect: "<value><base64>eW91IGNhbid0IHJlYWQgdGhpcyE=</base64></value>"},
		{Value: &xmlrpc.NullBytes{Value: []byte("you can't read this!"), Valid: false}, Expect: "<value><nil/></value>"},
		{Value: &xmlrpc.NullStruct[Hoge]{Value: Hoge{}, Valid: true}, Expect: "<value><struct><member><name>message</name><value><string></string></value></member></struct></value>"},
		{Value: &xmlrpc.NullStruct[Hoge]{Value: Hoge{}, Valid: false}, Expect: "<value><nil/></value>"},
	}

	for _, testcase := range testcases {
		s, err := testcase.Value.MarshalXMLRPC()
		if err != nil {
			t.Errorf("testcase: %#v => err: %s", testcase, err.Error())
		}
		if testcase.Expect != s {
			t.Errorf("testcase: %#v => %s", testcase, s)
		}
	}

}

func TestNullableTypesUnmarshal(t *testing.T) {
	var present xmlrpc.NullString
	reader, err := xmlrpc.NewRequestReader(strings.NewReader(
		"<methodCall><methodName>test</methodName><params><param><value><string>hello</string></value></param></params></methodCall>",
	))
	if err != nil {
		t.Fatal(err)
	}
	if err := reader.DecodeArgs(&present); err != nil {
		t.Fatal(err)
	}
	if !present.Valid || present.Value != "hello" {
		t.Fatalf("unexpected present value: %+v", present)
	}

	var nullValue xmlrpc.NullString
	reader, err = xmlrpc.NewRequestReader(strings.NewReader(
		"<methodCall><methodName>test</methodName><params><param><value><nil/></value></param></params></methodCall>",
	))
	if err != nil {
		t.Fatal(err)
	}
	if err := reader.DecodeArgs(&nullValue); err != nil {
		t.Fatal(err)
	}
	if nullValue.Valid {
		t.Fatalf("expected null value: %+v", nullValue)
	}
}
