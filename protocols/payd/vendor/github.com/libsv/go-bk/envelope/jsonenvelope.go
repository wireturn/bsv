// Package envelope supports the JSON Envelope Spec
// It can be found here https://github.com/bitcoin-sv-specs/brfc-misc/tree/master/jsonenvelope
//
// Standard for serialising a JSON document in order to have consistency when
// ECDSA signing the document.
// Any changes to a document being signed and verified, however minor they may be,
// will cause the signature verification to fail since the document will be converted into a string
// before being (hashed and then) signed. With JSON documents, the format permits changes to be made
// without compromising the validity of the format (eg. extra spaces, carriage returns, etc.).
package envelope

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/libsv/go-bk/bec"
)

// JSONEnvelope defines an envelop contain option signature and public key to use to
// verify the signature and payload.
// If no Signature or PublicKey are provided we do not validate the payload.
// The payload is usually an escaped JSON string.
type JSONEnvelope struct {
	Payload   string  `json:"payload"`
	Signature *string `json:"signature"`
	PublicKey *string `json:"publicKey"`
	Encoding  string  `json:"encoding"`
	MimeType  string  `json:"mimetype"`
}

// IsValid will check that a JSONEnvelope is valid by using the PublicKey and Signature to
// validate the payload. If the payload differs from the signature, false is returned.
// If the signature or public key are invalid an error is returned.
func (j *JSONEnvelope) IsValid() (bool, error) {
	// if no sig or pk we don't try to verify
	if j.Signature == nil && j.PublicKey == nil {
		return true, nil
	}
	// parse and validate public key
	pub, err := hex.DecodeString(*j.PublicKey)
	if err != nil {
		return false, fmt.Errorf("failed to decode json envelope publicKey %w", err)
	}
	verifyPubKey, err := bec.ParsePubKey(pub, bec.S256())
	if err != nil {
		return false, fmt.Errorf("failed to parse json envelope publicKey %w", err)
	}

	// parse and validate signature
	signature, err := hex.DecodeString(*j.Signature)
	if err != nil {
		return false, fmt.Errorf("failed to decode json envelope signature %w", err)
	}
	sig, err := bec.ParseSignature(signature, bec.S256())
	if err != nil {
		return false, fmt.Errorf("failed to parse json envelope signature %w", err)
	}
	var verifyHash [32]byte
	switch j.MimeType {
	case "application/json":
		verifyHash = sha256.Sum256([]byte(strings.Replace(j.Payload, `\`, "", -1)))
	case "base64":
		bb, err := base64.StdEncoding.DecodeString(j.Payload)
		if err != nil {
			return false, fmt.Errorf("failed to parse base64 payload %w", err)
		}
		verifyHash = sha256.Sum256(bb)
	default:
		verifyHash = sha256.Sum256([]byte(j.Payload))
	}
	return sig.Verify(verifyHash[:], verifyPubKey), nil
}

// NewJSONEnvelope will create and return a new JSONEnvelope with the provided
// payload serialised and signed.
func NewJSONEnvelope(payload interface{}) (*JSONEnvelope, error) {
	privateKey, err := bec.NewPrivateKey(bec.S256())
	if err != nil {
		return nil, fmt.Errorf("failed to generate new private key %w", err)
	}
	publicKey := privateKey.PubKey()

	pl, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode payload %w", err)
	}
	hash := sha256.Sum256(pl)
	signature, err := privateKey.Sign(hash[:])
	if err != nil {
		return nil, fmt.Errorf("failed to create signature %w", err)
	}
	signatureHex := hex.EncodeToString(signature.Serialise())
	publicKeyHex := hex.EncodeToString(publicKey.SerialiseCompressed())

	return &JSONEnvelope{
		Payload:   string(pl),
		Signature: &signatureHex,
		PublicKey: &publicKeyHex,
		Encoding:  "UTF-8",
		MimeType:  "application/json",
	}, nil
}
