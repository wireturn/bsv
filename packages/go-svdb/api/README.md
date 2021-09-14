##### 获取节点总体数据：

```
地址：/api/blockchaininfo/[:type]
请求方式：GET
参数：type,币种[BTC,BCH,BSV]
返回数据示例：/api/blockchaininfo/BSV
返回：
{
  "chain": "main",
  "blocks": 570463,
  "headers": 570523,
  "bestblockhash": "0000000000000000061af27641eeafb5e217730baefd4368c632087228ea5a8b",
  "difficulty": 117661313779.2578,
  "mediantime": 1550648426,
  "verificationprogress": 0.9999998628601205,
  "chainwork": "000000000000000000000000000000000000000000dc202cd879be6d922bbef7",
  "pruned": false,
  "softforks": [
  {
      "id": "bip34",
      "version": 2,
      "reject": {
        "status": true
      }
    },
    {
      "id": "bip66",
      "version": 3,
      "reject": {
        "status": true
      }
    },
    {
      "id": "bip65",
      "version": 4,
      "reject": {
        "status": true
      }
    }
  ],
  "bip9_softforks": {
    "csv": {
      "status": "active",
      "startTime": 1462060800,
      "timeout": 1493596800,
      "since": 419328
    }
  }
```

##### 根据高度，获取某个区块信息

```
地址：/api/blockByHeight/[:type]/[:height]
请求方式：GET
参数：type,币种[BTC,BCH,BSV]
参数：区块的高度值
返回数据示例：/api/blockByHeight/BSV/100 返回：
{ Hash: '000000007bc154e0fa7ea32218a72fe2c1bb9f86cf8c9ebf9a715ed27fdb229a',
  Confirmations: 570526,
  StrippedSize: 0,
  Size: 215,
  Weight: 0,
  Height: 100,
  Version: 1,
  VersionHex: '00000001',
  MerkleRoot: '2d05f0c9c3e1c226e63b5fac240137687544cf631cd616fd34fd188fc9020866',
  Tx:
   [ '2d05f0c9c3e1c226e63b5fac240137687544cf631cd616fd34fd188fc9020866' ],
  RawTx: null,
  Time: 1231660825,
  Nonce: 1573057331,
  Bits: '1d00ffff',
  Difficulty: 1,
  PreviousHash: '00000000cd9b12643e6854cb25939b39cd7a1ad0af31a9bd8b2efe67854b1995',
  NextHash: '00000000b69bd8e4dc60580117617a466d5c76ada85fb7b87e9baea01f9d9984',
  TxCnt: 1,
  CoinBaseInfo: 'BP//AB0BTQ==' }

```



##### 根据区块hash值，获取区块信息

```
地址：/api/blockByHash/[:type]/[:hash]
请求方式：GET
参数：type,币种[BTC,BCH,BSV]
参数:区块HASH值
返回数据示例：/api/blockByHash/BSV/0000000071966c2b1d065fd446b1e485b2c9d9594acd2007ccbd5441cfc89444
返回：
{ Hash: '000000007bc154e0fa7ea32218a72fe2c1bb9f86cf8c9ebf9a715ed27fdb229a',
  Confirmations: 570526,
  StrippedSize: 0,
  Size: 215,
  Weight: 0,
  Height: 100,
  Version: 1,
  VersionHex: '00000001',
  MerkleRoot: '2d05f0c9c3e1c226e63b5fac240137687544cf631cd616fd34fd188fc9020866',
  Tx:
   [ '2d05f0c9c3e1c226e63b5fac240137687544cf631cd616fd34fd188fc9020866' ],
  RawTx: null,
  Time: 1231660825,
  Nonce: 1573057331,
  Bits: '1d00ffff',
  Difficulty: 1,
  PreviousHash: '00000000cd9b12643e6854cb25939b39cd7a1ad0af31a9bd8b2efe67854b1995',
  NextHash: '00000000b69bd8e4dc60580117617a466d5c76ada85fb7b87e9baea01f9d9984',
  TxCnt: 1,
  CoinBaseInfo: 'BP//AB0BTQ=='
}
```


##### 获取最近的区块：
```
接口地址：/api/bestblock/[:type]/[:start]/[:limit]
请求方式：GET
参数：type,币种[BTC,BCH,BSV]
参数：start,区块的开始偏移
参数：limit,获取区块数
返回数据示例：/api/bestblock/BSV/0/15
返回:
[ { Hash: '000000000000000002208c6956fade5325e919164f961a986d65552b7f176011',
    Confirmations: 66,
    StrippedSize: 0,
    Size: 44739,
    Weight: 0,
    Height: 570624,
    Version: 536870912,
    VersionHex: '20000000',
    MerkleRoot: 'c037d1f9aee774ed3dbfdaadbd67bca2b654f226f7841748d24c5c73eadb2872',
    Tx:null,
    RawTx: null,
    Time: 1550755620,
    Nonce: 1726198927,
    Bits: '1809e76d',
    Difficulty: 111015153283.0872,
    PreviousHash: '000000000000000003555525c24def53c01ebc2f1ce88544c7f4573c7f6d4597',
    NextHash: '00000000000000000751a35e23038ccaa010c5d84d52310f4401f3fd793bd6f1',
    TxCnt: 79,
    CoinBaseInfo: 'AwC1CC9jb2luZ2Vlay5jb20vaHC/OOvwI3vIFn9kVZU=' } ]

```

