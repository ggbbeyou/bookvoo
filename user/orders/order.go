package orders

import (
	"fmt"
	"time"

	"xorm.io/xorm"
)

type OrderType string
type OrderSide string
type orderStatus int

const (
	OrderSideSell OrderSide = "sell"
	OrderSideBuy  OrderSide = "buy"

	OrderTypeLimit        OrderType = "limit"
	OrderTypeMarket       OrderType = "market"
	orderTypeMarketQty    OrderType = "market_qty"
	orderTypeMarketAmount OrderType = "market_amount"

	OrderStatusNew    orderStatus = 0
	OrderStatusDone   orderStatus = 1
	OrderStatusCancel orderStatus = 2
)

// 委托记录表
type TradeOrder struct {
	Id          int64       `xorm:"pk autoincr bigint" json:"-"`
	TradeSymbol string      `xorm:"-" json:"symbol"`
	TradingPair int         `xorm:"notnull index(pair_id) index(oa)" json:"-"`
	OrderId     string      `xorm:"varchar(30) unique(order_id) notnull" json:"order_id"`
	OrderSide   OrderSide   `xorm:"varchar(10) index(order_side) index(oa)" json:"order_side"`
	OrderType   OrderType   `xorm:"varchar(10)" json:"order_type"` //价格策略，市价单，限价单
	UserId      int64       `xorm:"bigint index(userid) index(oa) notnull" json:"-"`
	Price       string      `xorm:"decimal(40,20) index(oa) notnull default(0)" json:"price"`
	Quantity    string      `xorm:"decimal(40,20) notnull default(0)" json:"quantity"`
	FinishedQty string      `xorm:"decimal(40,20) notnull default(0)" json:"finished_qty"`
	FeeRate     string      `xorm:"decimal(40,20) notnull default(0)" json:"-"`
	Fee         string      `xorm:"decimal(40,20) notnull default(0)" json:"-"`
	TradeAmount string      `xorm:"decimal(40,20) notnull default(0)" json:"trade_amount"`
	TotalAmount string      `xorm:"decimal(40,20) notnull default(0)" json:"-"` //包含手续费
	Status      orderStatus `xorm:"tinyint(1)" json:"status"`
	CreateTime  int64       `xorm:"bigint" json:"create_time"`
	UpdateTime  time.Time   `xorm:"timestamp updated" json:"-"`
}

func (to *TradeOrder) Save(db *xorm.Session) error {
	//todo 频繁查询表是否存在，后面考虑缓存一下
	exist, err := db.IsTableExist(to.TableName())
	if err != nil {
		return err
	}
	if !exist {
		err := db.CreateTable(to)
		if err != nil {
			return err
		}

		err = db.CreateIndexes(to)
		if err != nil {
			return err
		}

		err = db.CreateUniques(to)
		if err != nil {
			return err
		}
	}

	to.CreateTime = time.Now().UnixNano()
	_, err = db.Table(to).Insert(to)
	if err != nil {
		return err
	}
	return nil
}

func (to *TradeOrder) TableName() string {
	return fmt.Sprintf("order_%s", to.TradeSymbol)
}
