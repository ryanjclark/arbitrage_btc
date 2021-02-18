package exchange

import (
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
	"github.com/ryanjclark/arbitrage_btc/forex/pkg/config"
)

// TraderMadeStreamer holds the URL, url Path, and API key to Connect to Trader Made forex stream.
type TraderMadeStreamer struct {
	forexConfig *config.ForexConfig
	Conn        *websocket.Conn
}

func (s *TraderMadeStreamer) initSocketConnection() (err error) {
	url := s.forexConfig.InitURL()
	log.Printf("Connecting to %s", url.String())

	c, resp, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		log.Printf("handshake failed with status %d", resp.StatusCode)
		return err
	}

	s.Conn = c
	return nil
}

func (s *TraderMadeStreamer) subscribe() (err error) {
	msg := fmt.Sprintf("{\"userKey\":\"%s\", \"symbol\":\"%s\"}", s.forexConfig.APIKey, s.forexConfig.Symbols[0])
	if err := s.Conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// GetPriceStream opens the WSS Connection and begins streaming forex data.
func (s *TraderMadeStreamer) GetPriceStream(c chan *WSSPayload, done chan struct{}) {
	// Connect to websocket
	if err := s.initSocketConnection(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer s.Conn.Close()

	// Subscribe to stream.
	if err := s.subscribe(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	for {
		message := WSSPayload{}
		if err := s.Conn.ReadJSON(&message); err != nil {
			log.Printf("error reading payload into message: %s", err)
		}

		c <- &message
		defer close(done)
	}
}

// NewTraderMadeStreamer creates a pointer to a new TraderMadeSocket.
func NewTraderMadeStreamer(forexConfig *config.ForexConfig) *TraderMadeStreamer {
	return &TraderMadeStreamer{
		forexConfig: forexConfig,
	}
}
