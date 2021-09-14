package whatsonchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

// mockHTTPTransactions for mocking requests
type mockHTTPTransactions struct{}

// Do is a mock http request
func (m *mockHTTPTransactions) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	//
	// Get Tx by Hash
	//

	// Valid
	if strings.Contains(req.URL.String(), "/tx/hash/c1d32f28baa27a376ba977f6a8de6ce0a87041157cef0274b20bfda2b0d8df96") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"hex":"","txid":"c1d32f28baa27a376ba977f6a8de6ce0a87041157cef0274b20bfda2b0d8df96","hash":"c1d32f28baa27a376ba977f6a8de6ce0a87041157cef0274b20bfda2b0d8df96","version":1,"size":113,"locktime":0,"vin":[{"coinbase":"03d7c6082f7376706f6f6c2e636f6d2f3edff034600055b8467f0040","txid":"","vout":0,"scriptSig":{"asm":"","hex":""},"sequence":4294967295}],"vout":[{"value":12.5000042,"n":0,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 492558fb8ca71a3591316d095afc0f20ef7d42f7 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a914492558fb8ca71a3591316d095afc0f20ef7d42f788ac","reqSigs":1,"type":"pubkeyhash","addresses":["17fm4xevwDh3XRHv9UoqYrVgPMbwcGHsUs"],"opReturn":null,"isTruncated":false}}],"blockhash":"0000000000000000091216c46973d82db057a6f9911352892b7769ed517681c3","confirmations":65369,"time":1553501874,"blocktime":1553501874}`)))
	}

	// Invalid - return an error
	if strings.Contains(req.URL.String(), "/tx/hash/error") {
		resp.StatusCode = http.StatusInternalServerError
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
		return resp, fmt.Errorf("missing request")
	}

	//
	// Get Merkle Proof
	//

	// Valid
	if strings.Contains(req.URL.String(), "/tx/c1d32f28baa27a376ba977f6a8de6ce0a87041157cef0274b20bfda2b0d8df96/proof") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`[{"blockHash":"0000000000000000091216c46973d82db057a6f9911352892b7769ed517681c3","branches":[{"hash":"7e0ba1980522125f1f40d19a249ab3ae036001b991776813d25aebe08e8b8a50","pos":"R"},{"hash":"1e3a5a8946e0caf07006f6c4f76773d7e474d4f240a276844f866bd09820adb3","pos":"R"}],"hash":"c1d32f28baa27a376ba977f6a8de6ce0a87041157cef0274b20bfda2b0d8df96","merkleRoot":"95a920b1002bed05379a0d2650bb13eb216138f28ee80172f4cf21048528dc60"}]`)))
	}

	// Invalid - invalid length
	if strings.Contains(req.URL.String(), "/tx/error/proof") {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`txid must be 64 hex characters in length`)))
		return resp, fmt.Errorf("txid must be 64 hex characters in length")
	}

	// Invalid - tx is not valid
	if strings.Contains(req.URL.String(), "/tx/c1d32f28baa27a376ba977f6a8de6ce0a87041157cef0274b20bfda2b0d8dfzz/proof") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`null`)))
	}

	//
	// Get Raw Tx Data
	//

	// Valid
	if strings.Contains(req.URL.String(), "/tx/c1d32f28baa27a376ba977f6a8de6ce0a87041157cef0274b20bfda2b0d8df96/hex") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff1c03d7c6082f7376706f6f6c2e636f6d2f3edff034600055b8467f0040ffffffff01247e814a000000001976a914492558fb8ca71a3591316d095afc0f20ef7d42f788ac00000000`)))
	}

	// Invalid - invalid length
	if strings.Contains(req.URL.String(), "/tx/error/hex") {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`txid must be 64 hex characters in length`)))
		return resp, fmt.Errorf("txid must be 64 hex characters in length")
	}

	// Invalid - tx is not valid
	if strings.Contains(req.URL.String(), "/tx/c1d32f28baa27a376ba977f6a8de6ce0a87041157cef0274b20bfda2b0d8dfzz/hex") {
		resp.StatusCode = http.StatusNotFound
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
	}

	//
	// Get Raw Tx Output Data
	//

	// Valid
	if strings.Contains(req.URL.String(), "/tx/c1d32f28baa27a376ba977f6a8de6ce0a87041157cef0274b20bfda2b0d8df96/out/0/hex") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`76a914492558fb8ca71a3591316d095afc0f20ef7d42f788ac`)))
	}

	// Invalid - invalid length
	if strings.Contains(req.URL.String(), "/tx/error/out/0/hex") {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`txid must be 64 hex characters in length`)))
		return resp, fmt.Errorf("txid must be 64 hex characters in length")
	}

	// Invalid - tx is not valid
	if strings.Contains(req.URL.String(), "/tx/c1d32f28baa27a376ba977f6a8de6ce0a87041157cef0274b20bfda2b0d8dfzz/out/0/hex") {
		resp.StatusCode = http.StatusBadGateway
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
	}

	//
	// Bulk Tx Data
	//

	var data TxHashes
	if strings.Contains(req.URL.String(), "/test/txs") {

		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&data)
		if err != nil {
			return resp, err
		}

		// Valid (two txs)
		if strings.Contains(data.TxIDs[0], "294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1aa") {
			resp.StatusCode = http.StatusOK
			resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`[{"hex":"","txid":"294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1aa","hash":"294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1aa","version":1,"size":1044,"locktime":0,"vin":[{"coinbase":"","txid":"57a79dbc0be8225f4d9dc705a91096e60d3da8f3a7fdf009715dc28975a8dbb8","vout":0,"scriptSig":{"asm":"304402206b9956a6dd39d7f081f6a2d3731b2b7c875ec563743d7620ea55b37f3c2d1ffd02202594d52b6df818d100a3e721b0d245625cbe200e51a7a846fd37436607e08a0a[ALL|FORKID] 03998489e31affb06deaa6890deea42447ac9b4b31c5fc93d023525de9f850b281","hex":"47304402206b9956a6dd39d7f081f6a2d3731b2b7c875ec563743d7620ea55b37f3c2d1ffd02202594d52b6df818d100a3e721b0d245625cbe200e51a7a846fd37436607e08a0a412103998489e31affb06deaa6890deea42447ac9b4b31c5fc93d023525de9f850b281"},"sequence":4294967295}],"vout":[{"value":0,"n":0,"scriptPubKey":{"asm":"OP_RETURN 6761746c696e6720626974636f696e626c6f636b732e6c697665","hex":"6a1a6761746c696e6720626974636f696e626c6f636b732e6c697665","type":"nulldata","opReturn":{"type":"OP_RETURN","action":"","text":"","parts":null},"isTruncated":false}},{"value":0.000013,"n":1,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 8432682cb8f5cbbb36913266fbe84176a999ac99 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a9148432682cb8f5cbbb36913266fbe84176a999ac9988ac","reqSigs":1,"type":"pubkeyhash","addresses":["1D3zaf652ajAPZVC1CrKgMQjesaPuLwcSW"],"opReturn":null,"isTruncated":false}},{"value":0.000013,"n":2,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 f6b68ada06df7323aca315ec54aa251e806c5168 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a914f6b68ada06df7323aca315ec54aa251e806c516888ac","reqSigs":1,"type":"pubkeyhash","addresses":["1PVVuTZCHdW83ucCviYFWkjqtDvjismrSV"],"opReturn":null,"isTruncated":false}},{"value":0.000013,"n":3,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 d47ab37ea2bcffe4a7bf98086a65d489acdf2fd6 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a914d47ab37ea2bcffe4a7bf98086a65d489acdf2fd688ac","reqSigs":1,"type":"pubkeyhash","addresses":["1LNVF6RV5dnSAzStK19ZM75EJMG8nfg4iv"],"opReturn":null,"isTruncated":false}},{"value":0.000013,"n":4,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 a20568923f3c65f91c0d7a3be3ca4f8b14c7d048 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a914a20568923f3c65f91c0d7a3be3ca4f8b14c7d04888ac","reqSigs":1,"type":"pubkeyhash","addresses":["1Fmgws8tbKmc6jcmqnhuShUSGFnLUfjpN3"],"opReturn":null,"isTruncated":false}},{"value":0.000013,"n":5,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 0e80337d887c653d2bd7228281aa7ecf4bf2c9dc OP_EQUALVERIFY OP_CHECKSIG","hex":"76a9140e80337d887c653d2bd7228281aa7ecf4bf2c9dc88ac","reqSigs":1,"type":"pubkeyhash","addresses":["12Kg3FtBJhZMZLMhDy6BQ12Numjzdp5Pxp"],"opReturn":null,"isTruncated":false}},{"value":0.000013,"n":6,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 7bc3cebb1dc5ee6356931add1f5f1196fce7b531 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a9147bc3cebb1dc5ee6356931add1f5f1196fce7b53188ac","reqSigs":1,"type":"pubkeyhash","addresses":["1CHQgYRLUXfh1kyXYgzVthjkD8u4L7Tppd"],"opReturn":null,"isTruncated":false}},{"value":0.000013,"n":7,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 a8854465a2bbefeb1d47c5c79b8aef737d62f06f OP_EQUALVERIFY OP_CHECKSIG","hex":"76a914a8854465a2bbefeb1d47c5c79b8aef737d62f06f88ac","reqSigs":1,"type":"pubkeyhash","addresses":["1GN4Ao1PW5A8D9Z7Mwvd1sHnBBqfVT3T9v"],"opReturn":null,"isTruncated":false}},{"value":0.000013,"n":8,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 6c376c6d05ff152e1537667cfb52e628547f620a OP_EQUALVERIFY OP_CHECKSIG","hex":"76a9146c376c6d05ff152e1537667cfb52e628547f620a88ac","reqSigs":1,"type":"pubkeyhash","addresses":["1AsCNcXyTJoYXoJ9X8TBtWe7aUuEujPHaC"],"opReturn":null,"isTruncated":false}},{"value":0.000013,"n":9,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 7c0bd6447ac246f100fec9392a396191a06fcae2 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a9147c0bd6447ac246f100fec9392a396191a06fcae288ac","reqSigs":1,"type":"pubkeyhash","addresses":["1CJtyEdJQYK5DE8skg6qZH2rfNmy2ZdG4F"],"opReturn":null,"isTruncated":false}},{"value":0.000013,"n":10,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 5221bd8f6901612cb1bf0618a398232292334749 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a9145221bd8f6901612cb1bf0618a39823229233474988ac","reqSigs":1,"type":"pubkeyhash","addresses":["18VGq4RXaKV1xXSPw9efaJkaST8BUFBx4G"],"opReturn":null,"isTruncated":false}},{"value":0.000013,"n":11,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 32c248707616cd030a3f346aec13c9d5ae8d99b0 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a91432c248707616cd030a3f346aec13c9d5ae8d99b088ac","reqSigs":1,"type":"pubkeyhash","addresses":["15dPXcJWgkskbpHppUrPkGStf3ZxCh7qaz"],"opReturn":null,"isTruncated":false}},{"value":0.000013,"n":12,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 c31eaefc851f6de4830f51d649e1731706159dc8 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a914c31eaefc851f6de4830f51d649e1731706159dc888ac","reqSigs":1,"type":"pubkeyhash","addresses":["1JnhXLXhSEB4oYJaAEHFL8iBd3isjcyw4s"],"opReturn":null,"isTruncated":false}},{"value":0.000013,"n":13,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 3d980ead104586cba895d944280f93ca8323b475 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a9143d980ead104586cba895d944280f93ca8323b47588ac","reqSigs":1,"type":"pubkeyhash","addresses":["16cgNq9FwK531hi8QPBjjNEM9V5iZXTqaM"],"opReturn":null,"isTruncated":false}},{"value":0.000013,"n":14,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 73721ce4189c5b996819f81f8cce7366a1bb7804 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a91473721ce4189c5b996819f81f8cce7366a1bb780488ac","reqSigs":1,"type":"pubkeyhash","addresses":["1BXRQm7ytETWWLrDy5M5c5qRuc1LXwbLu8"],"opReturn":null,"isTruncated":false}},{"value":0.000013,"n":15,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 1fd249a4cac1a2c95b95ad8ecbbe23129831232e OP_EQUALVERIFY OP_CHECKSIG","hex":"76a9141fd249a4cac1a2c95b95ad8ecbbe23129831232e88ac","reqSigs":1,"type":"pubkeyhash","addresses":["13uFrzSE6YfxH7nmf9JVAoTmQRY5FbeTxz"],"opReturn":null,"isTruncated":false}},{"value":0.000013,"n":16,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 ae4661a1c0f8192e79b5c304cadd39470077c8a4 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a914ae4661a1c0f8192e79b5c304cadd39470077c8a488ac","reqSigs":1,"type":"pubkeyhash","addresses":["1GtUtZby195ofFR1Wzwy1QacVYH2QCPfhj"],"opReturn":null,"isTruncated":false}},{"value":0.000013,"n":17,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 e99c1003b189bfe7ababb7cbe2e39b8eb2317ccf OP_EQUALVERIFY OP_CHECKSIG","hex":"76a914e99c1003b189bfe7ababb7cbe2e39b8eb2317ccf88ac","reqSigs":1,"type":"pubkeyhash","addresses":["1NJDQ7Z8w9CybHQwKJmCEP43rYKkoPmTCm"],"opReturn":null,"isTruncated":false}},{"value":0.000013,"n":18,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 e26673b87c7e644688b0a62f1e2cd5e7cc4db63c OP_EQUALVERIFY OP_CHECKSIG","hex":"76a914e26673b87c7e644688b0a62f1e2cd5e7cc4db63c88ac","reqSigs":1,"type":"pubkeyhash","addresses":["1Me6Sr4CRn766zaTKepxPo3Fo6MmhVoRzj"],"opReturn":null,"isTruncated":false}},{"value":0.000013,"n":19,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 5cce10f43fa0458b2f8c3982a9349e7a67e67365 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a9145cce10f43fa0458b2f8c3982a9349e7a67e6736588ac","reqSigs":1,"type":"pubkeyhash","addresses":["19Ti2P6ByGxhVR4rFdv6d82o8x8agHDwJU"],"opReturn":null,"isTruncated":false}},{"value":0.000013,"n":20,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 4f062ff6e6921a8b6f43c79fc0e89f582a222253 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a9144f062ff6e6921a8b6f43c79fc0e89f582a22225388ac","reqSigs":1,"type":"pubkeyhash","addresses":["18Cqo4UpM7bcN8Vsy8dg7RqKWyMYTcdWUX"],"opReturn":null,"isTruncated":false}},{"value":0.000013,"n":21,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 6dacfcede9e86678de202cb30d17bd1eb656cec3 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a9146dacfcede9e86678de202cb30d17bd1eb656cec388ac","reqSigs":1,"type":"pubkeyhash","addresses":["1AzutKxUnCTu9H8oETDCw7euH1oUpWnNu7"],"opReturn":null,"isTruncated":false}},{"value":0.000013,"n":22,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 fcdaf0d9fd36b2c1dbfa90e659a373290795230b OP_EQUALVERIFY OP_CHECKSIG","hex":"76a914fcdaf0d9fd36b2c1dbfa90e659a373290795230b88ac","reqSigs":1,"type":"pubkeyhash","addresses":["1Q3yZcm7nGKCabp5nWNNfsu9Fci2m1gNbP"],"opReturn":null,"isTruncated":false}},{"value":0.000013,"n":23,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 82f1a346d023fd8e5eb976bba6df7e6a20950032 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a91482f1a346d023fd8e5eb976bba6df7e6a2095003288ac","reqSigs":1,"type":"pubkeyhash","addresses":["1CwNKE9EF8eEid1Lt2NQQHtHqB9qh4i5eV"],"opReturn":null,"isTruncated":false}},{"value":0.000013,"n":24,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 2f82e5b519a351b646d23ceed7d43bd532132094 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a9142f82e5b519a351b646d23ceed7d43bd53213209488ac","reqSigs":1,"type":"pubkeyhash","addresses":["15LDZvLDJnt3hWqpDLXmsXX1YkH9HxcCbY"],"opReturn":null,"isTruncated":false}},{"value":0.000013,"n":25,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 3d0e5368bdadddca108a0fe44739919274c726c7 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a9143d0e5368bdadddca108a0fe44739919274c726c788ac","reqSigs":1,"type":"pubkeyhash","addresses":["16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA"],"opReturn":null,"isTruncated":false}}],"blockhash":"000000000000000004b5ce6670f2ff27354a1e87d0a01bf61f3307f4ccd358b5","confirmations":28395,"time":1575841517,"blocktime":1575841517},{"hex":"","txid":"91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba3258","hash":"91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba3258","version":1,"size":118,"locktime":0,"vin":[{"coinbase":"033f6b092f636f696e6765656b2e636f6d2f7759319af9d4f815e3a2fae5e60000","txid":"","vout":0,"scriptSig":{"asm":"","hex":""},"sequence":4294967295}],"vout":[{"value":12.5133703,"n":0,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 8460e9a972a8600766a1b38fac4a2cfb8692d3ad OP_EQUALVERIFY OP_CHECKSIG","hex":"76a9148460e9a972a8600766a1b38fac4a2cfb8692d3ad88ac","reqSigs":1,"type":"pubkeyhash","addresses":["1D4xHwLxA8E9vU87N1ELHtPEZdKeLhywY1"],"opReturn":null,"isTruncated":false}}],"blockhash":"000000000000000002e8d4b4c0385abd195709c82f16d9917f081b70000e8804","confirmations":23367,"time":1578837295,"blocktime":1578837295}]`)))
		}

		// Valid (1 bad tx, 1 good)
		if strings.Contains(data.TxIDs[0], "294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1ZZ") {
			resp.StatusCode = http.StatusOK
			resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`[{"hex":"","txid":"91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba3258","hash":"91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba3258","version":1,"size":118,"locktime":0,"vin":[{"coinbase":"033f6b092f636f696e6765656b2e636f6d2f7759319af9d4f815e3a2fae5e60000","txid":"","vout":0,"scriptSig":{"asm":"","hex":""},"sequence":4294967295}],"vout":[{"value":12.5133703,"n":0,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 8460e9a972a8600766a1b38fac4a2cfb8692d3ad OP_EQUALVERIFY OP_CHECKSIG","hex":"76a9148460e9a972a8600766a1b38fac4a2cfb8692d3ad88ac","reqSigs":1,"type":"pubkeyhash","addresses":["1D4xHwLxA8E9vU87N1ELHtPEZdKeLhywY1"],"opReturn":null,"isTruncated":false}}],"blockhash":"000000000000000002e8d4b4c0385abd195709c82f16d9917f081b70000e8804","confirmations":23367,"time":1578837295,"blocktime":1578837295}]`)))
		}

		// Invalid (two bad txs)
		if strings.Contains(data.TxIDs[0], "294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1VV") {
			resp.StatusCode = http.StatusOK
			resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`[]`)))
		}

		// Valid - for AddressDetails
		if strings.Contains(data.TxIDs[0], "33b9432a0ea203bbb6ec00592622cf6e90223849e4c9a76447a19a3ed43907d3") {
			resp.StatusCode = http.StatusOK
			resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`[{"hex":"","txid":"33b9432a0ea203bbb6ec00592622cf6e90223849e4c9a76447a19a3ed43907d3","hash":"33b9432a0ea203bbb6ec00592622cf6e90223849e4c9a76447a19a3ed43907d3","version":1,"size":440,"locktime":0,"vin":[{"coinbase":"","txid":"fabe0b5d0979e068dce986692d1c5620f37383657a2fe7969f1cfe4a81b7f517","vout":3,"scriptSig":{"asm":"30450221008f74bb75c331cb7902a4e7539ee60fafe2c9a73d325aba6fc3ff9105ed91e219022064e65a5662c0593086ab05a0131e5abac5ef249f5f33c74351c2bed653da269f[ALL|FORKID] 026d6fc8f05b630e637507084b1678ec753c75b9e050312919e1d973224c5c3103","hex":"4830450221008f74bb75c331cb7902a4e7539ee60fafe2c9a73d325aba6fc3ff9105ed91e219022064e65a5662c0593086ab05a0131e5abac5ef249f5f33c74351c2bed653da269f4121026d6fc8f05b630e637507084b1678ec753c75b9e050312919e1d973224c5c3103"},"sequence":4294967295}],"vout":[{"value":0,"n":0,"scriptPubKey":{"asm":"0 OP_RETURN 3150755161374b36324d694b43747373534c4b79316b683536575755374d74555235 5522771 7368801 746f6e6963706f77 1701869940 6f666665725f636c69636b 6f666665725f636f6e6669675f6964 56 6f666665725f73657373696f6e5f6964 33316566313830633732363465303032373836333261306131613830313835313336363236336537306361383233353138373664636436386563666163623365","hex":"006a223150755161374b36324d694b43747373534c4b79316b683536575755374d74555235035345540361707008746f6e6963706f7704747970650b6f666665725f636c69636b0f6f666665725f636f6e6669675f69640138106f666665725f73657373696f6e5f69644033316566313830633732363465303032373836333261306131613830313835313336363236336537306361383233353138373664636436386563666163623365","type":"nulldata","opReturn":{"type":"bitcom","action":"","text":"","parts":null},"isTruncated":false}},{"value":0.00000549,"n":1,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 09cc4559bdcb84cb35c107743f0dbb10d66679cc OP_EQUALVERIFY OP_CHECKSIG","hex":"76a91409cc4559bdcb84cb35c107743f0dbb10d66679cc88ac","reqSigs":1,"type":"pubkeyhash","addresses":["1tonicZQwN2BNKhVwPXqh8ez3q56y1EYw"],"opReturn":null,"isTruncated":false}},{"value":0.00005489,"n":2,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 b5195bf7db0652f536a7dddbe36a99a091125468 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a914b5195bf7db0652f536a7dddbe36a99a09112546888ac","reqSigs":1,"type":"pubkeyhash","addresses":["1HWZgiMKQKPSkLzT7hipS22AvkQZJsyxmT"],"opReturn":null,"isTruncated":false}},{"value":0.0245168,"n":3,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 bf49c6a5406675e174f4f6a83b3d94dd9d845398 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a914bf49c6a5406675e174f4f6a83b3d94dd9d84539888ac","reqSigs":1,"type":"pubkeyhash","addresses":["1JSSSgcyufLgbXFw6WAXyXgBrmgFpnqXWh"],"opReturn":null,"isTruncated":false}}],"blockhash":"0000000000000000026b9da3860e4c8ee351a7af46da6042eaa5d110113b9fad","confirmations":1348,"time":1592122768,"blocktime":1592122768},{"hex":"","txid":"4805041897a2ae59ffca85f0deb46e89d73d1ba4478bbd9c0fcd76ba0985ded2","hash":"4805041897a2ae59ffca85f0deb46e89d73d1ba4478bbd9c0fcd76ba0985ded2","version":1,"size":439,"locktime":0,"vin":[{"coinbase":"","txid":"5a45b8415e5c1740353cfb011d29e04ec104865be6560dff5bd6cb31db75d559","vout":3,"scriptSig":{"asm":"3044022008e2417d072cfbb95d4e04c7e6e6ab70e415a379fb912cb2e0503e3df0ae0d2002201f9fcbf6c65ba6624fe0669d08155ed7c0d19c28be72daf3e00de2613656f955[ALL|FORKID] 026d6fc8f05b630e637507084b1678ec753c75b9e050312919e1d973224c5c3103","hex":"473044022008e2417d072cfbb95d4e04c7e6e6ab70e415a379fb912cb2e0503e3df0ae0d2002201f9fcbf6c65ba6624fe0669d08155ed7c0d19c28be72daf3e00de2613656f9554121026d6fc8f05b630e637507084b1678ec753c75b9e050312919e1d973224c5c3103"},"sequence":4294967295}],"vout":[{"value":0,"n":0,"scriptPubKey":{"asm":"0 OP_RETURN 3150755161374b36324d694b43747373534c4b79316b683536575755374d74555235 5522771 7368801 746f6e6963706f77 1701869940 6f666665725f636c69636b 6f666665725f636f6e6669675f6964 56 6f666665725f73657373696f6e5f6964 38313865386165656339353733646431333439373334366135363464633461623035353062333039383830373563393733316631643063653731336536353335","hex":"006a223150755161374b36324d694b43747373534c4b79316b683536575755374d74555235035345540361707008746f6e6963706f7704747970650b6f666665725f636c69636b0f6f666665725f636f6e6669675f69640138106f666665725f73657373696f6e5f69644038313865386165656339353733646431333439373334366135363464633461623035353062333039383830373563393733316631643063653731336536353335","type":"nulldata","opReturn":{"type":"bitcom","action":"","text":"","parts":null},"isTruncated":false}},{"value":0.00000573,"n":1,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 09cc4559bdcb84cb35c107743f0dbb10d66679cc OP_EQUALVERIFY OP_CHECKSIG","hex":"76a91409cc4559bdcb84cb35c107743f0dbb10d66679cc88ac","reqSigs":1,"type":"pubkeyhash","addresses":["1tonicZQwN2BNKhVwPXqh8ez3q56y1EYw"],"opReturn":null,"isTruncated":false}},{"value":0.00005726,"n":2,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 819ae5a5cbb078e96379b8eb25c29d6f7b28c412 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a914819ae5a5cbb078e96379b8eb25c29d6f7b28c41288ac","reqSigs":1,"type":"pubkeyhash","addresses":["1CpHjBbHoWzbrqQsPeZ39GLUXejZce9mBs"],"opReturn":null,"isTruncated":false}},{"value":0.02744764,"n":3,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 bf49c6a5406675e174f4f6a83b3d94dd9d845398 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a914bf49c6a5406675e174f4f6a83b3d94dd9d84539888ac","reqSigs":1,"type":"pubkeyhash","addresses":["1JSSSgcyufLgbXFw6WAXyXgBrmgFpnqXWh"],"opReturn":null,"isTruncated":false}}],"blockhash":"000000000000000003d684082ab45014f89a7f8e5e35ec94fcb4aa8b5f00c01e","confirmations":1049,"time":1592307236,"blocktime":1592307236},{"hex":"","txid":"2493ff4cbca16b892ac641b7f2cb6d4388e75cb3f8963c291183f2bf0b27f415","hash":"2493ff4cbca16b892ac641b7f2cb6d4388e75cb3f8963c291183f2bf0b27f415","version":1,"size":439,"locktime":0,"vin":[{"coinbase":"","txid":"2ebc8f094fdc012f7d9a0ed39999dcf0318665830f7d5f113af0d1c79fba2f8e","vout":3,"scriptSig":{"asm":"30440220010a62c1d79afcc274b8db821cba1f093c316d67d505a3900c231ae6dfb2dd51022031fe80787c531e1c890754d2cafdc624f3446e4d1bdca18ade83cabd3a2317ac[ALL|FORKID] 026d6fc8f05b630e637507084b1678ec753c75b9e050312919e1d973224c5c3103","hex":"4730440220010a62c1d79afcc274b8db821cba1f093c316d67d505a3900c231ae6dfb2dd51022031fe80787c531e1c890754d2cafdc624f3446e4d1bdca18ade83cabd3a2317ac4121026d6fc8f05b630e637507084b1678ec753c75b9e050312919e1d973224c5c3103"},"sequence":4294967295}],"vout":[{"value":0,"n":0,"scriptPubKey":{"asm":"0 OP_RETURN 3150755161374b36324d694b43747373534c4b79316b683536575755374d74555235 5522771 7368801 746f6e6963706f77 1701869940 6f666665725f636c69636b 6f666665725f636f6e6669675f6964 56 6f666665725f73657373696f6e5f6964 35656237306231653930306535616437626335663961333663653861643435623664336435636337666466393437343762623364326461663732636631356533","hex":"006a223150755161374b36324d694b43747373534c4b79316b683536575755374d74555235035345540361707008746f6e6963706f7704747970650b6f666665725f636c69636b0f6f666665725f636f6e6669675f69640138106f666665725f73657373696f6e5f69644035656237306231653930306535616437626335663961333663653861643435623664336435636337666466393437343762623364326461663732636631356533","type":"nulldata","opReturn":{"type":"bitcom","action":"","text":"","parts":null},"isTruncated":false}},{"value":0.00000572,"n":1,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 09cc4559bdcb84cb35c107743f0dbb10d66679cc OP_EQUALVERIFY OP_CHECKSIG","hex":"76a91409cc4559bdcb84cb35c107743f0dbb10d66679cc88ac","reqSigs":1,"type":"pubkeyhash","addresses":["1tonicZQwN2BNKhVwPXqh8ez3q56y1EYw"],"opReturn":null,"isTruncated":false}},{"value":0.00005716,"n":2,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 0405a52b27214920873fa222071a8ec9610317a4 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a9140405a52b27214920873fa222071a8ec9610317a488ac","reqSigs":1,"type":"pubkeyhash","addresses":["1NGU17f9HTyv3LffW4zxukSEwsxwf4d53"],"opReturn":null,"isTruncated":false}},{"value":0.02568774,"n":3,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 bf49c6a5406675e174f4f6a83b3d94dd9d845398 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a914bf49c6a5406675e174f4f6a83b3d94dd9d84539888ac","reqSigs":1,"type":"pubkeyhash","addresses":["1JSSSgcyufLgbXFw6WAXyXgBrmgFpnqXWh"],"opReturn":null,"isTruncated":false}}],"blockhash":"00000000000000000087222006199927280a010d0db21c6d253409f3e960c7bf","confirmations":374,"time":1592698834,"blocktime":1592698834}]`)))
		}

		// Invalid - force an error
		if strings.Contains(data.TxIDs[0], "error") {
			resp.StatusCode = http.StatusBadRequest
			resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
			return resp, fmt.Errorf("unknown error")
		}
	}

	// Valid (download receipt)
	if strings.Contains(req.URL.String(), "/receipt/c1d32f28baa27a376ba977f6a8de6ce0a87041157cef0274b20bfda2b0d8df96") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`%PDF-1.4
%Óëéá
1 0 obj
<</Creator (Chromium)
/Producer (Skia/PDF m73)
/CreationDate (D:20200622155222+00'00')
/ModDate (D:20200622155222+00'00')>>
endobj
3 0 obj
<</ca 1
/BM /Normal>>
endobj
5 0 obj`)))
	}

	// Invalid (download receipt) (invalid address)
	if strings.Contains(req.URL.String(), "/receipt/invalid") {
		resp.StatusCode = http.StatusGatewayTimeout
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
		return resp, fmt.Errorf("gateway timeout")
	}

	// Default is valid
	return resp, nil
}

