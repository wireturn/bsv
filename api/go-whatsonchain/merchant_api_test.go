package whatsonchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

// txMerchantBroadcast is the struct for broadcasting a tx
type txMerchantBroadcast struct {
	RawTx string `json:"rawtx"`
}

// Testing miner id for our merchant tests
var testMinerID = "ba001df8"

// mockHTTPMerchantValid for mocking requests
type mockHTTPMerchantValid struct{}

// Do is a mock http request
func (m *mockHTTPMerchantValid) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Valid (quotes)
	if strings.Contains(req.URL.String(), "/mapi/feeQuotes") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"quotes":[{"providerName":"taal","providerId":"` + testMinerID + `","quote":{"apiVersion":"0.1.0","timestamp":"2020-06-23T18:10:46.571Z","expiryTime":"2020-06-23T18:20:46.571Z","minerId":"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270","currentHighestBlockHash":"000000000000000003051e36242b20fb96cb1338c4247480a3a652b90892c350","currentHighestBlockHeight":640651,"minerReputation":null,"fees":[{"feeType":"standard","miningFee":{"satoshis":5,"bytes":10},"relayFee":{"satoshis":25,"bytes":100}},{"feeType":"data","miningFee":{"satoshis":5,"bytes":10},"relayFee":{"satoshis":25,"bytes":100}}]},"payload":"{\"apiVersion\":\"0.1.0\",\"timestamp\":\"2020-06-23T18:10:46.571Z\",\"expiryTime\":\"2020-06-23T18:20:46.571Z\",\"minerId\":\"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270\",\"currentHighestBlockHash\":\"000000000000000003051e36242b20fb96cb1338c4247480a3a652b90892c350\",\"currentHighestBlockHeight\":640651,\"minerReputation\":null,\"fees\":[{\"feeType\":\"standard\",\"miningFee\":{\"satoshis\":5,\"bytes\":10},\"relayFee\":{\"satoshis\":25,\"bytes\":100}},{\"feeType\":\"data\",\"miningFee\":{\"satoshis\":5,\"bytes\":10},\"relayFee\":{\"satoshis\":25,\"bytes\":100}}]}","signature":"3045022100e391c86f3101458aad62312c238caf21198d5c85ba641d32d6e9bbaa9a0f2fac02204e94ccd4ce2319192cf4b9d68097ec05d2be0a30b371e2132204bab8b68c4608","publicKey":"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270","txSubmissionUrl":"/mapi/` + testMinerID + `/tx","txStatusUrl":"/mapi/` + testMinerID + `/tx/{hash:[0-9a-fA-F]+}"},{"providerName":"mempool","providerId":"ab398390","quote":{"apiVersion":"0.1.0","timestamp":"2020-06-23T18:10:48.319Z","expiryTime":"2020-06-23T18:20:48.319Z","minerId":null,"currentHighestBlockHash":"000000000000000003051e36242b20fb96cb1338c4247480a3a652b90892c350","currentHighestBlockHeight":640651,"minerReputation":null,"fees":[{"feeType":"standard","miningFee":{"satoshis":500,"bytes":1000},"relayFee":{"satoshis":250,"bytes":1000}},{"feeType":"data","miningFee":{"satoshis":500,"bytes":1000},"relayFee":{"satoshis":250,"bytes":1000}}]},"payload":"{\"apiVersion\":\"0.1.0\",\"timestamp\":\"2020-06-23T18:10:48.319Z\",\"expiryTime\":\"2020-06-23T18:20:48.319Z\",\"minerId\":null,\"currentHighestBlockHash\":\"000000000000000003051e36242b20fb96cb1338c4247480a3a652b90892c350\",\"currentHighestBlockHeight\":640651,\"minerReputation\":null,\"fees\":[{\"feeType\":\"standard\",\"miningFee\":{\"satoshis\":500,\"bytes\":1000},\"relayFee\":{\"satoshis\":250,\"bytes\":1000}},{\"feeType\":\"data\",\"miningFee\":{\"satoshis\":500,\"bytes\":1000},\"relayFee\":{\"satoshis\":250,\"bytes\":1000}}]}","signature":null,"publicKey":null,"txSubmissionUrl":"/mapi/ab398390/tx","txStatusUrl":"/mapi/ab398390/tx/{hash:[0-9a-fA-F]+}"}]}`)))
	}

	// Valid (submit)
	if strings.Contains(req.URL.String(), "/mapi/"+testMinerID+"/tx") && req.Method == http.MethodPost {

		decoder := json.NewDecoder(req.Body)
		var data txMerchantBroadcast
		err := decoder.Decode(&data)
		if err != nil {
			return resp, err
		}

		// Valid - but getting error from merchant
		if strings.Contains(data.RawTx, "020000000232900e9b0e359cb95ad3853e1450591fdf01e3efaa1b3b2b5ab5c5ef784946b8010000006a473044022055ec4f9b9cbdd97cf5f4893f921da75287483edb4ebba4f5b23231577212fd5f022007ed037ab7da039e0d4cf0fa35f620f2d7f71285959dcb6c885652d843d0038741210232b357c5309644cf4aa72b9b2d8bfe58bdf2515d40119318d5cb51ef378cae7efffffffff91d2feda79506806c5b0dc74b6fa6ae42fb7963da460ac256383fce498b9952020000006b483045022100ac75defcda55d644b6c095c2e7cbded92e38ea52e4d7f561f37962aa036a92fe0220059c0071d7c53cc964f641fa60ba2a076e3030ae6edc0f650792041c7691c66641210282d7e568e56f59e01a4edae297ac26caabc4684971ac6c7558c91c0fa84002f7ffffffff03322a093f000000001976a914c4263eb96d88849f498d139424b59a0cba1005e888ac2f92bb09000000001976a9146cbff9881ac47da8cb699e4543c28f9b3d6941da88ac404b4c00000000001976a914f7899faf1696892e6cb029b00c713f044761f03588ac00000000") {
			resp.StatusCode = http.StatusOK
			resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"providerName":"taal","providerId":"` + testMinerID + `","response":{"apiVersion":"0.1.0","timestamp":"2020-06-23T19:27:49.626Z","txid":"","returnResult":"failure","resultDescription":"ERROR: Missing inputs","minerId":"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270","currentHighestBlockHash":"0000000000000000014ea02d39fe38f5477ff330f5874778445bad57f7d380cc","currentHighestBlockHeight":640657,"txSecondMempoolExpiry":0},"error":{"status":0,"code":0,"error":""},"payload":"{\"apiVersion\":\"0.1.0\",\"timestamp\":\"2020-06-23T19:27:49.626Z\",\"txid\":\"\",\"returnResult\":\"failure\",\"resultDescription\":\"ERROR: Missing inputs\",\"minerId\":\"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270\",\"currentHighestBlockHash\":\"0000000000000000014ea02d39fe38f5477ff330f5874778445bad57f7d380cc\",\"currentHighestBlockHeight\":640657,\"txSecondMempoolExpiry\":0}","signature":"30440220745e4b3a794931a97b7735ceffd0a082d65edf768a2c7ba6086cb386189044bc0220052da10a0968b5d0f77ef39c96419899e10915a67f78d8732833633beed1b10a","publicKey":"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270"}`)))
		}

		// Valid - faked the response
		if strings.Contains(data.RawTx, "0100000001c871b60f86ad5d292ba51f9e94241fa2316aa4180dbdbae8286eec0bec679f01010000006a473044022100fafb74a5c6760ff2135cdf9ceee8e0342a24a6cc40804bd3245fe4aa055c8611021f6a07c86caa00d4bf92437d40a4b516aa6b7e870ff8e2a07782a5968a72ed99412103ce0ba24104223d4b1258833bacc147b27f6bd4361f403ec5cb1ee5efd0858759ffffffff0260dc2c00000000001976a914493e769b46c077954a376f9a0933f5bb9376482688ac4b807100000000001976a9141af4f2ec0ced0dece5c689e4afb28318e24b1f2d88ac00000000") {
			resp.StatusCode = http.StatusOK
			resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"providerName":"taal","providerId":"` + testMinerID + `","response":{"apiVersion":"0.1.0","timestamp":"2020-06-23T19:28:48.907Z","txid":"bbce41682bfb88d037de116996781c61235a15e6f530597075d5d69bb5a3789a","returnResult":"success","resultDescription":"","minerId":"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270","currentHighestBlockHash":"0000000000000000014ea02d39fe38f5477ff330f5874778445bad57f7d380cc","currentHighestBlockHeight":640657,"txSecondMempoolExpiry":0},"error":{"status":0,"code":0,"error":""},"payload":"{\"apiVersion\":\"0.1.0\",\"timestamp\":\"2020-06-23T19:28:48.907Z\",\"txid\":\"\",\"returnResult\":\"failure\",\"resultDescription\":\"ERROR: 257: txn-already-known\",\"minerId\":\"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270\",\"currentHighestBlockHash\":\"0000000000000000014ea02d39fe38f5477ff330f5874778445bad57f7d380cc\",\"currentHighestBlockHeight\":640657,\"txSecondMempoolExpiry\":0}","signature":"3045022100aa5ce43cd1c067a709d75afd41aded297d223c75e6ab8baf0b06c2c3b011de4502200eefb71189facc66dc96942dda487f94a5d826be891382ea040e6c7ed2796614","publicKey":"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270"}`)))
		}

	}

	// Valid (get status)
	if strings.Contains(req.URL.String(), "mapi/"+testMinerID+"/tx/294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1aa") && req.Method == http.MethodGet {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"providerName":"taal","providerId":"ba001df8","status":{"apiVersion":"0.1.0","timestamp":"2020-06-23T20:20:10.195Z","returnResult":"success","resultDescription":"","blockHash":"000000000000000004b5ce6670f2ff27354a1e87d0a01bf61f3307f4ccd358b5","blockHeight":612251,"confirmations":28408,"minerId":"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270","txSecondMempoolExpiry":0},"payload":"{\"apiVersion\":\"0.1.0\",\"timestamp\":\"2020-06-23T20:20:10.195Z\",\"returnResult\":\"success\",\"resultDescription\":\"\",\"blockHash\":\"000000000000000004b5ce6670f2ff27354a1e87d0a01bf61f3307f4ccd358b5\",\"blockHeight\":612251,\"confirmations\":28408,\"minerId\":\"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270\",\"txSecondMempoolExpiry\":0}","signature":"3045022100a4fc6b3bc0762f46e59b80e17072d3b46a01410926862f1bd05345a3c1ebee9502201ef7385f54567d811b33e800e0b76bd3d636c329bed1cdaadd7af50f8d157f57","publicKey":"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270"}`)))
	}

	// Default is valid
	return resp, nil
}

