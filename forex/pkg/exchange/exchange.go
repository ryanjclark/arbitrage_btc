package exchange

import "time"

// Exchange holds the data from a forex websocket stream.
type Exchange struct {
	symbol *string
	bid    *float32
	ask    *float32
	dt     time.Time
}

// NewExchange hydreates a new exchange struct from the websocket feed.
func NewExchange(symbol *string, bid *float32, ask *float32, dt time.Time) *Exchange {
	return &Exchange{
		symbol: symbol,
		bid:    bid,
		ask:    ask,
		dt:     dt,
	}
}
