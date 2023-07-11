package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Amenokal-Labs/mev-mantis.git/pkg/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func main() {
	httpsClient, err := ethclient.Dial("https://mainnet.infura.io/v3/" + utils.GetAPIKey("INFURA_KEY"))
	if err != nil {
		panic(err)
	}

	rpc, err := rpc.Dial("wss://mainnet.infura.io/ws/v3/" + utils.GetAPIKey("INFURA_KEY"))
	rpcClient := gethclient.New(rpc)
	fmt.Println("client connected")

	// txs := make(chan *types.Transaction)
	// _, err = rpcClient.SubscribeFullPendingTransactions(context.Background(), txs)
	// if err != nil {
	// 	panic(err)
	// }

	hashes := make(chan common.Hash)
	_, err = rpcClient.SubscribePendingTransactions(context.Background(), hashes)
	if err != nil {
		panic(err)
	}

	for {
		pendingCount, _ := httpsClient.PendingTransactionCount(context.Background())
		fmt.Println("Pending count:", pendingCount)

		hash := <-hashes
		fmt.Println(hash)
		tx, _, _ := httpsClient.TransactionByHash(context.Background(), hash)

		// tx, isPending, err := httpsClient.TransactionByHash(context.Background(), hash)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// fmt.Println(isPending)
		// if isPending != true {
		// 	log.Fatal("tx not pending anymore")
		// }
		data, _ := tx.MarshalJSON()
		// fmt.Println(string(data))

		type Tx struct {
			Input string `json:"input"`
		}
		var tx0 Tx
		json.Unmarshal(data, &tx0)
		fmt.Println("\n\ntx input: ", tx0.Input)
		time.Sleep(10 * time.Second)
	}
}
