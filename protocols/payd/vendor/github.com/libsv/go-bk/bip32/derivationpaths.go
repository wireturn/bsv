package bip32

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	numericPlusTick = regexp.MustCompile(`^[0-9]+'{0,1}$`)
)

// DerivePath given an uint64 number will generate a hardened BIP32 path 3 layers deep.
//
// This is achieved by the following process:
// We split the seed bits into 3 sections: (b63-b32|b32-b1|b1-b0)
// Each section is then added onto 2^31 and concatenated together which will give us the final path.
func DerivePath(i uint64) string {
	path := fmt.Sprintf("%d/", i>>33|1<<31)
	path += fmt.Sprintf("%d/", ((i<<31)>>33)|1<<31)
	path += fmt.Sprintf("%d", (i&3)|1<<31)
	return path
}

// DeriveNumber when given a derivation path of format 0/0/0 will
// reverse the DerivePath function and return the number used to generate
// the path.
func DeriveNumber(path string) (uint64, error) {
	ss := strings.Split(path, "/")
	if len(ss) != 3 {
		return 0, errors.New("path must have 3 levels ie 0/0/0")
	}
	d1, err := strconv.ParseUint(ss[0], 10, 32)
	if err != nil {
		return 0, err
	}
	seed := (d1 - 1<<31) << 33
	d2, err := strconv.ParseUint(ss[1], 10, 32)
	if err != nil {
		return 0, err
	}
	seed += (d2 - (1 << 31)) << 2
	d3, err := strconv.ParseUint(ss[2], 10, 32)
	if err != nil {
		return 0, err
	}
	seed += d3 - (1 << 31)
	return seed, nil
}

// DeriveChildFromPath will generate a new extended key derived from the key k using the
// bip32 path provided, ie "1234/0/123"
// Child keys must be ints or hardened keys followed by '.
// https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki
func (k *ExtendedKey) DeriveChildFromPath(derivationPath string) (*ExtendedKey, error) {
	if derivationPath == "" {
		return k, nil
	}
	key := k
	children := strings.Split(derivationPath, "/")
	for _, child := range children {
		if !numericPlusTick.MatchString(child) {
			return nil, fmt.Errorf("invalid path: %q", derivationPath)
		}
		childInt, err := childInt(child)
		if err != nil {
			return nil, fmt.Errorf("derive key failed %w", err)
		}
		key, err = key.Child(childInt)
		if err != nil {
			return nil, fmt.Errorf("derive key failed %w", err)
		}
	}
	return key, nil
}

// DerivePublicKeyFromPath will generate a new extended key derived from the key k using the
// bip32 path provided, ie "1234/0/123". It will then transform to an bec.PublicKey before
// serialising the bytes and returning.
func (k *ExtendedKey) DerivePublicKeyFromPath(derivationPath string) ([]byte, error) {
	key, err := k.DeriveChildFromPath(derivationPath)
	if err != nil {
		return nil, err
	}
	pubKey, err := key.ECPubKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate public key %w", err)
	}
	return pubKey.SerialiseCompressed(), nil
}

func childInt(child string) (uint32, error) {
	var suffix uint32
	if strings.HasSuffix(child, "'") {
		child = strings.TrimRight(child, "'")
		suffix = 2147483648 // 2^32
	}
	t, err := strconv.ParseUint(child, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to get child int %w", err)
	}
	return uint32(t) + suffix, nil
}
