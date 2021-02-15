package crypto

import "github.com/gorilla/websocket"

// CryptoStreamer implements a function to connect
// and a function to stream price data of a given ticker.
type CryptoStreamer interface {
	connectToSocket() *websocket.Conn
	GetPriceStream(string)
}
