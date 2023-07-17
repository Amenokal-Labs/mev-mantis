package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Amenokal-Labs/mev-mantis.git/pkg/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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
	if err != nil {
		panic(err)
	}
	subscriber := gethclient.New(rpc)
	fmt.Println("client connected..")

	// txs := make(chan *types.Transaction)
	// _, err = rpcClient.SubscribeFullPendingTransactions(context.Background(), txs)
	// if err != nil {
	// 	panic(err)
	// }

	hashes := make(chan common.Hash)
	_, err = subscriber.SubscribePendingTransactions(context.Background(), hashes)
	if err != nil {
		panic(err)
	}

	for {
		pc, _ := httpsClient.PendingTransactionCount(context.Background())
		fmt.Println("\nPending count:", pc)

		hash := <-hashes
		fmt.Println("      Tx hash:", hash)
		txn, _, err := httpsClient.TransactionByHash(context.Background(), hash)
		if err != nil {
			panic(err)
		}
		data, err := txn.MarshalJSON()
		if err != nil {
			panic(err)
		}

		from, err := types.Sender(types.NewLondonSigner(txn.ChainId()), txn)
		if err != nil {
			panic(err)
		}

		type Tx struct {
			Input string `json:"input"`
		}
		var tx0 Tx
		json.Unmarshal(data, &tx0)
		fmt.Println("  Tx Calldata:", tx0.Input)
		fmt.Println("   Tx address:", txn.To())

		// get contract code if any
		type Contract struct {
			JsonRpc string `json:"jsonrpc"`
			Id      int    `json:"id"`
			Result  string `json:"result"`
		}
		body := []byte(`{
			"jsonrpc":"2.0",
			"method":"eth_getCode",
			"params": ["` + txn.To().String() + `", "pending"],
			"id":1
		}`)
		r, err := http.NewRequest("POST", "https://mainnet.infura.io/v3/"+utils.GetAPIKey("INFURA_KEY"), bytes.NewBuffer(body))
		if err != nil {
			panic(err)
		}
		r.Header.Add("Content-Type", "application/json")
		client := &http.Client{}
		res, err := client.Do(r)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()
		contract := &Contract{}
		derr := json.NewDecoder(res.Body).Decode(contract)
		if derr != nil {
			panic(derr)
		}
		// fmt.Println("\nContract code:", contract.Result)
		// if res.StatusCode != http.StatusCreated {
		// 	panic(res.Status)
		// }

		if contract.Result == "0x" {
			continue
		}
		contractCode := contract.Result
		fmt.Println("\nContract code:", contractCode[:26], "...")

		calldata := tx0.Input
		PUSH21 := "74"    // opcode, push 21-byte value onto stack
		WORD_LENGTH := 64 // 32 bytes
		N_ZEROS := 22     //
		for i := 0; i < len(calldata); i = i + 2 {
			if (string(calldata[i])+string(calldata[i+1]) == PUSH21) && (calldata[i+1+N_ZEROS:i+1+WORD_LENGTH] == from.String()) {
				fmt.Println("replace address at index i+23")
			}
		}
		fmt.Println("___________________________")
		time.Sleep(4 * time.Second)
	}
}
