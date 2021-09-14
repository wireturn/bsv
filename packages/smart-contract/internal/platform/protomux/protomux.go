package protomux

import (
	"context"
	"errors"
	"sync"

	"github.com/tokenized/pkg/wire"
	"github.com/tokenized/smart-contract/pkg/inspector"
)

const (
	// SEE is used for broadcast txs
	SEE = "SEE"

	// SAW is used for replayed txs (reserved)
	SAW = "SAW"

	// LOST is used for reorgs
	LOST = "LOST"

	// STOLE is used for double spends
	STOLE = "STOLE"

	// END is used to call back finalization of a tx
	END = "END"

	// ANY_EVENT is used as the event, when the handler is for all events
	ANY_EVENT = "*"
)

// Handler is the interface for this Protocol Mux
type Handler interface {
	Respond(context.Context, wire.Message) error
	Reprocess(context.Context, *inspector.Transaction) error
	Trigger(context.Context, string, *inspector.Transaction) error
	SetResponder(ResponderFunc)
	SetReprocessor(ReprocessFunc)
}

// A Handler is a type that handles a protocol messages
type HandlerFunc func(ctx context.Context, itx *inspector.Transaction) error

// A ResponderFunc will handle responses
type ResponderFunc func(ctx context.Context, tx *wire.MsgTx) error

// A ReprocessFunc will handle responses
type ReprocessFunc func(ctx context.Context, itx *inspector.Transaction) error

type ProtoMux struct {
	Responder         ResponderFunc
	Reprocessor       ReprocessFunc
	SeeHandlers       map[string][]HandlerFunc
	LostHandlers      map[string][]HandlerFunc
	StoleHandlers     map[string][]HandlerFunc
	ReprocessHandlers map[string][]HandlerFunc

	SeeDefaultHandlers       []HandlerFunc
	LostDefaultHandlers      []HandlerFunc
	StoleDefaultHandlers     []HandlerFunc
	ReprocessDefaultHandlers []HandlerFunc

	lock, repLock sync.Mutex
}

func New() *ProtoMux {
	pm := &ProtoMux{
		SeeHandlers:       make(map[string][]HandlerFunc),
		LostHandlers:      make(map[string][]HandlerFunc),
		StoleHandlers:     make(map[string][]HandlerFunc),
		ReprocessHandlers: make(map[string][]HandlerFunc),
	}

	return pm
}

// Handle registers a new handler
func (p *ProtoMux) Handle(verb, event string, handler HandlerFunc) {
	switch verb {
	case SEE:
		p.SeeHandlers[event] = append(p.SeeHandlers[event], handler)
	case LOST:
		p.LostHandlers[event] = append(p.LostHandlers[event], handler)
	case STOLE:
		p.StoleHandlers[event] = append(p.StoleHandlers[event], handler)
	case END:
		p.ReprocessHandlers[event] = append(p.ReprocessHandlers[event], handler)
	default:
		panic("Unknown handler type")
	}
}

// Trigger fires a handler
func (p *ProtoMux) Trigger(ctx context.Context, verb string, itx *inspector.Transaction) error {

	if itx.MsgProto == nil {
		return errors.New("Not a protocol tx")
	}

	var group map[string][]HandlerFunc

	switch verb {
	case SEE:
		group = p.SeeHandlers
	case LOST:
		group = p.LostHandlers
	case STOLE:
		group = p.StoleHandlers
	case END:
		group = p.ReprocessHandlers
	default:
		return errors.New("Unknown handler type")
	}

	// Locate the handlers from the group
	event := itx.MsgProto.Code()
	handlers, exists := group[event]

	if !exists {
		// Fall back to ANY_EVENT
		handlers, exists = group[ANY_EVENT]
		if !exists {
			return nil
		}
	}

	// Notify the handlers
	for _, handler := range handlers {
		if err := handler(ctx, itx); err != nil {
			return err
		}
	}

	return nil
}

// SetResponder sets the function used for handling responses
func (p *ProtoMux) SetResponder(responder ResponderFunc) {
	p.lock.Lock()
	p.Responder = responder
	p.lock.Unlock()
}

// SetReprocessor sets the function used for reprocessing
func (p *ProtoMux) SetReprocessor(reprocessor ReprocessFunc) {
	p.repLock.Lock()
	p.Reprocessor = reprocessor
	p.repLock.Unlock()
}

func (p *ProtoMux) Reprocess(ctx context.Context, itx *inspector.Transaction) error {
	p.repLock.Lock()
	defer p.repLock.Unlock()
	return p.Reprocessor(ctx, itx)
}

// Respond handles a response via the responder
func (p *ProtoMux) Respond(ctx context.Context, m wire.Message) error {
	tx, ok := m.(*wire.MsgTx)
	if !ok {
		return errors.New("Could not assert as *wire.MsgTx")
	}

	p.lock.Lock()
	defer p.lock.Unlock()
	return p.Responder(ctx, tx)
}
