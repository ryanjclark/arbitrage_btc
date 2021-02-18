package exchange

// WSSPayload holds the payload from Coinbase's websocket
// channel called 'ticker'.
type WSSPayload struct {
	Type       string   `json:"type"`
	TradeID    int      `json:"trade_id"`
	Sequence   int      `json:"sequence"`
	Time       string   `json:"time"`
	ProductID  string   `json:"product_id"`
	ProductIDs []string `json:"product_ids"`
	Side       string   `json:"side"`
	LastSize   string   `json:"last_size"`
	BestBid    string   `json:"best_bid"`
	BestAsk    string   `json:"best_ask"`
}
