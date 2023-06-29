package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/" + GetAPIKey("INFURA_KEY"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("we have a connection")

	address := common.HexToAddress("0x3356c9a8f40f8e9c1d192a4347a76d18243fabc5")

	contract, err := NewMainCaller(address, client)
	if err != nil {
		log.Fatal(err)
	}

	data, err := contract.GetReserves(&bind.CallOpts{})
	fmt.Println(data)
}
