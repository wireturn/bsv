bip32
==========

Package bip32 provides an API for bitcoin hierarchical deterministic
extended keys (BIP0032).

## Feature Overview

- Full BIP0032 implementation
- Single type for private and public extended keys
- Convenient cryptograpically secure seed generation
- Simple creation of master nodes
- Support for multi-layer derivation
- Easy serialization and deserialization for both private and public extended
  keys
- Support for custom networks by registering them with chaincfg
- Obtaining the underlying EC pubkeys, EC privkeys, and associated bitcoin
  addresses ties in seamlessly with existing bsvec and bsvutil types which
  provide powerful tools for working with them to do things like sign
  transactions and generate payment scripts
- Uses the bec package which is highly optimized for secp256k1
- Code examples including:
  - Generating a cryptographically secure random seed and deriving a
    master node from it
  - Default HD wallet layout as described by BIP0032
  - Audits use case as described by BIP0032
- Comprehensive test coverage including the BIP0032 test vectors
- Benchmarks

## Installation and Updating

```bash
$ go get -u github.com/libsv/go-bk/bip32
```

## Examples

* [NewMaster Example](http://godoc.org/github.com/bitcoinsv/bsvutil/hdkeychain#example-NewMaster)  
  Demonstrates how to generate a cryptographically random seed then use it to
  create a new master node (extended key).
* [Default Wallet Layout Example](http://godoc.org/github.com/bitcoinsv/bsvutil/hdkeychain#example-package--DefaultWalletLayout)  
  Demonstrates the default hierarchical deterministic wallet layout as described
  in BIP0032.
* [Audits Use Case Example](http://godoc.org/github.com/bitcoinsv/bsvutil/hdkeychain#example-package--Audits)  
  Demonstrates the audits use case in BIP0032.

## License

Package hdkeychain is licensed under the [copyfree](http://copyfree.org) ISC
License.
