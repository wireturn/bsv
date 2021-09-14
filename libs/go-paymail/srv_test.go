package paymail

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestClient_GetSRVRecord will test the method GetSRVRecord()
func TestClient_GetSRVRecord(t *testing.T) {
	// t.Parallel() (turned off - race condition)

	client, err := newTestClient()
	assert.NoError(t, err)
	assert.NotNil(t, client)

	t.Run("valid cases", func(t *testing.T) {

		var tests = []struct {
			name             string
			service          string
			protocol         string
			domainName       string
			expectedTarget   string
			expectedPort     uint16
			expectedPriority uint16
			expectedWeight   uint16
		}{
			{
				"valid - test domain",
				DefaultServiceName,
				DefaultProtocol,
				testDomain,
				"www." + testDomain,
				443,
				10,
				10,
			},
			{
				"valid - relay",
				DefaultServiceName,
				DefaultProtocol,
				"relayx.io",
				"relayx.io",
				443,
				10,
				10,
			},
		}
		var srv *net.SRV
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				srv, err = client.GetSRVRecord(test.service, test.protocol, test.domainName)
				assert.NoError(t, err)
				assert.NotNil(t, srv)
				assert.Equal(t, test.expectedPort, srv.Port)
				assert.Equal(t, test.expectedPriority, srv.Priority)
				assert.Equal(t, test.expectedWeight, srv.Weight)
				assert.Equal(t, test.expectedTarget, srv.Target)
			})
		}
	})

	t.Run("invalid cases", func(t *testing.T) {

		var tests = []struct {
			name       string
			service    string
			protocol   string
			domainName string
		}{
			{"domain not found", DefaultServiceName, DefaultProtocol, "domain.com"},
			{"missing service and protocol", "", "", "domain.com"},
			{"missing service", "", DefaultProtocol, "domain.com"},
			{"missing protocol", DefaultServiceName, "", "domain.com"},
			{"all empty", "", "", ""},
			{"missing domain", DefaultServiceName, DefaultProtocol, ""},
			{"invalid service name", "bogus", DefaultProtocol, testDomain},
			{"invalid cname", "invalid", DefaultProtocol, testDomain},
			{"no records", DefaultServiceName, DefaultProtocol, "norecords.com"},
		}
		var srv *net.SRV
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				srv, err = client.GetSRVRecord(test.service, test.protocol, test.domainName)
				assert.Error(t, err)
				assert.Nil(t, srv)
			})
		}
	})
}

// ExampleClient_GetSRVRecord example using GetSRVRecord()
//
// See more examples in /examples/
func ExampleClient_GetSRVRecord() {
	client, _ := newTestClient()
	srv, _ := client.GetSRVRecord(DefaultServiceName, DefaultProtocol, testDomain)
	fmt.Printf("port: %d priority: %d weight: %d target: %s", srv.Port, srv.Priority, srv.Weight, srv.Target)
	// Output:port: 443 priority: 10 weight: 10 target: www.test.com
}

// BenchmarkClient_GetSRVRecord benchmarks the method GetSRVRecord()
func BenchmarkClient_GetSRVRecord(b *testing.B) {
	client, _ := newTestClient()
	for i := 0; i < b.N; i++ {
		_, _ = client.GetSRVRecord(DefaultServiceName, DefaultProtocol, testDomain)
	}
}

