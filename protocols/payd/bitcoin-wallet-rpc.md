# rpc

== Wallet ==
abandontransaction "txid"
addmultisigaddress nrequired ["key",...] ( "account" )
backupwallet "destination"
dumpprivkey "address"
dumpwallet "filename"
encryptwallet "passphrase"
getaccount "address"
getaccountaddress "account"
getaddressesbyaccount "account"
getbalance ( "account" minconf include_watchonly )
getnewaddress ( "account" )
getrawchangeaddress
getreceivedbyaccount "account" ( minconf )
getreceivedbyaddress "address" ( minconf )
gettransaction "txid" ( include_watchonly )
getunconfirmedbalance
getwalletinfo
importaddress "address" ( "label" rescan p2sh )
importmulti "requests" "options"
importprivkey "bitcoinprivkey" ( "label" ) ( rescan )
importprunedfunds
importpubkey "pubkey" ( "label" rescan )
importwallet "filename"
keypoolrefill ( newsize )
listaccounts ( minconf include_watchonly)
listaddressgroupings
listlockunspent
listreceivedbyaccount ( minconf include_empty include_watchonly)
listreceivedbyaddress ( minconf include_empty include_watchonly)
listsinceblock ( "blockhash" target_confirmations include_watchonly)
listtransactions ( "account" count skip include_watchonly)
listunspent ( minconf maxconf  ["addresses",...] [include_unsafe] )
listwallets
lockunspent unlock ([{"txid":"txid","vout":n},...])
move "fromaccount" "toaccount" amount ( minconf "comment" )
removeprunedfunds "txid"
sendfrom "fromaccount" "toaddress" amount ( minconf "comment" "comment_to" )
sendmany "fromaccount" {"address":amount,...} ( minconf "comment" ["address",...] )
sendtoaddress "address" amount ( "comment" "comment_to" subtractfeefromamount )
setaccount "address" "account"
settxfee amount
signmessage "address" "message"

## details

== Wallet ==
- abandontransaction "txid"
```
Mark in-wallet transaction <txid> as abandoned
This will mark this transaction and all its in-wallet descendants as abandoned which will allow
for their inputs to be respent.  It can be used to replace "stuck" or evicted transactions.
It only works on transactions which are not included in a block and are not currently in the mempool.
It has no effect on transactions which are already conflicted or abandoned.

Arguments:
1. "txid"    (string, required) The transaction id

Result:

Examples:
> bitcoin-cli abandontransaction "1075db55d416d3ca199f55b6084e2115b9345e16c5cf302fc80e9d5fbf5d48d"
> curl --user myusername --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "abandontransaction", "params": ["1075db55d416d3ca199f55b6084e2115b9345e16c5cf302fc80e9d5fbf5d48d"] }' -H 'content-type: text/plain;' http://127.0.0.1:8332/
```


- addmultisigaddress nrequired ["key",...] ( "account" )
```
Add a nrequired-to-sign multisignature address to the wallet.
Each key is a Bitcoin address or hex-encoded public key.
If 'account' is specified (DEPRECATED), assign address to that account.

Arguments:
1. nrequired        (numeric, required) The number of required signatures out of the n keys or addresses.
2. "keys"         (string, required) A json array of bitcoin addresses or hex-encoded public keys
     [
       "address"  (string) bitcoin address or hex-encoded public key
       ...,
     ]
3. "account"      (string, optional) DEPRECATED. An account to assign the addresses to.

Result:
"address"         (string) A bitcoin address associated with the keys.

Examples:

Add a multisig address from 2 addresses
> bitcoin-cli addmultisigaddress 2 "[\"16sSauSf5pF2UkUwvKGq4qjNRzBZYqgEL5\",\"171sgjn4YtPu27adkKGrdDwzRTxnRkBfKV\"]"

As json rpc call
> curl --user myusername --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "addmultisigaddress", "params": [2, "[\"16sSauSf5pF2UkUwvKGq4qjNRzBZYqgEL5\",\"171sgjn4YtPu27adkKGrdDwzRTxnRkBfKV\"]"] }' -H 'content-type: text/plain;' http://127.0.0.1:8332/
```

