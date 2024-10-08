package user

import (
	"context"
	"errors"
	"github.com/TON-Market/tma/server/datatype"
	"github.com/TON-Market/tma/server/datatype/deal"

	"sync"
)

type User struct {
	ID string
	*datatype.AccountInfo
	DealList []*deal.Deal
}

type Storage struct {
	sync.RWMutex
	m map[string]*User
}

var (
	once      sync.Once
	singleton *Storage
)

func GetStorage() *Storage {
	once.Do(func() {
		singleton = &Storage{
			RWMutex: sync.RWMutex{},
			m:       make(map[string]*User),
		}
	})
	return singleton
}

var (
	ErrUserAlreadyExists = errors.New("err user already exists")
	ErrUserDoesntExist   = errors.New("err user doesn't exist")
)

func (s *Storage) AddUser(_ context.Context, info *datatype.AccountInfo) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.m[info.Address.Raw]; ok {
		return ErrUserAlreadyExists
	}

	u := &User{
		ID:          info.Address.Raw,
		AccountInfo: info,
		DealList:    make([]*deal.Deal, 0),
	}

	s.m[info.Address.Raw] = u

	return nil
}

func (s *Storage) GetUser(_ context.Context, address string) (*User, error) {
	s.RLock()
	defer s.RUnlock()

	u, ok := s.m[address]
	if !ok {
		return nil, ErrUserDoesntExist
	}

	return u, nil
}

func (s *Storage) AddDeal(ctx context.Context, info *datatype.AccountInfo, d *deal.Deal) error {
	u, err := s.GetUser(ctx, info.Address.Raw)
	if err != nil {
		return err
	}

	if err := deal.GetStorage().AddDeal(ctx, d); err != nil {
		return err
	}

	u.DealList = append(u.DealList, d)

	return nil
}
