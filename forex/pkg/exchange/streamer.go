package exchange

import "github.com/gorilla/websocket"

// Streamer interface implements a function to connect
// and get the price of a given ticker.
type Streamer interface {
	connectToSocket() *websocket.Conn
	GetPrice(string)
}