// mockHTTPBroadcast for mocking requests
type mockHTTPBroadcast struct{}

// txBroadcast is the struct for broadcasting a tx
type txBroadcast struct {
	TxHex string `json:"txhex"`
}

// Do is a mock http request
func (m *mockHTTPBroadcast) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	decoder := json.NewDecoder(req.Body)
	var data txBroadcast
	err := decoder.Decode(&data)
	if err != nil {
		return resp, err
	}

	//
	// Broadcast
	//

	// Valid
	if strings.Contains(data.TxHex, "0100000001ace8d1d1d885039e2d58074188e1355f03ea782b382a2890ee4ed2fc76e51daf000000006a47304402206c408d858d8678666f112d07c6a6152121e9672ee6af45ea386592166efd0852022012a94bdba5d3efc66c695dad0fb5a3267a599f5ef4f144d10663aba34ac3f8234121036e4a5eff4cb231e0a3d34932d950a2a9f6476a01b41f38e8ddf036eae2cc9fcfffffffff0280d03b00000000001976a9146680fd90c9d68cb9cf5314c4e30fb5a3879440c988acbb360100000000001976a914847e32a2b2e0289912c3d2d0eed20feb0a8bf4c688ac00000000") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`45ed612b6a2e819af164eb55e8273b61f3bfbbf08b991db1bd1d9b48565a3297`)))
	}

	// Invalid - missing inputs
	if strings.Contains(data.TxHex, "0100000001d1bda0bde67183817b21af863adaa31fda8cafcf2083ca1eaba3054496cbde10010000006a47304402205fddd6abab6b8e94f36bfec51ba2e1f3a91b5327efa88264b5530d0c86538723022010e51693e3d52347d4d2ff142b85b460d3953e625d1e062a5fa2569623fb0ea94121029df3723daceb1fef64fa0558371bc48cc3a7a8e35d8e05b87137dc129a9d4598ffffffff0115d40000000000001976a91459cc95a8cde59ceda718dbf70e612dba4034552688ac00000000") {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`"Missing inputs"`)))
		return resp, fmt.Errorf("missing inputs")
	}

	// Invalid - already in mempool
	if strings.Contains(data.TxHex, "010000000220f177bc9f9d6311bf5a3ab3f95f1fb73904e91faec69e11f1a5780d717fb9da0e0400006b483045022100944ff1436e24ed37a0d644d6637ef36e1f79f217aefef2ab5c1fccbc474ccd2002200b0fb32081b31a27359f7e10760fe0c303a5877f49b323bd81629555f62b3987412103390e121888d94987dd77d74a40e643eaa24ae61410533046b363d602096312ceffffffff938d7f6d2600707128aedfe761b5895fa87aa72838f50efe9d8771ace92fbc6b020000006b483045022100ddecfb3b1caf072c2ed60cb3d032266f39e8858b03b3e3cd4b59ea881b20736b022052410db489955bd829878eae7d71489b30336acc3b716029046c1659615f8a76412103a38ede6ca418dce302f5241aa8a0b274046807f5527d2b4bf67d93cbf0c7fabaffffffff02d0b5b384000000001976a914501a793f77e008bdff95797738bed6c848e2bc3588ac68f032550d0000001976a91423cb7559a575b8f6d049372cba0f7cb4cfcdac3988ac00000000") {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`"257: txn-already-known"`)))
		return resp, fmt.Errorf("257: txn-already-known")
	}

	// Invalid - bad status code, no error
	if strings.Contains(data.TxHex, "bad-status-code") {
		resp.StatusCode = http.StatusExpectationFailed
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
		return resp, nil
	}

	//
	// Decode Tx
	//

	// Valid
	if strings.Contains(data.TxHex, "010000000110784fd521b55a303da0f8b4ea113a2a3b5fa71565bab86c4257cff83ab4a1b9010000006a4730440220113a56d87122f28d6b60931498951f9709527a5c095ae852dc7b17d3d7915ef802206fd0b026d2e8dd30a39daaef17f503d985e2306c8244bb0e09a957bf1c9b530441210269a7785783c12405a1eaecb3088a3d830ed7e2de6ac527f42374a55a8cc5aeadffffffff0266ba0200000000001976a914021e4ac858f0ee6e0dfdf4438857f602a900698988acde905606000000001976a914022a8c1a18378885db9054676f17a27f4219045e88ac00000000") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"txid": "cfcd9c342411592319442f705f9083938847e88f709096aa2cec15e6350b947e","hash": "cfcd9c342411592319442f705f9083938847e88f709096aa2cec15e6350b947e","version": 1,"size": 225,"locktime": 0,"vin": [{"txid": "b9a1b43af8cf57426cb8ba6515a75f3b2a3a11eab4f8a03d305ab521d54f7810","vout": 1,"scriptSig": {"asm": "30440220113a56d87122f28d6b60931498951f9709527a5c095ae852dc7b17d3d7915ef802206fd0b026d2e8dd30a39daaef17f503d985e2306c8244bb0e09a957bf1c9b5304[ALL|FORKID] 0269a7785783c12405a1eaecb3088a3d830ed7e2de6ac527f42374a55a8cc5aead","hex": "4730440220113a56d87122f28d6b60931498951f9709527a5c095ae852dc7b17d3d7915ef802206fd0b026d2e8dd30a39daaef17f503d985e2306c8244bb0e09a957bf1c9b530441210269a7785783c12405a1eaecb3088a3d830ed7e2de6ac527f42374a55a8cc5aead"},"sequence": 4294967295}],"vout": [{"value": 0.0017879,"n": 0,"scriptPubKey": {"asm": "OP_DUP OP_HASH160 021e4ac858f0ee6e0dfdf4438857f602a9006989 OP_EQUALVERIFY OP_CHECKSIG","hex": "76a914021e4ac858f0ee6e0dfdf4438857f602a900698988ac","reqSigs": 1,"type": "pubkeyhash","addresses": ["1CCe7ngEDcRY4drLGx3DSeTMfAhHTr4NS"]}},{"value": 1.06336478,"n": 1,"scriptPubKey": {"asm": "OP_DUP OP_HASH160 022a8c1a18378885db9054676f17a27f4219045e OP_EQUALVERIFY OP_CHECKSIG","hex": "76a914022a8c1a18378885db9054676f17a27f4219045e88ac","reqSigs": 1,"type": "pubkeyhash","addresses": ["1CTKcxmjZF9fk8mgtHA4tCGpnC7CvWyRw"]}}],"hex": "010000000110784fd521b55a303da0f8b4ea113a2a3b5fa71565bab86c4257cff83ab4a1b9010000006a4730440220113a56d87122f28d6b60931498951f9709527a5c095ae852dc7b17d3d7915ef802206fd0b026d2e8dd30a39daaef17f503d985e2306c8244bb0e09a957bf1c9b530441210269a7785783c12405a1eaecb3088a3d830ed7e2de6ac527f42374a55a8cc5aeadffffffff0266ba0200000000001976a914021e4ac858f0ee6e0dfdf4438857f602a900698988acde905606000000001976a914022a8c1a18378885db9054676f17a27f4219045e88ac00000000"}`)))
	}

	// Invalid
	if strings.Contains(data.TxHex, "zzzz0784fd521b55a303da0f8b4ea113a2a3b5fa71565bab86c4257cff83ab4a1b9010000006a4730440220113a56d87122f28d6b60931498951f9709527a5c095ae852dc7b17d3d7915ef802206fd0b026d2e8dd30a39daaef17f503d985e2306c8244bb0e09a957bf1c9b530441210269a7785783c12405a1eaecb3088a3d830ed7e2de6ac527f42374a55a8cc5aeadffffffff0266ba0200000000001976a914021e4ac858f0ee6e0dfdf4438857f602a900698988acde905606000000001976a914022a8c1a18378885db9054676f17a27f4219045e88ac00000000") {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`"TX decode failed"`)))
		return resp, fmt.Errorf("TX decode failed")
	}

	// Default is valid
	return resp, nil
}

// mockHTTPBroadcastBulk for mocking requests
type mockHTTPBroadcastBulk struct{}

// Do is a mock http request
func (m *mockHTTPBroadcastBulk) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	decoder := json.NewDecoder(req.Body)
	var data []string
	err := decoder.Decode(&data)
	if err != nil {
		return resp, err
	}

	//
	// Broadcast Bulk
	//

	// Valid (no feedback)
	if strings.Contains(data[0], "020000000232900e9b0e359cb95ad3853e1450591fdf01e3efaa1b3b2b5ab5c5ef784946b8010000006a473044022055ec4f9b9cbdd97cf5f4893f921da75287483edb4ebba4f5b23231577212fd5f022007ed037ab7da039e0d4cf0fa35f620f2d7f71285959dcb6c885652d843d0038741210232b357c5309644cf4aa72b9b2d8bfe58bdf2515d40119318d5cb51ef378cae7efffffffff91d2feda79506806c5b0dc74b6fa6ae42fb7963da460ac256383fce498b9952020000006b483045022100ac75defcda55d644b6c095c2e7cbded92e38ea52e4d7f561f37962aa036a92fe0220059c0071d7c53cc964f641fa60ba2a076e3030ae6edc0f650792041c7691c66641210282d7e568e56f59e01a4edae297ac26caabc4684971ac6c7558c91c0fa84002f7ffffffff03322a093f000000001976a914c4263eb96d88849f498d139424b59a0cba1005e888ac2f92bb09000000001976a9146cbff9881ac47da8cb699e4543c28f9b3d6941da88ac404b4c00000000001976a914f7899faf1696892e6cb029b00c713f044761f03588ac00000000") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
	}

	// Valid (with feedback)
	if strings.Contains(data[0], "0100000001d1bda0bde67183817b21af863adaa31fda8cafcf2083ca1eaba3054496cbde10010000006a47304402205fddd6abab6b8e94f36bfec51ba2e1f3a91b5327efa88264b5530d0c86538723022010e51693e3d52347d4d2ff142b85b460d3953e625d1e062a5fa2569623fb0ea94121029df3723daceb1fef64fa0558371bc48cc3a7a8e35d8e05b87137dc129a9d4598ffffffff0115d40000000000001976a91459cc95a8cde59ceda718dbf70e612dba4034552688ac00000000") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"statusUrl":"https://api.whatsonchain.com/v1/bsv/tx/broadcast/cxF3xOdSvR_JgoXWDYZ0RQ"}`)))
	}

	// Invalid - error
	if strings.Contains(data[0], "error") {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`unknown error`)))
		return resp, fmt.Errorf("unknown error")
	}

	// Default is valid
	return resp, nil
}

