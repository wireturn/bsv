package paymail

import (
	"context"
	"fmt"
	"net"
)

// resolver for mocking requests
type resolver struct {
	hosts        map[string][]string
	ipAddresses  map[string][]net.IPAddr
	liveResolver DNSResolver
	srvRecords   map[string][]*net.SRV
}

// newCustomResolver will return a custom resolver with specific records hard coded ,
func newCustomResolver(liveResolver DNSResolver, hosts map[string][]string,
	srvRecords map[string][]*net.SRV, ipAddresses map[string][]net.IPAddr) DNSResolver {
	return &resolver{
		hosts:        hosts,
		ipAddresses:  ipAddresses,
		liveResolver: liveResolver,
		srvRecords:   srvRecords,
	}
}

// LookupHost will lookup a host
func (r *resolver) LookupHost(ctx context.Context, host string) ([]string, error) {
	records, ok := r.hosts[host]
	if ok {
		return records, nil
	}
	return r.liveResolver.LookupHost(ctx, host)
}

// LookupIPAddr will lookup an ip address
func (r *resolver) LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error) {
	records, ok := r.ipAddresses[host]
	if ok {
		return records, nil
	}
	return r.liveResolver.LookupIPAddr(ctx, host)
}

// LookupSRV will lookup an SRV record
func (r *resolver) LookupSRV(ctx context.Context, service, proto, name string) (string, []*net.SRV, error) {
	records, ok := r.srvRecords[service+proto+name]
	if ok {
		if service == "invalid" { // Returns an invalid cname
			return fmt.Sprintf("_%s._%s", service, proto), records, nil
		}
		return fmt.Sprintf("_%s._%s.%s.", service, proto, name), records, nil
	}
	return r.liveResolver.LookupSRV(ctx, service, proto, name)
}
