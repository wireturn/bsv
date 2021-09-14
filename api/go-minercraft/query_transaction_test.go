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
)

const queryTestSignature = "3044022066a8a39ff5f5eae818636aa03fdfc386ea4f33f41993cf41d4fb6d4745ae032102206a8895a6f742d809647ad1a1df12230e9b480275853ed28bc178f4b48afd802a"
const queryTestPublicKey = "0211ccfc29e3058b770f3cf3eb34b0b2fd2293057a994d4d275121be4151cdf087"

// mockHTTPValidQuery for mocking requests
type mockHTTPValidQuery struct{}

// Do is a mock http request
func (m *mockHTTPValidQuery) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Valid response
	if strings.Contains(req.URL.String(), "/mapi/tx/"+testTx) {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{
    	"payload": "{\"apiVersion\":\"` + testAPIVersion + `\",\"timestamp\":\"2020-10-10T13:07:26.014Z\",\"returnResult\":\"success\",\"resultDescription\":\"\",\"blockHash\":\"0000000000000000050a09fe90b0e8542bba9e712edb8cc9349e61888fe45ac5\",\"blockHeight\":612530,\"confirmations\":43733,\"minerId\":\"0211ccfc29e3058b770f3cf3eb34b0b2fd2293057a994d4d275121be4151cdf087\",\"txSecondMempoolExpiry\":0}",
    	"signature": "3044022066a8a39ff5f5eae818636aa03fdfc386ea4f33f41993cf41d4fb6d4745ae032102206a8895a6f742d809647ad1a1df12230e9b480275853ed28bc178f4b48afd802a",
    	"publicKey": "0211ccfc29e3058b770f3cf3eb34b0b2fd2293057a994d4d275121be4151cdf087","encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	// Default is valid
	return resp, nil
}

// mockHTTPBadQuery for mocking requests
type mockHTTPBadQuery struct{}

// Do is a mock http request
func (m *mockHTTPBadQuery) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Valid response
	if strings.Contains(req.URL.String(), "/mapi/tx/"+testTx) {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{
    	"payload": "{}",
    	"signature": "3044022066a8a39ff5f5eae818636aa03fdfc386ea4f33f41993cf41d4fb6d4745ae032102206a8895a6f742d809647ad1a1df12230e9b480275853ed28bc178f4b48afd802a",
    	"publicKey": "0211ccfc29e3058b770f3cf3eb34b0b2fd2293057a994d4d275121be4151cdf087","encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	// Default is valid
	return resp, nil
}

// TestClient_QueryTransaction tests the method QueryTransaction()
func TestClient_QueryTransaction(t *testing.T) {
	t.Parallel()

	t.Run("query a valid transaction", func(t *testing.T) {
		// Create a client
		client := newTestClient(&mockHTTPValidQuery{})

		// Create a req
		response, err := client.QueryTransaction(context.Background(), client.MinerByName(MinerMatterpool), testTx)
		assert.NoError(t, err)
		assert.NotNil(t, response)

		// Check returned values
		assert.Equal(t, true, response.Validated)
		assert.Equal(t, queryTestSignature, *response.Signature)
		assert.Equal(t, queryTestPublicKey, *response.PublicKey)
		assert.Equal(t, testEncoding, response.Encoding)
		assert.Equal(t, testMimeType, response.MimeType)
	})

	t.Run("validate parsed values", func(t *testing.T) {
		client := newTestClient(&mockHTTPValidQuery{})
		response, err := client.QueryTransaction(context.Background(), client.MinerByName(MinerMatterpool), testTx)
		assert.NoError(t, err)
		assert.NotNil(t, response)

		// Test parsed values
		assert.Equal(t, MinerMatterpool, response.Miner.Name)
		assert.Equal(t, queryTestPublicKey, response.Miner.MinerID)
		assert.Equal(t, "2020-10-10T13:07:26.014Z", response.Query.Timestamp)
		assert.Equal(t, testAPIVersion, response.Query.APIVersion)
		assert.Equal(t, "0000000000000000050a09fe90b0e8542bba9e712edb8cc9349e61888fe45ac5", response.Query.BlockHash)
		assert.Equal(t, int64(612530), response.Query.BlockHeight)
		assert.Equal(t, QueryTransactionSuccess, response.Query.ReturnResult)
	})

	t.Run("invalid miner", func(t *testing.T) {
		client := newTestClient(&mockHTTPValidFeeQuote{})
		response, err := client.QueryTransaction(context.Background(), nil, testTx)
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("http error", func(t *testing.T) {
		client := newTestClient(&mockHTTPError{})
		response, err := client.QueryTransaction(context.Background(), client.MinerByName(MinerMatterpool), testTx)
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("bad request", func(t *testing.T) {
		client := newTestClient(&mockHTTPBadRequest{})
		response, err := client.QueryTransaction(context.Background(), client.MinerByName(MinerMatterpool), testTx)
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		client := newTestClient(&mockHTTPInvalidJSON{})
		response, err := client.QueryTransaction(context.Background(), client.MinerByName(MinerMatterpool), testTx)
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("invalid signature", func(t *testing.T) {
		client := newTestClient(&mockHTTPInvalidSignature{})
		response, err := client.QueryTransaction(context.Background(), client.MinerByName(MinerMatterpool), testTx)
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("bad query", func(t *testing.T) {
		client := newTestClient(&mockHTTPBadQuery{})
		response, err := client.QueryTransaction(context.Background(), client.MinerByName(MinerMatterpool), testTx)
		assert.Error(t, err)
		assert.Nil(t, response)
	})
}

// ExampleClient_QueryTransaction example using QueryTransaction()
func ExampleClient_QueryTransaction() {
	// Create a client (using a test client vs NewClient())
	client := newTestClient(&mockHTTPValidQuery{})

	// Create a req
	response, err := client.QueryTransaction(context.Background(), client.MinerByName(MinerTaal), testTx)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	fmt.Printf("got tx status %s from: %s", response.Query.ReturnResult, response.Miner.Name)
	// Output:got tx status success from: Taal
}

// BenchmarkClient_QueryTransaction benchmarks the method QueryTransaction()
func BenchmarkClient_QueryTransaction(b *testing.B) {
	client := newTestClient(&mockHTTPValidQuery{})
	miner := client.MinerByName(MinerTaal)
	for i := 0; i < b.N; i++ {
		_, _ = client.QueryTransaction(context.Background(), miner, testTx)
	}
}
