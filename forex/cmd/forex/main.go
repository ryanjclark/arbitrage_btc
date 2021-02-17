package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/ryanjclark/arbitrage_btc/forex/pkg/exchange"
)

var (
	forexURL = os.Getenv("TRADER_MADE_URL")
	forexKey = os.Getenv("TRADER_MADE_API")
)

func main() {
	err := godotenv.Load("dev.env")
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	s := exchange.NewTraderMadeSocket(forexURL, "/feed", forexKey)
	s.GetPriceStream("USD")
}
