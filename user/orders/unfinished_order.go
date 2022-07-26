package orders

//未完全成交的委托订单记录表
type UnfinishedOrder struct {
	TradeOrder TradeOrder `xorm:"extends"`
}

func (u *UnfinishedOrder) TableName() string {
	return "unfinished_order"
}
