package spv

import (
	"github.com/libsv/go-bc"
)

// Envelope is a struct which contains all information needed for a transaction to be verified.
//
// spec at https://tsc.bitcoinassociation.net/standards/spv-envelope/
type Envelope struct {
	TxID          string               `json:"txid,omitempty"`
	RawTx         string               `json:"rawTx,omitempty"`
	Proof         *bc.MerkleProof      `json:"proof,omitempty"`
	MapiResponses []bc.MapiCallback    `json:"mapiResponses,omitempty"`
	Parents       map[string]*Envelope `json:"parents,omitempty"`
}

// IsAnchored returns true if the envelope is the anchor tx.
func (e *Envelope) IsAnchored() bool {
	return e.Proof != nil
}
