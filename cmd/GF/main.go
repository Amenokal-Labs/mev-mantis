package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Amenokal-Labs/mev-mantis.git/pkg/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func main() {
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/" + utils.GetAPIKey("INFURA_KEY"))
	if err != nil {
		panic(err)
	}
	fmt.Println("client connected")

	test, err := client.PendingTransactionCount(context.Background())
	fmt.Println(test)

	rpc, err := rpc.Dial("https://mainnet.infura.io/ws/v3/" + utils.GetAPIKey("INFURA_KEY"))
	client2 := gethclient.New(rpc)

	hashes := make(chan common.Hash)
	_, err = client2.SubscribePendingTransactions(context.Background(), hashes)

	if err != nil {
		log.Fatal("failed to subscribe", err)
	}

	select {
	case hash := <-hashes:
		fmt.Println("")
		_ = hash
	}
}
