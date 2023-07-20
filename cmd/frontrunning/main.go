package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"
	"net/http"
	"time"

	"github.com/Amenokal-Labs/mev-mantis.git/pkg/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
)

func main() {
	ethclient, err := ethclient.Dial("https://goerli.infura.io/v3/" + utils.GetKey("INFURA_KEY"))
	if err != nil {
		log.Fatal("[0] ", err)

	}

	rpc, err := rpc.Dial("wss://goerli.infura.io/ws/v3/" + utils.GetKey("INFURA_KEY"))
	if err != nil {
		log.Fatal("[1] ", err)
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
		log.Fatal("[2] ", err)
	}

	for {
		pc, _ := ethclient.PendingTransactionCount(context.Background())
		fmt.Println("\nPending count:", pc)

		hash := <-hashes
		fmt.Println("      Tx hash:", hash)
		txn, _, err := ethclient.TransactionByHash(context.Background(), hash)
		if err != nil {
			log.Fatal("[3] ", err)
		}

		marshalledTxn, err := txn.MarshalJSON()
		if err != nil {
			log.Fatal("[4] ", err)
		}

		from, err := types.Sender(types.NewLondonSigner(txn.ChainId()), txn)
		if err != nil {
			log.Fatal("[5] ", err)
		}

		type Tx struct {
			Input string `json:"input"`
		}
		var tx Tx
		json.Unmarshal(marshalledTxn, &tx)
		fmt.Println("  Tx Calldata:", tx.Input)
		fmt.Println("   Tx address:", txn.To())

		// get contract code if any
		type Contract struct {
			JsonRpc string `json:"jsonrpc"`
			Id      int    `json:"id"`
			Result  string `json:"result"`
		}
		reciever := txn.To()
		if reciever == nil {
			log.Fatal("[14] transaction is a contract creation.")
		}
		body := []byte(`{
			"jsonrpc":"2.0",
			"method":"eth_getCode",
			"params": ["` + reciever.String() + `", "pending"],
			"id":1
		}`)
		r, err := http.NewRequest("POST", "https://mainnet.infura.io/v3/"+utils.GetKey("INFURA_KEY"), bytes.NewBuffer(body))
		if err != nil {
			log.Fatal("[6] ", err)
		}
		r.Header.Add("Content-Type", "application/json")
		client := &http.Client{}
		res, err := client.Do(r)
		if err != nil {
			log.Fatal("[7] ", err)
		}
		defer res.Body.Close()
		contract := &Contract{}
		derr := json.NewDecoder(res.Body).Decode(contract)
		if derr != nil {
			log.Fatal("[8] ", derr)
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

		calldata := tx.Input
		address := from.String()
		if replaceAddress(calldata, address) != calldata {
			sendTxn(ethclient, txn)
		}

		fmt.Println("___________________________")
		time.Sleep(4 * time.Second)
	}
}

func replaceAddress(_calldata, _address string) string {
	PUSH21 := "74"    // opcode, push 21-byte value onto stack
	WORD_LENGTH := 64 // 32 bytes
	N_ZEROS := 22     // number of zeros before first non null bit of the address
	calldata := _calldata

	for i := 0; i < len(calldata); i = i + 2 {
		if (string(calldata[i])+string(calldata[i+1]) == PUSH21) && (calldata[i+1+N_ZEROS:i+1+WORD_LENGTH] == _address) {
			calldata = calldata[0:i+1] + _address + calldata[i+2+WORD_LENGTH:]
		}
	}

	return calldata
}

func createTxn(_client *ethclient.Client, _originalTxn *types.Transaction) *types.Transaction {
	privateKey, err := crypto.HexToECDSA(utils.GetKey("PRIVATE_KEY"))
	if err != nil {
		log.Fatal("[9] ", err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := _client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal("[10] ", err)
	}
	toAddress := _originalTxn.To()
	value := _originalTxn.Value()
	gasLimit := uint64(math.MaxUint64)
	gasPrice := _originalTxn.GasPrice()
	data := _originalTxn.Data()
	tx := types.NewTransaction(nonce, *toAddress, value, gasLimit, gasPrice, data)

	GOERLI_ID := big.NewInt(5)
	signedTxn, err := types.SignTx(tx, types.NewLondonSigner(GOERLI_ID), privateKey)
	if err != nil {
		log.Fatal("[11] ", err)
	}
	rawTxnBytes, err := rlp.EncodeToBytes(signedTxn)
	if err != nil {
		log.Fatal("[12] ", err)
	}
	rawTxnHex := hex.EncodeToString(rawTxnBytes)
	fmt.Println("\ntransaction created:", rawTxnHex)

	return signedTxn
}

func sendTxn(_client *ethclient.Client, _originalTxn *types.Transaction) {
	txn := createTxn(_client, _originalTxn)
	err := _client.SendTransaction(context.Background(), txn)
	if err != nil {
		log.Fatal("[13] ", err)
	}
}
