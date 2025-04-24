package marshal

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

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
	if client == nil {
		client = http.DefaultClient
	}
	// Marshal the body into JSON
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	// Create a new request with the JSON body
	bodyReader := bytes.NewReader(bodyBytes)
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
		client = http.DefaultClient
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	return DecodeResponse[Response](resp, nil)
}

// Delete sends a DELETE request to the specified URL
// and decodes the JSON response into the target struct.
func Delete[Response any](
	client Client,
	url string,
) (*Response, error) {
	if client == nil {
		client = http.DefaultClient
	}
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return DecodeResponse[Response](resp, nil)
}

// DecodeSettings defines Settings for decoding JSON responses.
type DecodeSettings struct {
	// DisallowUnknownFields indicates whether
	// to disallow unknown fields in the JSON response.
	DisallowUnknownFields bool

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
