package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os/exec"

	"github.com/Amenokal-Labs/mev-mantis.git/pkg/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

const FOUNDRY_ACCOUNT_ADDRESS1 = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
const FOUNDRY_ACCOUNT_ADDRESS2 = "0x976EA74026E726554dB657fA54763abd0C3a0aa9"

func main() {
	ethclient, err := ethclient.Dial("https://sepolia.infura.io/v3/" + utils.GetKey("INFURA_KEY"))
	if err != nil {
		log.Fatal("Client not connected.", err)
	}

	// rpc, err := rpc.Dial("wss://sepolia.infura.io/ws/v3/" + utils.GetKey("INFURA_KEY"))
	// if err != nil {
	// 	log.Fatal("Client not connected. ", err)
	// }
	// subscriber := gethclient.New(rpc)
	fmt.Println("client connected..")

	// txs := make(chan *types.Transaction)
	// _, err = rpcClient.SubscribeFullPendingTransactions(context.Background(), txs)
	// if err != nil {
	// 	panic(err)
	// }
	from := common.HexToAddress(FOUNDRY_ACCOUNT_ADDRESS1)
	nonce, _ := ethclient.PendingNonceAt(context.Background(), from)
	suggestedGasPrice, _ := ethclient.SuggestGasPrice(context.Background())
	gasPrice := new(big.Int).Mul(suggestedGasPrice, big.NewInt(10))
	gas := uint64(30000000)
	to := common.HexToAddress(FOUNDRY_ACCOUNT_ADDRESS2)
	value := big.NewInt(1000000000000)
	data := []byte("hoy--")

	tx := createTx(nonce, gasPrice, gas, to, value, data)

	sepoliaID := big.NewInt(11155111)
	privateKey, err := crypto.HexToECDSA(utils.GetKey("PRIVATE_KEY"))
	if err != nil {
		log.Fatal("[9] ", err)
	}

	signedTx, err := signTx(tx, sepoliaID, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	hexTx := hexTx(signedTx)
	fmt.Println("raw tx:", hexTx)

	// for {
	// 	pc, _ := ethclient.PendingTransactionCount(context.Background())
	// 	fmt.Println("\nPending count:", pc)

	// 	txn, marshalledTxn, _ := getTxn(ethclient, subscriber)
	// 	calldata, contract := getCalldata(txn, marshalledTxn)
	// 	if contract == "0x" {
	// 		continue
	// 	}
	// 	fmt.Println("    Txn:", txn)

	// 	fmt.Println("    Tx Calldata:", calldata)
	// 	fmt.Println("     Tx address:", txn.To())
	// 	fmt.Println("\nContract code:", contract[:26], "...")

	// 	txn, rawTxnBytes := createTxn(ethclient, txn)
	// 	simulateTxn(rawTxnBytes)

	// 	// address := from.String()
	// 	// if replaceAddress(calldata, address) != calldata {
	// 	// 	sendTxn(ethclient, txn)
	// 	// }

	// 	fmt.Println("___________________________")
	// 	time.Sleep(4 * time.Second)
	// }
}

func getTxn(_ethclient *ethclient.Client, _subscriber *gethclient.Client) (txn *types.Transaction, marshalledTxn []byte, from common.Address) {
	hashes := make(chan common.Hash)
	_, err := _subscriber.SubscribePendingTransactions(context.Background(), hashes)
	if err != nil {
		log.Fatal("[2] ", err)
	}
	hash := <-hashes
	fmt.Println("      Tx hash:", hash)
	txn, _, err = _ethclient.TransactionByHash(context.Background(), hash)
	if err != nil {
		log.Fatal("[3] ", err)
	}

	marshalledTxn, err = txn.MarshalJSON()
	if err != nil {
		log.Fatal("[4] ", err)
	}

	from, err = types.Sender(types.NewLondonSigner(txn.ChainId()), txn)
	if err != nil {
		log.Fatal("[5] ", err)
	}

	return txn, marshalledTxn, from
}

func getCalldata(_txn *types.Transaction, _marshalledTxn []byte) (calldata string, contract string) {
	type Contract struct {
		JsonRpc string `json:"jsonrpc"`
		Id      int    `json:"id"`
		Code    string `json:"result"`
	}
	to := _txn.To()
	if to == nil {
		log.Fatal("[14] transaction is a contract creation.")
	}
	body := []byte(`{
		"jsonrpc":"2.0",
		"method":"eth_getCode",
		"params": ["` + to.String() + `", "pending"],
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
	c := &Contract{}
	derr := json.NewDecoder(res.Body).Decode(c)
	if derr != nil {
		log.Fatal("[8] ", derr)
	}
	// if res.StatusCode != http.StatusCreated {
	// 	panic(res.Status)
	// }

	type Txn struct {
		Input string `json:"input"`
	}
	var tx Txn
	json.Unmarshal(_marshalledTxn, &tx)
	return tx.Input, c.Code
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

func simulateTxn(_rawTxnBytes []byte) {
	__rpc_url := utils.GetKey("RPC_URL")
	cmd := exec.Command("cast", "publish", __rpc_url, string(_rawTxnBytes))
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

// func sendTxn(_client *ethclient.Client, _originalTxn *types.Transaction) {
// 	txn, _ := createTxn(_client, _originalTxn)
// 	to := _originalTxn.To()
// 	value := _originalTxn.Value()
// 	gasLimit := uint64(math.MaxUint64)
// 	gasPrice := _originalTxn.GasPrice()
// 	data := _originalTxn.Data()

// 	privateKey, err := crypto.HexToECDSA(utils.GetKey("PRIVATE_KEY"))
// 	if err != nil {
// 		log.Fatal("[9] ", err)
// 	}
// 	publicKey := privateKey.Public()
// 	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
// 	if !ok {
// 		log.Fatal("error casting public key to ECDSA")
// 	}
// 	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

// 	nonce, err := _client.PendingNonceAt(context.Background(), fromAddress)
// 	if err != nil {
// 		log.Fatal("[10] ", err)
// 	}

// 	err := _client.SendTransaction(context.Background(), txn)
// 	if err != nil {
// 		log.Fatal("[13] ", err)
// 	}
// }

func createTx(nonce uint64, gasPrice *big.Int, gas uint64, to *common.Address, value *big.Int, data []byte) (tx *types.Transaction) {
	return types.NewTx(&types.LegacyTx{
		Nonce:    nonce,    // nonce of sender account
		GasPrice: gasPrice, // wei per gas
		Gas:      gas,      // gas limit
		To:       to,       // nil means contract creation
		Value:    value,    // wei amount
		Data:     data,     // contract invocation input data
	})
}

func signTx(tx *types.Transaction, chainID *big.Int, privatekey *ecdsa.PrivateKey) (*types.Transaction, error) {
	signedTxn, err := types.SignTx(tx, types.NewLondonSigner(chainID), privatekey)
	if err != nil {
		log.Fatal("[11] ", err)
	}

	return signedTxn, err
}

func hexTx(tx *types.Transaction) string {
	rawTxnBytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		log.Fatal("[12] ", err)
	}

	return hex.EncodeToString(rawTxnBytes)
}
