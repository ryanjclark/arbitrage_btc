package crypto

import "github.com/gorilla/websocket"

// Streamer implements a function to connect
// and a function to stream price data of a given ticker.
type Streamer interface {
	connectToSocket() *websocket.Conn
	GetPriceStream(string)
}