- backupwallet "destination"
```
Safely copies current wallet file to destination, which can be a directory or a path with filename.

Arguments:
1. "destination"   (string) The destination directory or file

Examples:
> bitcoin-cli backupwallet "backup.dat"
> curl --user myusername --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "backupwallet", "params": ["backup.dat"] }' -H 'content-type: text/plain;' http://127.0.0.1:8332/
```

- dumpprivkey "address"
```
Reveals the private key corresponding to 'address'.
Then the importprivkey can be used with this output

Arguments:
1. "address"   (string, required) The bitcoin address for the private key

Result:
"key"                (string) The private key

Examples:
> bitcoin-cli dumpprivkey "myaddress"
> bitcoin-cli importprivkey "mykey"
> curl --user myusername --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "dumpprivkey", "params": ["myaddress"] }' -H 'content-type: text/plain;' http://127.0.0.1:8332/
```

- dumpwallet "filename"
```
Dumps all wallet keys in a human-readable format to a server-side file. This does not allow overwriting existing files.

Arguments:
1. "filename"    (string, required) The filename with path (either absolute or relative to bitcoind)

Result:
{                           (json object)
  "filename" : {        (string) The filename with full absolute path
}

Examples:
> bitcoin-cli dumpwallet "test"
> curl --user myusername --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "dumpwallet", "params": ["test"] }' -H 'content-type: text/plain;' http://127.0.0.1:8332/
```

- encryptwallet "passphrase"
```
Encrypts the wallet with 'passphrase'. This is for first time encryption.
After this, any calls that interact with private keys such as sending or signing
will require the passphrase to be set prior the making these calls.
Use the walletpassphrase call for this, and then walletlock call.
If the wallet is already encrypted, use the walletpassphrasechange call.
Note that this will shutdown the server.

Arguments:
1. "passphrase"    (string) The pass phrase to encrypt the wallet with. It must be at least 1 character, but should be long.

Examples:

Encrypt you wallet
> bitcoin-cli encryptwallet "my pass phrase"

Now set the passphrase to use the wallet, such as for signing or sending bitcoin
> bitcoin-cli walletpassphrase "my pass phrase"

Now we can so something like sign
> bitcoin-cli signmessage "address" "test message"

Now lock the wallet again by removing the passphrase
> bitcoin-cli walletlock

As a json rpc call
> curl --user myusername --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "encryptwallet", "params": ["my pass phrase"] }' -H 'content-type: text/plain;' http://127.0.0.1:8332/
```

- getaccount "address"
```
DEPRECATED.
```

- getaccountaddress "account"
```
DEPRECATED.
```

- getaddressesbyaccount "account"
```
DEPRECATED.
```

- getbalance ( "account" minconf include_watchonly )
```
If account is not specified, returns the server's total available balance.
If account is specified (DEPRECATED), returns the balance in the account.
Note that the account "" is not the same as leaving the parameter out.
The server total may be different to the balance in the default "" account.

Arguments:
1. "account"         (string, optional) DEPRECATED. The account string may be given as a
                     specific account name to find the balance associated with wallet keys in
                     a named account, or as the empty string ("") to find the balance
                     associated with wallet keys not in any named account, or as "*" to find
                     the balance associated with all wallet keys regardless of account.
                     When this option is specified, it calculates the balance in a different
                     way than when it is not specified, and which can count spends twice when
                     there are conflicting pending transactions temporarily resulting in low
                     or even negative balances.
                     In general, account balance calculation is not considered reliable and
                     has resulted in confusing outcomes, so it is recommended to avoid passing
                     this argument.
2. minconf           (numeric, optional, default=1) Only include transactions confirmed at least this many times.
3. include_watchonly (bool, optional, default=false) Also include balance in watch-only addresses (see 'importaddress')

Result:
amount              (numeric) The total amount in BSV received for this account.

Examples:

The total amount in the wallet
> bitcoin-cli getbalance

The total amount in the wallet at least 5 blocks confirmed
> bitcoin-cli getbalance "*" 6

As a json rpc call
> curl --user myusername --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "getbalance", "params": ["*", 6] }' -H 'content-type: text/plain;' http://127.0.0.1:8332/
```

