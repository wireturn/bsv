package whatsonchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockHTTPAddresses for mocking requests
type mockHTTPAddresses struct{}

// Do is a mock http request
func (m *mockHTTPAddresses) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	//
	// Address Info
	//

	// Valid (info)
	if strings.Contains(req.URL.String(), "/16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA/info") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"isvalid": true,"address": "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA","scriptPubKey": "76a9143d0e5368bdadddca108a0fe44739919274c726c788ac","ismine": false,"iswatchonly": false,"isscript": false}`)))
	}

	// Invalid (info) return an error
	if strings.Contains(req.URL.String(), "/error/info") {
		resp.StatusCode = http.StatusInternalServerError
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
		return resp, fmt.Errorf("missing request")
	}

	// Valid (but invalid bsv address)
	if strings.Contains(req.URL.String(), "/16ZqP5invalid/info") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"isvalid": false,"address": "","scriptPubKey": "","ismine": false,"iswatchonly": false,"isscript": false}`)))
	}

	//
	// Address Balance
	//

	// Valid (balance)
	if strings.Contains(req.URL.String(), "/16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA/balance") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"confirmed": 10102050381,"unconfirmed": 123}`)))
	}

	// Invalid (balance) return an error
	if strings.Contains(req.URL.String(), "/16ZqP5invalid/balance") {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
		return resp, fmt.Errorf("bad request")
	}

	//
	// Address History
	//

	// Valid (history)
	if strings.Contains(req.URL.String(), "/16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA/history") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`[{"tx_hash": "6b22c47e7956e5404e05c3dc87dc9f46e929acfd46c8dd7813a34e1218d2f9d1","height": 563052},{"tx_hash": "1c312435789754392f92ffcb64e1248e17da47bed179abfd27e6003c775e0e04","height": 565076}]`)))
	}

	// Valid (history) (no results)
	if strings.Contains(req.URL.String(), "/1NfHy82RqJVGEau9u5DwFRyGc6QKwDuQeT/history") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`[]`)))
	}

	// Invalid (history) return an error
	if strings.Contains(req.URL.String(), "/16ZqP5invalid/history") {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
		return resp, fmt.Errorf("bad request")
	}

	//
	// Address unspent
	//

	// Valid (unspent)
	if strings.Contains(req.URL.String(), "/16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA/unspent") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`[{"height": 639302,"tx_pos": 3,"tx_hash": "33b9432a0ea203bbb6ec00592622cf6e90223849e4c9a76447a19a3ed43907d3","value": 2451680},{"height": 639601,"tx_pos": 3,"tx_hash": "4805041897a2ae59ffca85f0deb46e89d73d1ba4478bbd9c0fcd76ba0985ded2","value": 2744764},{"height": 640276,"tx_pos": 3,"tx_hash": "2493ff4cbca16b892ac641b7f2cb6d4388e75cb3f8963c291183f2bf0b27f415","value": 2568774}]`)))
	}

	// Valid (unspent) (no results)
	if strings.Contains(req.URL.String(), "/1NfHy82RqJVGEau9u5DwFRyGc6QKwDuQeT/unspent") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`[]`)))
	}

	// Invalid (unspent) return an error
	if strings.Contains(req.URL.String(), "/16ZqP5invalid/unspent") {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
		return resp, fmt.Errorf("bad request")
	}

	//
	// Address bulk balance
	//

	// Valid (unspent)
	if strings.Contains(req.URL.String(), "/addresses/balance") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`[{"address":"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP","balance":{"confirmed":0,"unconfirmed":0},"error":""},{"address":"1KGHhLTQaPr4LErrvbAuGE62yPpDoRwrob","balance":{"confirmed":301995631,"unconfirmed":0},"error":""}]`)))
	}

	//
	// Address bulk utxo
	//

	// Valid (unspent)
	if strings.Contains(req.URL.String(), "/addresses/unspent") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`[{"address":"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP","unspent":[],"error":""},{"address":"1KGHhLTQaPr4LErrvbAuGE62yPpDoRwrob","unspent":[{"height":658677,"tx_pos":1,"tx_hash":"be97e63bf79a961c69bc09d73cbef18232c7962fdced58244ed4014ba7e342b9","value":39799008},{"height":658726,"tx_pos":1,"tx_hash":"e5a7bc2338287fc0f3e38dff696b9ba41b3950c12bd8b7b1d92f3b0c056b4255","value":19110368},{"height":661623,"tx_pos":1,"tx_hash":"adea0a707c16f712f7a8faacfc8759d0b9d148693545c83511be1e2ed7fab4aa","value":19599008},{"height":661746,"tx_pos":1,"tx_hash":"afd619887ba5de9eb0e6076b7a37e96625227791520092fe142366c5c631c79e","value":44764416},{"height":661989,"tx_pos":1,"tx_hash":"15dcf82c9c461f3cb430e5ada855483c9e6c01bf4bc6fe667f3b798bd9f44acb","value":16658528},{"height":662494,"tx_pos":1,"tx_hash":"db8872fb1315e7f62013657d68db1871859624991a3ed77265aa85b8fdc768e5","value":32237986},{"height":662783,"tx_pos":4,"tx_hash":"2fa2686c61b6df1796717ca6d5f1934f0c39a5f8d2e42a6f213e76cb2ae66b54","value":10000000},{"height":662789,"tx_pos":4,"tx_hash":"c868a0616836bb017f956ce846ca6f3c56a985955742bd0fea22840a9d0168df","value":10000000},{"height":662791,"tx_pos":2,"tx_hash":"95c69649798b1d66e37318f2d65374095f6e0cd1675d1214402bd8e6002bf424","value":10000000},{"height":662794,"tx_pos":4,"tx_hash":"69720b7e41ca113d5fa988f5f4fd635d398459ba8e6ebb5d0a3a8f42097f1dcd","value":10000000},{"height":662853,"tx_pos":1,"tx_hash":"c42a4752f2551c195b27d016f8e522e724ad99b81d7e2459630b53e5a06178f3","value":5045746},{"height":662897,"tx_pos":1,"tx_hash":"2ca9c744b857a46266e4c0ac827db254eef54fa19530d3f21e460fe8d9445844","value":11405228},{"height":662992,"tx_pos":1,"tx_hash":"01c33159e9e00a7cb07248926e1ff8ed2d6a2450565fc0c27c30600457b2e572","value":47748859},{"height":663033,"tx_pos":1,"tx_hash":"09c5c72f807e572a5ac96e809d1c10b5bf27d63099cb4a6d871b74d459778bde","value":13629728},{"height":663034,"tx_pos":1,"tx_hash":"c8d3137f13ce2a4b8bfd919210c233a14a565c87e7c1ef4a693e6576adcc0419","value":7393008},{"height":663095,"tx_pos":1,"tx_hash":"58f416f323ae5b4d104b6246fca84ec4b1a6bb5a26174a732801e48008d02bbc","value":4603748}],"error":""}]`)))
	}

	//
	// Address download statement
	//

	// Valid (download statement)
	if strings.Contains(req.URL.String(), "/statement/16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA") {
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

	// Valid (download statement) (invalid address)
	if strings.Contains(req.URL.String(), "/statement/invalid") {
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
invalid
5 0 obj`)))
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

		// Valid - for AddressDetails
		if strings.Contains(data.TxIDs[0], "33b9432a0ea203bbb6ec00592622cf6e90223849e4c9a76447a19a3ed43907d3") {
			resp.StatusCode = http.StatusOK
			resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`[{"hex":"","txid":"33b9432a0ea203bbb6ec00592622cf6e90223849e4c9a76447a19a3ed43907d3","hash":"33b9432a0ea203bbb6ec00592622cf6e90223849e4c9a76447a19a3ed43907d3","version":1,"size":440,"locktime":0,"vin":[{"coinbase":"","txid":"fabe0b5d0979e068dce986692d1c5620f37383657a2fe7969f1cfe4a81b7f517","vout":3,"scriptSig":{"asm":"30450221008f74bb75c331cb7902a4e7539ee60fafe2c9a73d325aba6fc3ff9105ed91e219022064e65a5662c0593086ab05a0131e5abac5ef249f5f33c74351c2bed653da269f[ALL|FORKID] 026d6fc8f05b630e637507084b1678ec753c75b9e050312919e1d973224c5c3103","hex":"4830450221008f74bb75c331cb7902a4e7539ee60fafe2c9a73d325aba6fc3ff9105ed91e219022064e65a5662c0593086ab05a0131e5abac5ef249f5f33c74351c2bed653da269f4121026d6fc8f05b630e637507084b1678ec753c75b9e050312919e1d973224c5c3103"},"sequence":4294967295}],"vout":[{"value":0,"n":0,"scriptPubKey":{"asm":"0 OP_RETURN 3150755161374b36324d694b43747373534c4b79316b683536575755374d74555235 5522771 7368801 746f6e6963706f77 1701869940 6f666665725f636c69636b 6f666665725f636f6e6669675f6964 56 6f666665725f73657373696f6e5f6964 33316566313830633732363465303032373836333261306131613830313835313336363236336537306361383233353138373664636436386563666163623365","hex":"006a223150755161374b36324d694b43747373534c4b79316b683536575755374d74555235035345540361707008746f6e6963706f7704747970650b6f666665725f636c69636b0f6f666665725f636f6e6669675f69640138106f666665725f73657373696f6e5f69644033316566313830633732363465303032373836333261306131613830313835313336363236336537306361383233353138373664636436386563666163623365","type":"nulldata","opReturn":{"type":"bitcom","action":"","text":"","parts":null},"isTruncated":false}},{"value":0.00000549,"n":1,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 09cc4559bdcb84cb35c107743f0dbb10d66679cc OP_EQUALVERIFY OP_CHECKSIG","hex":"76a91409cc4559bdcb84cb35c107743f0dbb10d66679cc88ac","reqSigs":1,"type":"pubkeyhash","addresses":["1tonicZQwN2BNKhVwPXqh8ez3q56y1EYw"],"opReturn":null,"isTruncated":false}},{"value":0.00005489,"n":2,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 b5195bf7db0652f536a7dddbe36a99a091125468 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a914b5195bf7db0652f536a7dddbe36a99a09112546888ac","reqSigs":1,"type":"pubkeyhash","addresses":["1HWZgiMKQKPSkLzT7hipS22AvkQZJsyxmT"],"opReturn":null,"isTruncated":false}},{"value":0.0245168,"n":3,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 bf49c6a5406675e174f4f6a83b3d94dd9d845398 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a914bf49c6a5406675e174f4f6a83b3d94dd9d84539888ac","reqSigs":1,"type":"pubkeyhash","addresses":["1JSSSgcyufLgbXFw6WAXyXgBrmgFpnqXWh"],"opReturn":null,"isTruncated":false}}],"blockhash":"0000000000000000026b9da3860e4c8ee351a7af46da6042eaa5d110113b9fad","confirmations":1348,"time":1592122768,"blocktime":1592122768},{"hex":"","txid":"4805041897a2ae59ffca85f0deb46e89d73d1ba4478bbd9c0fcd76ba0985ded2","hash":"4805041897a2ae59ffca85f0deb46e89d73d1ba4478bbd9c0fcd76ba0985ded2","version":1,"size":439,"locktime":0,"vin":[{"coinbase":"","txid":"5a45b8415e5c1740353cfb011d29e04ec104865be6560dff5bd6cb31db75d559","vout":3,"scriptSig":{"asm":"3044022008e2417d072cfbb95d4e04c7e6e6ab70e415a379fb912cb2e0503e3df0ae0d2002201f9fcbf6c65ba6624fe0669d08155ed7c0d19c28be72daf3e00de2613656f955[ALL|FORKID] 026d6fc8f05b630e637507084b1678ec753c75b9e050312919e1d973224c5c3103","hex":"473044022008e2417d072cfbb95d4e04c7e6e6ab70e415a379fb912cb2e0503e3df0ae0d2002201f9fcbf6c65ba6624fe0669d08155ed7c0d19c28be72daf3e00de2613656f9554121026d6fc8f05b630e637507084b1678ec753c75b9e050312919e1d973224c5c3103"},"sequence":4294967295}],"vout":[{"value":0,"n":0,"scriptPubKey":{"asm":"0 OP_RETURN 3150755161374b36324d694b43747373534c4b79316b683536575755374d74555235 5522771 7368801 746f6e6963706f77 1701869940 6f666665725f636c69636b 6f666665725f636f6e6669675f6964 56 6f666665725f73657373696f6e5f6964 38313865386165656339353733646431333439373334366135363464633461623035353062333039383830373563393733316631643063653731336536353335","hex":"006a223150755161374b36324d694b43747373534c4b79316b683536575755374d74555235035345540361707008746f6e6963706f7704747970650b6f666665725f636c69636b0f6f666665725f636f6e6669675f69640138106f666665725f73657373696f6e5f69644038313865386165656339353733646431333439373334366135363464633461623035353062333039383830373563393733316631643063653731336536353335","type":"nulldata","opReturn":{"type":"bitcom","action":"","text":"","parts":null},"isTruncated":false}},{"value":0.00000573,"n":1,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 09cc4559bdcb84cb35c107743f0dbb10d66679cc OP_EQUALVERIFY OP_CHECKSIG","hex":"76a91409cc4559bdcb84cb35c107743f0dbb10d66679cc88ac","reqSigs":1,"type":"pubkeyhash","addresses":["1tonicZQwN2BNKhVwPXqh8ez3q56y1EYw"],"opReturn":null,"isTruncated":false}},{"value":0.00005726,"n":2,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 819ae5a5cbb078e96379b8eb25c29d6f7b28c412 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a914819ae5a5cbb078e96379b8eb25c29d6f7b28c41288ac","reqSigs":1,"type":"pubkeyhash","addresses":["1CpHjBbHoWzbrqQsPeZ39GLUXejZce9mBs"],"opReturn":null,"isTruncated":false}},{"value":0.02744764,"n":3,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 bf49c6a5406675e174f4f6a83b3d94dd9d845398 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a914bf49c6a5406675e174f4f6a83b3d94dd9d84539888ac","reqSigs":1,"type":"pubkeyhash","addresses":["1JSSSgcyufLgbXFw6WAXyXgBrmgFpnqXWh"],"opReturn":null,"isTruncated":false}}],"blockhash":"000000000000000003d684082ab45014f89a7f8e5e35ec94fcb4aa8b5f00c01e","confirmations":1049,"time":1592307236,"blocktime":1592307236},{"hex":"","txid":"2493ff4cbca16b892ac641b7f2cb6d4388e75cb3f8963c291183f2bf0b27f415","hash":"2493ff4cbca16b892ac641b7f2cb6d4388e75cb3f8963c291183f2bf0b27f415","version":1,"size":439,"locktime":0,"vin":[{"coinbase":"","txid":"2ebc8f094fdc012f7d9a0ed39999dcf0318665830f7d5f113af0d1c79fba2f8e","vout":3,"scriptSig":{"asm":"30440220010a62c1d79afcc274b8db821cba1f093c316d67d505a3900c231ae6dfb2dd51022031fe80787c531e1c890754d2cafdc624f3446e4d1bdca18ade83cabd3a2317ac[ALL|FORKID] 026d6fc8f05b630e637507084b1678ec753c75b9e050312919e1d973224c5c3103","hex":"4730440220010a62c1d79afcc274b8db821cba1f093c316d67d505a3900c231ae6dfb2dd51022031fe80787c531e1c890754d2cafdc624f3446e4d1bdca18ade83cabd3a2317ac4121026d6fc8f05b630e637507084b1678ec753c75b9e050312919e1d973224c5c3103"},"sequence":4294967295}],"vout":[{"value":0,"n":0,"scriptPubKey":{"asm":"0 OP_RETURN 3150755161374b36324d694b43747373534c4b79316b683536575755374d74555235 5522771 7368801 746f6e6963706f77 1701869940 6f666665725f636c69636b 6f666665725f636f6e6669675f6964 56 6f666665725f73657373696f6e5f6964 35656237306231653930306535616437626335663961333663653861643435623664336435636337666466393437343762623364326461663732636631356533","hex":"006a223150755161374b36324d694b43747373534c4b79316b683536575755374d74555235035345540361707008746f6e6963706f7704747970650b6f666665725f636c69636b0f6f666665725f636f6e6669675f69640138106f666665725f73657373696f6e5f69644035656237306231653930306535616437626335663961333663653861643435623664336435636337666466393437343762623364326461663732636631356533","type":"nulldata","opReturn":{"type":"bitcom","action":"","text":"","parts":null},"isTruncated":false}},{"value":0.00000572,"n":1,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 09cc4559bdcb84cb35c107743f0dbb10d66679cc OP_EQUALVERIFY OP_CHECKSIG","hex":"76a91409cc4559bdcb84cb35c107743f0dbb10d66679cc88ac","reqSigs":1,"type":"pubkeyhash","addresses":["1tonicZQwN2BNKhVwPXqh8ez3q56y1EYw"],"opReturn":null,"isTruncated":false}},{"value":0.00005716,"n":2,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 0405a52b27214920873fa222071a8ec9610317a4 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a9140405a52b27214920873fa222071a8ec9610317a488ac","reqSigs":1,"type":"pubkeyhash","addresses":["1NGU17f9HTyv3LffW4zxukSEwsxwf4d53"],"opReturn":null,"isTruncated":false}},{"value":0.02568774,"n":3,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 bf49c6a5406675e174f4f6a83b3d94dd9d845398 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a914bf49c6a5406675e174f4f6a83b3d94dd9d84539888ac","reqSigs":1,"type":"pubkeyhash","addresses":["1JSSSgcyufLgbXFw6WAXyXgBrmgFpnqXWh"],"opReturn":null,"isTruncated":false}}],"blockhash":"00000000000000000087222006199927280a010d0db21c6d253409f3e960c7bf","confirmations":374,"time":1592698834,"blocktime":1592698834}]`)))
		}
	}

	// Default is valid
	return resp, nil
}

