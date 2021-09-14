package whatsonchain

// NetworkType is used internally to represent the possible values
// for network in queries to be submitted: {"main", "test", "stn"}
type NetworkType string

const (

	// NetworkMain is for main-net
	NetworkMain NetworkType = "main"

	// NetworkTest is for test-net
	NetworkTest NetworkType = "test"

	// NetworkStn is for the stn-net
	NetworkStn NetworkType = "stn"

	// MaxTransactionsUTXO is the max allowed in the request
	MaxTransactionsUTXO int = 20

	// MaxBroadcastTransactions is the max transactions for Bulk Broadcast
	MaxBroadcastTransactions = 100

	// MaxSingleTransactionSize is the max single TX size for Bulk Broadcast
	MaxSingleTransactionSize = 102400

	// MaxCombinedTransactionSize is the max of all transactions combined
	MaxCombinedTransactionSize = 1e+7

	// MaxAddressesForLookup is the max allowed in the request for Bulk requests
	MaxAddressesForLookup int = 20

	// MaxScriptsForLookup is the max allowed in the request for Bulk requests
	MaxScriptsForLookup int = 20

	// MaxRequestsPerSecond is the max requests allowed per second
	MaxRequestsPerSecond int = 3
)

// AddressInfo is the address info for a returned address request
type AddressInfo struct {
	Address      string `json:"address"`
	IsMine       bool   `json:"ismine"`
	IsScript     bool   `json:"isscript"`
	IsValid      bool   `json:"isvalid"`
	IsWatchOnly  bool   `json:"iswatchonly"`
	ScriptPubKey string `json:"scriptPubKey"`
}

// AddressBalance is the address balance (unconfirmed and confirmed)
type AddressBalance struct {
	Confirmed   int64 `json:"confirmed"`
	Unconfirmed int64 `json:"unconfirmed"`
}

// AddressBalanceRecord is the result from Bulk Balance request
type AddressBalanceRecord struct {
	Address string          `json:"address"`
	Error   string          `json:"error"`
	Balance *AddressBalance `json:"balance"`
}

// AddressList is used to create a Bulk Balance request
type AddressList struct {
	Addresses []string `json:"addresses"`
}

// AddressBalances is the response from Bulk Balance request
type AddressBalances []*AddressBalanceRecord

// AddressHistory is the history of transactions for an address
type AddressHistory []*HistoryRecord

// BlockInfo is the response info about a returned block
type BlockInfo struct {
	Bits              string         `json:"bits"`
	ChainWork         string         `json:"chainwork"`
	CoinbaseTx        CoinbaseTxInfo `json:"coinbaseTx"`
	Confirmations     int64          `json:"confirmations"`
	Difficulty        float64        `json:"difficulty"`
	Hash              string         `json:"hash"`
	Height            int64          `json:"height"`
	MedianTime        int64          `json:"mediantime"`
	MerkleRoot        string         `json:"merkleroot"`
	Miner             string         `json:"Bmgpool"`
	NextBlockHash     string         `json:"nextblockhash"`
	Nonce             int64          `json:"nonce"`
	Pages             Page           `json:"pages"`
	PreviousBlockHash string         `json:"previousblockhash"`
	Size              int64          `json:"size"`
	Time              int64          `json:"time"`
	TotalFees         float64        `json:"totalFees"`
	Tx                []string       `json:"tx"`
	TxCount           int64          `json:"txcount"`
	Version           int64          `json:"version"`
	VersionHex        string         `json:"versionHex"`
}

// BlockPagesInfo is the response from the page request
type BlockPagesInfo []string

// BulkBroadcastResponse is the response from a bulk broadcast request
type BulkBroadcastResponse struct {
	Feedback  bool   `json:"feedback"`
	StatusURL string `json:"statusUrl"`
}

// BulkUnspentResponse is the response from Bulk Unspent transactions
type BulkUnspentResponse []*BulkResponseRecord

// BulkResponseRecord is the record in the results for Bulk Unspent transactions
type BulkResponseRecord struct {
	Address string           `json:"address"`
	Error   string           `json:"error"`
	Utxos   []*HistoryRecord `json:"unspent"`
}

// BulkScriptUnspentResponse is the response from Bulk Unspent transactions
type BulkScriptUnspentResponse []*BulkScriptResponseRecord

// BulkScriptResponseRecord is the record in the results for Bulk Unspent transactions
type BulkScriptResponseRecord struct {
	Script string           `json:"script"`
	Error  string           `json:"error"`
	Utxos  []*HistoryRecord `json:"unspent"`
}

