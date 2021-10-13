package handcash

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// RequestResponse is the response from a request
type RequestResponse struct {
	BodyContents []byte `json:"body_contents"` // Raw body response
	Error        error  `json:"error"`         // If an error occurs
	Method       string `json:"method"`        // Method is the HTTP method used
	PostData     string `json:"post_data"`     // PostData is the post data submitted if POST/PUT request
	StatusCode   int    `json:"status_code"`   // StatusCode is the last code from the request
	URL          string `json:"url"`           // URL is used for the request
}

// httpPayload is used for a httpRequest
type httpPayload struct {
	Data           []byte `json:"data"`
	ExpectedStatus int    `json:"expected_status"`
	Method         string `json:"method"`
	URL            string `json:"url"`
}

// httpRequest is a generic request wrapper that can be used without constraints
func httpRequest(ctx context.Context, client *Client,
	payload *httpPayload, signedRequest *signedRequest) (response *RequestResponse) {

	// Set reader
	var bodyReader io.Reader

	// Start the response
	response = new(RequestResponse)

	// Add post data if applicable
	if payload.Method == http.MethodPost || payload.Method == http.MethodPut {
		bodyReader = bytes.NewBuffer(payload.Data)
		response.PostData = string(payload.Data)
	} else if payload.Method == http.MethodGet {
		// HandCash requires data even on a GET request (DO NOT REMOVE)
		if len(payload.Data) > 0 {
			bodyReader = bytes.NewBuffer(payload.Data) // empty: {}
		}
	}

	// Store for debugging purposes
	response.Method = payload.Method
	response.URL = payload.URL

	// Start the request
	var request *http.Request
	if request, response.Error = http.NewRequestWithContext(
		ctx, payload.Method, payload.URL, bodyReader,
	); response.Error != nil {
		return
	}

	// Change the header (user agent is in case they block default Go user agents)
	request.Header.Set("User-Agent", client.Options.UserAgent)

	// Set the content type on Method
	if payload.Method == http.MethodPost || payload.Method == http.MethodPut {
		request.Header.Set("Content-Type", "application/json")
	}

	// Set oAuth headers
	request.Header.Set("oauth-publickey", signedRequest.Headers.OauthPublicKey)
	request.Header.Set("oauth-signature", signedRequest.Headers.OauthSignature)
	request.Header.Set("oauth-timestamp", signedRequest.Headers.OauthTimestamp)

	// Fire the http request
	var resp *http.Response
	if resp, response.Error = client.httpClient.Do(request); response.Error != nil {
		if resp != nil {
			response.StatusCode = resp.StatusCode
		}
		return
	}

	// Close the response body
	defer func() {
		_ = resp.Body.Close()
	}()

	// Set the status
	response.StatusCode = resp.StatusCode

	// Read the body
	if response.BodyContents, response.Error = ioutil.ReadAll(resp.Body); response.Error != nil {
		return
	}

	// Status does not match as expected
	if resp.StatusCode != payload.ExpectedStatus {

		// Set the error message
		if len(response.BodyContents) > 0 {
			errorMsg := new(errorResponse)
			if response.Error = json.Unmarshal(
				response.BodyContents, &errorMsg,
			); response.Error != nil {
				return
			}
			response.Error = fmt.Errorf("%s", errorMsg.Message)
			return
		}

		// No error message found, set default error message
		response.Error = fmt.Errorf("request failed with status code: %d", resp.StatusCode)
		return
	}

	return
}
