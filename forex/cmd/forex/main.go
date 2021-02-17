package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/ryanjclark/arbitrage_btc/forex/pkg/exchange"
)

type forexConfig struct {
	exchange   string
	Symbols    []string
	APIKey     string
	needTicker bool
}

type traderMadeStreamer struct {
	url    *url.URL
	config *forexConfig
}

// ConnectToSocket establishes a websocket connection.
func (s *traderMadeStreamer) ConnectToSocket() *websocket.Conn {
	log.Printf("connecting to %s", s.url.String())

	c, resp, err := websocket.DefaultDialer.Dial(s.url.String(), nil)
	if err != nil {
		log.Printf("handshake failed with status %d", resp.StatusCode)
		log.Fatal("conn error: ", err)
	}
	return c
}

func getPriceStream(streamer exchange.Streamer, config *forexConfig, ch chan []byte) {
	messageOut := make(chan string)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c := streamer.ConnectToSocket()
	defer c.Close()

	done := make(chan struct{})

	go func(ch chan []byte) {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			ch <- message
			if err != nil {
				log.Println("read error: ", err)
				return
			}
			if string(message) == "Connected" {
				log.Printf("Send Sub Details: %s", message)
				messageOut <- fmt.Sprintf("{\"userKey\":\"%s\", \"symbol\":\"%s\"}", config.APIKey, config.Symbols[0])
			}
		}
	}(ch)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case m := <-messageOut:
			log.Printf("Send Message %s", m)
			err := c.WriteMessage(websocket.TextMessage, []byte(m))
			if err != nil {
				log.Println("write message error: ", err)
				return
			}
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write: ", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
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

func main() {
	err := godotenv.Load("dev.env")
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	var (
		forexURL = os.Getenv("TRADER_MADE_URL")
		forexKey = os.Getenv("TRADER_MADE_API")
	)

	ch := make(chan []byte)
	u := &url.URL{
		Scheme: "wss",
		Host:   forexURL,
		Path:   "/feed",
	}

	s := exchange.NewTraderMadeStreamer(u, forexKey)
	getPriceStream(
		s,
		&forexConfig{
			exchange:   "TraderMade",
			Symbols:    []string{"EURUSD"},
			APIKey:     forexKey,
			needTicker: true,
		},
		ch,
	)

	// exchangeTicker := make(chan []byte)
	// s.GetPriceStream("EURUSD", exchangeTicker)
	//
	// select {
	// case e := <-exchangeTicker:
	// 	log.Printf("received %s", e)
	// }
}
