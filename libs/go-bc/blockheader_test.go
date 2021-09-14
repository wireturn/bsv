package bc_test

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"testing"

	"github.com/libsv/go-bc"
	"github.com/stretchr/testify/assert"
)

func TestNewBlockHeader(t *testing.T) {
	ebh := &bc.BlockHeader{
		Version: 536870912,
		HashPrevBlock: func() []byte {
			t, _ := hex.DecodeString("784605133dff2eb242d5b9c5b6dc07d7b9677f2b127ed824910e89e79477a174")
			return t
		}(),
		HashMerkleRoot: func() []byte {
			t, _ := hex.DecodeString("fc9d931e8eecd947c870279571840a727924bf8cb1243587c0c22620a9afd8e9")
			return t
		}(),
		Time: 1614043423,
		Bits: func() []byte {
			t, _ := hex.DecodeString("207fffff")
			return t
		}(),
		Nonce: 0,
	}

	headerBytes := "0000002074a17794e7890e9124d87e122b7f67b9d707dcb6c5b9d542b22eff3d13054678e9d8afa92026c2c0873524b18cbf2479720a8471952770c847d9ec8e1e939dfc1f593460ffff7f2000000000"
	bh, err := bc.NewBlockHeaderFromStr(headerBytes)

	assert.NoError(t, err)
	assert.Equal(t, ebh, bh)
}

func TestBlockHeaderString(t *testing.T) {
	expectedHeader := "00000020fb9eacea87c1cc294a4f1633a45b9bfb21cf9878b439c6138d96b8ca3a856e3a37307cd123724eaa4ade23d29feea1358458d5c110275b6cca4e2b79cd14d98e39573460ffff7f2000000000"

	bh := &bc.BlockHeader{
		Version: 536870912,
		HashPrevBlock: func() []byte {
			t, _ := hex.DecodeString("3a6e853acab8968d13c639b47898cf21fb9b5ba433164f4a29ccc187eaac9efb")
			return t
		}(),
		HashMerkleRoot: func() []byte {
			t, _ := hex.DecodeString("8ed914cd792b4eca6c5b2710c1d5588435a1ee9fd223de4aaa4e7223d17c3037")
			return t
		}(),
		Time: 1614042937,
		Bits: func() []byte {
			t, _ := hex.DecodeString("207fffff")
			return t
		}(),
		Nonce: 0,
	}

	assert.Equal(t, expectedHeader, bh.String())
}

func TestBlockHeaderStringAndBytesMatch(t *testing.T) {
	headerStr := "0000002074a17794e7890e9124d87e122b7f67b9d707dcb6c5b9d542b22eff3d13054678e9d8afa92026c2c0873524b18cbf2479720a8471952770c847d9ec8e1e939dfc1f593460ffff7f2000000000"
	bh, err := bc.NewBlockHeaderFromStr(headerStr)
	assert.NoError(t, err)
	assert.Equal(t, hex.EncodeToString(bh.Bytes()), bh.String())
}

func TestBlockHeaderInvalid(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		expectedHeader string
		expErr         error
	}{
		"empty string": {
			expectedHeader: "",
			expErr:         errors.New("block header should be 80 bytes long"),
		},
		"too long": {
			expectedHeader: "00000020fb9eacea87c1cc294a4f1633a45b9bfb21cf9878b439c61123221312312312396b8ca3a856e3a37307cd123724eaa4ade23d29feea1358458d5c110275b6cca4e2b79cd14d98e39573460ffff7f2000000000",
			expErr:         errors.New("block header should be 80 bytes long"),
		},
		"too short": {
			expectedHeader: "00000020fb9eacea87c1c3a856e3a37307cd123724eaa4ade23d29feea1358458d5c110275b6cca4e2b79cd14d98e39573460ffff7f2000000000",
			expErr:         errors.New("block header should be 80 bytes long"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := bc.NewBlockHeaderFromStr(test.expectedHeader)
			assert.Error(t, err)
			assert.EqualError(t, err, test.expErr.Error())
		})
	}
}

func TestExtractMerkleRootFromBlockHeader(t *testing.T) {
	header := "000000208e33a53195acad0ab42ddbdbe3e4d9ca081332e5b01a62e340dbd8167d1a787b702f61bb913ac2063e0f2aed6d933d3386234da5c8eb9e30e498efd25fb7cb96fff12c60ffff7f2001000000"

	merkleRoot, err := bc.ExtractMerkleRootFromBlockHeader(header)

	assert.NoError(t, err)
	assert.Equal(t, merkleRoot, "96cbb75fd2ef98e4309eebc8a54d2386333d936ded2a0f3e06c23a91bb612f70")
}

func TestEncodeAndDecodeBlockHeader(t *testing.T) {
	// the genesis block
	genesisHex := "0100000000000000000000000000000000000000000000000000000000000000000000003ba3edfd7a7b12b27ac72c3e67768f617fc81bc3888a51323a9fb8aa4b1e5e4a29ab5f49ffff001d1dac2b7c"
	genesis, err := bc.NewBlockHeaderFromStr(genesisHex)
	assert.NoError(t, err)
	assert.Equal(t, genesisHex, genesis.String())
}

