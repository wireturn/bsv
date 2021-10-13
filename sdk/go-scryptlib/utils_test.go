package scryptlib


import (
    "testing"
    "math/big"

    "github.com/stretchr/testify/assert"
)

func TestStructCompare(t *testing.T) {
    personKeysInOrder := []string{"name", "nicknames", "height", "dog"}
    dogKeysInOrder := []string{"name", "breed"}

    dog0Name := Bytes{[]byte("Rex")}
    dog0Breed := Bytes{[]byte("Beagle")}
    dog0Values := map[string]ScryptType{
        "name": dog0Name,
        "breed": dog0Breed,
    }
    dog0 := Struct{
        keysInOrder: dogKeysInOrder,
        values: dog0Values,
    }

    person0Name := Bytes{[]byte("Alice")}
    var person0NicknamesVals []ScryptType
    person0NicknamesVals = append(person0NicknamesVals, Bytes{[]byte("Alie")})
    person0NicknamesVals = append(person0NicknamesVals, Bytes{[]byte("A")})
    person0Nicknames := Array{person0NicknamesVals}
    person0Height := Int{big.NewInt(192)}
    person0Values := map[string]ScryptType{
        "name": person0Name,
        "nicknames": person0Nicknames,
        "height": person0Height,
        "dog": dog0,
    }
    person0 := Struct{
        keysInOrder: personKeysInOrder,
        values: person0Values,
    }

    res0 := IsStructsSameStructure(dog0, person0)
    assert.Equal(t, res0, false)

    res1 := IsStructsSameStructure(person0, person0)
    assert.Equal(t, res1, true)

    res2 := IsStructsSameStructure(dog0, dog0)
    assert.Equal(t, res2, true)

    res3 := IsStructsSameStructure(person0, dog0)
    assert.Equal(t, res3, false)
}