// TestClient_ValidateSRVRecord will test the method ValidateSRVRecord()
func TestClient_ValidateSRVRecord(t *testing.T) {
	// t.Parallel() (turned off - race condition)

	client, err := newTestClient()
	assert.NoError(t, err)
	assert.NotNil(t, client)

	t.Run("valid cases", func(t *testing.T) {
		var tests = []struct {
			name     string
			srv      *net.SRV
			port     uint16
			priority uint16
			weight   uint16
		}{
			{
				"valid domain and parameters",
				&net.SRV{
					Target:   "domain.com",
					Port:     DefaultPort,
					Priority: DefaultPriority,
					Weight:   DefaultWeight,
				},
				DefaultPort,
				DefaultPriority,
				DefaultWeight,
			},
			{
				"use default ports",
				&net.SRV{
					Target:   "domain.com",
					Port:     DefaultPort,
					Priority: DefaultPriority,
					Weight:   DefaultWeight,
				},
				0,
				0,
				0,
			},
		}
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				err = client.ValidateSRVRecord(context.Background(), test.srv, test.port, test.priority, test.weight)
				assert.NoError(t, err)
			})
		}
	})

	t.Run("invalid cases", func(t *testing.T) {
		var tests = []struct {
			name     string
			srv      *net.SRV
			port     uint16
			priority uint16
			weight   uint16
		}{
			{
				"missing target",
				&net.SRV{
					Target:   "",
					Port:     DefaultPort,
					Priority: DefaultPriority,
					Weight:   DefaultWeight,
				},
				DefaultPort,
				DefaultPriority,
				DefaultWeight,
			},
			{
				"invalid domain name",
				&net.SRV{
					Target:   "domain",
					Port:     DefaultPort,
					Priority: DefaultPriority,
					Weight:   DefaultWeight,
				},
				DefaultPort,
				DefaultPriority,
				DefaultWeight,
			},
			{
				"invalid port",
				&net.SRV{
					Target:   "domain.com",
					Port:     123,
					Priority: DefaultPriority,
					Weight:   DefaultWeight,
				},
				DefaultPort,
				DefaultPriority,
				DefaultWeight,
			},
			{
				"invalid priority",
				&net.SRV{
					Target:   "domain.com",
					Port:     DefaultPort,
					Priority: 123,
					Weight:   DefaultWeight,
				},
				DefaultPort,
				DefaultPriority,
				DefaultWeight,
			},
			{
				"invalid weight",
				&net.SRV{
					Target:   "domain.com",
					Port:     DefaultPort,
					Priority: DefaultPriority,
					Weight:   123,
				},
				DefaultPort,
				DefaultPriority,
				DefaultWeight,
			},
			{
				"domain does not resolve",
				&net.SRV{
					Target:   "baddomain10901919.com",
					Port:     DefaultPort,
					Priority: DefaultPriority,
					Weight:   DefaultWeight,
				},
				DefaultPort,
				DefaultPriority,
				DefaultWeight,
			},
			{
				"no records",
				&net.SRV{
					Target:   "norecords.com",
					Port:     DefaultPort,
					Priority: DefaultPriority,
					Weight:   DefaultWeight,
				},
				DefaultPort,
				DefaultPriority,
				DefaultWeight,
			},
			{
				"srv is nil",
				nil,
				DefaultPort,
				DefaultPriority,
				DefaultWeight,
			},
		}
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				err = client.ValidateSRVRecord(context.Background(), test.srv, test.port, test.priority, test.weight)
				assert.Error(t, err)
			})
		}
	})
}

// ExampleClient_ValidateSRVRecord example using ValidateSRVRecord()
//
// See more examples in /examples/
func ExampleClient_ValidateSRVRecord() {
	client, _ := newTestClient()
	err := client.ValidateSRVRecord(
		context.Background(),
		&net.SRV{
			Target:   testDomain,
			Port:     DefaultPort,
			Priority: 1,
			Weight:   DefaultWeight,
		},
		DefaultPort,
		DefaultPriority,
		DefaultWeight,
	)
	if err != nil {
		fmt.Printf("error: %s", err.Error())
	}
	// Output:error: srv priority 1 does not match 10
}

// BenchmarkClient_ValidateSRVRecord benchmarks the method ValidateSRVRecord()
func BenchmarkClient_ValidateSRVRecord(b *testing.B) {
	client, _ := newTestClient()
	for i := 0; i < b.N; i++ {
		_ = client.ValidateSRVRecord(
			context.Background(),
			&net.SRV{
				Target:   testDomain,
				Port:     DefaultPort,
				Priority: DefaultPriority,
				Weight:   DefaultWeight,
			},
			DefaultPort,
			DefaultPriority,
			DefaultWeight,
		)
	}
}
