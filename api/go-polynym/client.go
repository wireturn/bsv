package polynym

import (
	"net"
	"net/http"
	"time"

	"github.com/gojektech/heimdall/v6"
	"github.com/gojektech/heimdall/v6/httpclient"
)

const (

	// version is the current version
	version = "v0.4.5"

	// defaultUserAgent is the default user agent for all requests
	defaultUserAgent string = "go-polynym: " + version

	// apiEndpoint is where we fire requests
	apiEndpoint string = "https://api.polynym.io"
)

// httpInterface is used for the http client (mocking heimdall)
type httpInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client is the parent struct that wraps the heimdall client
type Client struct {
	httpClient httpInterface // carries out the http operations (heimdall client)
	UserAgent  string        // (optional for changing user agents)
}

// Options holds all the configuration for connection, dialer and transport
type Options struct {
	BackOffExponentFactor          float64       `json:"back_off_exponent_factor"`
	BackOffInitialTimeout          time.Duration `json:"back_off_initial_timeout"`
	BackOffMaximumJitterInterval   time.Duration `json:"back_off_maximum_jitter_interval"`
	BackOffMaxTimeout              time.Duration `json:"back_off_max_timeout"`
	DialerKeepAlive                time.Duration `json:"dialer_keep_alive"`
	DialerTimeout                  time.Duration `json:"dialer_timeout"`
	RequestRetryCount              int           `json:"request_retry_count"`
	RequestTimeout                 time.Duration `json:"request_timeout"`
	TransportExpectContinueTimeout time.Duration `json:"transport_expect_continue_timeout"`
	TransportIdleTimeout           time.Duration `json:"transport_idle_timeout"`
	TransportMaxIdleConnections    int           `json:"transport_max_idle_connections"`
	TransportTLSHandshakeTimeout   time.Duration `json:"transport_tls_handshake_timeout"`
	UserAgent                      string        `json:"user_agent"`
}

// LastRequest is used to track what was submitted via the Request
type LastRequest struct {
	Method     string `json:"method"`      // Method is the HTTP method used
	StatusCode int    `json:"status_code"` // StatusCode is the last code from the request
	URL        string `json:"url"`         // URL is the url used for the request
}

// ClientDefaultOptions will return an clientOptions struct with the default settings
// Useful for starting with the base defaults and then modifying as needed
func ClientDefaultOptions() (clientOptions *Options) {
	return &Options{
		BackOffExponentFactor:          2.0,
		BackOffInitialTimeout:          2 * time.Millisecond,
		BackOffMaximumJitterInterval:   2 * time.Millisecond,
		BackOffMaxTimeout:              10 * time.Millisecond,
		DialerKeepAlive:                20 * time.Second,
		DialerTimeout:                  5 * time.Second,
		RequestRetryCount:              2,
		RequestTimeout:                 10 * time.Second,
		TransportExpectContinueTimeout: 3 * time.Second,
		TransportIdleTimeout:           20 * time.Second,
		TransportMaxIdleConnections:    10,
		TransportTLSHandshakeTimeout:   5 * time.Second,
		UserAgent:                      defaultUserAgent,
	}
}

// NewClient will make a new http client based on the options provided
func NewClient(options *Options) (c Client) {

	// Create a client
	c = Client{}

	// Set options (either default or user modified)
	if options == nil {
		options = ClientDefaultOptions()
	}

	// Set the user agent from options
	c.UserAgent = options.UserAgent

	// dial is the net dialer for clientDefaultTransport
	dial := &net.Dialer{KeepAlive: options.DialerKeepAlive, Timeout: options.DialerTimeout}

	// clientDefaultTransport is the default transport struct for the HTTP client
	clientDefaultTransport := &http.Transport{
		DialContext:           dial.DialContext,
		ExpectContinueTimeout: options.TransportExpectContinueTimeout,
		IdleConnTimeout:       options.TransportIdleTimeout,
		MaxIdleConns:          options.TransportMaxIdleConnections,
		Proxy:                 http.ProxyFromEnvironment,
		TLSHandshakeTimeout:   options.TransportTLSHandshakeTimeout,
	}

	// Determine the strategy for the http client (no retry enabled)
	if options.RequestRetryCount == 0 {
		c.httpClient = httpclient.NewClient(
			httpclient.WithHTTPTimeout(options.RequestTimeout),
			httpclient.WithHTTPClient(&http.Client{
				Transport: clientDefaultTransport,
				Timeout:   options.RequestTimeout,
			}),
		)
	} else { // Retry enabled
		// Create exponential back-off
		backOff := heimdall.NewExponentialBackoff(
			options.BackOffInitialTimeout,
			options.BackOffMaxTimeout,
			options.BackOffExponentFactor,
			options.BackOffMaximumJitterInterval,
		)

		c.httpClient = httpclient.NewClient(
			httpclient.WithHTTPTimeout(options.RequestTimeout),
			httpclient.WithRetrier(heimdall.NewRetrier(backOff)),
			httpclient.WithRetryCount(options.RequestRetryCount),
			httpclient.WithHTTPClient(&http.Client{
				Transport: clientDefaultTransport,
				Timeout:   options.RequestTimeout,
			}),
		)
	}

	return
}
