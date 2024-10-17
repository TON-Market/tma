package market

import (
	"context"
	"errors"
	"fmt"
	"github.com/TON-Market/tma/server/datatype/token"
	"github.com/TON-Market/tma/server/db"
	"github.com/TON-Market/tma/server/utils"
	"github.com/google/uuid"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"github.com/tonkeeper/tongo/wallet"
	"log"
	"math"
	"slices"
	"sort"
	"sync"
	"time"
)

type snapshot struct {
	sync.RWMutex
	list []EventDTO
}

type Market struct {
	wallet    wallet.Wallet
	WsCh      chan *EventDTO
	ch        chan *DepositReq
	client    *liteapi.Client
	snapshot  *snapshot
	persistor *persistor
	runtimer  *runtimer
}

var ErrMarketAddEventFailed = errors.New("market add event failed")

func (m *Market) SaveDealUnchecked(ctx context.Context, d *Deal) error {
	es := m.runtimer.getEventState(ctx, d.EventID)
	if !es.isActive {
		return ErrEventClosed
	}
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
	retryTicker := time.NewTicker(1 * time.Minute)

	for i := 0; i < 1; i++ {
		go m.verifyIncomeTransactions()
	}

	go func() {
		for range ticker.C {
			if err := m.makeSnapshot(ctx); err != nil {
				log.Println(err)
			}
		}
		defer ticker.Stop()
	}()

	go func() {
		for range retryTicker.C {
			if err := m.verifyPending(ctx); err != nil {
				log.Printf("[ERROR]: %s\n\n", err.Error())
			}
		}
	}()
}

func (m *Market) verifyPending(ctx context.Context) error {
	pendingDeals, err := m.persistor.getPendingDeals(ctx)
	if err != nil {
		return fmt.Errorf("verify pending deals failed: %w", err)
	}

	for _, deal := range pendingDeals {

		if err = m.Deposit(ctx, &DepositReq{
			deal.ID,
			OK,
			time.Now(),
		}); err != nil {
			log.Printf("[ERROR] retry deposit for dealID: %s failed: %w\n\n", deal.ID.String(), err)
		}
	}
	return nil
}

