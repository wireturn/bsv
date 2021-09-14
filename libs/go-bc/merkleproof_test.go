package bc_test

import (
	"encoding/hex"
	"testing"

	"github.com/libsv/go-bc"
	"github.com/stretchr/testify/assert"
)

func TestMerkleProofToBytes(t *testing.T) {
	t.Parallel()

	proofTests := map[string]struct {
		merkleProofJSON *bc.MerkleProof
		expected        string
	}{
		"test1": {
			&bc.MerkleProof{
				Index:  12,
				TxOrID: "ffeff11c25cde7c06d407490d81ef4d0db64aad6ab3d14393530701561a465ef",
				Target: "75edb0a69eb195cdd81e310553aa4d25e18450e08f168532a2c2e9cf447bf169",
				Nodes: []string{
					"b9ef07a62553ef8b0898a79c291b92c60f7932260888bde0dab2dd2610d8668e",
					"0fc1c12fb1b57b38140442927fbadb3d1e5a5039a5d6db355ea25486374f104d",
					"60b0e75dd5b8d48f2d069229f20399e07766dd651ceeed55ee3c040aa2812547",
					"c0d8dbda46366c2050b430a05508a3d96dc0ed55aea685bb3d9a993f8b97cc6f",
					"391e62b3419d8a943f7dbc7bddc90e30ec724c033000dc0c8872253c27b03a42",
				},
			},
			"000cef65a4611570303539143dabd6aa64dbd0f41ed89074406dc0e7cd251cf1efff69f17b44cfe9c2a23285168fe05084e1254daa5305311ed8cd95b19ea6b0ed7505008e66d81026ddb2dae0bd88082632790fc6921b299ca798088bef5325a607efb9004d104f378654a25e35dbd6a539505a1e3ddbba7f92420414387bb5b12fc1c10f00472581a20a043cee55edee1c65dd6677e09903f22992062d8fd4b8d55de7b060006fcc978b3f999a3dbb85a6ae55edc06dd9a30855a030b450206c3646dadbd8c000423ab0273c2572880cdc0030034c72ec300ec9dd7bbc7d3f948a9d41b3621e39",
		},
		"test2": {
			&bc.MerkleProof{
				Index:  5,
				TxOrID: "4848b9e94dd0e4f3173ebd6982ae7eb6b793de305d8450624b1d86c02a5c61d9",
				Target: "62ea2ebe6586c3c8f4b0a17806be932fe73816cd84c0d3ce9fe0976739e6cd46",
				Nodes: []string{
					"e3aa0230aa81abd483023886ad12790acf070e2a9f92d7f0ae3bebd90a904361",
					"f46309558d8701efa4b7c1b00b62af694e2a5c6719d21bd43f7167c8e9d12fc0",
					"39e5c80bad47d33ac369c2e4341b81b07821488dab75f4bc1651d3c3bf182a56",
				},
			},
			"0005d9615c2ac0861d4b6250845d30de93b7b67eae8269bd3e17f3e4d04de9b9484846cde6396797e09fced3c084cd1638e72f93be0678a1b0f4c8c38665be2eea6203006143900ad9eb3baef0d7929f2a0e07cf0a7912ad86380283d4ab81aa3002aae300c02fd1e9c867713fd41bd219675c2a4e69af620bb0c1b7a4ef01878d550963f400562a18bfc3d35116bcf475ab8d482178b0811b34e4c269c33ad347ad0bc8e539",
		},
	}

	for name, test := range proofTests {
		t.Run(name, func(t *testing.T) {

			expected, _ := hex.DecodeString(test.expected)

			proof, err := test.merkleProofJSON.ToBytes()

			assert.NoError(t, err)
			assert.Equal(t, expected, proof)
		})
	}
}
