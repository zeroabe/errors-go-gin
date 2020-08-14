package ginerrors

import (
	"errors"
	"fmt"
)

//New returns new error with passed message
func New(msg string) error {
	return errors.New(msg)
}

//Newf returns new error with message sprintf'ed by format with passed params
func Newf(format string, params ...interface{}) error {
	return fmt.Errorf(format, params...)
}

//HasErrors checks if error occurs in passed err
func HasErrors(err interface{}) bool {
	hasErrors := false
	switch e := err.(type) {
	case []error:
		hasErrors = len(e) > 0
	case map[string]error:
		hasErrors = len(e) > 0
	default:
		hasErrors = e != nil
	}

	return hasErrors
}

type FieldName string
type ValidationError string
type ErrorCode string
type ErrorType string
type DebugData interface{}

type Response struct {
	Error ErrorObject `json:"error"`
}

const (
	ErrorTypeError   ErrorType = "error"
	ErrorTypeWarning ErrorType = "warning"
	ErrorTypeInfo    ErrorType = "info"
)

type ErrorObject struct {
	Message    interface{}               `json:"message"`
	Type       *ErrorType                `json:"type,omitempty"`
	Code       *ErrorCode                `json:"code,omitempty"`
	Validation map[FieldName]interface{} `json:"validation,omitempty"`
	Debug      *DebugData                `json:"debug,omitempty"`
}
