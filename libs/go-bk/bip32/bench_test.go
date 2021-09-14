// Copyright (c) 2014 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package bip32_test

import (
	"testing"

	"github.com/bitcoinsv/bsvutil/hdkeychain"

	"github.com/libsv/go-bk/bip32"
)

// bip0032MasterPriv1 is the master private extended key from the first set of
// test vectors in BIP0032.
const bip0032MasterPriv1 = "xprv9s21ZrQH143K3QTDL4LXw2F7HEK3wJUD2nW2nRk4stbP" +
	"y6cq3jPPqjiChkVvvNKmPGJxWUtg6LnF5kejMRNNU3TGtRBeJgk33yuGBxrMPHi"

// BenchmarkDeriveHardened benchmarks how long it takes to derive a hardened
// child from a master private extended key.
func BenchmarkDeriveHardened(b *testing.B) {
	b.StopTimer()
	masterKey, err := bip32.NewKeyFromString(bip0032MasterPriv1)
	if err != nil {
		b.Errorf("Failed to decode master seed: %v", err)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		masterKey.Child(hdkeychain.HardenedKeyStart)
	}
}

// BenchmarkDeriveNormal benchmarks how long it takes to derive a normal
// (non-hardened) child from a master private extended key.
func BenchmarkDeriveNormal(b *testing.B) {
	b.StopTimer()
	masterKey, err := hdkeychain.NewKeyFromString(bip0032MasterPriv1)
	if err != nil {
		b.Errorf("Failed to decode master seed: %v", err)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		masterKey.Child(0)
	}
}

// BenchmarkPrivToPub benchmarks how long it takes to convert a private extended
// key to a public extended key.
func BenchmarkPrivToPub(b *testing.B) {
	b.StopTimer()
	masterKey, err := hdkeychain.NewKeyFromString(bip0032MasterPriv1)
	if err != nil {
		b.Errorf("Failed to decode master seed: %v", err)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		masterKey.Neuter()
	}
}

// BenchmarkDeserialise benchmarks how long it takes to deserialize a private
// extended key.
func BenchmarkDeserialise(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hdkeychain.NewKeyFromString(bip0032MasterPriv1)
	}
}

// BenchmarkSerialise benchmarks how long it takes to serialise a private
// extended key.
func BenchmarkSerialise(b *testing.B) {
	b.StopTimer()
	masterKey, err := hdkeychain.NewKeyFromString(bip0032MasterPriv1)
	if err != nil {
		b.Errorf("Failed to decode master seed: %v", err)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_ = masterKey.String()
	}
}
