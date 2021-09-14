package tonicpow

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// newTestRate creates a dummy profile for testing
func newTestRate() *Rate {
	return &Rate{
		Currency:        testRateCurrency,
		CurrencyAmount:  0.01,
		PriceInSatoshis: 4200,
	}
}

// TestClient_GetCurrentRate will test the method GetCurrentRate()
func TestClient_GetCurrentRate(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	t.Run("get current rate (success)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		rates := newTestRate()

		endpoint := fmt.Sprintf(
			"%s/%s/%s?%s=%f", EnvironmentDevelopment.apiURL,
			modelRates, testRateCurrency,
			fieldAmount, 0.0,
		)

		err = mockResponseData(http.MethodGet, endpoint, http.StatusOK, rates)
		assert.NoError(t, err)

		var currentRate *Rate
		var response *StandardResponse
		currentRate, response, err = client.GetCurrentRate(testRateCurrency, 0.00)
		assert.NoError(t, err)
		assert.NotNil(t, currentRate)
		assert.NotNil(t, response)
		assert.Equal(t, testRateCurrency, currentRate.Currency)
	})

	t.Run("missing currency", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		rates := newTestRate()

		endpoint := fmt.Sprintf(
			"%s/%s/%s?%s=%f", EnvironmentDevelopment.apiURL,
			modelRates, testRateCurrency,
			fieldAmount, 0.0,
		)

		err = mockResponseData(http.MethodGet, endpoint, http.StatusOK, rates)
		assert.NoError(t, err)

		var currentRate *Rate
		var response *StandardResponse
		currentRate, response, err = client.GetCurrentRate("", 0.00)
		assert.Error(t, err)
		assert.Nil(t, currentRate)
		assert.Nil(t, response)
	})

	t.Run("error from api (status code)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		rates := newTestRate()

		endpoint := fmt.Sprintf(
			"%s/%s/%s?%s=%f", EnvironmentDevelopment.apiURL,
			modelRates, testRateCurrency,
			fieldAmount, 0.0,
		)

		err = mockResponseData(http.MethodGet, endpoint, http.StatusBadRequest, rates)
		assert.NoError(t, err)

		var currentRate *Rate
		var response *StandardResponse
		currentRate, response, err = client.GetCurrentRate(testRateCurrency, 0.00)
		assert.Error(t, err)
		assert.Nil(t, currentRate)
		assert.NotNil(t, response)
	})

	t.Run("error from api (status code)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		endpoint := fmt.Sprintf(
			"%s/%s/%s?%s=%f", EnvironmentDevelopment.apiURL,
			modelRates, testRateCurrency,
			fieldAmount, 0.0,
		)

		apiError := &Error{
			Code:        400,
			Data:        "field_name",
			IPAddress:   "127.0.0.1",
			Message:     "some error message",
			Method:      http.MethodPut,
			RequestGUID: "7f3d97a8fd67ff57861904df6118dcc8",
			StatusCode:  http.StatusBadRequest,
			URL:         endpoint,
		}

		err = mockResponseData(http.MethodGet, endpoint, http.StatusBadRequest, apiError)
		assert.NoError(t, err)

		var currentRate *Rate
		var response *StandardResponse
		currentRate, response, err = client.GetCurrentRate(testRateCurrency, 0.00)
		assert.Error(t, err)
		assert.Nil(t, currentRate)
		assert.NotNil(t, response)
		assert.Equal(t, apiError.Message, err.Error())
	})
}

// ExampleClient_GetCurrentRate example using GetCurrentRate()
//
// See more examples in /examples/
func ExampleClient_GetCurrentRate() {

	// Load the client (using test client for example only)
	client, err := newTestClient()
	if err != nil {
		fmt.Printf("error loading client: %s", err.Error())
		return
	}

	// For mocking
	rates := newTestRate()

	// Mock response (for example only)
	_ = mockResponseData(
		http.MethodGet,
		fmt.Sprintf(
			"%s/%s/%s?%s=%f", EnvironmentDevelopment.apiURL,
			modelRates, testRateCurrency,
			fieldAmount, 0.0,
		),
		http.StatusOK,
		rates,
	)

	// Get rate (using mocking response)
	var currentRate *Rate
	if currentRate, _, err = client.GetCurrentRate(
		testRateCurrency, 0.00,
	); err != nil {
		fmt.Printf("error getting profile: " + err.Error())
		return
	}
	fmt.Printf("current rate: %s  %f usd is %d sats", currentRate.Currency, currentRate.CurrencyAmount, currentRate.PriceInSatoshis)
	// Output:current rate: usd  0.010000 usd is 4200 sats
}

// BenchmarkClient_GetCurrentRate benchmarks the method GetCurrentRate()
func BenchmarkClient_GetCurrentRate(b *testing.B) {
	client, _ := newTestClient()
	rate := newTestRate()
	_ = mockResponseData(
		http.MethodGet,
		fmt.Sprintf(
			"%s/%s/%s?%s=%f", EnvironmentDevelopment.apiURL,
			modelRates, testRateCurrency,
			fieldAmount, 0.0,
		),
		http.StatusOK,
		rate,
	)
	for i := 0; i < b.N; i++ {
		_, _, _ = client.GetCurrentRate(testRateCurrency, 0.00)
	}
}
