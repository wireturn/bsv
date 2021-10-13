package handcash

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/bitcoinsv/bsvd/bsvec"
)

// getRequestSignature will return the request signature
//
// Specs: https://github.com/HandCash/handcash-connect-sdk-js/blob/00300de6d225fa37fe2f4a5efe315dd08dd4beb9/src/api/http_request_factory.js#L16
func getRequestSignature(method, endpoint string, body interface{}, timestamp string,
	privateKey *bsvec.PrivateKey) ([]byte, error) {

	// Create the signature hash
	signatureHash, err := getRequestSignatureHash(method, endpoint, body, timestamp)
	if err != nil {
		return nil, err
	}

	// Sign using private key
	var sig *bsvec.Signature
	if sig, err = privateKey.Sign(signatureHash); err != nil {
		return nil, err
	}

	// Return the serialized signature string
	return sig.Serialize(), nil
}

// getRequestSignatureHash will return the signature hash
//
// Specs: https://github.com/HandCash/handcash-connect-sdk-js/blob/00300de6d225fa37fe2f4a5efe315dd08dd4beb9/src/api/http_request_factory.js#L34
func getRequestSignatureHash(method, endpoint string, body interface{},
	timestamp string) ([]byte, error) {

	// Default if not set
	bodyString := emptyBody
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body %w", err)
		}
		bodyString = string(bodyBytes)
	}

	// Set the signature string
	signatureString := fmt.Sprintf("%s\n%s\n%s\n%s", method, endpoint, timestamp, bodyString)
	hash := sha256.Sum256([]byte(signatureString))
	return hash[:], nil
}

// getSignedRequest returns the request with signature
//
// Specs: https://github.com/HandCash/handcash-connect-sdk-js/blob/00300de6d225fa37fe2f4a5efe315dd08dd4beb9/src/api/http_request_factory.js#L34
func (c *Client) getSignedRequest(method, endpoint, authToken string,
	body interface{}, timestamp string) (*signedRequest, error) {

	// Decode token
	tokenBytes, err := hex.DecodeString(authToken)
	if err != nil {
		return nil, err
	}

	// Get key pairs
	privateKey, publicKey := bsvec.PrivKeyFromBytes(bsvec.S256(), tokenBytes)

	// Get the request signature
	var requestSignature []byte
	if requestSignature, err = getRequestSignature(
		method, endpoint, body, timestamp, privateKey,
	); err != nil {
		return nil, err
	}

	// Return the signed request
	return &signedRequest{
		Body: body,
		Headers: oAuthHeaders{
			OauthPublicKey: hex.EncodeToString(publicKey.SerializeCompressed()),
			OauthSignature: hex.EncodeToString(requestSignature),
			OauthTimestamp: timestamp,
		},
		JSON:   true,
		Method: method,
		URI:    c.Environment.APIURL + endpoint,
	}, nil
}