// TestClient_GetTxByHash tests the GetTxByHash()
func TestClient_GetTxByHash(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPTransactions{})

	// Create the list of tests
	var tests = []struct {
		input         string
		expected      string
		expectedError bool
		statusCode    int
	}{
		{"c1d32f28baa27a376ba977f6a8de6ce0a87041157cef0274b20bfda2b0d8df96", "c1d32f28baa27a376ba977f6a8de6ce0a87041157cef0274b20bfda2b0d8df96", false, http.StatusOK},
		{"error", "", true, http.StatusInternalServerError},
	}

	// Test all
	for _, test := range tests {
		if output, err := client.GetTxByHash(test.input); err == nil && test.expectedError {
			t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted and [%s] expected", t.Name(), test.input, test.expected)
		} else if err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%v] error [%s]", t.Name(), test.input, test.expected, output, err.Error())
		} else if output != nil && output.TxID != test.expected && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%s]", t.Name(), test.input, test.expected, output.TxID)
		} else if client.LastRequest.StatusCode != test.statusCode {
			t.Errorf("%s Expected status code to be %d, got %d, [%s] inputted", t.Name(), test.statusCode, client.LastRequest.StatusCode, test.input)
		}
	}
}

// TestClient_GetMerkleProof tests the GetMerkleProof()
func TestClient_GetMerkleProof(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPTransactions{})

	// Create the list of tests
	var tests = []struct {
		input         string
		blockHash     string
		merkleRoot    string
		expectedError bool
		statusCode    int
	}{
		{"c1d32f28baa27a376ba977f6a8de6ce0a87041157cef0274b20bfda2b0d8df96", "0000000000000000091216c46973d82db057a6f9911352892b7769ed517681c3", "95a920b1002bed05379a0d2650bb13eb216138f28ee80172f4cf21048528dc60", false, http.StatusOK},
		{"error", "", "", true, http.StatusBadRequest},
		{"c1d32f28baa27a376ba977f6a8de6ce0a87041157cef0274b20bfda2b0d8dfzz", "", "", false, http.StatusOK},
	}

	// Test all
	for _, test := range tests {
		if output, err := client.GetMerkleProof(test.input); err == nil && test.expectedError {
			t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted", t.Name(), test.input)
		} else if err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted, received: [%v] error [%s]", t.Name(), test.input, output, err.Error())
		} else if output != nil && output[0].BlockHash != test.blockHash && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] blockhash expected, received: [%s]", t.Name(), test.input, test.blockHash, output[0].BlockHash)
		} else if output != nil && output[0].MerkleRoot != test.merkleRoot && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] merkle root expected, received: [%s]", t.Name(), test.input, test.merkleRoot, output[0].MerkleRoot)
		} else if client.LastRequest.StatusCode != test.statusCode {
			t.Errorf("%s Expected status code to be %d, got %d, [%s] inputted", t.Name(), test.statusCode, client.LastRequest.StatusCode, test.input)
		}
	}
}

