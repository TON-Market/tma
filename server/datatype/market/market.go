package market

import (
	"context"
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
	"sync"
	"time"
)

type DepositReq struct {
	ID            uuid.UUID
	DepositStatus DepositStatus
	Time          time.Time
}

type Market struct {
	wallet       wallet.Wallet
	WsCh         chan *EventDTO
	depositReqCh chan *DepositReq
	profitCh     chan *UserProfit
	client       *liteapi.Client
	snapshot     *snapshot
	persistor    *persistor
	runtimer     *runtimer
}

func (m *Market) SaveDealUnchecked(ctx context.Context, d *Deal) error {
	es := m.runtimer.getEventState(ctx, d.EventID)
	if !es.isActive {
		return ErrEventClosed
	}
	d.DealStatus = Unchecked
	d.Attempts = 0
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
		return fmt.Errorf("market add event failed: %w", err)
	}
	if err := m.runtimer.saveEvent(ctx, e); err != nil {
		return fmt.Errorf("market add event failed: %w", err)
	}
	return nil
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

func (m *Market) Deposit(_ context.Context, dr *DepositReq) error {
	dr.Time = time.Now()
	m.depositReqCh <- dr
	log.Printf("[INFO] deposit request registered, id: %s\n\n", dr.ID.String())
	return nil
}

func (m *Market) checkIncomeTransactions(ctx context.Context) {
	go func() {
		for depositReq := range m.depositReqCh {
			if depositReq.DepositStatus == ERROR {
				if err := m.persistor.declineDeal(ctx, depositReq.ID); err != nil {
					log.Printf("[ERROR] check income transaction failed, deal_id: %s, %s\n\n",
						depositReq.ID.String(), err.Error())
					m.depositReqCh <- depositReq
				}
				continue
			}

			if depositReq.Time.Add(20 * time.Second).After(time.Now()) {
				m.depositReqCh <- depositReq
				continue
			}

			isDepositDelivered, _ := m.tryCheckIfDepositDelivered(ctx, depositReq)
			if isDepositDelivered {
				if err := m.confirmSuccessTransaction(ctx, depositReq); err != nil {
					log.Printf("[ERROR] confirm success transaction failed: %s\n\n", err.Error())
				}
				continue
			}

			attempts, err := m.persistor.attemptDeal(ctx, depositReq.ID)
			if err != nil {
				m.depositReqCh <- depositReq
				continue
			}

			if attempts > 5 {
				log.Printf("[ERROR] deal %v out of attempts\n\n", attempts)
				continue
			}

			m.depositReqCh <- depositReq
		}
	}()
}

func (m *Market) tryCheckIfDepositDelivered(ctx context.Context, depositReq *DepositReq) (bool, error) {
	depositReq.Time = depositReq.Time.Add(20 * time.Second)
	userRawAddress, err := m.persistor.getUserAddressByUncheckedDealID(ctx, depositReq.ID)
	if err != nil {
		return false, err
	}

	trxList, err := m.getLastTransactions(ctx, userRawAddress)
	if err != nil {
		return false, err
	}

	for _, trx := range trxList {
		isDepositTransaction, err := m.isDepositTransaction(ctx, trx, depositReq)
		if err != nil {
			continue
		}
		if isDepositTransaction {
			return true, nil
		}
	}

	return false, nil
}

func (m *Market) isDepositTransaction(_ context.Context, trx ton.Transaction, depositReq *DepositReq) (bool, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("[ALARM] Recovered in isDepositTransaction", r)
		}
	}()
	var t wallet.TextComment
	if err := tlb.Unmarshal((*boc.Cell)(&trx.Msgs.OutMsgs.Values()[0].Value.Body.Value), &t); err != nil {
		return false, err
	}

	comment := wallet.TextComment(depositReq.ID.String())

	if t == comment {
		return true, nil
	}

	return false, nil
}

func (m *Market) confirmSuccessTransaction(ctx context.Context, depositReq *DepositReq) error {
	deal, err := m.persistor.verifyDealAndGet(ctx, depositReq.ID)
	if err != nil {
		return err
	}
	if err = m.runtimer.deposit(ctx, deal); err != nil {
		return err
	}
	if err = m.sendToSocket(ctx, deal.EventID); err != nil {
		return err
	}
	return nil
}

func (m *Market) getLastTransactions(ctx context.Context, userRawAddress string) ([]ton.Transaction, error) {
	accountID, err := ton.ParseAccountID(userRawAddress)
	if err != nil {
		return nil, fmt.Errorf("failed get last transactions: %w", err)
	}
	for i := 0; i < 10; i++ {
		l, err := m.client.GetLastTransactions(ctx, accountID, 10)
		if err != nil {
			time.Sleep(10 * time.Second)
			continue
		}
		return l, nil
	}
	return nil, fmt.Errorf("failed get last transactions: tryied 10 times")
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
			make(chan *UserProfit, 10000),
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

func (m *Market) Start(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)

	for i := 0; i < 1; i++ {
		m.checkIncomeTransactions(ctx)
	}

	go func() {
		for range ticker.C {
			if err := m.makeSnapshot(ctx); err != nil {
				log.Println(err)
			}
		}
		defer ticker.Stop()
	}()

	m.startResendProcess(ctx)
}
