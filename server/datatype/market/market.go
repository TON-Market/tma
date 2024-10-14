package market

import (
	"context"
	"errors"
	"fmt"
	"github.com/TON-Market/tma/server/db"
	"github.com/TON-Market/tma/server/utils"
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
	"sync"
	"time"
)

const BANK_ADDR = "EQBRW9rjhRUNL-Sy4swYbMzm2MgvlhC2DWIZFhYp2JnSoJaA"

type snapshot struct {
	sync.RWMutex
	list []EventDTO
}

type Market struct {
	ch        chan *DepositReq
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

	for i := 0; i < 1; i++ {
		go m.verifyTransaction()
	}

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

func (m *Market) GetUserAssets(ctx context.Context, addr string) ([]*AssetDTO, error) {
	assetList, err := m.persistor.getUserAssets(ctx, addr)
	if err != nil {
		return nil, fmt.Errorf("market get user assets failed: %v", err)
	}

	assetDtoList := make([]*AssetDTO, 0, len(assetList))
	for _, asset := range assetList {
		assetDtoList = append(assetDtoList, &AssetDTO{
			EventTitle:       asset.EventTitle,
			Token:            asset.Token,
			CollateralStaked: utils.GramsToStringInFloat(asset.CollateralStaked),
			Size:             utils.GramsToStringInFloat(asset.Size),
		})
	}

	return assetDtoList, nil
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

		collateral := utils.GramsToStringInFloat(state.collateral)
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
			make(chan *DepositReq, 100),
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

type DepositReq struct {
	ID            uuid.UUID
	DepositStatus DepositStatus
	Time          time.Time
}

func (m *Market) Deposit(_ context.Context, dr *DepositReq) error {
	dr.Time = time.Now()
	m.ch <- dr
	log.Printf("[INFO] deposit request registered, id: %s\n", dr.ID.String())
	return nil
}

var ErrVerifyTransaction = errors.New("[ERROR] verify transaction")

func (m *Market) verifyTransaction() {
	for dr := range m.ch {
		if dr.Time.Add(time.Minute).After(time.Now()) {
			m.ch <- dr
			continue
		}
		ctx := context.TODO()

		if dr.DepositStatus == ERROR {
			if err := m.persistor.deleteDeal(ctx, dr.ID); err != nil {
				log.Printf("%v, id: %s: err delete unsigned transaction: %v\n", ErrVerifyTransaction, dr.ID.String(), err)
			}

			continue
		}

		userRawAddress, err := m.persistor.getUserAddressByPendingDealID(ctx, dr.ID)
		if err != nil {
			log.Printf("%v, id: %s: %v\n", ErrVerifyTransaction, dr.ID.String(), err)
			continue
		}
		accountID, err := ton.ParseAccountID(userRawAddress)
		if err != nil {
			log.Printf("%v, id: %s: can't parse address: %s: %v\n", ErrVerifyTransaction, dr.ID.String(), userRawAddress, err)
			continue
		}

		getLastTransactions := func() ([]ton.Transaction, error) {
			for i := 0; i < 5; i++ {
				l, err := m.client.GetLastTransactions(ctx, accountID, 5)
				if err != nil {
					time.Sleep(5 * time.Second)
					continue
				}
				return l, nil
			}
			return nil, fmt.Errorf("transaction not found: %v", err)
		}

		trxList, err := getLastTransactions()

		for _, trx := range trxList {
			var t wallet.TextComment
			if err := tlb.Unmarshal((*boc.Cell)(&trx.Msgs.OutMsgs.Values()[0].Value.Body.Value), &t); err != nil {
				log.Printf("[WARNING] verify transaction, id: %s, user_raw_address: %s, can't unmarshal boc: %v\n", dr.ID.String(), userRawAddress, err)
				continue
			}

			idStr := wallet.TextComment(dr.ID.String())

			if t == idStr {
				deal, err := m.persistor.verifyDealAndGet(ctx, dr.ID)
				if err != nil {
					log.Printf("%v, id: %s: user_raw_address: %s: %v\n", ErrVerifyTransaction, dr.ID.String(), userRawAddress, err)
					continue
				}

				if err := m.runtimer.deposit(ctx, deal); err != nil {
					log.Printf("%v, id: %s: user_raw_address: %s: %v\n", ErrVerifyTransaction, dr.ID.String(), userRawAddress, err)
				}

				continue
			}
		}

		log.Printf("%v, id: %s: user_raw_address: %s: transaction with deal id wasn't found\n", ErrVerifyTransaction, dr.ID.String(), userRawAddress)
	}
}