- getnewaddress ( "account" )
```
Returns a new Bitcoin address for receiving payments.
If 'account' is specified (DEPRECATED), it is added to the address book
so payments received with the address will be credited to 'account'.

Arguments:
1. "account"        (string, optional) DEPRECATED. The account name for the address to be linked to. If not provided, the default account "" is used. It can also be set to the empty string "" to represent the default account. The account does not need to exist, it will be created if there is no account by the given name.

Result:
"address"    (string) The new bitcoin address

Examples:
> bitcoin-cli getnewaddress
> curl --user myusername --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "getnewaddress", "params": [] }' -H 'content-type: text/plain;' http://127.0.0.1:8332/
```

- getrawchangeaddress
```
Returns a new Bitcoin address, for receiving change.
This is for use with raw transactions, NOT normal use.

Result:
"address"    (string) The address

Examples:
> bitcoin-cli getrawchangeaddress
> curl --user myusername --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "getrawchangeaddress", "params": [] }' -H 'content-type: text/plain;' http://127.0.0.1:8332/
```

- getreceivedbyaccount "account" ( minconf )
```
DEPRECATED.
```

- getreceivedbyaddress "address" ( minconf )
```
Returns the total amount received by the given address in transactions with at least minconf confirmations.

Arguments:
1. "address"         (string, required) The bitcoin address for transactions.
2. minconf             (numeric, optional, default=1) Only include transactions confirmed at least this many times.

Result:
amount   (numeric) The total amount in BSV received at this address.

Examples:

The amount from transactions with at least 1 confirmation
> bitcoin-cli getreceivedbyaddress "1D1ZrZNe3JUo7ZycKEYQQiQAWd9y54F4XX"

The amount including unconfirmed transactions, zero confirmations
> bitcoin-cli getreceivedbyaddress "1D1ZrZNe3JUo7ZycKEYQQiQAWd9y54F4XX" 0

The amount with at least 6 confirmation, very safe
> bitcoin-cli getreceivedbyaddress "1D1ZrZNe3JUo7ZycKEYQQiQAWd9y54F4XX" 6

As a json rpc call
> curl --user myusername --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "getreceivedbyaddress", "params": ["1D1ZrZNe3JUo7ZycKEYQQiQAWd9y54F4XX", 6] }' -H 'content-type: text/plain;' http://127.0.0.1:8332/
```

- gettransaction "txid" ( include_watchonly )
```
Get detailed information about in-wallet transaction <txid>

Arguments:
1. "txid"                  (string, required) The transaction id
2. "include_watchonly"     (bool, optional, default=false) Whether to include watch-only addresses in balance calculation and details[]

Result:
{
  "amount" : x.xxx,        (numeric) The transaction amount in BSV
  "fee": x.xxx,            (numeric) The amount of the fee in BSV. This is negative and only available for the
                              'send' category of transactions.
  "confirmations" : n,     (numeric) The number of confirmations
  "blockhash" : "hash",  (string) The block hash
  "blockindex" : xx,       (numeric) The index of the transaction in the block that includes it
  "blocktime" : ttt,       (numeric) The time in seconds since epoch (1 Jan 1970 GMT)
  "txid" : "transactionid",   (string) The transaction id.
  "time" : ttt,            (numeric) The transaction time in seconds since epoch (1 Jan 1970 GMT)
  "timereceived" : ttt,    (numeric) The time received in seconds since epoch (1 Jan 1970 GMT)
  "details" : [
    {
      "account" : "accountname",      (string) DEPRECATED. The account name involved in the transaction, can be "" for the default account.
      "address" : "address",          (string) The bitcoin address involved in the transaction
      "category" : "send|receive",    (string) The category, either 'send' or 'receive'
      "amount" : x.xxx,                 (numeric) The amount in BSV
      "label" : "label",              (string) A comment for the address/transaction, if any
      "vout" : n,                       (numeric) the vout value
      "fee": x.xxx,                     (numeric) The amount of the fee in BSV. This is negative and only available for the
                                           'send' category of transactions.
      "abandoned": xxx                  (bool) 'true' if the transaction has been abandoned (inputs are respendable). Only available for the
                                           'send' category of transactions.
    }
    ,...
  ],
  "hex" : "data"         (string) Raw data for transaction
}

Examples:
> bitcoin-cli gettransaction "1075db55d416d3ca199f55b6084e2115b9345e16c5cf302fc80e9d5fbf5d48d"
> bitcoin-cli gettransaction "1075db55d416d3ca199f55b6084e2115b9345e16c5cf302fc80e9d5fbf5d48d" true
> curl --user myusername --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "gettransaction", "params": ["1075db55d416d3ca199f55b6084e2115b9345e16c5cf302fc80e9d5fbf5d48d"] }' -H 'content-type: text/plain;' http://127.0.0.1:8332/
```

