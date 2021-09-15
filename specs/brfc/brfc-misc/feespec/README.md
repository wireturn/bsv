# Fee Specification

|     BRFC     	|  title  	|  authors 	| version 	|
|:------------:	|:-------:	|:--------:	|:-------:	|
| fb567267440a 	| feeSpec 	| nChain 	|   0.1   	|

## Overview
The following specification describes the structure and functionality of fees that miners accept and advertise. Miners can accept different fees for different types of transactions. For example, a miner might be willing to mine data transactions (which use OP_RETURN) at a much lower fee than regular transactions, which actually take up space in the UTOS (unspent transaction output set). Also, miners can advertise their `relayFee`, which is effectively the minimum threshold where a miner would be willing to relay and hold a transaction in their mempool even though it is lower than the fee they accept to miner transactions. This would help in terms of double spend protection according to the first-seen rule, as well as make it easier for other miners to mine blocks with lower fee transactions.

The bytes field is essentally a unit denoting the number of bytes that the satoshi amount will cover. In the example below the `relayFee` for the `data` feeType would be interpreted as 2 satoshis per 1000 bytes. Making the unit variable future-proofs this standard against a scenario where more granularity is required. It is assumed by convention that miners will use values that are powers of 10 for the bytes field for simple interpretation however this is not a requirement.


## Fee Specification Example


```json
{
  "feeType": "standard",
  "miningFee": {
      "satoshis": 1,
      "bytes": 1
  },
  "relayFee": {
    "satoshis": 1,
    "bytes": 10
  }
}
```

```json
{
  "feeType": "data",
  "miningFee": {
      "satoshis": 2,
      "bytes": 1000
  },
  "relayFee": {
      "satoshis": 1,
      "bytes": 10000
  }
}
``` 

-----


## Deterministic Transaction Fee Calculation (DTFC)
To calculate the fee needed for a transaction of size `TX_SIZE` (in bytes), with the following fee:
```json 
  "miningFee": {
      "satoshis": FEE_SATS,
      "bytes": FEE_BYTES
  }
```
>where `FEE_SATS` and `FEE_BYTES` are integers  

### DTFC steps
1. **MULTIPLY**:  
`temp1` = `FEE_SATS` * `TX_SIZE` 

2. **DIVIDE**:   
 `temp2` = `temp1` / `FEE_BYTES`

3. **FLOOR**:   
`ans` = floor(`temp2`)
> **note:** no need for floor if integer maths used since the decimal remainder is just ignored/dropped

4. **CHECK IF ZERO**:  
if (`ans` == 0) {
  if (`FEE_SATS` != 0) {
    `ans` = 1
  }
}

### DTFC Rationale
#### Determinism
Transaction fee calculation must be done **deterministically** in order to avoid any discrepancies that could cause transactions to be rejected by a miner because the fee is too low. As a result the **order** of the steps is very important.

#### Integer maths
In order to avoid issues that arise when using floating point representations of numbers, the simpler solution is to use integer maths where those issues are not applicable.

#### Check if zero
Since the minimum fee possible is just 1 satoshi, if the algorithm evaluates to 0, the fee paid should just be 1 satoshi -- unless `FEE_SATS` is explicitly set to 0, meaning the the miner is accepting 0-fee transactions.

#### Bytes = power of 10 (recommendation)
In order to not be too restrictive, the `FEE_BYTES` does not have to be a power of 10, however it is recommended as good standard practice.


-----

## Calculating Tx Fee with Different FeeTypes:
Consider the following transcation:
```json
{
  "txid": "0f5a5fddda11dfb173bbdb6272cead1c2f4822f853e35a77421b17b6320b7081",
  "hash": "0f5a5fddda11dfb173bbdb6272cead1c2f4822f853e35a77421b17b6320b7081",
  "version": 2,
  "size": 389,
  "locktime": 0,
  "vin": [
    {
      "txid": "4685daa9da2a7ac6173da8899d4cd9f1bc2b4de76085a092710c3aefd2927004",
      "vout": 2,
      "scriptSig": {
        "asm": "3045022100918b998c14cfc43125be26208af6b3af6d4029d853a9e2101915e7fc3b4c1fbc02206b6ff5f5d04de56de479321a45e5786051ca4084602f1fabbba7c55f3a44a568[ALL|FORKID] 02c6391cb6c2b339a27ca03f4a38785ae3b72a8667ae0d0a2b51f87625711574fc",
        "hex": "483045022100918b998c14cfc43125be26208af6b3af6d4029d853a9e2101915e7fc3b4c1fbc02206b6ff5f5d04de56de479321a45e5786051ca4084602f1fabbba7c55f3a44a568412102c6391cb6c2b339a27ca03f4a38785ae3b72a8667ae0d0a2b51f87625711574fc"
      },
      "sequence": 4294967295
    }
  ],
  "vout": [
    {
      "value": 0.00003276,
      "n": 0,
      "scriptPubKey": {
        "asm": "OP_DUP OP_HASH160 f12529a2eaf638dd43cd77a760a2b43cc3739721 OP_EQUALVERIFY OP_CHECKSIG",
        "hex": "76a914f12529a2eaf638dd43cd77a760a2b43cc373972188ac",
        "reqSigs": 1,
        "type": "pubkeyhash",
        "addresses": [
          "n3W1fPcu4FAEwuxd9Vg6RPqSuEq4N1QVsh"
        ]
      }
    },
    {
      "value": 0.00,
      "n": 1,
      "scriptPubKey": {
        "asm": "0 OP_RETURN 436974796f6e636861696e2c4e657720596f726b7c63697479416c62756d7c75706c6f61647c3139302a2a2a2a71712e636f6d7c68747470733a2f2f636974796f6e636861696e2e6f73732d636e2d686f6e676b6f6e672e616c6979756e63732e636f6d2f3135383032313036303432333461613436646664643333643133336334326631316130633534653030363431632e6a7067",
        "hex": "006a4c96436974796f6e636861696e2c4e657720596f726b7c63697479416c62756d7c75706c6f61647c3139302a2a2a2a71712e636f6d7c68747470733a2f2f636974796f6e636861696e2e6f73732d636e2d686f6e676b6f6e672e616c6979756e63732e636f6d2f3135383032313036303432333461613436646664643333643133336334326631316130633534653030363431632e6a7067",
        "type": "nulldata"
      }
    },
    {
      "value": 0.04022046,
      "n": 2,
      "scriptPubKey": {
        "asm": "OP_DUP OP_HASH160 0cd5ece936542595830dd9ba8e614104b05d8ae9 OP_EQUALVERIFY OP_CHECKSIG",
        "hex": "76a9140cd5ece936542595830dd9ba8e614104b05d8ae988ac",
        "reqSigs": 1,
        "type": "pubkeyhash",
        "addresses": [
          "mggpg9yzGdrbrJKD1fsQmhS5CfFAGrR5Vq"
        ]
      }
    }
  ]
}
```

