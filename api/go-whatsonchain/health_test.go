package whatsonchain

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

// mockHTTPHealthValid for mocking requests
type mockHTTPHealthValid struct{}

// Do is a mock http request
func (m *mockHTTPHealthValid) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Valid
	if strings.Contains(req.URL.String(), "/woc") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`Whats On Chain`)))
	}

	// Default is valid
	return resp, nil
}

// mockHTTPHealthInvalid for mocking requests
type mockHTTPHealthInvalid struct{}

// Do is a mock http request
func (m *mockHTTPHealthInvalid) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Invalid
	if strings.Contains(req.URL.String(), "/woc") {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
		return resp, fmt.Errorf("bad request")
	}

	// Default is valid
	return resp, nil
}

// TestClient_GetHealth tests the GetHealth()
func TestClient_GetHealth(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPHealthValid{})

	// Test the valid response
	info, err := client.GetHealth()
	if err != nil {
		t.Errorf("%s Failed: error [%s]", t.Name(), err.Error())
	} else if info != "Whats On Chain" {
		t.Errorf("%s Failed: response was [%s] expected [%s]", t.Name(), info, "Whats On Chain")
	}

	// New invalid mock client
	client = newMockClient(&mockHTTPHealthInvalid{})

	// Test invalid response
	_, err = client.GetHealth()
	if err == nil {
		t.Errorf("%s Failed: error should have occurred", t.Name())
	}
}