- getunconfirmedbalance
```
Returns the server's total unconfirmed balance
```

- getwalletinfo
```
Returns an object containing various wallet state info.

Result:
{
  "walletname": xxxxx,             (string) the wallet name
  "walletversion": xxxxx,          (numeric) the wallet version
  "balance": xxxxxxx,              (numeric) the total confirmed balance of the wallet in BSV
  "unconfirmed_balance": xxx,      (numeric) the total unconfirmed balance of the wallet in BSV
  "immature_balance": xxxxxx,      (numeric) the total immature balance of the wallet in BSV
  "txcount": xxxxxxx,              (numeric) the total number of transactions in the wallet
  "keypoololdest": xxxxxx,         (numeric) the timestamp (seconds since Unix epoch) of the oldest pre-generated key in the key pool
  "keypoolsize": xxxx,             (numeric) how many new keys are pre-generated (only counts external keys)
  "keypoolsize_hd_internal": xxxx, (numeric) how many new keys are pre-generated for internal use (used for change outputs, only appears if the wallet is using this feature, otherwise external keys are used)
  "unlocked_until": ttt,           (numeric) the timestamp in seconds since epoch (midnight Jan 1 1970 GMT) that the wallet is unlocked for transfers, or 0 if the wallet is locked
  "paytxfee": x.xxxx,              (numeric) the transaction fee configuration, set in BSV/kB
  "hdmasterkeyid": "<hash160>"     (string) the Hash160 of the HD master pubkey
}

Examples:
> bitcoin-cli getwalletinfo
> curl --user myusername --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "getwalletinfo", "params": [] }' -H 'content-type: text/plain;' http://127.0.0.1:8332/
```

- importaddress "address" ( "label" rescan p2sh )
```
Adds a script (in hex) or address that can be watched as if it were in your wallet but cannot be used to spend.

Arguments:
1. "script"           (string, required) The hex-encoded script (or address)
2. "label"            (string, optional, default="") An optional label
3. rescan               (boolean, optional, default=true) Rescan the wallet for transactions
4. p2sh                 (boolean, optional, default=false) Add the P2SH version of the script as well

Note: This call can take minutes to complete if rescan is true.
If you have the full public key, you should call importpubkey instead of this.

Note: If you import a non-standard raw script in hex form, outputs sending to it will be treated
as change, and not show up in many RPCs.

Examples:

Import a script with rescan
> bitcoin-cli importaddress "myscript"

Import using a label without rescan
> bitcoin-cli importaddress "myscript" "testing" false

As a JSON-RPC call
> curl --user myusername --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "importaddress", "params": ["myscript", "testing", false] }' -H 'content-type: text/plain;' http://127.0.0.1:8332/
```

