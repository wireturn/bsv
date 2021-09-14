package bc

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"log"
	"sort"
)

// implement `Interface` in sort package.
type sortByteArrays [][]byte

func (b sortByteArrays) Len() int {
	return len(b)
}

func (b sortByteArrays) Less(i, j int) bool {
	// bytes package already implements Comparable for []byte.
	switch bytes.Compare(b[i], b[j]) {
	case -1:
		return true
	case 0, 1:
		return false
	default:
		log.Panic("not fail-able with `bytes.Comparable` bounded [-1, 1].")
		return false
	}
}

func (b sortByteArrays) Swap(i, j int) {
	b[j], b[i] = b[i], b[j]
}

// SortByteArrays comment TODO:
func SortByteArrays(src [][]byte) [][]byte {
	sorted := sortByteArrays(src)
	sort.Sort(sorted)
	return sorted
}

// ReverseHexString reverses the hex string (little endian/big endian).
// This is used when computing merkle trees in Bitcoin, for example.
func ReverseHexString(hex string) string {
	res := ""
	if len(hex)%2 != 0 {
		hex = "0" + hex
	}

	for i := len(hex); i >= 2; i -= 2 {

		res += hex[i-2 : i]
	}
	return res
}

// Equals checks if two byte arrays are equal.
func Equals(b1 []byte, b2 []byte) bool {
	if len(b1) != len(b2) {
		return false
	}
	for i, x := range b1 {
		if x != b2[i] {
			return false
		}
	}
	return true
}

// Decode32Byte decodes a hex string into a 32 byte array.
func Decode32Byte(hexStr string) ([32]byte, error) {
	var b32 [32]byte
	b, err := hex.DecodeString(hexStr)
	if err != nil {
		return b32, err
	}

	copy(b32[:], b[0:32])

	return b32, nil
}

// UInt32ToBytes converts a uint32 into an
// array of bytes.
func UInt32ToBytes(num uint32) []byte {
	t := make([]byte, 4)
	binary.LittleEndian.PutUint32(t, num)
	return t
}