##### 根据区块高度获取预览交易数据

```
接口地址：/api/block/browser_txs/[:type]/[:hash]/[:start]/[:end]
请求方式：GET
参数：type,币种[BTC,BCH,BSV]
参数：hash,区块高度
参数:start,开始交易位置
参数:end,结束交易位置
返回数据示例：/api/block/browser_txs/BSV/000000000000000002208c6956fade5325e919164f961a986d65552b7f176011/0/20
返回:
[ { Hash: 'fe28050b93faea61fa88c4c630f0e1f0a1c24d0082dd0e10d369e13212128f33',
    Height: 1000,
    VinAddrValue: [ '@0' ],
    VoutAddrValue: [ '1BW18n7MfpU35q4MTBSk8pse3XzQF8XvzT@5000000000' ],
    Fee: -5000000000,
    CoinDay: 0
  }
]
```

##### 根据交易hash获取rpc格式交易数据(快)：

```
地址：/api/tx/[:type]/[:txid]
请求方式：GET
参数：交易hash值
返回数据示例：/api/tx/BSV/07612439f1f75775377829d8394c7bb4264b3fe51712ddc65f81f28bdf4a392c
返回数据：
{ hex: '',
  txid: '07612439f1f75775377829d8394c7bb4264b3fe51712ddc65f81f28bdf4a392c',
  hash: '07612439f1f75775377829d8394c7bb4264b3fe51712ddc65f81f28bdf4a392c',
  size: 226,
  version: 2,
  locktime: 0,
  vin:
   [ { txid: '111a1f6a0e1775093bc7402e14218d76f185bd7369e5f72675eeb6340585fbd1',
       vout: 0,
       scriptSig: [Object],
       sequence: 4294967295 } ],
  vout:
   [ { value: 439.20127681, n: 0, scriptPubKey: [Object] },
     { value: 211.3124106, n: 1, scriptPubKey: [Object] } ] }
[ { txid: '111a1f6a0e1775093bc7402e14218d76f185bd7369e5f72675eeb6340585fbd1',
    vout: 0,
    scriptSig:
     { asm: '3045022100a96d252e46db05e3c73a8ccbd9084f70941f65f1aee5b31baa217fe7d08377ca02203b7965e202f002b1156a068200a4ba05a1367dba528afe75ee31048adc867dca[ALL|FORKID] 028f1e2506619857d053d2481e0d8f2c8231045acfe4eb745ff1f2bef6badada5e',
       hex: '483045022100a96d252e46db05e3c73a8ccbd9084f70941f65f1aee5b31baa217fe7d08377ca02203b7965e202f002b1156a068200a4ba05a1367dba528afe75ee31048adc867dca4121028f1e2506619857d053d2481e0d8f2c8231045acfe4eb745ff1f2bef6badada5e' },
    sequence: 4294967295 } ]
```

##### 根据交易hash获取浏览器交易数据(慢)：

