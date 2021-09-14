const bsv = require('bsv')
const { swapEndianness } = require('buffer-swap-endianness')

function VerifyMerkleProof (merkleProof, mapHashToHeader) {
  // flags:

  let txid
  if (merkleProof.txOrId.length === 64) {
    // The `txOrId` field contains a transaction ID
    txid = merkleProof.txOrId
  } else if (merkleProof.txOrId.length > 64) {
    // The `txOrId` field contains a full transaction
    const tx = new bsv.Transaction(merkleProof.txOrId)
    txid = tx.id
  } else {
    throw new Error('invalid txOrId length - must be at least 64 chars (32 bytes)')
  }

  let merkleRoot
  if (!merkleProof.targetType || merkleProof.targetType === 'hash') {
    // The `target` field contains a block hash

    if (merkleProof.target.length !== 64) {
      throw new Error('invalid target field')
    }

    // You will need to get the block header corresponding
    // to this block hash in order to get the merkle root
    // from it. You can get this from from the headers
    // store of an SPV client or from a third party
    // provider like WhatsOnChain
    const header = mapHashToHeader[merkleProof.target]
    if (!header) {
      throw new Error('block hash map to header not found in `mapHashToHeader`')
    }
    merkleRoot = extractMerkleRootFromBlockHeader(header)
  } else if (merkleProof.targetType === 'header' && merkleProof.target.length === 160) {
    // The `target` field contains a block header
    merkleRoot = extractMerkleRootFromBlockHeader(merkleProof.target)
  } else if (merkleProof.targetType === 'merkleRoot' && merkleProof.target.length === 64) {
    // the `target` field contains a merkle root
    merkleRoot = merkleProof.target
  } else {
    throw new Error('invalid targetType or target field')
  }

  if (merkleProof.proofType && merkleProof.proofType !== 'branch') {
    throw new Error('only merkle branch supported in this version') // merkle tree proof type not supported
  }

  if (merkleProof.composite === true) { // OR if (merkleProof.composite && merkleProof.composite !== false)
    throw new Error('only single proof supported in this version') // composite proof type not supported
  }

  if (!txid) {
    throw new Error('txid missing')
  }

  if (!merkleRoot) {
    throw new Error('merkleRoot missing')
  }

  const nodes = merkleProof.nodes // different nodes used in the merkle proof
  let index = merkleProof.index // index of node in current layer (will be changed on every iteration)
  let c = txid // first calculated node is the txid of the tx to prove
  let isLastInTree = true

  nodes.forEach(p => {
    // Check if the node is the left or the right child
    const cIsLeft = index % 2 === 0

    // Check for duplicate hash - this happens if the node (p) is
    // the last element of an uneven merkle tree layer
    if (p === '*') {
      if (!cIsLeft) { // this shouldn't happen...
        throw new Error('invalid duplicate on left hand side according to index value')
      }
      p = c
    }

    // This check fails at least once if it's not the last element
    if (cIsLeft && c !== p) {
      isLastInTree = false
    }

    // Calculate the parent node
    if (cIsLeft) {
      // Concatenate left leaf (c) with right leaf (p)
      c = getMerkleTreeParent(c, p)
    } else {
      // Concatenate left leaf (p) with right leaf (c)
      c = getMerkleTreeParent(p, c)
    }

    // We need integer division here with remainder dropped.
    // Javascript does floating point math by default so we
    // need to use Math.floor to drop the fraction.
    // In most languages we would use: i = i / 2;
    index = Math.floor(index / 2)
  })

  // c is now the calculated merkle root
  return {
    proofValid: c === merkleRoot,
    isLastInTree
  }
}

function getMerkleTreeParent (leftNode, rightNode) {
  // swap endianness before concatenating
  const leftConc = swapEndianness(Buffer.from(leftNode, 'hex'))
  const rightConc = swapEndianness(Buffer.from(rightNode, 'hex'))

  // concatenate leaves
  const concat = Buffer.concat([leftConc, rightConc])

  // hash the concatenation
  const hash = bsv.crypto.Hash.sha256sha256(concat)

  // swap endianness at the end and convert to hex string
  return swapEndianness(Buffer.from(hash, 'hex')).toString('hex')
}

function extractMerkleRootFromBlockHeader (blockHeader) {
  const blockHeaderBytes = Buffer.from(blockHeader, 'hex')

  // extract the 32 bytes that come after the version (4 bytes)
  // and the previous block hash (32 bytes).
  // https://en.bitcoin.it/wiki/Block_hashing_algorithm
  return swapEndianness(blockHeaderBytes.slice(36, 68)).toString('hex')
}

exports.VerifyMerkleProof = VerifyMerkleProof
