package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/TON-Market/tma/server/datatype"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/liteapi"
	"time"
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

func (h *handler) getAccountInfoFromCookie(ctx context.Context, c echo.Context) (*datatype.AccountInfo, error) {
	cookie, err := c.Cookie("AuthToken")
	if err != nil {
		return nil, err
	}

	signedToken := cookie.Value

	jwtToken, err := jwt.ParseWithClaims(signedToken, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.tonConnectMainNet.GetSecret()), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := jwtToken.Claims.(*jwtCustomClaims); ok && jwtToken.Valid {
		if time.Unix(claims.StandardClaims.ExpiresAt, 0).Before(time.Now()) {
			return nil, errors.New("jwt token expired")
		}

		net := networks["-239"]

		info, err := getAccountInfo(ctx, claims.Address, net)
		if err != nil {
			return nil, err
		}

		return info, nil
	}
	return nil, errors.New("invalid jwt claims")
}
