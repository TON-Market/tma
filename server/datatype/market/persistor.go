package market

import (
	"context"
	"errors"
	"fmt"
	"github.com/TON-Market/tma/server/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tonkeeper/tongo/tlb"
	"strconv"
	"sync"
)

type eventStorage struct {
	sync.RWMutex
	m map[uuid.UUID]*Event
}

var (
	ErrEventNotExist     = errors.New("event not exist")
	ErrEventAlreadyExist = errors.New("event already exist")
)

var ErrSaveEvent = errors.New("save event failed")

func (es *eventStorage) saveEvent(_ context.Context, e *Event) error {
	es.RLock()
	defer es.RUnlock()

	if _, ok := es.m[e.ID]; ok {
		return fmt.Errorf("%v: %v: %s", ErrSaveEvent, ErrEventAlreadyExist, e.ID.String())
	}

	es.m[e.ID] = e

	return nil
}

var ErrGetByIdEvent = errors.New("get by id event failed")

func (es *eventStorage) getCopyByID(_ context.Context, id uuid.UUID) (Event, error) {
	es.RLock()
	defer es.RUnlock()

	if _, ok := es.m[id]; !ok {
		return Event{}, fmt.Errorf("%v: %v: %s", ErrGetByIdEvent, ErrEventNotExist, id.String())
	}

	v := es.m[id]
	return *v, nil
}

type persistor struct {
	*eventStorage
	pool *pgxpool.Pool
}

var (
	ErrPersistDeal = errors.New("persist deal failed")
)

func (p *persistor) saveDeal(ctx context.Context, d *Deal) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%v: %v: %v", ErrPersistDeal, db.ErrOpenTransaction, err)
	}

	defer tx.Rollback(ctx)

	dq := `INSERT INTO public.deals (id, event_id, token, collateral, size, user_raw_addr, deal_status)
          VALUES ($1, $2, $3, $4, $5, $6, $7)`

	udq := `INSERT INTO public.user_deals (user_raw_addr, deal_id)
            VALUES ($1, $2)`

	colStr := strconv.FormatUint(uint64(d.Collateral), 10)

	if _, err := tx.Exec(ctx, dq, d.ID, d.EventID, d.Token, colStr, d.Size, d.UserRawAddr, d.DealStatus); err != nil {
		return fmt.Errorf("%v: %v: %v", ErrPersistDeal, db.ErrTransactionFailed, err)
	}

	if _, err := tx.Exec(ctx, udq, d.UserRawAddr, d.ID); err != nil {
		return fmt.Errorf("%v: %v: %v", ErrPersistDeal, db.ErrTransactionFailed, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%v: %v: %v", ErrPersistDeal, db.ErrCommitTransaction, err)
	}

	return nil
}

var ErrVerifyDealAndGet = errors.New("verify deal and get failed")

func (p *persistor) verifyDealAndGet(ctx context.Context, id uuid.UUID) (*Deal, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%v: %v: %v", ErrVerifyDealAndGet, db.ErrOpenTransaction, err)
	}

	defer tx.Rollback(ctx)

	qu := `UPDATE public.deals SET deal_status = 1 WHERE id = $1`
	qg := `SELECT id, event_id, token, collateral, size, user_raw_addr, deal_status
           FROM deals WHERE id = $1`

	if _, err := tx.Exec(ctx, qu, id); err != nil {
		return nil, fmt.Errorf("%v: %v: %v", ErrVerifyDealAndGet, db.ErrTransactionFailed, err)
	}

	var deal Deal
	if err := tx.QueryRow(ctx, qg, id).Scan(&deal.ID, &deal.EventID, &deal.Token, &deal.Collateral, &deal.Size, &deal.UserRawAddr, &deal.DealStatus); err != nil {
		return nil, fmt.Errorf("%v: %v: %v", ErrVerifyDealAndGet, db.ErrTransactionFailed, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%v: %v: %v", ErrVerifyDealAndGet, db.ErrCommitTransaction, err)
	}

	return &deal, nil
}

var ErrGetPendingDeals = errors.New("get pending deals failed")

func (p *persistor) getPendingDeals(ctx context.Context) ([]*Deal, error) {
	dealList := make([]*Deal, 0)

	q := `SELECT id, event_id, token, collateral, size, user_raw_addr, deal_status
		  FROM deals WHERE deal_status = $1`

	rows, err := p.pool.Query(ctx, q, PENDING)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", ErrGetPendingDeals, err)
	}

	defer rows.Close()

	for rows.Next() {
		var deal Deal
		var collateral string

		if err := rows.Scan(&deal.ID, &deal.EventID, &deal.UserRawAddr, &deal.Token, &collateral, &deal.Size, &deal.DealStatus); err != nil {
			return nil, fmt.Errorf("%v: %v", ErrGetPendingDeals, err)
		}

		grams, err := strconv.ParseUint(collateral, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%v: parse collateral failed: %v", ErrGetPendingDeals, err)
		}

		deal.Collateral = tlb.Grams(grams)
		dealList = append(dealList, &deal)
	}

	return dealList, nil
}

func (p *persistor) deleteDeal(ctx context.Context, id uuid.UUID) error {
	q := `DELETE FROM public.deals WHERE id = $1`

	if _, err := p.pool.Exec(ctx, q, id); err != nil {
		return fmt.Errorf("delete deal with id: %s failed: %v", id.String(), err)
	}

	return nil
}

func (p *persistor) getUserAddressByPendingDealID(ctx context.Context, id uuid.UUID) (string, error) {
	q := `SELECT user_raw_addr FROM deals WHERE id = $1`

	var userRawAddr string

	if err := p.pool.QueryRow(ctx, q, id).Scan(&userRawAddr); err != nil {
		return "", fmt.Errorf("select user_raw_address with deal id: %s falied: %v", id.String(), err)
	}

	return userRawAddr, nil
}
