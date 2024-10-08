package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/TON-Market/tma/server/datatype"
	"github.com/TON-Market/tma/server/datatype/event"
	"github.com/TON-Market/tma/server/datatype/token"
	"github.com/google/uuid"
	"github.com/tonkeeper/tongo"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/tonkeeper/tongo/tonconnect"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type jwtCustomClaims struct {
	Address string `json:"address"`
	jwt.StandardClaims
}

type handler struct {
	tonConnectMainNet *tonconnect.Server
	tonConnectTestNet *tonconnect.Server
}

func newHandler(tonConnectMainNet, tonConnectTestNet *tonconnect.Server) *handler {
	h := handler{
		tonConnectMainNet: tonConnectMainNet,
		tonConnectTestNet: tonConnectTestNet,
	}
	return &h
}

func (h *handler) ProofHandler(c echo.Context) error {
	log := log.WithContext(c.Request().Context()).WithField("prefix", "ProofHandler")

	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(HttpResErrorWithLog(err.Error(), http.StatusBadRequest, log))
	}
	var tp datatype.TonProof
	err = json.Unmarshal(b, &tp)
	if err != nil {
		return c.JSON(HttpResErrorWithLog(err.Error(), http.StatusBadRequest, log))
	}

	var tonConnect *tonconnect.Server
	switch tp.Network {
	case "-239":
		tonConnect = h.tonConnectMainNet
	case "-3":
		tonConnect = h.tonConnectTestNet
	default:
		return c.JSON(HttpResErrorWithLog(fmt.Sprintf("undefined network: %v", tp.Network), http.StatusBadRequest, log))
	}
	proof := tonconnect.Proof{
		Address: tp.Address,
		Proof: tonconnect.ProofData{
			Timestamp: tp.Proof.Timestamp,
			Domain:    tp.Proof.Domain.Value,
			Signature: tp.Proof.Signature,
			Payload:   tp.Proof.Payload,
			StateInit: tp.Proof.StateInit,
		},
	}
	verified, _, err := tonConnect.CheckProof(context.Background(), &proof,
		h.tonConnectMainNet.CheckPayload, func(string) (bool, error) {
			return true, nil
		})
	if err != nil || !verified {
		if err != nil {
			log.Errorln(err.Error())
		}
		return c.JSON(HttpResErrorWithLog("proof verification failed", http.StatusBadRequest, log))
	}

	claims := &jwtCustomClaims{
		tp.Address,
		jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(10, 0, 0).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(h.tonConnectMainNet.GetSecret()))
	if err != nil {
		return err
	}

	cookie := new(http.Cookie)
	cookie.Name = "AuthToken"
	cookie.Value = signedToken
	cookie.Expires = time.Now().Add(24 * 365 * time.Hour)
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = false
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, echo.Map{
		"token": signedToken,
	})
}

func (h *handler) PayloadHandler(c echo.Context) error {
	log := log.WithContext(c.Request().Context()).WithField("prefix", "PayloadHandler")

	payload, err := h.tonConnectMainNet.GeneratePayload()
	if err != nil {
		return c.JSON(HttpResErrorWithLog(err.Error(), http.StatusBadRequest, log))
	}

	return c.JSON(http.StatusOK, echo.Map{
		"payload": payload,
	})
}

func (h *handler) GetAccountInfo(c echo.Context) error {
	ctx := c.Request().Context()
	log := log.WithContext(ctx).WithField("prefix", "GetAccountInfo")

	// Получаем куки
	cookie, err := c.Cookie("AuthToken")
	if err != nil {
		return c.JSON(HttpResErrorWithLog(fmt.Sprintf("can't get cookie: %v", err), http.StatusBadRequest, log))
	}

	signedToken := cookie.Value

	token, err := jwt.ParseWithClaims(signedToken, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.tonConnectMainNet.GetSecret()), nil
	})
	if err != nil {
		return c.JSON(HttpResErrorWithLog(fmt.Sprintf("invalid token: %v", err), http.StatusUnauthorized, log))
	}

	if claims, ok := token.Claims.(*jwtCustomClaims); ok && token.Valid {
		if time.Unix(claims.StandardClaims.ExpiresAt, 0).Before(time.Now()) {
			return c.JSON(HttpResErrorWithLog("token has expired", http.StatusUnauthorized, log))
		}

		net := networks["-239"]

		addr, err := tongo.ParseAddress(claims.Address)
		if err != nil {
			log.Fatalln(err)
		}

		info, err := GetAccountInfo(ctx, addr.ID, net)
		if err != nil {
			return c.JSON(HttpResErrorWithLog(fmt.Sprintf("get account info error: %v", err), http.StatusBadRequest, log))
		}

		return c.JSON(http.StatusOK, info)
	} else {
		return c.JSON(HttpResErrorWithLog("invalid token claims", http.StatusUnauthorized, log))
	}
}

