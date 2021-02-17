package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/ryanjclark/arbitrage_btc/forex_streamer/extract"
)

func main() {
	err := godotenv.Load("dev.env")
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	var forexURL = os.Getenv("TRADER_MADE_URL")
	var forexKey = os.Getenv("TRADER_MADE_API")

	_ = extract.NewTraderMadeSocket(forexURL, "/feed", forexKey)
}