// ChainInfo is the structure response from getting info about the chain
type ChainInfo struct {
	BestBlockHash        string  `json:"bestblockhash"`
	Blocks               int64   `json:"blocks"`
	Chain                string  `json:"chain"`
	ChainWork            string  `json:"chainwork"`
	Difficulty           float64 `json:"difficulty"`
	Headers              int64   `json:"headers"`
	MedianTime           int64   `json:"mediantime"`
	Pruned               bool    `json:"pruned"`
	VerificationProgress float64 `json:"verificationprogress"`
}

// CirculatingSupply is the structure response
type CirculatingSupply float64

// CoinbaseTxInfo is the coinbase tx info inside the BlockInfo
type CoinbaseTxInfo struct {
	BlockHash     string     `json:"blockhash"`
	BlockTime     int64      `json:"blocktime"`
	Confirmations int64      `json:"confirmations"`
	Hash          string     `json:"hash"`
	Hex           string     `json:"hex"`
	LockTime      int64      `json:"locktime"`
	Size          int64      `json:"size"`
	Time          int64      `json:"time"`
	TxID          string     `json:"txid"`
	Version       int64      `json:"version"`
	Vin           []VinInfo  `json:"vin"`
	Vout          []VoutInfo `json:"vout"`
}

// ExchangeRate is the response from getting the current exchange rate
type ExchangeRate struct {
	Currency string `json:"currency"`
	Rate     string `json:"rate"`
}

// FeeQuote is the structure response for a fee in a quote
type FeeQuote struct {
	FeeType   string `json:"feeType"`
	MiningFee *Fee   `json:"miningFee"`
	RelayFee  *Fee   `json:"relayFee"`
}

// Fee is the actual fee (satoshis per byte)
type Fee struct {
	Bytes    int `json:"bytes"`
	Satoshis int `json:"satoshis"`
}

// FeeQuotes is the structure response from getting quotes from Merchant API
type FeeQuotes struct {
	Quotes []*QuoteProvider `json:"quotes"`
}

// HistoryRecord is an internal record of AddressHistory
type HistoryRecord struct {
	Height int64   `json:"height"`
	Info   *TxInfo `json:"info,omitempty"` // Custom for our wrapper
	TxHash string  `json:"tx_hash"`
	TxPos  int64   `json:"tx_pos"`
	Value  int64   `json:"value"`
}

// MempoolInfo is the response for the get mempool info request
type MempoolInfo struct {
	Bytes         int64 `json:"bytes"`
	MaxMempool    int64 `json:"maxmempool"`
	MempoolMinFee int64 `json:"mempoolminfee"`
	Size          int64 `json:"size"`
	Usage         int64 `json:"usage"`
}

// MerkleResults is the results from the proof request
type MerkleResults []*MerkleInfo

// MerkleInfo is the response for the get merkle request
type MerkleInfo struct {
	BlockHash  string          `json:"blockHash"`
	Branches   []*MerkleBranch `json:"branches"`
	Hash       string          `json:"hash"`
	MerkleRoot string          `json:"merkleRoot"`
}

// MerkleBranch is a merkle branch
type MerkleBranch struct {
	Hash string `json:"hash"`
	Pos  string `json:"pos"`
}

// MerchantResponse is the response from a tx submission
type MerchantResponse struct {
	APIVersion                string `json:"apiVersion"`
	CurrentHighestBlockHash   string `json:"currentHighestBlockHash"`
	CurrentHighestBlockHeight int64  `json:"currentHighestBlockHeight"`
	MinerID                   string `json:"minerId"`
	ResultDescription         string `json:"resultDescription"`
	ReturnResult              string `json:"returnResult"`
	Timestamp                 string `json:"timestamp"`
	TxID                      string `json:"txid"`
	TxSecondMempoolExpiry     int    `json:"txSecondMempoolExpiry"`
}

// MerchantError is the error response from a bad tx submission
type MerchantError struct {
	Code   int    `json:"code"`
	Error  string `json:"error"`
	Status int    `json:"status"`
}

// MerchantStatus is the response from a status request
type MerchantStatus struct {
	APIVersion            string `json:"apiVersion"`
	BlockHash             string `json:"blockHash"`
	BlockHeight           int64  `json:"blockHeight"`
	Confirmations         int64  `json:"confirmations"`
	MinerID               string `json:"minerId"`
	ResultDescription     string `json:"resultDescription"`
	ReturnResult          string `json:"returnResult"`
	Timestamp             string `json:"timestamp"`
	TxSecondMempoolExpiry int    `json:"txSecondMempoolExpiry"`
}

// Page is used as a sub-type for BlockInfo
type Page struct {
	Size int64    `json:"size"`
	URI  []string `json:"uri"`
}