// TestClient_GetRawTransactionData tests the GetRawTransactionData()
func TestClient_GetRawTransactionData(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPTransactions{})

	// Create the list of tests
	var tests = []struct {
		input         string
		expected      string
		expectedError bool
		statusCode    int
	}{
		{"c1d32f28baa27a376ba977f6a8de6ce0a87041157cef0274b20bfda2b0d8df96", "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff1c03d7c6082f7376706f6f6c2e636f6d2f3edff034600055b8467f0040ffffffff01247e814a000000001976a914492558fb8ca71a3591316d095afc0f20ef7d42f788ac00000000", false, http.StatusOK},
		{"c1d32f28baa27a376ba977f6a8de6ce0a87041157cef0274b20bfda2b0d8dfzz", "", false, http.StatusNotFound},
		{"error", "", true, http.StatusBadRequest},
	}

	// Test all
	for _, test := range tests {
		if output, err := client.GetRawTransactionData(test.input); err == nil && test.expectedError {
			t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted and [%s] expected", t.Name(), test.input, test.expected)
		} else if err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%v] error [%s]", t.Name(), test.input, test.expected, output, err.Error())
		} else if output != test.expected && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%s]", t.Name(), test.input, test.expected, output)
		} else if client.LastRequest.StatusCode != test.statusCode {
			t.Errorf("%s Expected status code to be %d, got %d, [%s] inputted", t.Name(), test.statusCode, client.LastRequest.StatusCode, test.input)
		}
	}
}

