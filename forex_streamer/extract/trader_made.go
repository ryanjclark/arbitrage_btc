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

// TraderMadeSocket holds the URL, url Path, and API key to
// connect to Trader Made forex stream.
type TraderMadeSocket struct {
	HostURL string
	Path    string
	Key     string
}

// GetPriceStream streams the price data from Trader Made given a ticker symbol.
func (s *TraderMadeSocket) GetPriceStream(symbol string) {
	messageOut := make(chan string)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c := s.connectToSocket()
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read error: ", err)
				return
			}
			log.Printf("recv: %s", message)
			if string(message) == "Connected" {
				log.Printf("Send Sub Details: %s", message)
				messageOut <- fmt.Sprintf("{\"userKey\":\"%s\", \"symbol\":\"%s\"}", s.Key, symbol)
			}
		}
	}()

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

func (s *TraderMadeSocket) connectToSocket() *websocket.Conn {
	u := url.URL{Scheme: "wss", Host: s.HostURL, Path: s.Path}
	log.Printf("connecting to %s", u.String())

	c, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("handshake failed with status %d", resp.StatusCode)
		log.Fatal("conn error: ", err)
	}
	return c
}

// NewTraderMadeSocket creates a pointer to a new TraderMadeSocket.
func NewTraderMadeSocket(hostURL string, path string, key string) *TraderMadeSocket {
	return &TraderMadeSocket{
		HostURL: hostURL,
		Path:    path,
		Key:     key,
	}
}
