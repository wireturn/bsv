package node

import (
	"context"

	"github.com/tokenized/pkg/bitcoin"
	"github.com/tokenized/pkg/txbuilder"
	"github.com/tokenized/pkg/wire"
	"github.com/tokenized/smart-contract/internal/platform/db"
	"github.com/tokenized/smart-contract/internal/platform/protomux"
)

type ResponseWriter struct {
	Inputs        []bitcoin.UTXO
	Outputs       []Output
	RejectInputs  []bitcoin.UTXO
	RejectOutputs []Output
	RejectAddress bitcoin.RawAddress
	Config        *Config
	MasterDB      *db.DB
	Mux           protomux.Handler
}

// AddChangeOutput is a helper to add a change output
func (w *ResponseWriter) AddChangeOutput(ctx context.Context, addr bitcoin.RawAddress) error {
	out := outputValue(ctx, w.Config, addr, 0, true)
	w.Outputs = append(w.Outputs, *out)
	return nil
}

// AddOutput is a helper to add a payment output
func (w *ResponseWriter) AddOutput(ctx context.Context, addr bitcoin.RawAddress, value uint64) error {
	out := outputValue(ctx, w.Config, addr, value, false)
	w.Outputs = append(w.Outputs, *out)
	return nil
}

// AddFee attaches the fee as the next output, if configured
func (w *ResponseWriter) AddContractFee(ctx context.Context, value uint64) error {
	if fee := outputFee(ctx, w.Config, value); fee != nil {
		changeAlreadySet := false
		for _, output := range w.Outputs {
			if output.Change {
				changeAlreadySet = true
				break
			}
		}
		if !changeAlreadySet {
			fee.Change = true
		}
		w.Outputs = append(w.Outputs, *fee)
	}
	return nil
}

// SetUTXOs is an optional function that allows explicit UTXOs to be spent in the response
// be sure to remember any remaining UTXOs so they can be spent later.
func (w *ResponseWriter) SetUTXOs(ctx context.Context, utxos []bitcoin.UTXO) error {
	w.Inputs = utxos
	return nil
}

// SetRejectUTXOs is an optional function that allows explicit UTXOs to be spent in the reject
// response be sure to remember any remaining UTXOs so they can be spent later.
func (w *ResponseWriter) SetRejectUTXOs(ctx context.Context, utxos []bitcoin.UTXO) error {
	w.RejectInputs = utxos
	return nil
}

// AddRejectValue is a helper to add a refund amount to an address. This is used for multi-party
//   value transfers when different users input different amounts of bitcoin and need that refunded
//   if the request is rejected.
func (w *ResponseWriter) AddRejectValue(ctx context.Context, addr bitcoin.RawAddress, value uint64) error {
	// Look for existing output to this address.
	for i, output := range w.RejectOutputs {
		if output.Address.Equal(addr) {
			w.RejectOutputs[i].Value += value
			return nil
		}
	}

	// Add a new output for this address. If it is the first output make it the change output.
	w.RejectOutputs = append(w.RejectOutputs, Output{
		Address: addr,
		Value:   value,
		Change:  len(w.RejectOutputs) == 0,
	})
	return nil
}

// ClearRejectOutputValues zeroizes the values of the reject outputs so they become only
//   notification outputs.
func (w *ResponseWriter) ClearRejectOutputValues(changeAddr bitcoin.RawAddress) {
	for i, _ := range w.RejectOutputs {
		w.RejectOutputs[i].Change = false
		w.RejectOutputs[i].Value = 0
	}
}

// Respond sends the prepared response to the protocol mux
func (w *ResponseWriter) Respond(ctx context.Context, tx *wire.MsgTx) error {
	return w.Mux.Respond(ctx, tx)
}

// Output is an output address for a response
type Output struct {
	Address bitcoin.RawAddress
	Value   uint64
	Change  bool
}

// outputValue returns a payment output ensuring the value is always above the dust limit
func outputValue(ctx context.Context, config *Config, addr bitcoin.RawAddress, value uint64, change bool) *Output {
	dustLimit, err := txbuilder.DustLimitForAddress(addr, config.DustFeeRate)
	if err != nil {
		dustLimit = txbuilder.DustLimit(txbuilder.P2PKHOutputSize, config.DustFeeRate)
	}
	if value < dustLimit {
		value = dustLimit
	}

	out := &Output{
		Address: addr,
		Value:   value,
		Change:  change,
	}

	return out
}

// outputFee prepares a special fee output based on node configuration
func outputFee(ctx context.Context, config *Config, value uint64) *Output {
	if value > 0 {
		return &Output{
			Address: config.FeeAddress,
			Value:   value,
		}
	}

	return nil
}
