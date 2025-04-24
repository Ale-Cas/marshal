package marshal_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ale-Cas/marshal"
)

func NewMockClient[Response any](response Response, statusCode int) *MockClient[Response] {
	return &MockClient[Response]{
		StatusCode: statusCode,
		Response:   response,
	}
}

type MockClient[Response any] struct {
	StatusCode int
	Response   Response
}

func (m *MockClient[Response]) Do(_ *http.Request) (*http.Response, error) {
	return m.mockResponse()
}

func (m *MockClient[Response]) Get(_ string) (*http.Response, error) {
	return m.mockResponse()
}

func (m *MockClient[Response]) Post(url string, _ string, _ io.Reader) (resp *http.Response, err error) {
	return m.mockResponse()
}

func (m *MockClient[Response]) mockResponse() (*http.Response, error) {
	recorder := httptest.NewRecorder()
	recorder.Code = m.StatusCode
	recorder.Header().Set("Content-Type", "application/json")
	resp, err := json.Marshal(m.Response)
	if err != nil {
		return nil, err
	}
	recorder.Body.Write(resp)
	return recorder.Result(), nil
}

type Body struct {
	ExampleStr string `json:"example_str"`
}

type Resp struct {
	ExampleInt int `json:"example_int"`
}

func TestGet(t *testing.T) {
	t.Parallel()
	expectedResp := Resp{
		ExampleInt: 42,
	}
	mockClient := NewMockClient(expectedResp, http.StatusOK)
	// Perform a GET request
	resp, err := marshal.Get[Resp](mockClient, httptest.DefaultRemoteAddr, nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExampleInt != expectedResp.ExampleInt {
		t.Errorf("Expected %d, got %d", expectedResp.ExampleInt, resp.ExampleInt)
	}
}

func TestGetDecodingError(t *testing.T) {
	t.Parallel()
	type ErrResp struct {
		ExampleInt string `json:"example_int"`
	}

	errorResp := ErrResp{
		ExampleInt: "error",
	}
	mockClient := NewMockClient(errorResp, http.StatusOK)
	// Perform a GET request
	_, err := marshal.Get[Resp](mockClient, httptest.DefaultRemoteAddr, nil)
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
	if err.Error() != "error decoding json: cannot unmarshal string into Go struct field Resp.example_int of type int, msg:\n{\"example_int\":\"error\"}" {
		t.Errorf("Expected a decoding error, got: %s", err.Error())
	}
}

func TestGetWithoutClient(t *testing.T) {
	t.Parallel()
	// Perform a GET request
	_, err := marshal.Get[Resp](nil, httptest.DefaultRemoteAddr, nil)
	if err != marshal.ErrNilHTTPClient {
		t.Fatal("Expected %w, got nil", marshal.ErrNilHTTPClient)
	}
}

func TestRequestWithoutClient(t *testing.T) {
	t.Parallel()
	tests := []struct {
		method marshal.HTTPMethod
	}{
		{method: marshal.MethodGet},
		{method: marshal.MethodPost},
		{method: marshal.MethodPut},
		{method: marshal.MethodPatch},
		{method: marshal.MethodDelete},
	}
	for _, test := range tests {
		t.Run(string(test.method), func(t *testing.T) {
			// Perform a request without a client
			body := Body{
				ExampleStr: "test",
			}
			// Perform a GET request
			_, err := marshal.Request[Body, Resp](nil, test.method, httptest.DefaultRemoteAddr, body, nil)
			if err != marshal.ErrNilHTTPClient {
				t.Fatal("Expected %w, got nil", marshal.ErrNilHTTPClient)
			}
		})
	}
}

func TestRequest(t *testing.T) {
	t.Parallel()
	tests := []struct {
		method marshal.HTTPMethod
	}{
		{method: marshal.MethodGet},
		{method: marshal.MethodPost},
		{method: marshal.MethodPut},
		{method: marshal.MethodPatch},
		{method: marshal.MethodDelete},
	}
	for _, test := range tests {
		expectedResp := Resp{
			ExampleInt: 42,
		}
		mockClient := NewMockClient(expectedResp, http.StatusOK)
		t.Run(string(test.method), func(t *testing.T) {
			// Perform a request without a client
			body := Body{
				ExampleStr: "test",
			}
			// Perform a GET request
			resp, err := marshal.Request[Body, Resp](mockClient, test.method, httptest.DefaultRemoteAddr, body, nil)
			if err != nil {
				t.Fatal(err)
			}
			if resp.ExampleInt != expectedResp.ExampleInt {
				t.Errorf("Expected %d, got %d", expectedResp.ExampleInt, resp.ExampleInt)
			}
		})
	}
}

func TestPost(t *testing.T) {
	t.Parallel()
	expectedResp := Resp{
		ExampleInt: 42,
	}
	mockClient := NewMockClient(expectedResp, http.StatusOK)
	// Perform a POST request without a body
	resp, err := marshal.Post[any, Resp](mockClient, httptest.DefaultRemoteAddr, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExampleInt != expectedResp.ExampleInt {
		t.Errorf("Expected %d, got %d", expectedResp.ExampleInt, resp.ExampleInt)
	}

	// Perform a POST request with a body
	resp, err = marshal.Post[Body, Resp](mockClient, httptest.DefaultRemoteAddr, Body{ExampleStr: "test"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExampleInt != expectedResp.ExampleInt {
		t.Errorf("Expected %d, got %d", expectedResp.ExampleInt, resp.ExampleInt)
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()
	expectedResp := Resp{
		ExampleInt: 42,
	}
	mockClient := NewMockClient(expectedResp, http.StatusOK)
	// Perform a DELETE request
	resp, err := marshal.Delete[Resp](mockClient, httptest.DefaultRemoteAddr, nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExampleInt != expectedResp.ExampleInt {
		t.Errorf("Expected %d, got %d", expectedResp.ExampleInt, resp.ExampleInt)
	}
}
