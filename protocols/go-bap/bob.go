package bap

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/bitcoinschema/go-bob"
)

// Bap is BAP data object from the bob.Tape
type Bap struct {
	Address  string          `json:"address,omitempty" bson:"address,omitempty"`
	IDKey    string          `json:"id_key,omitempty" bson:"id_key,omitempty"`
	Sequence uint64          `json:"sequence" bson:"sequence"`
	Type     AttestationType `json:"type,omitempty" bson:"type,omitempty"`
	URNHash  string          `json:"urn_hash,omitempty" bson:"urn_hash,omitempty"`
}

// FromTape takes a bob.Tape and returns a BAP data structure
func (b *Bap) FromTape(tape *bob.Tape) (err error) {
	b.Type = AttestationType(tape.Cell[1].S)

	// Invalid length
	if len(tape.Cell) < 4 {
		err = fmt.Errorf("invalid %s record %+v", b.Type, tape.Cell)
		return
	}

	switch b.Type {
	case REVOKE, ATTEST:
		b.URNHash = tape.Cell[2].S
		if b.Sequence, err = strconv.ParseUint(tape.Cell[3].S, 10, 64); err != nil {
			return err
		}
	case ID:
		b.Address = tape.Cell[3].S
		b.IDKey = tape.Cell[2].S
	}
	return
}

// NewFromTapes will create a new BAP object from a []bob.Tape
func NewFromTapes(tapes []bob.Tape) (*Bap, error) {
	// Loop tapes -> cells (only supporting 1 BAP record right now)
	for index, t := range tapes {
		for _, cell := range t.Cell {
			if cell.S == Prefix {
				return NewFromTape(&tapes[index])
			}
		}
	}
	return nil, errors.New("no BAP record found")
}

// NewFromTape takes a bob.Tape and returns a BAP data structure
func NewFromTape(tape *bob.Tape) (b *Bap, err error) {
	b = new(Bap)
	err = b.FromTape(tape)
	return
}
