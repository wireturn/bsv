package paymail

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

/*
Default:
{
    "avatar": "https://<domain><image>",
    "name": "<name>"
}
*/

// PublicProfile is the result returned from GetPublicProfile()
type PublicProfile struct {
	StandardResponse
	Avatar string `json:"avatar"` // A URL that returns a 180x180 image. It can accept an optional parameter `s` to return an image of width and height `s`. The image should be JPEG, PNG, or GIF.
	Name   string `json:"name"`   // A string up to 100 characters long. (name or nickname)
}

// GetPublicProfile will return a valid public profile
//
// Specs: https://github.com/bitcoin-sv-specs/brfc-paymail/pull/7/files
func (c *Client) GetPublicProfile(publicProfileURL, alias, domain string) (response *PublicProfile, err error) {

	// Require a valid url
	if len(publicProfileURL) == 0 || !strings.Contains(publicProfileURL, "https://") {
		err = fmt.Errorf("invalid url: %s", publicProfileURL)
		return
	}

	// Basic requirements for request
	if len(alias) == 0 {
		err = fmt.Errorf("missing alias")
		return
	} else if len(domain) == 0 {
		err = fmt.Errorf("missing domain")
		return
	}

	// Set the base url and path, assuming the url is from the prior GetCapabilities() request
	// https://<host-discovery-target>/public-profile/{alias}@{domain.tld}
	reqURL := strings.Replace(strings.Replace(publicProfileURL, "{alias}", alias, -1), "{domain.tld}", domain, -1)

	// Fire the GET request
	var resp StandardResponse
	if resp, err = c.getRequest(reqURL); err != nil {
		return
	}

	// Start the response
	response = &PublicProfile{StandardResponse: resp}

	// Test the status code (200 or 304 is valid)
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotModified {
		serverError := &ServerError{}
		if err = json.Unmarshal(resp.Body, serverError); err != nil {
			return
		}
		err = fmt.Errorf("bad response from paymail provider: code %d, message: %s", response.StatusCode, serverError.Message)
		return
	}

	// Decode the body of the response
	err = json.Unmarshal(resp.Body, &response)

	return
}
