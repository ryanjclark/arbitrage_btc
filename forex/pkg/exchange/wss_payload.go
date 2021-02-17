package exchange

import "time"

// WSSPayload holds the data from a forex websocket stream.
type WSSPayload struct {
	symbol *string
	bid    *float32
	ask    *float32
	dt     time.Time
}

// NewWSSPayload hydreates a new WSSPayload struct from the websocket feed.
func NewWSSPayload(symbol *string, bid *float32, ask *float32, dt time.Time) *WSSPayload {
	return &WSSPayload{
		symbol: symbol,
		bid:    bid,
		ask:    ask,
		dt:     dt,
	}
}
