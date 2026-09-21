package xmlrpc_test

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"log"
	"math"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/mamiya312/go-xmlrpc"
)

func TestValueScanAsInt(t *testing.T) {
	testcases := []struct {
		XML    string
		Expect int
	}{
		{
			XML: "<param><int>0</int></param>", Expect: 0,
		},
		{
			XML: "<param><int>123</int></param>", Expect: 123,
		},
		{
			XML: "<param><int>-123</int></param>", Expect: -123,
		},
		{
			XML: fmt.Sprintf("<param><int>%d</int></param>", math.MinInt), Expect: math.MinInt,
		},
		{
			XML: fmt.Sprintf("<param><int>%d</int></param>", math.MaxInt), Expect: math.MaxInt,
		},
		{
			XML: "<param><i4>0</i4></param>", Expect: 0,
		},
		{
			XML: "<param><i4>123</i4></param>", Expect: 123,
		},
		{
			XML: "<param><i4>-123</i4></param>", Expect: -123,
		},
		{
			XML: fmt.Sprintf("<param><i4>%d</i4></param>", math.MinInt), Expect: math.MinInt,
		},
		{
			XML: fmt.Sprintf("<param><i4>%d</i4></param>", math.MaxInt), Expect: math.MaxInt,
		},
	}
	for _, testcase := range testcases {
		var v int
		var value xmlrpc.XMLRPCValue
		if err := xml.Unmarshal([]byte(testcase.XML), &value); err != nil {
			t.Fatal(err)
		}
		field := reflect.ValueOf(&v).Elem()
		if err := value.Scan(&field); err != nil {
			t.Fatal(err)
		}
		if testcase.Expect != v {
			t.Errorf("expect:%d => %d", testcase.Expect, v)
		}
	}
}

func TestValueScanAsInt32(t *testing.T) {
	testcases := []struct {
		XML    string
		Expect int32
	}{
		{
			XML:    "<param><int>0</int></param>",
			Expect: 0,
		},
		{
			XML:    "<param><int>123</int></param>",
			Expect: 123,
		},
		{
			XML:    "<param><int>-123</int></param>",
			Expect: -123,
		},
		{
			XML:    fmt.Sprintf("<param><int>%d</int></param>", math.MinInt32),
			Expect: math.MinInt32,
		},
		{
			XML:    fmt.Sprintf("<param><int>%d</int></param>", math.MaxInt32),
			Expect: math.MaxInt32,
		},
		{
			XML:    "<param><i4>0</i4></param>",
			Expect: 0,
		},
		{
			XML:    "<param><i4>123</i4></param>",
			Expect: 123,
		},
		{
			XML:    "<param><i4>-123</i4></param>",
			Expect: -123,
		},
		{
			XML:    fmt.Sprintf("<param><i4>%d</i4></param>", math.MinInt32),
			Expect: math.MinInt32,
		},
		{
			XML:    fmt.Sprintf("<param><i4>%d</i4></param>", math.MaxInt32),
			Expect: math.MaxInt32,
		},
	}
	for _, testcase := range testcases {
		var v int32
		var value xmlrpc.XMLRPCValue
		if err := xml.Unmarshal([]byte(testcase.XML), &value); err != nil {
			t.Fatal(err)
		}
		field := reflect.ValueOf(&v).Elem()
		if err := value.Scan(&field); err != nil {
			t.Fatal(err)
		}
		if testcase.Expect != v {
			t.Errorf("expect:%d => %d", testcase.Expect, v)
		}
	}
}

func TestValueScanAsInt64(t *testing.T) {
	testcases := []struct {
		XML    string
		Expect int64
	}{
		{
			XML:    "<param><int>0</int></param>",
			Expect: 0,
		},
		{
			XML:    "<param><int>123</int></param>",
			Expect: 123,
		},
		{
			XML:    "<param><int>-123</int></param>",
			Expect: -123,
		},
		{
			XML:    fmt.Sprintf("<param><int>%d</int></param>", math.MinInt64),
			Expect: math.MinInt64,
		},
		{
			XML:    fmt.Sprintf("<param><int>%d</int></param>", math.MaxInt64),
			Expect: math.MaxInt64,
		},
		{
			XML:    "<param><i4>0</i4></param>",
			Expect: 0,
		},
		{
			XML:    "<param><i4>123</i4></param>",
			Expect: 123,
		},
		{
			XML:    "<param><i4>-123</i4></param>",
			Expect: -123,
		},
		{
			XML:    fmt.Sprintf("<param><i4>%d</i4></param>", math.MinInt64),
			Expect: math.MinInt64,
		},
		{
			XML:    fmt.Sprintf("<param><i4>%d</i4></param>", math.MaxInt64),
			Expect: math.MaxInt64,
		},
	}
	for _, testcase := range testcases {
		var v int64
		var value xmlrpc.XMLRPCValue
		if err := xml.Unmarshal([]byte(testcase.XML), &value); err != nil {
			t.Fatal(err)
		}
		field := reflect.ValueOf(&v).Elem()
		if err := value.Scan(&field); err != nil {
			t.Fatal(err)
		}
		if testcase.Expect != v {
			t.Errorf("expect:%d => %d", testcase.Expect, v)
		}
	}
}

