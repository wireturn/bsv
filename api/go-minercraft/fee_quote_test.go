package minercraft

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

const feeTestSignature = "3045022100eed49f6bf75d8f975f581271e3df658fbe8ec67e6301ea8fc25a72d18c92e30e022056af253f0d24db6a8fde4e2c1ee95e7a5ecf2c7cdc93246f8328c9e0ca582fc4"
const feeTestPublicKey = "03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270"

// mockHTTPValidFeeQuote for mocking requests
type mockHTTPValidFeeQuote struct{}

// Do is a mock http request
func (m *mockHTTPValidFeeQuote) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Valid response
	if strings.Contains(req.URL.String(), "/mapi/feeQuote") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{
    	"payload": "{\"apiVersion\":\"` + testAPIVersion + `\",\"timestamp\":\"2020-10-09T21:26:17.410Z\",\"expiryTime\":\"2020-10-09T21:36:17.410Z\",\"minerId\":\"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270\",\"currentHighestBlockHash\":\"0000000000000000035c5f8c0294802a01e500fa7b95337963bb3640da3bd565\",\"currentHighestBlockHeight\":656169,\"minerReputation\":null,\"fees\":[{\"id\":1,\"feeType\":\"standard\",\"miningFee\":{\"satoshis\":500,\"bytes\":1000},\"relayFee\":{\"satoshis\":250,\"bytes\":1000}},{\"id\":2,\"feeType\":\"data\",\"miningFee\":{\"satoshis\":500,\"bytes\":1000},\"relayFee\":{\"satoshis\":250,\"bytes\":1000}}]}",
   	 	"signature": "3045022100eed49f6bf75d8f975f581271e3df658fbe8ec67e6301ea8fc25a72d18c92e30e022056af253f0d24db6a8fde4e2c1ee95e7a5ecf2c7cdc93246f8328c9e0ca582fc4",
    	"publicKey": "03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270","encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	// Default is valid
	return resp, nil
}

// mockHTTPError for mocking requests
type mockHTTPError struct{}

// Do is a mock http request
func (m *mockHTTPError) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	if req.URL.String() != "" {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
		return resp, fmt.Errorf(`http timeout`)
	}

	// Default is valid
	return resp, nil
}

// mockHTTPBadRequest for mocking requests
type mockHTTPBadRequest struct{}

// Do is a mock http request
func (m *mockHTTPBadRequest) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	if req.URL.String() != "" {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
	}

	// Default is valid
	return resp, nil
}

// mockHTTPInvalidJSON for mocking requests
type mockHTTPInvalidJSON struct{}