// TestClient_GetRawTransactionOutputData tests the GetRawTransactionOutputData()
func TestClient_GetRawTransactionOutputData(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPTransactions{})

	// Create the list of tests
	var tests = []struct {
		input         string
		expected      string
		expectedError bool
		statusCode    int
	}{
		{"c1d32f28baa27a376ba977f6a8de6ce0a87041157cef0274b20bfda2b0d8df96", "76a914492558fb8ca71a3591316d095afc0f20ef7d42f788ac", false, http.StatusOK},
		{"c1d32f28baa27a376ba977f6a8de6ce0a87041157cef0274b20bfda2b0d8dfzz", "", false, http.StatusBadGateway},
		{"error", "", true, http.StatusBadRequest},
	}

	// Test all
	for _, test := range tests {
		if output, err := client.GetRawTransactionOutputData(test.input, 0); err == nil && test.expectedError {
			t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted and [%s] expected", t.Name(), test.input, test.expected)
		} else if err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%v] error [%s]", t.Name(), test.input, test.expected, output, err.Error())
		} else if output != test.expected && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%s]", t.Name(), test.input, test.expected, output)
		} else if client.LastRequest.StatusCode != test.statusCode {
			t.Errorf("%s Expected status code to be %d, got %d, [%s] inputted", t.Name(), test.statusCode, client.LastRequest.StatusCode, test.input)
		}
	}
}

