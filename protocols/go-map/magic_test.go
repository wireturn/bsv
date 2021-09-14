package magic

import (
	"testing"

	"github.com/bitcoinschema/go-bob"
)

const mapKey = "key"
const mapValue = "value"
const mapTestKey = "keyName1"
const mapTestValue = "something"

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func TestSelectDelete(t *testing.T) {
	tape := &bob.Tape{
		Cell: []bob.Cell{
			{S: Prefix},
			{S: Select},
			{S: "a9a4387d2baa2edcc53ec040b3affbc38778e9dd876f9a47e5c767c785aacf76"},
			{S: Delete},
			{S: mapTestKey},
			{S: mapTestValue},
		},
	}

	m, err := NewFromTape(tape)
	if err != nil {
		t.Fatalf("Failed to create magicTx from tape %s", err)
	}

	if m[Cmd] != Select || m[mapKey] != mapTestKey || m[mapValue] != mapTestValue {
		t.Fatalf("SELECT + DELETE Failed. command: %s, key: %s, value: %s", m[Cmd], m[mapKey], m[mapValue])
	}
}

func TestAdd(t *testing.T) {
	tape := bob.Tape{
		Cell: []bob.Cell{
			{S: Prefix},
			{S: Add},
			{S: "keyName"},
			{S: mapTestValue},
			{S: "something else"},
		},
	}
	m, err := NewFromTape(&tape)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	switch m["keyName"].(type) {
	case []string:
		if !contains(m["keyName"].([]string), mapTestValue) ||
			!contains(m["keyName"].([]string), "something else") {
			t.Fatalf("ADD Failed %s", m["keyName1"])
		}
	default:
		t.Fatalf("ADD Failed %s", m[mapTestKey])
	}
}

func TestGetValue(t *testing.T) {
	tape := bob.Tape{
		Cell: []bob.Cell{
			{S: Prefix},
			{S: Add},
			{S: "keyName"},
			{S: mapTestValue},
		},
	}
	m, err := NewFromTape(&tape)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	if val := m.getValue("keyName"); val != "something" {
		t.Fatalf("expected: [%v] but got: [%v]", "something", val)
	}
}

func TestGetValues(t *testing.T) {
	tape := bob.Tape{
		Cell: []bob.Cell{
			{S: Prefix},
			{S: Add},
			{S: "keyName"},
			{S: mapTestValue},
			{S: "another value"},
			{S: "third value"},
		},
	}
	m, err := NewFromTape(&tape)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	if val := m.getValues("keyName"); val[0] != "something" {
		t.Fatalf("expected: [%v] but got: [%v]", "something", val)
	} else if val[1] != "another value" {
		t.Fatalf("expected: [%v] but got: [%v]", "another value", val)
	}
}

func TestDelete(t *testing.T) {
	tape := bob.Tape{
		Cell: []bob.Cell{
			{S: Prefix},
			{S: Delete},
			{S: "keyName"},
			{S: mapTestValue},
		},
	}
	m, err := NewFromTape(&tape)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	if m[mapKey] != "keyName" || m[mapValue] != mapTestValue {
		t.Errorf("DELETE Failed %s %s", m[mapKey], m[mapValue])
	}

}

func TestSet(t *testing.T) {
	tape := bob.Tape{
		Cell: []bob.Cell{
			{S: Prefix},
			{S: Set},
			{S: mapTestKey},
			{S: mapTestValue},
			{S: "keyName2"},
			{S: "something else"},
		},
	}
	m, err := NewFromTape(&tape)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}
	if m[mapTestKey] != mapTestValue {
		t.Errorf("SET Failed %s", m[mapTestKey])
	}
}

func TestRemove(t *testing.T) {
	tape := bob.Tape{
		Cell: []bob.Cell{
			{S: Prefix},
			{S: Remove},
			{S: "keyName1"},
		},
	}
	m, err := NewFromTape(&tape)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}
	if m[mapKey] != "keyName1" {
		t.Errorf("REMOVE Failed %s", m[mapKey])
	}
}
