package mattercloud

import "testing"

// TestClient_Transaction tests the Transaction()
func TestClient_Transaction(t *testing.T) {

	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient(testAPIKey, NetworkMain, nil)
	if err != nil {
		t.Fatal(err)
	}

	var resp *Transaction
	tx := "96b3dc5941ce97046d4af6e7a69f4b38c48f05ef071c2a33f88807b89ab51da6"
	resp, err = client.Transaction(tx)
	if err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	if resp.TxID != tx {
		t.Fatal("tx returned does not match", resp.TxID, tx)
	}

	if resp.Hash != tx {
		t.Fatal("hash returned does not match", resp.Hash, tx)
	}

	if resp.Version != 1 {
		t.Fatal("version returned does not match", resp.Version)
	}

	if resp.BlockHash != "0000000000000000078f34d9cd3f48e4948aef4c79548ec777050e1c8953a85c" {
		t.Fatal("block hash returned does not match", resp.BlockHash)
	}

	if resp.Time != 1554007897 {
		t.Fatal("time returned does not match", resp.Time)
	}

	// todo: test all values!
}

// TestClient_TransactionBatch tests the TransactionBatch()
func TestClient_TransactionBatch(t *testing.T) {

	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient(testAPIKey, NetworkMain, nil)
	if err != nil {
		t.Fatal(err)
	}

	var resp []*Transaction
	resp, err = client.TransactionBatch([]string{"96b3dc5941ce97046d4af6e7a69f4b38c48f05ef071c2a33f88807b89ab51da6", "bdf6f49776faaa4790af3e41b8b474a7d0d47df540f8d71c3579dc0addd64c45"})
	if err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	if len(resp) == 0 {
		t.Fatal("missing balance results")
	} else if len(resp) != 2 {
		t.Fatal("should be two results")
	}

	// todo: test all the values in the result
}

// TestClient_Broadcast tests the Broadcast()
func TestClient_Broadcast(t *testing.T) {

	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient(testAPIKey, NetworkMain, nil)
	if err != nil {
		t.Fatal(err)
	}

	var resp *BroadcastResponse
	rawTx := "0100000001270e55963a167a2fae66307efa3565032402c1387d62e5276464295d2a6834d8010000008a4730440220132f6d484de9d34d314aec945865af5da95f35cf4c7cc271d40bc99f8d7f12e3022051fcb2ce4461d1c6e8a778f5e4dcb27c8461d18e0652f68a7a09a98e95df5cb74141044e2c1e2c055e7aefc291679882382c35894a6aa6dd95644f598e506c239f9d83b1d9671c1d9673e3c2b74f07e8032343f3adc21367bd4cffae92fe31efcd598affffffff020000000000000000456a2231394878696756345179427633744870515663554551797131707a5a56646f41757404617364660d746578742f6d61726b646f776e055554462d3807616e6f7468657240390000000000001976a91410bdcba3041b5e5517a58f2e405293c14a7c70c188ac00000000"
	resp, err = client.Broadcast(rawTx)
	if err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	if resp.TxID != "96b3dc5941ce97046d4af6e7a69f4b38c48f05ef071c2a33f88807b89ab51da6" {
		t.Fatal("expected tx id", resp.TxID)
	}
}
