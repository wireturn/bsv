# Protocol prefix for OP_RETURN

**Abstract:** As an _optional guideline_, we recommend that all present and future protocols to be implemented on Bitcoin SV implement a protocol prefix scheme, referred to as protocol identifiers, whenever they use `OP_RETURN`. This scheme will simplify interoperability between protocols, facilitate selective pruning of the blockchain and get built-in support from low-level infrastructure components. The proposal also provides a basic early stage process to let the ecosystem participants claim prefixes for their own protocols.

_In the following, Bitcoin always refers to Bitcoin SV._


## Overview

The Bitcoin ecosystem has gained renewed interest for overlay protocols built on top of the blockchain, leveraging the growing capacity of the `OP_RETURN` opcode of Bitcoin to embed arbitrary data in Bitcoin transaction. In order simplify parsing of structured data carried in `OP_RETURN` a simple message to identify data relevant to a use case or ignore data that is irrelevant a recommended method id prefixing messages carried through `OP_RETURN` with a protocol identifier to enable simple sorting of messages according to their respective protocols

An agreed unifying scheme will help to avoid collisions between protocols. While none of those collisions endanger Bitcoin itself, they will significantly and needlessly complicate the design of the software intended to operate those overlay protocols.

At this early stage there are at least four known protocol identifier schemes in use. However, they share a common property.  The protocol identifier is a [PUSHDATA data element](./01-PUSHDATA-data-element-framing.md) that appears immediately after the `OP_RETURN` op code. There is no length requirement so as to encapsulate several more specific protocol identifier standards currently in use under this specification.

## Specification

An arbitrary length data element containing the protocol identifier should be the first element of data following the `OP_RETURN` op code. The data element should be length prefixed as described in the [PUSHDATA framing specification](01-PUSHDATA-data-element-framing.md).

The `null` data element, that is, the byte `0x00` is reserved as a marker for protocols that do not wish to specify a protocol identifier.

A protocol that follows this specification MUST contain at least one other data element after the protocol identifier.  That element may be the `null` data element.  This is in order to provide a mechanism to distinguish protocols that comply with this protocol from `OP_RETURN` outputs that contain data embedded in a single push.

It is recommended (but not required) that protocol identifiers be less than 256 bytes to enable to `OP_RETURN` indexing services an easy mechanism to filter non-compliant data.

## Protocol identifier registry

As a courtesy to the community, we recommend to either submit a ticket or a pull request to the Git repository at https://github.com/bitcoin-sv-specs/op_return to claim your prefix with:

* A display name for the protocol
* An author (or list of authors)
* A URL pointing to the specification of the protocol
* A Bitcoin address (to avoid ambiguity and to later modify names / authors / URL).

This last step is _not_ a requirement. However, if you don’t try to make your prefix known to the community at large, be aware that it leaves you open to collisions with a useful protocol which just happens to have a lot more traction than yours.

See also _Annex: file format of /protocol_ids.csv

## Annex: file format of /protocol_ids.csv

In order to make known protocols easily available to the community at large, a simple file format is proposed to gather the protocol prefixes. 

This file `protocol_ids.csv` should be seen as an early stage effort to help various protocols gain traction within the Bitcoin community. If the number of active protocols becomes greater than a few hundred, we expect that the file `protocol_ids.csv` will be superseded with an approach more scalable than having a flat text file holding all known protocols.

The URL for the file is expected to be:

https://github.com/bitcoin-sv-specs/op_return/blob/master/etc/protocol_ids.csv 

The file is encoded in CSV as per [RFC 4180](https://tools.ietf.org/html/rfc4180) with the following options:

* UTF-8 encoding
* Unix line ending (\n)
* Comma delimiter
* Optional quote escaping for strings
* Quote escaped strings cannot contain newline (\n), returns (\r) or quotes (“)
* First line is the column headers

Then the columns themselves:

* `Prefix`: hexadecimal encoded (aka 0x01234567)
* `DisplayName`: string
* `Authors`: string
* `BitcoinAddress`: a valid Bitcoin address
* `SpecificationUrl`: string
* `TxidRedirectUrl` (optional): string (contains `{txid}`)

Each line should not be longer than 1KB (1024 bytes) in total.

The lines should be sorted in increasing order against their prefix.

The field `Prefix` must contain the hexadecimally encoded binary prefix, not including the length of the prefix.

The field `BitcoinAddress` is intended to help resolve any conflicting claims in the event where such claims were to arise.

The field `TxidRedirectUrl` is intended to help blockchain explorers making protocols more discoverable. For any transaction associated to the protocol - as identified through its prefix - a redirecting link can be inserted. The field `TxidRedirectUrl` should contain the substring `{txid}` to be replaced by the transaction identifier encoded in hexadecimal (64 characters). The landing page of the redirect is expected to be a human-readable version of the transaction aligned with the semantic of the protocol.
