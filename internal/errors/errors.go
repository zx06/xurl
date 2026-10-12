package errors

import (
	"fmt"
)

const (
	CodeCfgNotFound    = "CFG_NOT_FOUND"
	CodeCfgInvalid     = "CFG_INVALID"
	CodeSecretNotFound = "SECRET_NOT_FOUND"
	CodeHTTPFile       = "HTTP_FILE"
	CodeRequest        = "REQUEST"
)

type XError struct {
	Code    string
	Message string
	Details map[string]any
	Cause   error
}

func (e *XError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func New(code, message string, details map[string]any) *XError {
	return &XError{Code: code, Message: message, Details: details}
}

func Wrap(code, message string, details map[string]any, cause error) *XError {
	return &XError{Code: code, Message: message, Details: details, Cause: cause}
}
