package user

import (
	"context"
	"errors"
	"github.com/TON-Market/tma/server/datatype"
	"github.com/TON-Market/tma/server/datatype/token"
	"github.com/google/uuid"
	"sync"
)

type User struct {
	ID string
	*datatype.AccountInfo
	DealList []*Deal
}

type Storage struct {
	sync.RWMutex
	m map[string]*User
}

var (
	once      sync.Once
	singleton *Storage
)

func NewStorage() *Storage {
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
		DealList:    make([]*Deal, 0),
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

func (s *Storage) AddDeal(ctx context.Context, info *datatype.AccountInfo, deal *Deal) error {
	if err := NewDealStorage().addDeal(ctx, deal); err != nil {
		return err
	}

	u, err := s.GetUser(ctx, info.Address.Raw)
	if err != nil {
		return err
	}

	u.DealList = append(u.DealList, deal)

	return nil
}

type Deal struct {
	ID      uuid.UUID
	EventID uuid.UUID
	token.Token
	Collateral uint64
	Size       float64
}

type DealStorage struct {
	sync.RWMutex
	m map[uuid.UUID]*Deal
}

var (
	dealStorageOnce      sync.Once
	dealStorageSingleton *DealStorage
)

func NewDealStorage() *DealStorage {
	dealStorageOnce.Do(func() {
		dealStorageSingleton = &DealStorage{
			RWMutex: sync.RWMutex{},
			m:       make(map[uuid.UUID]*Deal),
		}
	})
	return dealStorageSingleton
}

func (s *DealStorage) addDeal(_ context.Context, deal *Deal) error {
	s.Lock()
	defer s.RUnlock()

	dealId := uuid.New()
	deal.ID = dealId
	s.m[dealId] = deal
	return nil
}
