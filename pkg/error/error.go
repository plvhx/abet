package error

import (
    "fmt"
)

type FieldError struct {
    Field  string `json:"field"`
    Reason string `json:"reason"`
}

func (fe FieldError) Error() string {
    return fe.Reason
}

type CustomError struct {
    StatusCode int          `json:"statusCode"`
    Message    string       `json:"message"`
    Errors     []FieldError `json:"errors,omitempty"`
}

func (ce CustomError) Error() string {
    return ce.Message
}

func Errorf(httpCode int, format string, args ...any) CustomError {
    return CustomError{
        StatusCode: httpCode,
        Message: fmt.Sprintf(format, args...),
    }
}

func Error(httpCode int, buf string) CustomError {
    return CustomError{
        StatusCode: httpCode,
        Message: buf,
    }
}
