package market

import (
	"context"
	"fmt"
	"github.com/TON-Market/tma/server/datatype/token"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"github.com/tonkeeper/tongo/wallet"
	"time"
)

func (m *Market) CloseEvent(ctx context.Context, id uuid.UUID, winToken token.Token) error {
	if err := m.runtimer.close(ctx, id); err != nil {
		return fmt.Errorf("close event failed: %w", err)
	}
	time.Sleep(1 * time.Minute)
	userTotalReturnMap, err := m.calcUsersProfit(ctx, id, winToken)
	if err != nil {
		return fmt.Errorf("close event failed: %w", err)
	}

	if err = m.profitUsers(ctx, userTotalReturnMap, id); err != nil {
		return fmt.Errorf("close event failed: %w", err)
	}
	return nil
}

type UserTotalReturnMap map[string]tlb.Grams

var baseFee = tlb.Grams(5000000)

func (m *Market) calcUsersProfit(ctx context.Context, id uuid.UUID, winToken token.Token) (UserTotalReturnMap, error) {
	assetList, err := m.persistor.getAssets(ctx, id)
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
		if returnGrams < 0 {
			continue
		}
		userTotalReturnMap[asset.UserRawAddress] = returnGrams
	}

	return userTotalReturnMap, nil
}

func (m *Market) profitUsers(ctx context.Context, userTotalReturnMap UserTotalReturnMap, eventId uuid.UUID) error {
	for address, grams := range userTotalReturnMap {
		recepient, err := ton.ParseAccountID(address)
		if err != nil {
			log.Printf("[ERROR] profit users, close event id: %s, failed parse account id for user: %s, err: %s\n", eventId.String(), recepient, err.Error())
			continue
		}

		message := wallet.SimpleTransfer{
			Amount:     grams,
			Address:    recepient,
			Comment:    "profit: " + eventId.String(),
			Bounceable: false,
		}

		if err = m.wallet.Send(ctx, message); err != nil {
			log.Printf("[ERROR] profit users, close event id: %s, send simple transfer failed for user: %s, err: %s\n", eventId.String(), recepient, err.Error())
		}
	}
	return nil
}
