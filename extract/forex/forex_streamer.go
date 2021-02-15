package forex

import "github.com/gorilla/websocket"

// ForexStreamer interface implements a function to connect
// and get the price of a given ticker.
type ForexStreamer interface {
	connectToSocket() *websocket.Conn
	GetPrice(string)
}
