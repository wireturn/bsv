package paymail

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestClient_CheckSSL will test the method CheckSSL()
func TestClient_CheckSSL(t *testing.T) {
	// t.Parallel() cannot use newTestClient() race condition

	// Integration test (requires internet connection)
	if testing.Short() {
		t.Skip("skipping integration testing in short mode")
	}

	client, err := newTestClient()
	assert.NoError(t, err)
	assert.NotNil(t, client)

	t.Run("valid ssl certs", func(t *testing.T) {
		var tests = []struct {
			host string
		}{
			{"google.com"},
			{"mrz1818.com"},
		}
		var valid bool
		for _, test := range tests {
			t.Run("checking: "+test.host, func(t *testing.T) {
				valid, err = client.CheckSSL(test.host)
				assert.NoError(t, err)
				assert.Equal(t, true, valid)
			})
		}
	})

	t.Run("invalid ssl certs", func(t *testing.T) {
		var tests = []struct {
			host string
		}{
			{"google"},
			{""},
			{"domaindoesntexistatall101910.co"},
		}
		var valid bool
		for _, test := range tests {
			t.Run("checking: "+test.host, func(t *testing.T) {
				valid, err = client.CheckSSL(test.host)
				assert.Error(t, err)
				assert.Equal(t, false, valid)
			})
		}
	})
}

// ExampleClient_CheckSSL example using CheckSSL()
//
// See more examples in /examples/
func ExampleClient_CheckSSL() {
	client, _ := NewClient()
	valid, _ := client.CheckSSL("google.com")
	if valid {
		fmt.Printf("valid SSL certificate found for: %s", "google.com")
	} else {
		fmt.Printf("invalid SSL certificate found for: %s", "google.com")
	}

	// Output:valid SSL certificate found for: google.com
}

// BenchmarkClient_CheckSSL benchmarks the method CheckSSL()
func BenchmarkClient_CheckSSL(b *testing.B) {
	client, _ := NewClient(nil, nil, nil)
	for i := 0; i < b.N; i++ {
		_, _ = client.CheckSSL("google.com")
	}
}
