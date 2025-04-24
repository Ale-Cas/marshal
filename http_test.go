package marshal_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ale-Cas/marshal"
)

type MockClient[Response any] struct {
	StatusCode int
	Response Response
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

type Resp struct {
	ExampleInt int `json:"example_int"`
}

func TestGet(t *testing.T) {
	expectedResp := Resp{
		ExampleInt: 42,
	}
	mockClient := &MockClient[Resp]{Response: expectedResp}
	// Perform a GET request
	resp, err := marshal.Get[Resp](mockClient, httptest.DefaultRemoteAddr)
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExampleInt != expectedResp.ExampleInt {
		t.Errorf("Expected %d, got %d", expectedResp.ExampleInt, resp.ExampleInt)
	}
}

func TestPost(t *testing.T) {
	expectedResp := Resp{
		ExampleInt: 42,
	}
	mockClient := &MockClient[Resp]{Response: expectedResp}
	// Perform a POST request without a body
	resp, err := marshal.Post[any, Resp](mockClient, httptest.DefaultRemoteAddr, nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExampleInt != expectedResp.ExampleInt {
		t.Errorf("Expected %d, got %d", expectedResp.ExampleInt, resp.ExampleInt)
	}

	type Body struct {
		ExampleStr string `json:"example_str"`
	}

	// Perform a POST request with a body
	resp, err = marshal.Post[Body, Resp](mockClient, httptest.DefaultRemoteAddr, Body{ExampleStr: "test"})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExampleInt != expectedResp.ExampleInt {
		t.Errorf("Expected %d, got %d", expectedResp.ExampleInt, resp.ExampleInt)
	}
}

func TestDelete(t *testing.T) {
	expectedResp := Resp{
		ExampleInt: 42,
	}
	mockClient := &MockClient[Resp]{Response: expectedResp}
	// Perform a DELETE request
	resp, err := marshal.Delete[Resp](mockClient, httptest.DefaultRemoteAddr)
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExampleInt != expectedResp.ExampleInt {
		t.Errorf("Expected %d, got %d", expectedResp.ExampleInt, resp.ExampleInt)
	}
}