// mockHTTPAddressesErrors for mocking requests
type mockHTTPAddressesErrors struct{}

// Do is a mock http request
func (m *mockHTTPAddressesErrors) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Invalid (info) return an error
	if strings.Contains(req.URL.String(), "/addresses/balance") {
		resp.StatusCode = http.StatusInternalServerError
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
		return resp, fmt.Errorf("missing request")
	}

	// Invalid (info) return an error
	if strings.Contains(req.URL.String(), "/addresses/unspent") {
		resp.StatusCode = http.StatusInternalServerError
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
		return resp, fmt.Errorf("missing request")
	}

	return nil, fmt.Errorf("no valid response found")
}

// TestClient_AddressInfo tests the AddressInfo()
func TestClient_AddressInfo(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPAddresses{})

	// Create the list of tests
	var tests = []struct {
		input         string
		expected      string
		expectedError bool
		statusCode    int
	}{
		{"16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA", "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA", false, http.StatusOK},
		{"16ZqP5invalid", "", false, http.StatusOK},
		{"error", "", true, http.StatusInternalServerError},
	}

	// Test all
	for _, test := range tests {
		if output, err := client.AddressInfo(test.input); err == nil && test.expectedError {
			t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted and [%s] expected", t.Name(), test.input, test.expected)
		} else if err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%v] error [%s]", t.Name(), test.input, test.expected, output, err.Error())
		} else if output != nil && output.Address != test.expected && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%s]", t.Name(), test.input, test.expected, output.Address)
		} else if client.LastRequest.StatusCode != test.statusCode {
			t.Errorf("%s Expected status code to be %d, got %d, [%s] inputted", t.Name(), test.statusCode, client.LastRequest.StatusCode, test.input)
		}
	}
}

