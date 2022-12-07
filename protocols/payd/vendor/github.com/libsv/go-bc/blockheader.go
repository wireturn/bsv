package bc

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/libsv/go-bk/crypto"
	"github.com/libsv/go-bt/v2"
)

/*
Field 													Purpose 									 														Size (Bytes)
----------------------------------------------------------------------------------------------------
Version 							Block version number 																									4
hashPrevBlock 				256-bit hash of the previous block header 	 													32
hashMerkleRoot 				256-bit hash based on all of the transactions in the block 	 					32
Time 									Current block timestamp as seconds since 1970-01-01T00:00 UTC 				4
Bits 									Current target in compact format 	 																		4
Nonce 								32-bit number (starts at 0) 	 																				4
*/

// A BlockHeader in the Bitcoin blockchain.
type BlockHeader struct {
	Version        uint32
	Time           uint32
	Nonce          uint32
	HashPrevBlock  []byte
	HashMerkleRoot []byte
	Bits           []byte
}

// HashPrevBlockStr returns the Block Header encoded as hex string.
func (bh *BlockHeader) HashPrevBlockStr() string {
	return hex.EncodeToString(bh.HashPrevBlock)
}

// HashMerkleRootStr returns the Block Header encoded as hex string.
func (bh *BlockHeader) HashMerkleRootStr() string {
	return hex.EncodeToString(bh.HashMerkleRoot)
}

// BitsStr returns the Block Header encoded as hex string.
func (bh *BlockHeader) BitsStr() string {
	return hex.EncodeToString(bh.Bits)
}

// String returns the Block Header encoded as hex string.
func (bh *BlockHeader) String() string {
	return hex.EncodeToString(bh.Bytes())
}

// Bytes will decode a bitcoin block header struct
// into a byte slice.
//
// See https://en.bitcoin.it/wiki/Block_hashing_algorithm
func (bh *BlockHeader) Bytes() []byte {
	bytes := []byte{}
	bytes = append(bytes, UInt32ToBytes(bh.Version)...)
	bytes = append(bytes, bt.ReverseBytes(bh.HashPrevBlock)...)
	bytes = append(bytes, bt.ReverseBytes(bh.HashMerkleRoot)...)
	bytes = append(bytes, UInt32ToBytes(bh.Time)...)
	bytes = append(bytes, bt.ReverseBytes(bh.Bits)...)
	bytes = append(bytes, UInt32ToBytes(bh.Nonce)...)
	return bytes
}

// Valid checks whether a blockheader satisfies the proof-of-work claimed
// in Bits. Wwe check whether its Hash256 read as a little endian number
// is less than the Bits written in expanded form.
func (bh *BlockHeader) Valid() bool {
	target, err := ExpandTargetFromAsInt(hex.EncodeToString(bh.Bits))
	if err != nil {
		return false
	}

	digest := bt.ReverseBytes(crypto.Sha256d(bh.Bytes()))
	var bn *big.Int = big.NewInt(0)
	bn.SetBytes(digest)

	return bn.Cmp(target) < 0
}

// NewBlockHeaderFromStr will encode a block header hash
// into the bitcoin block header structure.
//
// See https://en.bitcoin.it/wiki/Block_hashing_algorithm
func NewBlockHeaderFromStr(headerStr string) (*BlockHeader, error) {
	if len(headerStr) != 160 {
		return nil, errors.New("block header should be 80 bytes long")
	}

	headerBytes, err := hex.DecodeString(headerStr)
	if err != nil {
		return nil, err
	}

	return NewBlockHeaderFromBytes(headerBytes)
}

// NewBlockHeaderFromBytes will encode a block header byte slice
// into the bitcoin block header structure.
//
// See https://en.bitcoin.it/wiki/Block_hashing_algorithm
func NewBlockHeaderFromBytes(headerBytes []byte) (*BlockHeader, error) {
	if len(headerBytes) != 80 {
		return nil, errors.New("block header should be 80 bytes long")
	}

	return &BlockHeader{
		Version:        binary.LittleEndian.Uint32(headerBytes[:4]),
		HashPrevBlock:  bt.ReverseBytes(headerBytes[4:36]),
		HashMerkleRoot: bt.ReverseBytes(headerBytes[36:68]),
		Time:           binary.LittleEndian.Uint32(headerBytes[68:72]),
		Bits:           bt.ReverseBytes(headerBytes[72:76]),
		Nonce:          binary.LittleEndian.Uint32(headerBytes[76:]),
	}, nil
}

// ExtractMerkleRootFromBlockHeader will take an 80 byte Bitcoin block
// header hex string and return the Merkle Root from it.
func ExtractMerkleRootFromBlockHeader(header string) (string, error) {
	bh, err := NewBlockHeaderFromStr(header)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bh.HashMerkleRoot), nil
}
