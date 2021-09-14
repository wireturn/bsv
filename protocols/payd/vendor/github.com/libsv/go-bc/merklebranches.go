package bc

import (
	"encoding/hex"
	"fmt"

	"github.com/libsv/go-bk/crypto"
	"github.com/libsv/go-bt/v2"
)

func getHashes(txHashes []string) []string {
	hashes := make([]string, 0, len(txHashes))

	for i, tx := range txHashes {
		hashes[i] = (ReverseHexString(tx))
	}

	return hashes
}

// GetMerkleBranches comment TODO:
func GetMerkleBranches(template []string) []string {
	hashes := getHashes(template)

	var branches []string
	var walkBranch func(hashes []string) []string

	walkBranch = func(hashes []string) []string {
		var results []string

		tot := len(hashes)

		if len(hashes) < 2 {
			return make([]string, 0)
		}

		branches = append(branches, hashes[1])

		for i := 0; i < tot; i += 2 {
			var a, _ = hex.DecodeString(hashes[i])
			var b []byte
			if (i + 1) < tot {
				b, _ = hex.DecodeString(hashes[i+1])
			} else {
				b = a
			}

			concat := append(a, b...)
			hash := crypto.Sha256d(concat)
			results = append(results, hex.EncodeToString(hash[:]))
		}

		return walkBranch(results)
	}

	walkBranch(hashes)

	return branches
}

// MerkleRootFromBranches returns a Merkle root given a transaction hash (txid), the index in
// which it is positioned in the Merkle tree, and the branches needed along the way (Merkle path).
func MerkleRootFromBranches(txHash string, txIndex int, branches []string) (string, error) {
	hash, err := hex.DecodeString(txHash)
	if err != nil {
		return "", err
	}

	hash = bt.ReverseBytes(hash)

	for _, b := range branches {
		h, err := hex.DecodeString(b)
		if err != nil {
			return "", err
		}

		h = bt.ReverseBytes(h)

		if txIndex&1 > 0 {
			hash = crypto.Sha256d(append(h, hash...))
		} else {
			hash = crypto.Sha256d(append(hash, h...))
		}

		txIndex >>= 1
	}

	if txIndex > 0 {
		return "", fmt.Errorf("index %d out of range for proof of length %d", txIndex, len(branches))
	}

	return hex.EncodeToString(bt.ReverseBytes(hash)), nil

}
