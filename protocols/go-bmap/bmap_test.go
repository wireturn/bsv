package bmap

import (
	"testing"

	"github.com/bitcoinschema/go-b"
	"github.com/bitcoinschema/go-bap"
	"github.com/bitcoinschema/go-bob"
	magic "github.com/bitcoinschema/go-map"
)

func TestFromBob(t *testing.T) {

	dataSource := `{"_id":"5ebc04c7814c6a17a6c90b3b","tx":{"h":"ce7429a101b7aecdf1e5449151d0be17a3948cb5c22282832ae942107edb2272"},"in":[{"i":0,"tape":[{"cell":[{"b":"MEQCIDUGRtDdmf2I2p1vcA2s4fMuBcmnSi5kOI2chSiFrYQKAiAq8XSIx8EbM2oKJJC9t/SFXTGnJBfE7mRKdGOVR7zIB0E=","s":"0D\u0002 5\u0006F�ݙ��ڝop\r���.\u0005ɧJ.d8���(���\n\u0002 *�t���\u001b3j\n$����]1�$\u0017��dJtc�G��\u0007A","ii":0,"i":0},{"b":"A/20DJgUWAXROgZTKDRmcC0ja306xpg3SiMTPy3QKhqQ","s":"\u0003��\f�\u0014X\u0005�:\u0006S(4fp-#k}:Ƙ7J#\u0013?-�*\u001a�","ii":1,"i":1}],"i":0}],"e":{"h":"f8448e73fc7667b91f86cf152b9bc4d88c365174989d79871e64ca8c66c1e785","i":0,"a":"1P1dKk7BCB6iTUz13w1eXkLfcj8a8dC4iv"},"seq":4294967295}],"out":[{"i":0,"tape":[{"cell":[{"op":0,"ops":"OP_0","ii":0,"i":0},{"op":106,"ops":"OP_RETURN","ii":1,"i":1}],"i":0},{"cell":[{"b":"MVB1UWE3SzYyTWlLQ3Rzc1NMS3kxa2g1NldXVTdNdFVSNQ==","s":"1PuQa7K62MiKCtssSLKy1kh56WWU7MtUR5","ii":2,"i":0},{"b":"U0VU","s":"SET","ii":3,"i":1},{"b":"YXBw","s":"app","ii":4,"i":2},{"b":"MnBheW1haWw=","s":"2paymail","ii":5,"i":3},{"b":"cGF5bWFpbA==","s":"paymail","ii":6,"i":4},{"b":"aGFnYmFyZEBtb25leWJ1dHRvbi5jb20=","s":"hagbard@moneybutton.com","ii":7,"i":5},{"b":"cHVibGljX2tleQ==","s":"public_key","ii":8,"i":6},{"b":"MDJjODliNjc5MGViNjA1MDYyYTMxZjEyNDI1MDU5NGJkMGZkMDI5ODhkYTI1NDFiM2QyNWU3ZWYzOTM3ZmI0YWUw","s":"02c89b6790eb605062a31f124250594bd0fd02988da2541b3d25e7ef3937fb4ae0","ii":9,"i":7},{"b":"cGxhdGZvcm0=","s":"platform","ii":10,"i":8},{"b":"dHdpdHRlcg==","s":"twitter","ii":11,"i":9},{"b":"cHJvb2ZfdXJs","s":"proof_url","ii":12,"i":10},{"b":"aHR0cHM6Ly90d2l0dGVyLmNvbS9oYWdiYXJkZGQvc3RhdHVzLzEyMDUxODk1ODAzMDkzNzcwMjQ=","s":"https://twitter.com/hagbarddd/status/1205189580309377024","ii":13,"i":11},{"b":"cHJvb2ZfYm9keQ==","s":"proof_body","ii":14,"i":12},{"b":"SGkKCk15IHBheW1haWwgaXMgaGFnYmFyZEBtb25leWJ1dHRvbi5jb20=","s":"Hi\n\nMy paymail is hagbard@moneybutton.com","ii":15,"i":13},{"b":"cHJvb2ZfaWQ=","s":"proof_id","ii":16,"i":14},{"b":"Sms5dlFncGREcG9XMHFEWQ==","s":"Jk9vQgpdDpoW0qDY","ii":17,"i":15}],"i":1},{"cell":[{"b":"MXNpZ255Q2l6cDFWeUJzSjVTczJ0RUFndzd6Q1lOSnU0","s":"1signyCizp1VyBsJ5Ss2tEAgw7zCYNJu4","ii":19,"i":0},{"b":"SU5LRmIxNU1uQVhxTlFueStiNEtBVm5HTnR5bUcwZEhTdTEzKzg3MSt0aTBXTjVGQmVBLzdEZ1VuMXRsdzZGN29kYlc3SURyVmVQS1RMclRQQWlEcXlvPQ==","s":"INKFb15MnAXqNQny+b4KAVnGNtymG0dHSu13+871+ti0WN5FBeA/7DgUn1tlw6F7odbW7IDrVePKTLrTPAiDqyo=","ii":20,"i":1},{"b":"MDJjODliNjc5MGViNjA1MDYyYTMxZjEyNDI1MDU5NGJkMGZkMDI5ODhkYTI1NDFiM2QyNWU3ZWYzOTM3ZmI0YWUw","s":"02c89b6790eb605062a31f124250594bd0fd02988da2541b3d25e7ef3937fb4ae0","ii":21,"i":2},{"b":"aGFnYmFyZEBtb25leWJ1dHRvbi5jb20=","s":"hagbard@moneybutton.com","ii":22,"i":3}],"i":2}],"e":{"v":0,"i":0,"a":"false"}},{"i":1,"tape":[{"cell":[{"op":118,"ops":"OP_DUP","ii":0,"i":0},{"op":169,"ops":"OP_HASH160","ii":1,"i":1},{"b":"7njcAgMt9eekcx5JZXFDoaThy9M=","s":"�x�\u0002\u0003-��s\u001eIeqC�����","ii":2,"i":2},{"op":136,"ops":"OP_EQUALVERIFY","ii":3,"i":3},{"op":172,"ops":"OP_CHECKSIG","ii":4,"i":4}],"i":0}],"e":{"v":31202,"i":1,"a":"1Njvc7dj8UHG6hnV5k5ZjSJtPgTofknDmx"}}],"lock":0,"blk":{"i":618112,"h":"000000000000000001e1e1f2995c9ba2e316f6fb85c247c923c591e56ea00fb6","t":1579328162},"i":478}`

	bobData, err := bob.NewFromString(dataSource)
	if err != nil {
		t.Fatalf("failed to create bob tx %s", err)
	}

	var bmapData *Tx
	if bmapData, err = NewFromBob(bobData); err != nil {
		t.Fatalf("error occurred: %s", err)
	}

	if bmapData.Tx.H != "ce7429a101b7aecdf1e5449151d0be17a3948cb5c22282832ae942107edb2272" {
		t.Fatalf("inherited field failed %+v", bmapData.MAP)
	}

	mapData := bmapData.MAP
	if mapData["app"] != "2paymail" {
		t.Fatalf("test fromBob failed %+v", mapData)
	}

}

