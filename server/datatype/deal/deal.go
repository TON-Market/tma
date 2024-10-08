package deal

import (
	"context"
	"github.com/TON-Market/tma/server/datatype/token"
	"github.com/google/uuid"
	"sync"
)

type Deal struct {
	ID      uuid.UUID
	EventID uuid.UUID
	token.Token
	Collateral uint64
	Size       float64
}

type Storage struct {
	sync.RWMutex
	m map[uuid.UUID]*Deal
}

var (
	dealStorageOnce      sync.Once
	dealStorageSingleton *Storage
)

func GetStorage() *Storage {
	dealStorageOnce.Do(func() {
		dealStorageSingleton = &Storage{
			RWMutex: sync.RWMutex{},
			m:       make(map[uuid.UUID]*Deal),
		}
	})
	return dealStorageSingleton
}

func (s *Storage) AddDeal(_ context.Context, deal *Deal) error {
	s.Lock()
	defer s.RUnlock()

	dealId := uuid.New()
	deal.ID = dealId
	s.m[dealId] = deal
	return nil
}
