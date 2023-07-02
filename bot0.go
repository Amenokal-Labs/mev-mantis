package main

import (
	"fmt"
	"log"
	"math/big"

	UniswapV2Pair "github.com/Amenokal-Labs/mev-mantis.git/UniswapV2Pair"
	UniswapV3Pair "github.com/Amenokal-Labs/mev-mantis.git/UniswapV3Pair"
	"github.com/Amenokal-Labs/mev-mantis.git/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/" + utils.GetAPIKey("INFURA_KEY"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("we have a connection")

	const V2PairAddress string = "0x0d4a11d5eeaac28ec3f61d100daf4d40471f1852" // WETH/USDT
	cmnV2PairAddress := common.HexToAddress(V2PairAddress)
	v2PairContract, err := UniswapV2Pair.NewUniswapV2Pair(cmnV2PairAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	data, err := v2PairContract.GetReserves(&bind.CallOpts{})
	fmt.Println(new(big.Int).Div(data.Reserve0, data.Reserve1))

	const V3PairAddress string = "0x11b815efB8f581194ae79006d24E0d814B7697F6"
	cmnV3PairAddress := common.HexToAddress(V3PairAddress)
	v3PairContract, err := UniswapV3Pair.NewUniswapV3Pair(cmnV3PairAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	data2, err := v3PairContract.Slot0(&bind.CallOpts{})
	fmt.Println(data2.SqrtPriceX96)
}
