package paymail

import (
	"context"
	"net"
	"time"

	"github.com/go-resty/resty/v2"
)

// Client is the paymail client/configuration
type Client struct {
	options    *clientOptions // Options are all the default settings / configuration
	resolver   DNSResolver
	httpClient *resty.Client
}

// ClientOptions holds all the configuration for client requests and default resources
type clientOptions struct {
	brfcSpecs         []*BRFCSpec   // List of BRFC specifications
	dnsPort           string        // Default DNS port for SRV checks
	dnsTimeout        time.Duration // Default timeout in seconds for DNS fetching
	httpTimeout       time.Duration // Default timeout in seconds for GET requests
	nameServer        string        // Default name server for DNS checks
	nameServerNetwork string        // Default name server network
	requestTracing    bool          // If enabled, it will trace the request timing
	retryCount        int           // Default retry count for HTTP requests
	sslDeadline       time.Duration // Default timeout in seconds for SSL deadline
	sslTimeout        time.Duration // Default timeout in seconds for SSL timeout
	userAgent         string        // User agent for all outgoing requests

}

// ClientOps allow functional options to be supplied
// that overwrite default go-paymail client options.
type ClientOps func(c *clientOptions)

// WithDNSPort can be supplied with a custom dns port to perform SRV checks on.
// Default is 53.
func WithDNSPort(port string) ClientOps {
	return func(c *clientOptions) {
		c.dnsPort = port
	}
}

// WithDNSTimeout can be supplied to overwrite the default dns srv check timeout.
// The default is 5 seconds.
func WithDNSTimeout(timeout time.Duration) ClientOps {
	return func(c *clientOptions) {
		c.dnsTimeout = timeout
	}
}

// WithBRFCSpecs allows custom specs to be supplied to extend or replace the defaults.
func WithBRFCSpecs(specs []*BRFCSpec) ClientOps {
	return func(c *clientOptions) {
		c.brfcSpecs = specs
	}
}

// WithHTTPTimeout can be supplied to adjust the default http client timeouts.
// The http client is used when querying paymail services for capabilities
// Default timeout is 20 seconds.
func WithHTTPTimeout(timeout time.Duration) ClientOps {
	return func(c *clientOptions) {
		c.httpTimeout = timeout
	}
}

// WithNameServer can be supplied to overwrite the default name server used to resolve srv requests.
// default is 8.8.8.8.
func WithNameServer(ip string) ClientOps {
	return func(c *clientOptions) {
		c.nameServer = ip
	}
}

// WithNameServerNetwork can overwrite the default network protocol to use.
// The default is udp.
func WithNameServerNetwork(network string) ClientOps {
	return func(c *clientOptions) {
		c.nameServerNetwork = network
	}
}

// WithRequestTracing will enable tracing.
// Tracing is disabled by default.
func WithRequestTracing() ClientOps {
	return func(c *clientOptions) {
		c.requestTracing = true
	}
}

// WithRetryCount will overwrite the default retry count for http requests.
// Default retries is 2.
func WithRetryCount(retries int) ClientOps {
	return func(c *clientOptions) {
		c.retryCount = retries
	}
}

// WithSSLTimeout will overwrite the default ssl timeout.
// Default timeout is 10 seconds.
func WithSSLTimeout(timeout time.Duration) ClientOps {
	return func(c *clientOptions) {
		c.sslTimeout = timeout
	}
}

// WithSSLDeadline will overwrite the default ssl deadline.
// Default is 10 seconds.
func WithSSLDeadline(timeout time.Duration) ClientOps {
	return func(c *clientOptions) {
		c.sslDeadline = timeout
	}
}

// WithUserAgent will overwrite the default useragent.
// Default is go-paymail + version.
func WithUserAgent(userAgent string) ClientOps {
	return func(c *clientOptions) {
		c.userAgent = userAgent
	}
}

// WithCustomResolver will allow you to supply a custom  dns resolver,
// useful for testing etc.
func (c *Client) WithCustomResolver(resolver DNSResolver) *Client {
	c.resolver = resolver
	return c
}

