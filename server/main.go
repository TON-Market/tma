package main

import (
	"context"
	"fmt"
	"github.com/TON-Market/tma/server/config"
	"github.com/TON-Market/tma/server/datatype/event"
	_ "net/http/pprof"

	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tonconnect"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Tonproof is running")
	config.LoadConfig()

	e := echo.New()
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		Skipper:           nil,
		DisableStackAll:   true,
		DisablePrintStack: false,
	}))
	e.Use(middleware.Logger())
	e.Static("/", "./")

	mainNetClient, err := liteapi.NewClientWithDefaultMainnet()
	if err != nil {
		log.Fatalf("failed init mainnet liteapi client")
	}
	networks["-239"] = mainNetClient

	testNetClient, err := liteapi.NewClientWithDefaultTestnet()
	if err != nil {
		log.Fatalf("failed init testnet liteapi client")
	}
	networks["-3"] = testNetClient

	payloadLifeTime := config.Config.Proof.PayloadLifeTimeSec
	proofLifeTime := config.Config.Proof.ProofLifeTimeSec
	tonConnectMainNet, err := tonconnect.NewTonConnect(mainNetClient, config.Config.Proof.PayloadSignatureKey,
		tonconnect.WithLifeTimePayload(payloadLifeTime), tonconnect.WithLifeTimeProof(proofLifeTime))
	tonConnectTestNet, err := tonconnect.NewTonConnect(testNetClient, config.Config.Proof.PayloadSignatureKey,
		tonconnect.WithLifeTimePayload(payloadLifeTime), tonconnect.WithLifeTimeProof(proofLifeTime))

	h := newHandler(tonConnectMainNet, tonConnectTestNet)
	w := newSocket()

	event.Keeper().Start(context.Background())

	registerHandlers(e, h, w)

	log.Fatal(e.Start(fmt.Sprintf(":%v", config.Config.Port)))
}
