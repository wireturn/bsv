package minercraft

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

const (
	// MerkleFormatTSC can be set when calling SubmitTransaction to request a MerkleProof in TSC format.
	MerkleFormatTSC = "TSC"
)

/*
Example Transaction Submission (submitted in the body of the request)
{
  "rawtx":        "[transaction_hex_string]",
  "callBackUrl":  "https://your.service.callback/endpoint",
  "callBackToken" : <channel token>,
  "merkleProof" : true,
  "dsCheck" : true,
  "callBackEncryption" : <parameter>
}
*/

// Transaction is the body contents in the "submit transaction" request
type Transaction struct {
	RawTx              string `json:"rawtx"`
	CallBackURL        string `json:"callBackUrl,omitempty"`
	CallBackToken      string `json:"callBackToken,omitempty"`
	MerkleFormat       string `json:"merkleFormat,omitempty"`
	CallBackEncryption string `json:"callBackEncryption,omitempty"`
	MerkleProof        bool   `json:"merkleProof,omitempty"`
	DsCheck            bool   `json:"dsCheck,omitempty"`
}

/*
Example submit tx response from Merchant API:

{
  "payload": "{\"apiVersion\":\"0.1.0\",\"timestamp\":\"2020-01-15T11:40:29.826Z\",\"txid\":\"6bdbcfab0526d30e8d68279f79dff61fb4026ace8b7b32789af016336e54f2f0\",\"returnResult\":\"success\",\"resultDescription\":\"\",\"minerId\":\"03fcfcfcd0841b0a6ed2057fa8ed404788de47ceb3390c53e79c4ecd1e05819031\",\"currentHighestBlockHash\":\"71a7374389afaec80fcabbbf08dcd82d392cf68c9a13fe29da1a0c853facef01\",\"currentHighestBlockHeight\":207,\"txSecondMempoolExpiry\":0}",
  "signature": "3045022100f65ae83b20bc60e7a5f0e9c1bd9aceb2b26962ad0ee35472264e83e059f4b9be022010ca2334ff088d6e085eb3c2118306e61ec97781e8e1544e75224533dcc32379",
  "publicKey": "03fcfcfcd0841b0a6ed2057fa8ed404788de47ceb3390c53e79c4ecd1e05819031",
  "encoding": "UTF-8",
  "mimetype": "application/json"
}
*/

// SubmitTransactionResponse is the raw response from the Merchant API request
//
// Specs: https://github.com/bitcoin-sv-specs/brfc-merchantapi/tree/v1.2-beta#Submit-transaction
type SubmitTransactionResponse struct {
	JSONEnvelope
	Results *SubmissionPayload `json:"results"` // Custom field for unmarshalled payload data
}

/*
Example SubmitTransactionResponse.Payload (unmarshalled):

{
  "apiVersion": "1.2.3",
  "timestamp": "2020-01-15T11:40:29.826Z",
  "txid": "6bdbcfab0526d30e8d68279f79dff61fb4026ace8b7b32789af016336e54f2f0",
  "returnResult": "success",
  "resultDescription": "",
  "minerId": "03fcfcfcd0841b0a6ed2057fa8ed404788de47ceb3390c53e79c4ecd1e05819031",
  "currentHighestBlockHash": "71a7374389afaec80fcabbbf08dcd82d392cf68c9a13fe29da1a0c853facef01",
  "currentHighestBlockHeight": 207,
  "txSecondMempoolExpiry": 0,
  "conflictedWith": ""
}
*/

// SubmissionPayload is the unmarshalled version of the payload envelope
type SubmissionPayload struct {
	APIVersion                string            `json:"apiVersion"`
	Timestamp                 string            `json:"timestamp"`
	TxID                      string            `json:"txid"`
	ReturnResult              string            `json:"returnResult"`
	ResultDescription         string            `json:"resultDescription"`
	MinerID                   string            `json:"minerId"`
	CurrentHighestBlockHash   string            `json:"currentHighestBlockHash"`
	ConflictedWith            []*ConflictedWith `json:"conflictedWith"`
	CurrentHighestBlockHeight int64             `json:"currentHighestBlockHeight"`
	TxSecondMempoolExpiry     int64             `json:"txSecondMempoolExpiry"`
}

