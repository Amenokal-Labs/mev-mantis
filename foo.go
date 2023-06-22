package main

import (
	"log"
	"os"

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

// implement using ordered maps https://github.com/wk8/go-ordered-map/blob/v2.1.7/orderedmap.go#L23
func newFields(fields ...string) {
	return
}

func buildURL(base, endpoint string, fields map[string]string) {

}
