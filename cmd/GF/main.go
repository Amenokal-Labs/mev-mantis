package main

import (
	"context"
	"fmt"

	"github.com/Amenokal-Labs/mev-mantis.git/pkg/utils"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/" + utils.GetAPIKey("INFURA_KEY"))
	if err != nil {
		panic(err)
	}
	_ = client
	fmt.Println("client connected")

	test, err := client.PendingTransactionCount(context.Background())
	fmt.Println(test)
}
