package main

import (
	"context"
	"crypto/tls"
	"github.com/TON-Market/tma/server/config"
	"github.com/TON-Market/tma/server/datatype/market"
	"github.com/TON-Market/tma/server/datatype/token"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tonconnect"
	"net/http"
	_ "net/http/pprof"
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
	e.Static("/", "./static")

	mainNetClient, err := liteapi.NewClientWithDefaultMainnet()
	if err != nil {
		log.Fatalf("failed init mainnet liteapi client")
	}
	networks["-239"] = mainNetClient

	payloadLifeTime := config.Config.Proof.PayloadLifeTimeSec
	proofLifeTime := config.Config.Proof.ProofLifeTimeSec
	tonConnectMainNet, err := tonconnect.NewTonConnect(mainNetClient, config.Config.Proof.PayloadSignatureKey,
		tonconnect.WithLifeTimePayload(payloadLifeTime), tonconnect.WithLifeTimeProof(proofLifeTime))

	h := newHandler(tonConnectMainNet)
	w := newSocket()

	market.GetMarket().Start(context.TODO())
	registerHandlers(e, h, w)

	testData()

	s := http.Server{
		Addr:    ":8443",
		Handler: e,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	log.Fatal(s.ListenAndServeTLS("./server.crt", "./server.key"))
}

func testData() {
	e := &market.Event{
		Tag:      market.Crypto,
		LogoLink: "",
		Title:    "Will Ton Market win hackathon?",
		BetMap: map[token.Token]*market.Bet{
			token.A: {
				Token:    token.A,
				Title:    "Yes",
				LogoLink: "",
			},
			token.B: {
				Token:    token.B,
				Title:    "No",
				LogoLink: "",
			},
		},
	}
	elections := &market.Event{
		Tag:      market.Politic,
		LogoLink: "https://www.cft.org/sites/main/files/imagecache/medium/main-images/elections_2024_2_0.png?1718390367",
		Title:    "USA Elections 2024",
		BetMap: map[token.Token]*market.Bet{
			token.A: {
				Token:    token.A,
				Title:    "Trump",
				LogoLink: "https://upload.wikimedia.org/wikipedia/commons/5/56/Donald_Trump_official_portrait.jpg",
			},
			token.B: {
				Token:    token.B,
				Title:    "Harris",
				LogoLink: "https://www.whitehouse.gov/wp-content/uploads/2021/04/V20210305LJ-0043.jpg",
			},
		},
	}
	err := market.GetMarket().AddEvent(context.Background(), e)
	if err != nil {
		log.Fatalln(err)
	}
	err = market.GetMarket().AddEvent(context.Background(), elections)
	if err != nil {
		log.Fatalln(err)
	}
}
