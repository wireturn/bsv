package mapi

import (
	"context"
	"time"

	"github.com/libsv/go-bt/v2"
	"github.com/pkg/errors"
	"github.com/tonicpow/go-minercraft"

	gopayd "github.com/libsv/payd"
	"github.com/libsv/payd/config"
)

type minercraftMapi struct {
	client *minercraft.Client
	cfg    *config.MApi
	svrCfg *config.Server
	fq     *bt.FeeQuote
}

// NewMapi will setup and return a new MAPI minercraftMapi data store.
func NewMapi(cfg *config.MApi, svrCfg *config.Server, client *minercraft.Client) *minercraftMapi {
	return &minercraftMapi{client: client, cfg: cfg, svrCfg: svrCfg, fq: bt.NewFeeQuote()}
}

// Broadcast will submit a transaction to mapi for inclusion in a block.
// Any errors will be returned, no error denotes success.
func (m *minercraftMapi) Send(ctx context.Context, args gopayd.SendTransactionArgs, req gopayd.CreatePayment) error {
	resp, err := m.client.SubmitTransaction(ctx,
		m.client.MinerByName(m.cfg.MinerName),
		&minercraft.Transaction{
			RawTx:              req.Transaction,
			CallBackURL:        "http://" + m.svrCfg.Hostname + "/api/v1/proofs/" + args.TxID,
			CallBackToken:      "",
			MerkleFormat:       "TSC",
			CallBackEncryption: "",
			MerkleProof:        true,
			DsCheck:            true,
		})
	if err != nil {
		return errors.Wrap(err, "failed to submit transaction to minerpool")
	}
	if resp.Results.ReturnResult == minercraft.QueryTransactionSuccess {
		return nil
	}
	return errors.Errorf("failed to submit transaction %s", resp.Results.ResultDescription)
}

// Status will return the current network status of a transaction.
func (m *minercraftMapi) Status(ctx context.Context, args gopayd.TxStatusArgs) (*gopayd.TxStatus, error) {
	resp, err := m.client.QueryTransaction(ctx, m.client.MinerByName(m.cfg.MinerName), args.TxID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query Tx from mAPI")
	}
	if !resp.Validated {
		return nil, errors.Wrap(err, "invalid message payload received from mAPI")
	}
	errNum := 0
	if resp.Query.ReturnResult == minercraft.QueryTransactionFailure {
		errNum = 1
	}
	status := resp.Query.ResultDescription
	if resp.Query.ReturnResult == minercraft.QueryTransactionSuccess {
		status = minercraft.QueryTransactionSuccess
	}
	return &gopayd.TxStatus{
		TxID:          args.TxID,
		Status:        status,
		BlockHash:     resp.Query.BlockHash,
		BlockHeight:   uint64(resp.Query.BlockHeight),
		Confirmations: uint64(resp.Query.Confirmations),
		Error:         errNum,
	}, nil
}

// Fees will return the current fees for the configured miner. If the fee has not
// expired we will return the current memoized fee quote.
func (m *minercraftMapi) Fees(ctx context.Context) (*bt.FeeQuote, error) {
	if !m.fq.Expired() {
		return m.fq, nil
	}
	fq, err := m.client.FeeQuote(ctx, m.client.MinerByName(m.cfg.MinerName))
	if err != nil {
		return nil, errors.Wrap(err, "failed to read fees for ")
	}
	if !fq.Validated {
		return m.fq, nil
	}
	stdfee := fq.Quote.GetFee(string(bt.FeeTypeStandard))
	datafee := fq.Quote.GetFee(string(bt.FeeTypeData))
	m.fq.AddQuote(bt.FeeTypeStandard, &bt.Fee{
		FeeType: bt.FeeTypeStandard,
		MiningFee: bt.FeeUnit{
			Satoshis: stdfee.MiningFee.Satoshis,
			Bytes:    stdfee.MiningFee.Bytes,
		},
		RelayFee: bt.FeeUnit{
			Satoshis: stdfee.RelayFee.Satoshis,
			Bytes:    stdfee.RelayFee.Bytes,
		},
	})
	m.fq.AddQuote(bt.FeeTypeData, &bt.Fee{
		FeeType: bt.FeeTypeData,
		MiningFee: bt.FeeUnit{
			Satoshis: datafee.MiningFee.Satoshis,
			Bytes:    datafee.MiningFee.Bytes,
		},
		RelayFee: bt.FeeUnit{
			Satoshis: datafee.RelayFee.Satoshis,
			Bytes:    datafee.RelayFee.Bytes,
		},
	})
	exp, err := time.Parse(time.RFC3339, fq.Quote.ExpirationTime)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse expiration time when getting fee quote")
	}
	m.fq.UpdateExpiry(exp.UTC())
	return m.fq, nil
}
