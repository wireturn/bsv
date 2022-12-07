package minercraft

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

// QueryTransactionSuccess is on success
const QueryTransactionSuccess = "success"

// QueryTransactionFailure is on failure
const QueryTransactionFailure = "failure"

// QueryTransactionInMempoolFailure in mempool but not in a block yet
const QueryTransactionInMempoolFailure = "Transaction in mempool but not yet in block"

/*
Example query tx response from Merchant API:

{
  "payload": "{\"apiVersion\":\"1.2.3\",\"timestamp\":\"2020-01-15T11:41:29.032Z\",\"returnResult\":\"failure\",\"resultDescription\":\"Transaction in mempool but not yet in block\",\"blockHash\":\"\",\"blockHeight\":0,\"minerId\":\"03fcfcfcd0841b0a6ed2057fa8ed404788de47ceb3390c53e79c4ecd1e05819031\",\"confirmations\":0,\"txSecondMempoolExpiry\":0}",
  "signature": "3045022100f78a6ac49ef38fbe68db609ff194d22932d865d93a98ee04d2ecef5016872ba50220387bf7e4df323bf4a977dd22a34ea3ad42de1a2ec4e5af59baa13258f64fe0e5",
  "publicKey": "03fcfcfcd0841b0a6ed2057fa8ed404788de47ceb3390c53e79c4ecd1e05819031",
  "encoding": "UTF-8",
  "mimetype": "application/json"
}
*/

// QueryTransactionResponse is the raw response from the Merchant API request
//
// Specs: https://github.com/bitcoin-sv-specs/brfc-merchantapi/tree/v1.2-beta#Query-transaction-status
type QueryTransactionResponse struct {
	JSONEnvelope
	Query *QueryPayload `json:"query"` // Custom field for unmarshalled payload data
}

/*
Example QueryTransactionResponse.Payload (unmarshalled):

Failure - in mempool but not in block
{
  "apiVersion": "1.2.3",
  "timestamp": "2020-01-15T11:41:29.032Z",
  "txid": "6bdbcfab0526d30e8d68279f79dff61fb4026ace8b7b32789af016336e54f2f0",
  "returnResult": "failure",
  "resultDescription": "Transaction in mempool but not yet in block",
  "blockHash": "",
  "blockHeight": 0,
  "minerId": "03fcfcfcd0841b0a6ed2057fa8ed404788de47ceb3390c53e79c4ecd1e05819031",
  "confirmations": 0,
  "txSecondMempoolExpiry": 0
}

Success - added to block
{
  "apiVersion": "1.2.3",
  "timestamp": "2020-01-15T12:09:37.394Z",
  "txid": "6bdbcfab0526d30e8d68279f79dff61fb4026ace8b7b32789af016336e54f2f0",
  "returnResult": "success",
  "resultDescription": "",
  "blockHash": "745093bb0c80780092d4ce6926e0caa753fe3accdc09c761aee89bafa85f05f4",
  "blockHeight": 208,
  "minerId": "03fcfcfcd0841b0a6ed2057fa8ed404788de47ceb3390c53e79c4ecd1e05819031",
  "confirmations": 2,
  "txSecondMempoolExpiry": 0
}
*/

// QueryPayload is the unmarshalled version of the payload envelope
type QueryPayload struct {
	APIVersion            string `json:"apiVersion"`
	Timestamp             string `json:"timestamp"`
	TxID                  string `json:"txid"`
	ReturnResult          string `json:"returnResult"`
	ResultDescription     string `json:"resultDescription"`
	BlockHash             string `json:"blockHash"`
	BlockHeight           int64  `json:"blockHeight"`
	MinerID               string `json:"minerId"`
	Confirmations         int64  `json:"confirmations"`
	TxSecondMempoolExpiry int64  `json:"txSecondMempoolExpiry"`
}

// QueryTransaction will fire a Merchant API request to check the status of a transaction
//
// This endpoint is used to check the current status of a previously submitted transaction.
// It returns a JSONEnvelope with a payload that contains the transaction status.
// The purpose of the envelope is to ensure strict consistency in the message content for
// the purpose of signing responses.
//
// Specs: https://github.com/bitcoin-sv-specs/brfc-merchantapi/tree/v1.2-beta#Query-transaction-status
func (c *Client) QueryTransaction(ctx context.Context, miner *Miner, txID string) (*QueryTransactionResponse, error) {

	// Make sure we have a valid miner
	if miner == nil {
		return nil, errors.New("miner was nil")
	}

	// Make the HTTP request
	result := queryTransaction(ctx, c, miner, txID)
	if result.Response.Error != nil {
		return nil, result.Response.Error
	}

	// Parse the response
	response, err := result.parseQuery()
	if err != nil {
		return nil, err
	}

	// Valid?
	if response.Query == nil || len(response.Query.ReturnResult) == 0 {
		return nil, errors.New("failed getting query response from: " + miner.Name)
	}

	// Return the fully parsed response
	return &response, nil
}

// queryTransaction will fire the HTTP request to retrieve the tx status
func queryTransaction(ctx context.Context, client *Client, miner *Miner, txHash string) (result *internalResult) {
	result = &internalResult{Miner: miner}
	result.Response = httpRequest(ctx, client, &httpPayload{
		Method: http.MethodGet,
		URL:    miner.URL + routeQueryTx + txHash,
		Token:  miner.Token,
	})
	return
}

// parseQuery will convert the HTTP response into a struct and also unmarshal the payload JSON data
func (i *internalResult) parseQuery() (response QueryTransactionResponse, err error) {

	// Process the initial response payload
	if err = response.process(i.Miner, i.Response.BodyContents); err != nil {
		return
	}

	// If we have a valid payload
	if len(response.Payload) > 0 {
		err = json.Unmarshal([]byte(response.Payload), &response.Query)
	}
	return
}
