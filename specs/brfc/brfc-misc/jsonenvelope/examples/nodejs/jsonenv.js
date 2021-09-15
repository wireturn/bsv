const bsv = require('bsv')

const privateKey = bsv.PrivateKey.fromRandom()
const publicKey = privateKey.toPublicKey().toString()

const payload = {
  name: 'simon',
  colour: 'blue'
}

// SIGN
const payloadStr = JSON.stringify(payload)
const hash = bsv.crypto.Hash.sha256(Buffer.from(payloadStr))
const sig = bsv.crypto.ECDSA.sign(hash, privateKey).toString()

const jsonEnvelope = {
  payload: payloadStr,
  signature: sig,
  publicKey: publicKey,
  encoding: 'json',
  mimetype: 'application/json'
}

// VERIFY
const verifyHash = bsv.crypto.Hash.sha256(Buffer.from(jsonEnvelope.payload))
const signature = bsv.crypto.Signature.fromString(jsonEnvelope.signature)
const verifyPublicKey = bsv.PublicKey(jsonEnvelope.publicKey)
const verified = bsv.crypto.ECDSA.verify(verifyHash, signature, verifyPublicKey)

console.log('Signature Verified: ', verified)
