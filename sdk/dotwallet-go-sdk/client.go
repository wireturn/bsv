package dotwallet

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

// Client is the TonicPow client/configuration
type Client struct {
	options *clientOptions // Options are all the default settings / configuration
}

// clientOptions holds all the configuration for client requests and default resources
type clientOptions struct {
	clientID       string              // Your client ID
	clientSecret   string              // Your client secret
	customHeaders  map[string][]string // Custom headers on outgoing requests
	getTokenOnLoad bool                // If enabled, it will automatically fetch the application token
	host           string              // API host
	httpClient     *resty.Client       // Resty client for all HTTP requests
	httpTimeout    time.Duration       // Default timeout in seconds for GET requests
	redirectURI    string              // The redirect URI for oauth requests
	requestTracing bool                // If enabled, it will trace the request timing
	retryCount     int                 // Default retry count for HTTP requests
	token          *DotAccessToken     // The application authentication token (object)
	userAgent      string              // User agent for all outgoing requests
}

// ClientOps allow functional options to be supplied
// that overwrite default client options.
type ClientOps func(c *clientOptions)

// WithHTTPTimeout can be supplied to adjust the default http client timeouts.
// The http client is used when creating requests
// Default timeout is 10 seconds.
func WithHTTPTimeout(timeout time.Duration) ClientOps {
	return func(c *clientOptions) {
		c.httpTimeout = timeout
	}
}

// WithRequestTracing will enable tracing.
// Tracing is disabled by default.
func WithRequestTracing() ClientOps {
	return func(c *clientOptions) {
		c.requestTracing = true
	}
}

// WithAutoLoadToken will load the application token when spawning a new client
// By default, this is disabled, so you can make a new client without making an HTTP request
func WithAutoLoadToken() ClientOps {
	return func(c *clientOptions) {
		c.getTokenOnLoad = true
	}
}

// WithRetryCount will overwrite the default retry count for http requests.
// Default retries is 2.
func WithRetryCount(retries int) ClientOps {
	return func(c *clientOptions) {
		c.retryCount = retries
	}
}

// WithUserAgent will overwrite the default useragent.
// Default is package name + version.
func WithUserAgent(userAgent string) ClientOps {
	return func(c *clientOptions) {
		if len(userAgent) > 0 {
			c.userAgent = userAgent
		}
	}
}

// WithCredentials will set the ClientID and ClientSecret
// There is no default
func WithCredentials(clientID, clientSecret string) ClientOps {
	return func(c *clientOptions) {
		c.clientID = strings.TrimSpace(clientID)
		c.clientSecret = strings.TrimSpace(clientSecret)
	}
}

// WithHost will overwrite the default host.
// Default is located in definitions
func WithHost(host string) ClientOps {
	return func(c *clientOptions) {
		host = strings.TrimSpace(host)
		if len(host) > 0 {
			c.host = host
		}
	}
}

// WithRedirectURI will set the redirect URI for user authentication callbacks
// There is no default
func WithRedirectURI(redirectURI string) ClientOps {
	return func(c *clientOptions) {
		c.redirectURI = strings.TrimSpace(redirectURI)
	}
}

// WithCustomHeaders will add custom headers to outgoing requests
// Custom headers is empty by default
func WithCustomHeaders(headers map[string][]string) ClientOps {
	return func(c *clientOptions) {
		c.customHeaders = headers
	}
}

// WithCustomHTTPClient will overwrite the default client with a custom client.
func WithCustomHTTPClient(client *resty.Client) ClientOps {
	return func(c *clientOptions) {
		c.httpClient = client
	}
}

// WithApplicationAccessToken will override the token on the client
func WithApplicationAccessToken(token DotAccessToken) ClientOps {
	return func(c *clientOptions) {
		c.token = &token
	}
}

// NewClient creates a new client for all TonicPow requests
//
// If no options are given, it will use the DefaultClientOptions()
// If there is no client is supplied, it will use a default Resty HTTP client.
func NewClient(opts ...ClientOps) (*Client, error) {
	defaults := defaultClientOptions()

	// Create a new client
	client := &Client{
		options: defaults,
	}

	// overwrite defaults with any set by user
	for _, opt := range opts {
		opt(client.options)
	}

	// Check for Client ID
	if client.options.clientID == "" {
		return nil, errors.New("missing an client_id")
	}

	// Check for Client Secret
	if client.options.clientSecret == "" {
		return nil, errors.New("missing an client_secret")
	}

	// Set the Resty HTTP client
	if client.options.httpClient == nil {
		client.options.httpClient = resty.New()
		// Set defaults (for GET requests)
		client.options.httpClient.SetTimeout(client.options.httpTimeout)
		client.options.httpClient.SetRetryCount(client.options.retryCount)
	}

	// Autoload the token?
	if client.options.getTokenOnLoad {
		if err := client.UpdateApplicationAccessToken(); err != nil {
			return nil, err
		}
	}

	return client, nil
}

