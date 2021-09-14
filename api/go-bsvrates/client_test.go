package bsvrates

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestNewClient test new client
func TestNewClient(t *testing.T) {
	t.Parallel()

	t.Run("valid new client", func(t *testing.T) {
		client := NewClient(nil, nil)
		assert.NotNil(t, client)

		// Test default providers
		assert.Equal(t, ProviderCoinPaprika, client.Providers[0])
		assert.Equal(t, ProviderWhatsOnChain, client.Providers[1])
		assert.Equal(t, ProviderPreev, client.Providers[2])
	})

	t.Run("custom http client", func(t *testing.T) {
		client := NewClient(nil, http.DefaultClient, ProviderPreev)
		assert.NotNil(t, client)

		// Test custom providers
		assert.Equal(t, ProviderPreev, client.Providers[0])
	})

}

// BenchmarkNewClient benchmarks the NewClient method
func BenchmarkNewClient(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewClient(nil, nil)
	}
}

// TestDefaultClientOptions tests setting DefaultClientOptions()
func TestDefaultClientOptions(t *testing.T) {
	t.Parallel()

	t.Run("default client options", func(t *testing.T) {
		options := DefaultClientOptions()

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
		options.RequestRetryCount = 0
		client := NewClient(options, nil)
		assert.NotNil(t, client)
	})
}

// TestClientOptions_ToPreevOptions tests setting ToPreevOptions()
func TestClientOptions_ToPreevOptions(t *testing.T) {
	t.Parallel()

	t.Run("convert default options to Preev options", func(t *testing.T) {
		options := DefaultClientOptions()
		preevOptions := options.ToPreevOptions()

		assert.Contains(t, preevOptions.UserAgent, defaultUserAgent)
		assert.Equal(t, 2.0, preevOptions.BackOffExponentFactor)
		assert.Equal(t, 2*time.Millisecond, preevOptions.BackOffInitialTimeout)
		assert.Equal(t, 2*time.Millisecond, preevOptions.BackOffMaximumJitterInterval)
		assert.Equal(t, 10*time.Millisecond, preevOptions.BackOffMaxTimeout)
		assert.Equal(t, 20*time.Second, preevOptions.DialerKeepAlive)
		assert.Equal(t, 5*time.Second, preevOptions.DialerTimeout)
		assert.Equal(t, 2, preevOptions.RequestRetryCount)
		assert.Equal(t, 10*time.Second, preevOptions.RequestTimeout)
		assert.Equal(t, 3*time.Second, preevOptions.TransportExpectContinueTimeout)
		assert.Equal(t, 20*time.Second, preevOptions.TransportIdleTimeout)
		assert.Equal(t, 10, preevOptions.TransportMaxIdleConnections)
		assert.Equal(t, 5*time.Second, preevOptions.TransportTLSHandshakeTimeout)
	})
}

// TestClientOptions_ToWhatsOnChainOptions tests setting ToWhatsOnChainOptions()
func TestClientOptions_ToWhatsOnChainOptions(t *testing.T) {
	t.Parallel()

	t.Run("convert default options to WhatsOnChain options", func(t *testing.T) {
		options := DefaultClientOptions()
		whatsOnChainOptions := options.ToWhatsOnChainOptions()

		assert.Contains(t, whatsOnChainOptions.UserAgent, defaultUserAgent)
		assert.Equal(t, 2.0, whatsOnChainOptions.BackOffExponentFactor)
		assert.Equal(t, 2*time.Millisecond, whatsOnChainOptions.BackOffInitialTimeout)
		assert.Equal(t, 2*time.Millisecond, whatsOnChainOptions.BackOffMaximumJitterInterval)
		assert.Equal(t, 10*time.Millisecond, whatsOnChainOptions.BackOffMaxTimeout)
		assert.Equal(t, 20*time.Second, whatsOnChainOptions.DialerKeepAlive)
		assert.Equal(t, 5*time.Second, whatsOnChainOptions.DialerTimeout)
		assert.Equal(t, 2, whatsOnChainOptions.RequestRetryCount)
		assert.Equal(t, 10*time.Second, whatsOnChainOptions.RequestTimeout)
		assert.Equal(t, 3*time.Second, whatsOnChainOptions.TransportExpectContinueTimeout)
		assert.Equal(t, 20*time.Second, whatsOnChainOptions.TransportIdleTimeout)
		assert.Equal(t, 10, whatsOnChainOptions.TransportMaxIdleConnections)
		assert.Equal(t, 5*time.Second, whatsOnChainOptions.TransportTLSHandshakeTimeout)
	})
}
