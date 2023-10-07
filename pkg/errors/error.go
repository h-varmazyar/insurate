package errors

import (
	"fmt"
	"golang.org/x/text/language"
)

type Error struct {
	Message  string
	HttpCode int
	Code     int

	originalError error
	translates    map[language.Tag]string
	details       map[string]string
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) AddOriginalError(err error) {
	e.originalError = err
}

func (e *Error) Original() error {
	if e.originalError != nil {
		return e.originalError
	}
	return nil
}

func (e *Error) AddDetail(key string, detail interface{}) {
	if e.details == nil {
		e.details = make(map[string]string)
	}

	if v, ok := e.details[key]; !ok {
		e.details[key] = fmt.Sprintf("%v", detail)
	} else {
		e.details[key] = fmt.Sprintf("%v - %v", v, detail)
	}
}

func (e *Error) Details() map[string]string {
	if e.details == nil {
		e.details = make(map[string]string)
	}
	return e.details
}