func TestValueScanAsDateTime(t *testing.T) {
	testcases := []struct {
		XML    string
		Expect time.Time
	}{
		{
			XML:    "<param><dateTime.iso8601>19980717T14:08:55Z</dateTime.iso8601></param>",
			Expect: time.Date(1998, 7, 17, 14, 8, 55, 0, time.UTC).Local(),
		},
	}
	for _, testcase := range testcases {
		var v time.Time
		var value xmlrpc.XMLRPCValue
		if err := xml.Unmarshal([]byte(testcase.XML), &value); err != nil {
			t.Fatal(err)
		}
		field := reflect.ValueOf(&v).Elem()
		if err := value.Scan(&field); err != nil {
			t.Fatal(err)
		}
		if testcase.Expect != v {
			t.Errorf("expect:%v => %v", testcase.Expect, v)
		}
	}
}

type SubStructXml2Rpc struct {
	Foo  int
	Bar  string
	Data []int
}

type StructXml2Rpc struct {
	Int    int
	Float  float64
	Str    string
	Bool   bool
	Sub    SubStructXml2Rpc
	Time   time.Time
	Base64 []byte
}

func TestRequest(t *testing.T) {
	xml := "<methodCall><methodName>Some.Method</methodName><params><param><value><i4>123</i4></value></param><param><value><double>3.145926</double></value></param><param><value><string>Hello, World!</string></value></param><param><value><boolean>0</boolean></value></param><param><value><struct><member><name>Foo</name><value><int>42</int></value></member><member><name>Bar</name><value><string>I'm Bar</string></value></member><member><name>Data</name><value><array><data><value><int>1</int></value><value><int>2</int></value><value><int>3</int></value></data></array></value></member></struct></value></param><param><value><base64>eW91IGNhbid0IHJlYWQgdGhpcyE=</base64></value></param></params></methodCall>"
	r, err := http.NewRequest(http.MethodPost, "http://127.0.0.1/RPC2", strings.NewReader(xml))
	if err != nil {
		t.Fatal(err.Error())
	}
	request, err := xmlrpc.NewRequestReader(r.Body)
	if err != nil {
		t.Fatal(err.Error())
	}
	var i int
	var f float32
	var s string
	var b bool
	var sub SubStructXml2Rpc
	var b64 []byte
	if err := request.DecodeArgs(&i, &f, &s, &b, &sub, &b64); err != nil {
		t.Fatal(err.Error())
	}
	log.Printf("%d %g %s %v %+v %s", i, f, s, b, sub, base64.StdEncoding.EncodeToString(b64))
}

func TestResponse(t *testing.T) {
	xml := "<methodResponse><params><param><value><i4>123</i4></value></param><param><value><double>3.145926</double></value></param><param><value><string>Hello, World!</string></value></param><param><value><boolean>0</boolean></value></param><param><value><struct><member><name>Foo</name><value><int>42</int></value></member><member><name>Bar</name><value><string>I'm Bar</string></value></member><member><name>Data</name><value><array><data><value><int>1</int></value><value><int>2</int></value><value><int>3</int></value></data></array></value></member></struct></value></param><param><value><base64>eW91IGNhbid0IHJlYWQgdGhpcyE=</base64></value></param></params></methodResponse>"
	request, err := xmlrpc.NewResponseReader(strings.NewReader(xml))
	if err != nil {
		t.Fatal(err.Error())
	}
	var i int
	var f float32
	var s string
	var b bool
	var sub SubStructXml2Rpc
	var b64 []byte
	if err := request.Decode(&i, &f, &s, &b, &sub, &b64); err != nil {
		t.Fatal(err.Error())
	}
	log.Printf("%d %g %s %v %+v %s", i, f, s, b, sub, base64.StdEncoding.EncodeToString(b64))
}

func TestToXML(t *testing.T) {
	reply := StructXml2Rpc{
		Int:   1,
		Float: 0.2,
		Str:   "3",
		Bool:  false,
		Sub: SubStructXml2Rpc{
			Foo:  10,
			Bar:  "あいうえお",
			Data: []int{2, 3, 4},
		},
		Time:   time.Now(),
		Base64: []byte("you can't read this!"),
	}

	if s, err := xmlrpc.ToRequestXML("aaa", reply); err != nil {
		t.Error(err.Error())
	} else {
		log.Printf("%s", s)
	}
	if s, err := xmlrpc.ToResponseXML(reply); err != nil {
		t.Error(err.Error())
	} else {
		log.Printf("%s", s)
	}
}
