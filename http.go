package marshal

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// Client is an interface over the net/http Client methods.
type Client interface {
	Do(req *http.Request) (*http.Response, error)

	Get(url string) (resp *http.Response, err error)
	Post(url string, contentType string, body io.Reader) (resp *http.Response, err error)
}

// Post sends a POST request to the specified URL with the given body.
func Post[Body, Response any](
	client Client,
	url string,
	body Body,
) (*Response, error) {
	// Marshal the body into JSON
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	// Create a new request with the JSON body
	bodyReader := bytes.NewReader(bodyBytes)
	if client == nil {
		return nil, ErrNilHttpClient
	}
	resp, err := client.Post(url, "application/json", bodyReader)
	if err != nil {
		return nil, err
	}
	return DecodeResponse[Response](resp, nil)
}

// Get sends a GET request to the specified URL
// and decodes the JSON response into the target struct.
func Get[Response any](
	client Client,
	url string,
) (*Response, error) {
	if client == nil {
		return nil, ErrNilHttpClient
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	return DecodeResponse[Response](resp, nil)
}

// Put sends a PUT request to the specified URL
// and decodes the JSON response into the target struct.
func Put[Body, Response any](
	client Client,
	url string,
	body Body,
) (*Response, error) {
	return Request[Body, Response](client, MethodPut, url, body)
}

// Patch sends a PATCH request to the specified URL
// and decodes the JSON response into the target struct.
func Patch[Body, Response any](
	client Client,
	url string,
	body Body,
) (*Response, error) {
	return Request[Body, Response](client, MethodPut, url, body)
}

// Delete sends a DELETE request to the specified URL
// and decodes the JSON response into the target struct.
func Delete[Response any](
	client Client,
	url string,
) (*Response, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, ErrNilHttpClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return DecodeResponse[Response](resp, nil)
}

// MarshalBodyInRequest creates a new HTTP request with the specified method, URL, and body.
//
// It marshals the body into JSON and sets the appropriate headers.
// The function returns the created request and any error encountered during the process.
func MarshalBodyInRequest[Body any](
	client Client,
	method HTTPMethod,
	url string,
	body Body,
) (*http.Request, error) {
	// Marshal the body into JSON
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	// Create a new request with the JSON body
	bodyReader := bytes.NewReader(bodyBytes)
	req, err := http.NewRequest(string(method), url, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func Request[Body, Response any](
	client Client,
	method HTTPMethod,
	url string,
	body Body,
) (*Response, error) {
	req, err := MarshalBodyInRequest(client, method, url, body)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, ErrNilHttpClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return DecodeResponse[Response](resp, nil)
}

// DecodeSettings defines Settings for decoding JSON responses.
type DecodeSettings struct {
	// CheckStatusCode is the status code to check against.
	// If the response status code does not match this value, an error will be returned.
	// Default is http.StatusOK (200).
	CheckStatusCode int
}

// DecodeResponse reads and decodes a JSON response body into the target struct.
// It also closes the body automatically.
func DecodeResponse[T any](
	resp *http.Response,
	settings *DecodeSettings,
) (respBody *T, err error) {
	defer resp.Body.Close()
	if settings == nil {
		settings = &DecodeSettings{
			CheckStatusCode: http.StatusOK,
		}
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != settings.CheckStatusCode {
		return nil, &HttpError{resp.StatusCode, bodyBytes}
	}

	// Create a new instance of the target type
	var target T
	// Decode the JSON response into the target struct
	// and handle any decoding errors
	if err := json.Unmarshal(bodyBytes, &target); err != nil {
		return nil, &DecodingError{
			RawJson: bodyBytes,
			RawErr:  err,
		}
	}

	return &target, nil
}

// HTTPMethod wraps the standard net/http method constants in a string enum.
type HTTPMethod string

const (
	MethodGet     HTTPMethod = "GET"
	MethodHead    HTTPMethod = "HEAD"
	MethodPost    HTTPMethod = "POST"
	MethodPut     HTTPMethod = "PUT"
	MethodPatch   HTTPMethod = "PATCH" // RFC 5789
	MethodDelete  HTTPMethod = "DELETE"
	MethodConnect HTTPMethod = "CONNECT"
	MethodOptions HTTPMethod = "OPTIONS"
	MethodTrace   HTTPMethod = "TRACE"
)
