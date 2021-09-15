package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"github.com/btcsuite/btcd/btcec"
)

type jsonEnvelope struct {
	Payload   string  `json:"payload"`
	Signature *string `json:"signature"`
	PublicKey *string `json:"publicKey"`
	Encoding  string  `json:"encoding"`
	MimeType  string  `json:"mimetype"`
}

type test struct {
	Name   string `json:"name"`
	Colour string `json:"colour"`
}

func main() {
	privateKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		fmt.Println(err)
		return
	}

	publicKey := privateKey.PubKey()

	payload, err := json.Marshal(test{
		Name:   "simon",
		Colour: "blue",
	})
	if err != nil {
		log.Println(err)
		return
	}

	hash := sha256.Sum256([]byte(payload))

	// SIGN
	signature, err := privateKey.Sign(hash[:])
	if err != nil {
		fmt.Println(err)
		return
	}

	signatureHex := hex.EncodeToString(signature.Serialize())
	publicKeyHex := hex.EncodeToString(publicKey.SerializeCompressed())

	jsonEnvelope := &jsonEnvelope{
		Payload:   string(payload),
		Signature: &signatureHex,
		PublicKey: &publicKeyHex,
		MimeType:  "applicaton/json",
		Encoding:  "json",
	}

	pub, err := hex.DecodeString(*jsonEnvelope.PublicKey)
	if err != nil {
		fmt.Println(err)
		return
	}

	// VERIFY
	verifyPubKey, err := btcec.ParsePubKey(pub, btcec.S256())
	if err != nil {
		fmt.Println(err)
		return
	}
	verifyHash := sha256.Sum256([]byte(jsonEnvelope.Payload))

	verified := signature.Verify(verifyHash[:], verifyPubKey)
	fmt.Printf("Signature Verified: %v\n", verified)

}
