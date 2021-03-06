package bap

// Example BOB tx (valid signature and address)
const sampleValidBobTx = `{ "_id": "5f08ddb0f797435fbff1ddf0", "tx": { "h": "744a55a8637aa191aa058630da51803abbeadc2de3d65b4acace1f5f10789c5b" }, 
"out": [ { "i": 0, "tape": [ { "cell": [ { "op": 0, "ops": "OP_0", "i": 0, "ii": 0 }, { "op": 106, "ops": "OP_RETURN", "i": 1, "ii": 1 } ], "i": 0 }, 
{ "cell": [ { "s": "1BAPSuaPnfGnSBM3GLV9yhxUdYe4vGbdMT", "h": "31424150537561506e66476e53424d33474c56397968785564596534764762644d54", 
"b": "MUJBUFN1YVBuZkduU0JNM0dMVjl5aHhVZFllNHZHYmRNVA==", "i": 0, "ii": 2 }, { "s": "ATTEST", "h": "415454455354", "b": "QVRURVNU", "i": 1, "ii": 3 },
{ "s": "cf39fc55da24dc23eff1809e6e6cf32a0fe6aecc81296543e9ac84b8c501bac5", 
"h": "63663339666335356461323464633233656666313830396536653663663332613066653661656363383132393635343365396163383462386335303162616335",
"b": "Y2YzOWZjNTVkYTI0ZGMyM2VmZjE4MDllNmU2Y2YzMmEwZmU2YWVjYzgxMjk2NTQzZTlhYzg0YjhjNTAxYmFjNQ==", "i": 2, "ii": 4 },
{ "s": "0", "h": "30", "b": "MA==", "i": 3, "ii": 5 } ], "i": 1 }, { "cell": [ { "s": "15PciHG22SNLQJXMoSUaWVi7WSqc7hCfva", 
"h": "313550636948473232534e4c514a584d6f5355615756693757537163376843667661", "b": "MTVQY2lIRzIyU05MUUpYTW9TVWFXVmk3V1NxYzdoQ2Z2YQ==", "i": 0, "ii": 7 },
{ "s": "BITCOIN_ECDSA", "h": "424954434f494e5f4543445341", "b": "QklUQ09JTl9FQ0RTQQ==", "i": 1, "ii": 8 }, { "s": "134a6TXxzgQ9Az3w8BcvgdZyA5UqRL89da", 
"h": "31333461365458787a675139417a33773842637667645a7941355571524c38396461", "b": "MTM0YTZUWHh6Z1E5QXozdzhCY3ZnZFp5QTVVcVJMODlkYQ==", "i": 2, "ii": 9 },
{ "s": "\u001f�nm�3坨\u001b�{\u001f\t��\u0000��(ӏ��h�D��o\u000b\u0006�$�(\u001a�'i��_�\u0006YA\"\f��ޚ` + "`" + `/U.\u0012�^W�\n", 
"h": "1fe96e6df733e59da81bc07b1f098ff19fad00b3fe28d38f81e768ed44d7c16f0b06932480281ab42769bdbb5fef065941220ccfcdde9a602f552e12dc5e57d70a", 
"b": "H+lubfcz5Z2oG8B7HwmP8Z+tALP+KNOPgedo7UTXwW8LBpMkgCgatCdpvbtf7wZZQSIMz83emmAvVS4S3F5X1wo=", "i": 3, "ii": 10 } ], "i": 2 } ], "e": { "v": 0, "i": 0, "a": "false" } }, 
{ "i": 1, "tape": [ { "cell": [ { "op": 118, "ops": "OP_DUP", "i": 0, "ii": 0 }, { "op": 169, "ops": "OP_HASH160", "i": 1, "ii": 1 }, 
{ "s": "�\no;L˺��E\t^��{i\u0011}", "h": "d27f0a6f3b4ccbbacaf945095ed3eeb97b69117d", "b": "0n8KbztMy7rK+UUJXtPuuXtpEX0=", "i": 2, "ii": 2 },
{ "op": 136, "ops": "OP_EQUALVERIFY", "i": 3, "ii": 3 }, { "op": 172, "ops": "OP_CHECKSIG", "i": 4, "ii": 4 } ], "i": 0 } ], 
"e": { "v": 14492205, "i": 1, "a": "1LC16EQVsqVYGeYTCrjvNf8j28zr4DwBuk" } } ], "lock": 0, "timestamp": 1594416560292 }`

