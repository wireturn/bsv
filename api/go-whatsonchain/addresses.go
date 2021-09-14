package whatsonchain

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AddressInfo this endpoint retrieves various address info.
//
// For more information: https://developers.whatsonchain.com/#address
func (c *Client) AddressInfo(address string) (addressInfo *AddressInfo, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/address/<address>/info
	if resp, err = c.request(fmt.Sprintf("%s%s/address/%s/info", apiEndpoint, c.Network, address), http.MethodGet, nil); err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &addressInfo)
	return
}

// AddressBalance this endpoint retrieves confirmed and unconfirmed address balance.
//
// For more information: https://developers.whatsonchain.com/#get-balance
func (c *Client) AddressBalance(address string) (balance *AddressBalance, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/address/<address>/balance
	if resp, err = c.request(fmt.Sprintf("%s%s/address/%s/balance", apiEndpoint, c.Network, address), http.MethodGet, nil); err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &balance)
	return
}

// AddressHistory this endpoint retrieves confirmed and unconfirmed address transactions.
//
// For more information: https://developers.whatsonchain.com/#get-history
func (c *Client) AddressHistory(address string) (history AddressHistory, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/address/<address>/history
	if resp, err = c.request(fmt.Sprintf("%s%s/address/%s/history", apiEndpoint, c.Network, address), http.MethodGet, nil); err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &history)
	return
}

// AddressUnspentTransactions this endpoint retrieves ordered list of UTXOs.
//
// For more information: https://developers.whatsonchain.com/#get-unspent-transactions
func (c *Client) AddressUnspentTransactions(address string) (history AddressHistory, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/address/<address>/unspent
	if resp, err = c.request(fmt.Sprintf("%s%s/address/%s/unspent", apiEndpoint, c.Network, address), http.MethodGet, nil); err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &history)
	return
}

// AddressUnspentTransactionDetails this endpoint retrieves transaction details for a given address
// Use max transactions to filter if there are more UTXOs returned than needed by the user
//
// For more information: (custom request for this go package)
func (c *Client) AddressUnspentTransactionDetails(address string, maxTransactions int) (history AddressHistory, err error) {

	// Get the address UTXO history
	var utxos AddressHistory
	if utxos, err = c.AddressUnspentTransactions(address); err != nil {
		return
	} else if len(utxos) == 0 {
		return
	}

	// Do we have a "custom max" amount?
	if maxTransactions > 0 {
		total := len(utxos)
		if total > maxTransactions {
			utxos = utxos[:total-(total-maxTransactions)]
		}
	}

	// Break up the UTXOs into batches
	var batches []AddressHistory
	chunkSize := MaxTransactionsUTXO

	for i := 0; i < len(utxos); i += chunkSize {
		end := i + chunkSize

		if end > len(utxos) {
			end = len(utxos)
		}

		batches = append(batches, utxos[i:end])
	}

	// todo: use channels/wait group to fire all requests at the same time (rate limiting)

	// Loop Batches - and get each batch (multiple batches of MaxTransactionsUTXO)
	for _, batch := range batches {

		txHashes := new(TxHashes)

		// Loop the batch (max MaxTransactionsUTXO)
		for _, utxo := range batch {

			// Append to the list to send and return
			txHashes.TxIDs = append(txHashes.TxIDs, utxo.TxHash)
			history = append(history, utxo)
		}

		// Get the tx details (max of MaxTransactionsUTXO)
		var txList TxList
		if txList, err = c.BulkTransactionDetails(txHashes); err != nil {
			return
		}

		// Add to the history list
		for index, tx := range txList {
			for _, utxo := range history {
				if utxo.TxHash == tx.TxID {
					utxo.Info = txList[index]
					continue
				}
			}
		}
	}

	return
}

// DownloadStatement this endpoint downloads an address statement (PDF)
// The contents will be returned in plain-text and need to be converted to a file.pdf
//
// For more information: https://developers.whatsonchain.com/#download-statement
func (c *Client) DownloadStatement(address string) (string, error) {

	// https://<network>.whatsonchain.com/statement/<hash>
	// todo: this endpoint does not follow the convention of the WOC API v1
	return c.request(fmt.Sprintf("https://%s.whatsonchain.com/statement/%s", c.Network, address), http.MethodGet, nil)
}

// bulkRequest is the common parts of the bulk requests
func bulkRequest(list *AddressList) ([]byte, error) {

	// Max limit by WOC
	if len(list.Addresses) > MaxAddressesForLookup {
		return nil, fmt.Errorf("max limit of addresses is %d and you sent %d", MaxAddressesForLookup, len(list.Addresses))
	}

	// Convert to JSON
	return json.Marshal(list)
}

// BulkBalance this endpoint retrieves confirmed and unconfirmed address balances
// Max of 20 addresses at a time
//
// For more information: https://developers.whatsonchain.com/#bulk-balance
func (c *Client) BulkBalance(list *AddressList) (balances AddressBalances, err error) {

	// Get the JSON
	var postData []byte
	if postData, err = bulkRequest(list); err != nil {
		return
	}

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/addresses/balance
	if resp, err = c.request(fmt.Sprintf("%s%s/addresses/balance", apiEndpoint, c.Network), http.MethodPost, postData); err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &balances)
	return
}

// BulkUnspentTransactions will fetch UTXOs for multiple addresses in a single request
// Max of 20 addresses at a time
//
// For more information: https://developers.whatsonchain.com/#bulk-unspent-transactions
func (c *Client) BulkUnspentTransactions(list *AddressList) (response BulkUnspentResponse, err error) {

	// Get the JSON
	var postData []byte
	if postData, err = bulkRequest(list); err != nil {
		return
	}

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/addresses/unspent
	if resp, err = c.request(fmt.Sprintf("%s%s/addresses/unspent", apiEndpoint, c.Network), http.MethodPost, postData); err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &response)
	return
}
