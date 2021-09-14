package paymail

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

// TestClient_VerifyPubKey will test the method VerifyPubKey()
func TestClient_VerifyPubKey(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	t.Run("successful response", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		mockVerifyPubKey(http.StatusOK)

		var verification *Verification
		verification, err = client.VerifyPubKey(
			testServerURL+"verifypubkey/{alias}@{domain.tld}/{pubkey}",
			testAlias, testDomain, testPubKey,
		)
		assert.NoError(t, err)
		assert.NotNil(t, verification)
		assert.Equal(t, DefaultBsvAliasVersion, verification.BsvAlias)
		assert.Equal(t, http.StatusOK, verification.StatusCode)
		assert.Equal(t, testAlias+"@"+testDomain, verification.Handle)
		assert.Equal(t, testPubKey, verification.PubKey)
		assert.Equal(t, true, verification.Match)
	})

	t.Run("successful response - status not modified", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		mockVerifyPubKey(http.StatusNotModified)

		var verification *Verification
		verification, err = client.VerifyPubKey(
			testServerURL+"verifypubkey/{alias}@{domain.tld}/{pubkey}",
			testAlias, testDomain, testPubKey,
		)
		assert.NoError(t, err)
		assert.NotNil(t, verification)
		assert.Equal(t, DefaultBsvAliasVersion, verification.BsvAlias)
		assert.Equal(t, http.StatusNotModified, verification.StatusCode)
		assert.Equal(t, testAlias+"@"+testDomain, verification.Handle)
		assert.Equal(t, testPubKey, verification.PubKey)
		assert.Equal(t, true, verification.Match)
	})

	t.Run("missing url", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		mockVerifyPubKey(http.StatusNotModified)

		var verification *Verification
		verification, err = client.VerifyPubKey(
			"invalid-url",
			testAlias, testDomain, testPubKey,
		)
		assert.Error(t, err)
		assert.Nil(t, verification)
	})

	t.Run("missing alias", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		mockVerifyPubKey(http.StatusNotModified)

		var verification *Verification
		verification, err = client.VerifyPubKey(
			testServerURL+"verifypubkey/{alias}@{domain.tld}/{pubkey}",
			"", testDomain, testPubKey,
		)
		assert.Error(t, err)
		assert.Nil(t, verification)
	})

	t.Run("missing domain", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		mockVerifyPubKey(http.StatusNotModified)

		var verification *Verification
		verification, err = client.VerifyPubKey(
			testServerURL+"verifypubkey/{alias}@{domain.tld}/{pubkey}",
			testAlias, "", testPubKey,
		)
		assert.Error(t, err)
		assert.Nil(t, verification)
	})

	t.Run("missing pubkey", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		mockVerifyPubKey(http.StatusNotModified)

		var verification *Verification
		verification, err = client.VerifyPubKey(
			testServerURL+"verifypubkey/{alias}@{domain.tld}/{pubkey}",
			testAlias, testDomain, "",
		)
		assert.Error(t, err)
		assert.Nil(t, verification)
	})

	t.Run("bad request", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, testServerURL+"verifypubkey/"+testAlias+"@"+testDomain+"/"+testPubKey,
			httpmock.NewStringResponder(
				http.StatusBadRequest,
				`{"message": "request failed"}`,
			),
		)

		var verification *Verification
		verification, err = client.VerifyPubKey(
			testServerURL+"verifypubkey/{alias}@{domain.tld}/{pubkey}",
			testAlias, testDomain, testPubKey,
		)
		assert.Error(t, err)
		assert.NotNil(t, verification)
		assert.Equal(t, http.StatusBadRequest, verification.StatusCode)
	})

	t.Run("http error", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, testServerURL+"verifypubkey/"+testAlias+"@"+testDomain+"/"+testPubKey,
			httpmock.NewErrorResponder(fmt.Errorf("error in request")),
		)

		var verification *Verification
		verification, err = client.VerifyPubKey(
			testServerURL+"verifypubkey/{alias}@{domain.tld}/{pubkey}",
			testAlias, testDomain, testPubKey,
		)
		assert.Error(t, err)
		assert.Nil(t, verification)
	})

	t.Run("bad error", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, testServerURL+"verifypubkey/"+testAlias+"@"+testDomain+"/"+testPubKey,
			httpmock.NewStringResponder(
				http.StatusBadRequest,
				`{"message": request failed}`,
			),
		)

		var verification *Verification
		verification, err = client.VerifyPubKey(
			testServerURL+"verifypubkey/{alias}@{domain.tld}/{pubkey}",
			testAlias, testDomain, testPubKey,
		)
		assert.Error(t, err)
		assert.NotNil(t, verification)
		assert.Equal(t, http.StatusBadRequest, verification.StatusCode)
	})

	t.Run("bad json", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, testServerURL+"verifypubkey/"+testAlias+"@"+testDomain+"/"+testPubKey,
			httpmock.NewStringResponder(
				http.StatusOK,
				`{"bsvalias": 1.0,handle: `+testAlias+`@`+testDomain+`","pubkey": "`+testPubKey+`","match": true}`,
			),
		)

		var verification *Verification
		verification, err = client.VerifyPubKey(
			testServerURL+"verifypubkey/{alias}@{domain.tld}/{pubkey}",
			testAlias, testDomain, testPubKey,
		)
		assert.Error(t, err)
		assert.NotNil(t, verification)
		assert.Equal(t, http.StatusOK, verification.StatusCode)
		assert.Equal(t, "", verification.BsvAlias)
		assert.Equal(t, "", verification.PubKey)
		assert.Equal(t, "", verification.Handle)
	})

	t.Run("invalid handle", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, testServerURL+"verifypubkey/"+testAlias+"@"+testDomain+"/"+testPubKey,
			httpmock.NewStringResponder(
				http.StatusOK,
				`{"bsvalias": "`+DefaultBsvAliasVersion+`","handle": "","pubkey": "`+testPubKey+`","match": true}`,
			),
		)

		var verification *Verification
		verification, err = client.VerifyPubKey(
			testServerURL+"verifypubkey/{alias}@{domain.tld}/{pubkey}",
			testAlias, testDomain, testPubKey,
		)
		assert.Error(t, err)
		assert.NotNil(t, verification)
		assert.Equal(t, http.StatusOK, verification.StatusCode)
		assert.Equal(t, "", verification.Handle)
	})

	t.Run("invalid bsv alias", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, testServerURL+"verifypubkey/"+testAlias+"@"+testDomain+"/"+testPubKey,
			httpmock.NewStringResponder(
				http.StatusOK,
				`{"bsvalias": "","handle": "`+testAlias+`@`+testDomain+`","pubkey": "`+testPubKey+`","match": true}`,
			),
		)

		var verification *Verification
		verification, err = client.VerifyPubKey(
			testServerURL+"verifypubkey/{alias}@{domain.tld}/{pubkey}",
			testAlias, testDomain, testPubKey,
		)
		assert.Error(t, err)
		assert.NotNil(t, verification)
		assert.Equal(t, http.StatusOK, verification.StatusCode)
		assert.Equal(t, "", verification.BsvAlias)
	})

	t.Run("empty pubkey", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, testServerURL+"verifypubkey/"+testAlias+"@"+testDomain+"/"+testPubKey,
			httpmock.NewStringResponder(
				http.StatusOK,
				`{"bsvalias": "`+DefaultBsvAliasVersion+`","handle": "`+testAlias+`@`+testDomain+`","pubkey": "","match": true}`,
			),
		)

		var verification *Verification
		verification, err = client.VerifyPubKey(
			testServerURL+"verifypubkey/{alias}@{domain.tld}/{pubkey}",
			testAlias, testDomain, testPubKey,
		)
		assert.Error(t, err)
		assert.NotNil(t, verification)
		assert.Equal(t, http.StatusOK, verification.StatusCode)
		assert.Equal(t, "", verification.PubKey)
	})

	t.Run("invalid pubkey", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, testServerURL+"verifypubkey/"+testAlias+"@"+testDomain+"/"+testPubKey,
			httpmock.NewStringResponder(
				http.StatusOK,
				`{"bsvalias": "`+DefaultBsvAliasVersion+`","handle": "`+testAlias+`@`+testDomain+`","pubkey": "wrong","match": true}`,
			),
		)

		var verification *Verification
		verification, err = client.VerifyPubKey(
			testServerURL+"verifypubkey/{alias}@{domain.tld}/{pubkey}",
			testAlias, testDomain, testPubKey,
		)
		assert.Error(t, err)
		assert.NotNil(t, verification)
		assert.Equal(t, http.StatusOK, verification.StatusCode)
		assert.NotEqual(t, testPubKey, verification.PubKey)
	})
}

