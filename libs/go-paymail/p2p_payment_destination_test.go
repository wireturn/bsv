package paymail

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

// TestClient_GetP2PPaymentDestination will test the method GetP2PPaymentDestination()
func TestClient_GetP2PPaymentDestination(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	t.Run("successful response", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		mockP2PPaymentDestination(http.StatusOK)

		paymentRequest := &PaymentRequest{Satoshis: 100}

		var destination *PaymentDestination
		destination, err = client.GetP2PPaymentDestination(
			testServerURL+"p2p-payment-destination/{alias}@{domain.tld}",
			testAlias,
			testDomain,
			paymentRequest,
		)
		assert.NoError(t, err)
		assert.NotNil(t, destination)
		assert.Equal(t, http.StatusOK, destination.StatusCode)
		assert.NotEqual(t, 0, len(destination.Outputs))
		assert.NotEqual(t, 0, len(destination.Reference))
		assert.Equal(t, uint64(100), destination.Outputs[0].Satoshis)
	})

	t.Run("successful response - status not modified", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		mockP2PPaymentDestination(http.StatusNotModified)

		paymentRequest := &PaymentRequest{Satoshis: 100}

		var destination *PaymentDestination
		destination, err = client.GetP2PPaymentDestination(
			testServerURL+"p2p-payment-destination/{alias}@{domain.tld}",
			testAlias,
			testDomain,
			paymentRequest,
		)
		assert.NoError(t, err)
		assert.NotNil(t, destination)
		assert.Equal(t, http.StatusNotModified, destination.StatusCode)
		assert.NotEqual(t, 0, len(destination.Outputs))
		assert.NotEqual(t, 0, len(destination.Reference))
		assert.Equal(t, uint64(100), destination.Outputs[0].Satoshis)
	})

	t.Run("bad url", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		mockP2PPaymentDestination(http.StatusOK)

		paymentRequest := &PaymentRequest{Satoshis: 100}

		var destination *PaymentDestination
		destination, err = client.GetP2PPaymentDestination("invalid-url", testAlias, testDomain, paymentRequest)
		assert.Error(t, err)
		assert.Nil(t, destination)
	})

	t.Run("payment is nil", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		mockP2PPaymentDestination(http.StatusOK)

		var destination *PaymentDestination
		destination, err = client.GetP2PPaymentDestination(
			testServerURL+"p2p-payment-destination/{alias}@{domain.tld}",
			testAlias,
			testDomain,
			nil,
		)
		assert.Error(t, err)
		assert.Nil(t, destination)
	})

	t.Run("missing alias", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		mockP2PPaymentDestination(http.StatusOK)

		paymentRequest := &PaymentRequest{Satoshis: 100}

		var destination *PaymentDestination
		destination, err = client.GetP2PPaymentDestination(
			testServerURL+"p2p-payment-destination/{alias}@{domain.tld}",
			"",
			testDomain,
			paymentRequest,
		)
		assert.Error(t, err)
		assert.Nil(t, destination)
	})

	t.Run("missing domain", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		mockP2PPaymentDestination(http.StatusOK)

		paymentRequest := &PaymentRequest{Satoshis: 100}

		var destination *PaymentDestination
		destination, err = client.GetP2PPaymentDestination(
			testServerURL+"p2p-payment-destination/{alias}@{domain.tld}",
			testAlias,
			"",
			paymentRequest,
		)
		assert.Error(t, err)
		assert.Nil(t, destination)
	})

	t.Run("missing satoshis", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		mockP2PPaymentDestination(http.StatusOK)

		paymentRequest := &PaymentRequest{Satoshis: 0}

		var destination *PaymentDestination
		destination, err = client.GetP2PPaymentDestination(
			testServerURL+"p2p-payment-destination/{alias}@{domain.tld}",
			testAlias,
			testDomain,
			paymentRequest,
		)
		assert.Error(t, err)
		assert.Nil(t, destination)
	})

	t.Run("bad request", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodPost, testServerURL+"p2p-payment-destination/"+testAlias+"@"+testDomain,
			httpmock.NewStringResponder(
				http.StatusBadRequest,
				`{"message": "request failed"}`,
			),
		)

		paymentRequest := &PaymentRequest{Satoshis: 100}

		var destination *PaymentDestination
		destination, err = client.GetP2PPaymentDestination(
			testServerURL+"p2p-payment-destination/{alias}@{domain.tld}",
			testAlias,
			testDomain,
			paymentRequest,
		)
		assert.Error(t, err)
		assert.NotNil(t, destination)
		assert.Equal(t, http.StatusBadRequest, destination.StatusCode)
	})

	t.Run("http error", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodPost, testServerURL+"p2p-payment-destination/"+testAlias+"@"+testDomain,
			httpmock.NewErrorResponder(fmt.Errorf("error in request")),
		)

		paymentRequest := &PaymentRequest{Satoshis: 100}

		var destination *PaymentDestination
		destination, err = client.GetP2PPaymentDestination(
			testServerURL+"p2p-payment-destination/{alias}@{domain.tld}",
			testAlias,
			testDomain,
			paymentRequest,
		)
		assert.Error(t, err)
		assert.Nil(t, destination)
	})

	t.Run("address not found", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodPost, testServerURL+"p2p-payment-destination/"+testAlias+"@"+testDomain,
			httpmock.NewStringResponder(
				http.StatusNotFound,
				`{"message": "not found"}`,
			),
		)

		paymentRequest := &PaymentRequest{Satoshis: 100}

		var destination *PaymentDestination
		destination, err = client.GetP2PPaymentDestination(
			testServerURL+"p2p-payment-destination/{alias}@{domain.tld}",
			testAlias,
			testDomain,
			paymentRequest,
		)
		assert.Error(t, err)
		assert.NotNil(t, destination)
		assert.Equal(t, http.StatusNotFound, destination.StatusCode)
	})

	t.Run("error in json", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodPost, testServerURL+"p2p-payment-destination/"+testAlias+"@"+testDomain,
			httpmock.NewStringResponder(
				http.StatusBadRequest,
				`{"message": request failed}`,
			),
		)

		paymentRequest := &PaymentRequest{Satoshis: 100}

		var destination *PaymentDestination
		destination, err = client.GetP2PPaymentDestination(
			testServerURL+"p2p-payment-destination/{alias}@{domain.tld}",
			testAlias,
			testDomain,
			paymentRequest,
		)
		assert.Error(t, err)
		assert.NotNil(t, destination)
		assert.Equal(t, http.StatusBadRequest, destination.StatusCode)
	})

	t.Run("invalid json", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodPost, testServerURL+"p2p-payment-destination/"+testAlias+"@"+testDomain,
			httpmock.NewStringResponder(
				http.StatusOK,
				`{"outputs": [{script: 76a9143e2d1d795f8acaa7957045cc59376177eb04a3c588ac","satoshis": 100}],"reference": "z0bac4ec-6f15-42de-9ef4-e60bfdabf4f7"}`,
			),
		)

		paymentRequest := &PaymentRequest{Satoshis: 100}

		var destination *PaymentDestination
		destination, err = client.GetP2PPaymentDestination(
			testServerURL+"p2p-payment-destination/{alias}@{domain.tld}",
			testAlias,
			testDomain,
			paymentRequest,
		)
		assert.Error(t, err)
		assert.NotNil(t, destination)
		assert.Equal(t, http.StatusOK, destination.StatusCode)
		assert.Equal(t, 0, len(destination.Outputs))
	})

	t.Run("missing reference", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodPost, testServerURL+"p2p-payment-destination/"+testAlias+"@"+testDomain,
			httpmock.NewStringResponder(
				http.StatusOK,
				`{"outputs": [{"script": "76a9143e2d1d795f8acaa7957045cc59376177eb04a3c588ac","satoshis": 100}],"reference": ""}`,
			),
		)

		paymentRequest := &PaymentRequest{Satoshis: 100}

		var destination *PaymentDestination
		destination, err = client.GetP2PPaymentDestination(
			testServerURL+"p2p-payment-destination/{alias}@{domain.tld}",
			testAlias,
			testDomain,
			paymentRequest,
		)
		assert.Error(t, err)
		assert.NotNil(t, destination)
		assert.Equal(t, http.StatusOK, destination.StatusCode)
		assert.Equal(t, 0, len(destination.Reference))
	})

	t.Run("missing outputs", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodPost, testServerURL+"p2p-payment-destination/"+testAlias+"@"+testDomain,
			httpmock.NewStringResponder(
				http.StatusOK,
				`{"outputs": [],"reference": "12345678"}`,
			),
		)

		paymentRequest := &PaymentRequest{Satoshis: 100}

		var destination *PaymentDestination
		destination, err = client.GetP2PPaymentDestination(
			testServerURL+"p2p-payment-destination/{alias}@{domain.tld}",
			testAlias,
			testDomain,
			paymentRequest,
		)
		assert.Error(t, err)
		assert.NotNil(t, destination)
		assert.Equal(t, http.StatusOK, destination.StatusCode)
		assert.Equal(t, 0, len(destination.Outputs))
	})

	t.Run("invalid script", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodPost, testServerURL+"p2p-payment-destination/"+testAlias+"@"+testDomain,
			httpmock.NewStringResponder(
				http.StatusOK,
				`{"outputs": [{"script": "12345678","satoshis": 100}],"reference": "12345678"}`,
			),
		)

		paymentRequest := &PaymentRequest{Satoshis: 100}

		var destination *PaymentDestination
		destination, err = client.GetP2PPaymentDestination(
			testServerURL+"p2p-payment-destination/{alias}@{domain.tld}",
			testAlias,
			testDomain,
			paymentRequest,
		)
		assert.Error(t, err)
		assert.NotNil(t, destination)
		assert.Equal(t, http.StatusOK, destination.StatusCode)
		assert.Equal(t, 0, len(destination.Outputs[0].Address))
	})

	t.Run("empty script", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodPost, testServerURL+"p2p-payment-destination/"+testAlias+"@"+testDomain,
			httpmock.NewStringResponder(
				http.StatusOK,
				`{"outputs": [{"script": "","satoshis": 100}],"reference": "12345678"}`,
			),
		)

		paymentRequest := &PaymentRequest{Satoshis: 100}

		var destination *PaymentDestination
		destination, err = client.GetP2PPaymentDestination(
			testServerURL+"p2p-payment-destination/{alias}@{domain.tld}",
			testAlias,
			testDomain,
			paymentRequest,
		)
		assert.Error(t, err)
		assert.NotNil(t, destination)
		assert.Equal(t, http.StatusOK, destination.StatusCode)
		assert.Equal(t, 0, len(destination.Outputs[0].Address))
		assert.Equal(t, 0, len(destination.Outputs[0].Script))
	})

	t.Run("invalid hex in script", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodPost, testServerURL+"p2p-payment-destination/"+testAlias+"@"+testDomain,
			httpmock.NewStringResponder(
				http.StatusOK,
				`{"outputs": [{"script": "0","satoshis": 100}],"reference": "12345678"}`,
			),
		)

		paymentRequest := &PaymentRequest{Satoshis: 100}

		var destination *PaymentDestination
		destination, err = client.GetP2PPaymentDestination(
			testServerURL+"p2p-payment-destination/{alias}@{domain.tld}",
			testAlias,
			testDomain,
			paymentRequest,
		)
		assert.Error(t, err)
		assert.NotNil(t, destination)
		assert.Equal(t, http.StatusOK, destination.StatusCode)
		assert.Equal(t, 0, len(destination.Outputs[0].Address))
		assert.Equal(t, 1, len(destination.Outputs[0].Script))
	})
}