func (h *handler) validateUser(auth string, c echo.Context) (bool, error) {
	log := log.WithContext(context.Background()).WithField("prefix", "auth request")
	token, err := jwt.ParseWithClaims(auth, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.tonConnectMainNet.GetSecret()), nil
	})
	if err != nil {
		return false, c.JSON(HttpResErrorWithLog("token has expired", http.StatusUnauthorized, log))
	}

	if claims, ok := token.Claims.(*jwtCustomClaims); ok && token.Valid {
		if time.Unix(claims.StandardClaims.ExpiresAt, 0).Before(time.Now()) {
			return false, c.JSON(HttpResErrorWithLog("token has expired", http.StatusUnauthorized, log))
		}
		return true, nil
	} else {
		return false, c.JSON(HttpResErrorWithLog("invalid token claims", http.StatusUnauthorized, log))
	}
}

type GetEventsResponse struct {
	Items []*event.EventDTO `json:"items"`
	Pages int               `json:"pages"`
}

func (h *handler) GetEvents(c echo.Context) error {
	log := log.WithContext(context.Background()).WithField("prefix", "GetEvents")

	pageInput := c.QueryParam("page")
	page, err := strconv.Atoi(pageInput)
	if err != nil {
		return c.JSON(HttpResErrorWithLog("incorrect page passed", http.StatusBadRequest, log))
	}

	tagInput := c.QueryParam("tag")
	tag, err := strconv.Atoi(tagInput)
	if err != nil {
		return c.JSON(HttpResErrorWithLog("incorrect tag passed", http.StatusBadRequest, log))
	}

	list, totalPages, err := event.Keeper().GetSnapshot(context.Background(), event.EventTag(tag), page)
	if err != nil {
		return c.JSON(HttpResErrorWithLog(fmt.Sprintf("internal server error: %s", err.Error()), http.StatusInternalServerError, log))
	}

	getEventsResponse := &GetEventsResponse{
		Items: list,
		Pages: totalPages,
	}

	return c.JSON(http.StatusOK, getEventsResponse)
}

// AddDeposit - тестовая функция добавления депозита
func (h *handler) AddDeposit(c echo.Context) error {
	log := log.WithContext(context.Background()).WithField("prefix", "GetEvents")

	eventIDInput := c.QueryParam("eventId")

	u, err := uuid.Parse(eventIDInput)
	if err != nil {
		return c.JSON(HttpResErrorWithLog("incorrect event id passed", http.StatusBadRequest, log))
	}

	if err := event.Keeper().AddDeposit(context.Background(), u, 100, token.A); err != nil {
		return c.JSON(HttpResErrorWithLog(fmt.Sprintf("internal server error: %s", err.Error()), http.StatusInternalServerError, log))
	}

	return c.JSON(http.StatusOK, "ok")
}

type Tag struct {
	ID    event.EventTag `json:"id"`
	Title string         `json:"title"`
}

func (h *handler) getTags(c echo.Context) error {
	tagList := []*Tag{
		{
			ID:    event.No,
			Title: "No",
		},
		{
			ID:    event.Politic,
			Title: "Politics",
		},
		{
			ID:    event.Economics,
			Title: "Economics",
		},
		{
			ID:    event.Crypto,
			Title: "Crypto",
		},
		{
			ID:    event.Culture,
			Title: "Culture",
		},

		{
			ID:    event.Other,
			Title: "Other",
		},
	}

	return c.JSON(http.StatusOK, tagList)
}
