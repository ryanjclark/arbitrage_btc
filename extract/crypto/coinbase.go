package crypto

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

// CoinbaseSocket holds the socket URL to connect to.
type CoinbaseSocket struct {
	HostURL string
}

// GetPriceStream streams the price data from Coinbase given a ticker symbol.
func (s *CoinbaseSocket) GetPriceStream(symbol string) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c := s.connectToSocket()
	defer c.Close()

	msg := fmt.Sprintf("{\"type\": \"subscribe\", \"product_ids\": [\"%s\"], \"channels\": [{\"name\": \"ticker\"}]}", symbol)
	if err := c.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		log.Print(err)
	}

	done := make(chan struct{})

	go func() {
		defer close(done)
		message := TickerMessage{}
		for {
			if err := c.ReadJSON(&message); err != nil {
				log.Printf("error reading payload into message: %s", err)
				break
			}
		}
	}()

	for {
		select {
		case <-done:
			return
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

func (s *CoinbaseSocket) connectToSocket() *websocket.Conn {
	u := url.URL{Scheme: "wss", Host: s.HostURL}
	log.Printf("connecting to %s", u.String())

	c, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("handshake failed with status %d", resp.StatusCode)
		log.Fatal("conn error: ", err)
	}
	return c
}

// NewCoinbaseSocket creates a pointer to a new CoinbaseSocket.
func NewCoinbaseSocket(hostURL string) *CoinbaseSocket {
	return &CoinbaseSocket{
		HostURL: hostURL,
	}
}
