package polynym

import (
	"fmt"
	"net/http"
	"testing"
)

// newMockClient will create a new mock client for testing
func newMockClient(userAgent string) Client {
	return Client{
		httpClient: &mockHTTP{},
		UserAgent:  userAgent,
	}
}

// TestGetAddress tests the GetAddress()
func TestGetAddress(t *testing.T) {
	t.Parallel()

	// Create a mock client
	client := newMockClient(defaultUserAgent)

	// Create the list of tests
	var tests = []struct {
		input         string
		expected      string
		expectedError bool
		statusCode    int
	}{
		{"", "", true, http.StatusBadRequest},
		{"error", "", true, http.StatusBadRequest},
		{"bad-poly-response", "", true, http.StatusBadRequest},
		{"bad-poly-status", "", true, http.StatusBadRequest},
		{"doesnotexist@handcash.io", "", true, http.StatusBadRequest},
		{"$mr-z", "124dwBFyFtkcNXGfVWQroGcT9ybnpQ3G3Z", false, http.StatusOK},
		{"19gKzz8XmFDyrpk4qFobG7qKoqybe78v9h", "19gKzz8XmFDyrpk4qFobG7qKoqybe78v9h", false, http.StatusOK},
		{"1doesnotexisthandle", "", true, http.StatusBadRequest},
		{"1mrz", "1Lti3s6AQNKTSgxnTyBREMa6XdHLBnPSKa", false, http.StatusOK},
		{"bad@paymailaddress.com", "", true, http.StatusBadRequest},
		{"c6ZqP5Tb22KJuvSAbjNkoi", "", true, http.StatusBadRequest},
		{"mrz@handcash.io", "19gKzz8XmFDyrpk4qFobG7qKoqybe78v9h", false, http.StatusOK},
		{"@833", "19ksW6ueSw9nEj88X3QNJ9VkKPGf1zuKbQ", false, http.StatusOK},
	}

	// Test all
	for _, test := range tests {
		if output, err := GetAddress(client, test.input); err == nil && test.expectedError {
			t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted and [%s] expected", t.Name(), test.input, test.expected)
		} else if err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%v] error [%s]", t.Name(), test.input, test.expected, output, err.Error())
		} else if output != nil && output.Address != test.expected && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%s]", t.Name(), test.input, test.expected, output.Address)
		} else if output != nil && output.LastRequest.Method != http.MethodGet {
			t.Errorf("%s Expected method to be %s, got %s, [%s] inputted", t.Name(), http.MethodGet, output.LastRequest.Method, test.input)
		} else if output != nil && output.LastRequest.StatusCode != test.statusCode {
			t.Errorf("%s Expected status code to be %d, got %d, [%s] inputted", t.Name(), test.statusCode, output.LastRequest.StatusCode, test.input)
		}
	}
}

// ExampleGetAddress example using GetAddress()
func ExampleGetAddress() {
	client := newMockClient(defaultUserAgent)
	resp, _ := GetAddress(client, "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA")
	fmt.Println(resp.Address)
	// Output:16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA
}

// TestHandCashConvert will test the HandCashConvert() method
func TestHandCashConvert(t *testing.T) {
	t.Parallel()

	// Create the list of tests
	var tests = []struct {
		input    string
		beta     bool
		expected string
	}{
		{"$mr-z", false, "mr-z@handcash.io"},
		{"invalid$mr-z", false, "invalid$mr-z"},
		{"$", false, "@handcash.io"},
		{"$", true, "@beta.handcash.io"},
		{"1handle", false, "1handle"},
		{"$misterz", true, "misterz@beta.handcash.io"},
	}

	// Test all
	for _, test := range tests {
		if output := HandCashConvert(test.input, test.beta); output != test.expected {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%s]", t.Name(), test.input, test.expected, output)
		}
	}
}

// ExampleHandCashConvert example using HandCashConvert()
func ExampleHandCashConvert() {
	paymail := HandCashConvert("$mr-z", false)
	fmt.Println(paymail)
	// Output:mr-z@handcash.io
}

// BenchmarkHandCashConvert benchmarks the HandCashConvert method
func BenchmarkHandCashConvert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = HandCashConvert("$mr-z", false)
	}
}

// TestRelayXConvert will test the RelayXConvert() method
func TestRelayXConvert(t *testing.T) {
	t.Parallel()

	// Create the list of tests
	var tests = []struct {
		input    string
		expected string
	}{
		{"1mr-z", "mr-z@relayx.io"},
		{"invalid1mr-z", "invalid1mr-z"},
		{"1", "@relayx.io"},
		{"1mrz", "mrz@relayx.io"},
		{"1handle", "handle@relayx.io"},
		{"1misterz", "misterz@relayx.io"},
		{"mrz@relayx.io", "mrz@relayx.io"},
	}

	// Test all
	for _, test := range tests {
		if output := RelayXConvert(test.input); output != test.expected {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%s]", t.Name(), test.input, test.expected, output)
		}
	}
}

// ExampleHandCashConvert example using RelayXConvert()
func ExampleRelayXConvert() {
	paymail := RelayXConvert("1mr-z")
	fmt.Println(paymail)
	// Output:mr-z@relayx.io
}

// BenchmarkHandCashConvert benchmarks the RelayXConvert method
func BenchmarkRelayXConvert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = RelayXConvert("1mr-z")
	}
}
