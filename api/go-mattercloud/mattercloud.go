/*
Package mattercloud is the unofficial golang implementation for the MatterCloud API
*/
package mattercloud

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// NewClient creates a new client to submit requests
// Parameters values are set to the defaults defined by the API documentation.
//
// For more information: https://developers.mattercloud.net/#authentication
func NewClient(apiKey string, network NetworkType, clientOptions *Options) (c *Client, err error) {

	// Make sure we have an API key
	if len(apiKey) == 0 {
		err = fmt.Errorf("missing required api key")
		return
	}

	// Create a client using the given options
	c = createClient(clientOptions)

	// Set the key and network
	c.Parameters.apiKey = apiKey
	c.Parameters.Network = network

	return
}

// Request is a generic request wrapper that can be used without constraints
func (c *Client) Request(endpoint string, method string, payload []byte) (response string, err error) {

	// Set reader
	var bodyReader io.Reader

	// Add the network value
	endpoint = fmt.Sprintf("%s%s/%s", apiEndpoint, c.Parameters.Network, endpoint)

	// Switch on Methods
	switch method {
	case http.MethodPost, http.MethodPut:
		{
			bodyReader = bytes.NewBuffer(payload)
		}
	}

	// Store for debugging purposes
	c.LastRequest.Method = method
	c.LastRequest.URL = endpoint

	// Start the request
	var request *http.Request
	if request, err = http.NewRequestWithContext(
		context.Background(), method, endpoint, bodyReader,
	); err != nil {
		return
	}

	// Change the header (user agent is in case they block default Go user agents)
	request.Header.Set("User-Agent", c.Parameters.UserAgent)
	request.Header.Set(apiKeyField, c.Parameters.apiKey)

	// Set the content type on Method
	if method == http.MethodPost || method == http.MethodPut {
		request.Header.Set("Content-Type", "application/json")
	}

	// Fire the http request
	var resp *http.Response
	if resp, err = c.httpClient.Do(request); err != nil {
		return
	}

	// Close the response body
	defer func() {
		_ = resp.Body.Close()
	}()

	// Save the status
	c.LastRequest.StatusCode = resp.StatusCode

	// Read the body
	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	// Parse the response
	response = string(body)
	return
}
