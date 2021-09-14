# PUSHDATA data framing for OP_RETURN

This document describes a data framing scheme for `OP_RETURN` data.

## Introduction

It is commonly assumed that `PUSHDATA` operations are required to add data following an `OP_RETURN` op code.  This is not a requirement as no script execution happens after `OP_RETURN`.  It is valid in bitcoin for any sequence of bytes to follow the `OP_RETURN` op code.  However many existing bitcoin libraries attempt to parse scripts even when it is not necessary and will throw exceptions when encountering invalid scripts. `PUSHDATA` provides length framing of data elements and given it's use keeps the script in a parseable form for existing bitcoin libraries. It fulfills the requirements of a data framing standard with minimal drawbacks.

Considering uses cases where data is extracted from UTXOs then handed off to a consumer that may not be bitcoin aware, the semantics of `OP_PUSHDATA` operations are quite simple to implement.  Given a byte array `bytes` with the data element beginning at `offset` the following pseudo code will determine the length of the data that follows

```
byte b = bytes[offset];
int length = -1;
if (b < 0x4c) { //0x4c is PUSHDATA1 op code
    length = b;
    offset++;
} else if (b == 0x4c) {
    length = bytes[offset + 1];
    offset++;
} else if (b == 0x4d) { //0x4d is PUSHDATA2 op code 
    length = bytes[offset + 1] | bytes[offset + 2] << 8;
    offset += 2;
} else if (b == 0x4e { //ox4e is PUSHDATA4 op code
    length = bytes[offset + 1] | bytes[offset + 2] << 8 | bytes[offset + 3] << 16 | bytes[offset + 4] << 24;
    offset += 4;
} else {
    //this byte is not a PUSHDATA op code
    throw error;
}
byte[] data = bytes[offset] ... bytes[offset + length]
```


## Specification

This spec proposes the use of a simple data framing mechanism.  A data element is treated exactly like it is being pushed onto the stack in a running script.  That is `OP_PUSHDATAn` followed by the data itself. `OP_0` is used to designate the data element is `null` whilst maintaining positional order of following data elements.  The behaviour is described for different data lengths:

`0`: the value `0x00` is appended to the script as a length marker.  No data is appended.

`1-75`: the length value as a single byte is appended to the script followed by the data. e.g. to frame the data `0xffffffff` the length is `0x04` so `0x04ffffffff` is the framed data element.

`76-0xff`: the byte `0x4c` which is `OP_PUSHDATA1`, followed by the length as a single byte integer, followed by the actual data is appended to the script.

`0x0100-0xffff`: the byte `0x4d` which is `OP_PUSHDATA2`, followed by the length as a 2 byte integer in reverse byte order (little endian), followed by the actual data is appended to the script.

`0x010000-0xffffffff`: the byte `0x4e` which is `OP_PUSHDATA4`, followed by as a 4 byte integer in reverse byte order (little endian), followed by the actual data is appended to the script.


Examples:
To encode the ASCII string `http` which encodes to hex `0x68747470` the length field is `0x04` and the full `PUSHDATA` would be `0x0468747470`

To encode a 1024 byte data element the `PUSHDATA2` op code is `0x4d` the length field is `0x0400` and the full `PUSHDATA` would be `0x4d0400<data>`

#### _References_: 
[bitcoinj pushdata encoding implementation](https://github.com/bitcoinj-cash/bitcoinj/blob/c3e90e2e26e082d1b17f1940541dd1bda5feafd8/core/src/main/java/org/bitcoinj/script/ScriptChunk.java#L111)
[bitcoinj pushdata decoding implementation](https://github.com/bitcoinj-cash/bitcoinj/blob/c3e90e2e26e082d1b17f1940541dd1bda5feafd8/core/src/main/java/org/bitcoinj/script/Script.java#L186)

