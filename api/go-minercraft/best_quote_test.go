package minercraft

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

// mockHTTPValidBestQuote for mocking requests
type mockHTTPValidBestQuote struct{}

// Do is a mock http request
func (m *mockHTTPValidBestQuote) Do(req *http.Request) (*http.Response, error) {
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
    	"payload": "{\"apiVersion\":\"` + testAPIVersion + `\",\"timestamp\":\"2020-10-09T21:26:17.410Z\",\"expiryTime\":\"2020-10-09T21:36:17.410Z\",\"minerId\":\"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270\",\"currentHighestBlockHash\":\"0000000000000000035c5f8c0294802a01e500fa7b95337963bb3640da3bd565\",\"currentHighestBlockHeight\":656169,\"minerReputation\":null,\"fees\":[{\"id\":1,\"feeType\":\"standard\",\"miningFee\":{\"satoshis\":400,\"bytes\":1000},\"relayFee\":{\"satoshis\":250,\"bytes\":1000}},{\"id\":2,\"feeType\":\"data\",\"miningFee\":{\"satoshis\":500,\"bytes\":1000},\"relayFee\":{\"satoshis\":225,\"bytes\":1000}}]}",
   	 	"signature": "3045022100eed49f6bf75d8f975f581271e3df658fbe8ec67e6301ea8fc25a72d18c92e30e022056af253f0d24db6a8fde4e2c1ee95e7a5ecf2c7cdc93246f8328c9e0ca582fc4",
    	"publicKey": "03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270","encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	if req.URL.String() == feeQuoteURLMatterPool {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{
    	"payload": "{\"apiVersion\":\"` + testAPIVersion + `\",\"timestamp\":\"2020-10-09T22:08:26.236Z\",\"expiryTime\":\"2020-10-09T22:18:26.236Z\",\"minerId\":\"0211ccfc29e3058b770f3cf3eb34b0b2fd2293057a994d4d275121be4151cdf087\",\"currentHighestBlockHash\":\"0000000000000000028285a9168c95457521a743765f499de389c094e883f42a\",\"currentHighestBlockHeight\":656171,\"minerReputation\":null,\"fees\":[{\"feeType\":\"standard\",\"miningFee\":{\"satoshis\":400,\"bytes\":1000},\"relayFee\":{\"satoshis\":100,\"bytes\":1000}},{\"feeType\":\"data\",\"miningFee\":{\"satoshis\":430,\"bytes\":1000},\"relayFee\":{\"satoshis\":110,\"bytes\":1000}}]}",
    	"signature": "3044022011f90db2661726eb2659c3447ccaa9fd3368194f87d5d86a23e673c45d5d714502200c51eb600e3370b49d759aa4d441000286937b0803037a1d6de4c5a5c559d74c",
    	"publicKey": "0211ccfc29e3058b770f3cf3eb34b0b2fd2293057a994d4d275121be4151cdf087","encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	if req.URL.String() == feeQuoteURLMempool {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{
    	"payload": "{\"apiVersion\":\"` + testAPIVersion + `\",\"timestamp\":\"2020-10-09T22:09:04.433Z\",\"expiryTime\":\"2020-10-09T22:19:04.433Z\",\"minerId\":null,\"currentHighestBlockHash\":\"0000000000000000028285a9168c95457521a743765f499de389c094e883f42a\",\"currentHighestBlockHeight\":656171,\"minerReputation\":null,\"fees\":[{\"feeType\":\"standard\",\"miningFee\":{\"satoshis\":500,\"bytes\":1000},\"relayFee\":{\"satoshis\":250,\"bytes\":1000}},{\"feeType\":\"data\",\"miningFee\":{\"satoshis\":420,\"bytes\":1000},\"relayFee\":{\"satoshis\":150,\"bytes\":1000}}]}",
    	"signature": null,"publicKey": null,"encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	// Default is valid
	return resp, nil
}

// mockHTTPBadRate for mocking requests
type mockHTTPBadRate struct{}

