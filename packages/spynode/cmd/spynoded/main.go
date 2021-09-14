package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/tokenized/pkg/bitcoin"
	"github.com/tokenized/pkg/logger"
	"github.com/tokenized/pkg/rpcnode"
	"github.com/tokenized/pkg/storage"
	"github.com/tokenized/spynode/internal/platform/config"
	"github.com/tokenized/spynode/internal/spynode"
	"github.com/tokenized/spynode/pkg/client"

	"github.com/kelseyhightower/envconfig"
)

var (
	buildVersion = "unknown"
	buildDate    = "unknown"
	buildUser    = "unknown"
)

func main() {

	// -------------------------------------------------------------------------
	// Logging
	logConfig := logger.NewConfig(true, false, "")
	logConfig.EnableSubSystem(spynode.SubSystem)
	ctx := logger.ContextWithLogConfig(context.Background(), logConfig)

	// -------------------------------------------------------------------------
	// Config

	var cfg struct {
		Network string `default:"mainnet" envconfig:"BITCOIN_NETWORK"`
		IsTest  bool   `default:"true" envconfig:"IS_TEST"`
		Node    struct {
			Address        string `envconfig:"NODE_ADDRESS"`
			UserAgent      string `default:"/Tokenized:0.1.0/" envconfig:"NODE_USER_AGENT"`
			StartHash      string `envconfig:"START_HASH"`
			UntrustedNodes int    `default:"25" envconfig:"UNTRUSTED_NODES"`
			SafeTxDelay    int    `default:"2000" envconfig:"SAFE_TX_DELAY"`
			ShotgunCount   int    `default:"100" envconfig:"SHOTGUN_COUNT"`
			MaxRetries     int    `default:"25" envconfig:"MAX_RETRIES"`
			RetryDelay     int    `default:"2000" envconfig:"RETRY_DELAY"`
			RequestMempool bool   `default:"true" envconfig:"REQUEST_MEMPOOL" json:"REQUEST_MEMPOOL"`
		}
		NodeStorage struct {
			Region    string `default:"ap-southeast-2" envconfig:"NODE_STORAGE_REGION"`
			AccessKey string `envconfig:"NODE_STORAGE_ACCESS_KEY"`
			Secret    string `envconfig:"NODE_STORAGE_SECRET"`
			Bucket    string `default:"standalone" envconfig:"NODE_STORAGE_BUCKET"`
			Root      string `default:"./tmp" envconfig:"NODE_STORAGE_ROOT"`
		}
		RPC rpcnode.Config
	}

	if err := envconfig.Process("Node", &cfg); err != nil {
		logger.Info(ctx, "Parsing Config : %v", err)
	}

	logger.Info(ctx, "Started : Application Initializing")
	defer log.Println("Completed")

	cfgJSON, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		logger.Fatal(ctx, "Marshalling Config to JSON : %v", err)
	}

	logger.Info(ctx, "Build %v (%v on %v)\n", buildVersion, buildUser, buildDate)

	// TODO: Mask sensitive values
	logger.Info(ctx, "Config : %v\n", string(cfgJSON))

	// -------------------------------------------------------------------------
	// Storage
	storageConfig := storage.NewConfig(cfg.NodeStorage.Bucket, cfg.NodeStorage.Root)

	var store storage.Storage
	if strings.ToLower(storageConfig.Bucket) == "standalone" {
		store = storage.NewFilesystemStorage(storageConfig)
	} else {
		store = storage.NewS3Storage(storageConfig)
	}

	// -------------------------------------------------------------------------
	// Node Config
	nodeConfig, err := config.NewConfig(bitcoin.NetworkFromString(cfg.Network), cfg.IsTest,
		cfg.Node.Address, cfg.Node.UserAgent, cfg.Node.StartHash, cfg.Node.UntrustedNodes,
		cfg.Node.SafeTxDelay, cfg.Node.ShotgunCount, cfg.Node.MaxRetries, cfg.Node.RetryDelay,
		cfg.Node.RequestMempool)
	if err != nil {
		logger.Error(ctx, "Failed to create node config : %s\n", err)
		return
	}

	// -------------------------------------------------------------------------
	// RPC
	rpc, err := rpcnode.NewNode(&cfg.RPC)
	if err != nil {
		logger.Error(ctx, "Failed to create rpc node : %s\n", err)
		return
	}

	// -------------------------------------------------------------------------
	// Node

	node := spynode.NewNode(nodeConfig, store, rpc, rpc)

	logHandler := LogHandler{ctx: ctx}
	node.RegisterHandler(&logHandler)

	if err := node.SubscribeContracts(ctx); err != nil {
		logger.Error(ctx, "Failed to subscribe to contracts : %s", err)
		return
	}

	signals := make(chan os.Signal, 1)
	go func() {
		signal := <-signals
		logger.Info(ctx, "Received signal : %s", signal)
		if signal == os.Interrupt {
			logger.Info(ctx, "Stopping node")
			node.Stop(ctx)
		}
	}()

	// -------------------------------------------------------------------------
	// Start Node Service

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		logger.Info(ctx, "Node Running")
		serverErrors <- node.Run(ctx)
	}()

	// -------------------------------------------------------------------------
	// Shutdown

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		if err != nil {
			logger.Error(ctx, "Node failure : %s", err.Error())
		}

	case <-osSignals:
		logger.Info(ctx, "Start shutdown...")

		// Asking handler to shutdown and load shed.
		if err := node.Stop(ctx); err != nil {
			logger.Error(ctx, "Failed to stop spynode server: %s", err.Error())
		}
	}
}

type LogHandler struct {
	ctx   context.Context
	mutex sync.Mutex
}

// Spynode handler interface
func (handler *LogHandler) HandleHeaders(ctx context.Context, headers *client.Headers) {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()

	ctx = logger.ContextWithOutLogSubSystem(ctx)
	logger.Info(ctx, "New header (%d) : %s", headers.StartHeight, headers.Headers[0].BlockHash())
}

func (handler *LogHandler) HandleTx(ctx context.Context, tx *client.Tx) {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()

	ctx = logger.ContextWithOutLogSubSystem(ctx)
	logger.Info(ctx, "Tx : %s", tx.Tx.TxHash().String())
}

func (handler *LogHandler) HandleTxUpdate(ctx context.Context, update *client.TxUpdate) {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()

	ctx = logger.ContextWithOutLogSubSystem(ctx)
	logger.Info(ctx, "Tx update : %s", update.TxID)
}

func (handler *LogHandler) HandleInSync(ctx context.Context) {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()

	ctx = logger.ContextWithOutLogSubSystem(ctx)
	logger.Info(ctx, "In Sync")
}

func (handler *LogHandler) HandleMessage(ctx context.Context, payload client.MessagePayload) {}
