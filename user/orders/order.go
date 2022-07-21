package orders

import (
	"fmt"
	"time"

	"xorm.io/xorm"
)

type orderType int
type orderStatus int
type orderSide int

const (
	OrderSideAsk orderSide = 1
	OrderSideBid orderSide = 2

	OrderTypeLimit  orderType = 1
	OrderTypeMarket orderType = 2

	orderStatusNew  orderStatus = 0
	orderStatusDone orderStatus = 1
)

// 委托记录表
type TradeOrder struct {
	Id            int64       `xorm:"pk autoincr bigint"`
	TradeSymbol   string      `xorm:"-"`
	TradingPair   int         `xorm:"notnull index(pair_id) index(oa)"`
	OrderId       string      `xorm:"varchar(30) unique(order_id) notnull"`
	OrderSide     orderSide   `xorm:"index(order_side) index(oa)"`
	OrderType     orderType   `xorm:"tinyint(1) default(0)"` //价格策略，市价单，限价单
	UserId        int64       `xorm:"bigint index(userid) index(oa) notnull"`
	Price         string      `xorm:"decimal(40,20) index(oa) notnull default(0)"`
	Quantity      string      `xorm:"decimal(40,20) notnull default(0)"`
	UnfinishedQty string      `xorm:"decimal(40,20) notnull default(0)"`
	FeeRate       string      `xorm:"decimal(40,20) notnull default(0)"`
	Fee           string      `xorm:"decimal(40,20) notnull default(0)"`
	TradeAmount   string      `xorm:"decimal(40,20) notnull default(0)"`
	TotalAmount   string      `xorm:"decimal(40,20) notnull default(0)"` //包含手续费
	Status        orderStatus `xorm:"tinyint(1)"`
	CreateTime    int64       `xorm:"bigint"`
	UpdateTime    time.Time   `xorm:"timestamp updated"`
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
	return fmt.Sprintf("trade_order_%s", to.TradeSymbol)
}
