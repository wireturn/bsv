package paymail

import (
	"fmt"
	"testing"
)

/*
	Test cases from: http://bsvalias.org/01-02-brfc-id-assignment.html
*/

// TestBRFCSpec_Generate will test the method Generate()
func TestBRFCSpec_Generate(t *testing.T) {

	t.Parallel()

	var tests = []struct {
		brfc          *BRFCSpec
		expectedID    string
		expectedError bool
	}{
		// Test Case #1 from: http://bsvalias.org/01-02-brfc-id-assignment.html
		{&BRFCSpec{Author: "andy (nChain)", ID: "57dd1f54fc67", Title: "BRFC Specifications", Version: "1"}, "57dd1f54fc67", false},
		// Test Case #2 from: http://bsvalias.org/01-02-brfc-id-assignment.html
		{&BRFCSpec{Author: "andy (nChain)", ID: "74524c4d6274", Title: "bsvalias Payment Addressing (PayTo Protocol Prefix)", Version: "1"}, "74524c4d6274", false},
		// Test Case #3 from: http://bsvalias.org/01-02-brfc-id-assignment.html
		{&BRFCSpec{Author: "andy (nChain)", ID: "0036f9b8860f", Title: "bsvalias Integration with Simplified Payment Protocol", Version: "1"}, "0036f9b8860f", false},
		// Error cases:
		{&BRFCSpec{Author: "andy (nChain)", ID: "12345", Title: "", Version: "1"}, "", true},
		{&BRFCSpec{Author: "", ID: "12345", Title: "", Version: "1"}, "", true},
		{&BRFCSpec{Author: "", ID: "", Title: "", Version: "1"}, "", true},
		{&BRFCSpec{Author: "", ID: "", Title: "", Version: ""}, "", true},
		{&BRFCSpec{Author: "  andy (nChain)  ", ID: "0036f9b8860f", Title: "  bsvalias Integration with Simplified Payment Protocol  ", Version: "1"}, "0036f9b8860f", false},
	}

	for _, test := range tests {
		if err := test.brfc.Generate(); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%v] inputted, [%s] expected and error not expected but got: %s", t.Name(), test.brfc, test.expectedID, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%v] inputted, [%s] expected and error was expected", t.Name(), test.brfc, test.expectedID)
		} else if test.brfc.ID != test.expectedID {
			t.Errorf("%s Failed: [%v] inputted, [%s] expected and id did not match, got: %s", t.Name(), test.brfc, test.expectedID, test.brfc.ID)
		}
	}
}

// ExampleBRFCSpec_Generate example using Generate()
//
// See more examples in /examples/
func ExampleBRFCSpec_Generate() {
	// Start with a new BRFC specification
	newBRFC := &BRFCSpec{
		Author:  "MrZ",
		Title:   "New BRFC",
		Version: "1",
	}
	if err := newBRFC.Generate(); err != nil {
		fmt.Printf("error occurred: %s", err.Error())
	} else {
		fmt.Printf("id generated: %s", newBRFC.ID)
	}
	// Output:id generated: e898079d7d1a
}

// BenchmarkBRFCSpec_Generate benchmarks the method Generate()
func BenchmarkBRFCSpec_Generate(b *testing.B) {
	newBRFC := &BRFCSpec{Author: "MrZ", Title: "New BRFC", Version: "1"}
	for i := 0; i < b.N; i++ {
		_ = newBRFC.Generate()
	}
}

