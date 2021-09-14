// Package bob is a library for working with BOB formatted transactions
//
// Specs: https://bob.planaria.network/
//
// If you have any suggestions or comments, please feel free to open an issue on
// this GitHub repository!
//
// By BitcoinSchema Organization (https://bitcoinschema.org)
package bob

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/libsv/go-bt"
	"github.com/libsv/go-bt/bscript"
)

// Protocol delimiter constants
// OP_SWAP = 0x7c = 124 = "|"
const (
	ProtocolDelimiterAsm  = "OP_SAWP"
	ProtocolDelimiterInt  = 0x7c
	ProtocolDelimiterByte = byte(ProtocolDelimiterInt)
	ProtocolDelimiter     = string(rune(ProtocolDelimiterInt))
)

// TxUnmarshal is a BOB formatted Bitcoin transaction that includes
// interfaces where types may change
//
// DO NOT CHANGE ORDER - aligned for memory optimization (malign)
//
type TxUnmarshal struct {
	In   []Input           `json:"in"`
	Out  []OutputUnmarshal `json:"out"`
	ID   string            `json:"_id"`
	Tx   TxInfo            `json:"tx"`
	Blk  Blk               `json:"blk"`
	Lock uint32            `json:"lock"`
}

// OutputUnmarshal is a transaction output
type OutputUnmarshal struct {
	I    uint8      `json:"i"`
	Tape []Tape     `json:"tape"`
	E    EUnmarshal `json:"e,omitempty"`
}

// EUnmarshal has address and value information
type EUnmarshal struct {
	A interface{} `json:"a,omitempty" bson:"a,omitempty"`
	V uint64      `json:"v,omitempty" bson:"v,omitempty"`
	I uint32      `json:"i" bson:"i"`
	H string      `json:"h,omitempty" bson:"h,omitempty"`
}

// E has address and value information
type E struct {
	A string `json:"a,omitempty" bson:"a,omitempty"`
	V uint64 `json:"v,omitempty" bson:"v,omitempty"`
	I uint32 `json:"i" bson:"i"`
	H string `json:"h,omitempty" bson:"h,omitempty"`
}

// Cell is a single OP_RETURN protocol
type Cell struct {
	H   string `json:"h,omitempty" bson:"h,omitempty"`
	B   string `json:"b,omitempty" bson:"b,omitempty"`
	LB  string `json:"lb,omitempty" bson:"lb,omitempty"`
	S   string `json:"s,omitempty" bson:"s,omitempty"`
	LS  string `json:"ls,omitempty" bson:"ls,omitempty"`
	I   uint8  `json:"i" bson:"i"`
	II  uint8  `json:"ii" bson:"ii"`
	Op  uint16 `json:"op,omitempty" bson:"op,omitempty"`
	Ops string `json:"ops,omitempty" bson:"ops,omitempty"`
}

// Input is a transaction input
//
// DO NOT CHANGE ORDER - aligned for memory optimization (malign)
//
type Input struct {
	E    E      `json:"e" bson:"e"`
	Tape []Tape `json:"tape" bson:"tape"`
	Seq  uint32 `json:"seq" bson:"seq"`
	I    uint8  `json:"i" bson:"i"`
}

// Tape is a tape
type Tape struct {
	Cell []Cell `json:"cell"`
	I    uint8  `json:"i"`
}

// Output is a transaction output
type Output struct {
	I    uint8  `json:"i"`
	Tape []Tape `json:"tape"`
	E    E      `json:"e,omitempty"`
}

// Blk contains the block info
type Blk struct {
	I uint32 `json:"i"`
	T uint32 `json:"t"`
}

// TxInfo contains the transaction info
type TxInfo struct {
	H string `json:"h"`
}

// Tx is a BOB formatted Bitcoin transaction
//
// DO NOT CHANGE ORDER - aligned for memory optimization (malign)
//
type Tx struct {
	In   []Input  `json:"in"`
	Out  []Output `json:"out"`
	ID   string   `json:"_id"`
	Tx   TxInfo   `json:"tx"`
	Blk  Blk      `json:"blk"`
	Lock uint32   `json:"lock"`
}

// NewFromBytes creates a new BOB Tx from a NDJSON line representing a BOB transaction,
// as returned by the bitbus 2 API
func NewFromBytes(line []byte) (bobTx *Tx, err error) {
	bobTx = new(Tx)
	err = bobTx.FromBytes(line)
	return
}