// Example BOB tx (Invalid Address)
const sampleInvalidBobTx = `{ "_id": "5f08ddb0f797435fbff1ddf0", "tx": { "h": "744a55a8637aa191aa058630da51803abbeadc2de3d65b4acace1f5f10789c5b" }, 
"out": [ { "i": 0, "tape": [ { "cell": [ { "op": 0, "ops": "OP_0", "i": 0, "ii": 0 }, { "op": 106, "ops": "OP_RETURN", "i": 1, "ii": 1 } ], "i": 0 }, 
{ "cell": [ { "s": "invalid-prefix", "h": "31424150537561506e66476e53424d33474c56397968785564596534764762644d54", 
"b": "MUJBUFN1YVBuZkduU0JNM0dMVjl5aHhVZFllNHZHYmRNVA==", "i": 0, "ii": 2 }, { "s": "ATTEST", "h": "415454455354", "b": "QVRURVNU", "i": 1, "ii": 3 },
{ "s": "cf39fc55da24dc23eff1809e6e6cf32a0fe6aecc81296543e9ac84b8c501bac5", 
"h": "63663339666335356461323464633233656666313830396536653663663332613066653661656363383132393635343365396163383462386335303162616335",
"b": "Y2YzOWZjNTVkYTI0ZGMyM2VmZjE4MDllNmU2Y2YzMmEwZmU2YWVjYzgxMjk2NTQzZTlhYzg0YjhjNTAxYmFjNQ==", "i": 2, "ii": 4 },
{ "s": "0", "h": "30", "b": "MA==", "i": 3, "ii": 5 } ], "i": 1 }, { "cell": [ { "s": "15PciHG22SNLQJXMoSUaWVi7WSqc7hCfva", 
"h": "313550636948473232534e4c514a584d6f5355615756693757537163376843667661", "b": "MTVQY2lIRzIyU05MUUpYTW9TVWFXVmk3V1NxYzdoQ2Z2YQ==", "i": 0, "ii": 7 },
{ "s": "BITCOIN_ECDSA", "h": "424954434f494e5f4543445341", "b": "QklUQ09JTl9FQ0RTQQ==", "i": 1, "ii": 8 }, { "s": "invalid-address", 
"h": "31333461365458787a675139417a33773842637667645a7941355571524c38396461", "b": "MTM0YTZUWHh6Z1E5QXozdzhCY3ZnZFp5QTVVcVJMODlkYQ==", "i": 2, "ii": 9 },
{ "s": "\u001f�nm�3坨\u001b�{\u001f\t��\u0000��(ӏ��h�D��o\u000b\u0006�$�(\u001a�'i��_�\u0006YA\"\f��ޚ` + "`" + `/U.\u0012�^W�\n", 
"h": "1fe96e6df733e59da81bc07b1f098ff19fad00b3fe28d38f81e768ed44d7c16f0b06932480281ab42769bdbb5fef065941220ccfcdde9a602f552e12dc5e57d70a", 
"b": "H+lubfcz5Z2oG8B7HwmP8Z+tALP+KNOPgedo7UTXwW8LBpMkgCgatCdpvbtf7wZZQSIMz83emmAvVS4S3F5X1wo=", "i": 3, "ii": 10 } ], "i": 2 } ], "e": { "v": 0, "i": 0, "a": "false" } }, 
{ "i": 1, "tape": [ { "cell": [ { "op": 118, "ops": "OP_DUP", "i": 0, "ii": 0 }, { "op": 169, "ops": "OP_HASH160", "i": 1, "ii": 1 }, 
{ "s": "�\no;L˺��E\t^��{i\u0011}", "h": "d27f0a6f3b4ccbbacaf945095ed3eeb97b69117d", "b": "0n8KbztMy7rK+UUJXtPuuXtpEX0=", "i": 2, "ii": 2 },
{ "op": 136, "ops": "OP_EQUALVERIFY", "i": 3, "ii": 3 }, { "op": 172, "ops": "OP_CHECKSIG", "i": 4, "ii": 4 } ], "i": 0 } ], 
"e": { "v": 14492205, "i": 1, "a": "1LC16EQVsqVYGeYTCrjvNf8j28zr4DwBuk" } } ], "lock": 0, "timestamp": 1594416560292 }`