// TestClient_AddressBalance tests the AddressBalance()
func TestClient_AddressBalance(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPAddresses{})

	// Create the list of tests
	var tests = []struct {
		input         string
		confirmed     int64
		unconfirmed   int64
		expectedError bool
		statusCode    int
	}{
		{"16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA", 10102050381, 123, false, http.StatusOK},
		{"16ZqP5invalid", 0, 0, true, http.StatusBadRequest},
	}

	// Test all
	for _, test := range tests {
		if output, err := client.AddressBalance(test.input); err == nil && test.expectedError {
			t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted", t.Name(), test.input)
		} else if err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted, received: [%v] error [%s]", t.Name(), test.input, output, err.Error())
		} else if output != nil && output.Confirmed != test.confirmed && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%d] confirm expected, received: [%d]", t.Name(), test.input, test.confirmed, output.Confirmed)
		} else if output != nil && output.Unconfirmed != test.unconfirmed && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%d] unconfirmed expected, received: [%d]", t.Name(), test.input, test.unconfirmed, output.Unconfirmed)
		} else if client.LastRequest.StatusCode != test.statusCode {
			t.Errorf("%s Expected status code to be %d, got %d, [%s] inputted", t.Name(), test.statusCode, client.LastRequest.StatusCode, test.input)
		}
	}
}

