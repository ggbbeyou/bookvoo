package orders

import (
	"fmt"
	"time"

	"github.com/yzimhao/bookvoo/common/types"
	"xorm.io/xorm"
)

type TradeBy int

const (
	TradeBySell TradeBy = 1
	TradeByBuy  TradeBy = 2
)

// 成交记录表
type TradeRecord struct {
	Id     int64  `xorm:"pk autoincr bigint" json:"-"`
	Symbol string `xorm:"-" json:"-"`

	TradeId string `xorm:"varchar(30) unique(trade_id)" json:"trade_id"`
	Ask     string `xorm:"varchar(30) unique(trade)" json:"ask"`
	Bid     string `xorm:"varchar(30) unique(trade)" json:"bid"`

	TradeBy  TradeBy `xorm:"tinyint(1)" json:"trade_by"`
	AskUid   int64   `xorm:"bigint notnull" json:"-"`
	BidUid   int64   `xorm:"bigint notnull" json:"-"`
	Price    string  `xorm:"decimal(40,20) notnull default(0)" json:"price"`
	Quantity string  `xorm:"decimal(40,20) notnull default(0)" json:"quantity"`
	Amount   string  `xorm:"decimal(40,20) notnull default(0)" json:"amount"`

	AskFeeRate string `xorm:"decimal(40,20) notnull default(0)" json:"-"`
	AskFee     string `xorm:"decimal(40,20) notnull default(0)" json:"-"`

	BidFeeRate string `xorm:"decimal(40,20) notnull default(0)" json:"-"`
	BidFee     string `xorm:"decimal(40,20) notnull default(0)" json:"-"`

	CreateTime types.Time `xorm:"timestamp created" json:"create_time"`
	UpdateTime time.Time  `xorm:"timestamp updated" json:"-"`
}

func (tr *TradeRecord) Save(db *xorm.Session) error {
	if tr.Symbol == "" {
		return fmt.Errorf("symbol not set")
	}
	//todo 频繁查询表是否存在，后面考虑缓存一下
	exist, err := db.IsTableExist(tr.TableName())
	if err != nil {
		return err
	}
	if !exist {
		err := db.CreateTable(tr)
		if err != nil {
			return err
		}

		err = db.CreateIndexes(tr)
		if err != nil {
			return err
		}

		err = db.CreateUniques(tr)
		if err != nil {
			return err
		}
	}

	_, err = db.Table(tr).Insert(tr)
	if err != nil {
		return err
	}
	return nil
}

func (tr *TradeRecord) TableName() string {
	return tr.GetTableName(tr.Symbol)
}
func (tr *TradeRecord) GetTableName(symbol string) string {
	return fmt.Sprintf("order_%s_traderecord", symbol)
}
