package client

import (
	"io"
	"math"

	"github.com/tokenized/pkg/bitcoin"
	"github.com/tokenized/pkg/wire"
)

const (
	// UnconfirmedBlockHeight is the block height specified when the transaction is unconfirmed.
	UnconfirmedBlockHeight = math.MaxUint32

	// MessageTypeRegister is the type of a register message.
	MessageTypeRegister = 1

	// MessageTypeSubscribePushData is the type of a subscribe push data message.
	MessageTypeSubscribePushData = 11

	// MessageTypeUnsubscribePushData is the type of an unsubscribe push data message.
	MessageTypeUnsubscribePushData = 12

	// MessageTypeSubscribeTx is the type of a subscribe tx message.
	MessageTypeSubscribeTx = 13

	// MessageTypeUnsubscribeTx is the type of an unsubscribe tx message.
	MessageTypeUnsubscribeTx = 14

	// MessageTypeSubscribeOutputs is the type of a subscribe tx message.
	MessageTypeSubscribeOutputs = 15

	// MessageTypeUnsubscribeOutputs is the type of an unsubscribe tx message.
	MessageTypeUnsubscribeOutputs = 16

	// MessageTypeSubscribeHeaders is the type of a subscribe headers message.
	MessageTypeSubscribeHeaders = 17

	// MessageTypeUnsubscribeHeaders is the type of an unsubscribe headers message.
	MessageTypeUnsubscribeHeaders = 18

	// MessageTypeSubscribeContracts is the type of a subscribe contracts message.
	MessageTypeSubscribeContracts = 19

	// MessageTypeUnsubscribeContracts is the type of an unsubscribe contracts message.
	MessageTypeUnsubscribeContracts = 20

	// MessageTypeReady is the type of a ready message.
	MessageTypeReady = 30

	// MessageTypeGetChainTip requests chain tip info.
	MessageTypeGetChainTip = 41

	// MessageTypeGetHeaders requests headers.
	MessageTypeGetHeaders = 42

	// MessageTypeSendTx sends a tx to the Bitcoin network.
	MessageTypeSendTx = 43

	// MessageTypeGetTx requests a transaction.
	MessageTypeGetTx = 44

	// MessageTypeAcceptRegister is the type of an accept register message.
	MessageTypeAcceptRegister = 101

	// MessageTypeBaseTx is the type of a base tx message.
	MessageTypeBaseTx = 110

	// MessageTypeTx is the type of a tx message.
	MessageTypeTx = 111

	// MessageTypeTxUpdate is the type of a tx update message.
	MessageTypeTxUpdate = 112

	// MessageTypeInSync is in sync info.
	MessageTypeInSync = 121

	// MessageTypeChainTip is chain tip info.
	MessageTypeChainTip = 122

	// MessageTypeHeaders is headers.
	MessageTypeHeaders = 123

	// MessageTypeAccept is an accept of the previous request.
	MessageTypeAccept = 200

	// MessageTypeReject is a rejection of the previous request.
	MessageTypeReject = 201

	// MessageTypePing is a ping message to keep the connection alive.
	MessageTypePing = 301

	// ConnectionTypeFull is the normal connection type the allows control and receiving data
	// messages.
	ConnectionTypeFull = uint8(1)

	// ConnectionTypeControl is a control only connection type that does not receive data messages.
	ConnectionTypeControl = uint8(2)
)

type Message struct {
	Payload MessagePayload
}

type MessagePayload interface {
	// Deserialize reads the message from a reader.
	Deserialize(io.Reader) error

	// Serialize writes the message to a writer.
	Serialize(io.Writer) error

	// Type returns they type of the message.
	Type() uint64
}

// Client to Server Messages -----------------------------------------------------------------------

// Register is the first message received from the client. It can be from a previous connection or
// it can be a new client based on whether the Key has been seen before.
type Register struct {
	Version          uint8             // Version of communication protocol
	Key              bitcoin.PublicKey // Client's public key
	Hash             bitcoin.Hash32    // For deriving ephemeral keys for use during this connection.
	StartBlockHeight uint32            // For new clients this is the starting height for data.
	ChainTip         bitcoin.Hash32    // The client's current chain tip block hash.
	ConnectionType   uint8             // The type of the connection.
	Signature        bitcoin.Signature // Signature of this messaage to prove key ownership.
}

// SubscribePushData adds new push data hashes used to filter relevant transactions. These and UTXOs for
// relevant transaction outputs are automatically retained between sessions.
type SubscribePushData struct {
	PushDatas [][]byte
}

// UnsubscribePushData removes push data hashes used to filter relevant transactions.
type UnsubscribePushData struct {
	PushDatas [][]byte
}

