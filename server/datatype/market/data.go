package market

import (
	"github.com/TON-Market/tma/server/datatype/token"
	"github.com/google/uuid"
	"github.com/tonkeeper/tongo/tlb"
	"sync"
)

type DepositStatus int

const (
	OK DepositStatus = iota
	ERROR
	PENDING
)

type Tag int

const (
	Politic Tag = iota
	Economics
	Crypto
	Culture
	Other
	No
)

type BetDTO struct {
	Token      token.Token `json:"token"`
	Title      string      `json:"title"`
	Percentage string      `json:"percentage"`
}

type EventDTO struct {
	ID              uuid.UUID `json:"id"`
	Tag             Tag       `json:"tag"`
	LogoLink        string    `json:"logoLink"`
	Title           string    `json:"title"`
	Collateral      string    `json:"collateral"`
	CollateralGrams tlb.Grams
	Bets            []*BetDTO `json:"bets"`
}

type DealStatus int

const (
	Unchecked DealStatus = iota
	Verified
)

// Deal persist user deal
type Deal struct {
	ID          uuid.UUID
	EventID     uuid.UUID
	UserRawAddr string
	token.Token
	Collateral tlb.Grams
	Size       float64
	DealStatus DealStatus
}

// Bet persist data
type Bet struct {
	EventID uuid.UUID
	token.Token
	Title string
}

// Event persist data
type Event struct {
	ID       uuid.UUID
	Tag      Tag
	LogoLink string
	Title    string
	BetMap   map[token.Token]*Bet
}

type betState struct {
	collateral tlb.Grams
	percentage float64
}

// betRuntime runtime bet state
type betRuntime struct {
	sync.RWMutex
	token.Token
	collateral tlb.Grams
}

func (br *betRuntime) deposit(g tlb.Grams) {
	br.Lock()
	defer br.Unlock()
	br.collateral += g
}

func (br *betRuntime) getState() *betState {
	br.RLock()
	defer br.RUnlock()
	return &betState{br.collateral, 0}
}

type eventState struct {
	collateral  tlb.Grams
	betStateMap map[token.Token]*betState
}

type eventRuntime struct {
	eventID       uuid.UUID
	betRuntimeMap map[token.Token]*betRuntime
}

func (er *eventRuntime) deposit(t token.Token, g tlb.Grams) {
	er.betRuntimeMap[t].deposit(g)
}

func (er *eventRuntime) getState() *eventState {
	es := &eventState{
		betStateMap: make(map[token.Token]*betState),
	}

	for t, br := range er.betRuntimeMap {
		bs := br.getState()
		es.betStateMap[t] = bs
		es.collateral += bs.collateral
	}

	for _, bs := range es.betStateMap {
		if bs.collateral == tlb.Grams(0) {
			bs.percentage = 0
			continue
		}
		bs.percentage = float64(bs.collateral) / float64(es.collateral)
	}

	return es
}
