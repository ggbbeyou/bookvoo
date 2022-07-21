package orders

import (
	"time"
)

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
