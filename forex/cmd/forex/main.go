package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/ryanjclark/arbitrage_btc/forex/pkg/config"
	"github.com/ryanjclark/arbitrage_btc/forex/pkg/exchange"
)

func gRPCStub(msg *exchange.WSSPayload) {
	log.Println(msg)
}

func main() {
	if err := godotenv.Load("dev.env"); err != nil {
		log.Println("Warning: .env file not found")
	}

	var (
		forexURL = os.Getenv("TRADER_MADE_URL")
		forexKey = os.Getenv("TRADER_MADE_API")
	)

	fc := config.NewForexConfig(forexURL, "wss", "/feed", []string{"EURUSD"}, forexKey)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})
	c := make(chan *exchange.WSSPayload)

	streamer := exchange.NewTraderMadeStreamer(fc)
	go streamer.GetPriceStream(c, done)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := streamer.Conn.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write: ", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			err := streamer.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("error closing message: ", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		case m := <-c:
			gRPCStub(m)
		}
	}
}
