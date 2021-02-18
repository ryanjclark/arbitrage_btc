package exchange

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

// CoinbaseStreamer holds the URL to connect to and the websocket connection.
type CoinbaseStreamer struct {
	HostURL string
	Conn    *websocket.Conn
}

// GetPriceStream reads the JSON from the connection and hydrates the payload.
func (s *CoinbaseStreamer) GetPriceStream(c chan *WSSPayload, done chan struct{}) {
	for {
		message := WSSPayload{}
		if err := s.Conn.ReadJSON(&message); err != nil {
			log.Printf("error reading payload into message: %s", err)
		}
		c <- &message
		defer close(done)
	}
}

// Subscribe sends the subscription message to trigger the websocket channel
func (s *CoinbaseStreamer) Subscribe(symbol string) {
	msg := fmt.Sprintf("{\"type\": \"subscribe\", \"product_ids\": [\"%s\"], \"channels\": [{\"name\": \"ticker\"}]}", symbol)
	if err := s.Conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		log.Print(err)
	}
}

// InitConnection connects to the HostURL and gives the CoinbaseStreamer the connection.
func InitConnection(s *CoinbaseStreamer) {
	u := url.URL{Scheme: "wss", Host: s.HostURL}
	log.Printf("connecting to %s", u.String())

	conn, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("handshake failed with status %d", resp.StatusCode)
		log.Fatalf("conn error: %s", err)
	}
	s.Conn = conn
}

// NewCoinbaseStreamer creates a pointer to a new CoinbaseStreamer.
func NewCoinbaseStreamer(hostURL string) *CoinbaseStreamer {
	return &CoinbaseStreamer{
		HostURL: hostURL,
	}
}
