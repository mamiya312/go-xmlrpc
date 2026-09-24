package xmlrpc

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// XMLRPCRequest represents an XML-RPC method call.
type XMLRPCRequest struct {
	Name   xml.Name         `xml:"methodCall"`
	Method string           `xml:"methodName"`
	Params []XMLRPCParam    `xml:"params>param"`
	Fault  XMLRPCErrorValue `xml:"fault,omitempty"`
}

// XMLRPCResponse represents an XML-RPC method response.
type XMLRPCResponse struct {
	Name   xml.Name         `xml:"methodResponse"`
	Params []XMLRPCParam    `xml:"params>param"`
	Fault  XMLRPCErrorValue `xml:"fault,omitempty"`
}

// XMLRPCParam represents one XML-RPC parameter.
type XMLRPCParam struct {
	Value string `xml:",innerxml"`
}

// XMLRPCErrorValue represents the value inside an XML-RPC fault.
type XMLRPCErrorValue struct {
	Value XMLRPCValue `xml:"value"`
}

// XMLRPCValueStructMember represents one member of an XML-RPC struct.
type XMLRPCValueStructMember struct {
	Name  string      `xml:"name"`
	Value XMLRPCValue `xml:"value"`
}

// XMLRPCValue represents a decoded XML-RPC value.
type XMLRPCValue struct {
	Array    []XMLRPCValue             `xml:"array>data>value"`
	Struct   []XMLRPCValueStructMember `xml:"struct>member"`
	String   string                    `xml:"string"`
	Int      string                    `xml:"int"`
	Int4     string                    `xml:"i4"`
	Int8     string                    `xml:"i8"`
	Double   string                    `xml:"double"`
	Boolean  string                    `xml:"boolean"`
	DateTime string                    `xml:"dateTime.iso8601"`
	Base64   string                    `xml:"base64"`
	Raw      string                    `xml:",innerxml"`
}

// IsNull reports whether the value is XML-RPC nil.
func (c *XMLRPCValue) IsNull() bool {
	return strings.TrimSpace(c.Raw) == "<nil/>"
}

func (c *XMLRPCValue) scanAsInt(field *reflect.Value) error {
	if !field.CanSet() {
		return ErrorInternalError
	}
	var val interface{}
	var err error
	switch field.Kind() {
	case reflect.Int:
		switch {
		case c.Int != "":
			val, err = strconv.Atoi(c.Int)
			if err != nil {
				return err
			}
		case c.Int4 != "":
			val, err = strconv.Atoi(c.Int4)
			if err != nil {
				return err
			}
		case c.Int8 != "":
			val, err = strconv.Atoi(c.Int8)
			if err != nil {
				return err
			}
		default:
			return ErrorInvalidParams
		}
	case reflect.Int32:
		switch {
		case c.Int != "":
			v, err := strconv.ParseInt(c.Int, 10, 32)
			if err != nil {
				return err
			}
			val = int32(v)
		case c.Int4 != "":
			v, err := strconv.ParseInt(c.Int4, 10, 32)
			if err != nil {
				return err
			}
			val = int32(v)
		case c.Int8 != "":
			v, err := strconv.ParseInt(c.Int8, 10, 32)
			if err != nil {
				return err
			}
			val = int32(v)
		default:
			return ErrorInvalidParams
		}
	case reflect.Int64:
		switch {
		case c.Int != "":
			val, err = strconv.ParseInt(c.Int, 10, 64)
			if err != nil {
				return err
			}
		case c.Int4 != "":
			val, err = strconv.ParseInt(c.Int4, 10, 64)
			if err != nil {
				return err
			}
		case c.Int8 != "":
			val, err = strconv.ParseInt(c.Int8, 10, 64)
			if err != nil {
				return ErrorInvalidParams
			}
		default:
			return ErrorInvalidParams
		}
	case reflect.Int8, reflect.Int16, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return ErrorInvalidParams
	default:
		return ErrorInvalidParams
	}

	if val != nil {
		if reflect.TypeOf(val) != reflect.TypeOf(field.Interface()) {
			return ErrorInvalidParams
		}
		field.Set(reflect.ValueOf(val))
	}

	return err
}

// Scan decodes the value into a settable reflection value.
func (c *XMLRPCValue) Scan(field *reflect.Value) error {
	if field == nil || !field.IsValid() || !field.CanSet() {
		return ErrorInternalError
	}

	if field.Kind() == reflect.Pointer {
		if c.IsNull() {
			field.SetZero()
			return nil
		}
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		return c.ScanValue(field.Elem())
	}
	return c.ScanValue(*field)
}

