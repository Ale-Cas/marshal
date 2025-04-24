package marshal

import (
	"errors"
	"fmt"
)

var ErrNilHTTPClient = errors.New("nil http client error")

// HTTPError represents an error that occurred during an HTTP request.
type HTTPError struct {
	StatusCode int
	Body       []byte
}

// Error implements the error interface for HTTPError.
func (e *HTTPError) Error() string {
	return fmt.Sprintf("error in HTTP Response with StatusCode=%d and Body:\n%s", e.StatusCode, string(e.Body))
}

// DecodingError represents an error that occurred during JSON decoding.
type DecodingError struct {
	RawJson []byte
	RawErr  error
}

// Error implements the error interface for DecodingError.
func (e *DecodingError) Error() string {
	return fmt.Sprintf("error decoding %s, msg:\n%s", e.RawErr, e.RawJson)
}
