package main

import (
	"bytes"
	"io"
	"log"
	"net/http"

	xmlrpc "github.com/mamiya312/go-xmlrpc"
)

type Test struct {
}

type EchoArgs struct {
	Message string
}

type EchoReplies struct {
	Message string
}

func (c *Test) Echo(r *http.Request, args *EchoArgs, replies *EchoReplies) error {
	replies.Message = args.Message
	return nil
}

type SampleArgs struct {
	Arg1 string
	Arg2 int
	Arg3 bool
	Arg4 float64
	Arg5 []string
	Arg6 struct {
		Name string `xmlrpc:"name"`
	}
}

type SampleReplies struct {
	Reply struct {
		Arg1 string
		Arg2 int
		Arg3 bool
		Arg4 float64
		Arg5 []string
		Arg6 struct {
			Name string `xmlrpc:"name"`
		}
	}
}

func (c *Test) Sample(r *http.Request, args *SampleArgs, replies *SampleReplies) error {
	log.Printf("Request: %v", args)
	replies.Reply.Arg1 = args.Arg1
	replies.Reply.Arg2 = args.Arg2
	replies.Reply.Arg3 = args.Arg3
	replies.Reply.Arg4 = args.Arg4
	replies.Reply.Arg5 = args.Arg5
	replies.Reply.Arg6 = struct {
		Name string `xmlrpc:"name"`
	}{
		Name: "sample",
	}
	return nil
}

type interceptWriter struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
}

func (w *interceptWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *interceptWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func trace(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bufOfRequestBody, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(bufOfRequestBody))
		iw := &interceptWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: w,
		}
		w = iw
		log.Printf("Request: %s %s %s", r.Method, r.URL.Path, string(bufOfRequestBody))
		next.ServeHTTP(w, r)
		log.Printf("Response: %d %s", iw.status, iw.body.String())
	})

}

func main() {
	test := &Test{}
	xmlrpcServer := xmlrpc.NewServer()
	testService := xmlrpc.NewService("Test")
	if err := testService.Register("echo", test.Echo); err != nil {
		log.Fatal(err)
	}
	if err := testService.Register("sample", test.Sample); err != nil {
		log.Fatal(err)
	}
	xmlrpcServer.Register(testService)
	mux := http.NewServeMux()
	mux.Handle("/RPC2", trace(xmlrpcServer))
	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatalln(err.Error())
	}
}
