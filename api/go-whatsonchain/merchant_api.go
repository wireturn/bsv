package whatsonchain

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetFeeQuotes this endpoint provides fee quotes from multiple transaction processors.
// Each quote also contains transaction processor specific txSubmissionUrl and txStatusUrl.
// These unique URLs can be used to submit transactions to the selected transaction processor and check the status of the submitted transaction.
// Any post request to txSubmissionUrl is forwarded to the selected transaction processor ‘AS IS’ and is ‘NOT’ broadcast from any WoC nodes.
//
// For more information: https://developers.whatsonchain.com/#fee-quotes
func (c *Client) GetFeeQuotes() (quotes *FeeQuotes, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/mapi/feeQuotes
	if resp, err = c.request(fmt.Sprintf("%s%s/mapi/feeQuotes", apiEndpoint, c.Network), http.MethodGet, nil); err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &quotes)
	return
}

// SubmitTransaction this endpoint submits a transaction to a specific transaction processor using the
// txSubmissionUrl provided with each quote in the Fee quotes response.
//
// For more information: https://developers.whatsonchain.com/#submit-transaction
func (c *Client) SubmitTransaction(provider string, txHex string) (response *SubmissionResponse, err error) {

	// Start the post data
	postData := []byte(fmt.Sprintf(`{"rawtx":"%s"}`, txHex))

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/mapi/<providerId>/tx
	if resp, err = c.request(fmt.Sprintf("%s%s/mapi/%s/tx", apiEndpoint, c.Network, provider), http.MethodPost, postData); err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &response)
	return
}

// TransactionStatus gets a transaction's status from a specific transaction processor using
// the txStatusUrl provided with each quote in Fee quotes response.
//
// For more information: https://developers.whatsonchain.com/#transaction-status
func (c *Client) TransactionStatus(provider string, txID string) (status *StatusResponse, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/mapi/<providerId>/tx/<hash>
	if resp, err = c.request(fmt.Sprintf("%s%s/mapi/%s/tx/%s", apiEndpoint, c.Network, provider, txID), http.MethodGet, nil); err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &status)
	return
}
