package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Url struct {
	action, address, tag string
}

func newUrl(action, address, tag string) *Url {
	return &Url{
		action:  action,
		address: address,
		tag:     tag,
	}
}

func getEtherscanKey() string {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	// gets ETHERSCAN_KEY
	etherscanKey, exists := os.LookupEnv("ETHERSCAN_KEY")

	if exists {
		return etherscanKey
	}
	return ""
}

func buildUrl(action, address, tag string) string {
	u := newUrl(action, address, tag)
	return "https://api.etherscan.io/api?module=account&action=" + u.action + "&address=" + u.address + "&tag=" + u.tag + "&apikey=" + getEtherscanKey()
}

func call(url string) string {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(data)
}

// returns the Ether balance of a given address
func getBalance(address, tag string) string {
	return call(buildUrl("balance", address, tag))
}

func getBalances(address, tag string) string {
	return call(buildUrl("balancemulti", address, tag))
}

func main() {
	fmt.Println(getBalance("0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae", "latest"))
}