// TestClient_AddressHistory tests the AddressHistory()
func TestClient_AddressHistory(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPAddresses{})

	// Create the list of tests
	var tests = []struct {
		input         string
		txHash        string
		height        int64
		expectedError bool
		statusCode    int
	}{
		{"16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA", "6b22c47e7956e5404e05c3dc87dc9f46e929acfd46c8dd7813a34e1218d2f9d1", 563052, false, http.StatusOK},
		{"1NfHy82RqJVGEau9u5DwFRyGc6QKwDuQeT", "", 0, false, http.StatusOK},
		{"16ZqP5invalid", "", 0, true, http.StatusBadRequest},
	}

	// Test all
	for _, test := range tests {
		if output, err := client.AddressHistory(test.input); err == nil && test.expectedError {
			t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted", t.Name(), test.input)
		} else if err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted, received: [%v] error [%s]", t.Name(), test.input, output, err.Error())
		} else if output != nil && len(output) > 0 && output[0].TxHash != test.txHash && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] hash expected, received: [%s]", t.Name(), test.input, test.txHash, output[0].TxHash)
		} else if output != nil && len(output) > 0 && output[0].Height != test.height && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%d] height expected, received: [%d]", t.Name(), test.input, test.height, output[0].Height)
		} else if client.LastRequest.StatusCode != test.statusCode {
			t.Errorf("%s Expected status code to be %d, got %d, [%s] inputted", t.Name(), test.statusCode, client.LastRequest.StatusCode, test.input)
		}
	}
}

