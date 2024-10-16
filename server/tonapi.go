package main

import (
	"context"
	"fmt"
	"github.com/TON-Market/tma/server/datatype"
	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/liteapi"
)

var networks = map[string]*liteapi.Client{}

func getAccountInfo(ctx context.Context, addr string, net *liteapi.Client) (*datatype.AccountInfo, error) {
	address, err := tongo.ParseAddress(addr)
	if err != nil {
		return nil, fmt.Errorf("parse address failed: %s", err.Error())
	}

	accountId := address.ID

	account, err := net.GetAccountState(ctx, accountId)
	if err != nil {
		return nil, fmt.Errorf("get account state failed: %s", err.Error())
	}

	accountInfo := datatype.AccountInfo{
		Balance: int64(account.Account.Account.Storage.Balance.Grams),
		Status:  string(account.Account.Status()),
	}

	accountInfo.Address.Raw = accountId.ToRaw()
	accountInfo.Address.Bounceable = accountId.ToHuman(true, false)
	accountInfo.Address.NonBounceable = accountId.ToHuman(false, false)

	return &accountInfo, nil
}