// TestClient_BulkTransactionDetails tests the BulkTransactionDetails()
func TestClient_BulkTransactionDetails(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPTransactions{})

	// Create the list of tests
	var tests = []struct {
		input         *TxHashes
		tx1           string
		tx2           string
		expectedError bool
		statusCode    int
	}{
		{&TxHashes{TxIDs: []string{"294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1aa", "91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba3258"}}, "294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1aa", "91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba3258", false, http.StatusOK},
		{&TxHashes{TxIDs: []string{"294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1ZZ", "91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba3258"}}, "91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba3258", "", false, http.StatusOK},
		{&TxHashes{TxIDs: []string{"294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1VV", "91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba32VV"}}, "", "", false, http.StatusOK},
		{&TxHashes{TxIDs: []string{"294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1VV", "91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba32VV", "294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1VV", "91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba32VV", "294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1VV", "91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba32VV", "294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1VV", "91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba32VV", "294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1VV", "91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba32VV", "294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1VV", "91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba32VV", "294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1VV", "91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba32VV", "294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1VV", "91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba32VV", "294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1VV", "91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba32VV", "294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1VV", "91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba32VV", "294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1VV", "91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba32VV"}}, "", "", true, http.StatusOK},
		{&TxHashes{TxIDs: []string{"error"}}, "", "", true, http.StatusBadRequest},
	}

	// Test all
	for _, test := range tests {
		if output, err := client.BulkTransactionDetails(test.input); err == nil && test.expectedError {
			t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted", t.Name(), test.input)
		} else if err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted, received: [%v] error [%s]", t.Name(), test.input, output, err.Error())
		} else if output != nil && len(output) >= 1 && output[0].TxID != test.tx1 && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%s]", t.Name(), test.input, test.tx1, output[0].TxID)
		} else if output != nil && len(output) >= 2 && output[1].TxID != test.tx2 && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%s]", t.Name(), test.input, test.tx2, output[1].TxID)
		} else if client.LastRequest.StatusCode != test.statusCode {
			t.Errorf("%s Expected status code to be %d, got %d, [%s] inputted", t.Name(), test.statusCode, client.LastRequest.StatusCode, test.input)
		}
	}
}

// TestClient_BroadcastTx tests the BroadcastTx()
func TestClient_BroadcastTx(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPBroadcast{})

	// Create the list of tests
	var tests = []struct {
		input         string
		expected      string
		expectedError bool
		statusCode    int
	}{
		{"0100000001ace8d1d1d885039e2d58074188e1355f03ea782b382a2890ee4ed2fc76e51daf000000006a47304402206c408d858d8678666f112d07c6a6152121e9672ee6af45ea386592166efd0852022012a94bdba5d3efc66c695dad0fb5a3267a599f5ef4f144d10663aba34ac3f8234121036e4a5eff4cb231e0a3d34932d950a2a9f6476a01b41f38e8ddf036eae2cc9fcfffffffff0280d03b00000000001976a9146680fd90c9d68cb9cf5314c4e30fb5a3879440c988acbb360100000000001976a914847e32a2b2e0289912c3d2d0eed20feb0a8bf4c688ac00000000", "45ed612b6a2e819af164eb55e8273b61f3bfbbf08b991db1bd1d9b48565a3297", false, http.StatusOK},
		{"0100000001d1bda0bde67183817b21af863adaa31fda8cafcf2083ca1eaba3054496cbde10010000006a47304402205fddd6abab6b8e94f36bfec51ba2e1f3a91b5327efa88264b5530d0c86538723022010e51693e3d52347d4d2ff142b85b460d3953e625d1e062a5fa2569623fb0ea94121029df3723daceb1fef64fa0558371bc48cc3a7a8e35d8e05b87137dc129a9d4598ffffffff0115d40000000000001976a91459cc95a8cde59ceda718dbf70e612dba4034552688ac00000000", "", true, http.StatusBadRequest},
		{"010000000220f177bc9f9d6311bf5a3ab3f95f1fb73904e91faec69e11f1a5780d717fb9da0e0400006b483045022100944ff1436e24ed37a0d644d6637ef36e1f79f217aefef2ab5c1fccbc474ccd2002200b0fb32081b31a27359f7e10760fe0c303a5877f49b323bd81629555f62b3987412103390e121888d94987dd77d74a40e643eaa24ae61410533046b363d602096312ceffffffff938d7f6d2600707128aedfe761b5895fa87aa72838f50efe9d8771ace92fbc6b020000006b483045022100ddecfb3b1caf072c2ed60cb3d032266f39e8858b03b3e3cd4b59ea881b20736b022052410db489955bd829878eae7d71489b30336acc3b716029046c1659615f8a76412103a38ede6ca418dce302f5241aa8a0b274046807f5527d2b4bf67d93cbf0c7fabaffffffff02d0b5b384000000001976a914501a793f77e008bdff95797738bed6c848e2bc3588ac68f032550d0000001976a91423cb7559a575b8f6d049372cba0f7cb4cfcdac3988ac00000000", "", true, http.StatusBadRequest},
		{"bad-status-code", "", true, http.StatusExpectationFailed},
	}

	// Test all
	for _, test := range tests {
		if output, err := client.BroadcastTx(test.input); err == nil && test.expectedError {
			t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted and [%s] expected", t.Name(), test.input, test.expected)
		} else if err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%v] error [%s]", t.Name(), test.input, test.expected, output, err.Error())
		} else if output != test.expected && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%s]", t.Name(), test.input, test.expected, output)
		} else if client.LastRequest.StatusCode != test.statusCode {
			t.Errorf("%s Expected status code to be %d, got %d, [%s] inputted", t.Name(), test.statusCode, client.LastRequest.StatusCode, test.input)
		}
	}
}

