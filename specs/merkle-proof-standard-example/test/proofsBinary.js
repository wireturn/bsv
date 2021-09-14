const proofing = require('../verifyMerkleProofBinary')
var assert = require('assert')

describe('verifyMerkleProof in binary format', () => {
  describe('format conversion', () => {
    it('should be able to convert JSON format to binary format', () => {
      const binaryDataHex = '000cef65a4611570303539143dabd6aa64dbd0f41ed89074406dc0e7cd251cf1efff69f17b44cfe9c2a23285168fe05084e1254daa5305311ed8cd95b19ea6b0ed7505008e66d81026ddb2dae0bd88082632790fc6921b299ca798088bef5325a607efb9004d104f378654a25e35dbd6a539505a1e3ddbba7f92420414387bb5b12fc1c10f00472581a20a043cee55edee1c65dd6677e09903f22992062d8fd4b8d55de7b060006fcc978b3f999a3dbb85a6ae55edc06dd9a30855a030b450206c3646dadbd8c000423ab0273c2572880cdc0030034c72ec300ec9dd7bbc7d3f948a9d41b3621e39'
      const binaryData = Buffer.from(binaryDataHex, 'hex')

      const jsBinaryData = proofing.packObject({
        index: 12,
        txOrId: 'ffeff11c25cde7c06d407490d81ef4d0db64aad6ab3d14393530701561a465ef',
        target: '75edb0a69eb195cdd81e310553aa4d25e18450e08f168532a2c2e9cf447bf169',
        nodes: [
          'b9ef07a62553ef8b0898a79c291b92c60f7932260888bde0dab2dd2610d8668e',
          '0fc1c12fb1b57b38140442927fbadb3d1e5a5039a5d6db355ea25486374f104d',
          '60b0e75dd5b8d48f2d069229f20399e07766dd651ceeed55ee3c040aa2812547',
          'c0d8dbda46366c2050b430a05508a3d96dc0ed55aea685bb3d9a993f8b97cc6f',
          '391e62b3419d8a943f7dbc7bddc90e30ec724c033000dc0c8872253c27b03a42'
        ]
      }, 0x00)

      assert.strictEqual(Buffer.compare(binaryData, jsBinaryData), 0)
    })
  })

  describe('target type independent cases', () => {
    it('should reject unrecognized flags', () => {
      // Leading byte is forced to 0xFF
      const binaryDataHex = 'ff002069f17b44cfe9c2a23285168fe05084e1254daa5305311ed8cd95b19ea6b0ed75fecafecafecafecafecafecafecafecafecafecafecafecafecafecafecafeca00'
      const binaryData = Buffer.from(binaryDataHex, 'hex')
      assert.throws(() => {
        proofing.VerifyMerkleProof(binaryData, {})
      }, Error)

      // Present the associated js version and verify that it is the same when serialized.
      const jsBinaryData = proofing.packObject({
        index: 0,
        proofType: 'some made up thing',
        targetType: 'merkleRoot',
        txOrId: '75edb0a69eb195cdd81e310553aa4d25e18450e08f168532a2c2e9cf447bf169',
        target: 'cafecafecafecafecafecafecafecafecafecafecafecafecafecafecafecafe',
        nodes: [
        ]
      }, 0xff)
      assert.strictEqual(Buffer.compare(binaryData, jsBinaryData), 0)
    })

    it('should reject the "tree" proof type', () => {
      // Leading byte is 0x04 | 0x08 with the latter indicating that it is a composite proof
      const binaryDataHex = '0c0069f17b44cfe9c2a23285168fe05084e1254daa5305311ed8cd95b19ea6b0ed75fecafecafecafecafecafecafecafecafecafecafecafecafecafecafecafeca00'
      const binaryData = Buffer.from(binaryDataHex, 'hex')
      assert.throws(() => {
        proofing.VerifyMerkleProof(binaryData, {})
      }, Error)

      // Present the associated js version and verify that it is the same when serialized.
      const jsBinaryData = proofing.packObject({
        index: 0,
        proofType: 'tree',
        targetType: 'merkleRoot',
        txOrId: '75edb0a69eb195cdd81e310553aa4d25e18450e08f168532a2c2e9cf447bf169',
        target: 'cafecafecafecafecafecafecafecafecafecafecafecafecafecafecafecafe',
        nodes: [
        ]
      }, 0x08)
      assert.strictEqual(Buffer.compare(binaryData, jsBinaryData), 0)
    })

    it('should reject composite proofs', () => {
      // Leading byte is 0x04 | 0x10 with the latter indicating that it is a composite proof
      const binaryDataHex = '140069f17b44cfe9c2a23285168fe05084e1254daa5305311ed8cd95b19ea6b0ed75fecafecafecafecafecafecafecafecafecafecafecafecafecafecafecafeca00'
      const binaryData = Buffer.from(binaryDataHex, 'hex')
      assert.throws(() => {
        proofing.VerifyMerkleProof(binaryData, {})
      }, Error)

      // Present the associated js version and verify that it is the same when serialized.
      const jsBinaryData = proofing.packObject({
        index: 0,
        composite: true,
        targetType: 'merkleRoot',
        txOrId: '75edb0a69eb195cdd81e310553aa4d25e18450e08f168532a2c2e9cf447bf169',
        target: 'cafecafecafecafecafecafecafecafecafecafecafecafecafecafecafecafe',
        nodes: [
        ]
      }, 0x10)
      assert.strictEqual(Buffer.compare(binaryData, jsBinaryData), 0)
    })

    it('should reject an invalid merkle root', () => {
      const binaryDataHex = '040069f17b44cfe9c2a23285168fe05084e1254daa5305311ed8cd95b19ea6b0ed75fecafecafecafecafecafecafecafecafecafecafecafecafecafecafecafeca00'
      const binaryData = Buffer.from(binaryDataHex, 'hex')
      const result = proofing.VerifyMerkleProof(binaryData, {})
      assert.deepStrictEqual(result, {
        isLastInTree: true,
        proofValid: false
      })

      // Present the associated js version and verify that it is the same when serialized.
      const jsBinaryData = proofing.packObject({
        index: 0,
        targetType: 'merkleRoot',
        txOrId: '75edb0a69eb195cdd81e310553aa4d25e18450e08f168532a2c2e9cf447bf169',
        target: 'cafecafecafecafecafecafecafecafecafecafecafecafecafecafecafecafe',
        nodes: [
        ]
      })
      assert.strictEqual(Buffer.compare(binaryData, jsBinaryData), 0)
    })

    it('should accept a single candidate', () => {
      const binaryDataHex = '040069f17b44cfe9c2a23285168fe05084e1254daa5305311ed8cd95b19ea6b0ed7569f17b44cfe9c2a23285168fe05084e1254daa5305311ed8cd95b19ea6b0ed7500'
      const binaryData = Buffer.from(binaryDataHex, 'hex')
      const result = proofing.VerifyMerkleProof(binaryData, {})
      assert.deepStrictEqual(result, {
        isLastInTree: true,
        proofValid: true
      })

      // Present the associated js version and verify that it is the same when serialized.
      const jsBinaryData = proofing.packObject({
        index: 0,
        targetType: 'merkleRoot',
        txOrId: '75edb0a69eb195cdd81e310553aa4d25e18450e08f168532a2c2e9cf447bf169',
        target: '75edb0a69eb195cdd81e310553aa4d25e18450e08f168532a2c2e9cf447bf169',
        nodes: [
        ]
      })
      assert.strictEqual(Buffer.compare(binaryData, jsBinaryData), 0)
    })

    it('should accept a single pair where the transaction is on the left', () => {
      const binaryDataHex = '040069f17b44cfe9c2a23285168fe05084e1254daa5305311ed8cd95b19ea6b0ed75822f45f786ab6f17b52245bb4956e0a2016bd3613ba8115e97f3ef2ea6344ad00100fecafecafecafecafecafecafecafecafecafecafecafecafecafecafecafeca'
      const binaryData = Buffer.from(binaryDataHex, 'hex')
      const result = proofing.VerifyMerkleProof(binaryData, {})
      assert.deepStrictEqual(result, {
        isLastInTree: false,
        proofValid: true
      })

      // Present the associated js version and verify that it is the same when serialized.
      const jsBinaryData = proofing.packObject({
        index: 0,
        targetType: 'merkleRoot',
        txOrId: '75edb0a69eb195cdd81e310553aa4d25e18450e08f168532a2c2e9cf447bf169',
        target: 'd04a34a62eeff3975e11a83b61d36b01a2e05649bb4522b5176fab86f7452f82',
        nodes: [
          'cafecafecafecafecafecafecafecafecafecafecafecafecafecafecafecafe'
        ]
      })
      assert.strictEqual(Buffer.compare(binaryData, jsBinaryData), 0)
    })

    it('should accept a single pair where the transaction is on the right', () => {
      const binaryDataHex = '040169f17b44cfe9c2a23285168fe05084e1254daa5305311ed8cd95b19ea6b0ed752e89da4ddbcce991aa502edac9c7e9ec44a150c8b5290c4fefde8dd61df1d7510100fecafecafecafecafecafecafecafecafecafecafecafecafecafecafecafeca'
      const binaryData = Buffer.from(binaryDataHex, 'hex')
      const result = proofing.VerifyMerkleProof(binaryData, {})
      assert.deepStrictEqual(result, {
        isLastInTree: true,
        proofValid: true
      })

      // Present the associated js version and verify that it is the same when serialized.
      const jsBinaryData = proofing.packObject({
        index: 1,
        targetType: 'merkleRoot',
        txOrId: '75edb0a69eb195cdd81e310553aa4d25e18450e08f168532a2c2e9cf447bf169',
        target: '51d7f11dd68ddeef4f0c29b5c850a144ece9c7c9da2e50aa91e9ccdb4dda892e',
        nodes: [
          'cafecafecafecafecafecafecafecafecafecafecafecafecafecafecafecafe'
        ]
      })
      assert.strictEqual(Buffer.compare(binaryData, jsBinaryData), 0)
    })

    it('should accept a last element of an uneven tree if on the left', () => {
      const binaryDataHex = '040269f17b44cfe9c2a23285168fe05084e1254daa5305311ed8cd95b19ea6b0ed75a7b067a2cf49c96d83ab1cf216839a0e5915a191ad6f6174c4b0eb51bf859d6c030100fecafecafecafecafecafecafecafecafecafecafecafecafecafecafecafe0a00fecafecafecafecafecafecafecafecafecafecafecafecafecafecafecafe1a'
      const binaryData = Buffer.from(binaryDataHex, 'hex')
      const result = proofing.VerifyMerkleProof(binaryData, {})
      assert.deepStrictEqual(result, {
        isLastInTree: false,
        proofValid: true
      })

      // Present the associated js version and verify that it is the same when serialized.
      const jsBinaryData = proofing.packObject({
        index: 2,
        targetType: 'merkleRoot',
        txOrId: '75edb0a69eb195cdd81e310553aa4d25e18450e08f168532a2c2e9cf447bf169',
        target: '6c9d85bf51ebb0c474616fad91a115590e9a8316f21cab836dc949cfa267b0a7',
        nodes: [
          '*',
          '0afecafecafecafecafecafecafecafecafecafecafecafecafecafecafecafe',
          '1afecafecafecafecafecafecafecafecafecafecafecafecafecafecafecafe'
        ]
      })
      assert.strictEqual(Buffer.compare(binaryData, jsBinaryData), 0)
    })

    it('should reject a last element of an uneven tree if on the right', () => {
      const binaryDataHex = '040369f17b44cfe9c2a23285168fe05084e1254daa5305311ed8cd95b19ea6b0ed75cad5788aa7733560334238c0fd23c4f10c46080b7165f296c96bbcdbff298ef5030100fecafecafecafecafecafecafecafecafecafecafecafecafecafecafecafe2a00fecafecafecafecafecafecafecafecafecafecafecafecafecafecafecafe3a'
      const binaryData = Buffer.from(binaryDataHex, 'hex')
      assert.throws(() => {
        proofing.VerifyMerkleProof(binaryData, {})
      }, Error)

      // Present the associated js version and verify that it is the same when serialized.
      const jsBinaryData = proofing.packObject({
        index: 3,
        targetType: 'merkleRoot',
        txOrId: '75edb0a69eb195cdd81e310553aa4d25e18450e08f168532a2c2e9cf447bf169',
        target: 'f58e29ffdbbc6bc996f265710b08460cf1c423fdc0384233603573a78a78d5ca',
        nodes: [
          '*',
          '2afecafecafecafecafecafecafecafecafecafecafecafecafecafecafecafe',
          '3afecafecafecafecafecafecafecafecafecafecafecafecafecafecafecafe'
        ]
      })
      assert.strictEqual(Buffer.compare(binaryData, jsBinaryData), 0)
    })
  })

  describe('"header" target type related cases', () => {
    it('should reject an invalid header with explicit "header" target type', () => {
      const binaryDataHex = '030cc00200000001080e8558d7af4763fef68042ef1e723d521948a0fb465237d5fb21fafb61f0580000000049483045022100fb4c94dc29cfa7423775443f8d8bb49b5814dcf709553345fcfad240efce22920220558569f97acd0d2b7bbe1954d570b9629ddf5491d9341867d7c41a8e6ee4ed2a41feffffff0200e1f505000000001976a914e296a740f5d9ecc22e0a74f9799f54ec44ee215a88ac80dc4a1f000000001976a914c993ce218b406cb71c60bad1f2be9469d91593cd88ac85020000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000005008e66d81026ddb2dae0bd88082632790fc6921b299ca798088bef5325a607efb9004d104f378654a25e35dbd6a539505a1e3ddbba7f92420414387bb5b12fc1c10f00472581a20a043cee55edee1c65dd6677e09903f22992062d8fd4b8d55de7b060006fcc978b3f999a3dbb85a6ae55edc06dd9a30855a030b450206c3646dadbd8c000423ab0273c2572880cdc0030034c72ec300ec9dd7bbc7d3f948a9d41b3621e39'
      const binaryData = Buffer.from(binaryDataHex, 'hex')
      const result = proofing.VerifyMerkleProof(binaryData)
      assert.deepStrictEqual(result, {
        isLastInTree: false,
        proofValid: false
      })

      // Present the associated js version and verify that it is the same when serialized.
      const jsBinaryData = proofing.packObject({
        targetType: 'header',
        proofType: 'branch',
        composite: false,
        index: 12,
        txOrId: '0200000001080e8558d7af4763fef68042ef1e723d521948a0fb465237d5fb21fafb61f0580000000049483045022100fb4c94dc29cfa7423775443f8d8bb49b5814dcf709553345fcfad240efce22920220558569f97acd0d2b7bbe1954d570b9629ddf5491d9341867d7c41a8e6ee4ed2a41feffffff0200e1f505000000001976a914e296a740f5d9ecc22e0a74f9799f54ec44ee215a88ac80dc4a1f000000001976a914c993ce218b406cb71c60bad1f2be9469d91593cd88ac85020000',
        target: '0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000',
        nodes: [
          'b9ef07a62553ef8b0898a79c291b92c60f7932260888bde0dab2dd2610d8668e',
          '0fc1c12fb1b57b38140442927fbadb3d1e5a5039a5d6db355ea25486374f104d',
          '60b0e75dd5b8d48f2d069229f20399e07766dd651ceeed55ee3c040aa2812547',
          'c0d8dbda46366c2050b430a05508a3d96dc0ed55aea685bb3d9a993f8b97cc6f',
          '391e62b3419d8a943f7dbc7bddc90e30ec724c033000dc0c8872253c27b03a42'
        ]
      })
      assert.strictEqual(Buffer.compare(binaryData, jsBinaryData), 0)
    })
  })

  describe('"hash" target type related cases', () => {
    it('should accept a valid transaction as implicit target type', () => {
      const mapHashToHeader = {
        '75edb0a69eb195cdd81e310553aa4d25e18450e08f168532a2c2e9cf447bf169': '000000208e33a53195acad0ab42ddbdbe3e4d9ca081332e5b01a62e340dbd8167d1a787b702f61bb913ac2063e0f2aed6d933d3386234da5c8eb9e30e498efd25fb7cb96fff12c60ffff7f2001000000'
      }

      const binaryDataHex = '010cc00200000001080e8558d7af4763fef68042ef1e723d521948a0fb465237d5fb21fafb61f0580000000049483045022100fb4c94dc29cfa7423775443f8d8bb49b5814dcf709553345fcfad240efce22920220558569f97acd0d2b7bbe1954d570b9629ddf5491d9341867d7c41a8e6ee4ed2a41feffffff0200e1f505000000001976a914e296a740f5d9ecc22e0a74f9799f54ec44ee215a88ac80dc4a1f000000001976a914c993ce218b406cb71c60bad1f2be9469d91593cd88ac8502000069f17b44cfe9c2a23285168fe05084e1254daa5305311ed8cd95b19ea6b0ed7505008e66d81026ddb2dae0bd88082632790fc6921b299ca798088bef5325a607efb9004d104f378654a25e35dbd6a539505a1e3ddbba7f92420414387bb5b12fc1c10f00472581a20a043cee55edee1c65dd6677e09903f22992062d8fd4b8d55de7b060006fcc978b3f999a3dbb85a6ae55edc06dd9a30855a030b450206c3646dadbd8c000423ab0273c2572880cdc0030034c72ec300ec9dd7bbc7d3f948a9d41b3621e39'
      const binaryData = Buffer.from(binaryDataHex, 'hex')
      const result = proofing.VerifyMerkleProof(binaryData, mapHashToHeader)
      assert.deepStrictEqual(result, {
        isLastInTree: false,
        proofValid: true
      })

      // Present the associated js version and verify that it is the same when serialized.
      const jsBinaryData = proofing.packObject({
        index: 12,
        txOrId: '0200000001080e8558d7af4763fef68042ef1e723d521948a0fb465237d5fb21fafb61f0580000000049483045022100fb4c94dc29cfa7423775443f8d8bb49b5814dcf709553345fcfad240efce22920220558569f97acd0d2b7bbe1954d570b9629ddf5491d9341867d7c41a8e6ee4ed2a41feffffff0200e1f505000000001976a914e296a740f5d9ecc22e0a74f9799f54ec44ee215a88ac80dc4a1f000000001976a914c993ce218b406cb71c60bad1f2be9469d91593cd88ac85020000',
        target: '75edb0a69eb195cdd81e310553aa4d25e18450e08f168532a2c2e9cf447bf169',
        nodes: [
          'b9ef07a62553ef8b0898a79c291b92c60f7932260888bde0dab2dd2610d8668e',
          '0fc1c12fb1b57b38140442927fbadb3d1e5a5039a5d6db355ea25486374f104d',
          '60b0e75dd5b8d48f2d069229f20399e07766dd651ceeed55ee3c040aa2812547',
          'c0d8dbda46366c2050b430a05508a3d96dc0ed55aea685bb3d9a993f8b97cc6f',
          '391e62b3419d8a943f7dbc7bddc90e30ec724c033000dc0c8872253c27b03a42'
        ]
      })
      assert.strictEqual(Buffer.compare(binaryData, jsBinaryData), 0)
    })

    it('should reject an invalid transaction', () => {
      const mapHashToHeader = {
        '75edb0a69eb195cdd81e310553aa4d25e18450e08f168532a2c2e9cf447bf169': '000000208e33a53195acad0ab42ddbdbe3e4d9ca081332e5b01a62e340dbd8167d1a787b702f61bb913ac2063e0f2aed6d933d3386234da5c8eb9e30e498efd25fb7cb96fff12c60ffff7f2001000000'
      }

      const binaryDataHex = '010c3c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000069f17b44cfe9c2a23285168fe05084e1254daa5305311ed8cd95b19ea6b0ed7505008e66d81026ddb2dae0bd88082632790fc6921b299ca798088bef5325a607efb9004d104f378654a25e35dbd6a539505a1e3ddbba7f92420414387bb5b12fc1c10f00472581a20a043cee55edee1c65dd6677e09903f22992062d8fd4b8d55de7b060006fcc978b3f999a3dbb85a6ae55edc06dd9a30855a030b450206c3646dadbd8c000423ab0273c2572880cdc0030034c72ec300ec9dd7bbc7d3f948a9d41b3621e39'
      const binaryData = Buffer.from(binaryDataHex, 'hex')
      const result = proofing.VerifyMerkleProof(binaryData, mapHashToHeader)
      assert.deepStrictEqual(result, {
        isLastInTree: false,
        proofValid: false
      })

      // Present the associated js version and verify that it is the same when serialized.
      const jsBinaryData = proofing.packObject({
        index: 12,
        txOrId: '000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000',
        target: '75edb0a69eb195cdd81e310553aa4d25e18450e08f168532a2c2e9cf447bf169',
        nodes: [
          'b9ef07a62553ef8b0898a79c291b92c60f7932260888bde0dab2dd2610d8668e',
          '0fc1c12fb1b57b38140442927fbadb3d1e5a5039a5d6db355ea25486374f104d',
          '60b0e75dd5b8d48f2d069229f20399e07766dd651ceeed55ee3c040aa2812547',
          'c0d8dbda46366c2050b430a05508a3d96dc0ed55aea685bb3d9a993f8b97cc6f',
          '391e62b3419d8a943f7dbc7bddc90e30ec724c033000dc0c8872253c27b03a42'
        ]
      })
      assert.strictEqual(Buffer.compare(binaryData, jsBinaryData), 0)
    })
  })
})
