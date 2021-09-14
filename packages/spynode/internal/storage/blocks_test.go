package storage

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/tokenized/pkg/bitcoin"
	"github.com/tokenized/pkg/storage"
	"github.com/tokenized/pkg/wire"
	"github.com/tokenized/spynode/internal/platform/config"
)

func TestBlocks(test *testing.T) {
	testBlockCount := 2500
	testRevertHeights := [...]int{2400, 2000, 1999, 500}

	// Generate block hashes
	blocks := make([]*bitcoin.Hash32, 0, testBlockCount)
	seed := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(seed)

	// Setup config
	startHash := "0000000000000000000000000000000000000000000000000000000000000000"
	config, err := config.NewConfig(bitcoin.MainNet, true, "test", "Tokenized Test", startHash, 8,
		2000, 10, 10, 1000, true)
	if err != nil {
		test.Errorf("Failed to create config : %v", err)
	}

	ctx := context.Background()
	store := storage.NewMockStorage()
	repo := NewBlockRepository(config, store)

	t := uint32(time.Now().Unix())
	header := wire.BlockHeader{Version: 1}
	for i := 0; i < testBlockCount; i++ {
		header.Timestamp = time.Unix(int64(t), 0)
		header.Nonce = uint32(randGen.Int())
		repo.Add(ctx, &header)
		t += 600
		blocks = append(blocks, header.BlockHash())
		header.PrevBlock = *blocks[len(blocks)-1]
	}

	if err := repo.Save(ctx); err != nil {
		test.Errorf("Failed to save repo : %v", err)
	}

	for _, revertHeight := range testRevertHeights {
		test.Logf("Test revert to (%d) : %s", revertHeight, blocks[revertHeight].String())

		if err := repo.Revert(ctx, revertHeight); err != nil {
			test.Errorf("Failed to revert repo : %v", err)
		}

		if !repo.LastHash().Equal(blocks[revertHeight]) {
			test.Errorf("Failed to revert repo to height %d", revertHeight)
		}

		if err := repo.Load(ctx); err != nil {
			test.Errorf("Failed to load repo after revert to %d : %v", revertHeight, err)
		}
	}
}

func TestSpecificReorg(test *testing.T) {
	testBlockCount := 682964
	testRevertHeights := [...]int{682960}

	// Generate block hashes
	blocks := make([]*bitcoin.Hash32, 0, testBlockCount)
	seed := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(seed)

	// Setup config
	startHash := "0000000000000000000000000000000000000000000000000000000000000000"
	config, err := config.NewConfig(bitcoin.MainNet, true, "test", "Tokenized Test", startHash, 8,
		2000, 10, 10, 1000, true)
	if err != nil {
		test.Errorf("Failed to create config : %v", err)
	}

	ctx := context.Background()
	store := storage.NewMockStorage()
	repo := NewBlockRepository(config, store)

	t := uint32(time.Now().Unix())
	var prevBlock bitcoin.Hash32
	for i := 0; i < testBlockCount; i++ {
		header := wire.BlockHeader{
			Version:   1,
			PrevBlock: prevBlock,
			Timestamp: time.Unix(int64(t), 0),
			Nonce:     uint32(randGen.Int()),
		}

		repo.Add(ctx, &header)

		hash := header.BlockHash()
		prevBlock = *hash
		blocks = append(blocks, hash)

		t += 600
	}

	if err := repo.Save(ctx); err != nil {
		test.Errorf("Failed to save repo : %v", err)
	}

	for _, revertHeight := range testRevertHeights {
		test.Logf("Test revert to (%d) : %s", revertHeight, blocks[revertHeight].String())

		if err := repo.Revert(ctx, revertHeight); err != nil {
			test.Errorf("Failed to revert repo : %v", err)
		}

		if !repo.LastHash().Equal(blocks[revertHeight]) {
			test.Errorf("Failed to revert repo to height %d", revertHeight)
		}

		if err := repo.Load(ctx); err != nil {
			test.Errorf("Failed to load repo after revert to %d : %v", revertHeight, err)
		}
	}
}
