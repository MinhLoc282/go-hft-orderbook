package hftorderbook

// Single Order in an order book, as a node in a LimitOrder FIFO queue
type Order struct {
	Id       string
	UserId   string
	Volume   float64
	Next     *Order
	Prev     *Order
	Limit    *LimitOrder
	BidOrAsk bool
}
