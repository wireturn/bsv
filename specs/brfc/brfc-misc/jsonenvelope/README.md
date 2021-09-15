# JSON Envelope Specification

|     BRFC     |    title     | authors | version |
| :----------: | :----------: | :-----: | :-----: |
| 298e080a4598 | jsonEnvelope | nChain  |   0.1   |

## Overview

Standard for serializing a JSON document in order to have consistency when ECDSA signing the document. Any changes to a document being signed and verified, however minor they may be, will cause the signature verification to fail since the document will be converted into a string before being (hashed and then) signed. With JSON documents, the format permits changes to be made without compromising the validity of the format (eg. extra spaces, carriage returns, etc.).

This spec describes a technique to ensure consistency of the data being signed by encapsulting the JSON data _as a string_ in parent JSON object. That way, however the JSON is marshalled, the first element in the parent JSON, the payload, would remain the same and be signed/verified the way it is.

## JSON Envelope Example

```json
{
  "payload": "{\"name\":\"simon\",\"colour\":\"blue\"}",
  "signature": "30450221008209b19ffe2182d859ce36fdeff5ded4b3f70ad77e0e8715238a539db97c1282022043b1a5b260271b7c833ca7c37d1490f21b7bd029dbb8970570c7fdc3df5c93ab",
  "publicKey": "02b01c0c23ff7ff35f774e6d3b3491a123afb6c98965054e024d2320f7dbd25d8a",
  "encoding": "UTF-8",
  "mimetype": "application/json"
}
```

| field       | function                       |
| ----------- | ------------------------------ |
| `payload`   | payload of data being sent     |
| `signature` | signature on payload (string)  |
| `publicKey` | public key to verify signature |
| `encoding`  | encoding of the payload data   |
| `mimetype`  | mimetype of the payload data   |

> An API that returns a JSON envelope may not always want to sign the payload. In this scenario, the signature and publicKey fields may be set to _null_.

### Default recommendations

In case of binary data in the `payload`, directly convert `payload` to binary stream (be it binary, hex, base64, etc.), which is to be used as the input to the hash function.  
In case of string data in the `payload`, convert the characters to bytes by using UTF-8, which are then to be used as the input to the hash function.

## Extra Examples

### Image:

```json
{
  "payload": "/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAAYEBQYFBAYGBQYHBwYIChAKCgkJChQODwwQFxQYGBcUFhYaHSUfGhsjHBYWICwgIyYnKSopGR8tMC0oMCUoKSj/2wBDAQcHBwoIChMKChMoGhYaKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCj/wgARCAAKAAoDASIAAhEBAxEB/8QAFwAAAwEAAAAAAAAAAAAAAAAAAQMEB//EABUBAQEAAAAAAAAAAAAAAAAAAAAC/9oADAMBAAIQAxAAAAGjQS5H/8QAGBABAQEBAQAAAAAAAAAAAAAAAgMEAAX/2gAIAQEAAQUC9HbYaJJKTmE+/8QAFhEBAQEAAAAAAAAAAAAAAAAAAgAR/9oACAEDAQE/AScv/8QAFBEBAAAAAAAAAAAAAAAAAAAAAP/aAAgBAgEBPwF//8QAHxAAAgECBwAAAAAAAAAAAAAAAQIDABEEEBMUITFR/9oACAEBAAY/AsUdwY2jI04/aQslmI5oMyKWHRtl/8QAGhABAAMAAwAAAAAAAAAAAAAAAQARMRBBUf/aAAgBAQABPyFOB4HPtd3KY4ovGpvYACnH/9oADAMBAAIAAwAAABDb/8QAFhEBAQEAAAAAAAAAAAAAAAAAAREA/9oACAEDAQE/EESrbv/EABYRAQEBAAAAAAAAAAAAAAAAAAEAEf/aAAgBAgEBPxBdv//EABsQAQACAgMAAAAAAAAAAAAAAAEAIRARgaHB/9oACAEBAAE/EDyaVPB5UlLVgU4Z1DdcVNmP/9k=",
  "signature": "3045022100ebfde614a67d6f69c321664683b557a2eb605d7aa9357230684f49c1da4ccbef02203ab72beb9ffe1af76cb60b852b950baa2355c32ceb99715158e7e2d31a194f1d",
  "publicKey": "02aaee936deeb6d8296aa11d3134c624a2d8e72581ce49c73237f0359e4cf11949",
  "encoding": "base64",
  "mimetype": "image/jpeg"
}
```

## Implementations in different languages

### JavaScript

```javascript
const bsv = require('bsv')

const privateKey = bsv.PrivateKey.fromRandom()
const publicKey = privateKey.toPublicKey().toString()

const payload = {
  name: 'simon',
  colour: 'blue'
}

// SIGN
const payloadStr = JSON.stringify(payload)
const hash = bsv.crypto.Hash.sha256(Buffer.from(payloadStr))
const sig = bsv.crypto.ECDSA.sign(hash, privateKey).toString()

const jsonEnvelope = {
  payload: payloadStr,
  signature: sig,
  publicKey: publicKey,
  encoding: 'UTF-8',
  mimetype: 'application/json'
}

// VERIFY
const verifyHash = bsv.crypto.Hash.sha256(Buffer.from(jsonEnvelope.payload))
const signature = bsv.crypto.Signature.fromString(jsonEnvelope.signature)
const verifyPublicKey = bsv.PublicKey(jsonEnvelope.publicKey)
const verified = bsv.crypto.ECDSA.verify(verifyHash, signature, verifyPublicKey)

console.log('Signature Verified: ', verified)
```

### Golang

```go
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
		Encoding:  "UTF-8",
	}

	// VERIFY
	verifyPubKey, err := btcec.ParsePubKey([]byte(*jsonEnvelope.PublicKey), btcec.S256())
	if err != nil {
		fmt.Println(err)
		return
	}
	verifyHash := sha256.Sum256([]byte(jsonEnvelope.Payload))

	verified := signature.Verify(verifyHash[:], verifyPubKey)
	fmt.Printf("Signature Verified: %v\n", verified)

}

```
