package main

import (
	"log"
	"os"

	"github.com/arbitrage_btc/extract/crypto"
	"github.com/arbitrage_btc/extract/forex"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("dev.env")
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	var forexURL = os.Getenv("TRADER_MADE_URL")
	var forexKey = os.Getenv("TRADER_MADE_API")

	_ = forex.NewTraderMadeSocket(forexURL, "/feed", forexKey)

	var bitcoinURL = os.Getenv("COINBASE_URL")
	_ = crypto.NewCoinbaseSocket(bitcoinURL)
}