// TestClient_BulkBroadcastTx tests the BulkBroadcastTx()
func TestClient_BulkBroadcastTx(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPBroadcastBulk{})

	// Create the list of tests
	var tests = []struct {
		input         []string
		expected      string
		feedback      bool
		expectedError bool
		statusCode    int
	}{
		{[]string{"020000000232900e9b0e359cb95ad3853e1450591fdf01e3efaa1b3b2b5ab5c5ef784946b8010000006a473044022055ec4f9b9cbdd97cf5f4893f921da75287483edb4ebba4f5b23231577212fd5f022007ed037ab7da039e0d4cf0fa35f620f2d7f71285959dcb6c885652d843d0038741210232b357c5309644cf4aa72b9b2d8bfe58bdf2515d40119318d5cb51ef378cae7efffffffff91d2feda79506806c5b0dc74b6fa6ae42fb7963da460ac256383fce498b9952020000006b483045022100ac75defcda55d644b6c095c2e7cbded92e38ea52e4d7f561f37962aa036a92fe0220059c0071d7c53cc964f641fa60ba2a076e3030ae6edc0f650792041c7691c66641210282d7e568e56f59e01a4edae297ac26caabc4684971ac6c7558c91c0fa84002f7ffffffff03322a093f000000001976a914c4263eb96d88849f498d139424b59a0cba1005e888ac2f92bb09000000001976a9146cbff9881ac47da8cb699e4543c28f9b3d6941da88ac404b4c00000000001976a914f7899faf1696892e6cb029b00c713f044761f03588ac00000000", "0100000006279fdaff3d61920bc0b017c7b7ed8a8944b1ed1cc1998efc18b71f4828dd3a02010000006a47304402205ba16dccd88461f11e0a705b1015ae84d543acc92cb25012c14578d14cedf677022036bfdd96c6bbe57488b288a1dd3c1c4208154066fb7f8f9dfe59fb6ef53ca01b41210215b9d8176b6697758859e95ff43e61dfea94d44340b16d6f6ac4ae6d61a37ab3ffffffff3efb67c48b0a387f769219ee5c7431eaa49e11a3bfb4e7969d91b0973bcc3a30010000006b483045022100d1a25279f2bb720848717b85c518018c3a26ac584dce2956aeaa4ead86e43b1002200c2e5193fcc2250cc08407b22e7ee73ba48c8755d2aa193b543d01ad97dff85941210306e0678608241a4dc0bde0a39ab3d29dfdef6624ab5a20aa0b821dea85c9d9d4ffffffff9da533cd621baa3d07aca98f07c03115db993bf685afad6a3f321dac61db1731010000006a473044022009769824e84a1b9756aa9450249c19604ffe1474875e9389e851c562dcf272e502206dec4452cdec705bb4066883a2b9cd571768df2017638143bdec4a053b37abcb41210243b052e875c67900cebb6bcf816c7ee510f1a7c1068b0625a1828f1ab82c3806ffffffff4e6647164ee0edc0c7617e82575a88eca16056efb3af297e32fd03767c4e353d010000006a473044022012949fd73deca106804968fd50d86cb89872758f5c492b0938fdb8cee28b30df0220168f88d65400c68c5df348d60c800b1f9fb7e3135af05addac049c940da36bf9412102645b8993f1a183e9d37ec2fd9c5f950110b404aea7e9f01687ab42eb6f8b3563ffffffff49236581eb2ada13ac5972e7211771539ebe06061ac6a27a58acde5521b80949010000006a4730440220621f6dff81bfc0ab6c2b9e321c920f477f4ee5e5a01e3d24aecbb265ad264a1202205e99ac2fe08f6cecfcc8ddc98de09d4693c7f81f03bb478d1ee5936919e00b0d4121039925d3b23e560bd021d665e56f912735fcbb96303051b037a6816a55697a226effffffff1b6196327f18bbe6af3fabd46a5d3af568cc7c4447a340f4949824b6762a98f2010000006b483045022100f679c4408b7d893a6bbfb042685053369de044d4d12c6ddcc678d223684f7a320220669bd27d9e9e2be82c0f653d95cd74a4e188a66b976c5be7ba7d3b113692547041210235157690d4fc237b4aceb9e6f91ae0f0058628df36e4450f173ffae2d7d0d49affffffff0253290000000000001976a914dc57ab8a8365a7263fad7491e9a36601a786772388ac37d60000000000001976a914594d5717bd8f9ae5ae8c56646042082d3d6995f988ac00000000"}, "", false, false, http.StatusOK},
		{[]string{"0100000001d1bda0bde67183817b21af863adaa31fda8cafcf2083ca1eaba3054496cbde10010000006a47304402205fddd6abab6b8e94f36bfec51ba2e1f3a91b5327efa88264b5530d0c86538723022010e51693e3d52347d4d2ff142b85b460d3953e625d1e062a5fa2569623fb0ea94121029df3723daceb1fef64fa0558371bc48cc3a7a8e35d8e05b87137dc129a9d4598ffffffff0115d40000000000001976a91459cc95a8cde59ceda718dbf70e612dba4034552688ac00000000", "0100000006279fdaff3d61920bc0b017c7b7ed8a8944b1ed1cc1998efc18b71f4828dd3a02010000006a47304402205ba16dccd88461f11e0a705b1015ae84d543acc92cb25012c14578d14cedf677022036bfdd96c6bbe57488b288a1dd3c1c4208154066fb7f8f9dfe59fb6ef53ca01b41210215b9d8176b6697758859e95ff43e61dfea94d44340b16d6f6ac4ae6d61a37ab3ffffffff3efb67c48b0a387f769219ee5c7431eaa49e11a3bfb4e7969d91b0973bcc3a30010000006b483045022100d1a25279f2bb720848717b85c518018c3a26ac584dce2956aeaa4ead86e43b1002200c2e5193fcc2250cc08407b22e7ee73ba48c8755d2aa193b543d01ad97dff85941210306e0678608241a4dc0bde0a39ab3d29dfdef6624ab5a20aa0b821dea85c9d9d4ffffffff9da533cd621baa3d07aca98f07c03115db993bf685afad6a3f321dac61db1731010000006a473044022009769824e84a1b9756aa9450249c19604ffe1474875e9389e851c562dcf272e502206dec4452cdec705bb4066883a2b9cd571768df2017638143bdec4a053b37abcb41210243b052e875c67900cebb6bcf816c7ee510f1a7c1068b0625a1828f1ab82c3806ffffffff4e6647164ee0edc0c7617e82575a88eca16056efb3af297e32fd03767c4e353d010000006a473044022012949fd73deca106804968fd50d86cb89872758f5c492b0938fdb8cee28b30df0220168f88d65400c68c5df348d60c800b1f9fb7e3135af05addac049c940da36bf9412102645b8993f1a183e9d37ec2fd9c5f950110b404aea7e9f01687ab42eb6f8b3563ffffffff49236581eb2ada13ac5972e7211771539ebe06061ac6a27a58acde5521b80949010000006a4730440220621f6dff81bfc0ab6c2b9e321c920f477f4ee5e5a01e3d24aecbb265ad264a1202205e99ac2fe08f6cecfcc8ddc98de09d4693c7f81f03bb478d1ee5936919e00b0d4121039925d3b23e560bd021d665e56f912735fcbb96303051b037a6816a55697a226effffffff1b6196327f18bbe6af3fabd46a5d3af568cc7c4447a340f4949824b6762a98f2010000006b483045022100f679c4408b7d893a6bbfb042685053369de044d4d12c6ddcc678d223684f7a320220669bd27d9e9e2be82c0f653d95cd74a4e188a66b976c5be7ba7d3b113692547041210235157690d4fc237b4aceb9e6f91ae0f0058628df36e4450f173ffae2d7d0d49affffffff0253290000000000001976a914dc57ab8a8365a7263fad7491e9a36601a786772388ac37d60000000000001976a914594d5717bd8f9ae5ae8c56646042082d3d6995f988ac00000000"}, "https://api.whatsonchain.com/v1/bsv/tx/broadcast/cxF3xOdSvR_JgoXWDYZ0RQ", true, false, http.StatusOK},
		{[]string{"error", "error2"}, "", true, true, http.StatusBadRequest},
	}

	// Test all
	for _, test := range tests {
		if output, err := client.BulkBroadcastTx(test.input, test.feedback); err == nil && test.expectedError {
			t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted and [%s] expected", t.Name(), test.input, test.expected)
		} else if err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%v] error [%s]", t.Name(), test.input, test.expected, output, err.Error())
		} else if output != nil && output.StatusURL != test.expected && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%v]", t.Name(), test.input, test.expected, output)
		} else if client.LastRequest.StatusCode != test.statusCode {
			t.Errorf("%s Expected status code to be %d, got %d, [%s] inputted", t.Name(), test.statusCode, client.LastRequest.StatusCode, test.input)
		}
	}
}