// mockHTTPMerchantInvalid for mocking requests
type mockHTTPMerchantInvalid struct{}

// Do is a mock http request
func (m *mockHTTPMerchantInvalid) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Invalid
	if strings.Contains(req.URL.String(), "/mapi/feeQuotes") {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
		return resp, fmt.Errorf("bad request")
	}

	// Invalid
	if strings.Contains(req.URL.String(), "/mapi/"+testMinerID+"/tx") {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
		return resp, fmt.Errorf("bad request")
	}

	// Invalid
	if strings.Contains(req.URL.String(), "mapi/"+testMinerID+"/tx/error") {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
		return resp, fmt.Errorf("bad request")
	}

	// Default is valid
	return resp, nil
}

// TestClient_GetFeeQuotes tests the GetFeeQuotes()
func TestClient_GetFeeQuotes(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPMerchantValid{})

	// Test the valid response
	info, err := client.GetFeeQuotes()
	if err != nil {
		t.Errorf("%s Failed: error [%s]", t.Name(), err.Error())
	} else if info == nil {
		t.Errorf("%s Failed: info was nil", t.Name())
	} else if len(info.Quotes) != 2 {
		t.Errorf("%s Failed: expected 2 quotes, got: %d", t.Name(), len(info.Quotes))
	} else if info.Quotes[0].Quote.Fees[0].RelayFee.Satoshis != 25 {
		t.Errorf("%s Failed: expected 25 satoshis, got: %d", t.Name(), info.Quotes[0].Quote.Fees[0].RelayFee.Satoshis)
	}

	// New invalid mock client
	client = newMockClient(&mockHTTPMerchantInvalid{})

	// Test invalid response
	_, err = client.GetFeeQuotes()
	if err == nil {
		t.Errorf("%s Failed: error should have occurred", t.Name())
	}
}

