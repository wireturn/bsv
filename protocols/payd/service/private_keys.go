package service

import (
	"context"
	"regexp"

	"github.com/libsv/go-bk/bip32"
	"github.com/libsv/go-bk/chaincfg"
	"github.com/pkg/errors"

	gopayd "github.com/libsv/payd"
)

type privateKey struct {
	store           gopayd.PrivateKeyReaderWriter
	useMainNet      bool
	numericPlusTick *regexp.Regexp
}

// NewPrivateKeys will setup and return a new PrivateKey service.
func NewPrivateKeys(store gopayd.PrivateKeyReaderWriter, useMainNet bool) *privateKey {
	return &privateKey{
		store:           store,
		useMainNet:      useMainNet,
		numericPlusTick: regexp.MustCompile(`^[0-9]+'{0,1}$`),
	}
}

// Create creates a extended private key for a keyName.
func (svc *privateKey) Create(ctx context.Context, keyName string) error { // get keyname from settings in caller
	key, err := svc.store.PrivateKey(ctx, gopayd.KeyArgs{Name: keyName})
	if err != nil {
		return errors.Wrapf(err, "failed to get key %s by name", keyName)
	}
	if key != nil {
		return nil
	}
	seed, err := bip32.GenerateSeed(bip32.RecommendedSeedLen)
	if err != nil {
		return errors.Wrap(err, "failed to generate seed")
	}
	chain := &chaincfg.TestNet
	if svc.useMainNet {
		chain = &chaincfg.MainNet
	}
	xprv, err := bip32.NewMaster(seed, chain)
	if err != nil {
		return errors.Wrap(err, "failed to create master node for given seed and chain")
	}
	if _, err := svc.store.PrivateKeyCreate(ctx, gopayd.PrivateKey{
		Name: keyName,
		Xprv: xprv.String(),
	}); err != nil {
		return errors.Wrap(err, "failed to create private key")
	}
	return nil
}

// PrivateKey returns the extended private key for a keyname.
func (svc *privateKey) PrivateKey(ctx context.Context, keyName string) (*bip32.ExtendedKey, error) {
	key, err := svc.store.PrivateKey(ctx, gopayd.KeyArgs{Name: keyName})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get key %s by name", keyName)
	}
	if key == nil {
		return nil, errors.Wrap(err, "key not found")
	}

	xKey, err := bip32.NewKeyFromString(key.Xprv)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get extended key from xpriv")
	}
	return xKey, nil
}
