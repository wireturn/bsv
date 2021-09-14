package wallet

import (
	"bytes"
	"encoding/binary"

	"github.com/tokenized/pkg/bitcoin"

	"github.com/pkg/errors"
)

var (
	ErrKeyNotFound = errors.New("Key not found")
)

type KeyStore struct {
	Keys map[bitcoin.Hash20]*Key
}

func NewKeyStore() *KeyStore {
	return &KeyStore{
		Keys: make(map[bitcoin.Hash20]*Key),
	}
}

func (k KeyStore) Add(key *Key) error {
	hash, err := bitcoin.NewHash20(bitcoin.Hash160(key.Key.PublicKey().Bytes()))
	if err != nil {
		return err
	}
	k.Keys[*hash] = key
	return nil
}

func (k KeyStore) Remove(key *Key) error {
	hash, err := bitcoin.NewHash20(bitcoin.Hash160(key.Key.PublicKey().Bytes()))
	if err != nil {
		return err
	}
	delete(k.Keys, *hash)
	return nil
}

func (k KeyStore) RemoveAddress(ra bitcoin.RawAddress) error {
	hash, err := ra.Hash()
	if err != nil {
		return errors.Wrap(err, "address hash")
	}
	delete(k.Keys, *hash)
	return nil
}

// Get returns the key corresponding to the specified address.
func (k KeyStore) Get(address bitcoin.RawAddress) (*Key, error) {
	hash, err := address.Hash()
	if err != nil {
		return nil, err
	}
	key, ok := k.Keys[*hash]
	if !ok {
		return nil, ErrKeyNotFound
	}
	return key, nil
}

func (k KeyStore) GetAddresses() []bitcoin.RawAddress {
	result := make([]bitcoin.RawAddress, 0, len(k.Keys))
	for _, key := range k.Keys {
		result = append(result, key.Address)
	}
	return result
}

func (k KeyStore) GetAll() []*Key {
	result := make([]*Key, 0, len(k.Keys))
	for _, key := range k.Keys {
		result = append(result, key)
	}
	return result
}

func (k *KeyStore) Serialize(buf *bytes.Buffer) error {
	count := uint32(len(k.Keys))
	if err := binary.Write(buf, binary.LittleEndian, &count); err != nil {
		return err
	}

	for _, key := range k.Keys {
		if err := key.Write(buf); err != nil {
			return err
		}
	}

	return nil
}

func (k *KeyStore) Deserialize(buf *bytes.Reader) error {
	var count uint32
	if err := binary.Read(buf, binary.LittleEndian, &count); err != nil {
		return err
	}

	for i := uint32(0); i < count; i++ {
		var newKey Key
		if err := newKey.Read(buf, bitcoin.InvalidNet); err != nil {
			return err
		}

		hash, err := bitcoin.NewHash20(bitcoin.Hash160(newKey.Key.PublicKey().Bytes()))
		if err != nil {
			return err
		}
		k.Keys[*hash] = &newKey
	}

	return nil
}
