package bap

import (
	"fmt"
	"testing"

	"github.com/bitcoinschema/go-bob"
)

// TestFromTape will test the method NewFromTape()
func TestNewFromTape(t *testing.T) {

	// Get BOB data from string
	bobData, err := bob.NewFromString(sampleValidBobTx)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	// Get from tape
	var b *Bap
	b, err = NewFromTape(&bobData.Out[0].Tape[1])
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	} else if b.Type != ATTEST {
		t.Fatalf("expected: %s got: %s", ATTEST, b.Type)
	}

	// Wrong tape
	_, err = NewFromTape(&bobData.Out[0].Tape[0])
	if err == nil {
		t.Fatalf("error should have occurred")
	}

	// Revoke
	bobData.Out[0].Tape[1].Cell[1].S = string(REVOKE)
	_, err = NewFromTape(&bobData.Out[0].Tape[1])
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	// Bad Sequence
	bobData.Out[0].Tape[1].Cell[3].S = ""
	_, err = NewFromTape(&bobData.Out[0].Tape[1])
	if err == nil {
		t.Fatalf("error should have occurred")
	}

	// ID tape
	bobData.Out[0].Tape[1].Cell[1].S = string(ID)
	bobData.Out[0].Tape[1].Cell[2].S = "idKey"
	bobData.Out[0].Tape[1].Cell[3].S = "Address"
	_, err = NewFromTape(&bobData.Out[0].Tape[1])
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}
}

// ExampleNewFromTape example using NewFromTape()
func ExampleNewFromTape() {

	// Get BOB data from string
	bobData, err := bob.NewFromString(sampleValidBobTx)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Get from tape
	var b *Bap
	b, err = NewFromTape(&bobData.Out[0].Tape[1])
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("BAP type: %s", b.Type)
	// Output:BAP type: ATTEST
}

// BenchmarkNewFromTape benchmarks the method NewFromTape()
func BenchmarkNewFromTape(b *testing.B) {
	bobData, _ := bob.NewFromString(sampleValidBobTx)
	for i := 0; i < b.N; i++ {
		_, _ = NewFromTape(&bobData.Out[0].Tape[1])
	}
}

// TestFromTapePanic tests for nil case in NewFromTape()
func TestNewFromTapePanic(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("the code did not panic")
		}
	}()

	_, err := NewFromTape(nil)
	if err == nil {
		t.Fatalf("error expected")
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
			inputTapes       []bob.Tape
			expectedType     AttestationType
			expectedSequence uint64
			expectedURNHash  string
			expectedNil      bool
			expectedError    bool
		}{
			{
				bobValidData.Out[0].Tape,
				"ATTEST",
				0,
				"cf39fc55da24dc23eff1809e6e6cf32a0fe6aecc81296543e9ac84b8c501bac5",
				false,
				false,
			},
			{
				bobInvalidData.Out[0].Tape,
				"",
				0,
				"",
				true,
				true,
			},
		}
	)

	// Run tests
	var b *Bap
	for _, test := range tests {
		if b, err = NewFromTapes(test.inputTapes); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%v] inputted and error not expected but got: %s", t.Name(), test.inputTapes, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%v] inputted and error was expected", t.Name(), test.inputTapes)
		} else if b == nil && !test.expectedNil {
			t.Errorf("%s Failed: [%v] inputted and nil was not expected (bap)", t.Name(), test.inputTapes)
		} else if b != nil && test.expectedNil {
			t.Errorf("%s Failed: [%v] inputted and nil was expected (bap)", t.Name(), test.inputTapes)
		} else if b != nil && b.Type != test.expectedType {
			t.Errorf("%s Failed: [%v] inputted and expected [%s] but got [%s]", t.Name(), test.inputTapes, test.expectedType, b.Type)
		} else if b != nil && b.Sequence != test.expectedSequence {
			t.Errorf("%s Failed: [%v] inputted and expected [%d] but got [%d]", t.Name(), test.inputTapes, test.expectedSequence, b.Sequence)
		} else if b != nil && b.URNHash != test.expectedURNHash {
			t.Errorf("%s Failed: [%v] inputted and expected [%s] but got [%s]", t.Name(), test.inputTapes, test.expectedURNHash, b.URNHash)
		}
	}
}

// ExampleNewFromTapes example using NewFromTapes()
func ExampleNewFromTapes() {

	// Get BOB data from string
	bobData, err := bob.NewFromString(sampleValidBobTx)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Get from tapes
	var b *Bap
	b, err = NewFromTapes(bobData.Out[0].Tape)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("BAP type: %s", b.Type)
	// Output:BAP type: ATTEST
}

// BenchmarkNewFromTapes benchmarks the method NewFromTapes()
func BenchmarkNewFromTapes(b *testing.B) {
	bobData, _ := bob.NewFromString(sampleValidBobTx)
	for i := 0; i < b.N; i++ {
		_, _ = NewFromTapes(bobData.Out[0].Tape)
	}
}
