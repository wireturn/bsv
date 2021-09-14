package spvchannels

import "net/http"

// HTTPClient is the interface for http client.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// MockClient mocks the http client.
type MockClient struct {
	MockDo func(req *http.Request) (*http.Response, error)
}

// Do implement the mock of Do method for http client interface
func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.MockDo(req)
}
