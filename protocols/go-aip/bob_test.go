package aip

import (
	"fmt"
	"testing"

	"github.com/bitcoinschema/go-bitcoin"
	"github.com/bitcoinschema/go-bob"
)

// TestNewFromTape will test the method NewFromTape()
func TestNewFromTape(t *testing.T) {
	t.Parallel()

	// Parse from string into BOB
	bobValidData, err := bob.NewFromString(sampleValidBobTx)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}
	var bobInvalidData *bob.Tx
	if bobInvalidData, err = bob.NewFromString(sampleInvalidBobTx); err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	var (
		// Testing private methods
		tests = []struct {
			inputTapes         []bob.Tape
			inputIndex         int
			expectedSignature  string
			expectedAlgorithm  Algorithm
			expectedComponent  string
			expectedNil        bool
			expectedValidation bool
		}{
			{
				bobValidData.Out[0].Tape,
				2,
				"H+lubfcz5Z2oG8B7HwmP8Z+tALP+KNOPgedo7UTXwW8LBpMkgCgatCdpvbtf7wZZQSIMz83emmAvVS4S3F5X1wo=",
				BitcoinECDSA,
				"134a6TXxzgQ9Az3w8BcvgdZyA5UqRL89da",
				false,
				true,
			},
			{
				bobInvalidData.Out[0].Tape,
				2,
				"H+lubfcz5Z2oG8B7HwmP8Z+tALP+KNOPgedo7UTXwW8LBpMkgCgatCdpvbtf7wZZQSIMz83emmAvVS4S3F5X1wo=",
				BitcoinECDSA,
				"invalid-address",
				false,
				false,
			},
			{
				[]bob.Tape{*new(bob.Tape)},
				0,
				"",
				"",
				"",
				false,
				false,
			},
		}
	)

	// Run tests
	for idx, test := range tests {
		if a := NewFromTape(test.inputTapes[test.inputIndex]); a == nil && !test.expectedNil {
			t.Errorf("%d %s Failed: [%v] inputted and nil was not expected", idx, t.Name(), test.inputTapes[test.inputIndex])
		} else if a != nil && test.expectedNil {
			t.Errorf("%d %s Failed: [%v] inputted and nil was expected", idx, t.Name(), test.inputTapes[test.inputIndex])
		} else if a != nil && a.Signature != test.expectedSignature {
			t.Errorf("%d %s Failed: [%v] inputted and expected [%s] but got [%s]", idx, t.Name(), test.inputTapes[test.inputIndex], test.expectedSignature, a.Signature)
		} else if a != nil && a.Algorithm != test.expectedAlgorithm {
			t.Errorf("%d %s Failed: [%v] inputted and expected [%s] but got [%s]", idx, t.Name(), test.inputTapes[test.inputIndex], test.expectedAlgorithm, a.Algorithm)
		} else if a != nil && a.AlgorithmSigningComponent != test.expectedComponent {
			t.Errorf("%d %s Failed: [%v] inputted and expected [%s] but got [%s]", idx, t.Name(), test.inputTapes[test.inputIndex], test.expectedComponent, a.AlgorithmSigningComponent)
		} else if a != nil && len(test.inputTapes) > 1 {
			valid, err := ValidateTapes(test.inputTapes)
			if valid && !test.expectedValidation {
				t.Errorf("%d %s Failed: [%v] inputted and validation should have failed", idx, t.Name(), test.inputTapes)
			} else if !valid && test.expectedValidation && err != nil {
				t.Errorf("%d %s Failed: [%v] inputted and validation should have passed, error: %s", idx, t.Name(), test.inputTapes, err.Error())
			}
		}
	}
}

// ExampleNewFromTape example using NewFromTape()
func ExampleNewFromTape() {
	// Get BOB data from a TX
	bobValidData, err := bob.NewFromString(sampleValidBobTx)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Get from tape given the AIP index
	a := NewFromTape(bobValidData.Out[0].Tape[2])

	fmt.Printf("address: %s signature: %s", a.AlgorithmSigningComponent, a.Signature)
	// Output:address: 134a6TXxzgQ9Az3w8BcvgdZyA5UqRL89da signature: H+lubfcz5Z2oG8B7HwmP8Z+tALP+KNOPgedo7UTXwW8LBpMkgCgatCdpvbtf7wZZQSIMz83emmAvVS4S3F5X1wo=
}

