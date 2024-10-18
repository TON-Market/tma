package market

import (
	"context"
	"github.com/TON-Market/tma/server/utils"
	"math"
	"slices"
	"sort"
	"sync"
)

type snapshot struct {
	sync.RWMutex
	list []EventDTO
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

func (m *Market) makeSnapshot(ctx context.Context) error {
	list := make([]EventDTO, 0, len(m.snapshot.list))
	eventStates := m.runtimer.snapshot(ctx)

	for eventID, state := range eventStates {
		eventCopy, err := m.persistor.getCopyByID(ctx, eventID)
		if err != nil {
			continue
		}

		betDTOList := m.snapshotBets(ctx, eventCopy, state)

		collateral := utils.GramsToStringInFloat(state.collateral)

		eventDTO := EventDTO{
			ID:              eventID,
			Tag:             eventCopy.Tag,
			LogoLink:        eventCopy.LogoLink,
			Title:           eventCopy.Title,
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
