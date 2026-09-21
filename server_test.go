package xmlrpc_test

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/mamiya312/go-xmlrpc"
)

type EchoArgs struct {
	Message string
}
type EchoReplies struct {
	Message string
}

type testLogger struct {
	buffer bytes.Buffer
}

func (l *testLogger) Printf(format string, v ...any) {
	fmt.Fprintf(&l.buffer, format, v...)
	l.buffer.WriteByte('\n')
}

func Echo(args *EchoArgs, reply *EchoReplies) error {
	reply.Message = fmt.Sprintf("%s %s", args.Message, "fugafuga")
	return nil
}

func TestHoge(t *testing.T) {

	mValue := reflect.ValueOf(Echo)
	mType := mValue.Type()
	log.Printf("%d\n", mType.NumMethod())
	mArgsType := mType.In(0).Elem()
	mRepliesType := mType.In(1).Elem()

	args := reflect.New(mArgsType)

	xml := "<methodCall><methodName>Some.Method</methodName><params><param><value><string>hogehoge</string></value></param></params></methodCall>"

	reader, err := xmlrpc.NewRequestReader(strings.NewReader(xml))
	if err != nil {
		t.Fatalf("%+v", err)
	}

	if m, err := reader.Method(); err != nil {
		t.Fatalf("%+v", err)
	} else {
		log.Println(m)
	}
	if err := reader.Decode(args.Interface()); err != nil {
		t.Fatalf("%+v", err)
	}
	log.Printf("%+v\n", args)

	reply := reflect.New(mRepliesType)
	r := mValue.Call([]reflect.Value{
		args,
		reply,
	})
	if len(r) != 1 {
		t.Fatalf("len(r) != 1, %+v", r)
	}

	errInterface := r[0].Interface()
	if errInterface != nil {
		if err := errInterface.(error); err != nil {
			log.Printf("err => %+v\n", err)
		}
	}

	log.Printf("%+v %+v\n", args, reply)
}

func TestServerReturnsXMLRPCFaults(t *testing.T) {
	service := xmlrpc.NewService("Test")
	if err := service.Register("echo", func(_ *http.Request, args *EchoArgs, replies *EchoReplies) error {
		replies.Message = args.Message
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	server := xmlrpc.NewServer()
	server.Register(service)

	tests := []struct {
		name     string
		request  string
		wantCode int
	}{
		{name: "parse error", request: "<methodCall>", wantCode: -32700},
		{name: "invalid request", request: "<methodCall><methodName>.invalid</methodName><params/></methodCall>", wantCode: -32600},
		{name: "method not found", request: "<methodCall><methodName>Test.missing</methodName><params/></methodCall>", wantCode: -32601},
		{name: "invalid params", request: "<methodCall><methodName>Test.echo</methodName><params/></methodCall>", wantCode: -32602},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/RPC2", strings.NewReader(tt.request))
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, request)

			if recorder.Code != http.StatusOK {
				t.Fatalf("unexpected HTTP status: %d", recorder.Code)
			}
			_, err := xmlrpc.NewResponseReader(recorder.Body)
			if err == nil {
				t.Fatal("expected XML-RPC fault")
			}
			fault, ok := err.(xmlrpc.Error)
			if !ok {
				t.Fatalf("unexpected error type: %T", err)
			}
			if fault.Code != tt.wantCode {
				t.Fatalf("unexpected fault code: %d", fault.Code)
			}
		})
	}
}

func TestServerReturnsInternalErrorFault(t *testing.T) {
	service := xmlrpc.NewService("Test")
	if err := service.Register("fail", func(_ *http.Request, _ *EchoArgs, _ *EchoReplies) error {
		return fmt.Errorf("handler failed")
	}); err != nil {
		t.Fatal(err)
	}

	server := xmlrpc.NewServer()
	server.Register(service)

	request := httptest.NewRequest(http.MethodPost, "/RPC2", strings.NewReader(
		"<methodCall><methodName>Test.fail</methodName><params><param><value><string>hello</string></value></param></params></methodCall>",
	))
	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("unexpected HTTP status: %d", recorder.Code)
	}
	_, err := xmlrpc.NewResponseReader(recorder.Body)
	fault, ok := err.(xmlrpc.Error)
	if !ok {
		t.Fatalf("unexpected error: %v", err)
	}
	if fault.Code != -32603 {
		t.Fatalf("unexpected fault code: %d", fault.Code)
	}
}

func TestServerUsesCustomLogger(t *testing.T) {
	server := xmlrpc.NewServer()
	logger := &testLogger{}
	server.Logger = logger

	request := httptest.NewRequest(http.MethodGet, "/RPC2", nil)
	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("unexpected HTTP status: %d", recorder.Code)
	}
	if !strings.Contains(logger.buffer.String(), "405 MethodNotAllowed") {
		t.Fatalf("custom logger did not receive server log: %q", logger.buffer.String())
	}
}
