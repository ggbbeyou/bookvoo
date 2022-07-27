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

func (to *TradeRecord) Save(db *xorm.Session) error {
	if to.Symbol == "" {
		return fmt.Errorf("symbol not set")
	}
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

	_, err = db.Table(to).Insert(to)
	if err != nil {
		return err
	}
	return nil
}

func (to *TradeRecord) TableName() string {
	return fmt.Sprintf("order_trade_%s", to.Symbol)
}
func GetTradeRecordTableName(symbol string) string {
	return fmt.Sprintf("order_trade_%s", symbol)
}