func (c *XMLRPCValue) ScanValue(field reflect.Value) error {
	if !field.IsValid() || !field.CanSet() {
		return ErrorInternalError
	}

	var (
		err error
		val interface{}
	)

	if field.CanAddr() && field.Addr().CanInterface() {
		if u, ok := field.Addr().Interface().(Unmarshaler); ok {
			return u.UnmarshalXMLRPC(c)
		}
	}
	if field.Type().NumMethod() > 0 && field.CanInterface() {
		if u, ok := field.Interface().(Unmarshaler); ok {
			return u.UnmarshalXMLRPC(c)
		}
	}

	switch {
	case c.Int != "":
		err = c.scanAsInt(&field)
	case c.Int4 != "":
		err = c.scanAsInt(&field)
	case c.Int8 != "":
		err = c.scanAsInt(&field)
	case c.Double != "":
		switch reflect.TypeOf(field.Interface()).Kind() {
		case reflect.Float32:
			var val64 float64
			val64, err = strconv.ParseFloat(c.Double, 32)
			val = float32(val64)
		case reflect.Float64:
			val, err = strconv.ParseFloat(c.Double, 64)
		default:
			return ErrorInvalidParams
		}
	case c.String != "":
		if field.Kind() != reflect.String {
			return ErrorInvalidParams
		}
		val = c.String
	case c.Boolean != "":
		switch c.Boolean {
		case "1", "true", "TRUE", "True":
			val = true
		case "0", "false", "FALSE", "False":
			val = false
		default:
			return ErrorInvalidParams
		}
	case c.DateTime != "":
		if field.Type() != reflect.TypeOf(time.Time{}) {
			return ErrorInvalidParams
		}
		t, err := time.Parse("20060102T15:04:05Z", c.DateTime)
		if err == nil {
			val = t.Local()
		}
	case c.Base64 != "":
		if field.Type() != reflect.TypeOf([]byte{}) {
			return ErrorInvalidParams
		}
		val, err = base64.StdEncoding.DecodeString(c.Base64)
	case c.Struct != nil:
		if field.Kind() != reflect.Struct {
			return ErrorInvalidParams
		}
		s := c.Struct
		for i := 0; i < len(s); i++ {
			index := -1
			for j := 0; j < field.NumField(); j++ {
				structField := field.Type().Field(j)
				name := structField.Name
				if tag := structField.Tag.Get("xmlrpc"); tag != "" {
					name = strings.Split(tag, ",")[0]
				}
				if name == s[i].Name {
					index = j
					break
				}
			}
			if index == -1 {
				return ErrorInvalidParams
			}
			if err = s[i].Value.ScanValue(field.Field(index)); err != nil {
				return err
			}
		}
	case c.Array != nil:
		if field.Kind() != reflect.Slice && field.Kind() != reflect.Array {
			return ErrorInvalidParams
		}
		a := c.Array
		if field.Kind() == reflect.Array && len(a) != field.Len() {
			return ErrorInvalidParams
		}
		sliceType := reflect.SliceOf(field.Type().Elem())
		slice := reflect.MakeSlice(sliceType, len(a), len(a))
		for i := 0; i < len(a); i++ {
			if err = a[i].ScanValue(slice.Index(i)); err != nil {
				return err
			}
		}
		if field.Kind() == reflect.Array {
			reflect.Copy(field, slice)
		} else {
			field.Set(slice)
		}
		return nil
	default:
		if c.IsNull() {
			return nil
		}
		return ErrorInvalidParams
	}

	if val != nil {
		if !reflect.TypeOf(val).AssignableTo(field.Type()) {
			return ErrorInvalidParams
		}
		field.Set(reflect.ValueOf(val))
	}

	return err
}

// RequestReader reads and decodes an XML-RPC method call.
type RequestReader struct {
	xmlrpcRequest XMLRPCRequest
}

// NewRequestReader parses an XML-RPC method call from r.
func NewRequestReader(r io.Reader) (*RequestReader, error) {
	var xmlrpcRequest XMLRPCRequest
	decoder := xml.NewDecoder(r)
	if err := decoder.Decode(&xmlrpcRequest); err != nil {
		return nil, err
	}
	if len(xmlrpcRequest.Fault.Value.Struct) != 0 {
		var (
			code int
			str  string
		)
		for _, field := range xmlrpcRequest.Fault.Value.Struct {
			if field.Name == "faultCode" {
				code, _ = strconv.Atoi(field.Value.Int)
			} else if field.Name == "faultString" {
				str = field.Value.String
				if str == "" {
					str = field.Value.Raw
				}
			}
		}
		return &RequestReader{xmlrpcRequest: xmlrpcRequest}, Error{Code: code, String: str}
	}
	return &RequestReader{xmlrpcRequest: xmlrpcRequest}, nil
}

// Method returns the requested XML-RPC method name.
func (c *RequestReader) Method() (string, error) {
	return c.xmlrpcRequest.Method, nil
}

