package event

import (
	"context"
	"errors"
	"fmt"
	"github.com/TON-Market/tma/server/datatype/token"
	"github.com/google/uuid"
	"math"
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

type Bet struct {
	token.Token `json:"token"`
	Collateral  float64 `json:"collateral"`
	Title       string  `json:"title"`
	Percentage  int     `json:"percentage"`
}

// EventDTO - snapshot of Event
type EventDTO struct {
	ID         uuid.UUID `json:"id"`
	Tag        EventTag  `json:"tag"`
	Collateral float64   `json:"collateral"`
	LogoLink   string    `json:"logoLink"`
	Title      string    `json:"title"`
	Bets       []*Bet    `json:"bets"`
}

// Event Concurrent object
type Event struct {
	ID uuid.UUID
	sync.RWMutex
	Tag             EventTag
	LogoLink        string
	Title           string
	TotalCollateral float64
	Bets            map[token.Token]*Bet
}

func (e *Event) GetBets() []*Bet {
	e.RLock()
	defer e.RUnlock()

	betSlice := make([]*Bet, 0, 2)
	for _, b := range e.Bets {
		betSlice = append(betSlice, &Bet{
			Token:      b.Token,
			Collateral: b.Collateral,
			Title:      b.Title,
			Percentage: b.Percentage,
		})
	}

	return betSlice
}

func (e *Event) AddDeposit(d float64, t token.Token) {
	e.Lock()
	defer e.Unlock()

	b := e.Bets[t]
	bComplement := e.Bets[t.Complement()]
	b.Collateral += d

	total := b.Collateral + bComplement.Collateral

	e.TotalCollateral = total

	b.Percentage = int(b.Collateral / total)
	bComplement.Percentage = 1 - b.Percentage
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
	*snapshot
	*storage
}

func newEventKeeper() *EventKeeper {
	e := &Event{
		ID:              uuid.New(),
		RWMutex:         sync.RWMutex{},
		Tag:             Politic,
		LogoLink:        "https://tfj34t-88-201-232-88.ru.tuna.am/logo/elections2024.jpg",
		Title:           "2024 USA elections",
		TotalCollateral: 0,
		Bets: map[token.Token]*Bet{
			token.A: {
				Token:      token.A,
				Collateral: 0,
				Title:      "Yes",
				Percentage: 0,
			},

			token.B: {
				Token:      token.B,
				Collateral: 0,
				Title:      "No",
				Percentage: 0,
			},
		},
	}
	e1 := &Event{
		ID:              uuid.New(),
		RWMutex:         sync.RWMutex{},
		Tag:             Politic,
		LogoLink:        "https://tfj34t-88-201-232-88.ru.tuna.am/logo/elections2024.jpg",
		Title:           "2024 USA elections",
		TotalCollateral: 0,
		Bets: map[token.Token]*Bet{
			token.A: {
				Token:      token.A,
				Collateral: 0,
				Title:      "Yes",
				Percentage: 0,
			},

			token.B: {
				Token:      token.B,
				Collateral: 0,
				Title:      "No",
				Percentage: 0,
			},
		},
	}
	e2 := &Event{
		ID:              uuid.New(),
		RWMutex:         sync.RWMutex{},
		Tag:             Politic,
		LogoLink:        "https://tfj34t-88-201-232-88.ru.tuna.am/logo/elections2024.jpg",
		Title:           "2024 USA elections",
		TotalCollateral: 0,
		Bets: map[token.Token]*Bet{
			token.A: {
				Token:      token.A,
				Collateral: 0,
				Title:      "Yes",
				Percentage: 0,
			},

			token.B: {
				Token:      token.B,
				Collateral: 0,
				Title:      "No",
				Percentage: 0,
			},
		},
	}
	e3 := &Event{
		ID:              uuid.New(),
		RWMutex:         sync.RWMutex{},
		Tag:             Politic,
		LogoLink:        "https://tfj34t-88-201-232-88.ru.tuna.am/logo/elections2024.jpg",
		Title:           "2024 USA elections",
		TotalCollateral: 0,
		Bets: map[token.Token]*Bet{
			token.A: {
				Token:      token.A,
				Collateral: 0,
				Title:      "Yes",
				Percentage: 0,
			},

			token.B: {
				Token:      token.B,
				Collateral: 0,
				Title:      "No",
				Percentage: 0,
			},
		},
	}
	e4 := &Event{
		ID:              uuid.New(),
		RWMutex:         sync.RWMutex{},
		Tag:             Politic,
		LogoLink:        "https://tfj34t-88-201-232-88.ru.tuna.am/logo/elections2024.jpg",
		Title:           "2024 USA elections",
		TotalCollateral: 0,
		Bets: map[token.Token]*Bet{
			token.A: {
				Token:      token.A,
				Collateral: 0,
				Title:      "Yes",
				Percentage: 0,
			},

			token.B: {
				Token:      token.B,
				Collateral: 0,
				Title:      "No",
				Percentage: 0,
			},
		},
	}
	e5 := &Event{
		ID:              uuid.New(),
		RWMutex:         sync.RWMutex{},
		Tag:             Politic,
		LogoLink:        "https://tfj34t-88-201-232-88.ru.tuna.am/logo/elections2024.jpg",
		Title:           "2024 USA elections",
		TotalCollateral: 0,
		Bets: map[token.Token]*Bet{
			token.A: {
				Token:      token.A,
				Collateral: 0,
				Title:      "Yes",
				Percentage: 0,
			},

			token.B: {
				Token:      token.B,
				Collateral: 0,
				Title:      "No",
				Percentage: 0,
			},
		},
	}
	e6 := &Event{
		ID:              uuid.New(),
		RWMutex:         sync.RWMutex{},
		Tag:             Politic,
		LogoLink:        "https://tfj34t-88-201-232-88.ru.tuna.am/logo/elections2024.jpg",
		Title:           "2024 USA elections",
		TotalCollateral: 0,
		Bets: map[token.Token]*Bet{
			token.A: {
				Token:      token.A,
				Collateral: 0,
				Title:      "Yes",
				Percentage: 0,
			},

			token.B: {
				Token:      token.B,
				Collateral: 0,
				Title:      "No",
				Percentage: 0,
			},
		},
	}
	e7 := &Event{
		ID:              uuid.New(),
		RWMutex:         sync.RWMutex{},
		Tag:             Politic,
		LogoLink:        "https://tfj34t-88-201-232-88.ru.tuna.am/logo/elections2024.jpg",
		Title:           "2024 USA elections",
		TotalCollateral: 0,
		Bets: map[token.Token]*Bet{
			token.A: {
				Token:      token.A,
				Collateral: 0,
				Title:      "Yes",
				Percentage: 0,
			},

			token.B: {
				Token:      token.B,
				Collateral: 0,
				Title:      "No",
				Percentage: 0,
			},
		},
	}
	e8 := &Event{
		ID:              uuid.New(),
		RWMutex:         sync.RWMutex{},
		Tag:             Politic,
		LogoLink:        "https://tfj34t-88-201-232-88.ru.tuna.am/logo/elections2024.jpg",
		Title:           "2024 USA elections",
		TotalCollateral: 0,
		Bets: map[token.Token]*Bet{
			token.A: {
				Token:      token.A,
				Collateral: 0,
				Title:      "Yes",
				Percentage: 0,
			},

			token.B: {
				Token:      token.B,
				Collateral: 0,
				Title:      "No",
				Percentage: 0,
			},
		},
	}
	e9 := &Event{
		ID:              uuid.New(),
		RWMutex:         sync.RWMutex{},
		Tag:             Politic,
		LogoLink:        "https://tfj34t-88-201-232-88.ru.tuna.am/logo/elections2024.jpg",
		Title:           "2024 USA elections",
		TotalCollateral: 0,
		Bets: map[token.Token]*Bet{
			token.A: {
				Token:      token.A,
				Collateral: 0,
				Title:      "Yes",
				Percentage: 0,
			},

			token.B: {
				Token:      token.B,
				Collateral: 0,
				Title:      "No",
				Percentage: 0,
			},
		},
	}

	s := &storage{
		mu: sync.RWMutex{},
		m:  make(map[uuid.UUID]*Event),
	}
	s.m[e.ID] = e
	s.m[e1.ID] = e1
	s.m[e2.ID] = e2
	s.m[e3.ID] = e3
	s.m[e4.ID] = e4
	s.m[e5.ID] = e5
	s.m[e6.ID] = e6
	s.m[e7.ID] = e7
	s.m[e8.ID] = e8
	s.m[e9.ID] = e9
	return &EventKeeper{
		&snapshot{
			mu:   sync.RWMutex{},
			list: make([]*EventDTO, 0),
		},
		s,
	}
}

var (
	once      sync.Once
	singleton *EventKeeper
)

func Keeper() *EventKeeper {
	once.Do(func() {
		singleton = newEventKeeper()
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
		bets := v.GetBets()
		collateral := bets[0].Collateral + bets[1].Collateral
		l = append(l, &EventDTO{
			ID:         v.ID,
			Tag:        v.Tag,
			Collateral: collateral,
			LogoLink:   v.LogoLink,
			Title:      v.Title,
			Bets:       bets,
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

func (k *EventKeeper) GetSnapshot(_ context.Context, tag EventTag, page int) ([]*EventDTO, int, error) {
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

	pageSize := 4
	start := (page - 1) * pageSize
	end := start + pageSize

	totalPages := int(math.Ceil(float64(len(l)) / float64(pageSize)))

	if start >= len(l) {
		return []*EventDTO{}, totalPages, nil
	}

	if end > len(l) {
		end = len(l)
	}

	return l[start:end], totalPages, nil
}
