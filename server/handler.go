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
}

func newHandler(tonConnectMainNet *tonconnect.Server) *handler {
	h := handler{
		tonConnectMainNet: tonConnectMainNet,
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
	cookie.HttpOnly = true
	cookie.Secure = true
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

func (h *handler) Disconnect(c echo.Context) error {
	addr := c.Get("address").(string)
	log.Printf("[INFO] user: %s disconnected\n", addr)

	cookie := new(http.Cookie)
	cookie.Name = "AuthToken"
	cookie.Value = "null"
	cookie.Expires = time.Now().Add(-time.Hour)
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = true
	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, "ok")
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
			ID:    market.All,
			Title: "All",
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

type GetAssetsResp struct {
	AssetList     []*market.AssetDTO `json:"assetList"`
	TotalInMarket string             `json:"totalInMarket"`
}

func (h *handler) GetAssets(c echo.Context) error {
	ctx := context.TODO()
	lg := log.WithContext(ctx).WithField("prefix", "GetAssets")

	addr := c.Get("address").(string)

	assetDtoList, totalInMarket, err := market.GetMarket().GetUserAssets(ctx, addr)
	if err != nil {
		return c.JSON(HttpResErrorWithLog(err.Error(), http.StatusInternalServerError, lg))
	}

	resp := &GetAssetsResp{
		AssetList:     assetDtoList,
		TotalInMarket: totalInMarket,
	}

	return c.JSON(http.StatusOK, resp)
}
