package xmlrpc

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
)

type method struct {
	Method    reflect.Value
	ArgsType  reflect.Type
	ReplyType reflect.Type
}

type Service struct {
	name    string
	methods map[string]*method
}

// NewService creates a named XML-RPC service.
func NewService(name string) *Service {
	return &Service{
		name:    name,
		methods: map[string]*method{},
	}
}

// Register adds a method to the service.
//
// The function must have the signature
// func(*http.Request, *Args, *Replies) error.
func (s *Service) Register(name string, function any) error {
	fvalue := reflect.ValueOf(function)
	if !fvalue.IsValid() {
		return fmt.Errorf("function is invalid")
	}
	ftype := fvalue.Type()

	if ftype.Kind() != reflect.Func {
		return fmt.Errorf("function is not a function")
	}
	if ftype.NumIn() != 3 {
		return fmt.Errorf("function must have 3 arguments")
	}
	requestType := reflect.TypeOf((*http.Request)(nil))
	if ftype.In(0) != requestType {
		return fmt.Errorf("first argument must have type *http.Request")
	}
	if ftype.In(1).Kind() != reflect.Ptr || ftype.In(1).Elem().Kind() != reflect.Struct {
		return fmt.Errorf("second argument must be a pointer to a struct")
	}
	if ftype.In(2).Kind() != reflect.Ptr || ftype.In(2).Elem().Kind() != reflect.Struct {
		return fmt.Errorf("third argument must be a pointer to a struct")
	}
	if ftype.NumOut() != 1 {
		return fmt.Errorf("function must have one return value")
	}
	if !ftype.Out(0).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return fmt.Errorf("function must return error")
	}

	s.methods[name] = &method{
		Method:    fvalue,
		ArgsType:  ftype.In(1),
		ReplyType: ftype.In(2),
	}
	return nil
}

type Server struct {
	services map[string]*Service
	Logger   Logger
}

// Logger is the logging interface used by Server.
type Logger interface {
	Printf(format string, v ...any)
}

// NewServer creates an XML-RPC HTTP server.
func NewServer() *Server {
	return &Server{
		services: map[string]*Service{},
		Logger:   log.Default(),
	}
}

// Register adds a service to the server.
func (c *Server) Register(service *Service) {
	c.services[service.name] = service
}

func (s *Server) logf(format string, v ...any) {
	if s.Logger == nil {
		return
	}
	s.Logger.Printf(format, v...)
}

func (s *Server) writeFault(w http.ResponseWriter, fault Error) {
	body, err := ToResponseErrorXML(fault)
	if err != nil {
		s.logf("%d failed to encode XML-RPC fault: %s", http.StatusInternalServerError, err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	if _, err := w.Write([]byte(body)); err != nil {
		s.logf("fault response write failed: %s", err.Error())
	}
}

// ServeHTTP handles XML-RPC requests as an http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.logf("%d %s", http.StatusMethodNotAllowed, "MethodNotAllowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	reader, err := NewRequestReader(r.Body)
	if err != nil {
		s.logf("%d %s", http.StatusOK, err.Error())
		s.writeFault(w, ErrorParseError)
		return
	}

	serviceName := ""
	methodName, err := reader.Method()
	if err != nil {
		s.logf("%d %s", http.StatusOK, err.Error())
		s.writeFault(w, ErrorInvalidRequest)
		return
	}
	if strings.HasPrefix(methodName, ".") {
		s.logf("%d %s", http.StatusOK, "method name starts with `.`")
		s.writeFault(w, ErrorInvalidRequest)
		return
	}
	if idx := strings.LastIndex(methodName, "."); idx != -1 {
		// Hoge.fuga
		// => serviceName: Hoge
		//    methodName: fuga
		serviceName = methodName[0:idx]
		methodName = methodName[idx+1:]
	}

	service, ok := s.services[serviceName]
	if !ok {
		s.logf("%d %s", http.StatusOK, "service not found")
		s.writeFault(w, ErrorMethodNotFound)
		return
	}
	method, ok := service.methods[methodName]
	if !ok {
		s.logf("%d %s", http.StatusOK, "method not found")
		s.writeFault(w, ErrorMethodNotFound)
		return
	}
	args := reflect.New(method.ArgsType.Elem())
	if err := reader.Decode(args.Interface()); err != nil {
		s.logf("%d request arguments decode fail: %s", http.StatusOK, err.Error())
		s.writeFault(w, ErrorInvalidParams)
		return
	}

	replies := reflect.New(method.ReplyType.Elem())
	rv := method.Method.Call([]reflect.Value{
		reflect.ValueOf(r),
		args,
		replies,
	})
	if len(rv) != 1 {
		s.logf("%d unexpected handler result count", http.StatusOK)
		s.writeFault(w, ErrorInternalError)
		return
	}
	if ei := rv[0].Interface(); ei != nil {
		if err := ei.(error); err != nil {
			s.logf("%d %s", http.StatusOK, err.Error())
			s.writeFault(w, ErrorInternalError)
			return
		}
	}

	repliesValues := replies.Elem()
	if repliesValues.Kind() != reflect.Struct {
		s.logf("%d %s", http.StatusOK, "unexpected replies type")
		s.writeFault(w, ErrorInternalError)
		return
	}
	repliesSlice := []interface{}{}
	for i := 0; i < repliesValues.NumField(); i++ {
		repliesSlice = append(repliesSlice, repliesValues.Field(i))
	}

	res, err := ToResponseXML(repliesSlice...)
	if err != nil {
		s.logf("%d %s", http.StatusOK, err.Error())
		s.writeFault(w, ErrorInternalError)
		return
	}
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	if _, err := w.Write([]byte(res)); err != nil {
		s.logf("response write failed: %s", err.Error())
	}
}