- importmulti "requests" "options"
```
Import addresses/scripts (with private or public keys, redeem script (P2SH)), rescanning all addresses in one-shot-only (rescan can be disabled via options).

Arguments:
1. requests     (array, required) Data to be imported
  [     (array of json objects)
    {
      "scriptPubKey": "<script>" | { "address":"<address>" }, (string / json, required) Type of scriptPubKey (string for script, json for address)
      "timestamp": timestamp | "now"                        , (integer / string, required) Creation time of the key in seconds since epoch (Jan 1 1970 GMT),
                                                              or the string "now" to substitute the current synced blockchain time. The timestamp of the oldest
                                                              key will determine how far back blockchain rescans need to begin for missing wallet transactions.
                                                              "now" can be specified to bypass scanning, for keys which are known to never have been used, and
                                                              0 can be specified to scan the entire blockchain. Blocks up to 2 hours before the earliest key
                                                              creation time of all keys being imported by the importmulti call will be scanned.
      "redeemscript": "<script>"                            , (string, optional) Allowed only if the scriptPubKey is a P2SH address or a P2SH scriptPubKey
      "pubkeys": ["<pubKey>", ... ]                         , (array, optional) Array of strings giving pubkeys that must occur in the output or redeemscript
      "keys": ["<key>", ... ]                               , (array, optional) Array of strings giving private keys whose corresponding public keys must occur in the output or redeemscript
      "internal": <true>                                    , (boolean, optional, default: false) Stating whether matching outputs should be be treated as not incoming payments
      "watchonly": <true>                                   , (boolean, optional, default: false) Stating whether matching outputs should be considered watched even when they're not spendable, only allowed if keys are empty
      "label": <label>                                      , (string, optional, default: '') Label to assign to the address (aka account name, for now), only allowed with internal=false
    }
  ,...
  ]
2. options                 (json, optional)
  {
     "rescan": <false>,         (boolean, optional, default: true) Stating if should rescan the blockchain after all imports
  }

Examples:
> bitcoin-cli importmulti '[{ "scriptPubKey": { "address": "<my address>" }, "timestamp":1455191478 }, { "scriptPubKey": { "address": "<my 2nd address>" }, "label": "example 2", "timestamp": 1455191480 }]'
> bitcoin-cli importmulti '[{ "scriptPubKey": { "address": "<my address>" }, "timestamp":1455191478 }]' '{ "rescan": false}'

Response is an array with the same size as the input that has the execution result :
  [{ "success": true } , { "success": false, "error": { "code": -1, "message": "Internal Server Error"} }, ... ]
```

- importprivkey "bitcoinprivkey" ( "label" ) ( rescan )
```
Adds a private key (as returned by dumpprivkey) to your wallet.

Arguments:
1. "bitcoinprivkey"   (string, required) The private key (see dumpprivkey)
2. "label"            (string, optional, default="") An optional label
3. rescan               (boolean, optional, default=true) Rescan the wallet for transactions

Note: This call can take minutes to complete if rescan is true.

Examples:

Dump a private key
> bitcoin-cli dumpprivkey "myaddress"

Import the private key with rescan
> bitcoin-cli importprivkey "mykey"

Import using a label and without rescan
> bitcoin-cli importprivkey "mykey" "testing" false

Import using default blank label and without rescan
> bitcoin-cli importprivkey "mykey" "" false

As a JSON-RPC call
> curl --user myusername --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "importprivkey", "params": ["mykey", "testing", false] }' -H 'content-type: text/plain;' http://127.0.0.1:8332/
```

- importprunedfunds
```
Imports funds without rescan. Corresponding address or script must previously be included in wallet. Aimed towards pruned wallets. The end-user is responsible to import additional transactions that subsequently spend the imported outputs or rescan after the point in the blockchain the transaction is included.

Arguments:
1. "rawtransaction" (string, required) A raw transaction in hex funding an already-existing address in wallet
2. "txoutproof"     (string, required) The hex output from gettxoutproof that contains the transaction
```

