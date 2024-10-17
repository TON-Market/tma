package market

import (
	"context"
	"errors"
	"fmt"
	"github.com/TON-Market/tma/server/datatype/token"
	"github.com/google/uuid"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"github.com/tonkeeper/tongo/wallet"
	"log"
	"time"
)

func (m *Market) CloseEvent(ctx context.Context, id uuid.UUID, winToken token.Token) error {
	if err := m.runtimer.close(ctx, id); err != nil {
		return fmt.Errorf("close event failed: %w", err)
	}
	time.Sleep(3 * time.Minute)
	userTotalReturnMap, err := m.calcUsersProfit(ctx, id, winToken)
	if err != nil {
		return fmt.Errorf("close event failed: %w", err)
	}

	if err = m.profitUsers(ctx, userTotalReturnMap, id); err != nil {
		return fmt.Errorf("close event failed: %w", err)
	}

	//if err = m.persistor.deleteAssets(ctx, id); err != nil {
	//	return fmt.Errorf("delete assets failed: %w", err)
	//}

	if err = m.persistor.deleteEvent(ctx, id); err != nil {
		return fmt.Errorf("delete event failed: %w", err)
	}
	return nil
}

type UserTotalReturnMap map[string]tlb.Grams

var baseFee = tlb.Grams(7000000)

func (m *Market) calcUsersProfit(ctx context.Context, id uuid.UUID, winToken token.Token) (UserTotalReturnMap, error) {
	assetList, err := m.persistor.getEventAssets(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("calc users profit failed: %w", err)
	}

	var loseCollateral tlb.Grams
	var winCollateral tlb.Grams

	for _, asset := range assetList {
		if asset.Token == winToken {
			winCollateral += asset.CollateralStaked
		} else {
			loseCollateral += asset.CollateralStaked
		}
	}

	userTotalReturnMap := make(UserTotalReturnMap)

	for _, asset := range assetList {
		if asset.Token != winToken {
			continue
		}

		rest := float64(asset.CollateralStaked) / float64(winCollateral)
		profit := rest * float64(loseCollateral)
		userTotalReturn := profit + float64(asset.CollateralStaked)
		returnGrams := tlb.Grams(userTotalReturn) - baseFee

		log.Printf("[INFO] user_addr: %s, return_grams: %d\n\n", asset.UserRawAddress, returnGrams)

		if returnGrams < 0 {
			continue
		}
		userTotalReturnMap[asset.UserRawAddress] = returnGrams
	}

	return userTotalReturnMap, nil
}

var ErrCantSendSimpleTransfer = errors.New("can't send simple transfer")

func (m *Market) profitUsers(ctx context.Context, userTotalReturnMap UserTotalReturnMap, eventId uuid.UUID) error {
	for address, grams := range userTotalReturnMap {
		recepient := ton.MustParseAccountID(address)

		message := wallet.SimpleTransfer{
			Amount:     grams,
			Address:    recepient,
			Comment:    "profit: " + eventId.String(),
			Bounceable: false,
		}

		trySend := func() error {
			for i := 0; i < 20; i++ {
				if err := m.wallet.Send(ctx, message); err != nil {
					time.Sleep(10 * time.Second)
					continue
				}
				return nil
			}
			return ErrCantSendSimpleTransfer
		}

		log.Printf("[INFO] trying send user address: %s\n\n", address)
		if err := trySend(); err != nil {
			log.Printf("[ERROR] profit users, close event id: %s, send simple transfer failed for user: %s, has to get: %v err: %s\n\n",
				eventId.String(), recepient, grams, err.Error())
		}
	}
	return nil
}