// TestBRFCSpec_Validate will test the method Validate()
func TestBRFCSpec_Validate(t *testing.T) {

	t.Parallel()

	var tests = []struct {
		brfc          *BRFCSpec
		expectedID    string
		expectedError bool
		expectedValid bool
	}{
		// Test Case #1 from: http://bsvalias.org/01-02-brfc-id-assignment.html
		{&BRFCSpec{Author: "andy (nChain)", ID: "57dd1f54fc67", Title: "BRFC Specifications", Version: "1"}, "57dd1f54fc67", false, true},
		// Test Case #2 from: http://bsvalias.org/01-02-brfc-id-assignment.html
		{&BRFCSpec{Author: "andy (nChain)", ID: "74524c4d6274", Title: "bsvalias Payment Addressing (PayTo Protocol Prefix)", Version: "1"}, "74524c4d6274", false, true},
		// Test Case #3 from: http://bsvalias.org/01-02-brfc-id-assignment.html
		{&BRFCSpec{Author: "andy (nChain)", ID: "0036f9b8860f", Title: "bsvalias Integration with Simplified Payment Protocol", Version: "1"}, "0036f9b8860f", false, true},
		// Error cases:
		{&BRFCSpec{Author: "andy (nChain)", ID: "12345", Title: "", Version: "1"}, "", true, false},
		{&BRFCSpec{Author: "", ID: "12345", Title: "", Version: "1"}, "", true, false},
		{&BRFCSpec{Author: "", ID: "", Title: "", Version: "1"}, "", true, false},
		{&BRFCSpec{Author: "", ID: "", Title: "", Version: ""}, "", true, false},
		{&BRFCSpec{Author: "  andy (nChain)  ", ID: "0036f9b8860f", Title: "  bsvalias Integration with Simplified Payment Protocol  ", Version: "1"}, "0036f9b8860f", false, true},
		{&BRFCSpec{Author: "andy (nChain)", ID: "0036f9b8860z", Title: "  bsvalias Integration with Simplified Payment Protocol  ", Version: "1"}, "0036f9b8860f", false, false},
		{&BRFCSpec{Author: "nChain", ID: "a0a4c8b96133", Title: "spv_channels", Version: "1.0.0-beta"}, "a0a4c8b96133", false, true},
		{&BRFCSpec{Author: "nChain", ID: "b8930c2bbf5d", Title: "minerIdExt-blockBind", Version: "0.1"}, "b8930c2bbf5d", false, true},
		{&BRFCSpec{Author: "nChain", ID: "a224052ad433", Title: "minerIdExt-blockInfo", Version: "0.1"}, "a224052ad433", false, true},
		{&BRFCSpec{Author: "nChain", ID: "1b1d980b5b72", Title: "minerIdExt-minerParams", Version: "0.1"}, "1b1d980b5b72", false, true},
		{&BRFCSpec{Author: "nChain", ID: "298e080a4598", Title: "jsonEnvelope", Version: "0.1"}, "298e080a4598", false, true},
		{&BRFCSpec{Author: "nChain", ID: "62b21572ca46", Title: "minerIdExt-feeSpec", Version: "0.1"}, "62b21572ca46", false, true},
		{&BRFCSpec{Author: "nChain", ID: "fb567267440a", Title: "feeSpec", Version: "0.1"}, "fb567267440a", false, true},
		{&BRFCSpec{Author: "nChain", ID: "07f0786cdab6", Title: "minerId", Version: "0.1"}, "07f0786cdab6", false, true},
		{&BRFCSpec{Author: "nChain", ID: "ce852c4c2cd1", Title: "merchant_api", Version: "0.1"}, "eaad81dc6d4d", false, false},
		{&BRFCSpec{Author: "Fabriik", ID: "1300361cb2d4", Title: "Asset Information", Version: "1"}, "1300361cb2d4", false, true},
		{&BRFCSpec{Author: "Fabriik", ID: "189e32d93d28", Title: "Simple Fabriik Protocol for Tokens Build Action", Version: "1"}, "189e32d93d28", false, true},
		{&BRFCSpec{Author: "Fabriik", ID: "95dddb461bff", Title: "Simple Fabriik Protocol for Tokens Authorise Action", Version: "1"}, "95dddb461bff", false, true},
		{&BRFCSpec{Author: "Fabriik", ID: "f792b6eff07a", Title: "P2P Payment Destination with Tokens Support", Version: "1"}, "f792b6eff07a", false, true},
	}

	for _, test := range tests {
		if valid, id, err := test.brfc.Validate(); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%v] inputted, [%s] expected and error not expected but got: %s", t.Name(), test.brfc, test.expectedID, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%v] inputted, [%s] expected and error was expected", t.Name(), test.brfc, test.expectedID)
		} else if id != test.expectedID {
			t.Errorf("%s Failed: [%v] inputted, [%s] expected and id did not match, got: %s", t.Name(), test.brfc, test.expectedID, id)
		} else if valid != test.expectedValid || test.brfc.Valid != test.expectedValid {
			t.Errorf("%s Failed: [%v] inputted, [%s] expected and valid did not match", t.Name(), test.brfc, test.expectedID)
		}
	}
}

