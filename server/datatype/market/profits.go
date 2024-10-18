package market

import (
	"context"
	"errors"
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

type UserReturn struct {
	Address string
	Amount  tlb.Grams
	LastTry time.Time
	EventID string
}

func (m *Market) CloseEvent(ctx context.Context, id uuid.UUID, winToken token.Token) error {
	if err := m.runtimer.close(ctx, id); err != nil {
		return fmt.Errorf("close event failed: %w", err)
	}
	time.Sleep(1 * time.Minute)
	userTotalReturnMap, err := m.calcUsersProfit(ctx, id, winToken)
	if err != nil {
		return fmt.Errorf("close event failed: %w", err)
	}

	go func() {
		for addr, grams := range userTotalReturnMap {
			ur := UserReturn{
				Address: addr,
				Amount:  grams,
				LastTry: time.Now().Add(-2 * time.Minute),
				EventID: id.String(),
			}

			m.profitCh <- ur
		}
	}()

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

func (m *Market) startSendProcess(ctx context.Context) {
	for ur := range m.profitCh {
		if ur.LastTry.Add(2 * time.Minute).After(time.Now()) {
			m.profitCh <- ur
			continue
		}

		accountID, err := ton.ParseAccountID(ur.Address)
		if err != nil {
			log.Printf("[ERROR] err check outcome transaction for user: %s\n\n", ur.Address)
			m.profitCh <- ur
			continue
		}

		getLastTransactions := func() ([]ton.Transaction, error) {
			for i := 0; i < 10; i++ {
				l, err := m.client.GetLastTransactions(ctx, accountID, 10)
				if err != nil {
					time.Sleep(5 * time.Second)
					continue
				}
				return l, nil
			}
			return nil, fmt.Errorf("transaction not found: %v", err)
		}

		log.Printf("\n[INFO] getting outcome transaction list for user: %s\n\n", ur.Address)
		trxList, err := getLastTransactions()
		if err != nil {
			log.Printf("[ERROR] get outcome trx list for user: %s, err: %s\n\n", ur.Address, err.Error())
			continue
		}

		if err = iterateUsersIncomeTransactionList(trxList, ur); err != nil && errors.Is(err, ErrUserIncomeTransactionNotFound) {
			m.profitCh <- ur
		}
	}
}

var ErrUserIncomeTransactionNotFound = errors.New("user income transaction not found")

func iterateUsersIncomeTransactionList(trxList []ton.Transaction, ur UserReturn) error {
	for _, trx := range trxList {
		var t wallet.TextComment
		if err := tlb.Unmarshal((*boc.Cell)(&trx.Msgs.InMsg.Value.Value.Body.Value), &t); err != nil {
			log.Printf("[WARNING] unmarshalling boc for outcome trx, user: %s, %s", ur.Address, err.Error())
			continue
		}

		idStr := wallet.TextComment("profit " + ur.EventID)

		if t == idStr {
			log.Printf("[SUCCESS] user: %s, got grams: %v\n\n", ur.Address, ur.Amount)
			return nil
		}
	}
	return ErrUserIncomeTransactionNotFound
}

func (m *Market) trySend(ctx context.Context, address, eventId string, grams tlb.Grams) error {
	r := ton.MustParseAccountID(address)

	message := wallet.SimpleTransfer{
		Amount:     grams,
		Address:    r,
		Comment:    "profit: " + eventId,
		Bounceable: true,
	}

	if err := m.wallet.Send(ctx, message); err != nil {
		return err
	}

	return nil
}
