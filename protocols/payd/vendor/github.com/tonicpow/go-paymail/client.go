package paymail

import (
	"context"
	"net"
	"time"

	"github.com/go-resty/resty/v2"
)

// Client is the paymail client/configuration
type Client struct {
	Options  *ClientOptions    `json:"options"` // Options are all the default settings / configuration
	Resolver ResolverInterface `json:"-"`       // Resolver is used for DNS lookups
	Resty    *resty.Client     `json:"-"`       // Resty HTTP client for outgoing requests
}

// ClientOptions holds all the configuration for client requests and default resources
type ClientOptions struct {
	BRFCSpecs         []*BRFCSpec `json:"brfc_specs"`          // List of BRFC specifications
	DNSPort           string      `json:"dns_port"`            // Default DNS port for SRV checks
	DNSTimeout        int         `json:"dns_timeout"`         // Default timeout in seconds for DNS fetching
	GetTimeout        int         `json:"get_timeout"`         // Default timeout in seconds for GET requests
	NameServer        string      `json:"name_server"`         // Default name server for DNS checks
	NameServerNetwork string      `json:"name_server_network"` // Default name server network
	PostTimeout       int         `json:"post_timeout"`        // Default timeout in seconds for POST requests
	RequestTracing    bool        `json:"request_tracing"`     // If enabled, it will trace the request timing
	RetryCount        int         `json:"retry_count"`         // Default retry count for HTTP requests
	SSLDeadline       int         `json:"ssl_deadline"`        // Default timeout in seconds for SSL deadline
	SSLTimeout        int         `json:"ssl_timeout"`         // Default timeout in seconds for SSL timeout
	UserAgent         string      `json:"user_agent"`          // User agent for all outgoing requests
}

// DefaultClientOptions will return an Options struct with the default settings
//
// Useful for starting with the default and then modifying as needed
func DefaultClientOptions() (clientOptions *ClientOptions, err error) {

	// Set the default options
	clientOptions = &ClientOptions{
		DNSPort:           defaultDNSPort,
		DNSTimeout:        defaultDNSTimeout,
		GetTimeout:        defaultGetTimeout,
		NameServer:        defaultNameServer,
		NameServerNetwork: defaultNameServerNetwork,
		PostTimeout:       defaultPostTimeout,
		RequestTracing:    false,
		RetryCount:        defaultRetryCount,
		SSLDeadline:       defaultSSLDeadline,
		SSLTimeout:        defaultSSLTimeout,
		UserAgent:         defaultUserAgent,
	}

	// Load the default BRFC specs
	err = clientOptions.LoadBRFCs("")

	return
}

// ResolverInterface is a custom resolver interface for testing
type ResolverInterface interface {
	LookupHost(ctx context.Context, host string) ([]string, error)
	LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error)
	LookupSRV(ctx context.Context, service, proto, name string) (string, []*net.SRV, error)
}

// NewClient creates a new client for all outgoing paymail requests
//
// If no options are given, it will use the DefaultClientOptions()
// If no client is supplied it will use a default Resty HTTP client
func NewClient(clientOptions *ClientOptions, customClient *resty.Client,
	customResolver ResolverInterface) (client *Client, err error) {

	// Create a new client
	client = new(Client)

	// Set default options if none are provided
	if clientOptions == nil {
		clientOptions, err = DefaultClientOptions()
	} else {
		// Check for specs (if not set, use the defaults)
		if len(clientOptions.BRFCSpecs) == 0 {
			if err = clientOptions.LoadBRFCs(""); err != nil {
				// This error case should not occur since it's unmarshalling a JSON constant
				return
			}
		}
	}

	// Set the client options
	client.Options = clientOptions

	// Set the resolver
	if customResolver != nil {
		client.Resolver = customResolver
	} else {
		r := client.defaultResolver()
		client.Resolver = &r
	}

	// Set the Resty HTTP client
	if customClient != nil {
		client.Resty = customClient
	} else {
		client.Resty = resty.New()

		// Set defaults (for GET requests)
		client.Resty.SetTimeout(time.Duration(client.Options.GetTimeout) * time.Second)
		client.Resty.SetRetryCount(client.Options.RetryCount)
	}

	return
}

// getRequest is a standard GET request for all outgoing HTTP requests
func (c *Client) getRequest(requestURL string) (response StandardResponse, err error) {

	// Set the user agent
	req := c.Resty.R().SetHeader("User-Agent", c.Options.UserAgent)

	// Enable tracing
	if c.Options.RequestTracing {
		req.EnableTrace()
	}

	// Fire the request
	var resp *resty.Response
	if resp, err = req.Get(requestURL); err != nil {
		return
	}

	// Tracing enabled?
	if c.Options.RequestTracing {
		response.Tracing = resp.Request.TraceInfo()
	}

	// Set the status code
	response.StatusCode = resp.StatusCode()

	// Set the body
	response.Body = resp.Body()

	return
}

// postRequest is a standard PORT request for all outgoing HTTP requests
func (c *Client) postRequest(requestURL string, data interface{}) (response StandardResponse, err error) {

	// Set POST defaults
	c.Resty.SetTimeout(time.Duration(c.Options.PostTimeout) * time.Second)

	// Set the user agent
	req := c.Resty.R().SetBody(data).SetHeader("User-Agent", c.Options.UserAgent)

	// Enable tracing
	if c.Options.RequestTracing {
		req.EnableTrace()
	}

	// Fire the request
	var resp *resty.Response
	if resp, err = req.Post(requestURL); err != nil {
		return
	}

	// Tracing enabled?
	if c.Options.RequestTracing {
		response.Tracing = resp.Request.TraceInfo()
	}

	// Set the status code
	response.StatusCode = resp.StatusCode()

	// Set the body
	response.Body = resp.Body()

	return
}