// Do is a mock http request
func (m *mockHTTPBadRate) Do(req *http.Request) (*http.Response, error) {
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
    	"payload": "{\"apiVersion\":\"` + testAPIVersion + `\",\"timestamp\":\"2020-10-09T21:26:17.410Z\",\"expiryTime\":\"2020-10-09T21:36:17.410Z\",\"minerId\":\"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270\",\"currentHighestBlockHash\":\"0000000000000000035c5f8c0294802a01e500fa7b95337963bb3640da3bd565\",\"currentHighestBlockHeight\":656169,\"minerReputation\":null,\"fees\":[{\"id\":1,\"feeType\":\"standard\",\"miningFee\":{\"satoshis\":0,\"bytes\":1000},\"relayFee\":{\"satoshis\":0,\"bytes\":1000}},{\"id\":2,\"feeType\":\"data\",\"miningFee\":{\"satoshis\":0,\"bytes\":1000},\"relayFee\":{\"satoshis\":0,\"bytes\":1000}}]}",
   	 	"signature": "3045022100eed49f6bf75d8f975f581271e3df658fbe8ec67e6301ea8fc25a72d18c92e30e022056af253f0d24db6a8fde4e2c1ee95e7a5ecf2c7cdc93246f8328c9e0ca582fc4",
    	"publicKey": "03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270","encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	if req.URL.String() == feeQuoteURLMatterPool {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{
    	"payload": "{\"apiVersion\":\"` + testAPIVersion + `\",\"timestamp\":\"2020-10-09T22:08:26.236Z\",\"expiryTime\":\"2020-10-09T22:18:26.236Z\",\"minerId\":\"0211ccfc29e3058b770f3cf3eb34b0b2fd2293057a994d4d275121be4151cdf087\",\"currentHighestBlockHash\":\"0000000000000000028285a9168c95457521a743765f499de389c094e883f42a\",\"currentHighestBlockHeight\":656171,\"minerReputation\":null,\"fees\":[{\"feeType\":\"standard\",\"miningFee\":{\"satoshis\":0,\"bytes\":1000},\"relayFee\":{\"satoshis\":0,\"bytes\":1000}},{\"feeType\":\"data\",\"miningFee\":{\"satoshis\":0,\"bytes\":1000},\"relayFee\":{\"satoshis\":0,\"bytes\":1000}}]}",
    	"signature": "3044022011f90db2661726eb2659c3447ccaa9fd3368194f87d5d86a23e673c45d5d714502200c51eb600e3370b49d759aa4d441000286937b0803037a1d6de4c5a5c559d74c",
    	"publicKey": "0211ccfc29e3058b770f3cf3eb34b0b2fd2293057a994d4d275121be4151cdf087","encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	if req.URL.String() == feeQuoteURLMempool {
		resp.StatusCode = http.StatusBadRequest
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
	}

	// Default is valid
	return resp, nil
}

// mockHTTPBestQuoteTwoFailed for mocking requests
type mockHTTPBestQuoteTwoFailed struct{}

// Do is a mock http request
func (m *mockHTTPBestQuoteTwoFailed) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Valid response
	if req.URL.String() == feeQuoteURLTaal {
		resp.StatusCode = http.StatusBadRequest
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
	}

	if req.URL.String() == feeQuoteURLMatterPool {
		resp.StatusCode = http.StatusBadRequest
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
	}

	if req.URL.String() == feeQuoteURLMempool {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{
    	"payload": "{\"apiVersion\":\"` + testAPIVersion + `\",\"timestamp\":\"2020-10-09T22:09:04.433Z\",\"expiryTime\":\"2020-10-09T22:19:04.433Z\",\"minerId\":null,\"currentHighestBlockHash\":\"0000000000000000028285a9168c95457521a743765f499de389c094e883f42a\",\"currentHighestBlockHeight\":656171,\"minerReputation\":null,\"fees\":[{\"feeType\":\"standard\",\"miningFee\":{\"satoshis\":500,\"bytes\":1000},\"relayFee\":{\"satoshis\":250,\"bytes\":1000}},{\"feeType\":\"data\",\"miningFee\":{\"satoshis\":420,\"bytes\":1000},\"relayFee\":{\"satoshis\":150,\"bytes\":1000}}]}",
    	"signature": null,"publicKey": null,"encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	// Default is valid
	return resp, nil
}

// mockHTTPBestQuoteAllFailed for mocking requests
type mockHTTPBestQuoteAllFailed struct{}

// Do is a mock http request
func (m *mockHTTPBestQuoteAllFailed) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Valid response
	if req.URL.String() == feeQuoteURLTaal {
		resp.StatusCode = http.StatusBadRequest
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
	}

	if req.URL.String() == feeQuoteURLMatterPool {
		resp.StatusCode = http.StatusBadRequest
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
	}

	if req.URL.String() == feeQuoteURLMempool {
		resp.StatusCode = http.StatusBadRequest
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
	}

	// Default is valid
	return resp, nil
}

