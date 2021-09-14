package magic

import (
	"testing"

	"github.com/bitcoinschema/go-bob"
)

func TestNewFromTape(t *testing.T) {
	tape := bob.Tape{
		Cell: []bob.Cell{
			{S: Prefix},
			{S: Set},
			{S: "app"},
			{S: "myapp"},
		},
	}

	tx, err := NewFromTape(&tape)
	if err != nil {
		t.Errorf("Failed to create new magic from tape")
	}

	if tx["app"] != "myapp" {
		t.Errorf("Unexpected output %+v %s", tx, tx["app"])
	}
}
