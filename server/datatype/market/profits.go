package market

import (
	"context"
	"fmt"
	"github.com/TON-Market/tma/server/datatype/token"
	"github.com/google/uuid"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"github.com/tonkeeper/tongo/wallet"
	"log"
	"time"
)

type UserProfit struct {
	UserRawAddress string
	Grams          tlb.Grams
	LastTry        time.Time
	EventID        string
}

type TokenDeposits struct {
	WinCollateral  tlb.Grams
	LoseCollateral tlb.Grams
}

func (m *Market) CloseEvent(ctx context.Context, eventID uuid.UUID, winToken token.Token) error {
	if err := m.runtimer.close(ctx, eventID); err != nil {
		return fmt.Errorf("close event failed: %w", err)
	}
	//time.Sleep(5 * time.Minute)
	time.Sleep(20 * time.Second)
	userProfitList, err := m.buildUserProfitData(ctx, eventID, winToken)
	if err != nil {
		return fmt.Errorf("close event failed: %w", err)
	}

	m.broadcastUsersProfit(ctx, userProfitList)
	m.writeToChannel(ctx, userProfitList)
	return nil
}

func (m *Market) buildUserProfitData(ctx context.Context, eventID uuid.UUID, winToken token.Token) ([]*UserProfit, error) {
	assetList, err := m.persistor.getEventAssets(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("calc user profit failed: %w", err)
	}

	tokenDeposits := m.getTokenDeposits(assetList, winToken)

	userProfitList := make([]*UserProfit, 0)

	for _, asset := range assetList {
		if asset.Token != winToken {
			continue
		}

		profit := m.calcUserProfit(ctx, asset, tokenDeposits)
		if profit < 0 {
			continue
		}

		userProfit := &UserProfit{
			UserRawAddress: asset.UserRawAddress,
			Grams:          profit,
			EventID:        eventID.String(),
		}

		userProfitList = append(userProfitList, userProfit)
	}

	return userProfitList, nil
}

func (m *Market) getTokenDeposits(assetList []*Asset, winToken token.Token) *TokenDeposits {
	var loseCollateral tlb.Grams
	var winCollateral tlb.Grams

	for _, asset := range assetList {
		if asset.Token == winToken {
			winCollateral += asset.CollateralStaked
		} else {
			loseCollateral += asset.CollateralStaked
		}
	}

	return &TokenDeposits{
		WinCollateral:  winCollateral,
		LoseCollateral: loseCollateral,
	}
}

var baseFee = tlb.Grams(7000000)

func (m *Market) calcUserProfit(_ context.Context, asset *Asset, tokenDeposits *TokenDeposits) tlb.Grams {
	rest := float64(asset.CollateralStaked) / float64(tokenDeposits.WinCollateral)
	profit := rest * float64(tokenDeposits.LoseCollateral)
	userTotalReturn := profit + float64(asset.CollateralStaked)
	returnGrams := tlb.Grams(userTotalReturn) - baseFee
	return returnGrams
}

func (m *Market) broadcastUsersProfit(ctx context.Context, userProfitList []*UserProfit) {
	for _, userProfit := range userProfitList {
		if err := m.trySendProfit(ctx, userProfit); err != nil {
			log.Printf("[ERROR] try send profit for user: %s, grams: %v failed: %s\n\n",
				userProfit.UserRawAddress, userProfit.Grams, err.Error())
			continue
		}
	}
}

func (m *Market) writeToChannel(_ context.Context, userProfitList []*UserProfit) {
	for _, userProfit := range userProfitList {
		m.profitCh <- userProfit
	}
}

func (m *Market) startResendProcess(ctx context.Context) {
	go func() {
		for userProfit := range m.profitCh {
			if userProfit.LastTry.Add(3 * time.Minute).After(time.Now()) {
				m.profitCh <- userProfit
				continue
			}

			isProfitDelivered, err := m.checkIfProfitDelivered(ctx, userProfit)
			if err != nil {
				log.Printf("[ERROR] check profit delivered failed: %s\n\n", err.Error())
				continue
			}

			if !isProfitDelivered {
				log.Printf("[WARNING] profit for user: %s, grams: %v still pending\n\n",
					userProfit.UserRawAddress, userProfit.Grams)

				if err = m.trySendProfit(ctx, userProfit); err != nil {
					log.Printf("[ERROR] try send profit for user: %s, grams: %v failed: %s\n\n",
						userProfit.UserRawAddress, userProfit.Grams, err.Error())
				}
				m.profitCh <- userProfit
			}
			log.Printf("[SUCCESS] profit for user: %s, grams: %v delivered\n\n",
				userProfit.UserRawAddress, userProfit.Grams)
		}
	}()
}

func (m *Market) checkIfProfitDelivered(ctx context.Context, userProfit *UserProfit) (bool, error) {
	trxList, err := m.getLastTransactions(ctx, userProfit.UserRawAddress)
	if err != nil {
		return false, fmt.Errorf("check if profit delivered failed: %w", err)
	}

	for _, trx := range trxList {
		isProfitTransaction, err := m.isProfitTransaction(ctx, trx, userProfit)
		if err != nil {
			continue
		}
		if isProfitTransaction {
			return true, nil
		}
	}
	return false, nil
}

func (m *Market) isProfitTransaction(_ context.Context, trx ton.Transaction, userProfit *UserProfit) (bool, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("[ALARM] Recovered in isProfitTransaction", r)
		}
	}()
	var t wallet.TextComment
	if err := tlb.Unmarshal((*boc.Cell)(&trx.Msgs.InMsg.Value.Value.Body.Value), &t); err != nil {
		return false, fmt.Errorf("unmarshal transaction failed: %w", err)
	}

	comment := wallet.TextComment("event closed: " + userProfit.EventID)

	if t == comment {
		return true, nil
	}

	return false, nil
}

func (m *Market) trySendProfit(ctx context.Context, userProfit *UserProfit) error {
	recipient := ton.MustParseAccountID(userProfit.UserRawAddress)

	message := wallet.SimpleTransfer{
		Amount:     userProfit.Grams,
		Address:    recipient,
		Comment:    "event closed: " + userProfit.EventID,
		Bounceable: true,
	}

	userProfit.LastTry = time.Now()

	if err := m.wallet.Send(ctx, message); err != nil {
		return err
	}

	return nil
}
