package xmlrpc

import (
	"fmt"
	"reflect"
	"time"
)

type Error struct {
	Code   int    `xmlrpc:"faultCode"`
	String string `xmlrpc:"faultString"`
}

// Error is an XML-RPC fault.
func (f Error) Error() string {
	return fmt.Sprintf("%d: %s", f.Code, f.String)
}

var (
	// ErrorParseError reports malformed XML.
	ErrorParseError = Error{Code: -32700, String: "xml parse failed"}
	// ErrorInvalidRequest reports an invalid XML-RPC request.
	ErrorInvalidRequest = Error{Code: -32600, String: "invalid request"}
	// ErrorMethodNotFound reports an unknown XML-RPC method.
	ErrorMethodNotFound = Error{Code: -32601, String: "method not found"}
	// ErrorInvalidParams reports invalid method parameters.
	ErrorInvalidParams = Error{Code: -32602, String: "invalid params"}
	// ErrorInternalError reports an internal server error.
	ErrorInternalError = Error{Code: -32603, String: "internal error"}
)

// Unmarshaler decodes a value from an XML-RPC value.
type Unmarshaler interface {
	UnmarshalXMLRPC(*XMLRPCValue) error
}

// Marshaler encodes a value as an XML-RPC value.
type Marshaler interface {
	MarshalXMLRPC() (string, error)
}

// NullInt represents an optional int value.
type NullInt struct {
	Value int
	Valid bool // Valid is true if value it not Null
}

func (c *NullInt) UnmarshalXMLRPC(v *XMLRPCValue) error {
	if v.IsNull() {
		c.Valid = false
		return nil
	}
	c.Valid = true
	field := reflect.ValueOf(&c.Value).Elem()
	return v.Scan(&field)
}

func (c *NullInt) MarshalXMLRPC() (string, error) {
	if !c.Valid {
		return ToXML(nil)
	}
	return ToXML(c.Value)
}

// NullInt32 represents an optional int32 value.
type NullInt32 struct {
	Value int32
	Valid bool // Valid is true if value it not Null
}

func (c *NullInt32) UnmarshalXMLRPC(v *XMLRPCValue) error {
	if v.IsNull() {
		c.Valid = false
		return nil
	}
	c.Valid = true
	field := reflect.ValueOf(&c.Value).Elem()
	return v.Scan(&field)
}

func (c *NullInt32) MarshalXMLRPC() (string, error) {
	if !c.Valid {
		return ToXML(nil)
	}
	return ToXML(c.Value)
}

// NullInt64 represents an optional int64 value.
type NullInt64 struct {
	Value int64
	Valid bool // Valid is true if value it not Null
}

func (c *NullInt64) UnmarshalXMLRPC(v *XMLRPCValue) error {
	if v.IsNull() {
		c.Valid = false
		return nil
	}
	c.Valid = true
	field := reflect.ValueOf(&c.Value).Elem()
	return v.Scan(&field)
}

func (c *NullInt64) MarshalXMLRPC() (string, error) {
	if !c.Valid {
		return ToXML(nil)
	}
	return ToXML(c.Value)
}

// NullFloat32 represents an optional float32 value.
type NullFloat32 struct {
	Value float32
	Valid bool // Valid is true if value it not Null
}

func (c *NullFloat32) UnmarshalXMLRPC(v *XMLRPCValue) error {
	if v.IsNull() {
		c.Valid = false
		return nil
	}
	c.Valid = true
	field := reflect.ValueOf(&c.Value).Elem()
	return v.Scan(&field)
}

func (c *NullFloat32) MarshalXMLRPC() (string, error) {
	if !c.Valid {
		return ToXML(nil)
	}
	return ToXML(c.Value)
}

// NullFloat64 represents an optional float64 value.
type NullFloat64 struct {
	Value float64
	Valid bool // Valid is true if value it not Null
}

func (c *NullFloat64) UnmarshalXMLRPC(v *XMLRPCValue) error {
	if v.IsNull() {
		c.Valid = false
		return nil
	}
	c.Valid = true
	field := reflect.ValueOf(&c.Value).Elem()
	return v.Scan(&field)
}

func (c *NullFloat64) MarshalXMLRPC() (string, error) {
	if !c.Valid {
		return ToXML(nil)
	}
	return ToXML(c.Value)
}

// NullString represents an optional string value.
type NullString struct {
	Value string
	Valid bool // Valid is true if value it not Null
}

func (c *NullString) UnmarshalXMLRPC(v *XMLRPCValue) error {
	if v.IsNull() {
		c.Valid = false
		return nil
	}
	c.Valid = true
	field := reflect.ValueOf(&c.Value).Elem()
	return v.Scan(&field)
}

func (c *NullString) MarshalXMLRPC() (string, error) {
	if !c.Valid {
		return ToXML(nil)
	}
	return ToXML(c.Value)
}

// NullBool represents an optional bool value.
type NullBool struct {
	Value bool
	Valid bool // Valid is true if value it not Null
}

func (c *NullBool) UnmarshalXMLRPC(v *XMLRPCValue) error {
	if v.IsNull() {
		c.Valid = false
		return nil
	}
	c.Valid = true
	field := reflect.ValueOf(&c.Value).Elem()
	return v.Scan(&field)
}

func (c *NullBool) MarshalXMLRPC() (string, error) {
	if !c.Valid {
		return ToXML(nil)
	}
	return ToXML(c.Value)
}

// NullTime represents an optional time.Time value.
type NullTime struct {
	Value time.Time
	Valid bool // Valid is true if value it not Null
}

func (c *NullTime) UnmarshalXMLRPC(v *XMLRPCValue) error {
	if v.IsNull() {
		c.Valid = false
		return nil
	}
	c.Valid = true
	field := reflect.ValueOf(&c.Value).Elem()
	return v.Scan(&field)
}

func (c *NullTime) MarshalXMLRPC() (string, error) {
	if !c.Valid {
		return ToXML(nil)
	}
	return ToXML(c.Value)
}

// NullBytes represents an optional []byte value.
type NullBytes struct {
	Value []byte
	Valid bool // Valid is true if value it not Null
}

func (c *NullBytes) UnmarshalXMLRPC(v *XMLRPCValue) error {
	if v.IsNull() {
		c.Valid = false
		return nil
	}
	c.Valid = true
	field := reflect.ValueOf(&c.Value).Elem()
	return v.Scan(&field)
}

func (c *NullBytes) MarshalXMLRPC() (string, error) {
	if !c.Valid {
		return ToXML(nil)
	}
	return ToXML(c.Value)
}

// NullStruct represents an optional struct value.
type NullStruct[T interface{}] struct {
	Value T
	Valid bool // Valid is true if value it not Null
}

func (c *NullStruct[T]) UnmarshalXMLRPC(v *XMLRPCValue) error {
	if v.IsNull() {
		c.Valid = false
		return nil
	}
	c.Valid = true
	field := reflect.ValueOf(&c.Value).Elem()
	return v.Scan(&field)
}

func (c *NullStruct[T]) MarshalXMLRPC() (string, error) {
	if !c.Valid {
		return ToXML(nil)
	}
	return ToXML(c.Value)
}
