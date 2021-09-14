package bootstrap

import (
	"context"

	"github.com/tokenized/pkg/bitcoin"
	"github.com/tokenized/pkg/storage"
	"github.com/tokenized/spynode/internal/platform/config"
	"github.com/tokenized/spynode/internal/spynode"
	"github.com/tokenized/spynode/pkg/client"
)

const (
	// SubSystem is used by the logger package
	SubSystem = "SpyNode"
)

func NewConfig(net bitcoin.Network, isTest bool, host, useragent, starthash string,
	untrustedNodes, safeDelay, shotgunCount, maxRetries, retryDelay int,
	requestMempool bool) (config.Config, error) {
	return config.NewConfig(net, isTest, host, useragent, starthash, untrustedNodes, safeDelay,
		shotgunCount, maxRetries, retryDelay, requestMempool)
}

type SpyNodeEmbedded interface {
	client.Client

	Run(context.Context) error
	Stop(context.Context) error

	LastHeight(ctx context.Context) int
	Hash(ctx context.Context, height int) (*bitcoin.Hash32, error)
	Time(ctx context.Context, height int) (uint32, error)
}

func NewNode(config config.Config, store storage.Storage, txFetcher spynode.TxFetcher,
	outputFetcher spynode.OutputFetcher) SpyNodeEmbedded {
	return spynode.NewNode(config, store, txFetcher, outputFetcher)
}
