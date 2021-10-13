package handcash

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// mockHTTPDefaultClient for mocking requests
type mockHTTPDefaultClient struct{}

// Do is a mock http request
func (m *mockHTTPDefaultClient) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	resp.StatusCode = http.StatusOK
	resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"Message":"test"}`)))

	// Default is valid
	return resp, nil
}

// newTestClient returns a client for mocking (using a custom HTTP interface)
func newTestClient(httpClient httpInterface, environment string) *Client {
	client := NewClient(nil, nil, environment)
	client.httpClient = httpClient
	return client
}

// TestNewClient tests the method NewClient()
func TestNewClient(t *testing.T) {
	t.Parallel()

	t.Run("valid new client", func(t *testing.T) {
		client := NewClient(nil, nil, EnvironmentIAE)
		assert.NotNil(t, client)
		assert.NotNil(t, client.Options)
		assert.NotNil(t, client.httpClient)
	})

	t.Run("custom http client", func(t *testing.T) {
		client := NewClient(nil, http.DefaultClient, EnvironmentIAE)
		assert.NotNil(t, client)
		assert.NotNil(t, client.Options)
		assert.NotNil(t, client.httpClient)
	})

	t.Run("environment: iae", func(t *testing.T) {
		client := NewClient(nil, nil, EnvironmentIAE)
		assert.NotNil(t, client)
		assert.NotNil(t, client.Options)
		assert.NotNil(t, client.httpClient)
		assert.Equal(t, environments[EnvironmentIAE].APIURL, client.Environment.APIURL)
		assert.Equal(t, environments[EnvironmentIAE].ClientURL, client.Environment.ClientURL)
		assert.Equal(t, EnvironmentIAE, client.Environment.Environment)
	})

	t.Run("environment: beta", func(t *testing.T) {
		client := NewClient(nil, nil, EnvironmentBeta)
		assert.NotNil(t, client)
		assert.NotNil(t, client.Options)
		assert.NotNil(t, client.httpClient)
		assert.Equal(t, environments[EnvironmentBeta].APIURL, client.Environment.APIURL)
		assert.Equal(t, environments[EnvironmentBeta].ClientURL, client.Environment.ClientURL)
		assert.Equal(t, EnvironmentBeta, client.Environment.Environment)
	})

	t.Run("environment: production", func(t *testing.T) {
		client := NewClient(nil, nil, EnvironmentProduction)
		assert.NotNil(t, client)
		assert.NotNil(t, client.Options)
		assert.NotNil(t, client.httpClient)
		assert.Equal(t, environments[EnvironmentProduction].APIURL, client.Environment.APIURL)
		assert.Equal(t, environments[EnvironmentProduction].ClientURL, client.Environment.ClientURL)
		assert.Equal(t, EnvironmentProduction, client.Environment.Environment)
	})

	t.Run("environment: unknown", func(t *testing.T) {
		client := NewClient(nil, nil, "unknown")
		assert.NotNil(t, client)
		assert.NotNil(t, client.Options)
		assert.NotNil(t, client.httpClient)
		assert.Equal(t, environments[EnvironmentProduction].APIURL, client.Environment.APIURL)
		assert.Equal(t, environments[EnvironmentProduction].ClientURL, client.Environment.ClientURL)
		assert.Equal(t, EnvironmentProduction, client.Environment.Environment)
	})

	t.Run("environment: empty", func(t *testing.T) {
		client := NewClient(nil, nil, "")
		assert.NotNil(t, client)
		assert.NotNil(t, client.Options)
		assert.NotNil(t, client.httpClient)
		assert.Equal(t, environments[EnvironmentProduction].APIURL, client.Environment.APIURL)
		assert.Equal(t, environments[EnvironmentProduction].ClientURL, client.Environment.ClientURL)
		assert.Equal(t, EnvironmentProduction, client.Environment.Environment)
	})
}

// ExampleNewClient example using NewClient()
func ExampleNewClient() {
	client := NewClient(nil, nil, EnvironmentIAE)

	fmt.Printf("created new client: %s", client.Options.UserAgent)
	// Output:created new client: go-handcash-connect: v0.1.0
}

// BenchmarkNewClient benchmarks the method NewClient()
func BenchmarkNewClient(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewClient(nil, nil, "")
	}
}

// TestDefaultClientOptions tests setting DefaultClientOptions()
func TestDefaultClientOptions(t *testing.T) {
	t.Parallel()

	t.Run("default client options", func(t *testing.T) {
		options := DefaultClientOptions()
		assert.NotNil(t, options)
		assert.Equal(t, defaultUserAgent, options.UserAgent)
		assert.Equal(t, 2.0, options.BackOffExponentFactor)
		assert.Equal(t, 2*time.Millisecond, options.BackOffInitialTimeout)
		assert.Equal(t, 2*time.Millisecond, options.BackOffMaximumJitterInterval)
		assert.Equal(t, 10*time.Millisecond, options.BackOffMaxTimeout)
		assert.Equal(t, 20*time.Second, options.DialerKeepAlive)
		assert.Equal(t, 5*time.Second, options.DialerTimeout)
		assert.Equal(t, 2, options.RequestRetryCount)
		assert.Equal(t, 10*time.Second, options.RequestTimeout)
		assert.Equal(t, 3*time.Second, options.TransportExpectContinueTimeout)
		assert.Equal(t, 20*time.Second, options.TransportIdleTimeout)
		assert.Equal(t, 10, options.TransportMaxIdleConnections)
		assert.Equal(t, 5*time.Second, options.TransportTLSHandshakeTimeout)
	})

	t.Run("no retry", func(t *testing.T) {
		options := DefaultClientOptions()
		assert.NotNil(t, options)
		options.RequestRetryCount = 0
		client := NewClient(options, nil, "")
		assert.NotNil(t, client)
		assert.NotNil(t, client.Options)
	})
}

// ExampleDefaultClientOptions example using DefaultClientOptions()
func ExampleDefaultClientOptions() {
	options := DefaultClientOptions()
	options.UserAgent = "Custom UserAgent v1.0"
	client := NewClient(options, nil, "")

	fmt.Printf("created new client with user agent: %s", client.Options.UserAgent)
	// Output:created new client with user agent: Custom UserAgent v1.0
}

// BenchmarkDefaultClientOptions benchmarks the method DefaultClientOptions()
func BenchmarkDefaultClientOptions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = DefaultClientOptions()
	}
}
