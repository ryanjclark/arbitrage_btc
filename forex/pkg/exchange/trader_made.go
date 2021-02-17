package exchange

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

// TraderMadeStreamer holds the URL, url Path, and API key to
// connect to Trader Made forex stream.
type TraderMadeStreamer struct {
	url    *url.URL
	apiKey string
}

// ConnectToSocket establishes a websocket connection.
func (s *TraderMadeStreamer) ConnectToSocket() *websocket.Conn {
	log.Printf("connecting to %s", s.url.String())

	c, resp, err := websocket.DefaultDialer.Dial(s.url.String(), nil)
	if err != nil {
		log.Printf("handshake failed with status %d", resp.StatusCode)
		log.Fatal("conn error: ", err)
	}
	return c
}

// NewTraderMadeStreamer creates a pointer to a new TraderMadeSocket.
func NewTraderMadeStreamer(url *url.URL, key string) *TraderMadeStreamer {
	return &TraderMadeStreamer{
		url:    url,
		apiKey: key,
	}
}
