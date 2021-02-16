package main

import (
	"log"
	"os"

	"github.com/arbitrage_btc/crypto_streamer/extract"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("dev.env")
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	var bitcoinURL = os.Getenv("COINBASE_URL")

	_ = extract.NewCoinbaseSocket(bitcoinURL)
}
