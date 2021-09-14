package paymail

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

// TestClient_ResolveAddress will test the method ResolveAddress()
func TestClient_ResolveAddress(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	t.Run("successful response", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		mockResolveAddress(http.StatusOK)

		senderRequest := &SenderRequest{
			Dt:           time.Now().UTC().Format(time.RFC3339),
			SenderHandle: testAlias + "@" + testDomain,
			SenderName:   testName,
		}

		var resolution *Resolution
		resolution, err = client.ResolveAddress(
			testServerURL+"address/{alias}@{domain.tld}", testAlias, testDomain, senderRequest,
		)
		assert.NoError(t, err)
		assert.NotNil(t, resolution)
		assert.Equal(t, http.StatusOK, resolution.StatusCode)
		assert.Equal(t, testAddress, resolution.Address)
		assert.Equal(t, testOutput, resolution.Output)
	})

	t.Run("successful response - status not modified", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		mockResolveAddress(http.StatusNotModified)

		senderRequest := &SenderRequest{
			Dt:           time.Now().UTC().Format(time.RFC3339),
			SenderHandle: testAlias + "@" + testDomain,
			SenderName:   testName,
		}

		var resolution *Resolution
		resolution, err = client.ResolveAddress(
			testServerURL+"address/{alias}@{domain.tld}", testAlias, testDomain, senderRequest,
		)
		assert.NoError(t, err)
		assert.NotNil(t, resolution)
		assert.Equal(t, http.StatusNotModified, resolution.StatusCode)
		assert.Equal(t, testAddress, resolution.Address)
		assert.Equal(t, testOutput, resolution.Output)
	})

	t.Run("invalid url", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		mockResolveAddress(http.StatusOK)

		senderRequest := &SenderRequest{
			Dt:           time.Now().UTC().Format(time.RFC3339),
			SenderHandle: testAlias + "@" + testDomain,
			SenderName:   testName,
		}

		var resolution *Resolution
		resolution, err = client.ResolveAddress(
			"invalid-url", testAlias, testDomain, senderRequest,
		)
		assert.Error(t, err)
		assert.Nil(t, resolution)
	})

	t.Run("sender request is nil", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		mockResolveAddress(http.StatusOK)

		var resolution *Resolution
		resolution, err = client.ResolveAddress(
			testServerURL+"address/{alias}@{domain.tld}", testAlias, testDomain, nil,
		)
		assert.Error(t, err)
		assert.Nil(t, resolution)
	})

	t.Run("sender request - dt invalid", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		mockResolveAddress(http.StatusOK)

		senderRequest := &SenderRequest{
			Dt:           "",
			SenderHandle: testAlias + "@" + testDomain,
			SenderName:   testName,
		}

		var resolution *Resolution
		resolution, err = client.ResolveAddress(
			testServerURL+"address/{alias}@{domain.tld}", testAlias, testDomain, senderRequest,
		)
		assert.Error(t, err)
		assert.Nil(t, resolution)
	})

	t.Run("sender request - sender handle invalid", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		mockResolveAddress(http.StatusOK)

		senderRequest := &SenderRequest{
			Dt:           time.Now().UTC().Format(time.RFC3339),
			SenderHandle: "",
			SenderName:   testName,
		}

		var resolution *Resolution
		resolution, err = client.ResolveAddress(
			testServerURL+"address/{alias}@{domain.tld}", testAlias, testDomain, senderRequest,
		)
		assert.Error(t, err)
		assert.Nil(t, resolution)
	})

	t.Run("missing alias", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		mockResolveAddress(http.StatusOK)

		senderRequest := &SenderRequest{
			Dt:           time.Now().UTC().Format(time.RFC3339),
			SenderHandle: testAlias + "@" + testDomain,
			SenderName:   testName,
		}

		var resolution *Resolution
		resolution, err = client.ResolveAddress(
			testServerURL+"address/{alias}@{domain.tld}", "", testDomain, senderRequest,
		)
		assert.Error(t, err)
		assert.Nil(t, resolution)
	})

	t.Run("missing domain", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		mockResolveAddress(http.StatusOK)

		senderRequest := &SenderRequest{
			Dt:           time.Now().UTC().Format(time.RFC3339),
			SenderHandle: testAlias + "@" + testDomain,
			SenderName:   testName,
		}

		var resolution *Resolution
		resolution, err = client.ResolveAddress(
			testServerURL+"address/{alias}@{domain.tld}", testAlias, "", senderRequest,
		)
		assert.Error(t, err)
		assert.Nil(t, resolution)
	})

	t.Run("bad request", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodPost, testServerURL+"address/"+testAlias+"@"+testDomain,
			httpmock.NewStringResponder(
				http.StatusBadRequest,
				`{"message": "request failed"}`,
			),
		)

		senderRequest := &SenderRequest{
			Dt:           time.Now().UTC().Format(time.RFC3339),
			SenderHandle: testAlias + "@" + testDomain,
			SenderName:   testName,
		}

		var resolution *Resolution
		resolution, err = client.ResolveAddress(
			testServerURL+"address/{alias}@{domain.tld}", testAlias, testDomain, senderRequest,
		)
		assert.Error(t, err)
		assert.NotNil(t, resolution)
		assert.Equal(t, http.StatusBadRequest, resolution.StatusCode)
	})

	t.Run("http error", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodPost, testServerURL+"address/"+testAlias+"@"+testDomain,
			httpmock.NewErrorResponder(fmt.Errorf("error in request")),
		)

		senderRequest := &SenderRequest{
			Dt:           time.Now().UTC().Format(time.RFC3339),
			SenderHandle: testAlias + "@" + testDomain,
			SenderName:   testName,
		}

		var resolution *Resolution
		resolution, err = client.ResolveAddress(
			testServerURL+"address/{alias}@{domain.tld}", testAlias, testDomain, senderRequest,
		)
		assert.Error(t, err)
		assert.Nil(t, resolution)
	})

	t.Run("bad error", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodPost, testServerURL+"address/"+testAlias+"@"+testDomain,
			httpmock.NewStringResponder(
				http.StatusBadRequest,
				`{"message": request failed}`,
			),
		)

		senderRequest := &SenderRequest{
			Dt:           time.Now().UTC().Format(time.RFC3339),
			SenderHandle: testAlias + "@" + testDomain,
			SenderName:   testName,
		}

		var resolution *Resolution
		resolution, err = client.ResolveAddress(
			testServerURL+"address/{alias}@{domain.tld}", testAlias, testDomain, senderRequest,
		)
		assert.Error(t, err)
		assert.NotNil(t, resolution)
		assert.Equal(t, http.StatusBadRequest, resolution.StatusCode)
	})

	t.Run("paymail not found", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodPost, testServerURL+"address/"+testAlias+"@"+testDomain,
			httpmock.NewStringResponder(
				http.StatusNotFound,
				`{"message": "not found"}`,
			),
		)

		senderRequest := &SenderRequest{
			Dt:           time.Now().UTC().Format(time.RFC3339),
			SenderHandle: testAlias + "@" + testDomain,
			SenderName:   testName,
		}

		var resolution *Resolution
		resolution, err = client.ResolveAddress(
			testServerURL+"address/{alias}@{domain.tld}", testAlias, testDomain, senderRequest,
		)
		assert.Error(t, err)
		assert.NotNil(t, resolution)
		assert.Equal(t, http.StatusNotFound, resolution.StatusCode)
	})

	t.Run("invalid json", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodPost, testServerURL+"address/"+testAlias+"@"+testDomain,
			httpmock.NewStringResponder(
				http.StatusOK,
				`{"output: `+testOutput+`"}`,
			),
		)

		senderRequest := &SenderRequest{
			Dt:           time.Now().UTC().Format(time.RFC3339),
			SenderHandle: testAlias + "@" + testDomain,
			SenderName:   testName,
		}

		var resolution *Resolution
		resolution, err = client.ResolveAddress(
			testServerURL+"address/{alias}@{domain.tld}", testAlias, testDomain, senderRequest,
		)
		assert.Error(t, err)
		assert.NotNil(t, resolution)
		assert.Equal(t, http.StatusOK, resolution.StatusCode)
		assert.Equal(t, "", resolution.Output)
	})

	t.Run("missing output", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodPost, testServerURL+"address/"+testAlias+"@"+testDomain,
			httpmock.NewStringResponder(
				http.StatusOK,
				`{"output": ""}`,
			),
		)

		senderRequest := &SenderRequest{
			Dt:           time.Now().UTC().Format(time.RFC3339),
			SenderHandle: testAlias + "@" + testDomain,
			SenderName:   testName,
		}

		var resolution *Resolution
		resolution, err = client.ResolveAddress(
			testServerURL+"address/{alias}@{domain.tld}", testAlias, testDomain, senderRequest,
		)
		assert.Error(t, err)
		assert.NotNil(t, resolution)
		assert.Equal(t, http.StatusOK, resolution.StatusCode)
		assert.Equal(t, "", resolution.Output)
	})

	t.Run("invalid output", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodPost, testServerURL+"address/"+testAlias+"@"+testDomain,
			httpmock.NewStringResponder(
				http.StatusOK,
				`{"output": "12345678"}`,
			),
		)

		senderRequest := &SenderRequest{
			Dt:           time.Now().UTC().Format(time.RFC3339),
			SenderHandle: testAlias + "@" + testDomain,
			SenderName:   testName,
		}

		var resolution *Resolution
		resolution, err = client.ResolveAddress(
			testServerURL+"address/{alias}@{domain.tld}", testAlias, testDomain, senderRequest,
		)
		assert.Error(t, err)
		assert.NotNil(t, resolution)
		assert.Equal(t, http.StatusOK, resolution.StatusCode)
		assert.Equal(t, "12345678", resolution.Output)
	})

	t.Run("invalid output hex", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodPost, testServerURL+"address/"+testAlias+"@"+testDomain,
			httpmock.NewStringResponder(
				http.StatusOK,
				`{"output": "7e00bb007d4960727eb11d92a052502c"}`,
			),
		)

		senderRequest := &SenderRequest{
			Dt:           time.Now().UTC().Format(time.RFC3339),
			SenderHandle: testAlias + "@" + testDomain,
			SenderName:   testName,
		}

		var resolution *Resolution
		resolution, err = client.ResolveAddress(
			testServerURL+"address/{alias}@{domain.tld}", testAlias, testDomain, senderRequest,
		)
		assert.Error(t, err)
		assert.NotNil(t, resolution)
		assert.Equal(t, http.StatusOK, resolution.StatusCode)
		assert.Equal(t, "7e00bb007d4960727eb11d92a052502c", resolution.Output)
	})

	t.Run("invalid output hex length", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodPost, testServerURL+"address/"+testAlias+"@"+testDomain,
			httpmock.NewStringResponder(
				http.StatusOK,
				`{"output": "0"}`,
			),
		)

		senderRequest := &SenderRequest{
			Dt:           time.Now().UTC().Format(time.RFC3339),
			SenderHandle: testAlias + "@" + testDomain,
			SenderName:   testName,
		}

		var resolution *Resolution
		resolution, err = client.ResolveAddress(
			testServerURL+"address/{alias}@{domain.tld}", testAlias, testDomain, senderRequest,
		)
		assert.Error(t, err)
		assert.NotNil(t, resolution)
		assert.Equal(t, http.StatusOK, resolution.StatusCode)
		assert.Equal(t, "0", resolution.Output)
	})
}

