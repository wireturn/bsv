package b

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/bitcoinschema/go-bob"
)

// NewFromTape will create a new AIP object from a bob.Tape
// Using the FromTape() alone will prevent validation (data is needed via SetData to enable)
func NewFromTape(tape bob.Tape) (b *B, err error) {
	b = new(B)
	err = b.FromTape(tape)
	return
}

// NewFromTapes will create a new B object from a []bob.Tape
// Using the FromTapes() alone will prevent validation (data is needed via SetData to enable)
func NewFromTapes(tapes []bob.Tape) (b *B, err error) {
	// Loop tapes -> cells (only supporting 1 sig right now)
	for _, t := range tapes {
		for _, cell := range t.Cell {
			if cell.S == Prefix {
				b = new(B)
				err = b.FromTape(t)
				// b.SetDataFromTapes(tapes)
				return
			}
		}
	}
	err = errors.New("no b tape found")
	return
}

// todo: SetDataFromTapes()

// FromTape takes a BOB Tape and returns a B data structure
func (b *B) FromTape(tape bob.Tape) (err error) {
	if len(tape.Cell) < 4 {
		err = fmt.Errorf("invalid B tx Only %d pushdatas", len(tape.Cell))
		return
	}

	// Loop to find start of B
	var startIndex int
	for i, cell := range tape.Cell {
		if cell.S == Prefix {
			startIndex = i
			break
		}
	}

	// Media type is after data
	b.MediaType = tape.Cell[startIndex+2].S

	// Encoding is after media
	b.Encoding = tape.Cell[startIndex+3].S

	switch EncodingType(strings.ToLower(b.Encoding)) {
	case EncodingGzip:
		fallthrough
	case EncodingBinary:
		// Decode base64 data
		if b.Data.Bytes, err = base64.StdEncoding.DecodeString(tape.Cell[startIndex+1].B); err != nil {
			return
		}
	case EncodingUtf8:
		fallthrough
	case EncodingUtf8Alt:
		b.Data.UTF8 = tape.Cell[startIndex+1].S
	}

	// Filename is optional and last
	if len(tape.Cell) > startIndex+4 && len(tape.Cell[startIndex+4].S) != 0 {
		b.Filename = tape.Cell[startIndex+4].S
	}

	return
}
