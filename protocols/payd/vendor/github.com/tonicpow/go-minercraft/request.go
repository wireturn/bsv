package minercraft

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
	Method string `json:"method"`
	URL    string `json:"url"`
	Token  string `json:"token"`
	Data   []byte `json:"data"`
}

// httpRequest is a generic request wrapper that can be used without constraints
func httpRequest(ctx context.Context, client *Client,
	payload *httpPayload) (response *RequestResponse) {

	// Set reader
	var bodyReader io.Reader

	// Start the response
	response = new(RequestResponse)

	// Add post data if applicable
	if payload.Method == http.MethodPost || payload.Method == http.MethodPut {
		bodyReader = bytes.NewBuffer(payload.Data)
		response.PostData = string(payload.Data)
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

	// Set a token if supplied
	if len(payload.Token) > 0 {
		request.Header.Set("Authorization", payload.Token)
	}

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

	if resp.Body != nil {
		// Read the body
		response.BodyContents, response.Error = ioutil.ReadAll(resp.Body)
	}
	// Check status code
	if http.StatusOK == resp.StatusCode {
		return
	}
	// unexpected status, write an error.
	if response.BodyContents == nil {
		// There's no "body" present, so just echo status code.
		response.Error = fmt.Errorf(
			"status code: %d does not match %d",
			resp.StatusCode, http.StatusOK,
		)
		return
	}
	// Have a "body" so map to an error type and add to the error message.
	errBody := struct {
		Error string `json:"error"`
	}{}
	if err := json.Unmarshal(response.BodyContents, &errBody); err != nil {
		response.Error = fmt.Errorf("failed to unmarshal mapi error response: %w", err)
		return
	}
	response.Error = fmt.Errorf(
		"status code: %d does not match %d, error: %s",
		resp.StatusCode, http.StatusOK, errBody.Error,
	)
	return
}
