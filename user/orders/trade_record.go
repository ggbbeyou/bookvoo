package orders

import (
	"fmt"
	"time"

	"xorm.io/xorm"
)

type TradeBy int

const (
	TradeBySell TradeBy = 1
	TradeByBuy  TradeBy = 2
)

// 成交记录表
type TradeRecord struct {
	Id     int64  `xorm:"pk autoincr bigint"`
	Symbol string `xorm:"-"`

	TradeId string `xorm:"varchar(30) unique(trade_id)"`
	Ask     string `xorm:"varchar(30) unique(trade)"`
	Bid     string `xorm:"varchar(30) unique(trade)"`

	TradeBy  TradeBy `xorm:"tinyint(1)"`
	AskUid   int64   `xorm:"bigint notnull"`
	Biduid   int64   `xorm:"bigint notnull"`
	Price    string  `xorm:"decimal(40,20) notnull default(0)"`
	Quantity string  `xorm:"decimal(40,20) notnull default(0)"`
	Amount   string  `xorm:"decimal(40,20) notnull default(0)"`

	AskFeeRate string `xorm:"decimal(40,20) notnull default(0)"`
	AskFee     string `xorm:"decimal(40,20) notnull default(0)"`

	BidFeeRate string `xorm:"decimal(40,20) notnull default(0)"`
	BidFee     string `xorm:"decimal(40,20) notnull default(0)"`

	CreateTime time.Time `xorm:"timestamp created"`
	UpdateTime time.Time `xorm:"timestamp updated"`
}

func (t *TradeRecord) TableName() string {
	return fmt.Sprintf("order_trade_record_%s", t.Symbol)
}

func (t *TradeRecord) Save(db *xorm.Session) error {
	_, err := db.Table(new(TradeRecord)).Insert(t)
	if err != nil {
		return err
	}
	return nil
}