// BenchmarkNewFromTape benchmarks the method NewFromTape()
func BenchmarkNewFromTape(b *testing.B) {
	bobValidData, _ := bob.NewFromString(sampleValidBobTx)
	for i := 0; i < b.N; i++ {
		_ = NewFromTape(bobValidData.Out[0].Tape[2])
	}
}

// TestNewFromTapes will test the method NewFromTapes()
func TestNewFromTapes(t *testing.T) {
	t.Parallel()

	// Parse from string into BOB
	bobValidData, err := bob.NewFromString(sampleValidBobTx)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}
	var bobInvalidData *bob.Tx
	if bobInvalidData, err = bob.NewFromString(sampleInvalidBobTx); err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	var (
		// Testing private methods
		tests = []struct {
			inputTapes         []bob.Tape
			expectedSignature  string
			expectedAlgorithm  Algorithm
			expectedComponent  string
			expectedData0      string
			expectedData1      string
			expectedData2      string
			expectedNil        bool
			expectedValidation bool
		}{
			{
				bobValidData.Out[0].Tape,
				"H+lubfcz5Z2oG8B7HwmP8Z+tALP+KNOPgedo7UTXwW8LBpMkgCgatCdpvbtf7wZZQSIMz83emmAvVS4S3F5X1wo=",
				BitcoinECDSA,
				"134a6TXxzgQ9Az3w8BcvgdZyA5UqRL89da",
				opReturn,
				"1BAPSuaPnfGnSBM3GLV9yhxUdYe4vGbdMT",
				"ATTEST",
				false,
				true,
			},
			{
				bobInvalidData.Out[0].Tape,
				"H+lubfcz5Z2oG8B7HwmP8Z+tALP+KNOPgedo7UTXwW8LBpMkgCgatCdpvbtf7wZZQSIMz83emmAvVS4S3F5X1wo=",
				BitcoinECDSA,
				"invalid-address",
				opReturn,
				"1BAPSuaPnfGnSBM3GLV9yhxUdYe4vGbdMT",
				"ATTEST",
				false,
				false,
			},
			{
				[]bob.Tape{*new(bob.Tape)},
				"",
				"",
				"",
				"",
				"",
				"",
				true,
				false,
			},
		}
	)

	// Run tests
	for _, test := range tests {
		if a := NewFromTapes(test.inputTapes); a == nil && !test.expectedNil {
			t.Errorf("%s Failed: [%v] inputted and nil was not expected (aip)", t.Name(), test.inputTapes)
		} else if a != nil && test.expectedNil {
			t.Errorf("%s Failed: [%v] inputted and nil was expected (aip)", t.Name(), test.inputTapes)
		} else if a != nil && a.Signature != test.expectedSignature {
			t.Errorf("%s Failed: [%v] inputted and expected [%s] but got [%s]", t.Name(), test.inputTapes, test.expectedSignature, a.Signature)
		} else if a != nil && a.Algorithm != test.expectedAlgorithm {
			t.Errorf("%s Failed: [%v] inputted and expected [%s] but got [%s]", t.Name(), test.inputTapes, test.expectedAlgorithm, a.Algorithm)
		} else if a != nil && a.AlgorithmSigningComponent != test.expectedComponent {
			t.Errorf("%s Failed: [%v] inputted and expected [%s] but got [%s]", t.Name(), test.inputTapes, test.expectedComponent, a.AlgorithmSigningComponent)
		} else if a != nil && len(a.Data) > 0 && a.Data[0] != test.expectedData0 {
			t.Errorf("%s Failed: [%v] inputted and expected [%s] but got [%s]", t.Name(), test.inputTapes, test.expectedData0, a.Data[0])
		} else if a != nil && len(a.Data) > 0 && a.Data[1] != test.expectedData1 {
			t.Errorf("%s Failed: [%v] inputted and expected [%s] but got [%s]", t.Name(), test.inputTapes, test.expectedData1, a.Data[1])
		} else if a != nil && len(a.Data) > 0 && a.Data[2] != test.expectedData2 {
			t.Errorf("%s Failed: [%v] inputted and expected [%s] but got [%s]", t.Name(), test.inputTapes, test.expectedData2, a.Data[2])
		} else if a != nil && len(test.inputTapes) > 1 {
			var valid bool
			valid, err = ValidateTapes(test.inputTapes)
			if valid && !test.expectedValidation {
				t.Errorf("%s Failed: [%v] inputted and validation should have failed", t.Name(), test.inputTapes)
			} else if !valid && test.expectedValidation && err != nil {
				t.Errorf("%s Failed: [%v] inputted and validation should have passed, error: %s", t.Name(), test.inputTapes, err.Error())
			}
		}
	}
}

