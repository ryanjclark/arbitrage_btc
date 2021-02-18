package exchange

import (
	"fmt"
	"log"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

//
type CoinbaseSocket struct {
	HostURL string
	Mu      sync.Mutex
	Conn    *websocket.Conn
}

//
func (s *CoinbaseSocket) GetPriceStream(c chan *WSSPayload, done chan struct{}) {
	for {
		message := WSSPayload{}
		if err := s.Conn.ReadJSON(&message); err != nil {
			log.Printf("error reading payload into message: %s", err)
		}
		c <- &message
		defer close(done)
	}
}

func InitConnection(s *CoinbaseSocket) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	u := url.URL{Scheme: "wss", Host: s.HostURL}
	log.Printf("connecting to %s", u.String())

	conn, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("handshake failed with status %d", resp.StatusCode)
		log.Fatalf("conn error: %s", err)
	}
	s.Conn = conn
}

func (s *CoinbaseSocket) Subscribe(symbol string) {
	msg := fmt.Sprintf("{\"type\": \"subscribe\", \"product_ids\": [\"%s\"], \"channels\": [{\"name\": \"ticker\"}]}", symbol)
	if err := s.Conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		log.Print(err)
	}
}

// NewCoinbaseSocket creates a pointer to a new CoinbaseSocket.
func NewCoinbaseSocket(hostURL string) *CoinbaseSocket {
	return &CoinbaseSocket{
		HostURL: hostURL,
	}
}
