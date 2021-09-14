package bap

import (
	"github.com/bitcoinschema/go-bitcoin"
	"github.com/bitcoinsv/bsvutil/hdkeychain"
)

// deriveKeys will return the xPriv for the identity key and the corresponding address
func deriveKeys(xPrivateKey string, currentCounter uint32) (xPriv string, address string, err error) {

	// Get the raw private key from string into an HD key
	var hdKey *hdkeychain.ExtendedKey
	if hdKey, err = bitcoin.GenerateHDKeyFromString(xPrivateKey); err != nil {
		return
	}

	// Get id key
	var idKey *hdkeychain.ExtendedKey // m/0/N
	if idKey, err = bitcoin.GetHDKeyByPath(hdKey, 0, currentCounter); err != nil {
		return
	}

	// Get the address
	if address, err = bitcoin.GetAddressStringFromHDKey(idKey); err != nil {
		return
	}

	// Get the private key from the identity key
	xPriv, err = bitcoin.GetPrivateKeyStringFromHDKey(idKey)

	return
}
