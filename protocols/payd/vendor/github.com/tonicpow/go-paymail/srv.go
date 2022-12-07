package paymail

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"
)

// defaultResolver will return a custom dns resolver
//
// This uses client options to set the network and port
func (c *Client) defaultResolver() net.Resolver {
	return net.Resolver{
		PreferGo:     true,
		StrictErrors: false,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Second * time.Duration(c.Options.DNSTimeout),
			}
			return d.DialContext(
				ctx, c.Options.NameServerNetwork, c.Options.NameServer+":"+c.Options.DNSPort,
			)
		},
	}
}

// GetSRVRecord will get the SRV record for a given domain name
//
// Specs: http://bsvalias.org/02-01-host-discovery.html
func (c *Client) GetSRVRecord(service, protocol, domainName string) (srv *net.SRV, err error) {

	// Invalid parameters?
	if len(service) == 0 { // Use the default from paymail specs
		service = DefaultServiceName
	} else if len(protocol) == 0 { // Use the default from paymail specs
		protocol = DefaultProtocol
	} else if len(domainName) == 0 || len(domainName) > 255 {
		err = fmt.Errorf("invalid parameter: domainName")
		return
	}

	// Force the case
	protocol = strings.TrimSpace(strings.ToLower(protocol))

	// The computed cname to check against
	cnameCheck := fmt.Sprintf("_%s._%s.%s.", service, protocol, domainName)

	// Lookup the SRV record
	var cname string
	var records []*net.SRV
	if cname, records, err = c.Resolver.LookupSRV(
		context.Background(), service, protocol, domainName,
	); err != nil {
		return
	} else if len(records) == 0 {
		// todo: this check might not be needed if an error is always returned
		err = fmt.Errorf("zero SRV records found using cname: %s", cnameCheck)
		return
	}

	// Basic CNAME check (sanity check!)
	if cname != cnameCheck {
		err = fmt.Errorf(
			"srv cname was invalid or not found using: %s and expected: %s",
			cnameCheck, cname,
		)
		return
	}

	// Only return the first record (in case multiple are returned)
	srv = records[0]

	// Remove any period on the end
	srv.Target = strings.TrimSuffix(srv.Target, ".")

	return
}

// ValidateSRVRecord will check for a valid SRV record for paymail following specs
//
// Specs: http://bsvalias.org/02-01-host-discovery.html
func (c *Client) ValidateSRVRecord(ctx context.Context, srv *net.SRV, port, priority, weight uint16) error {

	// Check the parameters
	if srv == nil {
		return fmt.Errorf("invalid parameter: srv is missing or nil")
	}
	if port <= 0 { // Use the default(s) from paymail specs
		port = uint16(DefaultPort)
	}
	if priority <= 0 {
		priority = uint16(DefaultPriority)
	}
	if weight <= 0 {
		weight = uint16(DefaultWeight)
	}

	// Check the basics of the SRV record
	if len(srv.Target) == 0 {
		return fmt.Errorf("srv target is invalid or empty")
	} else if srv.Port != port {
		return fmt.Errorf("srv port %d does not match %d", srv.Port, port)
	} else if srv.Priority != priority {
		return fmt.Errorf("srv priority %d does not match %d", srv.Priority, priority)
	} else if srv.Weight != weight {
		return fmt.Errorf("srv weight %d does not match %d", srv.Weight, weight)
	}

	// Test resolving the target
	if addresses, err := c.Resolver.LookupHost(ctx, srv.Target); err != nil {
		return err
	} else if len(addresses) == 0 {
		return fmt.Errorf("srv target %s could not resolve a host", srv.Target)
	}

	return nil
}
