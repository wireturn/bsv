package aip

import (
	"encoding/hex"
	"errors"
	"strconv"
	"strings"

	"github.com/bitcoinschema/go-bob"
)

// NewFromTape will create a new AIP object from a bob.Tape
// Using the FromTape() alone will prevent validation (data is needed via SetData to enable)
func NewFromTape(tape bob.Tape) (a *Aip) {
	a = new(Aip)
	a.FromTape(tape)
	return
}

// FromTape takes a BOB Tape and returns an Aip data structure.
// Using the FromTape() alone will prevent validation (data is needed via SetData to enable)
func (a *Aip) FromTape(tape bob.Tape) {

	// Not a valid tape?
	if len(tape.Cell) < 4 {
		return
	}

	// Loop to find start of AIP
	var startIndex int
	for i, cell := range tape.Cell {
		if cell.S == Prefix {
			startIndex = i
			break
		}
	}

	// Set the AIP fields
	a.Algorithm = Algorithm(tape.Cell[startIndex+1].S)
	a.AlgorithmSigningComponent = tape.Cell[startIndex+2].S
	a.Signature = tape.Cell[startIndex+3].B

	// Final index count
	finalIndexCount := startIndex + 4

	// Store the indices
	if len(tape.Cell) > finalIndexCount {

		// TODO: Consider OP_RETURN is included in sig when processing a tx using indices
		// Loop over remaining indices if they exist and append to indices slice
		a.Indices = make([]int, len(tape.Cell)-finalIndexCount)
		for x := finalIndexCount - 1; x < len(tape.Cell); x++ {
			if index, err := strconv.Atoi(tape.Cell[x].S); err == nil {
				a.Indices = append(a.Indices, index)
			}
		}
	}
}

// NewFromTapes will create a new AIP object from a []bob.Tape
// Using the FromTapes() alone will prevent validation (data is needed via SetData to enable)
func NewFromTapes(tapes []bob.Tape) (a *Aip) {
	// Loop tapes -> cells (only supporting 1 sig right now)
	for _, t := range tapes {
		for _, cell := range t.Cell {
			if cell.S == Prefix {
				a = new(Aip)
				a.FromTape(t)
				a.SetDataFromTapes(tapes)
				return
			}
		}
	}
	return
}

// SetDataFromTapes sets the data the AIP signature is signing
func (a *Aip) SetDataFromTapes(tapes []bob.Tape) {

	// Set OP_RETURN to be consistent with BitcoinFiles SDK
	var data = []string{opReturn}

	if len(a.Indices) == 0 {

		// Walk over all output values and concatenate them until we hit the AIP prefix, then add in the separator
		for _, tape := range tapes {
			for _, cell := range tape.Cell {
				if cell.S != Prefix {
					// Skip the OPS
					if cell.Ops != "" {
						continue
					}
					data = append(data, strings.TrimSpace(cell.S))
				} else {
					data = append(data, pipe)
					a.Data = data
					return
				}
			}
		}

	} else {

		var indexCt = 0

		for _, tape := range tapes {
			for _, cell := range tape.Cell {
				if cell.S != Prefix && contains(a.Indices, indexCt) {
					data = append(data, cell.S)
				} else {
					data = append(data, pipe)
				}
				indexCt++
			}
		}

		a.Data = data
	}
}

// SignBobOpReturnData appends a signature to a BOB Tx by adding a
// protocol separator followed by AIP information
func SignBobOpReturnData(privateKey string, algorithm Algorithm, output bob.Output) (*bob.Output, *Aip, error) {

	// Parse the data to sign
	var dataToSign []string
	for _, tape := range output.Tape {
		for _, cell := range tape.Cell {
			if len(cell.S) > 0 {
				dataToSign = append(dataToSign, cell.S)
			} else {
				// TODO: Review this case. Should we assume the b64 is signed?
				//  Should protocol doc for AIP mention this?
				dataToSign = append(dataToSign, cell.B)
			}
		}
	}

	// Sign the data
	a, err := Sign(privateKey, algorithm, strings.Join(dataToSign, ""))
	if err != nil {
		return nil, nil, err
	}

	// Create the output tape
	output.Tape = append(output.Tape, bob.Tape{
		Cell: []bob.Cell{{
			H: hex.EncodeToString([]byte(Prefix)),
			S: Prefix,
		}, {
			H: hex.EncodeToString([]byte(algorithm)),
			S: string(algorithm),
		}, {
			H: hex.EncodeToString([]byte(a.AlgorithmSigningComponent)),
			S: a.AlgorithmSigningComponent,
		}, {
			H: hex.EncodeToString([]byte(a.Signature)),
			S: a.Signature,
		}},
	})

	return &output, a, nil
}

// ValidateTapes validates the AIP signature for a given []bob.Tape
func ValidateTapes(tapes []bob.Tape) (bool, error) {
	// Loop tapes -> cells (only supporting 1 sig right now)
	for _, tape := range tapes {
		for _, cell := range tape.Cell {

			// Once we hit AIP Prefix, stop
			if cell.S == Prefix {
				a := NewFromTape(tape)
				a.SetDataFromTapes(tapes)
				return a.Validate()
			}
		}

	}
	return false, errors.New("no AIP tape found")
}

// contains looks in a slice for a given value
func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