// NewFromRawTxString creates a new BobTx from a hex encoded raw tx string
func NewFromRawTxString(rawTxString string) (bobTx *Tx, err error) {
	bobTx = new(Tx)
	err = bobTx.FromRawTxString(rawTxString)
	return
}

// NewFromString creates a new BobTx from a BOB formatted string
func NewFromString(line string) (bobTx *Tx, err error) {
	bobTx = new(Tx)
	err = bobTx.FromString(line)
	return
}

// NewFromTx creates a new BobTx from a libsv Transaction
func NewFromTx(tx *bt.Tx) (bobTx *Tx, err error) {
	bobTx = new(Tx)
	err = bobTx.FromTx(tx)
	return
}

// FromBytes takes a BOB formatted tx string as bytes
func (t *Tx) FromBytes(line []byte) error {
	tu := new(TxUnmarshal)
	if err := json.Unmarshal(line, &tu); err != nil {
		return fmt.Errorf("error parsing line: %v, %w", line, err)
	}

	// The out.E.A field can be either a boolean or a string
	// So we need to unmarshal into an interface, and fix the normal struct the user
	// of this lib will work with (so they don't have to format the interface themselves)
	fixedOuts := make([]Output, 0)
	for _, out := range tu.Out {
		fixedOuts = append(fixedOuts, Output{
			I:    out.I,
			Tape: out.Tape,
			E: E{
				A: fmt.Sprintf("%s", out.E.A), // todo: test this with (string) and (bool)
				V: out.E.V,
				I: out.E.I,
				H: out.E.H,
			},
		})
	}
	t.Blk = tu.Blk
	t.ID = tu.ID
	t.In = tu.In
	t.Lock = tu.Lock
	t.Out = fixedOuts
	t.Tx = tu.Tx

	// Check for missing hex values and supply them
	for outIdx, out := range t.Out {
		for tapeIdx, tape := range out.Tape {
			for cellIdx, cell := range tape.Cell {
				if len(cell.H) == 0 && len(cell.B) > 0 {
					// base 64 decode cell.B and encode it to hex string
					cellBytes, err := base64.StdEncoding.DecodeString(cell.B)
					if err != nil {
						return err
					}
					t.Out[outIdx].Tape[tapeIdx].Cell[cellIdx].H = hex.EncodeToString(cellBytes)
				}
			}
		}
	}
	for inIdx, in := range t.In {
		for tapeIdx, tape := range in.Tape {
			for cellIdx, cell := range tape.Cell {
				if len(cell.H) == 0 && len(cell.B) > 0 {
					// base 64 decode cell.B and encode it to hex string
					cellBytes, err := base64.StdEncoding.DecodeString(cell.B)
					if err != nil {
						return err
					}
					t.In[inIdx].Tape[tapeIdx].Cell[cellIdx].H = hex.EncodeToString(cellBytes)
				}
			}
		}
	}

	return nil
}

// FromRawTxString takes a hex encoded tx string
func (t *Tx) FromRawTxString(rawTxString string) error {
	tx, err := bt.NewTxFromString(rawTxString)
	if err != nil {
		return err
	}
	return t.FromTx(tx)
}

// FromString takes a BOB formatted string
func (t *Tx) FromString(line string) (err error) {
	err = t.FromBytes([]byte(line))
	return
}

