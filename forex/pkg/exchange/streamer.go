package exchange

import "github.com/gorilla/websocket"

// Streamer interface implements methods to connect
// and get the price of a given ticker.
type Streamer interface {
	ConnectToSocket() *websocket.Conn
}
