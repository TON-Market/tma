package event

import (
	"context"
	"errors"
	"fmt"
	"github.com/TON-Market/tma/server/datatype/token"
	"github.com/google/uuid"
	"sort"
	"sync"
	"time"
)

type EventTag int

const (
	Politic EventTag = iota
	Economics
	Crypto
	Culture
	Other
	No
)

// EventDTO - snapshot of Event
type EventDTO struct {
	ID         uuid.UUID `json:"id"`
	Tag        EventTag  `json:"tag"`
	Collateral uint64    `json:"collateral"`
	LogoLink   string    `json:"logo_link"`
	Title      string    `json:"title"`
}

// Event Concurrent object
type Event struct {
	ID uuid.UUID
	sync.RWMutex
	Tag      EventTag
	collA    uint64
	collB    uint64
	LogoLink string
	Title    string
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

// snapshot sync of storage every 5 seconds
type snapshot struct {
	mu   sync.RWMutex
	list []*EventDTO
}

type storage struct {
	mu sync.RWMutex
	m  map[uuid.UUID]*Event
}

type EventKeeper struct {
	snapshot
	storage
}

func NewEventKeeper() *EventKeeper {
	return &EventKeeper{
		snapshot{
			mu:   sync.RWMutex{},
			list: make([]*EventDTO, 0),
		},
		storage{
			mu: sync.RWMutex{},
			m:  make(map[uuid.UUID]*Event),
		},
	}
}

var (
	once      sync.Once
	singleton *EventKeeper
)

func Keeper() *EventKeeper {
	once.Do(func() {
		singleton = NewEventKeeper()
	})
	return singleton
}

func (k *EventKeeper) Start(_ context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	k.storage.mu.RLock()
	events := make([]*Event, 0, len(k.storage.m))
	for _, v := range k.storage.m {
		events = append(events, v)
	}
	k.storage.mu.RUnlock()

	l := make([]*EventDTO, 0, len(events))
	for _, v := range events {
		v.RLock()
		l = append(l, &EventDTO{
			ID:         v.ID,
			Tag:        v.Tag,
			Collateral: v.GetTotalDeposit(),
			LogoLink:   v.LogoLink,
			Title:      v.Title,
		})
		v.RUnlock()
	}

	sort.Slice(l, func(i, j int) bool {
		return l[i].Collateral > l[j].Collateral
	})

	k.snapshot.mu.Lock()
	defer k.snapshot.mu.Unlock()

	k.snapshot.list = l
}

var ErrNoSuchEventDtoInSnapshot = errors.New("err no such event dto in snapshot")

func (k *EventKeeper) GetByID(_ context.Context, id uuid.UUID) (*EventDTO, error) {
	k.snapshot.mu.RLock()
	defer k.snapshot.mu.RUnlock()
	for _, e := range k.snapshot.list {
		if e.ID == id {
			return e, nil
		}
	}
	return nil, fmt.Errorf("%s, id: %s",
		ErrNoSuchEventDtoInSnapshot.Error(), id.String())
}

func (k *EventKeeper) GetSnapshot(_ context.Context, tag EventTag, page int) ([]*EventDTO, error) {
	k.snapshot.mu.RLock()
	defer k.snapshot.mu.RUnlock()

	l := make([]*EventDTO, 0)

	if tag == No {
		l = k.list
	} else {
		for _, e := range k.list {
			if e.Tag != tag {
				continue
			}
			l = append(l, e)
		}
	}

	pageSize := 10
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= len(l) {
		return []*EventDTO{}, nil
	}

	if end > len(l) {
		end = len(l)
	}

	return l[start:end], nil
}
