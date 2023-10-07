package errors

import (
	"net/http"
)

var (
	DefaultHttpCode = http.StatusInternalServerError
	DefaultCode     = 0
)

func NewWithHttp(msg string, code, httpCode int) *Error {
	err := &Error{
		Message:    msg,
		HttpCode:   httpCode,
		Code:       code,
		translates: translate(msg),
		details:    make(map[string]string),
	}
	return err
}

func NewWithCode(msg string, code int) *Error {
	return NewWithHttp(msg, code, DefaultHttpCode)
}

func New(msg string) *Error {
	return NewWithCode(msg, DefaultCode)
}
