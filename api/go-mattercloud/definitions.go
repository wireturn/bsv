package mattercloud

// NetworkType is used internally to represent the possible values
// for network in queries to be submitted: {"main", "test", "stn"}
type NetworkType string

const (
	// apiKeyField is the field for the api key in all requests
	apiKeyField string = "api_key"

	// NetworkMain is for main-net
	NetworkMain NetworkType = "main"

	// NetworkTest is for test-net
	NetworkTest NetworkType = "test"

	//NetworkStn is for the stn-net
	NetworkStn NetworkType = "stn"
)

// APIInternalError is for internal server errors (most requests)
type APIInternalError struct {
	Errors       []string `json:"errors,omitempty"`
	ErrorMessage string   `json:"message,omitempty"`
	ErrorName    string   `json:"name,omitempty"`
}

// AddressList is the list of addresses for batch
type AddressList struct {
	Addrs string `json:"addrs"`
}

// TransactionList is the list of tx ids for batch
type TransactionList struct {
	TxIDs string `json:"txids"`
}

// Balance is the response from the get balance request
type Balance struct {
	Address     string `json:"address"`
	Confirmed   int64  `json:"confirmed"`
	Unconfirmed int64  `json:"unconfirmed"`
}

// UnspentTransaction is a standard UTXO response
type UnspentTransaction struct {
	Address       string  `json:"address"`
	Amount        float64 `json:"amount"`
	Confirmations int64   `json:"confirmations"`
	Height        int64   `json:"height"`
	OutputIndex   int64   `json:"outputIndex"`
	Satoshis      int64   `json:"satoshis"`
	Script        string  `json:"script"`
	ScriptPubKey  string  `json:"scriptPubKey"`
	TxID          string  `json:"txid"`
	Value         int64   `json:"value"`
	Vout          int     `json:"vout"`
}

// History is the response from address history
type History struct {
	From    int            `json:"from"`
	Results []*HistoryItem `json:"results"`
	To      int            `json:"to"`
}

// HistoryItem is the individual history item
type HistoryItem struct {
	TxID   string `json:"txid"`
	Height int64  `json:"height"`
}

// Transaction is returned in the GetTransactionsResponse
type Transaction struct {
	APIInternalError
	BlockHash     string     `json:"blockhash"`
	BlockHeight   int64      `json:"blockheight"`
	BlockTime     int64      `json:"blocktime"`
	Confirmations int64      `json:"confirmations"`
	Fees          float64    `json:"fees"`
	Hash          string     `json:"hash"`
	LockTime      int64      `json:"locktime"`
	RawTx         string     `json:"rawtx"`
	Size          int64      `json:"size"`
	Time          int64      `json:"time"`
	TxID          string     `json:"txid"`
	ValueIn       float64    `json:"valueIn"`
	ValueOut      float64    `json:"valueOut"`
	Version       int        `json:"version"`
	Vin           []VinType  `json:"vin"`
	Vout          []VoutType `json:"vout"`
}

// VinType is the vin data
type VinType struct {
	Address       string        `json:"address"`
	AddressAddr   string        `json:"addr"`
	N             int           `json:"n"`
	ScriptSig     ScriptSigType `json:"scriptSig"`
	Sequence      int64         `json:"sequence"`
	TxID          string        `json:"txid"`
	Value         float64       `json:"value"`
	ValueSatoshis int64         `json:"valueSat"`
	Vout          int           `json:"vout"`
}

// VoutType is the vout data
type VoutType struct {
	N             int              `json:"n"`
	ScriptPubKey  ScriptPubKeyType `json:"scriptPubKey"`
	SpentHeight   int64            `json:"spentHeight"`
	SpentIndex    int64            `json:"spentIndex"`
	SpentTxID     string           `json:"spentTxId"`
	Value         float64          `json:"value"`
	ValueSatoshis int64            `json:"valueSat"`
}

// ScriptPubKeyType is the script pubkey data
type ScriptPubKeyType struct {
	Addresses          []string `json:"addresses"`
	Asm                string   `json:"asm"`
	Hex                string   `json:"hex"`
	RequiredSignatures int      `json:"reqSigs"`
	Type               string   `json:"type"`
}

// ScriptSigType is the script signature data
type ScriptSigType struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}

// BroadcastResponse is the response for the broadcast
type BroadcastResponse struct {
	TxID string `json:"txid"`
}
