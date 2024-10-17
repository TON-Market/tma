package market

import (
	"context"
	"errors"
	"fmt"
	"github.com/TON-Market/tma/server/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

var ErrVerifyDealAndGet = errors.New("verify deal and get failed")

func (p *persistor) verifyDealAndGet(ctx context.Context, id uuid.UUID) (*Deal, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%v: %v: %v", ErrVerifyDealAndGet, db.ErrOpenTransaction, err)
	}

	defer tx.Rollback(ctx)

	qa := `INSERT INTO assets (user_raw_address, event_id, collateral_staked, token, size)
           VALUES ($1, $2, $3, $4, $5)`

	qap := `UPDATE assets SET collateral_staked = $1, size = $2
            WHERE user_raw_address = $3 AND event_id = $4 AND token = $5`

	qag := `SELECT user_raw_address, event_id, collateral_staked, token, size
            FROM assets WHERE user_raw_address = $1 AND event_id = $2 AND token = $3`

	qu := `UPDATE public.deals SET deal_status = $1 WHERE id = $2`

	qg := `SELECT id, event_id, token, collateral, size, user_raw_addr, deal_status
           FROM deals WHERE id = $1`

	if _, err := tx.Exec(ctx, qu, Verified, id); err != nil {
		return nil, fmt.Errorf("%v: %v: %v", ErrVerifyDealAndGet, db.ErrTransactionFailed, err)
	}

	var deal Deal
	if err := tx.QueryRow(ctx, qg, id).Scan(&deal.ID, &deal.EventID, &deal.Token, &deal.Collateral, &deal.Size, &deal.UserRawAddr, &deal.DealStatus); err != nil {
		return nil, fmt.Errorf("%v: %v: %v", ErrVerifyDealAndGet, db.ErrTransactionFailed, err)
	}

	var asset Asset
	err = tx.QueryRow(ctx, qag, deal.UserRawAddr, deal.EventID, deal.Token).
		Scan(&asset.UserRawAddress, &asset.EventID, &asset.CollateralStaked, &asset.Token, &asset.Size)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%v: get existing asset: %v: %v", ErrVerifyDealAndGet, db.ErrTransactionFailed, err)
	}

	asset.UserRawAddress = deal.UserRawAddr
	asset.EventID = deal.EventID
	asset.Token = deal.Token
	asset.CollateralStaked += deal.Collateral
	asset.Size += deal.Size

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		if _, err := tx.Exec(ctx, qa, asset.UserRawAddress, asset.EventID, asset.CollateralStaked, asset.Token, asset.Size); err != nil {
			return nil, fmt.Errorf("%v: saving asset: %v: %v", ErrVerifyDealAndGet, db.ErrTransactionFailed, err)
		}
	}

	if err == nil {
		if _, err := tx.Exec(ctx, qap, asset.CollateralStaked, asset.Size, asset.UserRawAddress, asset.EventID, asset.Token); err != nil {
			return nil, fmt.Errorf("%v: updating asset: %v: %v", ErrVerifyDealAndGet, db.ErrTransactionFailed, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%v: %v: %v", ErrVerifyDealAndGet, db.ErrCommitTransaction, err)
	}

	return &deal, nil
}

func (p *persistor) declineDeal(ctx context.Context, id uuid.UUID) error {
	q := `UPDATE deals SET deal_status = $1 WHERE id = $2`

	if _, err := p.pool.Exec(ctx, q, Declined, id); err != nil {
		return fmt.Errorf("decline deal failed: %w", err)
	}
	return nil
}

func (p *persistor) deleteDeal(ctx context.Context, id uuid.UUID) error {
	q := `DELETE FROM deals WHERE id = $1`

	if _, err := p.pool.Exec(ctx, q, id); err != nil {
		return fmt.Errorf("delete deal with id: %s failed: %v", id.String(), err)
	}

	return nil
}

func (p *persistor) getUserAddressByUncheckedDealID(ctx context.Context, id uuid.UUID) (string, error) {
	q := `SELECT user_raw_addr FROM deals WHERE id = $1`

	var userRawAddr string

	if err := p.pool.QueryRow(ctx, q, id).Scan(&userRawAddr); err != nil {
		return "", fmt.Errorf("select user_raw_address with deal id: %s falied: %v", id.String(), err)
	}

	return userRawAddr, nil
}

func (p *persistor) getUserAssets(ctx context.Context, addr string) ([]*Asset, error) {
	q := `SELECT user_raw_address, event_id, collateral_staked, token, size
          FROM assets WHERE user_raw_address = $1`

	assetList := make([]*Asset, 0)

	rows, err := p.pool.Query(ctx, q, addr)
	if err != nil {
		return nil, fmt.Errorf("get user assets from db failed: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var asset Asset

		if err = rows.Scan(&asset.UserRawAddress, &asset.EventID, &asset.CollateralStaked, &asset.Token, &asset.Size); err != nil {
			return nil, fmt.Errorf("get user assets from db failed: %v", err)
		}

		assetList = append(assetList, &asset)
	}

	return assetList, nil
}

func (p *persistor) getEventAssets(ctx context.Context, id uuid.UUID) ([]*Asset, error) {
	q := `SELECT user_raw_address, event_id, collateral_staked, token, size
        FROM assets WHERE event_id = $1`

	assetList := make([]*Asset, 0)
	rows, err := p.pool.Query(ctx, q, id)
	if err != nil {
		return nil, fmt.Errorf("get assets failed: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var asset Asset
		if err = rows.Scan(&asset.UserRawAddress, &asset.EventID, &asset.CollateralStaked,
			&asset.Token, &asset.Size); err != nil {
			return nil, fmt.Errorf("get assets failed: %w", err)
		}

		assetList = append(assetList, &asset)
	}

	return assetList, nil
}

func (p *persistor) deleteAssets(ctx context.Context, id uuid.UUID) error {
	q := `DELETE FROM assets WHERE event_id = $1`

	if _, err := p.pool.Exec(ctx, q, id); err != nil {
		return fmt.Errorf("delete assets failed: %w", err)
	}

	return nil
}
