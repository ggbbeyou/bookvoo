package orders

var ChCancel chan TradeOrder

func service() {
	ChCancel = make(chan TradeOrder)
	for {
		if data, ok := <-ChCancel; ok {
			cancel_order(data.Symbol, data.OrderId)
		}
	}
}