// TestClient_AddressUnspentTransactions tests the AddressUnspentTransactions()
func TestClient_AddressUnspentTransactions(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPAddresses{})

	// Create the list of tests
	var tests = []struct {
		input         string
		txHash        string
		height        int64
		value         int64
		expectedError bool
		statusCode    int
	}{
		{"16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA", "33b9432a0ea203bbb6ec00592622cf6e90223849e4c9a76447a19a3ed43907d3", 639302, 2451680, false, http.StatusOK},
		{"1NfHy82RqJVGEau9u5DwFRyGc6QKwDuQeT", "", 0, 0, false, http.StatusOK},
		{"16ZqP5invalid", "", 0, 0, true, http.StatusBadRequest},
	}

	// Test all
	for _, test := range tests {
		if output, err := client.AddressUnspentTransactions(test.input); err == nil && test.expectedError {
			t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted", t.Name(), test.input)
		} else if err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted, received: [%v] error [%s]", t.Name(), test.input, output, err.Error())
		} else if output != nil && len(output) > 0 && !test.expectedError {
			if output[0].TxHash != test.txHash {
				t.Errorf("%s Failed: [%s] inputted and [%s] hash expected, received: [%s]", t.Name(), test.input, test.txHash, output[0].TxHash)
			} else if output[0].Height != test.height {
				t.Errorf("%s Failed: [%s] inputted and [%d] height expected, received: [%d]", t.Name(), test.input, test.height, output[0].Height)
			} else if output[0].Value != test.value {
				t.Errorf("%s Failed: [%s] inputted and [%d] value expected, received: [%d]", t.Name(), test.input, test.value, output[0].Value)
			}
		} else if client.LastRequest.StatusCode != test.statusCode {
			t.Errorf("%s Expected status code to be %d, got %d, [%s] inputted", t.Name(), test.statusCode, client.LastRequest.StatusCode, test.input)
		}
	}
}

