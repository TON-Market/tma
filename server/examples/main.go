package main

//func main() {
//tongoClient, err := liteapi.NewClientWithDefaultMainnet()
//if err != nil {
//	fmt.Printf("Unable to create tongo client: %v", err)
//}
//
//accountId := tongo.MustParseAccountID("0:E2D41ED396A9F1BA03839D63C5650FAFC6FCFB574FD03F2E67D6555B61A3ACD9")
//
//accountId := tongo.MustParseAddress("EQBRW9rjhRUNL-Sy4swYbMzm2MgvlhC2DWIZFhYp2JnSoJaA")
//
//accountId := ton.MustParseAccountID("EQBRW9rjhRUNL-Sy4swYbMzm2MgvlhC2DWIZFhYp2JnSoJaA")
//
//trxs, err := tongoClient.GetLastTransactions(context.Background(), accountId, 100)
//if err != nil {
//	panic(err)
//}
//state, err := tongoClient.GetAccountState(context.Background(), accountId)
//if err != nil {
//	fmt.Printf("Get account state error: %v", err)
//}
//for _, trx := range trxs {
//trx.Msgs.InMsg.Value.Value.Body.Value.
//b, err := trx.Msgs.InMsg.Value.Value.Body.Value
//if err != nil {
//	panic(err)
//}
//fmt.Println(b)
//tlb.Any{}
//fmt.Println(trx.Msgs.InMsg.Value.Value.Body.Value)
//c := &boc.Cell{}
//err := trx.Msgs.InMsg.Value.Value.Body.Value.UnmarshalTLB(c, tlb.NewDecoder())
//if err != nil {
//	fmt.Println(err)
//}
//fmt.Println(trx)
//
//trx.Msgs.InMsg
//}
//
//err := boc.NewCell().WriteBytes([]byte(uuid.New().String()))
//if //err != nil {
//	return
//}
//fmt.Println(trxs)
//
//fmt.Printf("Account status: %v\nBalance: %v\n", state.Account.Status(), state.Account.Account.Storage.Balance.Grams)
//}

//const SEED = "example consider fiscal mail guitar tiger duck exhibit ancient series differ wealth mix kitchen cactus upgrade unable yellow impact confirm denial mesh during dove"
//const my_ton_keeper_addr = "UQBbRSVWRlRH0D_OJ2pzj_Kaoeo5_Q3F-6GhDayX044Xr1fU"
//const sultan_ton_addr = "UQCh7h4Yx13eloWxHvQrerioosDDlp8WYRDH3sj6U3FCtZW5"
//const gleb_addr = "UQCc34X8ziXMgJyhmpZqgj-xX92h8Hfn-bfmT3gOSaDaacTi"
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
//	w, err := wallet.New(pk, wallet.HighLoadV2R2, client)
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	recepient, err := ton.AccountIDFromBase64Url(gleb_addr)
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	m := wallet.SimpleTransfer{
//		Amount:     1000000000,
//		Address:    recepient,
//		Comment:    uuid.NewString(),
//		Bounceable: false,
//	}
//
//	err = w.Send(ctx, m)
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	body := boc.NewCell()
//	_ = tlb.Marshal(body, wallet.TextComment(uuid.NewString()))
//
//}

//func main() {
//	tongoClient, err := liteapi.NewClientWithDefaultMainnet()
//	if err != nil {
//		fmt.Printf("Unable to create tongo client: %v", err)
//	}
//
//	//accountId := ton.MustParseAccountID("0:5b452556465447d03fce276a738ff29aa1ea39fd0dc5fba1a10dac97d38e17af")
//
//	//accountId, err := ton.AccountIDFromRaw("0:5b452556465447d03fce276a738ff29aa1ea39fd0dc5fba1a10dac97d38e17af")
//	//if err != nil {
//	//	log.Fatalln(err)
//	//}
//
//	accountId := ton.MustParseAccountID("0:5b452556465447d03fce276a738ff29aa1ea39fd0dc5fba1a10dac97d38e17af")
//
//	trxs, err := tongoClient.GetLastTransactions(context.Background(), accountId, 100)
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	trx := trxs[0]
//
//	var t wallet.TextComment
//	if err := tlb.Unmarshal((*boc.Cell)(&trx.Msgs.OutMsgs.Values()[0].Value.Body.Value), &t); err != nil {
//		log.Fatalln(err)
//	}
//
//	fmt.Println(t)
//}
