package market

import (
	"context"
	"errors"
	"fmt"
	"github.com/TON-Market/tma/server/db"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"github.com/tonkeeper/tongo/wallet"
	"math"
	"slices"
	"sort"
	"strconv"
	"sync"
	"time"
)

const BANK_ADDR = "EQBRW9rjhRUNL-Sy4swYbMzm2MgvlhC2DWIZFhYp2JnSoJaA"

type snapshot struct {
	sync.RWMutex
	list []EventDTO
}

type Market struct {
	//revisor   *revisor
	client    *liteapi.Client
	snapshot  *snapshot
	persistor *persistor
	runtimer  *runtimer
}

var ErrMarketAddEventFailed = errors.New("market add event failed")

func (m *Market) SaveDealUnchecked(ctx context.Context, d *Deal) error {
	d.DealStatus = Unchecked
	if err := m.persistor.saveDeal(ctx, d); err != nil {
		return fmt.Errorf("market save deal unchecked failed: %v", err)
	}
	return nil
}

func (m *Market) AddEvent(ctx context.Context, e *Event) error {
	e.ID = uuid.New()
	for _, b := range e.BetMap {
		b.EventID = e.ID
	}
	if err := m.persistor.saveEvent(ctx, e); err != nil {
		return fmt.Errorf("%v: %v", ErrMarketAddEventFailed, err)
	}
	if err := m.runtimer.saveEvent(ctx, e); err != nil {
		return fmt.Errorf("%v: %v", ErrMarketAddEventFailed, err)
	}
	return nil
}

func (m *Market) Start(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)

	go func() {
		for range ticker.C {
			if err := m.makeSnapshot(ctx); err != nil {
				log.Println(err)
			}
		}
		defer ticker.Stop()
	}()
}

func (m *Market) ReadFromSnapshot(_ context.Context, tag Tag, page int) ([]EventDTO, int, error) {
	m.snapshot.RLock()
	l := slices.Clone(m.snapshot.list)
	m.snapshot.RUnlock()

	r := make([]EventDTO, 0)

	if tag == No {
		r = l
	} else {
		for _, e := range l {
			if e.Tag != tag {
				r = append(r, e)
			}
		}
	}

	pgSz := 4
	start := (page - 1) * pgSz
	end := start + pgSz

	totalPages := int(math.Ceil(float64(len(r)) / float64(pgSz)))

	if start >= len(r) {
		return []EventDTO{}, totalPages, nil
	}

	if end > len(r) {
		end = len(r)
	}

	return r[start:end], totalPages, nil
}

func (m *Market) makeSnapshot(ctx context.Context) error {
	list := make([]EventDTO, 0, len(m.snapshot.list))
	eventStates := m.runtimer.snapshot(ctx)

	for id, state := range eventStates {
		event, err := m.persistor.getCopyByID(ctx, id)
		if err != nil && !errors.Is(err, ErrEventNotExist) {
			log.Printf(fmt.Sprintf("snapshot failed: %s\n", err.Error()))
			continue
		}

		if len(state.betStateMap) != len(event.BetMap) {
			log.Printf("snapshot failed: persist and runtime bet maps not equal\n")
			continue
		}

		betDTOList := m.snapshotBets(ctx, event, state)

		collateral := strconv.FormatUint(uint64(state.collateral), 10)

		eventDTO := EventDTO{
			ID:              id,
			Tag:             event.Tag,
			LogoLink:        event.LogoLink,
			Title:           event.Title,
			Collateral:      collateral,
			CollateralGrams: state.collateral,
			Bets:            betDTOList,
		}

		list = append(list, eventDTO)
	}

	m.snapshot.Lock()
	defer m.snapshot.Unlock()

	sort.Slice(list, func(i, j int) bool {
		return list[i].CollateralGrams > list[j].CollateralGrams
	})

	m.snapshot.list = list
	return nil
}

func (m *Market) snapshotBets(_ context.Context, e Event, state *eventState) []*BetDTO {
	betDTOList := make([]*BetDTO, 0, len(e.BetMap))

	for t, b := range e.BetMap {
		bDTO := &BetDTO{
			Token:      t,
			Title:      b.Title,
			Percentage: fmt.Sprintf("%.0f", state.betStateMap[t].percentage),
		}
		betDTOList = append(betDTOList, bDTO)
	}

	return betDTOList
}

var (
	once      sync.Once
	singleton *Market
)

func GetMarket() *Market {
	once.Do(func() {
		client, err := liteapi.NewClientWithDefaultMainnet()
		if err != nil {
			log.Fatalln(err)
		}
		singleton = &Market{
			client,
			&snapshot{
				sync.RWMutex{},
				make([]EventDTO, 0),
			},
			&persistor{
				&eventStorage{
					sync.RWMutex{},
					make(map[uuid.UUID]*Event),
				},
				db.Get(),
			},
			&runtimer{
				sync.RWMutex{},
				make(map[uuid.UUID]*eventRuntime),
			},
		}
	})
	return singleton
}

func (m *Market) _Deposit(ctx context.Context, d *Deal) error {
	d.ID = uuid.New()
	if err := m.persistor.saveDeal(ctx, d); err != nil {
		return fmt.Errorf("market deposit failed: %v", err)
	}
	if err := m.runtimer.deposit(ctx, d); err != nil {
		return fmt.Errorf("market deposit failed: %v", err)
	}
	return nil
}

var ErrTransactionNotVerified = errors.New("not verified transaction")

func (m *Market) Deposit(ctx context.Context, id uuid.UUID, ds DepositStatus) error {
	if ds == ERROR {
		if err := m.persistor.deleteDeal(ctx, id); err != nil {
			return fmt.Errorf("deposit failed: %v", err)
		}
		return nil
	}

	accountId := ton.MustParseAccountID(BANK_ADDR)

	trxList, err := m.client.GetLastTransactions(ctx, accountId, 1)
	if err != nil {
		return fmt.Errorf("deposit failed: %v", err)
	}

	trx := trxList[0]

	var t wallet.TextComment
	if err := tlb.Unmarshal((*boc.Cell)(&trx.Msgs.InMsg.Value.Value.Body.Value), &t); err != nil {
		return fmt.Errorf("deposit failed: %v", err)
	}

	idStr := wallet.TextComment(id.String())

	if t == idStr {
		return m.persistor.verifyDeal(ctx, id)
	}

	return ErrTransactionNotVerified
}
