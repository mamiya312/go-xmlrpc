# Copilot instructions

## Build and test

- Go version: `1.27` (see `go.mod`).
- Build both examples with `make build`. This writes `bin/server` and `bin/client`; the equivalent commands are `go build -o bin/server ./cmd/server` and `go build -o bin/client ./cmd/client`.
- Run all Go tests with `go test ./...`.
- Run one test or a subset with Go's `-run` filter, for example:
  `go test . -run '^TestNullableTypes$'`
- Run static checks with `go vet ./...`.
- No repository-specific lint command is configured. Format changed Go files with `gofmt` before submitting.
- The examples use the endpoint `http://127.0.0.1:3000/RPC2`. Start the server with `go run ./cmd/server`, then run the Go client with `go run ./cmd/client` or the Ruby interoperability client with `bundle exec ruby client.rb`.

## Architecture

This repository implements XML-RPC in the root package:

- `xmlrpc.go` contains the XML-RPC value model, XML encoding/decoding, request and response readers, and the `Marshaler`/`Unmarshaler` extension points.
- `types.go` defines XML-RPC faults and nullable wrapper types such as `NullString`, `NullInt64`, and `NullStruct[T]`.
- `server.go` adapts `net/http` to XML-RPC dispatch. A `Server` owns named `Service` registrations, splits method names at the final `.`, decodes request arguments, invokes the registered function through reflection, and encodes reply fields as response parameters.
- `client.go` provides the high-level HTTP client. `Client.Call` encodes a method call, sends it with a context, handles HTTP errors and XML-RPC faults, and decodes response parameters into caller-provided pointers.
- `cmd/server/main.go` is an executable example. It registers the `Test.echo` and `Test.sample` methods at `/RPC2`, listens on port `3000`, and includes request/response tracing.
- `cmd/client/main.go` is the matching Go client example. It calls the sample server's `Test.echo` and `Test.sample` methods and accepts `-endpoint` to override the default URL.
- `client_test.go` covers HTTP client/server round trips, XML-RPC faults, and HTTP errors. The other root-level `*_test.go` files cover value conversion, nullable values, request/response parsing, and reflection-based invocation; most API-level tests use the external package name `xmlrpc_test`.

## Repository-specific conventions

- XML-RPC service methods are registered directly and are expected to have the shape `func(*http.Request, *ArgsStruct, *ReplyStruct) error`.
- Service and method names are separate: register a service with `NewService("Test")`, register a method such as `"echo"`, and call it as `Test.echo`.
- `Client.Call` takes `context.Context`, the fully qualified method name, an `[]any` argument list, and one pointer per response parameter. Use `NewClientWithHTTPClient` when custom timeouts or transports are needed.
- Request and reply argument containers must be pointers to structs. Each request struct field maps to one XML-RPC parameter in declaration order; each reply struct field is emitted as one response parameter.
- Server-side XML-RPC processing errors use HTTP 200 with an XML-RPC `<fault>` response. Keep HTTP status errors for requests that cannot be handled as XML-RPC at the HTTP layer, such as non-POST requests.
- `Server.Logger` controls server logging and defaults to `log.Default()`. Custom loggers only need to implement `Printf(string, ...any)`; a nil logger disables server logging.
- Struct member names default to Go field names. Use an `xmlrpc:"name"` tag when the wire member name must differ.
- Supported wire mappings are implemented in `ToXML` and `XMLRPCValue.Scan`: signed integer forms, floating point, strings, booleans, `time.Time`, structs, arrays/slices, and `[]byte` as base64. Preserve these mappings when extending the codec.
- Character encoding support is intentionally limited to UTF-8. Do not add or plan non-UTF-8 XML character-set handling unless the project requirements explicitly change.
- Nullable types implement both `Marshaler` and `Unmarshaler`; use their `Valid` field to distinguish XML-RPC `<nil/>` from a present zero value.
- Keep protocol changes covered at the codec level and server behavior changes covered at the HTTP/dispatch level. Use the existing XML fixtures and table-driven style when adding conversion cases.