func TestVerifyBlockHeader(t *testing.T) {
	// the genesis block
	genesisHex := "0100000000000000000000000000000000000000000000000000000000000000000000003ba3edfd7a7b12b27ac72c3e67768f617fc81bc3888a51323a9fb8aa4b1e5e4a29ab5f49ffff001d1dac2b7c"
	header, err := hex.DecodeString(genesisHex)
	assert.NoError(t, err)
	genesis, err := bc.NewBlockHeaderFromBytes(header)
	assert.NoError(t, err)

	assert.Equal(t, genesisHex, genesis.String())
	assert.True(t, genesis.Valid())

	// change one letter
	header[0] = 222
	genesisInvalid, err := bc.NewBlockHeaderFromBytes(header)
	assert.NoError(t, err)
	assert.False(t, genesisInvalid.Valid())
}

func TestBlockHeader_MarshalJSON(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		bh      *bc.BlockHeader
		expJSON string
	}{
		"can be marshalled": {
			bh: func() *bc.BlockHeader {
				bh, _ := bc.NewBlockHeaderFromStr("00000020332bca82fa601bcee3cc00a703988b7126079a75a03ea87654d6b544df73156ae73339a84a36cc3d1c5dc452a6b6bc2ee33c3f583cab9bcf9542b199458fea4844183661ffff7f2000000000")
				return bh
			}(),
			expJSON: `{
	"version": 536870912,
	"time": 1630935108,
	"nonce": 0,
	"hashPrevBlock": "6a1573df44b5d65476a83ea0759a0726718b9803a700cce3ce1b60fa82ca2b33",
	"merkleRoot": "48ea8f4599b14295cf9bab3c583f3ce32ebcb6a652c45d1c3dcc364aa83933e7",
	"bits": "207fffff"
}`,
		},
		"can also be marshalled": {
			bh: func() *bc.BlockHeader {
				bh, _ := bc.NewBlockHeaderFromStr("00000020b21e96654f34d7b10e65d63a29f3978bc9057f1584dae6b80c0292981ece6461fd6af93252fb8b5aec241195eb3e438873969e7ad2b30dbd14149ae3d7b0a4591f183661ffff7f2002000000")
				return bh
			}(),
			expJSON: `{
	"version": 536870912,
	"time": 1630935071,
	"nonce": 2,
	"hashPrevBlock": "6164ce1e9892020cb8e6da84157f05c98b97f3293ad6650eb1d7344f65961eb2",
	"merkleRoot": "59a4b0d7e39a1414bd0db3d27a9e967388433eeb951124ec5a8bfb5232f96afd",
	"bits": "207fffff"
}`,
		},
		"nil data doesn't error": {
			bh: &bc.BlockHeader{
				Version: 0,
				Time:    0,
				Nonce:   0,
			},
			expJSON: `{
	"version": 0,
	"time": 0,
	"nonce": 0,
	"hashPrevBlock": "",
	"merkleRoot": "",
	"bits": ""
}`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			bhj, err := json.MarshalIndent(test.bh, "", "\t")
			assert.NoError(t, err)
			assert.Equal(t, test.expJSON, string(bhj))
		})
	}
}

func TestBlockHeader_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		bh  *bc.BlockHeader
		err error
	}{
		"valid zero data can be unmarshalled": {
			bh: &bc.BlockHeader{
				Version:        0,
				Nonce:          0,
				Time:           0,
				Bits:           []byte{},
				HashMerkleRoot: []byte{},
				HashPrevBlock:  []byte{},
			},
		},
		"valid data can be unmarshalled": {
			bh: &bc.BlockHeader{
				Version:        536870912,
				Nonce:          1630935071,
				Time:           1630935071,
				Bits:           []byte{0x20, 0x7f, 0xff, 0xff},
				HashMerkleRoot: []byte{0x59, 0xa4, 0xb0, 0xd7, 0xe3, 0x9a, 0x14, 0x14, 0xbd, 0x0d, 0xb3, 0xd2, 0x7a, 0x9e, 0x96, 0x73, 0x88, 0x43, 0x3e, 0xeb, 0x95, 0x11, 0x24, 0xec, 0x5a, 0x8b, 0xfb, 0x52, 0x32, 0xf9, 0x6a, 0xfd},
				HashPrevBlock:  []byte{0x61, 0x64, 0xce, 0x1e, 0x98, 0x92, 0x02, 0x0c, 0xb8, 0xe6, 0xda, 0x84, 0x15, 0x7f, 0x05, 0xc9, 0x8b, 0x97, 0xf3, 0x29, 0x3a, 0xd6, 0x65, 0x0e, 0xb1, 0xd7, 0x34, 0x4f, 0x65, 0x96, 0x1e, 0xb2},
			},
		},
		"nil data does not error": {
			bh: &bc.BlockHeader{
				Version: 536870912,
				Nonce:   1630935071,
				Time:    1630935071,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			b, err := json.Marshal(test.bh)
			if test.err != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, test.err.Error())
				return
			}
			assert.NoError(t, err)

			var bh *bc.BlockHeader
			assert.NoError(t, json.Unmarshal(b, &bh))
			assert.Equal(t, test.bh.String(), bh.String())
		})
	}
}
