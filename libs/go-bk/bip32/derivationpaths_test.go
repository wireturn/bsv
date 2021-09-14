package bip32

import (
	"errors"
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DerivePathAndDeriveSeed(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		counter      uint64
		startingPath string
		exp          string
	}{
		"successful run should return no errors": {
			counter:      0,
			startingPath: "",
			exp:          "2147483648/2147483648/2147483648",
		}, "172732732 counter with 0 path": {
			counter:      172732732,
			startingPath: "",
			exp:          "2147483648/2190666831/2147483648",
		}, "max int should return max int with root path": {
			counter:      math.MaxInt32 + 172732732,
			startingPath: "",
			exp:          "2147483648/2727537742/2147483651",
		}, "max int * 2 should return 1 with root path": {
			counter:      (math.MaxInt32 * 10000) + 172732732,
			startingPath: "",
			exp:          "2147486148/2190664331/2147483648",
		}, "max int squared should return 0/0 path": {
			counter:      (math.MaxInt32 * math.MaxInt32) + 172732732,
			startingPath: "",
			exp:          "2684354559/3264408655/2147483649",
		}, "max int squared + 100 should return correct path": {
			counter:      (math.MaxInt32*math.MaxInt32 + (math.MaxInt32 * 100)) + 172732732,
			startingPath: "",
			exp:          "2684354584/3264408630/2147483649",
		}, "max int squared plus two int32 should return correct path": {
			counter:      ((math.MaxInt32 * math.MaxInt32 * 1) + (math.MaxInt32 * 2)) + 172732732,
			startingPath: "",
			exp:          "2684354560/2190666830/2147483651",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.exp, DerivePath(test.counter))
			// assert the path can be converted correctly back to the seed.
			c, _ := DeriveNumber(test.exp)
			assert.Equal(t, test.counter, c)
		})
	}
}

func TestDeriveSeed(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		path string
		err  error
	}{
		"missing path should error": {
			err: errors.New("path must have 3 levels ie 0/0/0"),
		},
		"path too long should error": {
			path: "0/0/0/0",
			err:  errors.New("path must have 3 levels ie 0/0/0"),
		},
		"path too short should error": {
			path: "0/0",
			err:  errors.New("path must have 3 levels ie 0/0/0"),
		},
		"path length 3 should not error": {
			path: "0/0/0",
			err:  nil,
		},
		"path overflow uint32 should error": {
			path: "4294967296/0/0",
			err: &strconv.NumError{
				Func: "ParseUint",
				Num:  "4294967296",
				Err:  errors.New("value out of range"),
			},
		},
		"path less than min uint32 should error": {
			path: "-1/0/0",
			err: &strconv.NumError{
				Func: "ParseUint",
				Num:  "-1",
				Err:  errors.New("invalid syntax"),
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := DeriveNumber(test.path)
			assert.Equal(t, test.err, err)
		})
	}
}

