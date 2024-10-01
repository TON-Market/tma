package main

import (
	"context"
	"github.com/TON-Market/tma/server/datatype"

	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/liteapi"
)

var networks = map[string]*liteapi.Client{}

func GetAccountInfo(ctx context.Context, address tongo.AccountID, net *liteapi.Client) (*datatype.AccountInfo, error) {
	account, err := net.GetAccountState(ctx, address)
	if err != nil {
		return nil, err
	}
	accountInfo := datatype.AccountInfo{
		Balance: int64(account.Account.Account.Storage.Balance.Grams),
		Status:  string(account.Account.Status()),
	}
	accountInfo.Address.Raw = address.ToRaw()
	accountInfo.Address.Bounceable = address.ToHuman(true, false)
	accountInfo.Address.NonBounceable = address.ToHuman(false, false)

	return &accountInfo, nil
}
