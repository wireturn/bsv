const proofing = require('../verifyMerkleProofJSON')
var assert = require('assert')

// These tests are not intended to ensure correctness of this specific implementation. They are
// intended to document as comprehensively as possible all the different test vectors for the
// protocol so that any language author can take the set of inputs and outputs and use them
// to test their own implementation in their own language with the aim of ensuring consistent
// and correct behaviour.

// As valid JSON should use " rather than ' to deliminate strings, the "Standard" JS coding
// style is broken if we wish to retain JSON compatibility with the data structures used
// for testing below.

describe('verifyMerkleProof in JSON format', () => {
  describe('target type independent cases', () => {
    it('should reject unrecognized target types', () => {
      assert.throws(() => {
        proofing.VerifyMerkleProof({
          index: 0,
          targetType: 'some made up thing',
          txOrId: '75edb0a69eb195cdd81e310553aa4d25e18450e08f168532a2c2e9cf447bf169',
          target: 'cafecafecafecafecafecafecafecafecafecafecafecafecafecafecafecafe',
          nodes: [
          ]
        }, {})
      }, Error)
    })

    it('should reject unrecognized proof types', () => {
      assert.throws(() => {
        proofing.VerifyMerkleProof({
          index: 0,
          proofType: 'some made up thing',
          targetType: 'merkleRoot',
          txOrId: '75edb0a69eb195cdd81e310553aa4d25e18450e08f168532a2c2e9cf447bf169',
          target: 'cafecafecafecafecafecafecafecafecafecafecafecafecafecafecafecafe',
          nodes: [
          ]
        }, {})
      }, Error)
    })

    it('should reject the "tree" proof type', () => {
      assert.throws(() => {
        proofing.VerifyMerkleProof({
          index: 0,
          proofType: 'tree',
          targetType: 'merkleRoot',
          txOrId: '75edb0a69eb195cdd81e310553aa4d25e18450e08f168532a2c2e9cf447bf169',
          target: 'cafecafecafecafecafecafecafecafecafecafecafecafecafecafecafecafe',
          nodes: [
          ]
        }, {})
      }, Error)
    })

    it('should accept the "branch" proof type', () => {
      assert.doesNotThrow(() => {
        proofing.VerifyMerkleProof({
          index: 0,
          proofType: 'branch',
          targetType: 'merkleRoot',
          txOrId: '75edb0a69eb195cdd81e310553aa4d25e18450e08f168532a2c2e9cf447bf169',
          target: 'cafecafecafecafecafecafecafecafecafecafecafecafecafecafecafecafe',
          nodes: [
          ]
        }, {})
      }, Error)
    })

    it('should reject composite proofs', () => {
      assert.throws(() => {
        proofing.VerifyMerkleProof({
          index: 0,
          composite: true,
          targetType: 'merkleRoot',
          txOrId: '75edb0a69eb195cdd81e310553aa4d25e18450e08f168532a2c2e9cf447bf169',
          target: 'cafecafecafecafecafecafecafecafecafecafecafecafecafecafecafecafe',
          nodes: [
          ]
        }, {})
      }, Error)
    })

    it('should reject an invalid merkle root', () => {
      // 'index' can be any value as it is only relevant in the presence of node entries.
      const result = proofing.VerifyMerkleProof({
        index: 0,
        targetType: 'merkleRoot',
        txOrId: '75edb0a69eb195cdd81e310553aa4d25e18450e08f168532a2c2e9cf447bf169',
        target: 'cafecafecafecafecafecafecafecafecafecafecafecafecafecafecafecafe',
        nodes: [
        ]
      }, {})
      assert.notDeepStrictEqual(result, {
        isLastInTree: true,
        proofValid: true
      })
    })

    it('should accept a single candidate', () => {
      // 'index' can be any value as it is only relevant in the presence of node entries.
      const result = proofing.VerifyMerkleProof({
        index: 0,
        targetType: 'merkleRoot',
        txOrId: '75edb0a69eb195cdd81e310553aa4d25e18450e08f168532a2c2e9cf447bf169',
        target: '75edb0a69eb195cdd81e310553aa4d25e18450e08f168532a2c2e9cf447bf169',
        nodes: [
        ]
      }, {})
      assert.deepStrictEqual(result, {
        isLastInTree: true,
        proofValid: true
      })
    })

    it('should accept a single pair where the transaction is on the left', () => {
      const result = proofing.VerifyMerkleProof({
        index: 0,
        targetType: 'merkleRoot',
        txOrId: '75edb0a69eb195cdd81e310553aa4d25e18450e08f168532a2c2e9cf447bf169',
        target: 'd04a34a62eeff3975e11a83b61d36b01a2e05649bb4522b5176fab86f7452f82',
        nodes: [
          'cafecafecafecafecafecafecafecafecafecafecafecafecafecafecafecafe'
        ]
      }, {})
      assert.deepStrictEqual(result, {
        isLastInTree: false,
        proofValid: true
      })
    })

    it('should accept a single pair where the transaction is on the right', () => {
      const result = proofing.VerifyMerkleProof({
        index: 1,
        targetType: 'merkleRoot',
        txOrId: '75edb0a69eb195cdd81e310553aa4d25e18450e08f168532a2c2e9cf447bf169',
        target: '51d7f11dd68ddeef4f0c29b5c850a144ece9c7c9da2e50aa91e9ccdb4dda892e',
        nodes: [
          'cafecafecafecafecafecafecafecafecafecafecafecafecafecafecafecafe'
        ]
      }, {})
      assert.deepStrictEqual(result, {
        isLastInTree: true,
        proofValid: true
      })
    })

    it('should accept a last element of an uneven tree if on the left', () => {
      const result = proofing.VerifyMerkleProof({
        index: 2,
        targetType: 'merkleRoot',
        txOrId: '75edb0a69eb195cdd81e310553aa4d25e18450e08f168532a2c2e9cf447bf169',
        target: '6c9d85bf51ebb0c474616fad91a115590e9a8316f21cab836dc949cfa267b0a7',
        nodes: [
          '*',
          '0afecafecafecafecafecafecafecafecafecafecafecafecafecafecafecafe',
          '1afecafecafecafecafecafecafecafecafecafecafecafecafecafecafecafe'
        ]
      }, {})
      assert.deepStrictEqual(result, {
        isLastInTree: false,
        proofValid: true
      })
    })

    it('should reject a last element of an uneven tree if on the right', () => {
      assert.throws(() => {
        proofing.VerifyMerkleProof({
          index: 3,
          targetType: 'merkleRoot',
          txOrId: '75edb0a69eb195cdd81e310553aa4d25e18450e08f168532a2c2e9cf447bf169',
          target: 'f58e29ffdbbc6bc996f265710b08460cf1c423fdc0384233603573a78a78d5ca',
          nodes: [
            '*',
            '2afecafecafecafecafecafecafecafecafecafecafecafecafecafecafecafe',
            '3afecafecafecafecafecafecafecafecafecafecafecafecafecafecafecafe'
          ]
        }, {})
      })
    })
  })

  describe('"header" target type related cases', () => {
    it('should accept a valid header with explicit "header" target type', () => {
      const result = proofing.VerifyMerkleProof({
        targetType: 'header',
        proofType: 'branch',
        composite: false,
        index: 12,
        txOrId: '0200000001080e8558d7af4763fef68042ef1e723d521948a0fb465237d5fb21fafb61f0580000000049483045022100fb4c94dc29cfa7423775443f8d8bb49b5814dcf709553345fcfad240efce22920220558569f97acd0d2b7bbe1954d570b9629ddf5491d9341867d7c41a8e6ee4ed2a41feffffff0200e1f505000000001976a914e296a740f5d9ecc22e0a74f9799f54ec44ee215a88ac80dc4a1f000000001976a914c993ce218b406cb71c60bad1f2be9469d91593cd88ac85020000',
        target: '000000208e33a53195acad0ab42ddbdbe3e4d9ca081332e5b01a62e340dbd8167d1a787b702f61bb913ac2063e0f2aed6d933d3386234da5c8eb9e30e498efd25fb7cb96fff12c60ffff7f2001000000',
        nodes: [
          'b9ef07a62553ef8b0898a79c291b92c60f7932260888bde0dab2dd2610d8668e',
          '0fc1c12fb1b57b38140442927fbadb3d1e5a5039a5d6db355ea25486374f104d',
          '60b0e75dd5b8d48f2d069229f20399e07766dd651ceeed55ee3c040aa2812547',
          'c0d8dbda46366c2050b430a05508a3d96dc0ed55aea685bb3d9a993f8b97cc6f',
          '391e62b3419d8a943f7dbc7bddc90e30ec724c033000dc0c8872253c27b03a42'
        ]
      })
      assert.deepStrictEqual(result, {
        isLastInTree: false,
        proofValid: true
      })
    })

    it('should reject an invalid header with explicit "header" target type', () => {
      const result = proofing.VerifyMerkleProof({
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
      assert.deepStrictEqual(result, {
        isLastInTree: false,
        proofValid: false
      })
    })
  })

  describe('"hash" target type related cases', () => {
    it('should accept a valid transaction as implicit target type', () => {
      const mapHashToHeader = {
        '75edb0a69eb195cdd81e310553aa4d25e18450e08f168532a2c2e9cf447bf169': '000000208e33a53195acad0ab42ddbdbe3e4d9ca081332e5b01a62e340dbd8167d1a787b702f61bb913ac2063e0f2aed6d933d3386234da5c8eb9e30e498efd25fb7cb96fff12c60ffff7f2001000000'
      }

      const result = proofing.VerifyMerkleProof({
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
      }, mapHashToHeader)
      assert.deepStrictEqual(result, {
        isLastInTree: false,
        proofValid: true
      })
    })

    it('should accept a valid transaction with explicit "hash" target type', () => {
      const mapHashToHeader = {
        '75edb0a69eb195cdd81e310553aa4d25e18450e08f168532a2c2e9cf447bf169': '000000208e33a53195acad0ab42ddbdbe3e4d9ca081332e5b01a62e340dbd8167d1a787b702f61bb913ac2063e0f2aed6d933d3386234da5c8eb9e30e498efd25fb7cb96fff12c60ffff7f2001000000'
      }

      const result = proofing.VerifyMerkleProof({
        index: 12,
        txOrId: '0200000001080e8558d7af4763fef68042ef1e723d521948a0fb465237d5fb21fafb61f0580000000049483045022100fb4c94dc29cfa7423775443f8d8bb49b5814dcf709553345fcfad240efce22920220558569f97acd0d2b7bbe1954d570b9629ddf5491d9341867d7c41a8e6ee4ed2a41feffffff0200e1f505000000001976a914e296a740f5d9ecc22e0a74f9799f54ec44ee215a88ac80dc4a1f000000001976a914c993ce218b406cb71c60bad1f2be9469d91593cd88ac85020000',
        targetType: 'hash',
        target: '75edb0a69eb195cdd81e310553aa4d25e18450e08f168532a2c2e9cf447bf169',
        nodes: [
          'b9ef07a62553ef8b0898a79c291b92c60f7932260888bde0dab2dd2610d8668e',
          '0fc1c12fb1b57b38140442927fbadb3d1e5a5039a5d6db355ea25486374f104d',
          '60b0e75dd5b8d48f2d069229f20399e07766dd651ceeed55ee3c040aa2812547',
          'c0d8dbda46366c2050b430a05508a3d96dc0ed55aea685bb3d9a993f8b97cc6f',
          '391e62b3419d8a943f7dbc7bddc90e30ec724c033000dc0c8872253c27b03a42'
        ]
      }, mapHashToHeader)
      assert.deepStrictEqual(result, {
        isLastInTree: false,
        proofValid: true
      })
    })

    it('should reject an invalid transaction', () => {
      const mapHashToHeader = {
        '75edb0a69eb195cdd81e310553aa4d25e18450e08f168532a2c2e9cf447bf169': '000000208e33a53195acad0ab42ddbdbe3e4d9ca081332e5b01a62e340dbd8167d1a787b702f61bb913ac2063e0f2aed6d933d3386234da5c8eb9e30e498efd25fb7cb96fff12c60ffff7f2001000000'
      }

      const result = proofing.VerifyMerkleProof({
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
      }, mapHashToHeader)
      assert.deepStrictEqual(result, {
        isLastInTree: false,
        proofValid: false
      })
    })
  })
})