// TestClient_AddressUnspentTransactions tests the AddressUnspentTransactions()
func TestClient_AddressUnspentTransactionDetails(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPAddresses{})

	// Create the list of tests
	var tests = []struct {
		input         string
		txHash        string
		height        int64
		expectedError bool
		statusCode    int
	}{
		{"16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA", "33b9432a0ea203bbb6ec00592622cf6e90223849e4c9a76447a19a3ed43907d3", 639302, false, http.StatusOK},
	}

	// Test all
	for _, test := range tests {
		if output, err := client.AddressUnspentTransactionDetails(test.input, 5); err == nil && test.expectedError {
			t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted", t.Name(), test.input)
		} else if err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted, received: [%v] error [%s]", t.Name(), test.input, output, err.Error())
		} else if output != nil && len(output) > 0 && output[0].TxHash != test.txHash && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] hash expected, received: [%s]", t.Name(), test.input, test.txHash, output[0].TxHash)
		} else if output != nil && len(output) > 0 && output[0].Height != test.height && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%d] height expected, received: [%d]", t.Name(), test.input, test.height, output[0].Height)
		} else if client.LastRequest.StatusCode != test.statusCode {
			t.Errorf("%s Expected status code to be %d, got %d, [%s] inputted", t.Name(), test.statusCode, client.LastRequest.StatusCode, test.input)
		}
	}
}

