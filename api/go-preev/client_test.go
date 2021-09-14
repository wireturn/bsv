package preev

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestNewClient test new client
func TestNewClient(t *testing.T) {
	t.Parallel()

	client := NewClient(nil, nil)

	assert.NotEqual(t, "", client.UserAgent)
	assert.NotNil(t, client)
}

// TestNewClient_CustomHttpClient test new client with custom HTTP client
func TestNewClient_CustomHttpClient(t *testing.T) {
	t.Parallel()

	client := NewClient(nil, http.DefaultClient)

	assert.Equal(t, "", client.UserAgent)
	assert.NotNil(t, client)
}

// ExampleNewClient example using NewClient()
func ExampleNewClient() {
	client := NewClient(nil, nil)
	fmt.Println(client.UserAgent)
	// Output:go-preev: v0.2.4
}

// BenchmarkNewClient benchmarks the NewClient method
func BenchmarkNewClient(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewClient(nil, nil)
	}
}

// TestClientDefaultOptions tests setting ClientDefaultOptions()
func TestClientDefaultOptions(t *testing.T) {
	t.Parallel()

	options := ClientDefaultOptions()

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
}

// TestClientDefaultOptions_NoRetry will set 0 retry counts
func TestClientDefaultOptions_NoRetry(t *testing.T) {
	options := ClientDefaultOptions()
	options.RequestRetryCount = 0
	client := NewClient(options, nil)

	assert.NotNil(t, client)
	assert.NotNil(t, options)
	assert.Equal(t, defaultUserAgent, options.UserAgent)
}