Shown in hex below:
```
0200000001047092d2ef3a0c7192a08560e74d2bbcf1d94c9d89a83d17c67a2adaa9da8546020000006b483045022100918b998c14cfc43125be26208af6b3af6d4029d853a9e2101915e7fc3b4c1fbc02206b6ff5f5d04de56de479321a45e5786051ca4084602f1fabbba7c55f3a44a568412102c6391cb6c2b339a27ca03f4a38785ae3b72a8667ae0d0a2b51f87625711574fcffffffff03cc0c0000000000001976a914f12529a2eaf638dd43cd77a760a2b43cc373972188ac00000000000000009a006a4c96436974796f6e636861696e2c4e657720596f726b7c63697479416c62756d7c75706c6f61647c3139302a2a2a2a71712e636f6d7c68747470733a2f2f636974796f6e636861696e2e6f73732d636e2d686f6e676b6f6e672e616c6979756e63732e636f6d2f3135383032313036303432333461613436646664643333643133336334326631316130633534653030363431632e6a70671e5f3d00000000001976a9140cd5ece936542595830dd9ba8e614104b05d8ae988ac00000000
```

Bytes which are considered as `data` and therefore qualify for the `data` feeType rate are bytes of scripts that start with `6a` or `006a` (`OP_RETURN` or `OP_0 OP_RETURN`). This implies that only the following bytes are charged using the `data` fee rate:
```
006a4c96436974796f6e636861696e2c4e657720596f726b7c63697479416c62756d7c75706c6f61647c3139302a2a2a2a71712e636f6d7c68747470733a2f2f636974796f6e636861696e2e6f73732d636e2d686f6e676b6f6e672e616c6979756e63732e636f6d2f3135383032313036303432333461613436646664643333643133336334326631316130633534653030363431632e6a7067
```

while the rest of the transaction below (with the `data` omitted) is charged using the `standard` rate:
```
0200000001047092d2ef3a0c7192a08560e74d2bbcf1d94c9d89a83d17c67a2adaa9da8546020000006b483045022100918b998c14cfc43125be26208af6b3af6d4029d853a9e2101915e7fc3b4c1fbc02206b6ff5f5d04de56de479321a45e5786051ca4084602f1fabbba7c55f3a44a568412102c6391cb6c2b339a27ca03f4a38785ae3b72a8667ae0d0a2b51f87625711574fcffffffff03cc0c0000000000001976a914f12529a2eaf638dd43cd77a760a2b43cc373972188ac00000000000000009a

1e5f3d00000000001976a9140cd5ece936542595830dd9ba8e614104b05d8ae988ac00000000
```


-----


## Fiat Fee Pricing
Consider the case where a miner would like to set their transaction fees based on a fiat pricing system of `$0.000001` per byte:

Let `1₿` = `$102.58`

To get the `FEE_SATS` and `FEE_BYTES` needed:

`FEE_SATS` = 0.000001*10^8/102.58 = `0.9748488984`

Since `FEE_SATS` must be an integer, the miner must decide the amount of percision/accuracy needed through the use of `FEE_BYTES` and would charge:

```json 
  "miningFee": {
      "satoshis": 974,
      "bytes": 1000
  }
```
or
```json 
  "miningFee": {
      "satoshis": 974848,
      "bytes": 1000000
  }
```


-----

## Implementations in different languages

### C
```c
#include <stdio.h>
int main() {
  int satoshis = 974;
  int bytes = 1000;
  int tx_size = 500;
  int fee = tx_size * satoshis / bytes;
  printf("fee = %d\n", fee);
}
```
### Go
```go
package main
import "fmt"
func main() {
	var satoshis int = 974
	var bytes int = 1000
	var txsize int = 500
	fee := txsize * satoshis / bytes
	fmt.Printf("fee = %d\n", fee)
}
```
### Java
```java
class Fee {
  public static void main(String[] args) {
    int satoshis = 974;
    int bytes = 1000;
    int txsize = 500;
    int fee = txsize * satoshis / bytes;
    System.out.printf("fee = %d\n", fee);
  }
}
```
### JavaScript
```javascript
var satoshis = 974
var bytes = 1000
var txsize = 500
var fee = Math.floor(txsize * satoshis / bytes)
console.log('fee = ' + fee)
```
### Python
```python
satoshis = 974
bytes = 1000
txsize = 500
fee = txsize * satoshis // bytes
print("fee = %d" % fee)
```