// TestClient_SubmitTransaction tests the SubmitTransaction()
func TestClient_SubmitTransaction(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPMerchantValid{})

	// Create the list of tests
	var tests = []struct {
		input         string
		provider      string
		expected      string
		expectedError bool
		statusCode    int
	}{
		{"020000000232900e9b0e359cb95ad3853e1450591fdf01e3efaa1b3b2b5ab5c5ef784946b8010000006a473044022055ec4f9b9cbdd97cf5f4893f921da75287483edb4ebba4f5b23231577212fd5f022007ed037ab7da039e0d4cf0fa35f620f2d7f71285959dcb6c885652d843d0038741210232b357c5309644cf4aa72b9b2d8bfe58bdf2515d40119318d5cb51ef378cae7efffffffff91d2feda79506806c5b0dc74b6fa6ae42fb7963da460ac256383fce498b9952020000006b483045022100ac75defcda55d644b6c095c2e7cbded92e38ea52e4d7f561f37962aa036a92fe0220059c0071d7c53cc964f641fa60ba2a076e3030ae6edc0f650792041c7691c66641210282d7e568e56f59e01a4edae297ac26caabc4684971ac6c7558c91c0fa84002f7ffffffff03322a093f000000001976a914c4263eb96d88849f498d139424b59a0cba1005e888ac2f92bb09000000001976a9146cbff9881ac47da8cb699e4543c28f9b3d6941da88ac404b4c00000000001976a914f7899faf1696892e6cb029b00c713f044761f03588ac00000000", testMinerID, "", false, http.StatusOK},
		{"0100000001c871b60f86ad5d292ba51f9e94241fa2316aa4180dbdbae8286eec0bec679f01010000006a473044022100fafb74a5c6760ff2135cdf9ceee8e0342a24a6cc40804bd3245fe4aa055c8611021f6a07c86caa00d4bf92437d40a4b516aa6b7e870ff8e2a07782a5968a72ed99412103ce0ba24104223d4b1258833bacc147b27f6bd4361f403ec5cb1ee5efd0858759ffffffff0260dc2c00000000001976a914493e769b46c077954a376f9a0933f5bb9376482688ac4b807100000000001976a9141af4f2ec0ced0dece5c689e4afb28318e24b1f2d88ac00000000", testMinerID, "bbce41682bfb88d037de116996781c61235a15e6f530597075d5d69bb5a3789a", false, http.StatusOK},
	}

	// Test all
	for _, test := range tests {
		if output, err := client.SubmitTransaction(test.provider, test.input); err == nil && test.expectedError {
			t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted and [%s] expected", t.Name(), test.input, test.expected)
		} else if err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%v] error [%s]", t.Name(), test.input, test.expected, output, err.Error())
		} else if output != nil && output.Response.TxID != test.expected && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%s]", t.Name(), test.input, test.expected, output.Response.TxID)
		} else if client.LastRequest.StatusCode != test.statusCode {
			t.Errorf("%s Expected status code to be %d, got %d, [%s] inputted", t.Name(), test.statusCode, client.LastRequest.StatusCode, test.input)
		}
	}

	// New invalid mock client
	client = newMockClient(&mockHTTPMerchantInvalid{})

	// Test invalid response
	_, err := client.SubmitTransaction(testMinerID, "error")
	if err == nil {
		t.Errorf("%s Failed: error should have occurred", t.Name())
	}
}

