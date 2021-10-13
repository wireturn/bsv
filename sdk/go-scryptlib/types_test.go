package scryptlib


import (
    "testing"
    "math/big"

    "github.com/stretchr/testify/assert"
)


func TestTypesInt(t *testing.T) {
    bigIntObj := big.NewInt(0)
    intObj := Int { value: bigIntObj }
    hex, _ := intObj.Hex()
    assert.Equal(t, "00", hex)

    bigIntObj = big.NewInt(1)
    intObj = Int { value: bigIntObj }
    hex, _ = intObj.Hex()
    assert.Equal(t, "51", hex)

    bigIntObj = big.NewInt(16)
    intObj = Int { value: bigIntObj }
    hex, _ = intObj.Hex()
    assert.Equal(t, "60", hex)

    bigIntObj = big.NewInt(17)
    intObj = Int { value: bigIntObj }
    hex, _ = intObj.Hex()
    assert.Equal(t, "0111", hex)

    bigIntObj = big.NewInt(129)
    intObj = Int { value: bigIntObj }
    hex, _ = intObj.Hex()
    assert.Equal(t, "028100", hex)

    bigIntObj = big.NewInt(-243)
    intObj = Int { value: bigIntObj }
    hex, _ = intObj.Hex()
    assert.Equal(t, "02f380", hex)
}