- importpubkey "pubkey" ( "label" rescan )
```
Adds a public key (in hex) that can be watched as if it were in your wallet but cannot be used to spend.

Arguments:
1. "pubkey"           (string, required) The hex-encoded public key
2. "label"            (string, optional, default="") An optional label
3. rescan               (boolean, optional, default=true) Rescan the wallet for transactions

Note: This call can take minutes to complete if rescan is true.

Examples:

Import a public key with rescan
> bitcoin-cli importpubkey "mypubkey"

Import using a label without rescan
> bitcoin-cli importpubkey "mypubkey" "testing" false

As a JSON-RPC call
> curl --user myusername --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "importpubkey", "params": ["mypubkey", "testing", false] }' -H 'content-type: text/plain;' http://127.0.0.1:8332/
```

- importwallet "filename"
```
Imports keys from a wallet dump file (see dumpwallet).

Arguments:
1. "filename"    (string, required) The wallet file

Examples:

Dump the wallet
> bitcoin-cli dumpwallet "test"

Import the wallet
> bitcoin-cli importwallet "test"

Import using the json rpc call
> curl --user myusername --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "importwallet", "params": ["test"] }' -H 'content-type: text/plain;' http://127.0.0.1:8332/
```

- keypoolrefill ( newsize )
```
Fills the keypool.

Arguments
1. newsize     (numeric, optional, default=100) The new keypool size

Examples:
> bitcoin-cli keypoolrefill
> curl --user myusername --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "keypoolrefill", "params": [] }' -H 'content-type: text/plain;' http://127.0.0.1:8332/
```

- listaccounts ( minconf include_watchonly)
```
DEPRECATED.
```

- listaddressgroupings
```
Lists groups of addresses which have had their common ownership
made public by common use as inputs or as the resulting change
in past transactions

Result:
[
  [
    [
      "address",            (string) The bitcoin address
      amount,                 (numeric) The amount in BSV
      "account"             (string, optional) DEPRECATED. The account
    ]
    ,...
  ]
  ,...
]

Examples:
> bitcoin-cli listaddressgroupings
> curl --user myusername --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "listaddressgroupings", "params": [] }' -H 'content-type: text/plain;' http://127.0.0.1:8332/
```

- listlockunspent
```
Returns list of temporarily unspendable outputs.
See the lockunspent call to lock and unlock transactions for spending.

Result:
[
  {
    "txid" : "transactionid",     (string) The transaction id locked
    "vout" : n                      (numeric) The vout value
  }
  ,...
]

Examples:

List the unspent transactions
> bitcoin-cli listunspent

Lock an unspent transaction
> bitcoin-cli lockunspent false "[{\"txid\":\"a08e6907dbbd3d809776dbfc5d82e371b764ed838b5655e72f463568df1aadf0\",\"vout\":1}]"

List the locked transactions
> bitcoin-cli listlockunspent

Unlock the transaction again
> bitcoin-cli lockunspent true "[{\"txid\":\"a08e6907dbbd3d809776dbfc5d82e371b764ed838b5655e72f463568df1aadf0\",\"vout\":1}]"

As a json rpc call
> curl --user myusername --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "listlockunspent", "params": [] }' -H 'content-type: text/plain;' http://127.0.0.1:8332/
```

- listreceivedbyaccount ( minconf include_empty include_watchonly)
```
DEPRECATED.
```

- listreceivedbyaddress ( minconf include_empty include_watchonly)
```
List balances by receiving address.

Arguments:
1. minconf           (numeric, optional, default=1) The minimum number of confirmations before payments are included.
2. include_empty     (bool, optional, default=false) Whether to include addresses that haven't received any payments.
3. include_watchonly (bool, optional, default=false) Whether to include watch-only addresses (see 'importaddress').

Result:
[
  {
    "involvesWatchonly" : true,        (bool) Only returned if imported addresses were involved in transaction
    "address" : "receivingaddress",  (string) The receiving address
    "account" : "accountname",       (string) DEPRECATED. The account of the receiving address. The default account is "".
    "amount" : x.xxx,                  (numeric) The total amount in BSV received by the address
    "confirmations" : n,               (numeric) The number of confirmations of the most recent transaction included
    "label" : "label",               (string) A comment for the address/transaction, if any
    "txids": [
       n,                                (numeric) The ids of transactions received with the address
       ...
    ]
  }
  ,...
]

Examples:
> bitcoin-cli listreceivedbyaddress
> bitcoin-cli listreceivedbyaddress 6 true
> curl --user myusername --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "listreceivedbyaddress", "params": [6, true, true] }' -H 'content-type: text/plain;' http://127.0.0.1:8332/
```

