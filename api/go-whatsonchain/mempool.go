package whatsonchain

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetMempoolInfo this endpoint retrieves various info about the node's mempool for the selected network
//
// For more information: https://developers.whatsonchain.com/#get-mempool-info
func (c *Client) GetMempoolInfo() (info *MempoolInfo, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/mempool/info
	if resp, err = c.request(fmt.Sprintf("%s%s/mempool/info", apiEndpoint, c.Network), http.MethodGet, nil); err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &info)
	return
}

// GetMempoolTransactions this endpoint endpoint retrieve list of transaction ids from the node's mempool
// for the selected network
//
// For more information: https://developers.whatsonchain.com/#get-mempool-transactions
func (c *Client) GetMempoolTransactions() (transactions []string, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/mempool/raw
	if resp, err = c.request(fmt.Sprintf("%s%s/mempool/raw", apiEndpoint, c.Network), http.MethodGet, nil); err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &transactions)
	return
}