// ExampleNewFromTapes example using NewFromTapes()
func ExampleNewFromTapes() {
	// Get BOB data from a TX
	bobValidData, err := bob.NewFromString(sampleValidBobTx)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Get from tape given the AIP index
	a := NewFromTapes(bobValidData.Out[0].Tape)

	fmt.Printf("address: %s signature: %s", a.AlgorithmSigningComponent, a.Signature)
	// Output:address: 134a6TXxzgQ9Az3w8BcvgdZyA5UqRL89da signature: H+lubfcz5Z2oG8B7HwmP8Z+tALP+KNOPgedo7UTXwW8LBpMkgCgatCdpvbtf7wZZQSIMz83emmAvVS4S3F5X1wo=
}

// BenchmarkNewFromTapes benchmarks the method NewFromTapes()
func BenchmarkNewFromTapes(b *testing.B) {
	bobValidData, _ := bob.NewFromString(sampleValidBobTx)
	for i := 0; i < b.N; i++ {
		_ = NewFromTapes(bobValidData.Out[0].Tape)
	}
}

// TestValidateTapes will test the method ValidateTapes()
func TestValidateTapes(t *testing.T) {
	t.Parallel()

	// Parse from string into BOB
	bobValidData, err := bob.NewFromString(sampleValidBobTx)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}
	var bobInvalidData *bob.Tx
	if bobInvalidData, err = bob.NewFromString(sampleInvalidBobTx); err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	var (
		// Testing private methods
		tests = []struct {
			inputTapes         []bob.Tape
			expectedValidation bool
		}{
			{
				bobValidData.Out[0].Tape,
				true,
			},
			{
				bobInvalidData.Out[0].Tape,
				false,
			},
			{
				[]bob.Tape{*new(bob.Tape)},
				false,
			},
		}
	)

	// Run tests
	for idx, test := range tests {
		if valid, err := ValidateTapes(test.inputTapes); valid && !test.expectedValidation {
			t.Errorf("%d %s Failed: inputted and validation should have failed", idx, t.Name())
		} else if !valid && test.expectedValidation && err != nil {
			t.Errorf("%d %s Failed: inputted and validation should have passed, error: %s", idx, t.Name(), err.Error())
		}
	}
}

// ExampleValidateTapes example using ValidateTapes()
func ExampleValidateTapes() {
	// Get BOB data from a TX
	bobValidData, err := bob.NewFromString(sampleValidBobTx)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Get from tape
	if valid, err := ValidateTapes(bobValidData.Out[0].Tape); valid {
		fmt.Print("AIP is valid")
	} else if err != nil {
		fmt.Printf("AIP is invalid: %s", err.Error())
	}
	// Output:AIP is valid
}

// BenchmarkValidateTapes benchmarks the method ValidateTapes()
func BenchmarkValidateTapes(b *testing.B) {
	bobValidData, _ := bob.NewFromString(sampleValidBobTx)
	for i := 0; i < b.N; i++ {
		_, _ = ValidateTapes(bobValidData.Out[0].Tape)
	}
}

// getBobOutput helper to get op_return in BOB format
func getBobOutput() bob.Output {

	// Create op_return
	opReturn := bitcoin.OpReturnData{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}}

	// Create a transaction
	privateKey, _ := bitcoin.PrivateKeyFromString(examplePrivateKey)
	tx, _ := bitcoin.CreateTx(nil, nil, []bitcoin.OpReturnData{opReturn}, privateKey)

	// Create the bob tx from hex
	bobTx, _ := bob.NewFromRawTxString(tx.ToString())

	return bobTx.Out[0]
}

