// Package magic is a library for working with Magic Attribute Protocol and used in conjunction with a Tape
// from BOB transaction
//
// Protocol: https://github.com/rohenaz/MAP
// BOB: https://bob.planaria.network/
//
// If you have any suggestions or comments, please feel free to open an issue on
// this GitHub repository!
//
// By BitcoinSchema Organization (https://bitcoinschema.org)
package magic

import (
	"strings"

	"github.com/bitcoinschema/go-bob"
)

// MAP is Magic Attribute Protocol
type MAP map[string]interface{} // `json:"MAP,omitempty" bson:"MAP, omitempty"`

// Prefix is the Bitcom prefix for Magic Attribute Protocol
const Prefix = "1PuQa7K62MiKCtssSLKy1kh56WWU7MtUR5"

// MAP Commands
const (
	Cmd       = "CMD"
	Set       = "SET"
	Add       = "ADD"
	Delete    = "DELETE"
	Remove    = "REMOVE"
	Select    = "SELECT"
	TxID      = "TXID"
	SelectCmd = "SELECT_CMD"
)

// MAP Keys
const (
	MapKeyKey   = "key"
	MapValueKey = "value"
	MapAppKey   = "app"
	MapTypeKey  = "type"
)

// set is: MAP SET
func (m MAP) set(cells []bob.Cell) {
	for idx, cell := range cells {
		// Skip prefix (0) and command (1)
		if idx < 2 {
			continue
		}

		if idx%2 == 1 {
			m[cells[idx-1].S] = cell.S
		}
	}
}

// getValues will return all values in a slice of strings
func (m MAP) getValues(key string) (values []string) {
	values = m[key].([]string)
	return
}

// getValue will return all values in one concatenated string
func (m MAP) getValue(key string) (value string) {
	var data = m[key].([]string)
	value = strings.Join(data, "")
	return
}

// set is: MAP SET
func (m MAP) add(cells []bob.Cell) {
	keyValues := make([]string, 0)
	keyName := cells[2].S
	for idx, cell := range cells {
		// Skip prefix (0), command (1) and keyName (2)
		if idx < 3 {
			continue
		}
		keyValues = append(keyValues, cell.S)
	}
	m[keyName] = keyValues
}

// remove is: MAP REMOVE
func (m MAP) remove(cells []bob.Cell) {
	// Skip prefix (0) and command (1)
	m[MapKeyKey] = cells[2].S
}

// delete is: MAP DELETE
func (m MAP) delete(cells []bob.Cell) {
	// Skip prefix (0) and command (1)
	m[MapKeyKey] = cells[2].S
	m[MapValueKey] = cells[3].S
}