// mockVerifyPubKey is used for mocking the response
func mockVerifyPubKey(statusCode int) {
	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodGet,
		testServerURL+"verifypubkey/"+testAlias+"@"+testDomain+"/02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10",
		httpmock.NewStringResponder(
			statusCode,
			`{"`+DefaultServiceName+`": "`+DefaultBsvAliasVersion+`",
"handle": "`+testAlias+`@`+testDomain+`",
"pubkey": "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10","match": true}`,
		),
	)
}

// ExampleClient_VerifyPubKey example using VerifyPubKey()
//
// See more examples in /examples/
func ExampleClient_VerifyPubKey() {
	// Load the client
	client, err := newTestClient()
	if err != nil {
		fmt.Printf("error loading client: %s", err.Error())
		return
	}

	mockVerifyPubKey(http.StatusOK)

	// Verify PubKey
	var verification *Verification
	verification, err = client.VerifyPubKey(
		testServerURL+"verifypubkey/{alias}@{domain.tld}/{pubkey}",
		testAlias, testDomain, testPubKey,
	)
	if err != nil {
		fmt.Printf("error getting verification: " + err.Error())
		return
	}
	fmt.Printf("verified %s handle with pubkey: %s", verification.Handle, verification.PubKey)
	// Output:verified mrz@test.com handle with pubkey: 02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10
}

// BenchmarkClient_VerifyPubKey benchmarks the method VerifyPubKey()
func BenchmarkClient_VerifyPubKey(b *testing.B) {
	client, _ := newTestClient()
	mockVerifyPubKey(http.StatusOK)
	for i := 0; i < b.N; i++ {
		_, _ = client.VerifyPubKey(
			testServerURL+"verifypubkey/{alias}@{domain.tld}/{pubkey}",
			testAlias, testDomain, testPubKey,
		)
	}
}
