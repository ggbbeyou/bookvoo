package orders

import "time"

type OrderType int
type OrderStatus int
type OrderSide int

const (
	OrderSideAsk OrderSide = 0
	OrderSideBid OrderSide = 1

	OrderTypeLimit  OrderType = 0
	OrderTypeMarket OrderType = 1

	OrderStatusNew  OrderStatus = 0
	OrderStatusDone OrderStatus = 1
)

// 委托记录表
type TradeOrder struct {
	Id            int64       `xorm:"pk autoincr bigint"`
	TradingPair   int         `xorm:"notnull index(pair_id) index(oa)"`
	OrderId       string      `xorm:"varchar(30) unique(order_id) notnull"`
	OrderSide     OrderSide   `xorm:"index(order_side) index(oa)"`
	OrderType     OrderType   `xorm:"tinyint(1) default(0)"` //价格策略，市价单，限价单
	UserId        int64       `xorm:"bigint index(userid) index(oa) notnull"`
	Price         string      `xorm:"decimal(40,20) index(oa) notnull default(0)"`
	Quantity      string      `xorm:"decimal(40,20) notnull default(0)"`
	UnfinishedQty string      `xorm:"decimal(40,20) notnull default(0)"`
	FeeRate       string      `xorm:"decimal(40,20) notnull default(0)"`
	Fee           string      `xorm:"decimal(40,20) notnull default(0)"`
	TradeAmount   string      `xorm:"decimal(40,20) notnull default(0)"`
	TotalAmount   string      `xorm:"decimal(40,20) notnull default(0)"`
	Status        OrderStatus `xorm:"tinyint(1)"`
}

//未完全成交的委托订单记录表
type UnfinishedOrder struct {
	TradeOrder TradeOrder `xorm:"extends"`
}

// 成交记录表
type TradeRecord struct {
	Id          int64 `xorm:"pk autoincr bigint"`
	TradingPair int   `xorm:"notnull"`

	Ask         string `xorm:"varchar(30) unique(trade)"`
	Bid         string `xorm:"varchar(30) unique(trade)"`
	TradeId     string `xorm:"varchar(30) unique(trade)"`
	TradeByType int8   `xorm:"tinyint(1)"`
	AskUid      int64  `xorm:"bigint notnull"`
	Biduid      int64  `xorm:"bigint notnull"`
	Price       string `xorm:"decimal(40,20) notnull default(0)"`
	Quantity    string `xorm:"decimal(40,20) notnull default(0)"`

	AskFeeRate string `xorm:"decimal(40,20) notnull default(0)"`
	AskFee     string `xorm:"decimal(40,20) notnull default(0)"`

	BidFeeRate string `xorm:"decimal(40,20) notnull default(0)"`
	BidFee     string `xorm:"decimal(40,20) notnull default(0)"`

	CreateTime time.Time `xorm:"timestamp created"`
	UpdateTime time.Time `xorm:"timestamp updated"`
}
