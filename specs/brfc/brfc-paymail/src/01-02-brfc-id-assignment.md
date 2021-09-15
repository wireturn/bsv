# BRFC ID Assignment

It is desirable that a BRFC document be uniquely identified. Without a central authority to issue an identification number, we have chosen to borrow inspiration from Bitcoin and use hashes of content.

## ID Construction

To construct a BRFC ID from a specification, take the UTF8 string value of the `title`, `author` and `version` metadata fields (omit those not present), trim leading and trailing whitespace (leaving whitespace mid-way through the value), concatenate each value, then reinterpret the string as a byte array, and apply a double SHA256 hash.

```js
let hash = sha256d(
  spec.title.trim() +
  (spec.author || '').trim() +
  (spec.version || '').trim()
);
```

Hex-format the hash as per Bitcoin conventions (usually this means reversing the bytes before converting to hex).

```js
let bitcoinDisplayHash = hash
  .reverse()
  .toString('hex');
```

Take the first 12 characters of the Bitcoin-style display hash (representing the last six bytes of the underlying `sha256d` value):

```js
let brfcId = bitcoinDisplayHash.substring(0, 12);
```

## Considerations

Hashing the title, author and version metadata of a specification allows us to generate a unique ID without central authority. Hashing the entire specification was considered, however this was discounted due to the following drawbacks:

* Any change, however minor (like typo fixes) would create an entirely new specification id
* Different platforms handle line endings differently, and different source control and editor software can replace these without warning. This leads to unstable hashes across seemingly identical documents
* Some file formats update metadata even when content remains unchanged. Again, this would lead to unstable hashes over otherwise stable content

## Test Cases

```yaml
title: BRFC Specifications
author: andy (nChain)
version: 1
```

Expected BRFC ID: `57dd1f54fc67`

```yaml
title: bsvalias Payment Addressing (PayTo Protocol Prefix)
author: andy (nChain)
version: 1
```

Expected BRFC ID: `74524c4d6274`

```yaml
title: bsvalias Integration with Simplified Payment Protocol
author: andy (nChain)
version: 1
```

Expected BRFC ID: `0036f9b8860f`
