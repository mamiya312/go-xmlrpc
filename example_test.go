package xmlrpc_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	xmlrpc "github.com/mamiya312/go-xmlrpc"
)

func ExampleClient_Call() {
	type EchoArgs struct {
		Message string
	}
	type EchoReplies struct {
		Message string
	}

	service := xmlrpc.NewService("Example")
	if err := service.Register("echo", func(_ *http.Request, args *EchoArgs, replies *EchoReplies) error {
		replies.Message = args.Message
		return nil
	}); err != nil {
		panic(err)
	}

	server := xmlrpc.NewServer()
	server.Register(service)
	httpServer := httptest.NewServer(server)
	defer httpServer.Close()

	client := xmlrpc.NewClient(httpServer.URL)
	var reply string
	if err := client.Call(context.Background(), "Example.echo", []any{"hello"}, &reply); err != nil {
		panic(err)
	}
	fmt.Println(reply)

	// Output:
	// hello
}

func ExampleServer() {
	type EchoArgs struct {
		Message string
	}
	type EchoReplies struct {
		Message string
	}

	service := xmlrpc.NewService("Example")
	if err := service.Register("echo", func(_ *http.Request, args *EchoArgs, replies *EchoReplies) error {
		replies.Message = args.Message
		return nil
	}); err != nil {
		panic(err)
	}

	server := xmlrpc.NewServer()
	server.Register(service)

	mux := http.NewServeMux()
	mux.Handle("/RPC2", server)
	_ = mux
}
