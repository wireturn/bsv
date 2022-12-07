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
  "bsvalias": "1.0",
  "capabilities": {
	"6745385c3fc0": false,
	"pki": "https://bsvalias.example.org/{alias}@{domain.tld}/id",
	"paymentDestination": "https://bsvalias.example.org/{alias}@{domain.tld}/payment-destination"
  }
}
*/

// Capabilities is the result returned (plus some custom features)
type Capabilities struct {
	StandardResponse
	BsvAlias     string                 `json:"bsvalias"`     // Version of the bsvalias
	Capabilities map[string]interface{} `json:"capabilities"` // Raw list of the capabilities
}

// Has will check if a BRFC ID (or alternate) is found in the list of capabilities
//
// Alternate is used for example: "pki" is also BRFC "0c4339ef99c2"
func (c *Capabilities) Has(brfcID, alternateID string) bool {
	for key := range c.Capabilities {
		if key == brfcID || (len(alternateID) > 0 && key == alternateID) {
			return true
		}
	}
	return false
}

// getValue will return the value (if found) from the capability (url or bool)
//
// Alternate is used for IE: pki (it breaks convention of using the BRFC ID)
func (c *Capabilities) getValue(brfcID, alternateID string) (bool, interface{}) {
	for key, val := range c.Capabilities {
		if key == brfcID || (len(alternateID) > 0 && key == alternateID) {
			return true, val
		}
	}
	return false, nil
}

// GetString will perform getValue() but cast to a string if found
//
// Returns an empty string if not found
func (c *Capabilities) GetString(brfcID, alternateID string) string {
	if ok, val := c.getValue(brfcID, alternateID); ok {
		return val.(string)
	}
	return ""
}

// GetBool will perform getValue() but cast to a bool if found
//
// Returns false if not found
func (c *Capabilities) GetBool(brfcID, alternateID string) bool {
	if ok, val := c.getValue(brfcID, alternateID); ok {
		return val.(bool)
	}
	return false
}

// GetCapabilities will return a list of capabilities for a given domain & port
//
// Specs: http://bsvalias.org/02-02-capability-discovery.html
func (c *Client) GetCapabilities(target string, port int) (response *Capabilities, err error) {

	// Basic requirements for the request
	if len(target) == 0 {
		err = fmt.Errorf("missing target")
		return
	} else if port == 0 {
		err = fmt.Errorf("missing port")
		return
	}

	// Set the base url and path
	// https://<host-discovery-target>:<host-discovery-port>/.well-known/bsvalias
	reqURL := fmt.Sprintf("https://%s:%d/.well-known/%s", target, port, DefaultServiceName)

	// Fire the GET request
	var resp StandardResponse
	if resp, err = c.getRequest(reqURL); err != nil {
		return
	}

	// Start the response
	response = &Capabilities{StandardResponse: resp}

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
	if err = json.Unmarshal(resp.Body, &response); err != nil {

		// Invalid character (sometimes quote related: U+0022 vs U+201C)
		if strings.Contains(err.Error(), "invalid character") {

			// Replace any invalid quotes
			bodyString := strings.Replace(strings.Replace(string(resp.Body), `“`, `"`, -1), `”`, `"`, -1)

			// Parse again after fixing quotes
			if err = json.Unmarshal([]byte(bodyString), &response); err != nil {
				return
			}
		}

		// Still have an error?
		if err != nil {
			return
		}
	}

	// Invalid version detected
	if len(response.BsvAlias) == 0 {
		err = fmt.Errorf("missing %s version", DefaultServiceName)
	}

	return
}
