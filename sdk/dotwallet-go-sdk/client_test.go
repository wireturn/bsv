package dotwallet

import (
	"fmt"
	"net/http"
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
		WithCredentials(testClientID, testClientSecret),
		WithCustomHeaders(headers),
		WithCustomHTTPClient(client),
		WithHost(testHost),
		WithRedirectURI(testRedirectURI),
		WithRequestTracing(),
	)
	if err != nil {
		return nil, err
	}

	// Return the mocking client
	return newClient, nil
}

// TestNewClient will test the method NewClient()
func TestNewClient(t *testing.T) {

	t.Run("missing client_id", func(t *testing.T) {
		customHTTPClient := resty.New()
		customHTTPClient.SetTimeout(defaultHTTPTimeout)
		httpmock.ActivateNonDefault(customHTTPClient.GetClient())

		newClient, err := NewClient(
			WithCredentials("", testClientSecret),
			WithCustomHTTPClient(customHTTPClient),
		)
		assert.Error(t, err)
		assert.Nil(t, newClient)
	})

	t.Run("missing client_secret", func(t *testing.T) {
		customHTTPClient := resty.New()
		customHTTPClient.SetTimeout(defaultHTTPTimeout)
		httpmock.ActivateNonDefault(customHTTPClient.GetClient())

		newClient, err := NewClient(
			WithCredentials(testClientID, ""),
			WithCustomHTTPClient(customHTTPClient),
		)
		assert.Error(t, err)
		assert.Nil(t, newClient)
	})

	t.Run("test client, get access token, defaults", func(t *testing.T) {
		c, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, c)

		mockUpdateApplicationAccessToken()
		err = c.UpdateApplicationAccessToken()
		assert.NoError(t, err)
		assert.NotNil(t, c.options.token)

		assert.Equal(t, defaultHTTPTimeout, c.options.httpTimeout)
		assert.Equal(t, defaultRetryCount, c.options.retryCount)
		assert.Equal(t, defaultUserAgent, c.options.userAgent)

		assert.Equal(t, testTokenType, c.options.token.TokenType)
		assert.Equal(t, testExpiresIn, c.options.token.ExpiresIn)
		assert.Equal(t, time.Now().UTC().Unix()+c.options.token.ExpiresIn, c.options.token.ExpiresAt)
		assert.Equal(t, testAccessToken, c.options.token.AccessToken)
		assert.Equal(t, testClientID, c.options.clientID)
		assert.Equal(t, testClientSecret, c.options.clientSecret)
		assert.Equal(t, testRedirectURI, c.options.redirectURI)
		assert.Equal(t, testHost, c.options.host)
	})

	t.Run("custom HTTP client", func(t *testing.T) {
		customHTTPClient := resty.New()
		customHTTPClient.SetTimeout(defaultHTTPTimeout)
		httpmock.ActivateNonDefault(customHTTPClient.GetClient())

		newClient, err := NewClient(
			WithCredentials(testClientID, testClientSecret),
			WithCustomHTTPClient(customHTTPClient),
		)
		assert.NoError(t, err)
		assert.NotNil(t, newClient)
	})

	t.Run("custom HTTP timeout", func(t *testing.T) {
		newClient, err := NewClient(
			WithCredentials(testClientID, testClientSecret),
			WithHTTPTimeout(15*time.Second),
		)
		assert.NoError(t, err)
		assert.NotNil(t, newClient)
		assert.Equal(t, 15*time.Second, newClient.options.httpTimeout)
	})

	t.Run("custom host", func(t *testing.T) {
		newClient, err := NewClient(
			WithCredentials(testClientID, testClientSecret),
			WithHost(testHost),
		)
		assert.NoError(t, err)
		assert.NotNil(t, newClient)
		assert.Equal(t, testHost, newClient.options.host)
	})

	t.Run("custom retry count", func(t *testing.T) {
		newClient, err := NewClient(
			WithCredentials(testClientID, testClientSecret),
			WithRetryCount(4),
		)
		assert.NoError(t, err)
		assert.NotNil(t, newClient)
		assert.Equal(t, 4, newClient.options.retryCount)
	})

	t.Run("custom user agent", func(t *testing.T) {
		newClient, err := NewClient(
			WithCredentials(testClientID, testClientSecret),
			WithUserAgent("custom-user-agent"),
		)
		assert.NoError(t, err)
		assert.NotNil(t, newClient)
		assert.Equal(t, "custom-user-agent", newClient.options.userAgent)
	})

	t.Run("custom headers", func(t *testing.T) {
		customHTTPClient := resty.New()
		httpmock.ActivateNonDefault(customHTTPClient.GetClient())

		headers := make(map[string][]string)
		headers["custom_header_1"] = append(headers["custom_header_1"], "value_1")
		headers["custom_header_2"] = append(headers["custom_header_2"], "value_1")
		headers["custom_header_2"] = append(headers["custom_header_2"], "value_2")
		client, err := NewClient(
			WithCredentials(testClientID, testClientSecret),
			WithCustomHTTPClient(customHTTPClient),
			WithCustomHeaders(headers),
		)
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, 2, len(client.options.customHeaders))
		assert.Equal(t, []string{"value_1"}, client.options.customHeaders["custom_header_1"])
	})

	t.Run("auto-load the application token", func(t *testing.T) {
		customHTTPClient := resty.New()
		httpmock.ActivateNonDefault(customHTTPClient.GetClient())
		mockUpdateApplicationAccessToken()

		client, err := NewClient(
			WithHost(testHost),
			WithCredentials(testClientID, testClientSecret),
			WithCustomHTTPClient(customHTTPClient),
			WithAutoLoadToken(),
		)
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.NotNil(t, client.Token())
	})

	t.Run("auto-load fails", func(t *testing.T) {
		customHTTPClient := resty.New()
		httpmock.ActivateNonDefault(customHTTPClient.GetClient())
		mockAccessTokenFailed(http.StatusBadRequest)

		client, err := NewClient(
			WithHost(testHost),
			WithCredentials(testClientID, testClientSecret),
			WithCustomHTTPClient(customHTTPClient),
			WithAutoLoadToken(),
		)
		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Nil(t, client.Token())
	})

	t.Run("override the application access token", func(t *testing.T) {
		newClient, err := NewClient(
			WithCredentials(testClientID, testClientSecret),
			WithHTTPTimeout(15*time.Second),
			WithApplicationAccessToken(DotAccessToken{
				AccessToken: testAccessToken,
				ExpiresAt:   time.Now().UTC().Unix() + testExpiresIn,
				ExpiresIn:   testExpiresIn,
				TokenType:   testTokenType,
			}),
		)
		assert.NoError(t, err)
		assert.NotNil(t, newClient)

		token := newClient.Token()
		assert.Equal(t, testAccessToken, token.AccessToken)
	})
}

