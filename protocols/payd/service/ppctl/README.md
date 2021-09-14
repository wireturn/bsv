## PPCTL (Payment Protocol)

This contains structs and interfaces that allow payments to be made between a user and a merchant running a Bip-70 or BIP-270 server.

These have been implemented in the Service package which contains the business logic and some data packages that implement data store layers.

![img.png](https://github.com/moneybutton/bips/raw/master/bip-0270/Protocol_Sequence.png)

## Packages

The root `ppctl` package contains the public struct and interface definitions.

### service

Implements the service interfaces defined in the root `ppctl` package to enforce the business logic of the various domains within the payment protocol.

They are built to support paymail and wallet payments.

The services then call the storer interfaces to store and retrieve data.

### data

Contains implementations of the Storer interfaces defined in the root `ppctl` package.

At the moment we support sqlite and a few noop implementations.

In future, we may support other stores such as postgres and mysql.
