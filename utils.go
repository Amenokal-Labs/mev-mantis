package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

func HandleTransaction(method, hash string) {
	if method != "POST" && method != "DELETE" {
		log.Fatal("bad request")
	}
	body := []byte(`
		{
			"apiKey":` + GetAPIKey("BLOCKNATIVE_KEY") + `,
			"hash":` + hash + `,
			"blockchain":"ethereum",
			"network":"main"
		}
	`)

	r, err := http.NewRequest("POST", "https://api.blocknative.com/transaction", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}
	r.Header.Add("Content-Type", "application/json")
	b, _ := ioutil.ReadAll(r.Body) // Log the request body
	log.Print(string(b))

	// implement client to get response
}

func HandleAddress(method, address string) {
	if method != "POST" && method != "DELETE" {
		log.Fatal("bad request")
	}
	body := []byte(`
		{
			"apiKey":` + GetAPIKey("BLOCKNATIVE_KEY") + `,
			"address":` + address + `,
			"blockchain":"ethereum",
			"networks":"main"
		}
	`)

	r, err := http.NewRequest("POST", "https://api.blocknative.com/address", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}
	r.Header.Add("Content-Type", "application/json")
	b, _ := ioutil.ReadAll(r.Body) // Log the request body
	log.Print(string(b))

	// implement client to get response
}

func GetAddresses(blockchain, network string) {
	url := "https://api.blocknative.com/address/" + GetAPIKey("BLOCKNATIVE_KEY") + "/" + blockchain + "/" + network + "/"
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(string(data))
}

func GetTransactions(blockchain, network string) {
	url := "https://api.blocknative.com/transaction/" + GetAPIKey("BLOCKNATIVE_KEY") + "/" + blockchain + "/" + network + "/"
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(string(data))
}