// mockResolveAddress is used for mocking the response
func mockResolveAddress(statusCode int) {
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, testServerURL+"address/"+testAlias+"@"+testDomain,
		httpmock.NewStringResponder(
			statusCode,
			`{"output": "`+testOutput+`"}`,
		),
	)
}

// ExampleClient_ResolveAddress example using ResolveAddress()
//
// See more examples in /examples/
func ExampleClient_ResolveAddress() {
	// Load the client
	client, err := newTestClient()
	if err != nil {
		fmt.Printf("error loading client: %s", err.Error())
		return
	}

	mockResolveAddress(http.StatusOK)

	// Sender Request
	senderRequest := &SenderRequest{
		Dt:           time.Now().UTC().Format(time.RFC3339),
		SenderHandle: testAlias + "@" + testDomain,
		SenderName:   testName,
	}

	// Fire the request
	var resolution *Resolution
	resolution, err = client.ResolveAddress(
		testServerURL+"address/{alias}@{domain.tld}",
		testAlias, testDomain, senderRequest,
	)
	if err != nil {
		fmt.Printf("error occurred in ResolveAddress: %s", err.Error())
		return
	}
	fmt.Printf("address found: %s", resolution.Address)
	// Output:address found: 1Cat862cjhp8SgLLMvin5gyk5UScasg1P9
}

// BenchmarkClient_ResolveAddress benchmarks the method ResolveAddress()
func BenchmarkClient_ResolveAddress(b *testing.B) {
	client, _ := newTestClient()

	// Sender Request
	senderRequest := &SenderRequest{
		Dt:           time.Now().UTC().Format(time.RFC3339),
		SenderHandle: testAlias + "@" + testDomain,
		SenderName:   testName,
	}

	for i := 0; i < b.N; i++ {
		_, _ = client.ResolveAddress(
			testServerURL+"address/{alias}@{domain.tld}",
			testAlias, testDomain, senderRequest,
		)
	}
}