// FromTx takes a bt.Tx
func (t *Tx) FromTx(tx *bt.Tx) error {

	// Set the transaction ID
	t.Tx.H = tx.GetTxID()

	// Set the inputs
	for inIdx, i := range tx.Inputs {

		bobInput := Input{
			I: uint8(inIdx),
			Tape: []Tape{{
				Cell: []Cell{{
					H: hex.EncodeToString(i.ToBytes(false)),
					B: base64.RawStdEncoding.EncodeToString(i.ToBytes(false)),
					S: i.String(),
				}},
				I: 0,
			}},
			E: E{
				H: i.PreviousTxID,
			},
		}

		t.In = append(t.In, bobInput)
	}

	// Process outputs
	for idxOut, o := range tx.Outputs {
		var adr string

		// Try to get a pub_key hash (ignore fail when this is not a locking script)
		outPubKeyHash, _ := o.LockingScript.GetPublicKeyHash()
		if len(outPubKeyHash) > 0 {
			outAddress, err := bscript.NewAddressFromPublicKeyHash(outPubKeyHash, true)
			if err != nil {
				return fmt.Errorf("failed to get address from pubkeyhash %x: %w", outPubKeyHash, err)
			}
			adr = outAddress.AddressString
		}

		// Initialize out tapes and locking script asm
		asm, err := o.LockingScript.ToASM()
		if err != nil {
			return err
		}
		pushDatas := strings.Split(asm, " ")

		var outTapes []Tape
		bobOutput := Output{
			I:    uint8(idxOut),
			Tape: outTapes,
			E: E{
				A: adr,
			},
		}

		var currentTape Tape
		if len(pushDatas) > 0 {

			for pdIdx, pushData := range pushDatas {

				// Ignore error if it fails, use empty
				pushDataBytes, _ := hex.DecodeString(pushData)
				b64String := base64.StdEncoding.EncodeToString(pushDataBytes)

				if pushData != ProtocolDelimiterAsm {
					currentTape.Cell = append(currentTape.Cell, Cell{
						B:  b64String,
						H:  pushData,
						S:  string(pushDataBytes),
						I:  uint8(idxOut),
						II: uint8(pdIdx),
					})
				}
				// Note: OP_SWAP is 0x7c which is also ascii "|" which is our protocol separator.
				// This is not used as OP_SWAP at all since this is in the script after the OP_FALSE
				if "OP_RETURN" == pushData || ProtocolDelimiterAsm == pushData {
					outTapes = append(outTapes, currentTape)
					currentTape = Tape{}
				}
			}
		}

		// Add the trailing tape
		outTapes = append(outTapes, currentTape)
		bobOutput.Tape = outTapes

		t.Out = append(t.Out, bobOutput)
	}

	return nil
}

// ToRawTxString converts the BOBTx to a libsv.transaction, and outputs the raw hex
func (t *Tx) ToRawTxString() (string, error) {
	tx, err := t.ToTx()
	if err != nil {
		return "", err
	}
	return tx.ToString(), nil
}

// ToString returns a json string of bobTx
func (t *Tx) ToString() (string, error) {
	// Create JSON from the instance data.
	b, err := json.Marshal(t)
	return string(b), err

}

// ToTx returns a bt.Tx
func (t *Tx) ToTx() (*bt.Tx, error) {
	tx := bt.NewTx()

	tx.LockTime = t.Lock

	for _, in := range t.In {

		if len(in.Tape) == 0 || len(in.Tape[0].Cell) == 0 {
			return nil, fmt.Errorf("failed to process inputs. More tapes or cells than expected. %+v", in.Tape)
		}

		prevTxScript, _ := bscript.NewP2PKHFromAddress(in.E.A)

		var scriptAsm []string
		for _, cell := range in.Tape[0].Cell {
			cellData := cell.H
			scriptAsm = append(scriptAsm, cellData)
		}

		builtUnlockScript, err := bscript.NewFromASM(strings.Join(scriptAsm, " "))
		if err != nil {
			return nil, fmt.Errorf("failed to get script from asm: %v error: %w", scriptAsm, err)
		}

		// add inputs
		i := &bt.Input{
			PreviousTxID:       in.E.H,
			PreviousTxOutIndex: in.E.I,
			PreviousTxSatoshis: in.E.V,
			PreviousTxScript:   prevTxScript,
			UnlockingScript:    builtUnlockScript,
			SequenceNumber:     in.Seq,
		}

		tx.AddInput(i)
	}

	// add outputs
	for _, out := range t.Out {
		// Build the locking script
		var lockScriptAsm []string
		for tapeIdx, tape := range out.Tape {
			for cellIdx, cell := range tape.Cell {
				if cellIdx == 0 && tapeIdx > 1 {
					// add the separator back in
					lockScriptAsm = append(lockScriptAsm, ProtocolDelimiterAsm)
				}

				if len(cell.H) > 0 {
					lockScriptAsm = append(lockScriptAsm, cell.H)
				} else if len(cell.Ops) > 0 {
					lockScriptAsm = append(lockScriptAsm, cell.Ops)
				}
			}
		}

		lockingScript, _ := bscript.NewFromASM(strings.Join(lockScriptAsm, " "))
		o := &bt.Output{
			Satoshis:      out.E.V,
			LockingScript: lockingScript,
		}

		tx.AddOutput(o)
	}

	return tx, nil
}
