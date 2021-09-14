// Package aip is a library for working with Author Identity Protocol (AIP) in Go
//
// If you have any suggestions or comments, please feel free to open an issue on
// this GitHub repository!
//
// By BitcoinSchema Organization (https://bitcoinschema.org)
package aip

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/bitcoinschema/go-bitcoin"
	"github.com/bitcoinsv/bsvutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/libsv/go-bt"
)

// Prefix is the Bitcom prefix used by AIP
const Prefix = "15PciHG22SNLQJXMoSUaWVi7WSqc7hCfva"
const pipe = "|"
const opReturn = string(rune(txscript.OP_RETURN)) // creates: j

// Algorithm is an enum for the different possible signature algorithms
type Algorithm string

// Algorithm names
const (
	BitcoinECDSA         Algorithm = "BITCOIN_ECDSA"        // Backwards compatible for BitcoinSignedMessage
	BitcoinSignedMessage Algorithm = "BitcoinSignedMessage" // New algo name
	Paymail              Algorithm = "paymail"              // Using "pubkey" as aip.Address
)

// Aip is an Author Identity Protocol object
type Aip struct {
	Algorithm                 Algorithm `json:"algorithm"`                   // Known AIP algorithm type
	AlgorithmSigningComponent string    `json:"algorithm_signing_component"` // Changes based on the Algorithm
	Data                      []string  `json:"data"`                        // Data to be signed or validated
	Indices                   []int     `json:"indices,omitempty"`           // BOB indices
	Signature                 string    `json:"signature"`                   // AIP generated signature
}

// Validate returns true if the given AIP signature is valid for given data
func (a *Aip) Validate() (bool, error) {

	// Both data and component are required
	if len(a.Data) == 0 || len(a.AlgorithmSigningComponent) == 0 {
		return false, errors.New("missing data or signing component")
	}

	// Check to be sure OP_RETURN was prepended before trying to validate
	if a.Data[0] != opReturn {
		return false, fmt.Errorf("the first item in payload is always OP_RETURN, got: %s", a.Data[0])
	}

	// Convert pubkey to address
	if a.Algorithm == Paymail {
		// Detect whether this key was compressed when sig was made
		_, wasCompressed, err := bitcoin.PubKeyFromSignature(a.Signature, strings.Join(a.Data, ""))
		if err != nil {
			return false, err
		}

		// Get the public address for this paymail from pki
		var addr *bsvutil.LegacyAddressPubKeyHash
		if addr, err = bitcoin.GetAddressFromPubKeyString(a.AlgorithmSigningComponent, wasCompressed); err != nil {
			return false, err
		}
		a.AlgorithmSigningComponent = addr.String()
	}

	// You get the address associated with the pki instead of the current address
	err := bitcoin.VerifyMessage(a.AlgorithmSigningComponent, a.Signature, strings.Join(a.Data, ""))
	return err == nil, err
}

// Sign will provide an AIP signature for a given private key and message using
// the provided algorithm. It prepends an OP_RETURN to the payload
func Sign(privateKey string, algorithm Algorithm, message string) (a *Aip, err error) {

	// Prepend the OP_RETURN to keep consistent with BitcoinFiles SDK
	// data = append(data, []byte{byte(txscript.OP_RETURN)})
	prependedData := []string{opReturn, message}

	// Create the base AIP object
	a = &Aip{Algorithm: algorithm, Data: prependedData}

	// Sign using the private key and the message
	if a.Signature, err = bitcoin.SignMessage(privateKey, strings.Join(prependedData, ""), false); err != nil {
		return
	}

	// Store address vs pubkey
	switch algorithm {
	case BitcoinECDSA, BitcoinSignedMessage:
		// Signing component = bitcoin address
		// Get the address of the private key
		if a.AlgorithmSigningComponent, err = bitcoin.GetAddressFromPrivateKeyString(privateKey, false); err != nil {
			return
		}
	case Paymail:
		// Signing component = paymail identity key
		// Get pubKey from private key and overload the address field in AIP
		if a.AlgorithmSigningComponent, err = bitcoin.PubKeyFromPrivateKeyString(privateKey, false); err != nil {
			return
		}
	}

	return
}

// SignOpReturnData will append the given data and return an bt.Output
func SignOpReturnData(privateKey string, algorithm Algorithm,
	data [][]byte) (out *bt.Output, outData [][]byte, a *Aip, err error) {

	// Sign with AIP
	if a, err = Sign(privateKey, algorithm, string(bytes.Join(data, []byte{}))); err != nil {
		return
	}

	// Add AIP signature
	outData = append(
		data,
		[]byte(Prefix),
		[]byte(a.Algorithm),
		[]byte(a.AlgorithmSigningComponent),
		[]byte(a.Signature),
	)

	// Create the output
	out, err = bt.NewOpReturnPartsOutput(outData)
	return
}
