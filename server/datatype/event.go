package datatype

import (
	"context"
	"github.com/TON-Market/tma/server/datatype/token"
	"github.com/google/uuid"
	"sync"
)

type EventDTO struct {
	ID         uuid.UUID `json:"id"`
	Collateral uint64    `json:"collateral"`
	LogoLink   string    `json:"logo-link"`
	Title      string    `json:"title"`
}

type Event struct {
	ID uuid.UUID
	sync.RWMutex
	collA    uint64
	collB    uint64
	LogoLink string
	Title    string
}

func NewEvent(logoLink, title string) *Event {
	return &Event{
		ID:       uuid.New(),
		RWMutex:  sync.RWMutex{},
		collA:    0,
		collB:    0,
		LogoLink: logoLink,
		Title:    title,
	}
}

func (e *Event) GetDeposit(t token.Token) uint64 {
	e.RLock()
	defer e.RUnlock()
	if t == token.A {
		return e.collA
	}
	return e.collB
}

func (e *Event) GetTotalDeposit() uint64 {
	e.RLock()
	defer e.RUnlock()
	return e.collA + e.collB
}

func (e *Event) AddDeposit(d uint64, t token.Token) {
	e.Lock()
	defer e.Unlock()
	if t == token.A {
		e.collA += d
	} else {
		e.collB += d
	}
}

type EventStorage interface {
	GetEventsPaginate(context.Context, int) []*Event
	GetEventByID(context.Context, uuid.UUID) (*Event, error)
	SaveEvent(context.Context, *Event) error
}

type EventMap map[string]*Event

func NewEventStorage() EventStorage {
	return &EventMap{}
}

func (s *EventMap) GetEventByID(ctx context.Context, u uuid.UUID) (*Event, error) {
	return nil, nil
}

func (s *EventMap) SaveEvent(ctx context.Context, e *Event) error {
	panic(nil)
}

func (s *EventMap) GetEventsPaginate(ctx context.Context, page int) []*Event {
	panic(nil)
}
