package errors

import "fmt"

type Error struct {
	Code        int    `json:"code"`
	Name        string `json:"error"`
	Description string `json:"-"`
}

func (e Error) Error() string {
	if e.Description == "" {
		return e.Name
	} else {
		return fmt.Sprintf("%s (%s)", e.Name, e.Description)
	}
}

func Internal() Error {
	return Error{Code: 500, Name: "internal"}
}

func Invalid(key string) Error {
	return Error{Code: 400, Name: fmt.Sprintf("invalid_%s", key)}
}

func New(text string) Error {
	return Error{Code: 500, Name: text}
}

func From(code int, name string) Error {
	return Error{Code: code, Name: name}
}

func FromError(err error) Error {
	if e, ok := err.(Error); ok {
		return e
	} else {
		e = Internal()
		e.Description = err.Error()
		return e
	}
}

func IsRetryable(err error) bool {
	return FromError(err).Code >= 500
}