// TestClient_BestQuote tests the method BestQuote()
func TestClient_BestQuote(t *testing.T) {

	t.Run("get a valid best quote", func(t *testing.T) {

		defer goleak.VerifyNone(t)

		// Create a client
		client := newTestClient(&mockHTTPValidBestQuote{})

		// Create a req
		response, err := client.BestQuote(context.Background(), FeeCategoryMining, FeeTypeData)
		assert.NoError(t, err)
		assert.NotNil(t, response)

		// Check returned values
		assert.Equal(t, testEncoding, response.Encoding)
		assert.Equal(t, testMimeType, response.MimeType)

		// Check that we got fees
		assert.Equal(t, 2, len(response.Quote.Fees))
	})

	t.Run("http error", func(t *testing.T) {
		defer goleak.VerifyNone(t)

		client := newTestClient(&mockHTTPError{})
		response, err := client.BestQuote(context.Background(), FeeCategoryMining, FeeTypeData)
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("bad request", func(t *testing.T) {
		defer goleak.VerifyNone(t)

		client := newTestClient(&mockHTTPBadRequest{})
		response, err := client.BestQuote(context.Background(), FeeCategoryMining, FeeTypeData)
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		defer goleak.VerifyNone(t)

		client := newTestClient(&mockHTTPInvalidJSON{})
		response, err := client.BestQuote(context.Background(), FeeCategoryMining, FeeTypeData)
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("invalid category", func(t *testing.T) {
		defer goleak.VerifyNone(t)

		client := newTestClient(&mockHTTPValidBestQuote{})
		response, err := client.BestQuote(context.Background(), "invalid", FeeTypeData)
		assert.Error(t, err)
		assert.Nil(t, response)

		// Create a req
		response, err = client.BestQuote(context.Background(), FeeCategoryMining, "invalid")
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("better rate", func(t *testing.T) {
		defer goleak.VerifyNone(t)

		client := newTestClient(&mockHTTPBetterRate{})
		response, err := client.BestQuote(context.Background(), FeeCategoryRelay, FeeTypeData)
		assert.NoError(t, err)
		assert.NotNil(t, response)

		// Check that we got fees
		assert.Equal(t, 2, len(response.Quote.Fees))

		var fee uint64
		fee, err = response.Quote.CalculateFee(FeeCategoryRelay, FeeTypeData, 1000)
		assert.NoError(t, err)
		assert.Equal(t, uint64(100), fee)

		fee, err = response.Quote.CalculateFee(FeeCategoryMining, FeeTypeData, 1000)
		assert.NoError(t, err)
		assert.Equal(t, uint64(500), fee)
	})

	t.Run("bad rate", func(t *testing.T) {
		defer goleak.VerifyNone(t)

		client := newTestClient(&mockHTTPBadRate{})
		response, err := client.BestQuote(context.Background(), FeeCategoryRelay, FeeTypeData)
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("best quote - two failed", func(t *testing.T) {
		defer goleak.VerifyNone(t)

		// Create a client
		client := newTestClient(&mockHTTPBestQuoteTwoFailed{})

		// Create a req
		response, err := client.BestQuote(context.Background(), FeeCategoryMining, FeeTypeData)
		assert.NoError(t, err)
		assert.NotNil(t, response)

		// Check returned values
		assert.Equal(t, testEncoding, response.Encoding)
		assert.Equal(t, testMimeType, response.MimeType)

		// Check that we got fees
		assert.Equal(t, 2, len(response.Quote.Fees))
		assert.Equal(t, MinerMempool, response.Miner.Name)
	})

	t.Run("best quote - all failed", func(t *testing.T) {

		defer goleak.VerifyNone(t)

		// Create a client
		client := newTestClient(&mockHTTPBestQuoteAllFailed{})

		// Create a req
		response, err := client.BestQuote(context.Background(), FeeCategoryMining, FeeTypeData)
		assert.Error(t, err)
		assert.Nil(t, response)
	})
}

// ExampleClient_BestQuote example using BestQuote()
func ExampleClient_BestQuote() {
	// Create a client (using a test client vs NewClient())
	client := newTestClient(&mockHTTPValidBestQuote{})

	// Create a req
	_, err := client.BestQuote(context.Background(), FeeCategoryMining, FeeTypeData)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Note: cannot show response since the miner might be different each time
	fmt.Printf("got best quote!")
	// Output:got best quote!
}

// BenchmarkClient_BestQuote benchmarks the method BestQuote()
func BenchmarkClient_BestQuote(b *testing.B) {
	client := newTestClient(&mockHTTPValidBestQuote{})
	for i := 0; i < b.N; i++ {
		_, _ = client.BestQuote(context.Background(), FeeCategoryMining, FeeTypeData)
	}
}
