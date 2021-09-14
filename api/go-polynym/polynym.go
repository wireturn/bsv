/*
Package polynym is the unofficial golang implementation for the Polynym API

Example:

// Create a new client
client, _ := polynym.NewClient(nil)

// Get address
resp, _ := client.GetAddress("1mrz")

// What's the address?
log.Println("address:", resp.Address)
*/
package polynym

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GetAddressResponse is what polynym returns (success or fail)
type GetAddressResponse struct {
	Address      string       `json:"address"`
	ErrorMessage string       `json:"error"`
	LastRequest  *LastRequest `json:"-"`
}

// GetAddress returns the address of a given 1handle, $handcash, paymail, Twetch user id or BitcoinSV address
func GetAddress(client Client, handleOrPaymail string) (response *GetAddressResponse, err error) {

	// Convert handle to paymail if detected
	if strings.Contains(handleOrPaymail, "$") {
		handleOrPaymail = HandCashConvert(handleOrPaymail, false)
	} else if strings.HasPrefix(handleOrPaymail, "1") && len(handleOrPaymail) < 25 {
		handleOrPaymail = RelayXConvert(handleOrPaymail)
	}

	// Set the API url
	// todo: beta is temporary, and only used via the method directly
	reqURL := fmt.Sprintf("%s/%s/%s", apiEndpoint, "getAddress", handleOrPaymail)

	// Store for debugging purposes
	response = &GetAddressResponse{
		LastRequest: &LastRequest{
			Method: http.MethodGet,
			URL:    reqURL,
		},
	}

	// Check for a value
	if len(handleOrPaymail) == 0 {
		response.LastRequest.StatusCode = http.StatusBadRequest
		err = fmt.Errorf("missing handle or paymail to resolve")
		return
	}

	// Start the request
	var req *http.Request
	if req, err = http.NewRequestWithContext(context.Background(), http.MethodGet, reqURL, nil); err != nil {
		return
	}

	// Set the header (user agent is in case they block default Go user agents)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", client.UserAgent)

	// Fire the request
	var resp *http.Response
	if resp, err = client.httpClient.Do(req); err != nil {
		if resp != nil {
			response.LastRequest.StatusCode = resp.StatusCode
		}
		return
	}

	// Cleanup
	defer func() {
		if resp.Body != nil {
			_ = resp.Body.Close()
		}

	}()

	// Set the status
	response.LastRequest.StatusCode = resp.StatusCode

	// Handle errors
	if resp.StatusCode != http.StatusOK {

		// Decode the error message
		if resp.StatusCode == http.StatusBadRequest {
			if resp.Body == nil {
				err = fmt.Errorf("no response body found")
				return
			}
			if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
				return
			}
			if len(response.ErrorMessage) == 0 {
				response.ErrorMessage = "unknown error resolving address"
			}
			err = fmt.Errorf("%s", response.ErrorMessage)
		} else {
			err = fmt.Errorf("bad response from polynym: %d", resp.StatusCode)
		}

		return
	}

	// Try and decode the response
	err = json.NewDecoder(resp.Body).Decode(&response)

	return
}

// HandCashConvert now converts $handle to paymail: handle@handcash.io or handle@beta.handcash.io
func HandCashConvert(handle string, isBeta bool) string {
	if strings.HasPrefix(handle, "$") {
		if isBeta {
			return strings.ToLower(strings.Replace(handle, "$", "", -1)) + "@beta.handcash.io"
		}
		return strings.ToLower(strings.Replace(handle, "$", "", -1)) + "@handcash.io"
	}
	return handle
}

// RelayXConvert now converts 1handle to paymail: handle@relayx.io
func RelayXConvert(handle string) string {
	if strings.HasPrefix(handle, "1") && len(handle) < 25 {
		return strings.ToLower(strings.Replace(handle, "1", "", -1)) + "@relayx.io"
	}
	return handle
}
