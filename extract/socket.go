package extract

import "github.com/gorilla/websocket"

// Socket interface implements a function to connect
// and get the price of a given ticker.
type Socket interface {
	connectToSocket() *websocket.Conn
	GetPrice(string)
}
