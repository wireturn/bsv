package mattercloud

import "testing"

const (
	testAddress = `12XXBHkRNrBEb7GCvAP4G8oUs5SoDREkVX`
)

// TestClient_AddressBalance tests the AddressBalance()
func TestClient_AddressBalance(t *testing.T) {

	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient(testAPIKey, NetworkMain, nil)
	if err != nil {
		t.Fatal(err)
	}

	var resp *Balance
	resp, err = client.AddressBalance(testAddress)
	if err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	if resp.Address != testAddress {
		t.Fatal("failed to get the address", testAddress, resp.Address)
	}

	// Can't test for confirmed or unconfirmed, might change!
}

// TestClient_AddressBalanceBatch tests the AddressBalanceBatch()
func TestClient_AddressBalanceBatch(t *testing.T) {

	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient(testAPIKey, NetworkMain, nil)
	if err != nil {
		t.Fatal(err)
	}

	var resp []*Balance
	resp, err = client.AddressBalanceBatch([]string{"1GJ3x5bcEnKMnzNFPPELDfXUCwKEaLHM5H", testAddress})
	if err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	if len(resp) == 0 {
		t.Fatal("missing balance results")
	} else if len(resp) != 2 {
		t.Fatal("should be two results")
	}

	// Can't test for confirmed or unconfirmed, might change!
}

// TestClient_AddressUtxos tests the AddressUtxos()
func TestClient_AddressUtxos(t *testing.T) {

	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient(testAPIKey, NetworkMain, nil)
	if err != nil {
		t.Fatal(err)
	}

	var resp []*UnspentTransaction
	resp, err = client.AddressUtxos(testAddress)
	if err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	if len(resp) == 0 {
		t.Fatal("there are no utxos")
	}

	// Can't test for the UTXOs as they might change
}

// TestClient_AddressUtxosBatch tests the AddressUtxosBatch()
func TestClient_AddressUtxosBatch(t *testing.T) {

	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient(testAPIKey, NetworkMain, nil)
	if err != nil {
		t.Fatal(err)
	}

	var resp []*UnspentTransaction
	resp, err = client.AddressUtxosBatch([]string{"1GJ3x5bcEnKMnzNFPPELDfXUCwKEaLHM5H", testAddress})
	if err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	if len(resp) == 0 {
		t.Fatal("no utxos found")
	}

	// Can't test for UTXOs since they may change
}

// TestClient_AddressHistory tests the AddressHistory()
func TestClient_AddressHistory(t *testing.T) {

	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient(testAPIKey, NetworkMain, nil)
	if err != nil {
		t.Fatal(err)
	}

	var resp *History
	resp, err = client.AddressHistory(testAddress)
	if err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	if resp.From > 0 {
		t.Fatal("from should be zero", resp.From)
	}

	if resp.To != 20 {
		t.Fatal("to should be 20 by default", resp.To)
	}

	if len(resp.Results) == 0 {
		t.Fatal("this address has history but its missing", resp.Results)
	}

	// Can't test History as that may change or grow
}

// TestClient_AddressHistoryBatch tests the AddressHistoryBatch()
func TestClient_AddressHistoryBatch(t *testing.T) {

	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient(testAPIKey, NetworkMain, nil)
	if err != nil {
		t.Fatal(err)
	}

	var resp *History
	resp, err = client.AddressHistoryBatch([]string{"1GJ3x5bcEnKMnzNFPPELDfXUCwKEaLHM5H", testAddress})
	if err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	if resp.From > 0 {
		t.Fatal("from should be zero", resp.From)
	}

	if resp.To != 20 {
		t.Fatal("to should be 20 by default", resp.To)
	}

	if len(resp.Results) == 0 {
		t.Fatal("this address has history but its missing", resp.Results)
	}

	// Can't test History as that may change or grow
}