// TestClient_GetUserAgent will test the method GetUserAgent()
func TestClient_GetUserAgent(t *testing.T) {
	t.Run("get user agent", func(t *testing.T) {
		c, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, c)
		userAgent := c.GetUserAgent()
		assert.Equal(t, defaultUserAgent, userAgent)
	})
}

// TestClient_NewState will test the method NewState()
func TestClient_NewState(t *testing.T) {
	t.Run("generate a new state", func(t *testing.T) {
		c, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, c)

		var state string
		state, err = c.NewState()
		assert.NoError(t, err)
		assert.Equal(t, 64, len(state))
	})
}

// TestClient_Token will test the method Token()
func TestClient_Token(t *testing.T) {
	t.Run("token is present", func(t *testing.T) {
		c, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, c)

		mockUpdateApplicationAccessToken()
		err = c.UpdateApplicationAccessToken()
		assert.NoError(t, err)

		token := c.Token()
		assert.NotNil(t, token)
	})

	t.Run("token is not present", func(t *testing.T) {
		c, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, c)

		mockAccessTokenFailed(http.StatusBadRequest)
		err = c.UpdateApplicationAccessToken()
		assert.Error(t, err)

		token := c.Token()
		assert.Nil(t, token)
	})
}

// ExampleNewClient example using NewClient()
//
// See more examples in /examples/
func ExampleNewClient() {
	client, err := NewClient(
		WithCredentials("your-client-id", "your-secret-id"),
		WithRedirectURI("http://localhost:3000/v1/auth/dotwallet/callback"),
	)
	if err != nil {
		fmt.Printf("error loading client: %s", err.Error())
		return
	}
	fmt.Printf("loaded client: %s", client.options.userAgent)
	// Output:loaded client: dotwallet-go-sdk: v0.0.3
}

// BenchmarkNewClient benchmarks the method NewClient()
func BenchmarkNewClient(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewClient(
			WithCredentials(testClientID, testClientSecret),
		)
	}
}

// TestDefaultClientOptions will test the method defaultClientOptions()
func TestDefaultClientOptions(t *testing.T) {
	options := defaultClientOptions()
	assert.NotNil(t, options)

	assert.Equal(t, defaultHost, options.host)
	assert.Equal(t, defaultHTTPTimeout, options.httpTimeout)
	assert.Equal(t, defaultRetryCount, options.retryCount)
	assert.Equal(t, defaultUserAgent, options.userAgent)
	assert.Equal(t, false, options.requestTracing)
}

// BenchmarkDefaultClientOptions benchmarks the method defaultClientOptions()
func BenchmarkDefaultClientOptions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = defaultClientOptions()
	}
}
