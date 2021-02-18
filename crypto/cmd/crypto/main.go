package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/ryanjclark/arbitrage_btc/crypto/pkg/exchange"
)

func main() {
	err := godotenv.Load("dev.env")
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	var bitcoinURL = os.Getenv("COINBASE_URL")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})

	c := make(chan *exchange.WSSPayload)

	btcStreamer := exchange.NewCoinbaseSocket(bitcoinURL)

	exchange.InitConnection(btcStreamer)
	defer btcStreamer.Conn.Close()

	btcStreamer.Subscribe("BTC-USD")

	go btcStreamer.GetPriceStream(c, done)
	go futureGRPCFunc(c)

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			err := btcStreamer.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("error closing message: ", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func futureGRPCFunc(c chan *exchange.WSSPayload) {
	for m := range c {
		log.Println(m)
	}
}
