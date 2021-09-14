package gopayd

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/libsv/go-bc"
	"github.com/libsv/go-bk/envelope"
	"github.com/libsv/go-bt/v2"
	"github.com/pkg/errors"
	validator "github.com/theflyingcodr/govalidator"
)

var reTargetType = regexp.MustCompile(`^(header|hash|merkleRoot)$`)

// ProofCreateArgs are used to create a proof.
type ProofCreateArgs struct {
	// TxID will be used to validate the proof envelope.
	TxID string `json:"txId" param:"txid"`
}

// ProofWrapper represents a mapi callback payload for a merkleproof.
// mAPI returns proofs in a JSONEnvelope with a payload. This represents the
// Payload format which contains a parent object with tx meta and a nested object
// which is the TSC format merkleProof.
type ProofWrapper struct {
	CallbackPayload *bc.MerkleProof `json:"callbackPayload"`
	BlockHash       string          `json:"blockHash"`
	BlockHeight     uint32          `json:"blockHeight"`
	CallbackTxID    string          `json:"callbackTxID"`
	CallbackReason  string          `json:"callbackReason"`
}

// Validate will ensure the ProofWrapper is valid.
func (p ProofWrapper) Validate(args ProofCreateArgs) error {
	vl := validator.New().Validate("blockhash",
		validator.NotEmpty(p.BlockHash)).
		Validate("callbackReason", func() error {
			if strings.ToLower(p.CallbackReason) != "merkleproof" {
				return errors.New("invalid callback received, should be of type merkleProof")
			}
			return nil
		}).Validate("callbackTxID", func() error {
		if args.TxID != p.CallbackTxID {
			return fmt.Errorf("proof txid does not match expected txid %s", args.TxID)
		}
		return nil
	}).Validate("callbackPayload", validator.NotEmpty(p.CallbackPayload))
	if p.CallbackPayload == nil {
		return vl.Err()
	}
	vl = vl.Validate("callbackPayload.targetType", validator.MatchString(p.CallbackPayload.TargetType, reTargetType)).
		Validate("callbackPayload.target", validator.NotEmpty(p.CallbackPayload.Target)).
		Validate("callbackPayload.proofType", func() error {
			if p.CallbackPayload.ProofType == "" {
				return nil
			}
			if strings.ToLower(p.CallbackPayload.ProofType) != "branch" && strings.ToLower(p.CallbackPayload.ProofType) != "tree" {
				return errors.New("only branch or tree are allowed as proofType")
			}
			return nil
		}).Validate("callbackPayload.txOrId", func() error {
		if p.CallbackPayload.TxOrID == args.TxID {
			return nil
		}
		if len(p.CallbackPayload.TxOrID) == 64 {
			return fmt.Errorf("txId provided in callbackPayload doesn't match expected txID %s", args.TxID)
		}
		tx, err := bt.NewTxFromString(p.CallbackPayload.TxOrID)
		if err != nil {
			return errors.Wrap(err, "failed to parse txhex")
		}
		if tx.TxID() != args.TxID {
			return fmt.Errorf("tx provided in callbackPayload doesn't match expected txID %s", args.TxID)
		}
		return nil
	})
	return vl.Err()
}

// ProofsService enforces business rules and validation when handling merkle proofs.
type ProofsService interface {
	// Create will store a JSONEnvelope that contains a merkleproof. The envelope should
	// be validated to not be tampered with and the Envelope should be opened to check the payload
	// is indeed a MerkleProof.
	Create(ctx context.Context, args ProofCreateArgs, req envelope.JSONEnvelope) error
}

// ProofsWriter is used to persist a proof to a data store.
type ProofsWriter interface {
	// ProofCreate can be used to persist a merkle proof in TSC format.
	ProofCreate(ctx context.Context, req ProofWrapper) error
}