```
地址：/api/browser_tx/[:type]/[:txid]
请求方式：GET
参数：交易hash值
返回数据示例：/api/browser_tx/BSV/07612439f1f75775377829d8394c7bb4264b3fe51712ddc65f81f28bdf4a392c
返回数据：
{
"txid":"987ea317eb5a7946067203f423fdd6d41e7c855d20f83bedbce9dc7211708197",
"version":1,
"locktime":0,
"vin":[
    {
        "txid":"9c206b9ab77aa304de16d6544091ffab8ae71ddb76492df3a83d6c4cf3d98c69",
        "vout":1,
        "sequence":4294967295,
        "n":0,
        "scriptSig":{
            "hex":"483045022100c500255080bab8bf5e82420890450d94b03443a70e267ba08953b67432726c8002201517d0dc7d07003ce6c29e4c7e16710cde225e1aaad24c588b87ea7ee6cb534e412102f354fcf38c5110c19bb007ce02d1c74db51476be77722fd8c616ad5ff57b2783",
            "asm":"3045022100c500255080bab8bf5e82420890450d94b03443a70e267ba08953b67432726c8002201517d0dc7d07003ce6c29e4c7e16710cde225e1aaad24c588b87ea7ee6cb534e[ALL|FORKID] 02f354fcf38c5110c19bb007ce02d1c74db51476be77722fd8c616ad5ff57b2783"
        },
        "addr":"161JK2mQ1ToGeaxbr83rtDmyAoFFj6F6Lr",
        "valueSat":1260799332,
        "value":12.60799332,
        "doubleSpentTxID":null
    }
],
"vout":[
    {
        "value":"7.31348411",
        "n":0,
        "scriptPubKey":{
            "hex":"76a9144649864610d227935996b16468652e179fce4c2988ac",
            "asm":"OP_DUP OP_HASH160 4649864610d227935996b16468652e179fce4c29 OP_EQUALVERIFY OP_CHECKSIG",
            "addresses":[
                "17QeP12NhT61R4sNhYTRJAFDmCjjZaM8y2"
            ]
        },
        "spentTxId":"615e5e0df447975c1ffacef4376d018c5af4ffba01fef860c1b6b0529f4e00cf",
        "spentIndex":0,
        "spentHeight":-1
    },
    {
        "value":"5.29450695",
        "n":1,
        "scriptPubKey":{
            "hex":"76a914acf337989ce75f7a56cd1857c617385d47aa4c3388ac",
            "asm":"OP_DUP OP_HASH160 acf337989ce75f7a56cd1857c617385d47aa4c33 OP_EQUALVERIFY OP_CHECKSIG",
            "addresses":[
                "1GmUb2B4hcCqUn5vLLud3g3kYRKmr5MsyP"
            ]
        },
        "spentTxId":null,
        "spentIndex":null,
        "spentHeight":null
    }
],
"blockhash":"000000000000000001427bb7e503686ac2a22466daaf4b09ddc70aefe8be662f",
"blockheight":529443,
"confirmations":1,
"time":1525851383,
"blocktime":1525851383,
"valueOut":12.60799106,
"size":1,
"valueIn":12.60799332,
"fees":0.00000226
}
```
##### 根据地址获取账户余额：

```
地址：/api/unspentinfo/[:type]/[:addr]
请求方式：GET
参数：type,币种[BTC,BCH,BSV]
参数：addr,收款地址
返回数据示例：/api/unspentinfo/BSV/1AiBYt8XbsdyPAELFpcSwRpu45eb2bArMf
返回
{ addrStr: '1AiBYt8XbsdyPAELFpcSwRpu45eb2bArMf',
  balance: 5000000000,
  totalReceived: 5000000000,
  totalSent: 0,
  txApperances: 1,
  unconfirmedBalance: 0,
  unconfirmedTxApperances: 0
 }

```

##### 根据地址查询该地址交易信息
```
地址：/api/address/browser_txs/[:type]/[:addrs]/[:start]/[:end]
请求方式：GET
参数：type,币种[BTC,BCH,BSV]
参数：[:addrs] 地址
参数:start,开始交易位置
参数:end,结束交易位置
返回参数示例：/api/utxo/BSV/1C5epYhuhg2UDHgR3BQmEkcMRUBxwbyUe4/0/20
返回：
[ { Hash: '7ea1d2304f1f95fae773ed8ef67b51cfd5ab33ea8b6ab0a932ee3e248b7ba74c',
    Height: 78,
    VinAddrValue: [ '@0' ],
    VoutAddrValue: [ '1AiBYt8XbsdyPAELFpcSwRpu45eb2bArMf@5000000000' ],
    Fee: -5000000000,
    CoinDay: 0
   }
]
```

##### 根据地址获取可以花费的交易
```
地址：/api/utxo/[:type]/[:addrs]
请求方式：GET
参数：type,币种[BTC,BCH,BSV]
参数：[:addrs] 地址
返回参数示例：/api/utxo/BSV/1C5epYhuhg2UDHgR3BQmEkcMRUBxwbyUe4 返回：
[
{
    "address":"1C5epYhuhg2UDHgR3BQmEkcMRUBxwbyUe4",
    "txid":"fe0c8ffb40c05b319dc8d08d1981934acac591e17e7e5060fa82753f44e0d043",
    "vout":1,
    "scriptPubKey":"76a914798a9e299f2971169a5259b6048bf1defe03582088ac",
    "amount":1.88361101,
    "satoshis":188361101,
    "confirmations":0,
    "ts":1525839539
},
{
    "address":"1C5epYhuhg2UDHgR3BQmEkcMRUBxwbyUe4",
    "txid":"fe0c8ffb40c05b319dc8d08d1981934acac591e17e7e5060fa82753f44e0d043",
    "vout":1,
    "scriptPubKey":"76a914798a9e299f2971169a5259b6048bf1defe03582088ac",
    "amount":1.88361101,
    "satoshis":188361101,
    "height":529414,
    "confirmations":30
}
]
```

