package main

import (
	"context"
	"fmt"
	"github.com/TON-Market/tma/server/config"
	"github.com/TON-Market/tma/server/datatype/event"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tonconnect"
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

//const SEED = "example consider fiscal mail guitar tiger duck exhibit ancient series differ wealth mix kitchen cactus upgrade unable yellow impact confirm denial mesh during dove"
//const my_ton_keeper_addr = "UQBbRSVWRlRH0D_OJ2pzj_Kaoeo5_Q3F-6GhDayX044Xr1fU"
//
//func main() {
//	ctx := context.Background()
//
//	client, err := liteapi.NewClientWithDefaultMainnet()
//	if err != nil {
//		log.Fatalf("Unable to create lite client: %v", err)
//	}
//
//	pk, err := wallet.SeedToPrivateKey(SEED)
//	if err != nil {
//		log.Fatalln(err.Error())
//	}
//
//	if err != nil {
//		fmt.Println("Ошибка шифрования:", err)
//		return
//	}
//
//	fmt.Printf("Privte Key: %v\n", pk)
//
//	w, err := wallet.New(pk, wallet.HighLoadV2R2, client)
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	addrHuman := w.GetAddress().ToHuman(false, false)
//	fmt.Printf("Human address: %s\n", addrHuman)
//
//	balance, err := w.GetBalance(ctx)
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	fmt.Printf("Balance: %d\n", balance)
//
//	s, err := w.StateInit()
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	fmt.Println(s)
//
//	recepient, err := ton.AccountIDFromBase64Url(my_ton_keeper_addr)
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	m := wallet.SimpleTransfer{
//		Amount:     10000,
//		Address:    recepient,
//		Comment:    "Платеж из банка",
//		Bounceable: false,
//	}
//
//	err = w.Send(ctx, m)
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//}
