package paymail

import (
	"time"

	"github.com/go-resty/resty/v2"
)

// Defaults for paymail functions
const (
	defaultDNSPort           = "53"                     // Default port for DNS / NameServer checks
	defaultDNSTimeout        = 5 * time.Second          // In seconds
	defaultHTTPTimeout       = 20 * time.Second         // Default timeout for all GET requests in seconds
	defaultNameServer        = "8.8.8.8"                // Default DNS NameServer
	defaultNameServerNetwork = "udp"                    // Default for NS dialer
	defaultRetryCount        = 2                        // Default retry count for HTTP requests
	defaultSSLDeadline       = 10 * time.Second         // Default deadline in seconds
	defaultSSLTimeout        = 10 * time.Second         // Default timeout in seconds
	defaultUserAgent         = "go-paymail: " + version // Default user agent
	version                  = "v0.2.11"                // Go-Paymail version
)

// Public defaults for paymail specs
/*
	http://bsvalias.org/02-01-host-discovery.html

	Service	  bsvalias
	Proto	  tcp
	Name	  <domain>.<tld>.
	TTL	      3600 (see notes)
	Class	  IN
	Priority  10
	Weight	  10
	Port	  443
	Target	  <endpoint-discovery-host>

	Max SRV Records:  1
*/
const (
	DefaultBsvAliasVersion = "1.0"      // Default version number for bsvalias
	DefaultPort            = 443        // Default port (from specs)
	DefaultPriority        = 10         // Default priority (from specs)
	DefaultProtocol        = "tcp"      // Default protocol (from specs)
	DefaultServiceName     = "bsvalias" // Default service name (from specs)
	DefaultWeight          = 10         // Default weight (from specs)
	PubKeyLength           = 66         // Required length for a valid PubKey (pki)
)

// StandardResponse is the standard fields returned on all responses
type StandardResponse struct {
	Body       []byte          `json:"-"` // Body of the response request
	StatusCode int             `json:"-"` // Status code returned on the request
	Tracing    resty.TraceInfo `json:"-"` // Trace information if enabled on the request
}

/*
Example error response
{
    "code": "not-found",
    "message": "Paymail not found: mrz@mneybutton.com"
}
*/

// ServerError is the standard error response from a paymail server
type ServerError struct {
	Code    string `json:"code"`    // Shows the corresponding code
	Message string `json:"message"` // Shows the error message returned by the server
}