// mockP2PPaymentDestination is used for mocking the response
func mockP2PPaymentDestination(statusCode int) {
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, testServerURL+"p2p-payment-destination/"+testAlias+"@"+testDomain,
		httpmock.NewStringResponder(
			statusCode,
			`{"outputs": [{"script": "76a9143e2d1d795f8acaa7957045cc59376177eb04a3c588ac","satoshis": 100}],"reference": "z0bac4ec-6f15-42de-9ef4-e60bfdabf4f7"}`,
		),
	)
}

// ExampleClient_GetP2PPaymentDestination example using GetP2PPaymentDestination()
//
// See more examples in /examples/
func ExampleClient_GetP2PPaymentDestination() {
	// Load the client
	client, err := newTestClient()
	if err != nil {
		fmt.Printf("error loading client: %s", err.Error())
		return
	}

	mockP2PPaymentDestination(http.StatusOK)

	// Payment Request
	paymentRequest := &PaymentRequest{Satoshis: 100}

	// Fire the request
	var destination *PaymentDestination
	destination, err = client.GetP2PPaymentDestination(
		testServerURL+"p2p-payment-destination/{alias}@{domain.tld}",
		testAlias,
		testDomain,
		paymentRequest,
	)
	if err != nil {
		fmt.Printf("error occurred in GetP2PPaymentDestination: %s", err.Error())
		return
	}
	if len(destination.Outputs) > 0 {
		fmt.Printf("payment destination: " + destination.Outputs[0].Script)
	}
	// Output:payment destination: 76a9143e2d1d795f8acaa7957045cc59376177eb04a3c588ac
}

// BenchmarkClient_GetP2PPaymentDestination benchmarks the method GetP2PPaymentDestination()
func BenchmarkClient_GetP2PPaymentDestination(b *testing.B) {
	client, _ := newTestClient()
	mockP2PPaymentDestination(http.StatusOK)

	// Payment Request
	paymentRequest := &PaymentRequest{Satoshis: 100}

	for i := 0; i < b.N; i++ {
		_, _ = client.GetP2PPaymentDestination(""+
			testServerURL+"p2p-payment-destination/{alias}@{domain.tld}",
			testAlias,
			testDomain,
			paymentRequest,
		)
	}
}