// DecodeArgs decodes request parameters into individual pointers.
func (c *RequestReader) DecodeArgs(v ...interface{}) error {
	if len(v) != len(c.xmlrpcRequest.Params) {
		return ErrorInvalidParams
	}
	for i, param := range c.xmlrpcRequest.Params {
		var value XMLRPCValue
		decoder := xml.NewDecoder(bytes.NewReader([]byte(string(param.Value))))
		// TODO: Determin encoding and contert to UTF-8
		//decoder.CharsetReader = nil
		if err := decoder.Decode(&value); err != nil {
			return err
		}
		field := reflect.ValueOf(v[i])
		if !field.IsValid() || field.Kind() != reflect.Pointer || field.IsNil() {
			return ErrorInvalidParams
		}
		if err := value.ScanValue(field.Elem()); err != nil {
			return err
		}
	}
	return nil
}

// Decode decodes request parameters into the fields of a struct pointer.
func (c *RequestReader) Decode(args interface{}) error {
	argsValue := reflect.ValueOf(args)
	if !argsValue.IsValid() || argsValue.Kind() != reflect.Pointer || argsValue.IsNil() ||
		argsValue.Elem().Kind() != reflect.Struct {
		return ErrorInvalidParams
	}
	argsValue = argsValue.Elem()
	if argsValue.NumField() != len(c.xmlrpcRequest.Params) {
		return ErrorInvalidParams
	}
	for i := 0; i < argsValue.NumField(); i++ {
		var value XMLRPCValue
		decoder := xml.NewDecoder(bytes.NewReader([]byte(string(c.xmlrpcRequest.Params[i].Value))))
		// TODO: Determin encoding and contert to UTF-8
		//decoder.CharsetReader = nil
		if err := decoder.Decode(&value); err != nil {
			return err
		}
		if err := value.ScanValue(argsValue.Field(i)); err != nil {
			return err
		}
	}
	return nil
}

// ResponseReader reads and decodes an XML-RPC method response.
type ResponseReader struct {
	xmlrpcResponse XMLRPCResponse
}

// NewResponseReader parses an XML-RPC method response from r.
func NewResponseReader(r io.Reader) (*ResponseReader, error) {
	var xmlrpcResponse XMLRPCResponse
	decoder := xml.NewDecoder(r)
	if err := decoder.Decode(&xmlrpcResponse); err != nil {
		return nil, err
	}
	if len(xmlrpcResponse.Fault.Value.Struct) != 0 {
		var (
			code int
			str  string
		)
		for _, field := range xmlrpcResponse.Fault.Value.Struct {
			if field.Name == "faultCode" {
				code, _ = strconv.Atoi(field.Value.Int)
			} else if field.Name == "faultString" {
				str = field.Value.String
				if str == "" {
					str = field.Value.Raw
				}
			}
		}
		return &ResponseReader{xmlrpcResponse: xmlrpcResponse}, Error{Code: code, String: str}
	}
	return &ResponseReader{xmlrpcResponse: xmlrpcResponse}, nil
}

// Decode decodes response parameters into individual pointers.
func (c *ResponseReader) Decode(v ...any) error {
	if len(v) != len(c.xmlrpcResponse.Params) {
		return ErrorInvalidParams
	}

	for i, param := range c.xmlrpcResponse.Params {
		var value XMLRPCValue
		decoder := xml.NewDecoder(bytes.NewReader([]byte(string(param.Value))))
		// TODO: UTF-8 以外に対応する？
		//decoder.CharsetReader = nil
		if err := decoder.Decode(&value); err != nil {
			return err
		}

		field := reflect.ValueOf(v[i])
		if !field.IsValid() || field.Kind() != reflect.Pointer || field.IsNil() {
			return ErrorInvalidParams
		}
		if err := value.ScanValue(field.Elem()); err != nil {
			return err
		}
	}

	return nil
}

// ToRequestXML encodes an XML-RPC method call.
func ToRequestXML(method string, values ...interface{}) (string, error) {
	var sb strings.Builder
	sb.WriteString("<methodCall><methodName>")
	if err := xml.EscapeText(&sb, []byte(method)); err != nil {
		return sb.String(), err
	}
	sb.WriteString("</methodName><params>")
	for _, value := range values {
		s, err := ToXML(value)
		if err != nil {
			return sb.String(), err
		}
		sb.WriteString("<param>")
		sb.WriteString(s)
		sb.WriteString("</param>")
	}
	sb.WriteString("</params></methodCall>")
	return sb.String(), nil
}

// ToResponseXML encodes an XML-RPC method response.
func ToResponseXML(values ...interface{}) (string, error) {
	var sb strings.Builder
	sb.WriteString("<?xml version=\"1.0\" ?><methodResponse><params>")
	for _, value := range values {
		s, err := ToXML(value)
		if err != nil {
			return sb.String(), err
		}
		sb.WriteString("<param>")
		sb.WriteString(s)
		sb.WriteString("</param>")
	}
	sb.WriteString("</params></methodResponse>")
	return sb.String(), nil
}

