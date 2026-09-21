package xmlrpc_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	xmlrpc "github.com/mamiya312/go-xmlrpc"
)

type clientEchoArgs struct {
	Message string
}

type clientEchoReplies struct {
	Message string
}

func TestClientCall(t *testing.T) {
	service := xmlrpc.NewService("Test")
	if err := service.Register("echo", func(_ *http.Request, args *clientEchoArgs, replies *clientEchoReplies) error {
		replies.Message = strings.ToUpper(args.Message)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	server := xmlrpc.NewServer()
	server.Register(service)
	httpServer := httptest.NewServer(server)
	defer httpServer.Close()

	client := xmlrpc.NewClient(httpServer.URL)
	var reply string
	if err := client.Call(context.Background(), "Test.echo", []any{"hello"}, &reply); err != nil {
		t.Fatal(err)
	}
	if reply != "HELLO" {
		t.Fatalf("unexpected reply: %q", reply)
	}
}

func TestClientCallFault(t *testing.T) {
	httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		body, err := xmlrpc.ToResponseErrorXML(xmlrpc.Error{Code: 42, String: "failed"})
		if err != nil {
			t.Fatal(err)
		}
		w.Header().Set("Content-Type", "text/xml")
		_, _ = w.Write([]byte(body))
	}))
	defer httpServer.Close()

	var reply string
	err := xmlrpc.NewClient(httpServer.URL).Call(
		context.Background(), "Test.fail", nil, &reply,
	)
	if err == nil || !strings.Contains(err.Error(), "42: failed") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClientCallHTTPError(t *testing.T) {
	httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "server failed", http.StatusInternalServerError)
	}))
	defer httpServer.Close()

	err := xmlrpc.NewClient(httpServer.URL).Call(
		context.Background(), "Test.fail", nil,
	)
	if err == nil || !strings.Contains(err.Error(), "500") {
		t.Fatalf("unexpected error: %v", err)
	}
}