// Do is a mock http request
func (m *mockHTTPInvalidJSON) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	if req.URL.String() != "" {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{invalid:json}`)))
		resp.StatusCode = http.StatusOK
	}

	// Default is valid
	return resp, nil
}

// mockHTTPMissingFees for mocking requests
type mockHTTPMissingFees struct{}

// Do is a mock http request
func (m *mockHTTPMissingFees) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Valid response
	if strings.Contains(req.URL.String(), "/mapi/feeQuote") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{
    	"payload": "{\"apiVersion\":\"` + testAPIVersion + `\",\"timestamp\":\"2020-10-09T21:26:17.410Z\",\"expiryTime\":\"2020-10-09T21:36:17.410Z\",\"minerId\":\"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270\",\"currentHighestBlockHash\":\"0000000000000000035c5f8c0294802a01e500fa7b95337963bb3640da3bd565\",\"currentHighestBlockHeight\":656169,\"minerReputation\":null,\"fees\":[]}",
   	 	"signature": "3045022100eed49f6bf75d8f975f581271e3df658fbe8ec67e6301ea8fc25a72d18c92e30e022056af253f0d24db6a8fde4e2c1ee95e7a5ecf2c7cdc93246f8328c9e0ca582fc4",
    	"publicKey": "03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270","encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	// Default is valid
	return resp, nil
}

// mockHTTPInvalidSignature for mocking requests
type mockHTTPInvalidSignature struct{}

// Do is a mock http request
func (m *mockHTTPInvalidSignature) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Invalid sig response
	if strings.Contains(req.URL.String(), "/mapi/feeQuote") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{
    	"payload": "{\"apiVersion\":\"` + testAPIVersion + `\",\"timestamp\":\"2020-10-09T21:26:17.410Z\",\"expiryTime\":\"2020-10-09T21:36:17.410Z\",\"minerId\":\"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270\",\"currentHighestBlockHash\":\"0000000000000000035c5f8c0294802a01e500fa7b95337963bb3640da3bd565\",\"currentHighestBlockHeight\":656169,\"minerReputation\":null,\"fees\":[{\"id\":1,\"feeType\":\"standard\",\"miningFee\":{\"satoshis\":500,\"bytes\":1000},\"relayFee\":{\"satoshis\":250,\"bytes\":1000}},{\"id\":2,\"feeType\":\"data\",\"miningFee\":{\"satoshis\":500,\"bytes\":1000},\"relayFee\":{\"satoshis\":250,\"bytes\":1000}}]}",
   	 	"signature": "03045022100eed49f6bf75d8f975f581271e3df658fbe8ec67e6301ea8fc25a72d18c92e30e022056af253f0d24db6a8fde4e2c1ee95e7a5ecf2c7cdc93246f8328c9e0ca582fc40",
    	"publicKey": "03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270","encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	// Invalid sig response
	if strings.Contains(req.URL.String(), "/mapi/tx/"+testTx) {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{
    	"payload": "{\"apiVersion\":\"` + testAPIVersion + `\",\"timestamp\":\"2020-10-11T15:41:29.269Z\",\"returnResult\":\"success\",\"resultDescription\":\"\",\"blockHash\":\"0000000000000000050a09fe90b0e8542bba9e712edb8cc9349e61888fe45ac5\",\"blockHeight\":612530,\"confirmations\":43923,\"minerId\":\"0211ccfc29e3058b770f3cf3eb34b0b2fd2293057a994d4d275121be4151cdf087\",\"txSecondMempoolExpiry\":0}",
   	 	"signature": "03045022100eed49f6bf75d8f975f581271e3df658fbe8ec67e6301ea8fc25a72d18c92e30e022056af253f0d24db6a8fde4e2c1ee95e7a5ecf2c7cdc93246f8328c9e0ca582fc40",
    	"publicKey": "03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270","encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	// Invalid sig response
	if strings.Contains(req.URL.String(), "/mapi/tx") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{
    	"payload": "{\"apiVersion\":\"` + testAPIVersion + `\",\"timestamp\":\"2020-01-15T11:40:29.826Z\",\"txid\":\"6bdbcfab0526d30e8d68279f79dff61fb4026ace8b7b32789af016336e54f2f0\",\"returnResult\":\"success\",\"resultDescription\":\"\",\"minerId\":\"03fcfcfcd0841b0a6ed2057fa8ed404788de47ceb3390c53e79c4ecd1e05819031\",\"currentHighestBlockHash\":\"71a7374389afaec80fcabbbf08dcd82d392cf68c9a13fe29da1a0c853facef01\",\"currentHighestBlockHeight\":207,\"txSecondMempoolExpiry\":0}",
    	"signature": "03045022100f65ae83b20bc60e7a5f0e9c1bd9aceb2b26962ad0ee35472264e83e059f4b9be022010ca2334ff088d6e085eb3c2118306e61ec97781e8e1544e75224533dcc323790",
    	"publicKey": "03fcfcfcd0841b0a6ed2057fa8ed404788de47ceb3390c53e79c4ecd1e05819031","encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	// Default is valid
	return resp, nil
}

// mockHTTPBetterRate for mocking requests
type mockHTTPBetterRate struct{}

const (
	feeQuoteURLMatterPool = "https://merchantapi.matterpool.io/mapi/feeQuote"
	feeQuoteURLMempool    = "https://www.ddpurse.com/openapi/mapi/feeQuote"
	feeQuoteURLTaal       = "https://merchantapi.taal.com/mapi/feeQuote"
)

// Do is a mock http request
func (m *mockHTTPBetterRate) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Valid response
	if req.URL.String() == feeQuoteURLTaal {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{
    	"payload": "{\"apiVersion\":\"` + testAPIVersion + `\",\"timestamp\":\"2020-10-09T21:26:17.410Z\",\"expiryTime\":\"2020-10-09T21:36:17.410Z\",\"minerId\":\"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270\",\"currentHighestBlockHash\":\"0000000000000000035c5f8c0294802a01e500fa7b95337963bb3640da3bd565\",\"currentHighestBlockHeight\":656169,\"minerReputation\":null,\"fees\":[{\"id\":1,\"feeType\":\"standard\",\"miningFee\":{\"satoshis\":475,\"bytes\":1000},\"relayFee\":{\"satoshis\":150,\"bytes\":1000}},{\"id\":2,\"feeType\":\"data\",\"miningFee\":{\"satoshis\":500,\"bytes\":1000},\"relayFee\":{\"satoshis\":250,\"bytes\":1000}}]}",
   	 	"signature": "3045022100eed49f6bf75d8f975f581271e3df658fbe8ec67e6301ea8fc25a72d18c92e30e022056af253f0d24db6a8fde4e2c1ee95e7a5ecf2c7cdc93246f8328c9e0ca582fc4",
    	"publicKey": "03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270","encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	if req.URL.String() == feeQuoteURLMatterPool {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{
    	"payload": "{\"apiVersion\":\"` + testAPIVersion + `\",\"timestamp\":\"2020-10-09T22:08:26.236Z\",\"expiryTime\":\"2020-10-09T22:18:26.236Z\",\"minerId\":\"0211ccfc29e3058b770f3cf3eb34b0b2fd2293057a994d4d275121be4151cdf087\",\"currentHighestBlockHash\":\"0000000000000000028285a9168c95457521a743765f499de389c094e883f42a\",\"currentHighestBlockHeight\":656171,\"minerReputation\":null,\"fees\":[{\"feeType\":\"standard\",\"miningFee\":{\"satoshis\":405,\"bytes\":1000},\"relayFee\":{\"satoshis\":100,\"bytes\":1000}},{\"feeType\":\"data\",\"miningFee\":{\"satoshis\":500,\"bytes\":1000},\"relayFee\":{\"satoshis\":100,\"bytes\":1000}}]}",
    	"signature": "3044022011f90db2661726eb2659c3447ccaa9fd3368194f87d5d86a23e673c45d5d714502200c51eb600e3370b49d759aa4d441000286937b0803037a1d6de4c5a5c559d74c",
    	"publicKey": "0211ccfc29e3058b770f3cf3eb34b0b2fd2293057a994d4d275121be4151cdf087","encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	if req.URL.String() == feeQuoteURLMempool {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{
    	"payload": "{\"apiVersion\":\"` + testAPIVersion + `\",\"timestamp\":\"2020-10-09T22:09:04.433Z\",\"expiryTime\":\"2020-10-09T22:19:04.433Z\",\"minerId\":null,\"currentHighestBlockHash\":\"0000000000000000028285a9168c95457521a743765f499de389c094e883f42a\",\"currentHighestBlockHeight\":656171,\"minerReputation\":null,\"fees\":[{\"feeType\":\"standard\",\"miningFee\":{\"satoshis\":350,\"bytes\":1000},\"relayFee\":{\"satoshis\":250,\"bytes\":1000}},{\"feeType\":\"data\",\"miningFee\":{\"satoshis\":430,\"bytes\":1000},\"relayFee\":{\"satoshis\":175,\"bytes\":1000}}]}",
    	"signature": null,"publicKey": null,"encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	// Default is valid
	return resp, nil
}

// mockHTTPMissingFeeType for mocking requests
type mockHTTPMissingFeeType struct{}

// Do is a mock http request
func (m *mockHTTPMissingFeeType) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Valid response
	if strings.Contains(req.URL.String(), "/mapi/feeQuote") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{
    	"payload": "{\"apiVersion\":\"` + testAPIVersion + `\",\"timestamp\":\"2020-10-09T21:26:17.410Z\",\"expiryTime\":\"2020-10-09T21:36:17.410Z\",\"minerId\":\"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270\",\"currentHighestBlockHash\":\"0000000000000000035c5f8c0294802a01e500fa7b95337963bb3640da3bd565\",\"currentHighestBlockHeight\":656169,\"minerReputation\":null,\"fees\":[{\"id\":2,\"feeType\":\"data\",\"miningFee\":{\"satoshis\":500,\"bytes\":1000},\"relayFee\":{\"satoshis\":250,\"bytes\":1000}}]}",
   	 	"signature": "3045022100eed49f6bf75d8f975f581271e3df658fbe8ec67e6301ea8fc25a72d18c92e30e022056af253f0d24db6a8fde4e2c1ee95e7a5ecf2c7cdc93246f8328c9e0ca582fc4",
    	"publicKey": "03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270","encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	// Default is valid
	return resp, nil
}

// TestClient_FeeQuote tests the method FeeQuote()
func TestClient_FeeQuote(t *testing.T) {

	t.Run("get a valid fee quote", func(t *testing.T) {

		defer goleak.VerifyNone(t)

		// Create a client
		client := newTestClient(&mockHTTPValidFeeQuote{})

		// Create a req
		response, err := client.FeeQuote(context.Background(), client.MinerByName(MinerTaal))
		assert.NoError(t, err)
		assert.NotNil(t, response)

		// Check returned values
		assert.Equal(t, true, response.Validated)
		assert.Equal(t, feeTestSignature, *response.Signature)
		assert.Equal(t, feeTestPublicKey, *response.PublicKey)
		assert.Equal(t, testEncoding, response.Encoding)
		assert.Equal(t, testMimeType, response.MimeType)
	})

	t.Run("valid parse values", func(t *testing.T) {

		defer goleak.VerifyNone(t)

		// Create a client
		client := newTestClient(&mockHTTPValidFeeQuote{})

		// Create a req
		response, err := client.FeeQuote(context.Background(), client.MinerByName(MinerTaal))
		assert.NoError(t, err)
		assert.NotNil(t, response)

		// Test parsed values
		assert.Equal(t, MinerTaal, response.Miner.Name)
		assert.Equal(t, feeTestPublicKey, response.Quote.MinerID)
		assert.Equal(t, "2020-10-09T21:36:17.410Z", response.Quote.ExpirationTime)
		assert.Equal(t, "2020-10-09T21:26:17.410Z", response.Quote.Timestamp)
		assert.Equal(t, testAPIVersion, response.Quote.APIVersion)
		assert.Equal(t, "0000000000000000035c5f8c0294802a01e500fa7b95337963bb3640da3bd565", response.Quote.CurrentHighestBlockHash)
		assert.Equal(t, uint64(656169), response.Quote.CurrentHighestBlockHeight)
		assert.Equal(t, 2, len(response.Quote.Fees))
	})

	t.Run("get actual rates", func(t *testing.T) {

		defer goleak.VerifyNone(t)

		// Create a client
		client := newTestClient(&mockHTTPValidFeeQuote{})

		// Create a req
		response, err := client.FeeQuote(context.Background(), client.MinerByName(MinerTaal))
		assert.NoError(t, err)
		assert.NotNil(t, response)

		// Test getting rate from request
		var rate uint64
		rate, err = response.Quote.CalculateFee(FeeCategoryMining, FeeTypeData, 1000)
		assert.NoError(t, err)
		assert.Equal(t, uint64(500), rate)

		// Test relay rate
		rate, err = response.Quote.CalculateFee(FeeCategoryRelay, FeeTypeData, 1000)
		assert.NoError(t, err)
		assert.Equal(t, uint64(250), rate)
	})

	t.Run("invalid miner", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		client := newTestClient(&mockHTTPValidFeeQuote{})
		response, err := client.FeeQuote(context.Background(), nil)
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("http error", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		client := newTestClient(&mockHTTPError{})
		response, err := client.FeeQuote(context.Background(), client.MinerByName(MinerTaal))
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("bad request", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		client := newTestClient(&mockHTTPBadRequest{})
		response, err := client.FeeQuote(context.Background(), client.MinerByName(MinerTaal))
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		client := newTestClient(&mockHTTPInvalidJSON{})
		response, err := client.FeeQuote(context.Background(), client.MinerByName(MinerTaal))
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("invalid signature", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		client := newTestClient(&mockHTTPInvalidSignature{})
		response, err := client.FeeQuote(context.Background(), client.MinerByName(MinerTaal))
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("missing fees", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		client := newTestClient(&mockHTTPMissingFees{})
		response, err := client.FeeQuote(context.Background(), client.MinerByName(MinerTaal))
		assert.Error(t, err)
		assert.Nil(t, response)
	})
}

// ExampleClient_FeeQuote example using FeeQuote()
func ExampleClient_FeeQuote() {
	// Create a client (using a test client vs NewClient())
	client := newTestClient(&mockHTTPValidFeeQuote{})

	// Create a req
	response, err := client.FeeQuote(context.Background(), client.MinerByName(MinerTaal))
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	fmt.Printf("got quote from: %s", response.Miner.Name)
	// Output:got quote from: Taal
}

// BenchmarkClient_FeeQuote benchmarks the method FeeQuote()
func BenchmarkClient_FeeQuote(b *testing.B) {
	client := newTestClient(&mockHTTPValidFeeQuote{})
	for i := 0; i < b.N; i++ {
		_, _ = client.FeeQuote(context.Background(), client.MinerByName(MinerTaal))
	}
}

// TestFeePayload_CalculateFee tests the method CalculateFee()
func TestFeePayload_CalculateFee(t *testing.T) {
	t.Parallel()

	t.Run("calculate valid fees", func(t *testing.T) {

		// Create a client
		client := newTestClient(&mockHTTPValidFeeQuote{})

		// Create a req
		response, err := client.FeeQuote(context.Background(), client.MinerByName(MinerTaal))
		assert.NoError(t, err)
		assert.NotNil(t, response)

		// Mining & Data
		var fee uint64
		fee, err = response.Quote.CalculateFee(FeeCategoryMining, FeeTypeData, 1000)
		assert.NoError(t, err)
		assert.Equal(t, uint64(500), fee)

		// Mining and standard
		fee, err = response.Quote.CalculateFee(FeeCategoryMining, FeeTypeStandard, 1000)
		assert.NoError(t, err)
		assert.Equal(t, uint64(500), fee)

		// Relay & Data
		fee, err = response.Quote.CalculateFee(FeeCategoryRelay, FeeTypeData, 1000)
		assert.NoError(t, err)
		assert.Equal(t, uint64(250), fee)

		// Relay and standard
		fee, err = response.Quote.CalculateFee(FeeCategoryRelay, FeeTypeStandard, 1000)
		assert.NoError(t, err)
		assert.Equal(t, uint64(250), fee)
	})

	t.Run("calculate zero fee", func(t *testing.T) {
		client := newTestClient(&mockHTTPValidFeeQuote{})
		response, err := client.FeeQuote(context.Background(), client.MinerByName(MinerTaal))
		assert.NoError(t, err)
		assert.NotNil(t, response)

		// Zero tx size produces 0 fee and error
		var fee uint64
		fee, err = response.Quote.CalculateFee(FeeCategoryMining, FeeTypeData, 0)
		assert.Error(t, err)
		assert.Equal(t, uint64(1), fee)
	})

	t.Run("missing fee type", func(t *testing.T) {
		client := newTestClient(&mockHTTPMissingFeeType{})
		response, err := client.FeeQuote(context.Background(), client.MinerByName(MinerTaal))
		assert.NoError(t, err)
		assert.NotNil(t, response)

		// Zero tx size produces 0 fee and error
		var fee uint64
		fee, err = response.Quote.CalculateFee(FeeCategoryRelay, FeeTypeStandard, 1000)
		assert.Error(t, err)
		assert.Equal(t, uint64(1), fee)
	})

}

// ExampleFeePayload_CalculateFee example using CalculateFee()
func ExampleFeePayload_CalculateFee() {
	// Create a client (using a test client vs NewClient())
	client := newTestClient(&mockHTTPValidBestQuote{})

	// Create a req
	response, err := client.BestQuote(context.Background(), FeeCategoryMining, FeeTypeData)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Calculate fee for tx
	var fee uint64
	fee, err = response.Quote.CalculateFee(FeeCategoryMining, FeeTypeData, 1000)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Note: cannot show response since the miner might be different each time
	fmt.Printf("got best quote and fee for 1000 byte tx is: %d", fee)
	// Output:got best quote and fee for 1000 byte tx is: 420
}

// BenchmarkFeePayload_CalculateFee benchmarks the method CalculateFee()
func BenchmarkFeePayload_CalculateFee(b *testing.B) {
	client := newTestClient(&mockHTTPValidBestQuote{})
	response, _ := client.BestQuote(context.Background(), FeeCategoryMining, FeeTypeData)
	for i := 0; i < b.N; i++ {
		_, _ = response.Quote.CalculateFee(FeeCategoryMining, FeeTypeData, 1000)
	}
}

// TestFeePayload_GetFee tests the method GetFee()
func TestFeePayload_GetFee(t *testing.T) {
	t.Parallel()

	t.Run("get valid fees", func(t *testing.T) {

		// Create a client
		client := newTestClient(&mockHTTPValidFeeQuote{})

		// Create a req
		response, err := client.FeeQuote(context.Background(), client.MinerByName(MinerTaal))
		assert.NoError(t, err)
		assert.NotNil(t, response)

		// Standard
		fee := response.Quote.GetFee(FeeTypeStandard)
		assert.NotNil(t, fee)
		assert.Equal(t, "standard", fee.FeeType)

		// Data
		fee = response.Quote.GetFee(FeeTypeData)
		assert.NotNil(t, fee)
		assert.Equal(t, "data", fee.FeeType)
	})

	t.Run("missing fee type", func(t *testing.T) {
		client := newTestClient(&mockHTTPMissingFeeType{})
		response, err := client.FeeQuote(context.Background(), client.MinerByName(MinerTaal))
		assert.NoError(t, err)
		assert.NotNil(t, response)

		// Standard
		fee := response.Quote.GetFee("")
		assert.Nil(t, fee)
	})
}

// ExampleFeePayload_GetFee example using GetFee()
func ExampleFeePayload_GetFee() {
	// Create a client (using a test client vs NewClient())
	client := newTestClient(&mockHTTPValidBestQuote{})

	// Create a req
	response, err := client.BestQuote(context.Background(), FeeCategoryMining, FeeTypeData)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Get the fee
	fee := response.Quote.GetFee(FeeTypeStandard)

	fmt.Printf(
		"got best quote and fee for %d byte tx is %d sats",
		fee.MiningFee.Bytes, fee.MiningFee.Satoshis,
	)
	// Output:got best quote and fee for 1000 byte tx is 500 sats
}

// BenchmarkFeePayload_GetFee benchmarks the method GetFee()
func BenchmarkFeePayload_GetFee(b *testing.B) {
	client := newTestClient(&mockHTTPValidBestQuote{})
	response, _ := client.BestQuote(context.Background(), FeeCategoryMining, FeeTypeData)
	for i := 0; i < b.N; i++ {
		_ = response.Quote.GetFee(FeeTypeStandard)
	}
}
