package whatsonchain

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

// mockHTTPExchangeValid for mocking requests
type mockHTTPExchangeValid struct{}

// Do is a mock http request
func (m *mockHTTPExchangeValid) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Valid (exchange rate)
	if strings.Contains(req.URL.String(), "/exchangerate") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"currency": "USD","rate": "178.59733333333335"}`)))
	}

	// Default is valid
	return resp, nil
}

// mockHTTPExchangeInvalid for mocking requests
type mockHTTPExchangeInvalid struct{}

// Do is a mock http request
func (m *mockHTTPExchangeInvalid) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Invalid (exchange rate)
	if strings.Contains(req.URL.String(), "/exchangerate") {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
		return resp, fmt.Errorf("bad request")
	}

	// Default is valid
	return resp, nil
}

// TestClient_GetExchangeRate tests the GetExchangeRate()
func TestClient_GetExchangeRate(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPExchangeValid{})

	// Test the valid response
	info, err := client.GetExchangeRate()
	if err != nil {
		t.Errorf("%s Failed: error [%s]", t.Name(), err.Error())
	} else if info == nil {
		t.Errorf("%s Failed: info was nil", t.Name())
	} else if info.Currency != "USD" {
		t.Errorf("%s Failed: currency was [%s] expected [%s]", t.Name(), info.Currency, "USD")
	} else if info.Rate != "178.59733333333335" {
		t.Errorf("%s Failed: currency was [%s] expected [%s]", t.Name(), info.Rate, "178.59733333333335")
	}

	// New invalid mock client
	client = newMockClient(&mockHTTPExchangeInvalid{})

	// Test invalid response
	_, err = client.GetExchangeRate()
	if err == nil {
		t.Errorf("%s Failed: error should have occurred", t.Name())
	}
}
