package market

import (
	"errors"
	"github.com/TON-Market/tma/server/datatype/token"
	"github.com/google/uuid"
	"github.com/tonkeeper/tongo/tlb"
	"sync"
)

type DepositStatus int

const (
	BANK_ADDR = "EQBRW9rjhRUNL-Sy4swYbMzm2MgvlhC2DWIZFhYp2JnSoJaA"
	SEED      = "example consider fiscal mail guitar tiger duck exhibit ancient series differ wealth mix kitchen cactus upgrade unable yellow impact confirm denial mesh during dove"
)

const (
	OK DepositStatus = iota
	ERROR
)

type Tag int

const (
	Politic Tag = iota
	Economics
	Crypto
	Culture
	Other
	All
)

type AssetDTO struct {
	EventTitle       string `json:"eventTitle"`
	BetTitle         string `json:"betTitle"`
	CollateralStaked string `json:"collateralStaked"`
	Size             string `json:"size"`
}

type Asset struct {
	UserRawAddress   string
	EventID          uuid.UUID
	CollateralStaked tlb.Grams
	Token            token.Token
	Size             tlb.Grams
}

type BetDTO struct {
	Token      token.Token `json:"token"`
	Title      string      `json:"title"`
	Percentage string      `json:"percentage"`
	LogoLink   string      `json:"logoLink"`
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
	Declined
)

// Deal persist user deal
type Deal struct {
	ID          uuid.UUID
	EventID     uuid.UUID
	UserRawAddr string
	token.Token
	Collateral tlb.Grams
	Size       tlb.Grams
	DealStatus DealStatus
	Attempts   int
}

// Bet persist data
type Bet struct {
	EventID uuid.UUID
	token.Token
	Title    string
	LogoLink string
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
	isActive    bool
	collateral  tlb.Grams
	betStateMap map[token.Token]*betState
}

type eventRuntime struct {
	sync.RWMutex
	isActive      bool
	eventID       uuid.UUID
	betRuntimeMap map[token.Token]*betRuntime
}

var ErrEventClosed = errors.New("event closed")

func (er *eventRuntime) deposit(t token.Token, g tlb.Grams) {
	er.RLock()
	defer er.RUnlock()
	er.betRuntimeMap[t].deposit(g)
}

func (er *eventRuntime) getState() *eventState {
	er.RLock()
	defer er.RUnlock()

	es := &eventState{
		isActive:    er.isActive,
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
		bs.percentage = (float64(bs.collateral) / float64(es.collateral)) * 100
	}

	return es
}

func (er *eventRuntime) close() {
	er.Lock()
	defer er.Unlock()

	er.isActive = false
}
