package tonicpow

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

// newTestClient will return a client for testing purposes
func newTestClient() (*Client, error) {
	// Create a Resty Client
	client := resty.New()

	// Get the underlying HTTP Client and set it to Mock
	httpmock.ActivateNonDefault(client.GetClient())

	// Add custom headers in request
	headers := make(map[string][]string)
	headers["custom_header_1"] = append(headers["custom_header_1"], "value_1")

	// Create a new client
	newClient, err := NewClient(
		WithRequestTracing(),
		WithAPIKey(testAPIKey),
		WithEnvironment(EnvironmentDevelopment),
		WithCustomHeaders(headers),
	)
	if err != nil {
		return nil, err
	}
	newClient.WithCustomHTTPClient(client)

	// Return the mocking client
	return newClient, nil
}

// TestNewClient will test the method NewClient()
func TestNewClient(t *testing.T) {
	t.Parallel()

	t.Run("default client", func(t *testing.T) {
		client, err := NewClient(WithAPIKey(testAPIKey))
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, defaultHTTPTimeout, client.options.httpTimeout)
		assert.Equal(t, defaultRetryCount, client.options.retryCount)
		assert.Equal(t, defaultUserAgent, client.options.userAgent)
		assert.Equal(t, false, client.options.requestTracing)
		assert.Equal(t, EnvironmentLive.apiURL, client.options.env.URL())
		assert.Equal(t, EnvironmentLive.name, client.options.env.Name())
	})

	t.Run("missing api key", func(t *testing.T) {
		client, err := NewClient(WithAPIKey(""))
		assert.Error(t, err)
		assert.Nil(t, client)
	})

	t.Run("custom http client", func(t *testing.T) {
		customHTTPClient := resty.New()
		customHTTPClient.SetTimeout(defaultHTTPTimeout)
		client, err := NewClient(WithAPIKey(testAPIKey))
		assert.NoError(t, err)
		assert.NotNil(t, client)
		client.WithCustomHTTPClient(customHTTPClient)
	})

	t.Run("custom http timeout", func(t *testing.T) {
		client, err := NewClient(WithAPIKey(testAPIKey), WithHTTPTimeout(10*time.Second))
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, 10*time.Second, client.options.httpTimeout)
	})

	t.Run("custom retry count", func(t *testing.T) {
		client, err := NewClient(WithAPIKey(testAPIKey), WithRetryCount(3))
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, 3, client.options.retryCount)
	})

	t.Run("custom headers", func(t *testing.T) {
		headers := make(map[string][]string)
		headers["custom_header_1"] = append(headers["custom_header_1"], "value_1")
		headers["custom_header_2"] = append(headers["custom_header_2"], "value_1")
		headers["custom_header_2"] = append(headers["custom_header_2"], "value_2")
		client, err := NewClient(WithAPIKey(testAPIKey), WithCustomHeaders(headers))
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, 2, len(client.options.customHeaders))
		assert.Equal(t, []string{"value_1"}, client.options.customHeaders["custom_header_1"])
	})

	t.Run("custom options", func(t *testing.T) {
		client, err := NewClient(WithAPIKey(testAPIKey), WithUserAgent("custom user agent"))
		assert.NotNil(t, client)
		assert.NoError(t, err)
		assert.Equal(t, "custom user agent", client.GetUserAgent())
	})

	t.Run("custom Environment (live)", func(t *testing.T) {
		client, err := NewClient(WithAPIKey(testAPIKey), WithEnvironment(EnvironmentLive))
		assert.NotNil(t, client)
		assert.NoError(t, err)
		assert.Equal(t, liveAPIURL, client.options.env.URL())
		assert.Equal(t, environmentLiveName, client.options.env.Name())
		assert.Equal(t, environmentLiveAlias, client.options.env.Alias())
	})

	t.Run("custom Environment (staging)", func(t *testing.T) {
		client, err := NewClient(WithAPIKey(testAPIKey), WithEnvironment(EnvironmentStaging))
		assert.NotNil(t, client)
		assert.NoError(t, err)
		assert.Equal(t, stagingAPIURL, client.options.env.URL())
		assert.Equal(t, environmentStagingName, client.options.env.Name())
		assert.Equal(t, environmentStagingAlias, client.options.env.Alias())
	})

	t.Run("custom Environment (development)", func(t *testing.T) {
		client, err := NewClient(WithAPIKey(testAPIKey), WithEnvironment(EnvironmentDevelopment))
		assert.NotNil(t, client)
		assert.NoError(t, err)
		assert.Equal(t, developmentURL, client.options.env.URL())
		assert.Equal(t, environmentDevelopmentName, client.options.env.Name())
		assert.Equal(t, environmentDevelopmentAlias, client.options.env.Alias())
	})

	t.Run("custom Environment (custom)", func(t *testing.T) {
		client, err := NewClient(
			WithAPIKey(testAPIKey),
			WithCustomEnvironment("custom", "alias", "http://localhost:5000"),
		)
		assert.NotNil(t, client)
		assert.NoError(t, err)
		assert.Equal(t, "http://localhost:5000", client.options.env.URL())
		assert.Equal(t, "custom", client.options.env.Name())
		assert.Equal(t, "alias", client.options.env.Alias())
	})

	t.Run("default no Environment", func(t *testing.T) {
		client, err := NewClient(WithAPIKey(testAPIKey), WithEnvironmentString(""))
		assert.NotNil(t, client)
		assert.NoError(t, err)
		assert.Equal(t, liveAPIURL, client.options.env.URL())
		assert.Equal(t, environmentLiveName, client.options.env.Name())
	})
}