// TestClient_TransactionStatus tests the TransactionStatus()
func TestClient_TransactionStatus(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPMerchantValid{})

	// Create the list of tests
	var tests = []struct {
		input         string
		provider      string
		expected      string
		expectedError bool
		statusCode    int
	}{
		{"294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1aa", testMinerID, "000000000000000004b5ce6670f2ff27354a1e87d0a01bf61f3307f4ccd358b5", false, http.StatusOK},
	}

	// Test all
	for _, test := range tests {
		if output, err := client.TransactionStatus(test.provider, test.input); err == nil && test.expectedError {
			t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted and [%s] expected", t.Name(), test.input, test.expected)
		} else if err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%v] error [%s]", t.Name(), test.input, test.expected, output, err.Error())
		} else if output != nil && output.Status.BlockHash != test.expected && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%s]", t.Name(), test.input, test.expected, output.Status.BlockHash)
		} else if client.LastRequest.StatusCode != test.statusCode {
			t.Errorf("%s Expected status code to be %d, got %d, [%s] inputted", t.Name(), test.statusCode, client.LastRequest.StatusCode, test.input)
		}
	}

	// New invalid mock client
	client = newMockClient(&mockHTTPMerchantInvalid{})

	// Test invalid response
	_, err := client.TransactionStatus(testMinerID, "error")
	if err == nil {
		t.Errorf("%s Failed: error should have occurred", t.Name())
	}
}