// TestClient_BulkBroadcastTxMaxValues tests the BulkBroadcastTx()
func TestClient_BulkBroadcastTxMaxValues(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPBroadcastBulk{})

	// Create an input with more than the max txs
	var bulkTransactions []string
	for i := 1; i < MaxBroadcastTransactions+2; i++ {
		bulkTransactions = append(bulkTransactions, "0100000001d1bda0bde67183817b21af863adaa31fda8cafcf2083ca1eaba3054496cbde10010000006a47304402205fddd6abab6b8e94f36bfec51ba2e1f3a91b5327efa88264b5530d0c86538723022010e51693e3d52347d4d2ff142b85b460d3953e625d1e062a5fa2569623fb0ea94121029df3723daceb1fef64fa0558371bc48cc3a7a8e35d8e05b87137dc129a9d4598ffffffff0115d40000000000001976a91459cc95a8cde59ceda718dbf70e612dba4034552688ac00000000")
	}

	// Test the max
	_, err := client.BulkBroadcastTx(bulkTransactions, true)
	if err == nil {
		t.Errorf("%s Failed: expected to throw an error, no error, total txs %d", t.Name(), len(bulkTransactions))
	}

	// Test the max size of an allowed tx
	maxSizeTx := []string{""}
	txString := ""
	for i := 1; i < MaxSingleTransactionSize+2; i++ {
		txString = txString + "AA"
	}
	maxSizeTx = append(maxSizeTx, txString)

	// Test the max
	_, err = client.BulkBroadcastTx(maxSizeTx, true)
	if err == nil {
		t.Errorf("%s Failed: expected to throw an error, no error, total txs %d", t.Name(), len(bulkTransactions))
	}

	// Test the max size of all txs
	for i := 1; i < MaxBroadcastTransactions-1; i++ {
		maxSizeTx = append(maxSizeTx, txString)
	}

	// Test the max
	_, err = client.BulkBroadcastTx(maxSizeTx, true)
	if err == nil {
		t.Errorf("%s Failed: expected to throw an error, no error, total txs %d", t.Name(), len(bulkTransactions))
	}
}

// TestClient_DecodeTransaction tests the DecodeTransaction()
func TestClient_DecodeTransaction(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPBroadcast{})

	// Create the list of tests
	var tests = []struct {
		input         string
		expected      string
		expectedError bool
		statusCode    int
	}{
		{"010000000110784fd521b55a303da0f8b4ea113a2a3b5fa71565bab86c4257cff83ab4a1b9010000006a4730440220113a56d87122f28d6b60931498951f9709527a5c095ae852dc7b17d3d7915ef802206fd0b026d2e8dd30a39daaef17f503d985e2306c8244bb0e09a957bf1c9b530441210269a7785783c12405a1eaecb3088a3d830ed7e2de6ac527f42374a55a8cc5aeadffffffff0266ba0200000000001976a914021e4ac858f0ee6e0dfdf4438857f602a900698988acde905606000000001976a914022a8c1a18378885db9054676f17a27f4219045e88ac00000000", "cfcd9c342411592319442f705f9083938847e88f709096aa2cec15e6350b947e", false, http.StatusOK},
		{"zzzz0784fd521b55a303da0f8b4ea113a2a3b5fa71565bab86c4257cff83ab4a1b9010000006a4730440220113a56d87122f28d6b60931498951f9709527a5c095ae852dc7b17d3d7915ef802206fd0b026d2e8dd30a39daaef17f503d985e2306c8244bb0e09a957bf1c9b530441210269a7785783c12405a1eaecb3088a3d830ed7e2de6ac527f42374a55a8cc5aeadffffffff0266ba0200000000001976a914021e4ac858f0ee6e0dfdf4438857f602a900698988acde905606000000001976a914022a8c1a18378885db9054676f17a27f4219045e88ac00000000", "", true, http.StatusBadRequest},
	}

	// Test all
	for _, test := range tests {
		if output, err := client.DecodeTransaction(test.input); err == nil && test.expectedError {
			t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted and [%s] expected", t.Name(), test.input, test.expected)
		} else if err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%v] error [%s]", t.Name(), test.input, test.expected, output, err.Error())
		} else if output != nil && output.TxID != test.expected && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%s]", t.Name(), test.input, test.expected, output.TxID)
		} else if client.LastRequest.StatusCode != test.statusCode {
			t.Errorf("%s Expected status code to be %d, got %d, [%s] inputted", t.Name(), test.statusCode, client.LastRequest.StatusCode, test.input)
		}
	}
}

// TestClient_DownloadReceipt tests the DownloadReceipt()
func TestClient_DownloadReceipt(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPTransactions{})

	// Create the list of tests
	var tests = []struct {
		input         string
		expected      string
		expectedError bool
		statusCode    int
	}{
		{"c1d32f28baa27a376ba977f6a8de6ce0a87041157cef0274b20bfda2b0d8df96", "PDF", false, http.StatusOK},
		{"invalid", "invalid", true, http.StatusGatewayTimeout},
	}

	// Test all
	for _, test := range tests {
		if output, err := client.DownloadReceipt(test.input); err == nil && test.expectedError {
			t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted", t.Name(), test.input)
		} else if err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted, received: [%v] error [%s]", t.Name(), test.input, output, err.Error())
		} else if !strings.Contains(output, test.expected) && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%s]", t.Name(), test.input, test.expected, output)
		} else if client.LastRequest.StatusCode != test.statusCode {
			t.Errorf("%s Expected status code to be %d, got %d, [%s] inputted", t.Name(), test.statusCode, client.LastRequest.StatusCode, test.input)
		}
	}
}

// TestClient_BulkTransactionDetailsProcessor tests the BulkTransactionDetailsProcessor()
func TestClient_BulkTransactionDetailsProcessor(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPTransactions{})

	var tests = []struct {
		name          string
		input         *TxHashes
		tx1           string
		tx2           string
		expectedError bool
		statusCode    int
	}{
		{"valid transactions",
			&TxHashes{TxIDs: []string{
				"294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1aa",
				"91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba3258",
			}},
			"294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1aa",
			"91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba3258",
			false,
			http.StatusOK,
		},
		{"one real tx, one wrong",
			&TxHashes{TxIDs: []string{
				"294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1ZZ",
				"91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba3258",
			}},
			"91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba3258",
			"",
			false,
			http.StatusOK,
		},
		{"both txs are not found",
			&TxHashes{TxIDs: []string{
				"294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1VV",
				"91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba32VV",
			}},
			"",
			"",
			false,
			http.StatusOK},
		{"using 20 transactions",
			&TxHashes{TxIDs: []string{
				"294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1VV",
				"91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba32VV",
				"294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1VV",
				"91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba32VV",
				"294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1VV",
				"91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba32VV",
				"294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1VV",
				"91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba32VV",
				"294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1VV",
				"91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba32VV",
				"294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1VV",
				"91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba32VV",
				"294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1VV",
				"91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba32VV",
				"294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1VV",
				"91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba32VV",
				"294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1VV",
				"91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba32VV",
				"294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1VV",
				"91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba32VV",
				"294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1VV",
				"91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba32VV",
			}},
			"",
			"",
			false,
			http.StatusOK,
		},
		{"invalid tx",
			&TxHashes{TxIDs: []string{
				"error",
			}},
			"",
			"",
			true,
			http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if output, err := client.BulkTransactionDetailsProcessor(test.input); err == nil && test.expectedError {
				t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted", t.Name(), test.input)
			} else if err != nil && !test.expectedError {
				t.Errorf("%s Failed: [%s] inputted, received: [%v] error [%s]", t.Name(), test.input, output, err.Error())
			} else if output != nil && len(output) >= 1 && output[0].TxID != test.tx1 && !test.expectedError {
				t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%s]", t.Name(), test.input, test.tx1, output[0].TxID)
			} else if output != nil && len(output) >= 2 && output[1].TxID != test.tx2 && !test.expectedError {
				t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%s]", t.Name(), test.input, test.tx2, output[1].TxID)
			} else if client.LastRequest.StatusCode != test.statusCode {
				t.Errorf("%s Expected status code to be %d, got %d, [%s] inputted", t.Name(), test.statusCode, client.LastRequest.StatusCode, test.input)
			}
		})
	}
}
