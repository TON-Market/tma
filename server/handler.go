package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TON-Market/tma/server/datatype"
	"github.com/TON-Market/tma/server/datatype/market"
	"github.com/TON-Market/tma/server/datatype/user"
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
	lg := log.WithContext(c.Request().Context()).WithField("prefix", "ProofHandler")

	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(HttpResErrorWithLog(err.Error(), http.StatusBadRequest, lg))
	}
	var tp datatype.TonProof
	err = json.Unmarshal(b, &tp)
	if err != nil {
		return c.JSON(HttpResErrorWithLog(err.Error(), http.StatusBadRequest, lg))
	}

	tonConnect := h.tonConnectMainNet

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
			lg.Errorln(err.Error())
		}
		return c.JSON(HttpResErrorWithLog("proof verification failed", http.StatusBadRequest, lg))
	}

	claims := &jwtCustomClaims{
		tp.Address,
		jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(10, 0, 0).Unix(),
		},
	}

	jwtTkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := jwtTkn.SignedString([]byte(h.tonConnectMainNet.GetSecret()))
	if err != nil {
		return err
	}

	ctx := context.TODO()

	info, err := getAccountInfo(ctx, tp.Address, networks["-239"])
	if err != nil {
		return c.JSON(HttpResErrorWithLog(fmt.Errorf("get account info failed: %s", err.Error()).Error(), http.StatusInternalServerError, lg))
	}

	if _, err = user.UserStorage().Get(ctx, info.Address.Raw); err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			if err = user.UserStorage().Save(ctx, &user.User{
				RawAddr:  info.Address.Raw,
				DealList: make([]*market.Deal, 0),
			}); err != nil {
				lg.Println(err)
				return c.JSON(HttpResErrorWithLog(fmt.Errorf("save user failed: %v", err).Error(), http.StatusInternalServerError, lg))
			}
		} else {
			lg.Println(err)
			return c.JSON(HttpResErrorWithLog(fmt.Errorf("check user exists failed: %v", err).Error(), http.StatusInternalServerError, lg))
		}
	}

	cookie := new(http.Cookie)
	cookie.Name = "AuthToken"
	cookie.Value = signedToken
	cookie.Expires = time.Now().Add(24 * 365 * time.Hour)
	cookie.Path = "/"
	cookie.HttpOnly = false
	cookie.Secure = false
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, echo.Map{
		"token": signedToken,
	})
}

func (h *handler) PayloadHandler(c echo.Context) error {
	lg := log.WithContext(c.Request().Context()).WithField("prefix", "PayloadHandler")

	payload, err := h.tonConnectMainNet.GeneratePayload()
	if err != nil {
		return c.JSON(HttpResErrorWithLog(err.Error(), http.StatusBadRequest, lg))
	}

	return c.JSON(http.StatusOK, echo.Map{
		"payload": payload,
	})
}

func (h *handler) GetAccountInfo(c echo.Context) error {
	ctx := c.Request().Context()
	lg := log.WithContext(ctx).WithField("prefix", "getAccountInfo")

	cookie, err := c.Cookie("AuthToken")
	if err != nil {
		return c.JSON(HttpResErrorWithLog(fmt.Sprintf("can't get cookie: %v", err), http.StatusBadRequest, lg))
	}

	signedToken := cookie.Value

	jwtToken, err := jwt.ParseWithClaims(signedToken, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.tonConnectMainNet.GetSecret()), nil
	})
	if err != nil {
		return c.JSON(HttpResErrorWithLog(fmt.Sprintf("invalid jwtToken: %v", err), http.StatusUnauthorized, lg))
	}

	if claims, ok := jwtToken.Claims.(*jwtCustomClaims); ok && jwtToken.Valid {
		if time.Unix(claims.StandardClaims.ExpiresAt, 0).Before(time.Now()) {
			return c.JSON(HttpResErrorWithLog("jwtToken has expired", http.StatusUnauthorized, lg))
		}

		net := networks["-239"]

		info, err := getAccountInfo(ctx, claims.Address, net)
		if err != nil {
			return c.JSON(HttpResErrorWithLog(fmt.Sprintf("get account info error: %v", err), http.StatusBadRequest, lg))
		}

		return c.JSON(http.StatusOK, info)
	} else {
		return c.JSON(HttpResErrorWithLog("invalid jwtToken claims", http.StatusUnauthorized, lg))
	}
}

func (h *handler) validateUser(auth string, c echo.Context) (bool, error) {
	lg := log.WithContext(context.Background()).WithField("prefix", "auth request")
	jwtTkn, err := jwt.ParseWithClaims(auth, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.tonConnectMainNet.GetSecret()), nil
	})
	if err != nil {
		return false, c.JSON(HttpResErrorWithLog("jwtTkn has expired", http.StatusUnauthorized, lg))
	}

	if claims, ok := jwtTkn.Claims.(*jwtCustomClaims); ok && jwtTkn.Valid {
		if time.Unix(claims.StandardClaims.ExpiresAt, 0).Before(time.Now()) {
			return false, c.JSON(HttpResErrorWithLog("jwtTkn has expired", http.StatusUnauthorized, lg))
		}
		c.Set("address", claims.Address)
		return true, nil
	} else {
		return false, c.JSON(HttpResErrorWithLog("invalid jwtTkn claims", http.StatusUnauthorized, lg))
	}
}

type GetEventsResponse struct {
	Items []market.EventDTO `json:"items"`
	Pages int               `json:"pages"`
}

func (h *handler) GetEvents(c echo.Context) error {
	lg := log.WithContext(context.Background()).WithField("prefix", "GetEvents")

	ctx := context.TODO()

	pageInput := c.QueryParam("page")
	page, err := strconv.Atoi(pageInput)
	if err != nil {
		return c.JSON(HttpResErrorWithLog("incorrect page passed", http.StatusBadRequest, lg))
	}

	tagInput := c.QueryParam("tag")
	tag, err := strconv.Atoi(tagInput)
	if err != nil {
		return c.JSON(HttpResErrorWithLog("incorrect tag passed", http.StatusBadRequest, lg))
	}

	list, totalPages, _ := market.GetMarket().ReadFromSnapshot(ctx, market.Tag(tag), page)

	getEventsResponse := &GetEventsResponse{
		Items: list,
		Pages: totalPages,
	}

	return c.JSON(http.StatusOK, getEventsResponse)
}

type Tag struct {
	ID    market.Tag `json:"id"`
	Title string     `json:"title"`
}

func (h *handler) GetTags(c echo.Context) error {
	tagList := []*Tag{
		{
			ID:    market.No,
			Title: "No",
		},
		{
			ID:    market.Politic,
			Title: "Politics",
		},
		{
			ID:    market.Economics,
			Title: "Economics",
		},
		{
			ID:    market.Crypto,
			Title: "Crypto",
		},
		{
			ID:    market.Culture,
			Title: "Culture",
		},

		{
			ID:    market.Other,
			Title: "Other",
		},
	}

	return c.JSON(http.StatusOK, tagList)
}
