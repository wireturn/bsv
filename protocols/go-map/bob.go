package magic

import (
	"fmt"

	"github.com/bitcoinschema/go-bob"
)

// NewFromTape takes a tape and returns a new MAP
func NewFromTape(tape *bob.Tape) (magicTx MAP, err error) {
	magicTx = make(MAP)
	err = magicTx.FromTape(tape)
	return
}

// FromTape sets a MAP object from a BOB Tape
func (m MAP) FromTape(tape *bob.Tape) error {

	if len(tape.Cell) < 3 {
		return fmt.Errorf("invalid MAP record - missing required parameters %d", len(tape.Cell))
	}

	if tape.Cell[0].S == Prefix {
		m[Cmd] = tape.Cell[1].S

		switch m[Cmd] {
		case Delete:
			m.delete(tape.Cell)
		case Add:
			m.add(tape.Cell)
		case Remove:
			m.remove(tape.Cell)
		case Set:
			m.set(tape.Cell)
		case Select:
			if len(tape.Cell) < 5 {
				return fmt.Errorf("missing required parameters in MAP SELECT statement - cell length: %d", len(tape.Cell))
			}
			if len(tape.Cell[2].S) != 64 {
				return fmt.Errorf("syntax error - invalid Txid in SELECT command: %d", len(tape.Cell))
			}
			m[TxID] = tape.Cell[2].S
			m[SelectCmd] = tape.Cell[3].S

			// Build new command from SELECT
			newCells := []bob.Cell{{S: Prefix}, {S: tape.Cell[3].S}}
			newCells = append(newCells, tape.Cell[4:]...)
			switch m[SelectCmd] {
			case Add:
				m.add(newCells)
			case Delete:
				m.delete(newCells)
			case Set:
				m.set(newCells)
			case Remove:
				m.remove(newCells)
			}
		}
	}
	return nil
}