// ScriptsList is used to create a Bulk UTXO request
type ScriptsList struct {
	Scripts []string `json:"scripts"`
}

// ScriptList is the list of script history records
type ScriptList []*ScriptRecord

// ScriptRecord is the script history record
type ScriptRecord struct {
	Height int64  `json:"height"`
	TxHash string `json:"tx_hash"`
	TxPos  int64  `json:"tx_pos"`
	Value  int64  `json:"value"`
}

// ScriptSigInfo is the scriptSig info inside the VinInfo
type ScriptSigInfo struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}

// SearchResults is the response from searching for explorer links
type SearchResults struct {
	Results []*SearchResult `json:"results"`
}

// SearchResult is the actual result for the search (included in SearchResults)
type SearchResult struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

// ScriptPubKeyInfo is the scriptPubKey info inside the VoutInfo
type ScriptPubKeyInfo struct {
	Addresses   []string `json:"addresses"`
	Asm         string   `json:"asm"`
	Hex         string   `json:"hex"`
	IsTruncated bool     `json:"isTruncated"`
	OpReturn    string   `json:"-"` // todo: support this (can be an object of key/vals based on the op return data)
	ReqSigs     int64    `json:"reqSigs"`
	Type        string   `json:"type"`
}

// StatusResponse is the response from requesting a status update
type StatusResponse struct {
	Payload      string          `json:"payload"`
	ProviderID   string          `json:"providerId"`
	ProviderName string          `json:"providerName"`
	PublicKey    string          `json:"publicKey"`
	Signature    string          `json:"signature"`
	Status       *MerchantStatus `json:"status"`
}

// QuoteProvider is the structure response for a quote provider (which has quotes)
type QuoteProvider struct {
	Payload         string `json:"payload"`
	ProviderID      string `json:"providerId"`
	ProviderName    string `json:"providerName"`
	PublicKey       string `json:"publicKey"`
	Quote           *Quote `json:"quote"`
	Signature       string `json:"signature"`
	TxStatusURL     string `json:"txStatusUrl"`
	TxSubmissionURL string `json:"txSubmissionUrl"`
}

// Quote is the structure response for a quote
type Quote struct {
	APIVersion                string      `json:"apiVersion"`
	CurrentHighestBlockHash   string      `json:"currentHighestBlockHash"`
	CurrentHighestBlockHeight int64       `json:"currentHighestBlockHeight"`
	ExpiryTime                string      `json:"expiryTime"`
	Fees                      []*FeeQuote `json:"fees"`
	MinerID                   string      `json:"minerId"`
	MinerReputation           interface{} `json:"minerReputation"`
	Timestamp                 string      `json:"timestamp"`
}

// SubmissionResponse is the response from submitting a tx via Merchant API
type SubmissionResponse struct {
	Error        *MerchantError    `json:"error"`
	Payload      string            `json:"payload"`
	ProviderID   string            `json:"providerId"`
	ProviderName string            `json:"providerName"`
	PublicKey    string            `json:"publicKey"`
	Response     *MerchantResponse `json:"response"`
	Signature    string            `json:"signature"`
}

// TxInfo is the response info about a returned tx
type TxInfo struct {
	BlockHash     string     `json:"blockhash"`
	BlockTime     int64      `json:"blocktime"`
	Confirmations int64      `json:"confirmations"`
	Hash          string     `json:"hash"`
	Hex           string     `json:"hex"`
	LockTime      int64      `json:"locktime"`
	Size          int64      `json:"size"`
	Time          int64      `json:"time"`
	TxID          string     `json:"txid"`
	Version       int64      `json:"version"`
	Vin           []VinInfo  `json:"vin"`
	Vout          []VoutInfo `json:"vout"`
}

// TxList is the list of tx info structs returned from the /txs post response
type TxList []*TxInfo

// TxHashes is the list of tx hashes for the post request
type TxHashes struct {
	TxIDs []string `json:"txids"`
}

// VinInfo is the vin info inside the CoinbaseTxInfo
type VinInfo struct {
	Coinbase  string        `json:"coinbase"`
	ScriptSig ScriptSigInfo `json:"scriptSig"`
	Sequence  int64         `json:"sequence"`
	TxID      string        `json:"txid"`
	Vout      int64         `json:"vout"`
}

// VoutInfo is the vout info inside the CoinbaseTxInfo
type VoutInfo struct {
	N            int64            `json:"n"`
	ScriptPubKey ScriptPubKeyInfo `json:"scriptPubKey"`
	Value        float64          `json:"value"`
}
