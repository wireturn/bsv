package spv_test

import (
	"context"
	"testing"

	"github.com/libsv/go-bc"
	"github.com/libsv/go-bc/spv"
	"github.com/stretchr/testify/assert"
)

type mockBlockHeaderChain struct{}

func (bhc *mockBlockHeaderChain) BlockHeader(ctx context.Context, blockHash string) (blockHeader *bc.BlockHeader, err error) {
	return bc.NewBlockHeaderFromStr("000000208e33a53195acad0ab42ddbdbe3e4d9ca081332e5b01a62e340dbd8167d1a787b702f61bb913ac2063e0f2aed6d933d3386234da5c8eb9e30e498efd25fb7cb96fff12c60ffff7f2001000000")
}

func TestVerifyMerkleProof(t *testing.T) {
	t.Parallel()

	proofJSON := &bc.MerkleProof{
		Index:  12,
		TxOrID: "ffeff11c25cde7c06d407490d81ef4d0db64aad6ab3d14393530701561a465ef",
		Target: "75edb0a69eb195cdd81e310553aa4d25e18450e08f168532a2c2e9cf447bf169",
		Nodes: []string{
			"b9ef07a62553ef8b0898a79c291b92c60f7932260888bde0dab2dd2610d8668e",
			"0fc1c12fb1b57b38140442927fbadb3d1e5a5039a5d6db355ea25486374f104d",
			"60b0e75dd5b8d48f2d069229f20399e07766dd651ceeed55ee3c040aa2812547",
			"c0d8dbda46366c2050b430a05508a3d96dc0ed55aea685bb3d9a993f8b97cc6f",
			"391e62b3419d8a943f7dbc7bddc90e30ec724c033000dc0c8872253c27b03a42",
		},
	}
	hcm := &mockBlockHeaderChain{}

	v, _ := spv.NewMerkleProofVerifier(hcm)

	t.Run("JSON", func(t *testing.T) {
		valid, isLastInTree, err := v.VerifyMerkleProofJSON(context.Background(), proofJSON)

		assert.NoError(t, err)
		assert.False(t, isLastInTree)
		assert.True(t, valid)
	})

	t.Run("Bytes", func(t *testing.T) {
		proof, _ := proofJSON.ToBytes()
		valid, isLastInTree, err := v.VerifyMerkleProof(context.Background(), proof)

		assert.NoError(t, err)
		assert.False(t, isLastInTree)
		assert.True(t, valid)
	})
}
