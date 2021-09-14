package client

import "testing"

func TestInterface(t *testing.T) {
	rc, err := NewRemoteClient(&Config{})
	if err != nil {
		t.Fatalf("Failed to create remote client : %s", err)
	}

	testInterface(rc)
}

// Used to ensure remote client fulfills interface.
func testInterface(c Client) {}
