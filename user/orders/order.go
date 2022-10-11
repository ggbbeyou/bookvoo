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
	OrderSideSell    OrderSide = "sell"
	OrderSideBuy     OrderSide = "buy"
	OrderSideUnknown OrderSide = "unknown"

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
	Id        int64     `xorm:"pk autoincr bigint" json:"-"`
	Symbol    string    `xorm:"-" json:"symbol"`
	PairId    int       `xorm:"notnull index(pair_id)" json:"-"`
	OrderId   string    `xorm:"varchar(30) unique(order_id) notnull" json:"order_id"`
	OrderSide OrderSide `xorm:"varchar(10) index(order_side)" json:"order_side"`
	OrderType OrderType `xorm:"varchar(10)" json:"order_type"` //价格策略，市价单，限价单
	UserId    int64     `xorm:"bigint index(userid) notnull" json:"-"`
	FeeRate   string    `xorm:"decimal(40,20) notnull default(0)" json:"-"`
	//用户委托原始信息
	OriginalPrice    string `xorm:"decimal(40,20) notnull default(0)" json:"original_price"`
	OriginalQuantity string `xorm:"decimal(40,20) notnull default(0)" json:"original_quantity"`
	OriginalAmount   string `xorm:"decimal(40,20) notnull default(0)" json:"-"`
	//根据订单方向不同，冻结的资产也不同
	FreezeAsset string `xorm:"decimal(40,20) notnull default(0)" json:"-"`

	//成交的部分信息
	TradeAvgPrice string `xorm:"decimal(40,20) notnull default(0)" json:"trade_avg_price"`
	TradeQty      string `xorm:"decimal(40,20) notnull default(0)" json:"trade_qty"`
	TradeAmount   string `xorm:"decimal(40,20) notnull default(0)" json:"trade_amount"`

	Fee        string      `xorm:"decimal(40,20) notnull default(0)" json:"-"`
	Status     orderStatus `xorm:"tinyint(1)" json:"status"`
	CreateTime int64       `xorm:"bigint" json:"create_time"` //时间戳 精确到纳秒
	UpdateTime time.Time   `xorm:"timestamp updated" json:"-"`
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
	return fmt.Sprintf("order_%s", to.Symbol)
}

func GetOrderTableName(symbol string) string {
	return fmt.Sprintf("order_%s", symbol)
}
