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

const submitTestSignature = "3045022100f65ae83b20bc60e7a5f0e9c1bd9aceb2b26962ad0ee35472264e83e059f4b9be022010ca2334ff088d6e085eb3c2118306e61ec97781e8e1544e75224533dcc32379"
const submitTestPublicKey = "03fcfcfcd0841b0a6ed2057fa8ed404788de47ceb3390c53e79c4ecd1e05819031"
const submitTestExampleTx = "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff1c03d7c6082f7376706f6f6c2e636f6d2f3edff034600055b8467f0040ffffffff01247e814a000000001976a914492558fb8ca71a3591316d095afc0f20ef7d42f788ac00000000"

// mockHTTPValidSubmission for mocking requests
type mockHTTPValidSubmission struct{}

// Do is a mock http request
func (m *mockHTTPValidSubmission) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Valid response
	if strings.Contains(req.URL.String(), "/mapi/tx") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{
    	"payload": "{\"apiVersion\":\"` + testAPIVersion + `\",\"timestamp\":\"2020-01-15T11:40:29.826Z\",\"txid\":\"6bdbcfab0526d30e8d68279f79dff61fb4026ace8b7b32789af016336e54f2f0\",\"returnResult\":\"success\",\"resultDescription\":\"\",\"minerId\":\"03fcfcfcd0841b0a6ed2057fa8ed404788de47ceb3390c53e79c4ecd1e05819031\",\"currentHighestBlockHash\":\"71a7374389afaec80fcabbbf08dcd82d392cf68c9a13fe29da1a0c853facef01\",\"currentHighestBlockHeight\":207,\"txSecondMempoolExpiry\":0}",
    	"signature": "3045022100f65ae83b20bc60e7a5f0e9c1bd9aceb2b26962ad0ee35472264e83e059f4b9be022010ca2334ff088d6e085eb3c2118306e61ec97781e8e1544e75224533dcc32379",
    	"publicKey": "03fcfcfcd0841b0a6ed2057fa8ed404788de47ceb3390c53e79c4ecd1e05819031","encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	// Default is valid
	return resp, nil
}

// mockHTTPBadSubmission for mocking requests
type mockHTTPBadSubmission struct{}

// Do is a mock http request
func (m *mockHTTPBadSubmission) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Valid response
	if strings.Contains(req.URL.String(), "/mapi/tx") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{
    	"payload": "{}",
    	"signature": "3044022066a8a39ff5f5eae818636aa03fdfc386ea4f33f41993cf41d4fb6d4745ae032102206a8895a6f742d809647ad1a1df12230e9b480275853ed28bc178f4b48afd802a",
    	"publicKey": "0211ccfc29e3058b770f3cf3eb34b0b2fd2293057a994d4d275121be4151cdf087","encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	// Default is valid
	return resp, nil
}

// TestClient_SubmitTransaction tests the method SubmitTransaction()
func TestClient_SubmitTransaction(t *testing.T) {

	tx := &Transaction{RawTx: submitTestExampleTx}

	t.Run("submit a valid transaction", func(t *testing.T) {

		defer goleak.VerifyNone(t)

		// Create a client
		client := newTestClient(&mockHTTPValidSubmission{})

		// Create a req
		response, err := client.SubmitTransaction(context.Background(), client.MinerByName(MinerMatterpool), tx)
		assert.NoError(t, err)
		assert.NotNil(t, response)

		// Check returned values
		assert.Equal(t, true, response.Validated)
		assert.Equal(t, submitTestSignature, *response.Signature)
		assert.Equal(t, submitTestPublicKey, *response.PublicKey)
		assert.Equal(t, testEncoding, response.Encoding)
		assert.Equal(t, testMimeType, response.MimeType)
	})

	t.Run("validate parsed values", func(t *testing.T) {

		defer goleak.VerifyNone(t)

		client := newTestClient(&mockHTTPValidSubmission{})
		response, err := client.SubmitTransaction(context.Background(), client.MinerByName(MinerMatterpool), tx)
		assert.NoError(t, err)
		assert.NotNil(t, response)

		// Test parsed values
		assert.Equal(t, MinerMatterpool, response.Miner.Name)
		assert.Equal(t, submitTestPublicKey, response.Results.MinerID)
		assert.Equal(t, "2020-01-15T11:40:29.826Z", response.Results.Timestamp)
		assert.Equal(t, testAPIVersion, response.Results.APIVersion)
		assert.Equal(t, QueryTransactionSuccess, response.Results.ReturnResult)
		assert.Equal(t, "6bdbcfab0526d30e8d68279f79dff61fb4026ace8b7b32789af016336e54f2f0", response.Results.TxID)
	})

	t.Run("invalid miner", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		client := newTestClient(&mockHTTPValidSubmission{})
		response, err := client.SubmitTransaction(context.Background(), nil, tx)
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("http error", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		client := newTestClient(&mockHTTPError{})
		response, err := client.SubmitTransaction(context.Background(), client.MinerByName(MinerMatterpool), tx)
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("bad request", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		client := newTestClient(&mockHTTPBadRequest{})
		response, err := client.SubmitTransaction(context.Background(), client.MinerByName(MinerMatterpool), tx)
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		client := newTestClient(&mockHTTPInvalidJSON{})
		response, err := client.SubmitTransaction(context.Background(), client.MinerByName(MinerMatterpool), tx)
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("invalid signature", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		client := newTestClient(&mockHTTPInvalidSignature{})
		response, err := client.SubmitTransaction(context.Background(), client.MinerByName(MinerMatterpool), tx)
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("bad submission", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		client := newTestClient(&mockHTTPBadSubmission{})
		response, err := client.SubmitTransaction(context.Background(), client.MinerByName(MinerMatterpool), tx)
		assert.Error(t, err)
		assert.Nil(t, response)
	})
}

// ExampleClient_SubmitTransaction example using SubmitTransaction()
func ExampleClient_SubmitTransaction() {
	// Create a client (using a test client vs NewClient())
	client := newTestClient(&mockHTTPValidSubmission{})

	tx := &Transaction{RawTx: submitTestExampleTx}

	// Create a req
	response, err := client.SubmitTransaction(context.Background(), client.MinerByName(MinerTaal), tx)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	fmt.Printf("submitted tx to: %s", response.Miner.Name)
	// Output:submitted tx to: Taal
}

// BenchmarkClient_SubmitTransaction benchmarks the method SubmitTransaction()
func BenchmarkClient_SubmitTransaction(b *testing.B) {
	client := newTestClient(&mockHTTPValidSubmission{})
	miner := client.MinerByName(MinerTaal)
	tx := &Transaction{RawTx: submitTestExampleTx}
	for i := 0; i < b.N; i++ {
		_, _ = client.SubmitTransaction(context.Background(), miner, tx)
	}
}
