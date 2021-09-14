package polynym

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// mockHTTP for mocking requests
type mockHTTP struct{}

// Do is a mock http request
func (m *mockHTTP) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	if strings.Contains(req.URL.String(), "/19gKzz8XmFDyrpk4qFobG7qKoqybe78v9h") {

		// Valid BSV Address
		resp.StatusCode = http.StatusOK
		resp.Body = validResponse("19gKzz8XmFDyrpk4qFobG7qKoqybe78v9h", req.URL.String(), resp.StatusCode)

	} else if strings.Contains(req.URL.String(), "/error") {

		// Return an error
		return resp, fmt.Errorf("req error")

	} else if strings.Contains(req.URL.String(), "/bad-poly-response") {

		// Return a bad error response from Polynym (empty)
		resp.Body = invalidResponse("", req.URL.String(), resp.StatusCode)

	} else if strings.Contains(req.URL.String(), "/bad-poly-status") {

		// Return a bad error response from Polynym (empty)
		resp.Body = invalidResponse("Some error message", req.URL.String(), http.StatusUnavailableForLegalReasons)

	} else if strings.Contains(req.URL.String(), "/16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA") {

		// Valid BSV Address
		resp.StatusCode = http.StatusOK
		resp.Body = validResponse("16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA", req.URL.String(), resp.StatusCode)

	} else if strings.Contains(req.URL.String(), "/c6ZqP5Tb22KJuvSAbjNkoi") {

		// Invalid BSV Address
		resp.Body = invalidResponse("Unable to resolve to address", req.URL.String(), resp.StatusCode)

	} else if strings.Contains(req.URL.String(), "/1doesnotexisthandle") {

		// Invalid handle
		resp.Body = invalidResponse("1handle not found", req.URL.String(), resp.StatusCode)

	} else if strings.Contains(req.URL.String(), "/doesnotexist@handcash.io") {

		// Invalid handle
		resp.Body = invalidResponse("$handle not found", req.URL.String(), resp.StatusCode)

	} else if strings.Contains(req.URL.String(), "/bad@paymailaddress.com") {

		// Invalid paymail
		resp.Body = invalidResponse("PayMail not found", req.URL.String(), resp.StatusCode)

	} else if strings.Contains(req.URL.String(), "/mrz@relayx.io") {

		// Valid 1handle
		resp.StatusCode = http.StatusOK
		resp.Body = validResponse("1Lti3s6AQNKTSgxnTyBREMa6XdHLBnPSKa", req.URL.String(), resp.StatusCode)

	} else if strings.Contains(req.URL.String(), "/mr-z@handcash.io") {

		// Valid $handle / paymail
		resp.StatusCode = http.StatusOK
		resp.Body = validResponse("124dwBFyFtkcNXGfVWQroGcT9ybnpQ3G3Z", req.URL.String(), resp.StatusCode)

	} else if strings.Contains(req.URL.String(), "/mrz@handcash.io") {

		// Valid paymail
		resp.StatusCode = http.StatusOK
		resp.Body = validResponse("19gKzz8XmFDyrpk4qFobG7qKoqybe78v9h", req.URL.String(), resp.StatusCode)
	} else if strings.Contains(req.URL.String(), "/@833") {

		// Valid Twetch ID
		resp.StatusCode = http.StatusOK
		resp.Body = validResponse("19ksW6ueSw9nEj88X3QNJ9VkKPGf1zuKbQ", req.URL.String(), resp.StatusCode)
	}

	return resp, nil
}

// validResponse returns a valid polynym response
func validResponse(address, url string, status int) io.ReadCloser {
	result := &GetAddressResponse{
		Address: address,
		LastRequest: &LastRequest{
			Method:     http.MethodGet,
			StatusCode: status,
			URL:        url,
		},
	}

	b, _ := json.Marshal(result)
	return ioutil.NopCloser(bytes.NewBuffer(b))
}

// invalidResponse returns an invalid polynym response (error)
func invalidResponse(errorMessage, url string, status int) io.ReadCloser {
	result := &GetAddressResponse{
		ErrorMessage: errorMessage,
		LastRequest: &LastRequest{
			Method:     http.MethodGet,
			StatusCode: status,
			URL:        url,
		},
	}

	b, _ := json.Marshal(result)
	return ioutil.NopCloser(bytes.NewBuffer(b))
}
