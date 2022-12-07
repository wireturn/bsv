package bc

// MapiCallback is the body contents posted to the provided callback url from Merchant API
type MapiCallback struct {
	CallbackPayload string `json:"callbackPayload"`
	APIVersion      string `json:"apiVersion"`
	Timestamp       string `json:"timestamp"`
	MinerID         string `json:"minerId"`
	BlockHash       string `json:"blockHash"`
	BlockHeight     uint64 `json:"blockHeight"`
	CallbackTxID    string `json:"callbackTxId"`
	CallbackReason  string `json:"callbackReason"`
}
