package response

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

var (
	ErrParameterInvalid = errors.New("parameter invalid")

	ErrInvalidInput = errors.New("invalid input")
)

type CommonError struct {
	Errors []CommonErrorObject `json:"errors"`
}

type CommonErrorObject struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

// type Error struct {
// 	statusCode int
// 	Errors     []string `json:"errors"`
// }

// func NewError(err error, status int) *Error {
// 	return &Error{
// 		statusCode: status,
// 		Errors:     []string{err.Error()},
// 	}
// }

// Wrap the error info in a object
func NewError(key string, code int, err error, traceID string) CommonError {
	res := CommonError{}
	var errors []CommonErrorObject
	errorObject := CommonErrorObject{
		Code:    code,
		Message: err.Error(),
		Type:    key,
	}
	errors = append(errors, errorObject)
	res.Errors = errors
	return res
}

// func NewErrorMessage(messages []string, status int) *Error {
// 	return &Error{
// 		statusCode: status,
// 		Errors:     messages,
// 	}
// }

func (e CommonError) Send(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Errors[0].Code)
	return json.NewEncoder(w).Encode(e)
}
