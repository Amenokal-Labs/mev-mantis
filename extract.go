package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func get_etherscan_key() string {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	// get ETHERSCAN_KEY
	etherscanKey, exists := os.LookupEnv("ETHERSCAN_KEY")

	if exists {
		return etherscanKey
	}
	return ""
}

// returns the Ether balance of a given address
func get_balance(address, tag string) string {
	url := "https://api.etherscan.io/api?module=account&action=balance&address=" + address + "&tag=" + tag + "&apikey=" + get_etherscan_key()

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