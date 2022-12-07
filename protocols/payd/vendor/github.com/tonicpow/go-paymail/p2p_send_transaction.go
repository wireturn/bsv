package paymail

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

/*
Example:
{
  "hex": "01000000012adda020db81f2155ebba69e7.........154888ac00000000",
  "metadata": {
	"sender": "someone@example.tld",
	"pubkey": "<sender-pubkey>",
	"signature": "signature(txid)",
	"note": "Human readable information related to the tx."
  },
  "reference": "someRefId"
}
*/

// P2PTransaction is the request body for the P2P transaction request
type P2PTransaction struct {
	Hex       string       `json:"hex"`       // The raw transaction, encoded as a hexadecimal string
	MetaData  *P2PMetaData `json:"metadata"`  // An object containing data associated with the transaction
	Reference string       `json:"reference"` // Reference for the payment (from previous P2P Destination request)
}

// P2PMetaData is an object containing data associated with the P2P transaction
type P2PMetaData struct {
	Note      string `json:"note,omitempty"`      // A human readable bit of information about the payment
	PubKey    string `json:"pubkey,omitempty"`    // Public key to validate the signature (if signature is given)
	Sender    string `json:"sender,omitempty"`    // The paymail of the person that originated the transaction
	Signature string `json:"signature,omitempty"` // A signature of the tx id made by the sender
}

// P2PTransactionResponse is the response to the request
type P2PTransactionResponse struct {
	StandardResponse
	Note string `json:"note"` // Some human readable note
	TxID string `json:"txid"` // The txid of the broadcasted tx
}

// SendP2PTransaction will submit a transaction hex string (tx_hex) to a paymail provider
//
// Specs: https://docs.moneybutton.com/docs/paymail-06-p2p-transactions.html
func (c *Client) SendP2PTransaction(p2pURL, alias, domain string, transaction *P2PTransaction) (response *P2PTransactionResponse, err error) {

	// Require a valid url
	if len(p2pURL) == 0 || !strings.Contains(p2pURL, "https://") {
		err = fmt.Errorf("invalid url: %s", p2pURL)
		return
	} else if len(alias) == 0 {
		err = fmt.Errorf("missing alias")
		return
	} else if len(domain) == 0 {
		err = fmt.Errorf("missing domain")
		return
	}

	// Basic requirements for request
	if transaction == nil {
		err = fmt.Errorf("transaction cannot be nil")
		return
	} else if len(transaction.Hex) == 0 {
		err = fmt.Errorf("hex is required")
		return
	} else if len(transaction.Reference) == 0 {
		err = fmt.Errorf("reference is required")
		return
	}

	// Set the base url and path, assuming the url is from the prior GetCapabilities() request
	// https://<host-discovery-target>/api/rawtx/{alias}@{domain.tld}
	// https://<host-discovery-target>/api/receive-transaction/{alias}@{domain.tld}
	reqURL := strings.Replace(strings.Replace(p2pURL, "{alias}", alias, -1), "{domain.tld}", domain, -1)

	// Fire the POST request
	var resp StandardResponse
	if resp, err = c.postRequest(reqURL, transaction); err != nil {
		return
	}

	// Start the response
	response = &P2PTransactionResponse{StandardResponse: resp}

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

	// Check for a reference number
	if len(response.TxID) == 0 {
		err = fmt.Errorf("missing a returned txid")
		return
	}

	return
}