- listsinceblock ( "blockhash" target_confirmations include_watchonly)
```
Get all transactions in blocks since block [blockhash], or all transactions if omitted

Arguments:
1. "blockhash"            (string, optional) The block hash to list transactions since
2. target_confirmations:    (numeric, optional) The confirmations required, must be 1 or more
3. include_watchonly:       (bool, optional, default=false) Include transactions to watch-only addresses (see 'importaddress')
Result:
{
  "transactions": [
    "account":"accountname",       (string) DEPRECATED. The account name associated with the transaction. Will be "" for the default account.
    "address":"address",    (string) The bitcoin address of the transaction. Not present for move transactions (category = move).
    "category":"send|receive",     (string) The transaction category. 'send' has negative amounts, 'receive' has positive amounts.
    "amount": x.xxx,          (numeric) The amount in BSV. This is negative for the 'send' category, and for the 'move' category for moves
                                          outbound. It is positive for the 'receive' category, and for the 'move' category for inbound funds.
    "vout" : n,               (numeric) the vout value
    "fee": x.xxx,             (numeric) The amount of the fee in BSV. This is negative and only available for the 'send' category of transactions.
    "confirmations": n,       (numeric) The number of confirmations for the transaction. Available for 'send' and 'receive' category of transactions.
                                          When it's < 0, it means the transaction conflicted that many blocks ago.
    "blockhash": "hashvalue",     (string) The block hash containing the transaction. Available for 'send' and 'receive' category of transactions.
    "blockindex": n,          (numeric) The index of the transaction in the block that includes it. Available for 'send' and 'receive' category of transactions.
    "blocktime": xxx,         (numeric) The block time in seconds since epoch (1 Jan 1970 GMT).
    "txid": "transactionid",  (string) The transaction id. Available for 'send' and 'receive' category of transactions.
    "time": xxx,              (numeric) The transaction time in seconds since epoch (Jan 1 1970 GMT).
    "timereceived": xxx,      (numeric) The time received in seconds since epoch (Jan 1 1970 GMT). Available for 'send' and 'receive' category of transactions.
    "abandoned": xxx,         (bool) 'true' if the transaction has been abandoned (inputs are respendable). Only available for the 'send' category of transactions.
    "comment": "...",       (string) If a comment is associated with the transaction.
    "label" : "label"       (string) A comment for the address/transaction, if any
    "to": "...",            (string) If a comment to is associated with the transaction.
  ],
  "lastblock": "lastblockhash"     (string) The hash of the last block
}

Examples:
> bitcoin-cli listsinceblock
> bitcoin-cli listsinceblock "000000000000000bacf66f7497b7dc45ef753ee9a7d38571037cdb1a57f663ad" 6
> curl --user myusername --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "listsinceblock", "params": ["000000000000000bacf66f7497b7dc45ef753ee9a7d38571037cdb1a57f663ad", 6] }' -H 'content-type: text/plain;' http://127.0.0.1:8332/
```

- listtransactions ( "account" count skip include_watchonly)
```
```

- listunspent ( minconf maxconf  ["addresses",...] [include_unsafe] )
```
```

- listwallets
```
```

- lockunspent unlock ([{"txid":"txid","vout":n},...])
```
```

- move "fromaccount" "toaccount" amount ( minconf "comment" )
```
```

- removeprunedfunds "txid"
```
```

- sendfrom "fromaccount" "toaddress" amount ( minconf "comment" "comment_to" )
```
```

- sendmany "fromaccount" {"address":amount,...} ( minconf "comment" ["address",...] )
```
```

- sendtoaddress "address" amount ( "comment" "comment_to" subtractfeefromamount )
```
```

- setaccount "address" "account"
```
```

- settxfee amount
```
```

- signmessage "address" "message"
```
```