func Test_DeriveChildFromPath(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		key     *ExtendedKey
		path    string
		expPriv string
		expPub  string
		err     error
	}{
		"successful run, 1 level child, should return no errors": {
			key: func() *ExtendedKey {
				k, err := NewKeyFromString("xprv9s21ZrQH143K3QTDL4LXw2F7HEK3wJUD2nW2nRk4stbPy6cq3jPPqjiChkVvvNKmPGJxWUtg6LnF5kejMRNNU3TGtRBeJgk33yuGBxrMPHi")
				assert.NoError(t, err)
				return k
			}(),
			path:    "0/1",
			expPriv: "xprv9ww7sMFLzJMzy7bV1qs7nGBxgKYrgcm3HcJvGb4yvNhT9vxXC7eX7WVULzCfxucFEn2TsVvJw25hH9d4mchywguGQCZvRgsiRaTY1HCqN8G",
			expPub:  "xpub6AvUGrnEpfvJBbfx7sQ89Q8hEMPM65UteqEX4yUbUiES2jHfjexmfJoxCGSwFMZiPBaKQT1RiKWrKfuDV4vpgVs4Xn8PpPTR2i79rwHd4Zr",
			err:     nil,
		}, "successful run, 2 level child, should return no errors": {
			key: func() *ExtendedKey {
				k, err := NewKeyFromString("xprv9s21ZrQH143K3QTDL4LXw2F7HEK3wJUD2nW2nRk4stbPy6cq3jPPqjiChkVvvNKmPGJxWUtg6LnF5kejMRNNU3TGtRBeJgk33yuGBxrMPHi")
				assert.NoError(t, err)
				return k
			}(),
			path:    "0/1/100000",
			expPriv: "xprv9xrdP7iD2MKJthXr1NiyGJ5KqmD2sLbYYFi49AMq9bXrKJGKBnjx5ivSzXRfLhXxzQNsqCi51oUjniwGemvfAZpzpAGohpzFkat42ohU5bR",
			expPub:  "xpub6BqyndF6risc7BcK7QFydS24Po3XGoKPuUdewYmShw4qC6bTjL4CdXEvqow6yhsfAtvU8e6kHPNFM2LzeWwKQoJm6hrYttTcxVQrk42WRE3",
			err:     nil,
		}, "successful run, 10 level child, should return no errors": {
			key: func() *ExtendedKey {
				k, err := NewKeyFromString("xprv9s21ZrQH143K3QTDL4LXw2F7HEK3wJUD2nW2nRk4stbPy6cq3jPPqjiChkVvvNKmPGJxWUtg6LnF5kejMRNNU3TGtRBeJgk33yuGBxrMPHi")
				assert.NoError(t, err)
				return k
			}(),
			path:    "0/1/1/1/1/1/1/1/1/2147483647",
			expPriv: "xprvAD89K3nZjaG8NqELN8Ce2ATWTcRADLH6JTbrXoVJT6eBRbMwbG7J75v3ym4tGC7X3Mih5krQF77pGi6GNdvxfNcr6WqYacHCSa6uzotoAx2",
			expPub:  "xpub6S7ViZKTZwpRbKJoU9jePJQF1eFecnzwfgXTLBtv1SBAJPh68oRYetEXq1RvGzsYnTzeikfdM5UM3WDrSZxuBrJi5nLpGxsuSE6cDE8pB2o",
			err:     nil,
		}, "successful run, 1 level, hardened": {
			key: func() *ExtendedKey {
				k, err := NewKeyFromString("xprv9s21ZrQH143K3QTDL4LXw2F7HEK3wJUD2nW2nRk4stbPy6cq3jPPqjiChkVvvNKmPGJxWUtg6LnF5kejMRNNU3TGtRBeJgk33yuGBxrMPHi")
				assert.NoError(t, err)
				return k
			}(),
			path:    "0/1'",
			expPriv: "xprv9ww7sMFVKxty8iXvY7Yn2NyvHZ2CgEoAYXmvf2a4XvkhzBUBmYmaMWyjyAhSxgyKK4zYzbJT6hT4JeGW5fFcNaYsBsBR9a8TxVX1LJQiZ1P",
			expPub:  "xpub6AvUGrnPALTGMCcPe95nPWveqarh5hX1ukhXTQyg6GHgryoLK65puKJDpTcMBKJKdtXQYVwbK3zMgydcTcf5qpLpJcULu9hKUxx5rzgYhrk",
			err:     nil,
		}, "successful run, 3 level, hardened": {
			key: func() *ExtendedKey {
				k, err := NewKeyFromString("xprv9s21ZrQH143K3QTDL4LXw2F7HEK3wJUD2nW2nRk4stbPy6cq3jPPqjiChkVvvNKmPGJxWUtg6LnF5kejMRNNU3TGtRBeJgk33yuGBxrMPHi")
				assert.NoError(t, err)
				return k
			}(),
			path:    "10/1'/1000'/15'",
			expPriv: "xprvA1bKm9LnkQbMvUW6kwKDLFapT9V9vTeh9D9VnVSJhRf8KmqQTc9W5YboNYcUUkZLreNq1NmeuPpw8x86C87gGyxyV6jNBV4kztFrPdSWz2t",
			expPub:  "xpub6EagAesgan9f8xaZrxrDhPXZ1BKeKvNYWS56asqvFmC7CaAZ19TkdLvHDrzubSMiC6tAqTMcumVFkgT2duhZncV3KieshEDHNc4jPWkRMGD",
			err:     nil,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			k, err := test.key.DeriveChildFromPath(test.path)
			assert.NoError(t, err)
			assert.Equal(t, test.expPriv, k.String())
			pubKey, err := k.Neuter()
			assert.NoError(t, err)
			assert.Equal(t, test.expPub, pubKey.String())
		})
	}
}
