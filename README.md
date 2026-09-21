# go-xmlrpc

Go library for XML-RPC values, requests, responses, HTTP clients, and HTTP servers built on Go's standard `net/http` package.

## Client

```go
client := xmlrpc.NewClient("http://127.0.0.1:3000/RPC2")

var reply string
err := client.Call(ctx, "Test.echo", []any{"hello"}, &reply)
```

Pass one pointer for each response value after the method name and arguments. `Call` uses `context.Context`. Use `NewClientWithHTTPClient` when you need a custom timeout or transport.

## Development

```sh
go test ./...
make build
```

`make build` creates `bin/server` and `bin/client`.

The example server listens on port `3000` and accepts XML-RPC requests at `/RPC2`.

```sh
go run ./cmd/server
go run ./cmd/client
bundle exec ruby client.rb
```

`cmd/client` calls the `Test.echo` and `Test.sample` methods from `cmd/server`.

```sh
go run ./cmd/client -endpoint http://127.0.0.1:3000/RPC2
```

The Ruby client is an interoperability example. It uses the same endpoint and needs the gems in `Gemfile`.

## Server

`Server` implements `http.Handler`, so you can use it with `net/http` directly:

```go
server := xmlrpc.NewServer()
service := xmlrpc.NewService("Test")

if err := service.Register("echo", test.Echo); err != nil {
    log.Fatal(err)
}

server.Register(service)

http.Handle("/RPC2", server)
log.Fatal(http.ListenAndServe(":3000", nil))
```

You can register the server on any route and combine it with `http.ServeMux` or other HTTP middleware.

## Server methods

Register methods with this function signature:

```go
func (s *Test) Method(r *http.Request, args *Args, replies *Replies) error
```

Pass the function directly when you register it. Fields in the argument struct match request parameters in declaration order. Fields in the reply struct become response parameters.

`Server` uses the standard Go logger by default. Set a logger with a `Printf(string, ...any)` method to use custom logging.

```go
server := xmlrpc.NewServer()
server.Logger = log.New(os.Stderr, "xmlrpc: ", log.LstdFlags)
```

## Errors

Errors that can be handled as XML-RPC return an HTTP 200 response with an XML-RPC `<fault>`. Requests that cannot be handled as XML-RPC at the HTTP layer, such as non-POST requests, return an HTTP error status.
