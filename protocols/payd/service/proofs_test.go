package service

import (
	"context"
	"errors"
	"testing"

	"github.com/libsv/go-bc"
	"github.com/libsv/go-bk/envelope"
	"github.com/stretchr/testify/assert"

	gopayd "github.com/libsv/payd"
	"github.com/libsv/payd/mocks"
)

func Test_Proofs_create(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		args           gopayd.ProofCreateArgs
		req            envelope.JSONEnvelope
		proofsCreateFn func(ctx context.Context, req gopayd.ProofWrapper) error
		err            error
	}{
		"successful run should return no errors": {
			args: gopayd.ProofCreateArgs{
				TxID: "2f8d0ac044aa2fd8fc7675809f5d17acac4e9bf63dd0ea4eb58f43b66ccc70ca",
			},
			req: func() envelope.JSONEnvelope {
				e, err := envelope.NewJSONEnvelope(gopayd.ProofWrapper{
					CallbackPayload: &bc.MerkleProof{
						Index:      0,
						TxOrID:     "2f8d0ac044aa2fd8fc7675809f5d17acac4e9bf63dd0ea4eb58f43b66ccc70ca",
						Target:     "abc123",
						Nodes:      nil,
						TargetType: "header",
						ProofType:  "",
						Composite:  false,
					},
					BlockHash:      "abc123",
					BlockHeight:    0,
					CallbackTxID:   "2f8d0ac044aa2fd8fc7675809f5d17acac4e9bf63dd0ea4eb58f43b66ccc70ca",
					CallbackReason: "merkleProof",
				})
				assert.NoError(t, err)
				assert.NotEmpty(t, e)
				return *e
			}(),
			proofsCreateFn: func(ctx context.Context, req gopayd.ProofWrapper) error {
				return nil
			},
			err: nil,
		}, "mismatch txid should return error": {
			args: gopayd.ProofCreateArgs{
				TxID: "2f8d0ac044aa2fd8fc7675809f5d17acac4e9bf63dd0ea4eb58f43b66ccc70cb",
			},
			req: func() envelope.JSONEnvelope {
				e, err := envelope.NewJSONEnvelope(gopayd.ProofWrapper{
					CallbackPayload: &bc.MerkleProof{
						Index:      0,
						TxOrID:     "2f8d0ac044aa2fd8fc7675809f5d17acac4e9bf63dd0ea4eb58f43b66ccc70ca",
						Target:     "abc123",
						Nodes:      nil,
						TargetType: "header",
						ProofType:  "",
						Composite:  false,
					},
					BlockHash:      "abc123",
					BlockHeight:    123,
					CallbackTxID:   "2f8d0ac044aa2fd8fc7675809f5d17acac4e9bf63dd0ea4eb58f43b66ccc70ca",
					CallbackReason: "merkleProof",
				})
				assert.NoError(t, err)
				assert.NotEmpty(t, e)
				return *e
			}(),
			proofsCreateFn: func(ctx context.Context, req gopayd.ProofWrapper) error {
				return nil
			},
			err: errors.New("[callbackPayload.txOrId: txId provided in callbackPayload doesn't match expected txID 2f8d0ac044aa2fd8fc7675809f5d17acac4e9bf63dd0ea4eb58f43b66ccc70cb], [callbackTxID: proof txid does not match expected txid 2f8d0ac044aa2fd8fc7675809f5d17acac4e9bf63dd0ea4eb58f43b66ccc70cb]"),
		}, "mismatch txid in proof only should return single error": {
			args: gopayd.ProofCreateArgs{
				TxID: "2f8d0ac044aa2fd8fc7675809f5d17acac4e9bf63dd0ea4eb58f43b66ccc70cb",
			},
			req: func() envelope.JSONEnvelope {
				e, err := envelope.NewJSONEnvelope(gopayd.ProofWrapper{
					CallbackPayload: &bc.MerkleProof{
						Index:      0,
						TxOrID:     "2f8d0ac044aa2fd8fc7675809f5d17acac4e9bf63dd0ea4eb58f43b66ccc70ca",
						Target:     "abc123",
						Nodes:      nil,
						TargetType: "header",
						ProofType:  "",
						Composite:  false,
					},
					BlockHash:      "abc123",
					BlockHeight:    123,
					CallbackTxID:   "2f8d0ac044aa2fd8fc7675809f5d17acac4e9bf63dd0ea4eb58f43b66ccc70cb",
					CallbackReason: "merkleProof",
				})
				assert.NoError(t, err)
				assert.NotEmpty(t, e)
				return *e
			}(),
			proofsCreateFn: func(ctx context.Context, req gopayd.ProofWrapper) error {
				return nil
			},
			err: errors.New("[callbackPayload.txOrId: txId provided in callbackPayload doesn't match expected txID 2f8d0ac044aa2fd8fc7675809f5d17acac4e9bf63dd0ea4eb58f43b66ccc70cb]"),
		}, "empty payload should return error": {
			args: gopayd.ProofCreateArgs{
				TxID: "2f8d0ac044aa2fd8fc7675809f5d17acac4e9bf63dd0ea4eb58f43b66ccc70cb",
			},
			req: func() envelope.JSONEnvelope {
				e, err := envelope.NewJSONEnvelope(gopayd.ProofWrapper{
					BlockHash:      "abc123",
					BlockHeight:    0,
					CallbackTxID:   "2f8d0ac044aa2fd8fc7675809f5d17acac4e9bf63dd0ea4eb58f43b66ccc70cb",
					CallbackReason: "merkleProof",
				})
				assert.NoError(t, err)
				assert.NotEmpty(t, e)
				return *e
			}(),
			proofsCreateFn: func(ctx context.Context, req gopayd.ProofWrapper) error {
				return nil
			},
			err: errors.New("[callbackPayload: value cannot be empty]"),
		}, "invalid envelope sig should return error": {
			args: gopayd.ProofCreateArgs{
				TxID: "2f8d0ac044aa2fd8fc7675809f5d17acac4e9bf63dd0ea4eb58f43b66ccc70cb",
			},
			req: func() envelope.JSONEnvelope {
				e, err := envelope.NewJSONEnvelope(gopayd.ProofWrapper{
					BlockHash:      "abc123",
					BlockHeight:    1,
					CallbackTxID:   "2f8d0ac044aa2fd8fc7675809f5d17acac4e9bf63dd0ea4eb58f43b66ccc70cb",
					CallbackReason: "merkleProof",
				})
				assert.NoError(t, err)
				assert.NotEmpty(t, e)
				e.Signature = func() *string { s := "aaaaaaaa"; return &s }()
				return *e
			}(),
			proofsCreateFn: func(ctx context.Context, req gopayd.ProofWrapper) error {
				return nil
			},
			err: errors.New("[jsonEnvelope: invalid merkleproof envelope: failed to parse json envelope signature malformed signature: too short]"),
		}, "invalid callback reason should error": {
			args: gopayd.ProofCreateArgs{
				TxID: "2f8d0ac044aa2fd8fc7675809f5d17acac4e9bf63dd0ea4eb58f43b66ccc70ca",
			},
			req: func() envelope.JSONEnvelope {
				e, err := envelope.NewJSONEnvelope(gopayd.ProofWrapper{
					CallbackPayload: &bc.MerkleProof{
						Index:      0,
						TxOrID:     "2f8d0ac044aa2fd8fc7675809f5d17acac4e9bf63dd0ea4eb58f43b66ccc70ca",
						Target:     "abc123",
						Nodes:      nil,
						TargetType: "header",
						ProofType:  "",
						Composite:  false,
					},
					BlockHash:      "abc123",
					BlockHeight:    2,
					CallbackTxID:   "2f8d0ac044aa2fd8fc7675809f5d17acac4e9bf63dd0ea4eb58f43b66ccc70ca",
					CallbackReason: "mine",
				})
				assert.NoError(t, err)
				assert.NotEmpty(t, e)
				return *e
			}(),
			proofsCreateFn: func(ctx context.Context, req gopayd.ProofWrapper) error {
				return nil
			},
			err: errors.New("[callbackReason: invalid callback received, should be of type merkleProof]"),
		}, "txhex with correct txid should validate with no error": {
			args: gopayd.ProofCreateArgs{
				TxID: "2f8d0ac044aa2fd8fc7675809f5d17acac4e9bf63dd0ea4eb58f43b66ccc70ca",
			},
			req: func() envelope.JSONEnvelope {
				e, err := envelope.NewJSONEnvelope(gopayd.ProofWrapper{
					CallbackPayload: &bc.MerkleProof{
						Index:      0,
						TxOrID:     "02000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0502a5000101ffffffff0100f9029500000000232103e59d8dfdd42499ceaf97362c64b997abaafc872efacc7a2011df98ee41fe84b7ac00000000",
						Target:     "abc123",
						Nodes:      nil,
						TargetType: "header",
						ProofType:  "",
						Composite:  false,
					},
					BlockHash:      "abc123",
					BlockHeight:    7,
					CallbackTxID:   "2f8d0ac044aa2fd8fc7675809f5d17acac4e9bf63dd0ea4eb58f43b66ccc70ca",
					CallbackReason: "merkleProof",
				})
				assert.NoError(t, err)
				assert.NotEmpty(t, e)
				return *e
			}(),
			proofsCreateFn: func(ctx context.Context, req gopayd.ProofWrapper) error {
				return nil
			},
		}, "invalid targetType should error": {
			args: gopayd.ProofCreateArgs{
				TxID: "2f8d0ac044aa2fd8fc7675809f5d17acac4e9bf63dd0ea4eb58f43b66ccc70ca",
			},
			req: func() envelope.JSONEnvelope {
				e, err := envelope.NewJSONEnvelope(gopayd.ProofWrapper{
					CallbackPayload: &bc.MerkleProof{
						Index:      0,
						TxOrID:     "2f8d0ac044aa2fd8fc7675809f5d17acac4e9bf63dd0ea4eb58f43b66ccc70ca",
						Target:     "abc123",
						Nodes:      nil,
						TargetType: "mine",
						ProofType:  "",
						Composite:  false,
					},
					BlockHash:      "abc123",
					BlockHeight:    0,
					CallbackTxID:   "2f8d0ac044aa2fd8fc7675809f5d17acac4e9bf63dd0ea4eb58f43b66ccc70ca",
					CallbackReason: "merkleProof",
				})
				assert.NoError(t, err)
				assert.NotEmpty(t, e)
				return *e
			}(),
			proofsCreateFn: func(ctx context.Context, req gopayd.ProofWrapper) error {
				return nil
			},
			err: errors.New("[callbackPayload.targetType: value mine failed to meet requirements]"),
		}, "error from proof create should be echoed back": {
			args: gopayd.ProofCreateArgs{
				TxID: "2f8d0ac044aa2fd8fc7675809f5d17acac4e9bf63dd0ea4eb58f43b66ccc70ca",
			},
			req: func() envelope.JSONEnvelope {
				e, err := envelope.NewJSONEnvelope(gopayd.ProofWrapper{
					CallbackPayload: &bc.MerkleProof{
						Index:      0,
						TxOrID:     "2f8d0ac044aa2fd8fc7675809f5d17acac4e9bf63dd0ea4eb58f43b66ccc70ca",
						Target:     "abc123",
						Nodes:      nil,
						TargetType: "header",
						ProofType:  "",
						Composite:  false,
					},
					BlockHash:      "abc123",
					BlockHeight:    0,
					CallbackTxID:   "2f8d0ac044aa2fd8fc7675809f5d17acac4e9bf63dd0ea4eb58f43b66ccc70ca",
					CallbackReason: "merkleProof",
				})
				assert.NoError(t, err)
				assert.NotEmpty(t, e)
				return *e
			}(),
			proofsCreateFn: func(ctx context.Context, req gopayd.ProofWrapper) error {
				return errors.New("I failed")
			},
			err: errors.New("failed to save proof: I failed"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockProofWrtr := &mocks.ProofsWriterMock{ProofCreateFunc: test.proofsCreateFn}
			err := NewProofsService(mockProofWrtr).Create(context.Background(), test.args, test.req)
			if test.err != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, test.err.Error())
				return
			}
			assert.NoError(t, err)
		})
	}
}
