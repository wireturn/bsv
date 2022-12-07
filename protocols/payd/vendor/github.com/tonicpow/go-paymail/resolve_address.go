package paymail

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/bitcoinschema/go-bitcoin"
)

// Resolution is the response from the ResolveAddress() request
type Resolution struct {
	StandardResponse
	Address   string `json:"address,omitempty"`   // Legacy BSV address derived from the output script (custom for our Go package)
	Output    string `json:"output"`              // hex-encoded Bitcoin script, which the sender MUST use during the construction of a payment transaction
	Signature string `json:"signature,omitempty"` // This is used if SenderValidation is enforced (signature of "output" value)
}

// ResolveAddress will return a hex-encoded Bitcoin script if successful
//
// Specs: http://bsvalias.org/04-01-basic-address-resolution.html
func (c *Client) ResolveAddress(resolutionURL, alias, domain string, senderRequest *SenderRequest) (response *Resolution, err error) {

	// Require a valid url
	if len(resolutionURL) == 0 || !strings.Contains(resolutionURL, "https://") {
		err = fmt.Errorf("invalid url: %s", resolutionURL)
		return
	}

	// Basic requirements for the request
	if len(alias) == 0 {
		err = fmt.Errorf("missing alias")
		return
	} else if len(domain) == 0 {
		err = fmt.Errorf("missing domain")
		return
	}

	// Basic requirements for request
	if senderRequest == nil {
		err = fmt.Errorf("senderReqeuest cannot be nil")
		return
	} else if len(senderRequest.Dt) == 0 {
		err = fmt.Errorf("time is required on senderReqeuest")
		return
	} else if len(senderRequest.SenderHandle) == 0 {
		err = fmt.Errorf("sender handle is required on senderReqeuest")
		return
	}

	// Set the base url and path, assuming the url is from the prior GetCapabilities() request
	// https://<host-discovery-target>/{alias}@{domain.tld}/payment-destination
	reqURL := strings.Replace(strings.Replace(resolutionURL, "{alias}", alias, -1), "{domain.tld}", domain, -1)

	// Fire the POST request
	var resp StandardResponse
	if resp, err = c.postRequest(reqURL, senderRequest); err != nil {
		return
	}

	// Start the response
	response = &Resolution{StandardResponse: resp}

	// Test the status code
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotModified {

		// Paymail address not found?
		if response.StatusCode == http.StatusNotFound {
			err = fmt.Errorf("paymail address not found")
		} else {
			serverError := &ServerError{}
			if err = json.Unmarshal(resp.Body, serverError); err != nil {
				return
			}
			err = fmt.Errorf("bad response from paymail provider: code %d, message: %s", response.StatusCode, serverError.Message)
		}

		return
	}

	// Decode the body of the response
	if err = json.Unmarshal(resp.Body, &response); err != nil {
		return
	}

	// Check for an output
	if len(response.Output) == 0 {
		err = fmt.Errorf("missing an output value")
		return
	}

	// Extract the address
	response.Address, err = bitcoin.GetAddressFromScript(response.Output)

	return
}
