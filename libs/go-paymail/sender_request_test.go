package paymail

import (
	"fmt"
	"testing"
	"time"

	"github.com/bitcoinschema/go-bitcoin"
	"github.com/stretchr/testify/assert"
)

// TestSenderRequest_Sign will test the method Sign()
func TestSenderRequest_Sign(t *testing.T) {

	// Create key
	key, err := bitcoin.CreatePrivateKeyString()
	assert.NoError(t, err)
	assert.NotNil(t, key)

	// Create the request / message
	senderRequest := &SenderRequest{
		Dt:           time.Now().UTC().Format(time.RFC3339),
		SenderHandle: testAlias + "@" + testDomain,
		SenderName:   testName,
		Purpose:      testMessage,
	}

	var signature string

	t.Run("invalid key - empty", func(t *testing.T) {
		signature, err = senderRequest.Sign("")
		assert.Error(t, err)
		assert.Equal(t, len(signature), 0)
	})

	t.Run("invalid key - 0", func(t *testing.T) {
		signature, err = senderRequest.Sign("0")
		assert.Error(t, err)
		assert.Equal(t, len(signature), 0)
	})

	t.Run("invalid dt", func(t *testing.T) {
		senderRequest.Dt = ""
		signature, err = senderRequest.Sign(key)
		assert.Error(t, err)
		assert.Equal(t, len(signature), 0)
	})

	t.Run("invalid sender handle", func(t *testing.T) {
		senderRequest.Dt = time.Now().UTC().Format(time.RFC3339)
		senderRequest.SenderHandle = ""
		signature, err = senderRequest.Sign(key)
		assert.Error(t, err)
		assert.Equal(t, len(signature), 0)
	})

	t.Run("valid signature", func(t *testing.T) {
		senderRequest.SenderHandle = testAlias + "@" + testDomain
		signature, err = senderRequest.Sign(key)
		assert.NoError(t, err)
		assert.NotEqual(t, len(signature), 0)

		// Get address for verification
		var address string
		address, err = bitcoin.GetAddressFromPrivateKeyString(key, false)
		assert.NoError(t, err)

		// Verify the signature
		err = senderRequest.Verify(address, signature)
		assert.NoError(t, err)
	})
}

// ExampleSenderRequest_Sign example using Sign()
//
// See more examples in /examples/
func ExampleSenderRequest_Sign() {

	// Test private key
	key := "54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd"

	// Create the request / message
	senderRequest := &SenderRequest{
		Dt:           time.Now().UTC().Format(time.RFC3339),
		SenderHandle: testAlias + "@" + testDomain,
		SenderName:   testName,
		Purpose:      testMessage,
	}

	// Sign the sender request
	signature, err := senderRequest.Sign(key)
	if err != nil {
		fmt.Printf("error occurred in sign: %s", err.Error())
		return
	} else if len(signature) == 0 {
		fmt.Printf("signature was empty")
		return
	}

	// Cannot display signature as it changes because of the "dt" field
	fmt.Printf("signature created!")
	// Output:signature created!
}

// BenchmarkSenderRequest_Sign benchmarks the method Sign()
func BenchmarkSenderRequest_Sign(b *testing.B) {

	// Create the request / message
	senderRequest := &SenderRequest{
		Dt:           time.Now().UTC().Format(time.RFC3339),
		SenderHandle: testAlias + "@" + testDomain,
		SenderName:   testName,
		Purpose:      testMessage,
	}

	for i := 0; i < b.N; i++ {
		_, _ = senderRequest.Sign("54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd")
	}
}

// TestSenderRequest_Verify will test the method Verify()
func TestSenderRequest_Verify(t *testing.T) {

	// Create key
	key, err := bitcoin.CreatePrivateKeyString()
	assert.NoError(t, err)
	assert.NotNil(t, key)

	// Create the request / message
	senderRequest := &SenderRequest{
		Dt:           time.Now().UTC().Format(time.RFC3339),
		SenderHandle: testAlias + "@" + testDomain,
		SenderName:   testName,
		Purpose:      testMessage,
	}

	// Sign
	var signature string
	signature, err = senderRequest.Sign(key)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(signature))

	// Get address from private key
	var address string
	address, err = bitcoin.GetAddressFromPrivateKeyString(key, false)
	assert.NoError(t, err)
	assert.NotNil(t, address)

	t.Run("valid verification", func(t *testing.T) {
		err = senderRequest.Verify(address, signature)
		assert.NoError(t, err)
	})

	t.Run("invalid - empty address", func(t *testing.T) {
		err = senderRequest.Verify("", signature)
		assert.Error(t, err)
	})

	t.Run("invalid - empty signature", func(t *testing.T) {
		err = senderRequest.Verify(address, "")
		assert.Error(t, err)
	})

	t.Run("invalid - wrong signature - hex short", func(t *testing.T) {
		err = senderRequest.Verify(address, "0")
		assert.Error(t, err)
	})

	t.Run("invalid - wrong signature", func(t *testing.T) {
		err = senderRequest.Verify(address, "73646661736466736466617364667364666173646673646661736466")
		assert.Error(t, err)
	})
}

// ExampleSenderRequest_Verify example using Verify()
//
// See more examples in /examples/
func ExampleSenderRequest_Verify() {

	// Example sender request
	senderRequest := &SenderRequest{
		Dt:           "2020-10-02T16:43:39Z",
		SenderHandle: testAlias + "@" + testDomain,
		SenderName:   testName,
		Purpose:      testMessage,
	}

	// Try verifying (valid) (using an address and a signature - previously generated for example)
	if err := senderRequest.Verify(
		"1LqWAxSaKdXRKATAj7ELk34ioyT1T8gXgU",
		"G70DPE2p8xtCehUjRkQF2gI26kDu59JsQ6KKUmJyHi1XFGkeoIokgzN/kiMy+lujpXOi+C35sZUwgSMqOYRDXPQ=",
	); err != nil {
		fmt.Printf("error occurred in Verify: %s", err.Error())
		return
	}

	fmt.Printf("signature verified!")
	// Output:signature verified!
}

// BenchmarkSenderRequest_Verify benchmarks the method Verify()
func BenchmarkSenderRequest_Verify(b *testing.B) {

	// Example sender request
	senderRequest := &SenderRequest{
		Dt:           "2020-10-02T16:43:39Z",
		SenderHandle: testAlias + "@" + testDomain,
		SenderName:   testName,
		Purpose:      testMessage,
	}

	for i := 0; i < b.N; i++ {
		_ = senderRequest.Verify(
			"1LqWAxSaKdXRKATAj7ELk34ioyT1T8gXgU",
			"G70DPE2p8xtCehUjRkQF2gI26kDu59JsQ6KKUmJyHi1XFGkeoIokgzN/kiMy+lujpXOi+C35sZUwgSMqOYRDXPQ=",
		)
	}
}
