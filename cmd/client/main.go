package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"time"

	xmlrpc "github.com/mamiya312/go-xmlrpc"
)

type SampleReply struct {
	Arg1 string
	Arg2 int
	Arg3 bool
	Arg4 float64
	Arg5 []string
	Arg6 struct {
		Name string `xmlrpc:"name"`
	}
}

func main() {
	endpoint := flag.String("endpoint", "http://127.0.0.1:3000/RPC2", "XML-RPC endpoint")
	flag.Parse()

	client := xmlrpc.NewClient(*endpoint)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var echoReply string
	if err := client.Call(ctx, "Test.echo", []any{"hogehoge"}, &echoReply); err != nil {
		log.Fatal(err)
	}
	printJSON("Test.echo", echoReply)

	var sampleReply SampleReply
	if err := client.Call(ctx, "Test.sample", []any{
		"hogehoge",
		1,
		true,
		3.14,
		[]string{"hoge", "fuga"},
		struct {
			Name string `xmlrpc:"name"`
		}{Name: "hoge"},
	}, &sampleReply); err != nil {
		log.Fatal(err)
	}
	printJSON("Test.sample", sampleReply)
}

func printJSON(method string, value any) {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		log.Fatalf("%s: encode result: %s", method, err)
	}
	log.Printf("%s result:\n%s", method, data)
}