// ToResponseErrorXML encodes an XML-RPC fault response.
func ToResponseErrorXML(v Error) (string, error) {
	var sb strings.Builder
	sb.WriteString("<?xml version=\"1.0\" ?><methodResponse><fault>")
	s, err := ToXML(v)
	if err != nil {
		return sb.String(), err
	}
	sb.WriteString(s)
	sb.WriteString("</fault></methodResponse>")
	return sb.String(), nil
}

// ToXML encodes a Go value as an XML-RPC value.
func ToXML(value interface{}) (string, error) {
	if value == nil {
		return "<value><nil/></value>", nil
	}
	rv, ok := value.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(value)
	}
	if !rv.IsValid() {
		return "<value><nil/></value>", nil
	}
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return "<value><nil/></value>", nil
		}
		return ToXML(rv.Elem())
	}
	rt := rv.Type()
	if value != nil {
		//rv := reflect.ValueOf(value)
		if rv.NumMethod() > 0 && rv.CanInterface() {
			if u, ok := rv.Interface().(Marshaler); ok {
				return u.MarshalXMLRPC()
			}
		}
	}
	var sb strings.Builder
	sb.WriteString("<value>")

	switch rt.Kind() {
	case reflect.Int:
		sb.WriteString(fmt.Sprintf("<int>%d</int>", rv.Interface().(int)))
	case reflect.Int32:
		sb.WriteString(fmt.Sprintf("<int>%d</int>", rv.Interface().(int32)))
	case reflect.Int64:
		sb.WriteString(fmt.Sprintf("<i8>%d</i8>", rv.Interface().(int64)))
	case reflect.Float32:
		sb.WriteString(fmt.Sprintf("<double>%g</double>", rv.Interface().(float32)))
	case reflect.Float64:
		sb.WriteString(fmt.Sprintf("<double>%g</double>", rv.Interface().(float64)))
	case reflect.String:
		s := rv.Interface().(string)
		s = strings.ReplaceAll(s, "&", "&amp;")
		s = strings.ReplaceAll(s, "\"", "&quot;")
		s = strings.ReplaceAll(s, "<", "&lt;")
		s = strings.ReplaceAll(s, ">", "&gt;")
		sb.WriteString("<string>")
		sb.WriteString(s)
		sb.WriteString("</string>")
	case reflect.Bool:
		if rv.Interface().(bool) {
			sb.WriteString("<boolean>1</boolean>")
		} else {
			sb.WriteString("<boolean>0</boolean>")
		}
	case reflect.Struct:
		switch rv.Interface().(type) {
		case time.Time:
			sb.WriteString("<dateTime.iso8601>")
			sb.WriteString(rv.Interface().(time.Time).UTC().Format("20060102T15:04:05Z"))
			sb.WriteString("</dateTime.iso8601>")
		default:
			sb.WriteString("<struct>")
			for i := 0; i < rt.NumField(); i++ {
				fieldValue := rv.Field(i)
				fieldType := rt.Field(i)
				fieldName := fieldType.Name
				if fieldType.Tag.Get("xmlrpc") != "" {
					fieldName = strings.Split(fieldType.Tag.Get("xmlrpc"), ",")[0]
				}
				if fieldName == "-" || !fieldValue.CanInterface() {
					continue
				}
				v, err := ToXML(fieldValue.Interface())
				if err != nil {
					return "", err
				}
				sb.WriteString("<member><name>")
				sb.WriteString(fieldName)
				sb.WriteString("</name>")
				sb.WriteString(v)
				sb.WriteString("</member>")
			}
			sb.WriteString("</struct>")
		}
	case reflect.Slice, reflect.Array:
		if rt.Elem().Kind() != reflect.Uint8 {
			sb.WriteString("<array><data>")
			for i := 0; i < rv.Len(); i++ {
				s, err := ToXML(rv.Index(i).Interface())
				if err != nil {
					return "", err
				}
				sb.WriteString(s)
			}
			sb.WriteString("</data></array>")
		} else if rv.Kind() == reflect.Slice {
			sb.WriteString("<base64>")
			sb.WriteString(base64.StdEncoding.EncodeToString(rv.Interface().([]byte)))
			sb.WriteString("</base64>")
		} else {
			buf := make([]byte, rv.Len())
			for i := range buf {
				buf[i] = byte(rv.Index(i).Uint())
			}
			sb.WriteString("<base64>")
			sb.WriteString(base64.StdEncoding.EncodeToString(buf))
			sb.WriteString("</base64>")
		}
	default:
		return "", fmt.Errorf("unexpected value type `%s`", rt.Kind())
	}

	sb.WriteString("</value>")
	return sb.String(), nil
}
