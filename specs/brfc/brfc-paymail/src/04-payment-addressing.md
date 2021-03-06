# Payment Addressing

Payment Addressing is the mechanism through which a wallet can, on behalf of a user making a payment transaction, discover the preferred Bitcoin output script of a receiver given only their paymail handle, in the form `<alias>@<domain>.<tld>`.

Payment Addressing is specified across a number of BRFCs:

* [Basic Address Resolution](./04-01-basic-address-resolution.md)
* [Sender Validation](./04-02-sender-validation.md)
* [Receiver Approvals](./04-03-receiver-approvals.md)
* [PayTo Protocol Prefix](./04-04-payto-protocol-prefix.md)

Only _Basic Address Resolution_ is required in order for an implementation to carry the paymail branding.
