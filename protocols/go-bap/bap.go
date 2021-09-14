// Package bap is a library for working with Bitcoin Attestation Protocol (BAP) in Go
//
// Protocol: https://github.com/icellan/bap
//
// If you have any suggestions or comments, please feel free to open an issue on
// this GitHub repository!
//
// By BitcoinSchema Organization (https://bitcoinschema.org)
package bap

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/bitcoinschema/go-aip"
	"github.com/libsv/go-bt"
)

// Prefix is the bitcom prefix for Bitcoin Attestation Protocol (BAP)
const Prefix = "1BAPSuaPnfGnSBM3GLV9yhxUdYe4vGbdMT"
const pipe string = "|"

// AttestationType is an enum for BAP Type Constants
type AttestationType string

// BAP attestation type constants
const (
	ATTEST AttestationType = "ATTEST"
	ID     AttestationType = "ID"
	REVOKE AttestationType = "REVOKE"
)

// CreateIdentity creates an identity from a private key, an id key, and a counter
//
// Source: https://github.com/icellan/bap
func CreateIdentity(privateKey, idKey string, currentCounter uint32) (*bt.Tx, error) {

	// Test for id key
	if len(idKey) == 0 {
		return nil, fmt.Errorf("missing required field: %s", "idKey")
	}

	// Derive the keys
	newSigningPrivateKey, newAddress, err := deriveKeys(privateKey, currentCounter+1) // Increment the next key
	if err != nil {
		return nil, err
	}

	// Create the identity attestation op_return data
	var data [][]byte
	data = append(
		data,
		[]byte(Prefix),
		[]byte(ID),
		[]byte(idKey),
		[]byte(newAddress),
		[]byte(pipe),
	)

	// Generate a signature from this point
	var finalOutput *bt.Output
	if finalOutput, _, _, err = aip.SignOpReturnData(newSigningPrivateKey, aip.BitcoinECDSA, data); err != nil {
		return nil, err
	}

	// Return the transaction
	return returnTx(finalOutput), nil
}

// CreateAttestation creates an attestation transaction from an id key, signing key, and signing address
//
// Source: https://github.com/icellan/bap
func CreateAttestation(idKey, attestorSigningKey, attributeName,
	attributeValue, identityAttributeSecret string) (*bt.Tx, error) {

	// ID key is required
	if len(idKey) == 0 {
		return nil, errors.New("missing required field: idKey")
	}

	// Attribute secret and name
	if len(attributeName) == 0 {
		return nil, errors.New("missing required field: attributeName")
	} else if len(identityAttributeSecret) == 0 {
		return nil, errors.New("missing required field: identityAttributeSecret")
	}

	// Attest that an internal wallet address is associated with our identity key
	idUrn := fmt.Sprintf("urn:bap:id:%s:%s:%s", attributeName, attributeValue, identityAttributeSecret)
	attestationUrn := fmt.Sprintf("urn:bap:attest:%v:%s", sha256.Sum256([]byte(idUrn)), idKey)
	attestationHash := sha256.Sum256([]byte(attestationUrn))

	// Create op_return attestation
	var data [][]byte
	data = append(
		data,
		[]byte(Prefix),
		[]byte(ATTEST),
		attestationHash[0:],
		[]byte(pipe),
	)

	// Generate a signature from this point
	finalOutput, _, _, err := aip.SignOpReturnData(attestorSigningKey, aip.BitcoinECDSA, data)
	if err != nil {
		return nil, err
	}

	// Return the transaction
	return returnTx(finalOutput), nil
}

// returnTx will add the output and return a new tx
func returnTx(out *bt.Output) (t *bt.Tx) {
	t = bt.NewTx()
	t.AddOutput(out)
	return
}
