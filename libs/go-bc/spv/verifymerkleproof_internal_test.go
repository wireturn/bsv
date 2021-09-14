package spv

import (
	"testing"

	"github.com/libsv/go-bc"
	"github.com/stretchr/testify/assert"
)

func TestParseBinaryMerkleProof(t *testing.T) {
	var proofJSON = &bc.MerkleProof{
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

	proof, _ := proofJSON.ToBytes()

	mpb, err := parseBinaryMerkleProof(proof)

	assert.NoError(t, err)
	assert.Equal(t, proofJSON.Index, mpb.index)
	assert.Equal(t, proofJSON.TxOrID, mpb.txOrID)
	assert.Equal(t, proofJSON.Target, mpb.target)
	assert.Equal(t, proofJSON.Nodes, mpb.nodes)
}

func TestVerifyProof(t *testing.T) {
	c := "ffeff11c25cde7c06d407490d81ef4d0db64aad6ab3d14393530701561a465ef"
	merkleRoot := "96cbb75fd2ef98e4309eebc8a54d2386333d936ded2a0f3e06c23a91bb612f70"
	index := 12
	nodes := []string{"b9ef07a62553ef8b0898a79c291b92c60f7932260888bde0dab2dd2610d8668e", "0fc1c12fb1b57b38140442927fbadb3d1e5a5039a5d6db355ea25486374f104d", "60b0e75dd5b8d48f2d069229f20399e07766dd651ceeed55ee3c040aa2812547", "c0d8dbda46366c2050b430a05508a3d96dc0ed55aea685bb3d9a993f8b97cc6f", "391e62b3419d8a943f7dbc7bddc90e30ec724c033000dc0c8872253c27b03a42"}

	valid, isLastInTree, err := verifyProof(c, merkleRoot, uint64(index), nodes)

	assert.NoError(t, err)
	assert.False(t, isLastInTree)
	assert.True(t, valid)
}
