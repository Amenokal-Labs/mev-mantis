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

	const V2PairAddress string = "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc" // WETH/USDT
	cmnV2PairAddress := common.HexToAddress(V2PairAddress)
	v2PairContract, err := UniswapV2Pair.NewUniswapV2Pair(cmnV2PairAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	data, err := v2PairContract.GetReserves(&bind.CallOpts{})

	price0 := new(big.Float).Quo(new(big.Float).SetInt(data.Reserve1), new(big.Float).SetInt(data.Reserve0))

	adjustmentFactor0 := new(big.Int).Exp(big.NewInt(10), big.NewInt(12), nil)
	adjustedPrice0 := new(big.Float).Quo(price0, new(big.Float).SetInt(adjustmentFactor0))
	roundedPrice0 := fmt.Sprintf("V2 Price: 1 USDC = %.10f ETH", adjustedPrice0)
	fmt.Println(roundedPrice0)

	const V3PairAddress string = "0x88e6A0c2dDD26FEEb64F039a2c41296FcB3f5640"
	cmnV3PairAddress := common.HexToAddress(V3PairAddress)
	v3PairContract, err := UniswapV3Pair.NewUniswapV3Pair(cmnV3PairAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	data2, err := v3PairContract.Slot0(&bind.CallOpts{})
	// fmt.Println("SqrtPriceX96:", data2.SqrtPriceX96)

	q := new(big.Int).Exp(big.NewInt(2), big.NewInt(96), nil)
	// fmt.Println("Q96:", q)

	sqrtPrice := new(big.Float).Quo(new(big.Float).SetInt(data2.SqrtPriceX96), new(big.Float).SetInt(q))
	// fmt.Println("sqrtPrice:", sqrtPrice)

	price := new(big.Float).Mul(sqrtPrice, sqrtPrice)
	// fmt.Println("Price:", price)

	adjustmentFactor := new(big.Int).Exp(big.NewInt(10), big.NewInt(12), nil)
	adjustedPrice := new(big.Float).Quo(price, new(big.Float).SetInt(adjustmentFactor))

	roundedPrice := fmt.Sprintf("V3 Price: 1 USDC = %.10f ETH", adjustedPrice)
	fmt.Println(roundedPrice)
}
