package logger

import (
	"encoding/json"
)

// Error is a standardized error struct
type Error struct {
	K string      `json:"key"`
	V interface{} `json:"value"`
}

// Key implements the Logger KeyValue interface
func (e *Error) Key() string {
	return e.K
}

// Value implements the Logger KeyValue interface
func (e *Error) Value() interface{} {
	return e.V
}

// Error json encodes the error, implements the std error interface
func (e *Error) Error() string {
	data, _ := json.Marshal(e) // disregard error
	return string(data)
}

// MakeError creates a new Error (key/value)
func MakeError(key string, value interface{}) *Error {
	return &Error{
		K: key,
		V: value,
	}
}
