package paymail

import (
	"fmt"
	"testing"
)

// TestClient_CheckDNSSEC will test the method CheckDNSSEC()
func TestClient_CheckDNSSEC(t *testing.T) {

	// t.Parallel() (turned off - race condition)

	// Integration test (requires internet connection)
	if testing.Short() {
		t.Skip("skipping integration testing in short mode")
	}

	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	var tests = []struct {
		host          string
		expectedError bool
	}{
		{"", true},
		{"---", true},
		{"---.---", true},
		{"*.---", true},
		{"moneybutton", true},
		{"asdfadfasdfasdfasdf10909.com", true},
		{"google.com", false},
		{"moneybutton.com", false},
		// {"relayx.io", false}, // Disabled for timeout issues
		{"cloudflare.com", false},
		{"mrz1836.com", false},
		{"handcash-cloud-production.herokuapp.com", true},
	}

	for _, test := range tests {
		if result := client.CheckDNSSEC(test.host); len(result.ErrorMessage) > 0 && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and error not expected but got: %s", t.Name(), test.host, result.ErrorMessage)
		} else if len(result.ErrorMessage) == 0 && test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and error was expected", t.Name(), test.host)
		}
	}

	// todo: test results, test using mock interfaces for DNS resolving
}

// ExampleClient_CheckDNSSEC example using CheckDNSSEC()
//
// See more examples in /examples/
func ExampleClient_CheckDNSSEC() {
	client, _ := NewClient()
	results := client.CheckDNSSEC("google.com")
	if len(results.ErrorMessage) == 0 {
		fmt.Printf("valid DNSSEC found for: %s", "google.com")
	} else {
		fmt.Printf("invalid DNSSEC found for: %s error: %s", "google.com", results.ErrorMessage)
	}

	// Output:valid DNSSEC found for: google.com
}

// BenchmarkClient_CheckDNSSEC benchmarks the method CheckDNSSEC()
func BenchmarkClient_CheckDNSSEC(b *testing.B) {
	client, _ := NewClient(nil, nil, nil)
	for i := 0; i < b.N; i++ {
		_ = client.CheckDNSSEC("google.com")
	}
}
