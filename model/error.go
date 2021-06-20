package model

import (
	"fmt"
	"strings"
)

// Code type errors of the API
type StatusCode string

// Error records an error, the http code and the message that caused it.
type Error struct {
	code  StatusCode
	err   error
	where string
	who   string
	// it allows to overwrite the default status
	status int
	// it is used for send additional data in the response
	data       interface{}
	apiMessage string
}

func NewError() *Error {
	return &Error{}
}

func (e *Error) Error() string {
	return e.err.Error()
}

func (e *Error) SetError(err error) { e.err = err }

func (e *Error) Code() StatusCode { return e.code }

func (e *Error) SetCode(code StatusCode) { e.code = code }

func (e *Error) Where() string { return e.where }

func (e *Error) SetWhere(where string) { e.where = where }

func (e *Error) Who() string { return e.who }

func (e *Error) SetWho(who string) { e.who = who }

func (e *Error) APIMessage() string { return e.apiMessage }

func (e *Error) SetAPIMessage(message string) {
	e.apiMessage = message
}

func (e *Error) SetErrorAsAPIMessage() {
	e.apiMessage = fmt.Sprintf("%v", e.err)
}

func (e *Error) Status() int { return e.status }

func (e *Error) SetStatus(status int) { e.status = status }

func (e *Error) Data() interface{} { return e.data }

func (e *Error) SetData(data interface{}) {
	e.data = data
}

func (e *Error) IsFailureError() bool { return strings.ToLower(string(e.code)) == "failure" }

func (e *Error) HasCode() bool { return e.code != "" }

func (e *Error) HasStatus() bool { return e.status != 0 }

func (e *Error) HasAPIMessage() bool { return e.apiMessage != "" }

func (e *Error) HasData() bool { return e.data != nil }

func (e *Error) HasWhere() bool { return e.where != "" }

func (e *Error) HasWho() bool { return e.who != "" }
