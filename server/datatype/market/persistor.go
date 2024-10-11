package market

import (
	"context"
	"errors"
	"fmt"
	"github.com/TON-Market/tma/server/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
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

func (p *persistor) verifyDeal(ctx context.Context, id uuid.UUID) error {
	q := `UPDATE public.deals SET deal_status = 1 WHERE id = $1`

	if _, err := p.pool.Exec(ctx, q, id); err != nil {
		return fmt.Errorf("verify deal failed: %v", err)
	}

	return nil
}

func (p *persistor) deleteDeal(ctx context.Context, id uuid.UUID) error {
	q := `DELETE FROM public.deals WHERE id = $1`

	if _, err := p.pool.Exec(ctx, q, id); err != nil {
		return fmt.Errorf("delete deal with id: %s failed: %v", id.String(), err)
	}

	return nil
}