// TestClient_DownloadStatement tests the DownloadStatement()
func TestClient_DownloadStatement(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPAddresses{})

	// Create the list of tests
	var tests = []struct {
		input         string
		expected      string
		expectedError bool
		statusCode    int
	}{
		{"16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA", "PDF", false, http.StatusOK},
		{"invalid", "invalid", false, http.StatusOK},
	}

	// Test all
	for _, test := range tests {
		if output, err := client.DownloadStatement(test.input); err == nil && test.expectedError {
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

// TestClient_BulkBalance tests the BulkBalance()
func TestClient_BulkBalance(t *testing.T) {
	t.Parallel()

	t.Run("valid response", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})

		balances, err := client.BulkBalance(&AddressList{Addresses: []string{"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP", "1KGHhLTQaPr4LErrvbAuGE62yPpDoRwrob"}})
		assert.NoError(t, err)
		assert.NotNil(t, balances)
		assert.Equal(t, 2, len(balances))
	})

	t.Run("max addresses (error)", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})

		balances, err := client.BulkBalance(&AddressList{Addresses: []string{
			"1",
			"2",
			"3",
			"4",
			"5",
			"6",
			"7",
			"8",
			"9",
			"10",
			"11",
			"12",
			"13",
			"14",
			"15",
			"16",
			"17",
			"18",
			"19",
			"20",
			"21",
		}})
		assert.Error(t, err)
		assert.Nil(t, balances)
	})

	t.Run("bad response (error)", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddressesErrors{})

		balances, err := client.BulkBalance(&AddressList{Addresses: []string{
			"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP",
		}})
		assert.Error(t, err)
		assert.Nil(t, balances)
	})
}

// TestClient_BulkUnspentTransactions tests the BulkUnspentTransactions()
func TestClient_BulkUnspentTransactions(t *testing.T) {
	t.Parallel()

	t.Run("valid response", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})

		balances, err := client.BulkUnspentTransactions(&AddressList{Addresses: []string{"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP", "1KGHhLTQaPr4LErrvbAuGE62yPpDoRwrob"}})
		assert.NoError(t, err)
		assert.NotNil(t, balances)
		assert.Equal(t, 2, len(balances))
	})

	t.Run("max addresses (error)", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})

		balances, err := client.BulkUnspentTransactions(&AddressList{Addresses: []string{
			"1",
			"2",
			"3",
			"4",
			"5",
			"6",
			"7",
			"8",
			"9",
			"10",
			"11",
			"12",
			"13",
			"14",
			"15",
			"16",
			"17",
			"18",
			"19",
			"20",
			"21",
		}})
		assert.Error(t, err)
		assert.Nil(t, balances)
	})

	t.Run("bad response (error)", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddressesErrors{})

		balances, err := client.BulkUnspentTransactions(&AddressList{Addresses: []string{
			"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP",
		}})
		assert.Error(t, err)
		assert.Nil(t, balances)
	})

}
