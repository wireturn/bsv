package preev

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockHTTPPairsValid for mocking requests
type mockHTTPPairsValid struct{}

// Do is a mock http request
func (m *mockHTTPPairsValid) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Valid
	if req.URL.String() == apiEndpoint+"/pairs" {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"BSV:USD":{"id":"12eLTxv1vyUeJtp5zqWbqpdWvfLdZ7dGf8","name":"Bitcoin SV","base":"BSV","quote":"USD","sources":{"Bitfinex":"bsvusd","Bittrex":"USD-BSV","OkCoin":"bsv_usd","Poloniex":"USDC_BCHSV"},"status":{"active":true,"balance":2.16774049,"tx_size":438,"total_broadcasts":338883,"last_funded":1591604580,"days_remaining":687.4}},"BSV:EUR":{"id":"184UFRtgg8BR2Uebt3BVBvqBLoj81ZwDK1","name":"Bitcoin SV","base":"BSV","quote":"EUR","sources":{"Bitvavo":"BSV-EUR"},"status":{"active":false,"balance":0,"days_remaining":0,"tx_size":341,"total_broadcasts":225782,"last_funded":1585981860}},"BSV:AUD":{"id":"1NsdpaqdWcFrzsaDEM55t7372tEmbrTsZx","name":"Bitcoin SV","base":"BSV","quote":"AUD","sources":{"BTCMarkets":"BSV-AUD"},"status":{"active":false,"balance":0,"days_remaining":0,"tx_size":344,"total_broadcasts":220943,"last_funded":1585868820}},"BSV:BTC":{"id":"1LtHpUkTMULQKNnA4ReAY3EXBXyeMv5U2r","name":"Bitcoin SV","base":"BSV","quote":"BTC","sources":{"Bitfinex":"bsvbtc","Bittrex":"BTC-BSV","Poloniex":"BTC_BCHSV"},"status":{"active":false,"balance":0,"days_remaining":0,"tx_size":418,"total_broadcasts":232885,"last_funded":1587570180}}}`)))
	}

	// Valid (by id)
	if req.URL.String() == apiEndpoint+"/pairs/12eLTxv1vyUeJtp5zqWbqpdWvfLdZ7dGf8" {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"id":"12eLTxv1vyUeJtp5zqWbqpdWvfLdZ7dGf8","name":"Bitcoin SV","base":"BSV","quote":"USD","sources":{"Bitfinex":"bsvusd","Bittrex":"USD-BSV","OkCoin":"bsv_usd","Poloniex":"USDC_BCHSV"},"status":{"active":true,"balance":2.16774049,"tx_size":438,"total_broadcasts":338883,"last_funded":1591604580,"days_remaining":687.4}}`)))
	}

	// Invalid (by id)
	if req.URL.String() == apiEndpoint+"/pairs/bad-pair-id" {
		resp.StatusCode = http.StatusNotFound
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"success": false,"message": "Not Found"}`)))
		return resp, fmt.Errorf(`{"success": false,"message": "Not Found"}`)
	}

	// Invalid (by id)
	if req.URL.String() == apiEndpoint+"/pairs/bad-status-code" {
		resp.StatusCode = http.StatusExpectationFailed
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
	}

	// Default is valid
	return resp, nil
}

// mockHTTPPairsInvalid for mocking requests
type mockHTTPPairsInvalid struct{}

// Do is a mock http request
func (m *mockHTTPPairsInvalid) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Invalid
	if req.URL.String() == apiEndpoint+"/pairs" {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"success": false,"message": "Not Found"}`)))
		return resp, fmt.Errorf(`{"success": false,"message": "Not Found"}`)
	}

	// Default is valid
	return resp, nil
}

// TestClient_GetPairs tests the GetPairs()
func TestClient_GetPairs(t *testing.T) {
	t.Parallel()

	t.Run("valid case", func(t *testing.T) {
		client := newMockClient(&mockHTTPPairsValid{})

		pairs, err := client.GetPairs(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, pairs)
		assert.Equal(t, "12eLTxv1vyUeJtp5zqWbqpdWvfLdZ7dGf8", pairs.BsvUsd.ID)
	})

	t.Run("invalid case", func(t *testing.T) {
		client := newMockClient(&mockHTTPPairsInvalid{})
		assert.NotNil(t, client)

		pairs, err := client.GetPairs(context.Background())
		assert.Error(t, err)
		assert.Nil(t, pairs)
	})
}

// TestClient_GetPair tests the GetPair()
func TestClient_GetPair(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPPairsValid{})

	// Create the list of tests
	var tests = []struct {
		input         string
		expected      string
		expectedError bool
		statusCode    int
	}{
		{"12eLTxv1vyUeJtp5zqWbqpdWvfLdZ7dGf8", "12eLTxv1vyUeJtp5zqWbqpdWvfLdZ7dGf8", false, http.StatusOK},
		{"bad-pair-id", "", true, http.StatusNotFound},
		{"bad-status-code", "", true, http.StatusExpectationFailed},
	}

	// Test all
	for _, test := range tests {
		if output, err := client.GetPair(context.Background(), test.input); err == nil && test.expectedError {
			t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted", t.Name(), test.input)
		} else if err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted, received: [%v] error [%s]", t.Name(), test.input, output, err.Error())
		} else if output != nil && output.ID != test.expected && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%s]", t.Name(), test.input, test.expected, output.ID)
		} else if client.LastRequest.StatusCode != test.statusCode {
			t.Errorf("%s Expected status code to be %d, got %d, [%s] inputted", t.Name(), test.statusCode, client.LastRequest.StatusCode, test.input)
		}
	}

}