// TestClient_GetUserAgent will test the method GetUserAgent()
func TestClient_GetUserAgent(t *testing.T) {
	t.Parallel()

	t.Run("get user agent", func(t *testing.T) {
		client, err := NewClient(WithAPIKey(testAPIKey))
		assert.NoError(t, err)
		assert.NotNil(t, client)
		userAgent := client.GetUserAgent()
		assert.Equal(t, defaultUserAgent, userAgent)
	})
}

// TestClient_GetEnvironment will test the method GetEnvironment()
func TestClient_GetEnvironment(t *testing.T) {
	t.Parallel()

	t.Run("get client Environment", func(t *testing.T) {
		client, err := NewClient(WithAPIKey(testAPIKey))
		assert.NoError(t, err)
		assert.NotNil(t, client)
		env := client.GetEnvironment()
		assert.Equal(t, environmentLiveName, env.Name())
		assert.Equal(t, environmentLiveAlias, env.Alias())
	})
}

// ExampleNewClient example using NewClient()
//
// See more examples in /examples/
func ExampleNewClient() {
	client, err := NewClient(WithAPIKey(testAPIKey))
	if err != nil {
		fmt.Printf("error loading client: %s", err.Error())
		return
	}
	fmt.Printf("loaded client: %s", client.options.userAgent)
	// Output:loaded client: go-tonicpow: v0.6.8
}

// BenchmarkNewClient benchmarks the method NewClient()
func BenchmarkNewClient(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewClient(WithAPIKey(testAPIKey))
	}
}

// TestDefaultClientOptions will test the method defaultClientOptions()
func TestDefaultClientOptions(t *testing.T) {
	t.Parallel()

	options := defaultClientOptions()
	assert.NotNil(t, options)

	assert.Equal(t, defaultHTTPTimeout, options.httpTimeout)
	assert.Equal(t, defaultRetryCount, options.retryCount)
	assert.Equal(t, defaultUserAgent, options.userAgent)
	assert.Equal(t, environmentLiveAlias, options.env.Alias())
	assert.Equal(t, environmentLiveName, options.env.Name())
	assert.Equal(t, false, options.requestTracing)
	assert.Equal(t, liveAPIURL, options.env.URL())
}

// BenchmarkDefaultClientOptions benchmarks the method defaultClientOptions()
func BenchmarkDefaultClientOptions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = defaultClientOptions()
	}
}
