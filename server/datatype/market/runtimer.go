package market

import (
	"context"
	"errors"
	"fmt"
	"github.com/TON-Market/tma/server/datatype/token"
	"github.com/google/uuid"
	"github.com/tonkeeper/tongo/tlb"
	"sync"
)

type runtimer struct {
	sync.RWMutex
	eventRuntimeMap map[uuid.UUID]*eventRuntime
}

var (
	ErrRuntimeEventAlreadyExists = errors.New("err runtime event already exists")
	ErrRuntimeEventNotExist      = errors.New("err runtime event not exist")
)

func (r *runtimer) saveEvent(_ context.Context, e *Event) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.eventRuntimeMap[e.ID]; ok {
		return fmt.Errorf("save runtime event failed: %v", ErrRuntimeEventAlreadyExists)
	}

	er := &eventRuntime{
		isActive:      true,
		eventID:       e.ID,
		betRuntimeMap: make(map[token.Token]*betRuntime),
	}

	for t := range e.BetMap {
		er.betRuntimeMap[t] = &betRuntime{
			sync.RWMutex{},
			t,
			tlb.Grams(0),
		}
	}

	r.eventRuntimeMap[e.ID] = er

	return nil
}

func (r *runtimer) deposit(_ context.Context, d *Deal) error {
	r.RLock()
	defer r.RUnlock()

	if _, ok := r.eventRuntimeMap[d.EventID]; !ok {
		return fmt.Errorf("runtimer deposit failed: %v: id: %s", ErrRuntimeEventNotExist, d.EventID.String())
	}

	r.eventRuntimeMap[d.EventID].deposit(d.Token, d.Collateral)

	return nil
}

func (r *runtimer) snapshot(_ context.Context) map[uuid.UUID]*eventState {
	r.RLock()
	defer r.RUnlock()

	m := make(map[uuid.UUID]*eventState, len(r.eventRuntimeMap))

	for k, v := range r.eventRuntimeMap {
		m[k] = v.getState()
	}

	return m
}

func (r *runtimer) getEventState(_ context.Context, id uuid.UUID) *eventState {
	r.RLock()
	defer r.RUnlock()
	s := r.eventRuntimeMap[id].getState()
	return s
}

func (r *runtimer) close(_ context.Context, id uuid.UUID) error {
	r.RLock()
	defer r.RUnlock()

	er, ok := r.eventRuntimeMap[id]
	if !ok {
		return fmt.Errorf("runtimer close event: %s failed: %w", id.String(), ErrRuntimeEventNotExist)
	}

	er.close()
	return nil
}
