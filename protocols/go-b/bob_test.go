package b

import (
	"testing"

	"github.com/bitcoinschema/go-bob"
)

// TestNewFromTape tests for nil case in NewFromTape()
func TestNewFromTape(t *testing.T) {
	expectedTx := "10afc796d06fec11a4b6077012a1522355c82e5de316f4dd5c42ddccd6d61cdb"
	expectedEncoding := "binary"

	tx, err := bob.NewFromString(exampleBobTx)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}
	if tx.Tx.H != expectedTx {
		t.Fatalf("expected Tx.h: %s got: %s", expectedTx, tx.Tx.H)
	}

	// Invalid tape
	var b *B
	_, err = NewFromTape(tx.Out[0].Tape[0])
	if err == nil {
		t.Fatalf("error should have occurred")
	}

	// Valid tape
	b, err = NewFromTape(tx.Out[0].Tape[1])
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}
	if b.Encoding != expectedEncoding {
		t.Fatalf("expected Encoding: %s got: %s", expectedEncoding, b.Encoding)
	}

	// todo: finish tests (error cases) and examples (need BOB txs)
}

// TestNewFromTapes tests for nil case in NewFromTapes()
func TestNewFromTapes(t *testing.T) {

	expectedTx := "10afc796d06fec11a4b6077012a1522355c82e5de316f4dd5c42ddccd6d61cdb"
	expectedEncoding := "binary"

	tx, err := bob.NewFromString(exampleBobTx)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}
	if tx.Tx.H != expectedTx {
		t.Fatalf("expected Tx.h: %s got: %s", expectedTx, tx.Tx.H)
	}

	// Valid
	var b *B
	b, err = NewFromTapes(tx.Out[0].Tape)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}
	if b.Encoding != expectedEncoding {
		t.Fatalf("expected Encoding: %s got: %s", expectedEncoding, b.Encoding)
	}

	// Change type
	tx.Out[0].Tape[1].Cell[3].S = string(EncodingGzip)
	b, err = NewFromTapes(tx.Out[0].Tape)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}
	if b.Encoding != string(EncodingGzip) {
		t.Fatalf("expected Encoding: %s got: %s", EncodingGzip, b.Encoding)
	}

	// Invalid
	_, err = NewFromTapes(tx.Out[1].Tape)
	if err == nil {
		t.Fatalf("error should have occurred")
	}

	// todo: finish tests and examples (need BOB tx)
}
