package hftorderbook

import "github.com/shopspring/decimal"

// Limit price orders combined as a FIFO queue
type LimitOrder struct {
	Price float64

	orders      *ordersQueue
	totalVolume float64
}

func NewLimitOrder(price float64) LimitOrder {
	q := NewOrdersQueue()
	return LimitOrder{
		Price:  price,
		orders: &q,
	}
}

func (this *LimitOrder) EnqueueBulk(orders []*Order) {
	this.orders.BulkEnqueue(orders)
	for _, o := range orders {
		o.Limit = this
		this.totalVolume += o.Volume
	}
}

func (this *LimitOrder) TotalVolume() float64 {
	return this.totalVolume
}

func (this *LimitOrder) Size() int {
	return this.orders.Size()
}

func (this *LimitOrder) Enqueue(o *Order) {
	this.orders.Enqueue(o)
	o.Limit = this
	this.totalVolume += o.Volume
}

func (this *LimitOrder) Dequeue() *Order {
	if this.orders.IsEmpty() {
		return nil
	}

	o := this.orders.Dequeue()
	this.totalVolume -= o.Volume
	return o
}

func (this *LimitOrder) Delete(o *Order) {
	if o.Limit != this {
		panic("order does not belong to the limit")
	}

	this.orders.Delete(o)
	o.Limit = nil
	this.totalVolume -= o.Volume
}

func (this *LimitOrder) Clear() {
	q := NewOrdersQueue()
	this.orders = &q
	this.totalVolume = 0
}

func (this *LimitOrder) Peek(index int) *Order {
	return this.orders.Peek(index)
}

func (lo *LimitOrder) AddVolume(volumeToAdd float64) {
	lo.totalVolume += volumeToAdd
}

func (lo *LimitOrder) SubtractVolume(volumeToSubtract decimal.Decimal) {
	loTotalVolume := decimal.NewFromFloat(lo.totalVolume)

	if loTotalVolume.LessThan(volumeToSubtract) {
		panic("volume to subtract exceeds total volume")
	}

	lo.totalVolume, _ = loTotalVolume.Sub(volumeToSubtract).Float64()
}
