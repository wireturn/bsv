package bmap

import (
	"github.com/bitcoinschema/go-aip"
	"github.com/bitcoinschema/go-b"
	"github.com/bitcoinschema/go-bap"
	"github.com/bitcoinschema/go-bob"
	magic "github.com/bitcoinschema/go-map"
)

// Tx is a Bmap formatted tx
type Tx struct {
	AIP *aip.Aip     `json:"AIP,omitempty" bson:"AIP,omitempty"`
	B   *b.B         `json:"B,omitempty" bson:"B,omitempty"`
	BAP *bap.Bap     `json:"BAP,omitempty" bson:"BAP,omitempty"`
	Blk bob.Blk      `json:"blk,omitempty" bson:"blk,omitempty"`
	In  []bob.Input  `json:"in,omitempty" bson:"in,omitempty"`
	MAP magic.MAP    `json:"MAP,omitempty" bson:"MAP,omitempty"`
	Out []bob.Output `json:"out,omitempty" bson:"out,omitempty"`
	Tx  bob.TxInfo   `json:"tx,omitempty" bson:"tx,omitempty"`
}

// NewFromBob returns a new BmapTx from a BobTx
func NewFromBob(bobTx *bob.Tx) (bmapTx *Tx, err error) {
	bmapTx = new(Tx)
	err = bmapTx.FromBob(bobTx)
	return
}

// FromBob returns a BmapTx from a BobTx
func (t *Tx) FromBob(bobTx *bob.Tx) (err error) {
	for _, out := range bobTx.Out {
		for index, tape := range out.Tape {
			if len(tape.Cell) > 0 {
				prefixData := tape.Cell[0].S
				switch prefixData {
				case aip.Prefix:
					t.AIP = aip.NewFromTape(tape)
					t.AIP.SetDataFromTapes(out.Tape)
				case bap.Prefix:
					if t.BAP, err = bap.NewFromTape(&out.Tape[index]); err != nil {
						return
					}
				case magic.Prefix:
					if t.MAP, err = magic.NewFromTape(&out.Tape[index]); err != nil {
						return
					}
				case b.Prefix:
					if t.B, err = b.NewFromTape(out.Tape[index]); err != nil {
						return
					}
				}
			}
		}

		// Set inherited fields
		t.Blk = bobTx.Blk
		t.In = bobTx.In
		t.Out = bobTx.Out
		t.Tx = bobTx.Tx
	}
	return
}
