package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/TON-Market/tma/server/datatype/market"
	"github.com/TON-Market/tma/server/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tonkeeper/tongo/tlb"
	"strconv"
	"sync"
)

type User struct {
	RawAddr  string
	DealList []*market.Deal
}

type Storage struct {
	pool *pgxpool.Pool
}

var ErrSaveUser = errors.New("save user failed")

func (s *Storage) Save(ctx context.Context, u *User) error {
	q := `INSERT INTO users (raw_addr) VALUES ($1)`

	if _, err := s.pool.Exec(ctx, q, u.RawAddr); err != nil {
		return fmt.Errorf("%v: %v", ErrSaveUser, err)
	}
	return nil
}

var (
	ErrUserNotFound = errors.New("user not found in db")
	ErrGetUser      = errors.New("get user failed")
)

func (s *Storage) Get(ctx context.Context, addr string) (*User, error) {
	q := `SELECT raw_addr FROM users WHERE raw_addr = $1`

	var user User

	err := s.pool.QueryRow(ctx, q, addr).Scan(&user.RawAddr)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Join(ErrUserNotFound, err)
		}
		return nil, fmt.Errorf("%v: %v", ErrGetUser, err)
	}

	user.DealList = []*market.Deal{}

	qDeals := `SELECT d.id, d.event_id, d.token, d.collateral, d.size
               FROM deals d
               JOIN user_deals ud ON d.id = ud.deal_id
               WHERE ud.user_raw_addr = $1`

	rows, err := s.pool.Query(ctx, qDeals, addr)
	if err != nil {
		return nil, fmt.Errorf("get user deals failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var deal market.Deal
		var collateral string

		if err := rows.Scan(&deal.ID, &deal.EventID, &deal.Token, &collateral, &deal.Size); err != nil {
			return nil, fmt.Errorf("scan deal failed: %v", err)
		}

		grams, err := strconv.ParseUint(collateral, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse collateral failed: %v", err)
		}
		deal.Collateral = tlb.Grams(grams)

		user.DealList = append(user.DealList, &deal)
	}

	return &user, nil
}

var (
	once      sync.Once
	singleton *Storage
)

func UserStorage() *Storage {
	once.Do(func() {
		singleton = &Storage{
			pool: db.Get(),
		}
	})
	return singleton
}
