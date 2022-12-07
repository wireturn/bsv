package bc

import (
	"encoding/hex"
	"errors"

	"github.com/libsv/go-bt"
)

/*
Field 													Purpose 									 														Size (Bytes)
----------------------------------------------------------------------------------------------------
block_header 					Block Header				 																									80
txn_count 						Total number of txs in this block, including the coinbase tx 	 				VarInt
txns 									Every tx in this block, one after another, in raw tx format 	 				-
*/

// A Block in the Bitcoin blockchain.
type Block struct {
	BlockHeader *BlockHeader
	Txs         []*bt.Tx
}

// String returns the Block Header encoded as hex string.
func (b *Block) String() string {
	return hex.EncodeToString(b.Bytes())
}

// Bytes will decode a bitcoin block struct into a byte slice.
//
// See https://btcinformation.org/en/developer-reference#serialized-blocks
func (b *Block) Bytes() []byte {
	bytes := []byte{}

	bytes = append(bytes, b.BlockHeader.Bytes()...)

	txCount := uint64(len(b.Txs))
	bytes = append(bytes, bt.VarInt(txCount)...)

	for _, tx := range b.Txs {
		bytes = append(bytes, tx.ToBytes()...)
	}

	return bytes
}

// NewBlockFromStr will encode a block header hash
// into the bitcoin block header structure.
//
// See https://btcinformation.org/en/developer-reference#serialized-blocks
func NewBlockFromStr(blockStr string) (*Block, error) {
	blockBytes, err := hex.DecodeString(blockStr)
	if err != nil {
		return nil, err
	}

	return NewBlockFromBytes(blockBytes)
}

// NewBlockFromBytes will encode a block header byte slice
// into the bitcoin block header structure.
//
// See https://btcinformation.org/en/developer-reference#serialized-blocks
func NewBlockFromBytes(b []byte) (*Block, error) {
	if len(b) == 0 {
		return nil, errors.New("block cannot be empty")
	}

	var offset int
	bh, err := NewBlockHeaderFromBytes(b[:80])
	if err != nil {
		return nil, err
	}
	offset += 80

	txCount, size := bt.DecodeVarInt(b[offset:])
	offset += size

	var txs []*bt.Tx
	for i := 0; i < int(txCount); i++ {
		tx, size, err := bt.NewTxFromStream(b[offset:])
		if err != nil {
			return nil, err
		}
		txs = append(txs, tx)
		offset += size
	}

	return &Block{
		BlockHeader: bh,
		Txs:         txs,
	}, nil
}
