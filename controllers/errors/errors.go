// Package errors handles all structures related with errors.
package errors

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	// ErrorNotFound shows that requested object is not found
	ErrorNotFound = errors.New("not found")

	// ErrorUnsupportedMethod shows that requested method is not allowed
	ErrorUnsupportedMethod = errors.New("this method is not allowed in current url")

	// ErrorInvalidBody shows request body is invalid
	ErrorInvalidBody = errors.New("request body is invalid")

	// ErrorContentType shows only application/json allowed for content-type
	ErrorContentType = errors.New("unsupported content-type, only application/json is allowed")

	// ErrorMongoAggregate shows there is an error in aggregation
	ErrorMongoAggregate = errors.New("error in mongo aggregate")
)

// ErrorResponse represents error response body.
type ErrorResponse struct {
	Err        string `json:"error"`
	StatusCode int    `json:"StatusCode"`
}

// WriteError writes error to http response and gives case.
func WriteError(w http.ResponseWriter, err error) {
	var statusCode int
	switch err {
	case ErrorNotFound:
		statusCode = http.StatusNotFound
	case ErrorUnsupportedMethod:
		statusCode = http.StatusMethodNotAllowed
	case ErrorInvalidBody:
		statusCode = http.StatusBadRequest
	case ErrorMongoAggregate:
		statusCode = http.StatusBadRequest
	case ErrorContentType:
		statusCode = http.StatusBadRequest
	default:
		statusCode = http.StatusInternalServerError
	}
	errorResponse := ErrorResponse{Err: err.Error(), StatusCode: statusCode}
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(&errorResponse)
}