// defaultClientOptions will return a clientOptions struct with the default settings
//
// Useful for starting with the default and then modifying as needed
func defaultClientOptions() (opts *clientOptions) {
	// Set the default options
	opts = &clientOptions{
		host:           defaultHost,
		httpTimeout:    defaultHTTPTimeout,
		requestTracing: false,
		retryCount:     defaultRetryCount,
		userAgent:      defaultUserAgent,
	}
	return
}

// GetUserAgent will return the user agent string of the client
func (c *Client) GetUserAgent() string {
	return c.options.userAgent
}

// Token will return the token for the application
func (c *Client) Token() *DotAccessToken {
	// Skip panics
	if c == nil || c.options == nil {
		return nil
	}
	return c.options.token
}

// NewState returns a new state key for CSRF protection
// It is not required to use this method, you can use your own UUID() for state validation
//
// For more information: https://developers.dotwallet.com/documents/en/#user-authorization
func (c *Client) NewState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// Request is a standard GET / POST / PUT / DELETE request for all outgoing HTTP requests
// Omit the data attribute if using a GET request
func (c *Client) Request(httpMethod string, requestEndpoint string,
	data interface{}, expectedCode int, authorization *DotAccessToken) (response *StandardResponse, err error) {

	// Set the user agent
	req := c.options.httpClient.R().SetHeader("User-Agent", c.options.userAgent)

	// Set the body if (PUT || POST)
	if httpMethod != http.MethodGet && httpMethod != http.MethodDelete {
		var j []byte
		if j, err = json.Marshal(data); err != nil {
			return
		}
		req = req.SetBody(string(j))
		req.Header.Add("Content-Length", strconv.Itoa(len(j)))
		req.Header.Set("Content-Type", "application/json")
	}

	// Enable tracing
	if c.options.requestTracing {
		req.EnableTrace()
	}

	// Set the authorization and content type (assuming it's application)
	if authorization != nil && len(authorization.RefreshToken) == 0 {

		// If the token is nil or expired, get a new token
		if c.IsTokenExpired(authorization) {

			// Update the access token
			if err = c.UpdateApplicationAccessToken(); err != nil {
				return
			}
		}

		// Set the header authorization for the application
		req.Header.Set(headerAuthorization, authorization.TokenType+" "+authorization.AccessToken)
	} else if authorization != nil && len(authorization.RefreshToken) > 0 {

		// If the token is expired, get a new token if we had a previous token
		if c.IsTokenExpired(authorization) {

			// Try to update the access token first
			if authorization, err = c.RefreshUserToken(authorization); err != nil {
				return
			}
		}

		// Set the header authorization for the application
		req.Header.Set(headerAuthorization, authorization.TokenType+" "+authorization.AccessToken)
	}

	// Custom headers?
	for key, headers := range c.options.customHeaders {
		for _, value := range headers {
			req.Header.Set(key, value)
		}
	}

	// Fire the request
	var resp *resty.Response
	switch httpMethod {
	case http.MethodPost:
		resp, err = req.Post(c.options.host + requestEndpoint)
	// case http.MethodPut:
	// 	resp, err = req.Put(c.options.host + requestEndpoint)
	// case http.MethodDelete:
	// 	resp, err = req.Delete(c.options.host + requestEndpoint)
	case http.MethodGet:
		resp, err = req.Get(c.options.host + requestEndpoint)
	}
	if err != nil {
		return
	}

	// Start the response
	response = new(StandardResponse)

	// Tracing enabled?
	if c.options.requestTracing {
		response.Tracing = resp.Request.TraceInfo()
	}

	// Set the status code & body
	response.StatusCode = resp.StatusCode()
	response.Body = resp.Body()

	// Check expected code if set
	if expectedCode > 0 && response.StatusCode != expectedCode {
		response.Error = new(Error)
		if err = json.Unmarshal(response.Body, &response.Error); err != nil {
			return
		}
		err = fmt.Errorf("%s", response.Error.Message)
		response.Error.Method = httpMethod
		response.Error.URL = requestEndpoint
	}

	return
}
