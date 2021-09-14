package whatsonchain

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

// mockHTTPMempoolValid for mocking requests
type mockHTTPMempoolValid struct{}

// Do is a mock http request
func (m *mockHTTPMempoolValid) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Valid
	if strings.Contains(req.URL.String(), "/mempool/info") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"size": 520,"bytes": 108095,"usage": 549776,"maxmempool": 64000000000,"mempoolminfee": 0}`)))
	}

	// Valid
	if strings.Contains(req.URL.String(), "/mempool/raw") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`["86806b3587956552ea0e3f09dfd14f485fc870fa319ab37e98289a5043234644","bd9e6c83f8fdcaa3b66b214a4fbf910976bd16ec926ab983a2367edfa3e2bbd9","9cf4450a20f91419623d9b461d4e47647ce3812f0fd2e2d2904c5f5a24e45bba"]`)))
	}

	// Default is valid
	return resp, nil
}

// mockHTTPMempoolInvalid for mocking requests
type mockHTTPMempoolInvalid struct{}

// Do is a mock http request
func (m *mockHTTPMempoolInvalid) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Invalid
	if strings.Contains(req.URL.String(), "/mempool/info") {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
		return resp, fmt.Errorf("bad request")
	}

	// Invalid
	if strings.Contains(req.URL.String(), "/mempool/raw") {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
		return resp, fmt.Errorf("bad request")
	}

	// Default is valid
	return resp, nil
}

// TestClient_GetMempoolInfo tests the GetMempoolInfo()
func TestClient_GetMempoolInfo(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPMempoolValid{})

	// Test the valid response
	info, err := client.GetMempoolInfo()
	if err != nil {
		t.Errorf("%s Failed: error [%s]", t.Name(), err.Error())
	} else if info == nil {
		t.Errorf("%s Failed: info was nil", t.Name())
	} else if info.Size != 520 {
		t.Errorf("%s Failed: size was [%d] expected [%d]", t.Name(), info.Size, 520)
	} else if info.Bytes != 108095 {
		t.Errorf("%s Failed: bytes was [%d] expected [%d]", t.Name(), info.Bytes, 108095)
	} else if info.MaxMempool != 64000000000 {
		t.Errorf("%s Failed: max mempool was [%d] expected [%d]", t.Name(), info.MaxMempool, 64000000000)
	}

	// New invalid mock client
	client = newMockClient(&mockHTTPMempoolInvalid{})

	// Test invalid response
	_, err = client.GetMempoolInfo()
	if err == nil {
		t.Errorf("%s Failed: error should have occurred", t.Name())
	}
}

// TestClient_GetMempoolTransactions tests the GetMempoolTransactions()
func TestClient_GetMempoolTransactions(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPMempoolValid{})

	// Test the valid response
	info, err := client.GetMempoolTransactions()
	if err != nil {
		t.Errorf("%s Failed: error [%s]", t.Name(), err.Error())
	} else if info == nil {
		t.Errorf("%s Failed: info was nil", t.Name())
	} else if info[0] != "86806b3587956552ea0e3f09dfd14f485fc870fa319ab37e98289a5043234644" {
		t.Errorf("%s Failed: tx was [%s] expected [%s]", t.Name(), info[0], "86806b3587956552ea0e3f09dfd14f485fc870fa319ab37e98289a5043234644")
	} else if info[1] != "bd9e6c83f8fdcaa3b66b214a4fbf910976bd16ec926ab983a2367edfa3e2bbd9" {
		t.Errorf("%s Failed: tx was [%s] expected [%s]", t.Name(), info[1], "bd9e6c83f8fdcaa3b66b214a4fbf910976bd16ec926ab983a2367edfa3e2bbd9")
	}

	// New invalid mock client
	client = newMockClient(&mockHTTPMempoolInvalid{})

	// Test invalid response
	_, err = client.GetMempoolTransactions()
	if err == nil {
		t.Errorf("%s Failed: error should have occurred", t.Name())
	}
}