func TestMap(t *testing.T) {
	tape := bob.Tape{
		Cell: []bob.Cell{
			{S: magic.Prefix},
			{S: magic.Set},
			{S: "keyName1"},
			{S: "something"},
			{S: "keyName2"},
			{S: "something else"},
		},
	}
	m, err := magic.NewFromTape(&tape)
	if err != nil {
		t.Fatalf("error occurred: %s", err)
	} else if m["keyName1"] != "something" {
		t.Fatalf("SET Failed %s", m["keyName1"])
	}
}

func TestB(t *testing.T) {
	tape := bob.Tape{
		Cell: []bob.Cell{
			{S: b.Prefix},
			{S: "Hello world"},
			{S: "text/plain"},
			{S: "utf8"},
		},
	}
	bTx, err := b.NewFromTape(tape)
	if err != nil {
		t.Fatalf("error occurred: %s", err)
	} else if bTx.Data.UTF8 != "Hello world" {
		t.Fatalf("Unexpected data %s %s", bTx.Data.UTF8, err)
	}
}

func TestNewFromBob(t *testing.T) {
	bobTx, err := bob.NewFromString(sampleValidBobTx)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}
	var bMap *Tx
	bMap, err = NewFromBob(bobTx)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}
	if bMap.BAP.Type != bap.ATTEST {
		t.Fatalf("expected: %s but got: %s", bap.ATTEST, bMap.BAP.Type)
	}
	if bMap.AIP.Signature != "H+lubfcz5Z2oG8B7HwmP8Z+tALP+KNOPgedo7UTXwW8LBpMkgCgatCdpvbtf7wZZQSIMz83emmAvVS4S3F5X1wo=" {
		t.Fatalf("expected: %s but got: %s", "H+lubfcz5Z2oG8B7HwmP8Z+tALP+KNOPgedo7UTXwW8LBpMkgCgatCdpvbtf7wZZQSIMz83emmAvVS4S3F5X1wo=", bMap.AIP.Signature)
	}
}