// ExampleBRFCSpec_Validate example using Validate()
//
// See more examples in /examples/
func ExampleBRFCSpec_Validate() {
	// Start with an existing BRFC specification
	newBRFC := &BRFCSpec{
		Author:  "MrZ",
		ID:      "e898079d7d1a",
		Title:   "New BRFC",
		Version: "1",
	}
	if valid, id, err := newBRFC.Validate(); err != nil {
		fmt.Printf("error occurred: %s", err.Error())
	} else if !valid {
		fmt.Printf("id is invalid: %s vs %s", newBRFC.ID, id)
	} else {
		fmt.Printf("brfc is valid: %s", id)
	}
	// Output:brfc is valid: e898079d7d1a
}

// BenchmarkBRFCSpec_Validate benchmarks the method Validate()
func BenchmarkBRFCSpec_Validate(b *testing.B) {
	newBRFC := &BRFCSpec{Author: "MrZ", ID: "e898079d7d1a", Title: "New BRFC", Version: "1"}
	for i := 0; i < b.N; i++ {
		_, _, _ = newBRFC.Validate()
	}
}

// TestClientOptions_LoadBRFCs will test the method LoadBRFCs()
func TestClientOptions_LoadBRFCs(t *testing.T) {
	// t.Parallel() cannot use newTestClient() race condition

	// Create a client with options
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	var tests = []struct {
		specJSON       string
		expectedLength int
		expectedError  bool
	}{
		{`[{"author": "andy (nChain)","id": "57dd1f54fc67","title": "BRFC Specifications","url": "http://bsvalias.org/01-02-brfc-id-assignment.html","version": "1"}]`, len(client.options.brfcSpecs) + 1, false},
		{`[{"invalid:1}]`, len(client.options.brfcSpecs), true},
		{`[{"author": "andy (nChain), Ryan X. Charles (Money Button)","title":"invalid-spec","id": "17dd1f54fc66"}]`, len(client.options.brfcSpecs), true},
		{`[{"author": "andy (nChain), Ryan X. Charles (Money Button)","title":""}]`, len(client.options.brfcSpecs), true},
	}

	for _, test := range tests {
		if err = client.options.LoadBRFCs(test.specJSON); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted, [%d] expected specs and error not expected but got: %s", t.Name(), test.specJSON, test.expectedLength, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%s] inputted, [%d] expected specs and error was expected", t.Name(), test.specJSON, test.expectedLength)
		} else if len(client.options.brfcSpecs) != test.expectedLength {
			t.Errorf("%s Failed: [%s] inputted, [%d] expected specs but got: %d", t.Name(), test.specJSON, test.expectedLength, len(client.options.brfcSpecs))
		}
	}
}

// ExampleClientOptions_LoadBRFCs example using LoadBRFCs()
//
// See more examples in /examples/
// nolint:govet // options are now private but example still useful.
func ExampleClientOptions_LoadBRFCs() {
	// Create a client with options
	client, err := NewClient()
	if err != nil {
		fmt.Printf("error loading client: %s", err.Error())
		return
	}

	// Load additional specification(s)
	additionalSpec := `[{"author": "andy (nChain)","id": "57dd1f54fc67","title": "BRFC Specifications","url": "http://bsvalias.org/01-02-brfc-id-assignment.html","version": "1"}]`
	if err = client.options.LoadBRFCs(additionalSpec); err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("total specifications found: %d", len(client.options.brfcSpecs))

	// Output:total specifications found: 24
}

// BenchmarkClientOptions_LoadBRFCs benchmarks the method LoadBRFCs()
func BenchmarkClientOptions_LoadBRFCs(b *testing.B) {
	client, _ := NewClient(nil, nil, nil)
	additionalSpec := `[{"author": "andy (nChain)","id": "57dd1f54fc67","title": "BRFC Specifications","url": "http://bsvalias.org/01-02-brfc-id-assignment.html","version": "1"}]`
	for i := 0; i < b.N; i++ {
		_ = client.options.LoadBRFCs(additionalSpec)
	}
}
