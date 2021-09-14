package b

import (
	"fmt"
	"testing"

	"github.com/bitcoinschema/go-bob"
)

const exampleFileURLKong = "https://x.bitfs.network/6ce94f75b88a6c24815d480437f4f06ae895afdab8039ddec10748660c29f910.out.0.3"
const exampleFileURLGif = "https://x.bitfs.network/10afc796d06fec11a4b6077012a1522355c82e5de316f4dd5c42ddccd6d61cdb.out.0.3"
const exampleTxKong = "6ce94f75b88a6c24815d480437f4f06ae895afdab8039ddec10748660c29f910"

// TestBitFsURL tests for nil case in BitFsURL()
func TestBitFsURL(t *testing.T) {
	bitURL := BitFsURL(exampleTxKong, 0, 3)

	if bitURL != exampleFileURLKong {
		t.Fatalf("failed url: %s", bitURL)
	}
}

// ExampleBitFsURL example using BitFsURL()
func ExampleBitFsURL() {
	bitURL := BitFsURL(exampleTxKong, 0, 3)
	fmt.Printf("url: %s", bitURL)
	// Output:url: https://x.bitfs.network/6ce94f75b88a6c24815d480437f4f06ae895afdab8039ddec10748660c29f910.out.0.3
}

// BenchmarkBitFsURL benchmarks the method BitFsURL()
func BenchmarkBitFsURL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = BitFsURL(exampleTxKong, 0, 3)
	}
}

// TestDataURI tests for nil case in DataURI()
func TestDataURI(t *testing.T) {

	// Convert from string
	bobData, err := bob.NewFromString(exampleBobTx)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	// NOTE: The above data does not contain the actual B information from the tx it has been stripped out

	// Start from a tape
	var bData *B
	bData, err = NewFromTape(bobData.Out[0].Tape[1])
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	// Get the data URI
	dataURI := bData.DataURI()
	if dataURI != "data:binary;base64," {
		t.Fatalf("expected dataURI: %s got: %s", "", dataURI)
	}

	// Create the bitfs URL
	bitfsURL := BitFsURL(bobData.Tx.H, 0, 3)
	if bitfsURL != exampleFileURLGif {
		t.Fatalf("expected bitfsURL: %s got: %s", "", bitfsURL)
	}
}

// ExampleB_DataURI example using DataURI()
func ExampleB_DataURI() {
	// Convert from string
	bobData, _ := bob.NewFromString(exampleBobTx)
	bData, _ := NewFromTape(bobData.Out[0].Tape[1])

	// NOTE: The above data does not contain the actual B information from the tx it has been stripped out

	fmt.Printf("data URI: %s", bData.DataURI())
	// Output:data URI: data:binary;base64,
}

// BenchmarkB_DataURI benchmarks the method DataURI()
func BenchmarkB_DataURI(b *testing.B) {
	bobData, _ := bob.NewFromString(exampleBobTx)
	bData, _ := NewFromTape(bobData.Out[0].Tape[1])
	for i := 0; i < b.N; i++ {
		bData.DataURI()
	}
}
