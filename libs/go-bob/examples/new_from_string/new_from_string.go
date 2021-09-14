package main

import (
	"log"

	"github.com/bitcoinschema/go-bob"
)

func main() {

	const sampleBobTx = `{ "_id": "5ed082db57cd6b1658b88400", "tx": { "h": "207eaadc096849e037b8944df21a8bba6d91d8445848db047c0a3f963121e19d" }, 
"in": [ { "i": 0, "tape": [ { "cell": [ { "b": "MEUCIQDwEsO9N4EJGqjlPKsv/LkKzO2MZVALQQhv0iXkjJjB1wIgC4/xF7js0rLX6VVRvFobO7zKgEmGRHmii+2dyEKoaARB", 
"s": "0E\u0002!\u0000�\u0012ý7�\t\u001a��<�/��\n��eP\u000bA\bo�%䌘��\u0002 \u000b��\u0017��Ҳ��UQ�Z\u001b;�ʀI�Dy����B�h\u0004A", "ii": 0, "i": 0 }, 
{ "b": "A+9bsilk1SnAr3SNmmOBQy8FKY56Zu0v4i55dbFQJSin", "s": "\u0003�[�)d�)��t��c�C/\u0005)�zf�/�.yu�P%(�", "ii": 1, "i": 1 } ], "i": 0 } ], 
"e": { "h": "3d1fc854830cb7f5cf4e89459f1e2f4331ffed09ad66a02ce1140c553c9d5af1", "i": 1, "a": "1FFuYLM8a66GddCG25nUbarazeMr5dnUwC" },
"seq": 4294967295 } ], "out": [ { "i": 0, "tape": [ { "cell": [ { "op": 0, "ops": "OP_0", "ii": 0, "i": 0 }, 
{ "op": 106, "ops": "OP_RETURN", "ii": 1, "i": 1 } ], "i": 0 }, { "cell": [ { "b": "5LiA54Gv6IO96Zmk5Y2D5bm05pqX", "s": "一灯能除千年暗", "ii": 2, "i": 0 },
{ "b": "NThhNTk3", "s": "58a597", "ii": 3, "i": 1 } ], "i": 1 } ], "e": { "v": 0, "i": 0, "a": "false" } }, 
{ "i": 1, "tape": [ { "cell": [ { "op": 118, "ops": "OP_DUP", "ii": 0, "i": 0 }, { "op": 169, "ops": "OP_HASH160", "ii": 1, "i": 1 },
{ "b": "nGNxXG0fpsYbMdKRFRbhw9s736g=", "s": "�cq\\m\u001f��\u001b1ґ\u0015\u0016���;ߨ", "ii": 2, "i": 2 }, 
{ "op": 136, "ops": "OP_EQUALVERIFY", "ii": 3, "i": 3 }, { "op": 172, "ops": "OP_CHECKSIG", "ii": 4, "i": 4 } ], "i": 0 } ],
"e": { "v": 111411, "i": 1, "a": "1FFuYLM8a66GddCG25nUbarazeMr5dnUwC" } } ], "lock": 0, 
"blk": { "i": 635140, "h": "0000000000000000031d01ce0a8471d6cfab81d403ba10c878f671eac28d5d39", "t": 1589607858 }, "i": 4042 }`

	b, err := bob.NewFromString(sampleBobTx)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}
	log.Printf("found tx: %s", b.Tx.H)
}
