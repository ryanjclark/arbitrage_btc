package extract

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type CoinbaseSocket struct {
	HostURL string
}

func (s *CoinbaseSocket) GetPrice(symbol string) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	u := url.URL{Scheme: "wss", Host: s.HostURL}
	log.Printf("connecting to %s", u.String())

	dialer := websocket.Dialer{HandshakeTimeout: 120 * time.Second}
	c, resp, err := dialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("handshake failed with status %d", resp.StatusCode)
		log.Fatal("dial:", err)
	}
	defer c.Close()

	msg := fmt.Sprintf("{\"type\": \"subscribe\", \"product_ids\": [\"%s\"], \"channels\": [{\"name\": \"ticker\"}]}", symbol)
	if err := c.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		log.Print(err)
	}

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read: ", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
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

func NewCoinbaseSocket(hostURL string) *CoinbaseSocket {
	return &CoinbaseSocket{
		HostURL: hostURL,
	}
}