func (m *Market) ReadFromSnapshot(_ context.Context, tag Tag, page int) ([]EventDTO, int, error) {
	m.snapshot.RLock()
	l := slices.Clone(m.snapshot.list)
	m.snapshot.RUnlock()

	r := make([]EventDTO, 0)

	if tag == All {
		r = l
	} else {
		for _, e := range l {
			if e.Tag == tag {
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

func (m *Market) GetUserAssets(ctx context.Context, addr string) ([]*AssetDTO, string, error) {
	assetList, err := m.persistor.getUserAssets(ctx, addr)
	if err != nil {
		return nil, "", fmt.Errorf("market get user assets failed: %w", err)
	}

	assetDtoList := make([]*AssetDTO, 0, len(assetList))
	totalInMarket := tlb.Grams(0)

	for _, asset := range assetList {
		totalInMarket += asset.CollateralStaked

		eventCopy, err := m.persistor.getCopyByID(ctx, asset.EventID)
		if err != nil {
			return nil, "", fmt.Errorf("market get user assets failed: %w", err)
		}

		assetDto := &AssetDTO{
			EventTitle:       eventCopy.Title,
			BetTitle:         eventCopy.BetMap[asset.Token].Title,
			CollateralStaked: utils.GramsToStringInFloat(asset.CollateralStaked),
			Size:             utils.GramsToStringInFloat(asset.Size),
		}

		assetDtoList = append(assetDtoList, assetDto)
	}
	return assetDtoList, utils.GramsToStringInFloat(totalInMarket), nil
}

func (m *Market) makeSnapshot(ctx context.Context) error {
	list := make([]EventDTO, 0, len(m.snapshot.list))
	eventStates := m.runtimer.snapshot(ctx)

	for id, state := range eventStates {
		event, err := m.persistor.getCopyByID(ctx, id)
		if err != nil && !errors.Is(err, ErrEventNotExist) {
			log.Printf("[ERROR] snapshot failed: %s\n\n", err.Error())
			continue
		}

		if len(state.betStateMap) != len(event.BetMap) {
			log.Printf("[ERROR] snapshot failed: persist and runtime bet maps not equal\n\n")
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
			Percentage: utils.FloatToString(state.betStateMap[t].percentage),
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
		pk, err := wallet.SeedToPrivateKey(SEED)
		if err != nil {
			log.Fatalln(err)
		}

		w, err := wallet.New(pk, wallet.HighLoadV2R2, client)
		if err != nil {
			log.Fatalln(err)
		}

		singleton = &Market{
			w,
			make(chan *EventDTO),
			make(chan *DepositReq, 10000),
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
	log.Printf("[INFO] deposit request registered, id: %s\n\n", dr.ID.String())
	return nil
}

var ErrVerifyTransaction = errors.New("[ERROR] verify transaction")

func (m *Market) verifyIncomeTransactions() {
	for dr := range m.ch {
		ctx := context.TODO()

		if dr.DepositStatus == ERROR {
			if err := m.persistor.declineDeal(ctx, dr.ID); err != nil {
				log.Printf("[INFO] %v, id: %s: err decline deal: %v\n\n", ErrVerifyTransaction, dr.ID.String(), err)
			}
			continue
		}

		if dr.Time.Add(time.Minute).After(time.Now()) {
			m.ch <- dr
			continue
		}

		userRawAddress, err := m.persistor.getUserAddressByUncheckedDealID(ctx, dr.ID)
		if err != nil {
			log.Printf("%v, id: %s: %v\n\n", ErrVerifyTransaction, dr.ID.String(), err)
			continue
		}
		accountID, err := ton.ParseAccountID(userRawAddress)
		if err != nil {
			log.Printf("%v, id: %s: can't parse address: %s: %v\n\n", ErrVerifyTransaction, dr.ID.String(), userRawAddress, err)
			continue
		}

		getLastTransactions := func() ([]ton.Transaction, error) {
			for i := 0; i < 10; i++ {
				l, err := m.client.GetLastTransactions(ctx, accountID, 10)
				if err != nil {
					time.Sleep(5 * time.Second)
					continue
				}
				return l, nil
			}
			return nil, fmt.Errorf("transaction not found: %v", err)
		}

		log.Printf("[INFO] getting transaction list for user: %s\n\n", userRawAddress)
		trxList, err := getLastTransactions()

		if err != nil {
			log.Printf("[ERROR] %v, id: %s: user_raw_address: %s: %v\n\n", ErrVerifyTransaction, dr.ID.String(), userRawAddress, err)
			continue
		}

		deal, err := m.iterateTransactionList(ctx, dr, userRawAddress, trxList)
		if err != nil {
			log.Printf("[ERROR] %v\n\n", err)
			continue
		}

		if err := m.sendToSocket(ctx, deal.EventID); err != nil {
			log.Printf("[ERROR] %v\n\n", err)
			continue
		}
	}
}

var ErrTransactionNotFound = errors.New("transaction not found in block chain")

func (m *Market) iterateTransactionList(ctx context.Context, dr *DepositReq, userRawAddress string, trxList []ton.Transaction) (*Deal, error) {
	for _, trx := range trxList {
		var t wallet.TextComment
		if err := tlb.Unmarshal((*boc.Cell)(&trx.Msgs.OutMsgs.Values()[0].Value.Body.Value), &t); err != nil {
			log.Printf("[WARNING] verify transaction, id: %s, user_raw_address: %s, can't unmarshal boc: %v\n\n", dr.ID.String(), userRawAddress, err)
			continue
		}

		idStr := wallet.TextComment(dr.ID.String())

		if t == idStr {
			deal, err := m.persistor.verifyDealAndGet(ctx, dr.ID)
			if err != nil {
				return nil, fmt.Errorf("iterate transaction list failed: id: %s: user_raw_address: %s: %w", dr.ID.String(), userRawAddress, err)
			}

			if err = m.runtimer.deposit(ctx, deal); err != nil {
				return nil, fmt.Errorf("iterate transaction list failed: id: %s: user_raw_address: %s: %w", dr.ID.String(), userRawAddress, err)
			}
			return deal, nil
		}
	}
	return nil, fmt.Errorf("iterate transaction list failed: id: %s: user_raw_address: %s: %w", dr.ID.String(), userRawAddress, ErrTransactionNotFound)
}

func (m *Market) sendToSocket(ctx context.Context, id uuid.UUID) error {
	eventCopy, err := m.persistor.getCopyByID(ctx, id)
	if err != nil {
		return fmt.Errorf("send to socket failed: %w", err)
	}

	eventState := m.runtimer.getEventState(ctx, id)

	eventDTO := &EventDTO{
		ID:              eventCopy.ID,
		Tag:             eventCopy.Tag,
		LogoLink:        eventCopy.LogoLink,
		Title:           eventCopy.Title,
		Collateral:      utils.GramsToStringInFloat(eventState.collateral),
		CollateralGrams: eventState.collateral,
		Bets: []*BetDTO{
			{
				Token:      token.A,
				Title:      eventCopy.BetMap[token.A].Title,
				Percentage: utils.FloatToString(eventState.betStateMap[token.A].percentage),
			},
			{
				Token:      token.B,
				Title:      eventCopy.BetMap[token.B].Title,
				Percentage: utils.FloatToString(eventState.betStateMap[token.B].percentage),
			},
		},
	}

	m.WsCh <- eventDTO
	return nil
}
