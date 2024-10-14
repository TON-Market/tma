package main

import (
	"context"
	"encoding/json"
	"github.com/TON-Market/tma/server/datatype/market"
	"github.com/TON-Market/tma/server/datatype/token"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/wallet"
	"io"
	"net/http"
)

type Message struct {
	Addr    string `json:"address"`
	Amount  string `json:"amount"`
	Payload []byte `json:"payload"`
}

type PayResp struct {
	Message   *Message `json:"message"`
	DepositID string   `json:"depositID"`
}

type PayReq struct {
	EventID    string      `json:"eventID"`
	Collateral float64     `json:"collateral"`
	Token      token.Token `json:"token"`
}

func (h *handler) Pay(c echo.Context) error {
	ctx := context.TODO()
	lg := log.WithContext(ctx).WithField("prefix", "Pay")

	addr := c.Get("address").(string)
	if addr == "" {
		return c.JSON(HttpResErrorWithLog("address is empty", http.StatusUnauthorized, lg))
	}

	log.Printf("[INFO] client addr: %s\n", addr)

	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(HttpResErrorWithLog(err.Error(), http.StatusBadRequest, lg))
	}

	var payReq PayReq
	if err := json.Unmarshal(b, &payReq); err != nil {
		return c.JSON(HttpResErrorWithLog(err.Error(), http.StatusBadRequest, lg))
	}

	dealId := uuid.New()
	eventId, err := uuid.Parse(payReq.EventID)
	if err != nil {
		return c.JSON(HttpResErrorWithLog(err.Error(), http.StatusBadRequest, lg))
	}
	grams := FloatToGrams(payReq.Collateral)

	d := &market.Deal{
		ID:          dealId,
		EventID:     eventId,
		Token:       payReq.Token,
		Collateral:  grams,
		UserRawAddr: addr,
		Size:        0,
	}

	if err := market.GetMarket().SaveDealUnchecked(ctx, d); err != nil {
		return c.JSON(HttpResErrorWithLog(err.Error(), http.StatusInternalServerError, lg))
	}

	body := boc.NewCell()
	if err := tlb.Marshal(body, wallet.TextComment(dealId.String())); err != nil {
		return c.JSON(HttpResErrorWithLog(err.Error(), http.StatusInternalServerError, lg))
	}

	payload, err := body.ToBoc()
	if err != nil {
		return c.JSON(HttpResErrorWithLog(err.Error(), http.StatusInternalServerError, lg))
	}

	gramsStr := GramsToString(grams)

	payResp := &PayResp{
		Message: &Message{
			Addr:    market.BANK_ADDR,
			Amount:  gramsStr,
			Payload: payload,
		},
		DepositID: dealId.String(),
	}

	return c.JSON(http.StatusOK, payResp)
}

type DepositReq struct {
	DepositStatus market.DepositStatus `json:"depositStatus"`
	DepositID     string               `json:"depositID"`
}

func (h *handler) Deposit(c echo.Context) error {
	ctx := context.TODO()
	lg := log.WithContext(ctx).WithField("prefix", "Deposit")

	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(HttpResErrorWithLog(err.Error(), http.StatusBadRequest, lg))
	}

	var depositReq DepositReq
	if err := json.Unmarshal(b, &depositReq); err != nil {
		return c.JSON(HttpResErrorWithLog(err.Error(), http.StatusBadRequest, lg))
	}

	depositUid, err := uuid.Parse(depositReq.DepositID)
	if err != nil {
		return c.JSON(HttpResErrorWithLog(err.Error(), http.StatusBadRequest, lg))
	}

	dr := &market.DepositReq{
		ID:            depositUid,
		DepositStatus: depositReq.DepositStatus,
	}
	if err := market.GetMarket().Deposit(ctx, dr); err != nil {
		return c.JSON(HttpResErrorWithLog(err.Error(), http.StatusInternalServerError, lg))
	}

	return c.JSON(http.StatusOK, "ok")
}
