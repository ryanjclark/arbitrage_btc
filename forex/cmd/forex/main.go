package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/ryanjclark/arbitrage_btc/forex/pkg/exchange"
)

func main() {
	err := godotenv.Load("dev.env")
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	var (
		forexURL = os.Getenv("TRADER_MADE_URL")
		forexKey = os.Getenv("TRADER_MADE_API")
	)

	s := exchange.NewTraderMadeSocket(forexURL, "/feed", forexKey)
	exchangeTicker := make(chan []byte)
	go s.GetPriceStream("EURUSD", exchangeTicker)

	select {
	case e := <-exchangeTicker:
		log.Printf("received %s", e)
	}
}
