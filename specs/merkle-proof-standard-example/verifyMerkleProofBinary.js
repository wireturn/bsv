const bsv = require('bsv')
const { swapEndianness } = require('buffer-swap-endianness')

function VerifyMerkleProof (data, mapHashToHeader) {
  let buffer
  if (Buffer.isBuffer(data)) {
    buffer = data
  } else {
    // In theory we should check that it is an { key: value } js thing but who knows how to do that.
    buffer = packObject(data)
  }

  let txHash
  let txLength = 32
  let bufferIndex = 0
  const flags = buffer.readUInt8(bufferIndex)
  bufferIndex = bufferIndex + 1

  if ((flags & 0x08) !== 0) {
    throw new Error('only merkle branch supported in this version') // merkle tree proof type not supported
  }

  if ((flags & 0x10) !== 0) {
    throw new Error('only single proof supported in this version') // composite proof type not supported
  }

  let index
  [index, bufferIndex] = readVarInt(buffer, bufferIndex)
  if ((flags & 1) !== 0) {
    ;[txLength, bufferIndex] = readVarInt(buffer, bufferIndex)
  }
  if (txLength < 32) {
    throw new Error('invalid txOrId length - must be at least 32 bytes')
  }
  const txBuffer = buffer.slice(bufferIndex, bufferIndex + txLength)
  bufferIndex = bufferIndex + txLength
  if (txLength === 32) {
    txHash = txBuffer
  } else {
    const tx = new bsv.Transaction(txBuffer)
    txHash = tx._getHash()
  }

  let targetLength = 32 // block hash (0x00) or merkle root (0x04)
  const targetType = flags & (0x02 | 0x04)
  if (targetType === 0x02) {
    targetLength = 80 // block header (0x02)
  } else if ((targetType !== 0x00) && (targetType !== 0x04)) {
    throw new Error('invalid target type flags')
  }
  const targetBuffer = buffer.slice(bufferIndex, bufferIndex + targetLength)
  bufferIndex = bufferIndex + targetLength

  let rootBuffer
  if (targetType === 0x04) {
    rootBuffer = targetBuffer
  } else if (targetType === 0x02) {
    rootBuffer = extractMerkleRootFromBlockHeader(targetBuffer)
  } else if (targetType === 0x00) {
    const headerHashHex = swapEndianness(targetBuffer).toString('hex')
    const headerHex = mapHashToHeader[headerHashHex]
    if (!headerHex) {
      throw new Error('block hash map to header not found in `mapHashToHeader`')
    }
    const headerBuffer = Buffer.from(headerHex, 'hex')
    rootBuffer = extractMerkleRootFromBlockHeader(headerBuffer)
  }

  let c = txHash // first calculated node is the txHash of the tx to prove
  let isLastInTree = true

  let nodeCount
  [nodeCount, bufferIndex] = readVarInt(buffer, bufferIndex)

  for (let i = 0; i < nodeCount; i = i + 1) {
    const nodeType = buffer.readUInt8(bufferIndex)
    bufferIndex = bufferIndex + 1

    // Check if the node is the left or the right child
    const cIsLeft = index % 2 === 0
    let p

    // Check for duplicate hash - this happens if the node (p) is
    // the last element of an uneven merkle tree layer
    if (nodeType === 0) {
      p = buffer.slice(bufferIndex, bufferIndex + 32)
      bufferIndex = bufferIndex + 32
    } else if (nodeType === 1) {
      if (!cIsLeft) { // this shouldn't happen...
        throw new Error('invalid duplicate on left hand side according to index value')
      }
      p = c
    } else if (nodeType === 2) {
      throw new Error('what does it mean that indexes are here, I have no idea')
    } else {
      throw new Error('invalid type ' + nodeType + ' for node ' + i)
    }

    // This check fails at least once if it's not the last element
    if (cIsLeft && (Buffer.compare(c, p) !== 0)) {
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
  }
  // c is now the calculated merkle root
  return {
    proofValid: Buffer.compare(c, rootBuffer) === 0,
    isLastInTree
  }
}

function readVarInt (buffer, bufferIndex) {
  let n = buffer.readUInt8(bufferIndex)
  bufferIndex = bufferIndex + 1
  if (n < 253) { return [n, bufferIndex] }
  if (n === 253) {
    n = buffer.readUInt16(bufferIndex)
    bufferIndex = bufferIndex + 2
  } else if (n === 254) {
    n = buffer.readUInt32(bufferIndex)
    bufferIndex = bufferIndex + 4
  } else {
    n = buffer.readUInt64(bufferIndex)
    bufferIndex = bufferIndex + 8
  }
  return [n, bufferIndex]
}

function packObject (merkleProof, forcedFlagValues) {
  let flags = 0x00

  let txData = Buffer.from(merkleProof.txOrId, 'hex')
  if (txData.length === 32) {
    // The `txOrId` field contains a transaction ID
    txData = swapEndianness(txData)
  } else if (txData.length > 32) {
    // The `txOrId` field contains a full transaction
    flags |= 0x01
  } else {
    throw new Error('invalid txOrId length - must be at least 64 chars (32 bytes)')
  }

  let proofData = Buffer.from(merkleProof.target, 'hex')
  if (!merkleProof.targetType || merkleProof.targetType === 'hash') {
    if (proofData.length !== 32) {
      throw new Error('invalid target field')
    }
    proofData = swapEndianness(proofData)
  } else if (merkleProof.targetType === 'header' && proofData.length === 80) {
    flags |= 0x02
  } else if (merkleProof.targetType === 'merkleRoot' && proofData.length === 32) {
    flags |= 0x04
    proofData = swapEndianness(proofData)
  } else {
    throw new Error('invalid targetType or target field')
  }

  const forceTree = forcedFlagValues !== undefined && (forcedFlagValues & 0x08) === 0x08
  if (!forceTree && merkleProof.proofType && merkleProof.proofType !== 'branch') {
    throw new Error('only merkle branch supported in this version') // merkle tree proof type not supported
  }

  const forceComposite = forcedFlagValues !== undefined && (forcedFlagValues & 0x10) === 0x10
  if (!forceComposite && merkleProof.composite === true) { // OR if (merkleProof.composite && merkleProof.composite !== false)
    throw new Error('only single proof supported in this version') // composite proof type not supported
  }

  if (forcedFlagValues !== undefined) {
    flags = flags | forcedFlagValues
  }

  const writer = bsv.encoding.BufferWriter()
  writer.writeUInt8(flags)
  writer.writeVarintNum(merkleProof.index)
  if ((flags & 1) !== 0) {
    writer.writeVarintNum(txData.length)
  }
  writer.write(txData)
  writer.write(proofData)

  const nodeCount = merkleProof.nodes.length
  writer.writeVarintNum(nodeCount)
  merkleProof.nodes.forEach(p => {
    if (p === '*') {
      writer.writeUInt8(1)
    } else {
      writer.writeUInt8(0)
      writer.write(swapEndianness(Buffer.from(p, 'hex')))
    }
  })
  return writer.toBuffer()
}

function getMerkleTreeParent (leftBuffer, rightBuffer) {
  // concatenate leaves
  const concat = Buffer.concat([leftBuffer, rightBuffer])

  // hash the concatenation
  return bsv.crypto.Hash.sha256sha256(concat)
}

function extractMerkleRootFromBlockHeader (buffer) {
  // extract the 32 bytes that come after the version (4 bytes)
  // and the previous block hash (32 bytes).
  // https://en.bitcoin.it/wiki/Block_hashing_algorithm
  return buffer.slice(36, 68)
}

exports.VerifyMerkleProof = VerifyMerkleProof
exports.packObject = packObject