// TestSignBobOpReturnData tests for nil case in SignBobOpReturnData()
func TestSignBobOpReturnData(t *testing.T) {
	t.Parallel()

	var (
		// Testing private methods
		tests = []struct {
			inputPrivateKey    string
			inputAlgorithm     Algorithm
			inputData          bob.Output
			expectedSignature  string
			expectedAipNil     bool
			expectedOutNil     bool
			expectedError      bool
			expectedValidation bool
		}{
			{
				"80699541455b59a8a8a33b85892319de8b8e8944eb8b48e9467137825ae192e59f01",
				BitcoinECDSA,
				getBobOutput(),
				"G/tmf2aRfTqri7pCF793/xDZO2COIzcY2BRpb5P0o+zGIG0tCF3JJLadQwCCC+Lu+Xcanv+Fl82lrk3hVlo8bXY=",
				false,
				false,
				false,
				true,
			},
			{
				"",
				BitcoinECDSA,
				getBobOutput(),
				"",
				true,
				true,
				true,
				false,
			},
			{
				"80699541455b59a8a8a33b85892319de8b8e8944eb8b48e9467137825ae192e59f01",
				Paymail,
				getBobOutput(),
				"G/tmf2aRfTqri7pCF793/xDZO2COIzcY2BRpb5P0o+zGIG0tCF3JJLadQwCCC+Lu+Xcanv+Fl82lrk3hVlo8bXY=",
				false,
				false,
				false,
				true,
			},
		}
	)

	// Run tests
	for _, test := range tests {
		if out, a, err := SignBobOpReturnData(test.inputPrivateKey, test.inputAlgorithm, test.inputData); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] [%s] [%v] inputted and error not expected but got: %s", t.Name(), test.inputPrivateKey, test.inputAlgorithm, test.inputData, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%s] [%s] [%v] inputted and error was expected", t.Name(), test.inputPrivateKey, test.inputAlgorithm, test.inputData)
		} else if a == nil && !test.expectedAipNil {
			t.Errorf("%s Failed: [%s] [%s] [%v] inputted and nil was not expected (aip)", t.Name(), test.inputPrivateKey, test.inputAlgorithm, test.inputData)
		} else if a != nil && test.expectedAipNil {
			t.Errorf("%s Failed: [%s] [%s] [%v] inputted and nil was expected (aip)", t.Name(), test.inputPrivateKey, test.inputAlgorithm, test.inputData)
		} else if out == nil && !test.expectedOutNil {
			t.Errorf("%s Failed: [%s] [%s] [%v] inputted and nil was not expected (out)", t.Name(), test.inputPrivateKey, test.inputAlgorithm, test.inputData)
		} else if out != nil && test.expectedOutNil {
			t.Errorf("%s Failed: [%s] [%s] [%v] inputted and nil was expected (out)", t.Name(), test.inputPrivateKey, test.inputAlgorithm, test.inputData)
		} else if a != nil && a.Signature != test.expectedSignature {
			t.Errorf("%s Failed: [%s] [%s] [%v] inputted and expected signature [%s] but got [%s]", t.Name(), test.inputPrivateKey, test.inputAlgorithm, test.inputData, test.expectedSignature, a.Signature)
		} else if a != nil {
			var valid bool
			if valid, err = a.Validate(); valid && !test.expectedValidation {
				t.Errorf("%s Failed: [%s] [%s] [%v] inputted and validation should have failed", t.Name(), test.inputPrivateKey, test.inputAlgorithm, test.inputData)
			} else if !valid && test.expectedValidation && err != nil {
				t.Errorf("%s Failed: [%s] [%s] [%v] inputted and validation should have passed, error: %s", t.Name(), test.inputPrivateKey, test.inputAlgorithm, test.inputData, err.Error())
			}
		}
	}
}

// ExampleSignBobOpReturnData example using SignBobOpReturnData()
func ExampleSignBobOpReturnData() {
	_, a, err := SignBobOpReturnData(examplePrivateKey, BitcoinECDSA, getBobOutput())
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("signature: %s", a.Signature)
	// Output:signature: G1oBzIBhluBmdu2U6xfb75ACIPpDVmjCV4OV1hVAlnovG1IN3jyIOYK/HJQuH5UOjf2rfaI45SZqPxx9llN+Mgs=
}

// BenchmarkSignBobOpReturnData benchmarks the method SignBobOpReturnData()
func BenchmarkSignBobOpReturnData(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, _ = SignBobOpReturnData(examplePrivateKey, BitcoinECDSA, getBobOutput())
	}
}

// todo: test a TX with signature of specific indexes