// WithCustomHTTPClient will overwrite the default client with a custom client.
func (c *Client) WithCustomHTTPClient(client *resty.Client) *Client {
	c.httpClient = client
	return c
}

// GetBRFCs will return the list of specs
func (c *Client) GetBRFCs() []*BRFCSpec {
	return c.options.brfcSpecs
}

// GetUserAgent will return the user agent string of the client
func (c *Client) GetUserAgent() string {
	return c.options.userAgent
}

// defaultClientOptions will return an clientOptions struct with the default settings
//
// Useful for starting with the default and then modifying as needed
func defaultClientOptions() (opts *clientOptions, err error) {
	// Set the default options
	opts = &clientOptions{
		dnsPort:           defaultDNSPort,
		dnsTimeout:        defaultDNSTimeout,
		httpTimeout:       defaultHTTPTimeout,
		nameServer:        defaultNameServer,
		nameServerNetwork: defaultNameServerNetwork,
		requestTracing:    false,
		retryCount:        defaultRetryCount,
		sslDeadline:       defaultSSLDeadline,
		sslTimeout:        defaultSSLTimeout,
		userAgent:         defaultUserAgent,
	}
	// Load the default BRFC specs
	err = opts.LoadBRFCs("")
	return
}

// DNSResolver is a custom resolver interface for testing
type DNSResolver interface {
	LookupHost(ctx context.Context, host string) ([]string, error)
	LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error)
	LookupSRV(ctx context.Context, service, proto, name string) (string, []*net.SRV, error)
}

// NewClient creates a new client for all paymail requests
//
// If no options are given, it will use the defaultClientOptions()
// If no client is supplied it will use a default Resty HTTP client
func NewClient(opts ...ClientOps) (*Client, error) {

	// Start with the defaults
	defaults, err := defaultClientOptions()
	if err != nil {
		return nil, err
	}

	// Create a new client
	client := &Client{
		options: defaults,
	}

	// Overwrite defaults with any set by user
	for _, opt := range opts {
		opt(client.options)
	}

	// Check for specs (if not set, use the defaults)
	if len(client.options.brfcSpecs) == 0 {
		if err = client.options.LoadBRFCs(""); err != nil {
			return nil, err
		}
	}

	// Set the resolver
	if client.resolver == nil {
		r := client.defaultResolver()
		client.resolver = &r
	}

	// Set the Resty HTTP client
	if client.httpClient == nil {
		client.httpClient = resty.New()
		// Set defaults (for GET requests)
		client.httpClient.SetTimeout(client.options.httpTimeout)
		client.httpClient.SetRetryCount(client.options.retryCount)
	}
	return client, nil
}

// getRequest is a standard GET request for all outgoing HTTP requests
func (c *Client) getRequest(requestURL string) (response StandardResponse, err error) {
	// Set the user agent
	req := c.httpClient.R().SetHeader("User-Agent", c.options.userAgent)

	// Enable tracing
	if c.options.requestTracing {
		req.EnableTrace()
	}

	// Fire the request
	var resp *resty.Response
	if resp, err = req.Get(requestURL); err != nil {
		return
	}

	// Tracing enabled?
	if c.options.requestTracing {
		response.Tracing = resp.Request.TraceInfo()
	}

	// Set the status code
	response.StatusCode = resp.StatusCode()

	// Set the body
	response.Body = resp.Body()
	return
}

// postRequest is a standard POST request for all outgoing HTTP requests
func (c *Client) postRequest(requestURL string, data interface{}) (response StandardResponse, err error) {
	// Set the user agent
	req := c.httpClient.R().SetBody(data).SetHeader("User-Agent", c.options.userAgent)

	// Enable tracing
	if c.options.requestTracing {
		req.EnableTrace()
	}

	// Fire the request
	var resp *resty.Response
	if resp, err = req.Post(requestURL); err != nil {
		return
	}

	// Tracing enabled?
	if c.options.requestTracing {
		response.Tracing = resp.Request.TraceInfo()
	}

	// Set the status code
	response.StatusCode = resp.StatusCode()

	// Set the body
	response.Body = resp.Body()
	return
}