/*
Example callback from Merchant API:
{
   "callbackPayload": "{\"index\":1,\"txOrId\":\"e7b3eefab33072e62283255f193ef5d22f26bbcfc0a80688fe2cc5178a32dda6\",\"targetType\":\"header\",\"target\":\"00000020a552fb757cf80b7341063e108884504212da2f1e1ce2ad9ffc3c6163955a27274b53d185c6b216d9f4f8831af1249d7b4b8c8ab16096cb49dda5e5fbd59517c775ba8b60ffff7f2000000000\",\"nodes\":[\"30361d1b60b8ca43d5cec3efc0a0c166d777ada0543ace64c4034fa25d253909\",\"e7aa15058daf38236965670467ade59f96cfc6ec6b7b8bb05c9a7ed6926b884d\",\"dad635ff856c81bdba518f82d224c048efd9aae2a045ad9abc74f2b18cde4322\",\"6f806a80720b0603d2ad3b6dfecc3801f42a2ea402789d8e2a77a6826b50303a\"]}",
   "apiVersion":"1.3.0",
   "timestamp":"2021-04-30T08:06:13.4129624Z",
   "minerId":"030d1fe5c1b560efe196ba40540ce9017c20daa9504c4c4cec6184fc702d9f274e",
   "blockHash":"2ad8af91739e9dc41ea155a9ab4b14ab88fe2a0934f14420139867babf5953c4",
   "blockHeight":105,
   "callbackTxId":"e7b3eefab33072e62283255f193ef5d22f26bbcfc0a80688fe2cc5178a32dda6",
   "callbackReason":"merkleProof"
}
*/

// Callback is the body contents posted to the provided callback url from Merchant API
type Callback struct {
	CallbackPayload string `json:"callbackPayload"`
	APIVersion      string `json:"apiVersion"`
	Timestamp       string `json:"timestamp"`
	MinerID         string `json:"minerId"`
	BlockHash       string `json:"blockHash"`
	BlockHeight     uint64 `json:"blockHeight"`
	CallbackTxID    string `json:"callbackTxId"`
	CallbackReason  string `json:"callbackReason"`
}

// ConflictedWith contains the information about the transactions that conflict
// with the transaction submitted to mAPI. A conflict could arise if multiple
// transactions attempt to spend the same UTXO (double spend).
type ConflictedWith struct {
	TxID string `json:"txid"`
	Size int    `json:"size"`
	Hex  string `json:"hex"`
}

// SubmitTransaction will fire a Merchant API request to submit a given transaction
//
// This endpoint is used to send a raw transaction to a miner for inclusion in the next block
// that the miner creates. It returns a JSONEnvelope with a payload that contains the response to the
// transaction submission. The purpose of the envelope is to ensure strict consistency in the
// message content for the purpose of signing responses.
//
// Specs: https://github.com/bitcoin-sv-specs/brfc-merchantapi/tree/v1.2-beta#Submit-transaction
func (c *Client) SubmitTransaction(ctx context.Context, miner *Miner, tx *Transaction) (*SubmitTransactionResponse, error) {

	// Make sure we have a valid miner
	if miner == nil {
		return nil, errors.New("miner was nil")
	}

	// Make the HTTP request
	result := submitTransaction(ctx, c, miner, tx)
	if result.Response.Error != nil {
		return nil, result.Response.Error
	}

	// Parse the response
	response, err := result.parseSubmission()
	if err != nil {
		return nil, err
	}

	// Valid query?
	if response.Results == nil || len(response.Results.ReturnResult) == 0 {
		return nil, errors.New("failed getting submission response from: " + miner.Name)
	}

	// Return the fully parsed response
	return &response, nil
}

// parseSubmission will convert the HTTP response into a struct and also unmarshal the payload JSON data
func (i *internalResult) parseSubmission() (response SubmitTransactionResponse, err error) {

	// Process the initial response payload
	if err = response.process(i.Miner, i.Response.BodyContents); err != nil {
		return
	}

	// If we have a valid payload
	if len(response.Payload) > 0 {
		err = json.Unmarshal([]byte(response.Payload), &response.Results)
	}
	return
}

// submitTransaction will fire the HTTP request to submit a transaction
func submitTransaction(ctx context.Context, client *Client, miner *Miner, tx *Transaction) (result *internalResult) {
	result = &internalResult{Miner: miner}
	data, _ := json.Marshal(tx) // Ignoring error - if it fails, the submission would also fail
	result.Response = httpRequest(ctx, client, &httpPayload{
		Method: http.MethodPost,
		URL:    miner.URL + routeSubmitTx,
		Token:  miner.Token,
		Data:   data,
	})
	return
}
