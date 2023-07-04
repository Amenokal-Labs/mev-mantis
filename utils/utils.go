package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func GetAPIKey(keyName string) string {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	// gets API_KEY
	apiKey, exists := os.LookupEnv(keyName)
	if exists {
		return apiKey
	}
	return ""
}

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

	r, err := http.NewRequest(method, "https://api.blocknative.com/transaction", bytes.NewBuffer(body))
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

	r, err := http.NewRequest(method, "https://api.blocknative.com/address", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}
	r.Header.Add("Content-Type", "application/json")

	// implement client to get response
	// httpClient := &http.Client{}
	c := http.Client{Timeout: time.Duration(1) * time.Second}
	resp, err := c.Do(r)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	fmt.Println(r)
	b, _ := ioutil.ReadAll(resp.Body) // Log the request body
	fmt.Println(b)
	log.Print(string(b))
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
