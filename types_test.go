package xmlrpc_test

import (
	"reflect"
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
	testNullableUnmarshal(t, "int/value", "<int>42</int>", xmlrpc.NullInt{Value: 42, Valid: true}, nil)
	testNullableUnmarshal(t, "int/null", "<nil/>", xmlrpc.NullInt{}, nil)
	testNullableUnmarshal(t, "int32/value", "<int>100</int>", xmlrpc.NullInt32{Value: 100, Valid: true}, nil)
	testNullableUnmarshal(t, "int32/null", "<nil/>", xmlrpc.NullInt32{}, nil)
	testNullableUnmarshal(t, "int64/value", "<i8>9999999999</i8>", xmlrpc.NullInt64{Value: 9999999999, Valid: true}, nil)
	testNullableUnmarshal(t, "int64/null", "<nil/>", xmlrpc.NullInt64{}, nil)
	testNullableUnmarshal(t, "float32/value", "<double>3.14</double>", xmlrpc.NullFloat32{Value: 3.14, Valid: true}, nil)
	testNullableUnmarshal(t, "float32/null", "<nil/>", xmlrpc.NullFloat32{}, nil)
	testNullableUnmarshal(t, "float64/value", "<double>2.71828</double>", xmlrpc.NullFloat64{Value: 2.71828, Valid: true}, nil)
	testNullableUnmarshal(t, "float64/null", "<nil/>", xmlrpc.NullFloat64{}, nil)
	testNullableUnmarshal(t, "string/value", "<string>hello</string>", xmlrpc.NullString{Value: "hello", Valid: true}, nil)
	testNullableUnmarshal(t, "string/null", "<nil/>", xmlrpc.NullString{}, nil)
	testNullableUnmarshal(t, "bool/true", "<boolean>1</boolean>", xmlrpc.NullBool{Value: true, Valid: true}, nil)
	testNullableUnmarshal(t, "bool/false", "<boolean>0</boolean>", xmlrpc.NullBool{Value: false, Valid: true}, nil)
	testNullableUnmarshal(t, "bool/null", "<nil/>", xmlrpc.NullBool{}, nil)
	testNullableUnmarshal(t, "time/value", "<dateTime.iso8601>19980717T14:08:55Z</dateTime.iso8601>",
		xmlrpc.NullTime{Value: time.Date(1998, 7, 17, 14, 8, 55, 0, time.UTC), Valid: true},
		func(got, want xmlrpc.NullTime) bool {
			return got.Valid == want.Valid && got.Value.Equal(want.Value)
		})
	testNullableUnmarshal(t, "time/null", "<nil/>", xmlrpc.NullTime{}, nil)
	testNullableUnmarshal(t, "bytes/value", "<base64>aGVsbG8gd29ybGQ=</base64>", xmlrpc.NullBytes{Value: []byte("hello world"), Valid: true}, nil)
	testNullableUnmarshal(t, "bytes/null", "<nil/>", xmlrpc.NullBytes{}, nil)
	testNullableUnmarshal(t, "struct/value",
		"<struct><member><name>message</name><value><string>hello</string></value></member></struct>",
		xmlrpc.NullStruct[Hoge]{Value: Hoge{Message: "hello"}, Valid: true}, nil)
	testNullableUnmarshal(t, "struct/null", "<nil/>", xmlrpc.NullStruct[Hoge]{}, nil)
}

func testNullableUnmarshal[T any](t *testing.T, name, valueXML string, want T, equal func(T, T) bool) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		var got T
		requestXML := "<methodCall><methodName>test</methodName><params><param><value>" +
			valueXML + "</value></param></params></methodCall>"
		reader, err := xmlrpc.NewRequestReader(strings.NewReader(requestXML))
		if err != nil {
			t.Fatalf("NewRequestReader failed: %v", err)
		}
		if err := reader.DecodeArgs(&got); err != nil {
			t.Fatalf("DecodeArgs failed: %v", err)
		}
		matches := reflect.DeepEqual(got, want)
		if equal != nil {
			matches = equal(got, want)
		}
		if !matches {
			t.Errorf("decoded value: got %#v, want %#v", got, want)
		}
	})
}
