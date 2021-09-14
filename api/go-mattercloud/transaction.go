package mattercloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Transaction this endpoint retrieves specific transaction
//
// For more information: https://developers.mattercloud.net/#get-transaction
func (c *Client) Transaction(tx string) (transaction *Transaction, err error) {

	var resp string

	// GET https://api.mattercloud.net/api/v3/main/tx/<txid>
	if resp, err = c.Request("tx/"+tx, http.MethodGet, nil); err != nil {
		return
	}

	// Check for error
	if c.LastRequest.StatusCode != http.StatusOK {
		var apiError APIInternalError
		if err = json.Unmarshal([]byte(resp), &apiError); err != nil {
			return
		}
		err = fmt.Errorf("error: %s", apiError.ErrorMessage)
		return
	}

	// Process the response
	err = json.Unmarshal([]byte(resp), &transaction)
	return
}

// TransactionBatch this endpoint retrieves details for multiple transactions at same time
//
// For more information: https://developers.mattercloud.net/#get-transaction-batch
func (c *Client) TransactionBatch(txIDs []string) (transactions []*Transaction, err error) {

	// Check ids
	if len(txIDs) == 0 {
		err = fmt.Errorf("missing tx ids")
		return
	}

	// Marshall into JSON
	var data []byte
	if data, err = json.Marshal(&TransactionList{TxIDs: strings.Join(txIDs, ",")}); err != nil {
		return
	}

	var resp string

	// POST https://api.mattercloud.net/api/v3/main/tx
	if resp, err = c.Request("tx", http.MethodPost, data); err != nil {
		return
	}

	// Check for error
	if c.LastRequest.StatusCode != http.StatusOK {
		var apiError APIInternalError
		if err = json.Unmarshal([]byte(resp), &apiError); err != nil {
			return
		}
		err = fmt.Errorf("error: %s", apiError.ErrorMessage)
		return
	}

	// Process the response
	err = json.Unmarshal([]byte(resp), &transactions)
	return
}

// Broadcast this endpoint broadcasts a raw transaction to the network
//
// For more information: https://developers.mattercloud.net/#broadcast-transaction
func (c *Client) Broadcast(rawTx string) (response *BroadcastResponse, err error) {

	var resp string
	// POST https://api.mattercloud.net/api/v3/main/tx/send

	resp, err = c.Request("tx/send", http.MethodPost, []byte(fmt.Sprintf(`{"rawtx":"%s"}`, rawTx)))
	if err != nil {
		return
	}

	// Check for error
	if c.LastRequest.StatusCode != http.StatusOK {
		var apiError APIInternalError
		if err = json.Unmarshal([]byte(resp), &apiError); err != nil {
			return
		}
		err = fmt.Errorf("error: %s", apiError.ErrorMessage)
		return
	}

	// Process the response
	err = json.Unmarshal([]byte(resp), &response)
	return
}
