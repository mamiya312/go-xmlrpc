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

func (f Error) Error() string {
	return fmt.Sprintf("%d: %s", f.Code, f.String)
}

var (
	ErrorParseError     = Error{Code: -32700, String: "xml parse failed"}
	ErrorInvalidRequest = Error{Code: -32600, String: "invalid request"}
	ErrorMethodNotFound = Error{Code: -32601, String: "method not found"}
	ErrorInvalidParams  = Error{Code: -32602, String: "invalid params"}
	ErrorInternalError  = Error{Code: -32603, String: "internal error"}
)

type Unmarshaler interface {
	UnmarshalXMLRPC(*XMLRPCValue) error
}

type Marshaler interface {
	MarshalXMLRPC() (string, error)
}

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