// SubscribeTx adds a new txid used to filter relevant transactions. Indexes are the indexes of the
// outputs that need to be monitored for spending.
type SubscribeTx struct {
	TxID    bitcoin.Hash32
	Indexes []uint32
}

// UnsubscribeTx removes the txid used to filter relevant transactions.
type UnsubscribeTx struct {
	TxID    bitcoin.Hash32
	Indexes []uint32
}

// SubscribeOutputs adds outputs that need to be monitored for spending.
type SubscribeOutputs struct {
	Outputs []*wire.OutPoint
}

// UnsubscribeOutputs removes the outputs used to filter relevant transactions.
type UnsubscribeOutputs struct {
	Outputs []*wire.OutPoint
}

// SubscribeHeaders requests that all new headers be sent as they are known.
type SubscribeHeaders struct{}

// UnsubscribeHeaders requests that headers no longer be sent automatically. They can still be
// directly requested.
type UnsubscribeHeaders struct{}

// SubscribeContracts requests that all contract-wide transactions be sent.
// ContractFormations/AssetCreations
type SubscribeContracts struct{}

// UnsubscribeContracts requests that all contract-wide transactions no longer be sent.
// ContractFormations/AssetCreations
type UnsubscribeContracts struct{}

// Ready tells the server that it can start syncing the client. This is sent after all initial
// Subscribe/Unsubscribe messages.
type Ready struct {
	NextMessageID uint64
}

// GetChainTip requests the height and hash of the lastest block.
type GetChainTip struct{}

// GetHeaders requests a set of headers.
type GetHeaders struct {
	RequestHeight int32  // -1 for most recent
	MaxCount      uint32 // max number of headers to return
}

// SendTx requests that tx be broadcast to the Bitcoin network. Indexes are the indexes of the
// outputs that need to be monitored for spending.
type SendTx struct {
	Tx      *wire.MsgTx
	Indexes []uint32
}

// GetTx requests a tx by its hash.
type GetTx struct {
	TxID bitcoin.Hash32
}

// Server to Client Messages -----------------------------------------------------------------------

type AcceptRegister struct {
	Key           bitcoin.PublicKey // Server's public key
	PushDataCount uint64
	UTXOCount     uint64
	MessageCount  uint64
	Signature     bitcoin.Signature // Signature of public key and the hash to prove key ownership.
}

// BaseTx is a just the basic transaction.
type BaseTx struct {
	Tx *wire.MsgTx
}

// Tx is a new transaction that is relevant to the client.
type Tx struct {
	ID      uint64 // message id to uniquely identify this message and the order of messages.
	Tx      *wire.MsgTx
	Outputs []*wire.TxOut // outputs being spent by inputs in Tx
	State   TxState       // initial state
}

// TxUpdate is an updated state for a transaction.
type TxUpdate struct {
	ID    uint64 // message id to uniquely identify this message and the order of messages.
	TxID  bitcoin.Hash32
	State TxState
}

// Headers is a list of block headers.
type Headers struct {
	RequestHeight int32  // height of request. zero if not a response to a request.
	StartHeight   uint32 // height of the first header, other headers are consecutive.
	Headers       []*wire.BlockHeader
}

// InSync is a notification that the messages are "up to date" with the network.
type InSync struct{}

// ChainTip is information about the latest block in the most POW chain.
type ChainTip struct {
	Height uint32
	Hash   bitcoin.Hash32
}

// Accept is an accept of the previous request
type Accept struct {
	MessageType uint8           // type of the message being rejected
	Hash        *bitcoin.Hash32 // optional identifier for the rejected item (tx)
}

// Reject is a rejection of the previous request
type Reject struct {
	MessageType uint8           // type of the message being rejected
	Hash        *bitcoin.Hash32 // optional identifier for the rejected item (tx)
	Code        uint32          // code representing the reason for the reject
	Message     string
}

// Ping is a ping to keep the connection live.
type Ping struct {
	TimeStamp uint64 // Current time
}

// Sub structures ----------------------------------------------------------------------------------

// TxState is state of a transaction.
type TxState struct {
	Safe             bool         // initial acceptance after checking for double spends
	UnSafe           bool         // transaction has known double spends or other unsafe attributes
	Cancelled        bool         // transaction has had a conflicting transaction confirmed
	UnconfirmedDepth uint32       // mempool chain depth
	MerkleProof      *MerkleProof // proof the txid is in the block.
}

// MerkleProof is the proof a txid is in the tree referenced by the merkle root of a block header.
type MerkleProof struct {
	Index             uint64 // Index of tx in block
	Path              []bitcoin.Hash32
	BlockHeader       wire.BlockHeader
	DuplicatedIndexes []uint64
}
