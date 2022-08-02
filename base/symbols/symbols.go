package symbols

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/shopspring/decimal"
	"xorm.io/xorm"
)

var (
	db_engine *xorm.Engine
)

type status int

const (
	StatusDisable status = 0
	StatusEnable  status = 1
)

type SymbolInfo struct {
	Id           int       `xorm:"pk autoincr int"`
	Symbol       string    `xorm:"varchar(100) notnull unique(symbol)"`
	Name         string    `xorm:"varchar(250) notnull"`
	ShowPrec     int       `xorm:"default(0)"`
	MinPrecision int       `xorm:"default(0)"`
	Standard     bool      `xorm:"default(0)"` //是否为本位币
	Status       status    `xorm:"default(0) notnull"`
	CreateTime   time.Time `xorm:"timestamp created"`
	UpdateTime   time.Time `xorm:"timestamp updated"`
}

type TradePairOpt struct {
	Id     int    `xorm:"pk autoincr int"`
	Symbol string `xorm:"varchar(100) notnull unique(symbol)"`
	Name   string `xorm:"varchar(250) notnull"`

	SymbolId         int `xorm:"default(0) unique(symbol_base)"` //交易物品
	StandardSymbolId int `xorm:"default(0) unique(symbol_base)"` //支付货币

	PricePrec      int    `xorm:"default(2)"`
	QtyPrec        int    `xorm:"default(0)"`
	AllowMinQty    string `xorm:"decimal(40,20) notnull"`
	AllowMaxQty    string `xorm:"decimal(40,20) notnull"`
	AllowMinAmount string `xorm:"decimal(40,20) notnull"`
	AllowMaxAmount string `xorm:"decimal(40,20) notnull"`
	FeeRate        string `xorm:"decimal(40,20) notnull default(0)"`

	Status     status    `xorm:"default(0) notnull"`
	CreateTime time.Time `xorm:"timestamp created"`
	UpdateTime time.Time `xorm:"timestamp updated"`
}

func (t *TradePairOpt) TableName() string {
	return "trade_pair_option"
}

func (t *TradePairOpt) FormatAmount(a string) string {
	q, _ := decimal.NewFromString(a)
	return q.StringFixedBank(int32(t.PricePrec))
}

func (t *TradePairOpt) FormatQty(qty string) string {
	q, _ := decimal.NewFromString(qty)
	return q.StringFixedBank(int32(t.QtyPrec))
}

func GetTradePairBySymbol(symbol string) (*TradePairOpt, error) {
	db := db_engine.NewSession()
	defer db.Close()

	item := TradePairOpt{}
	has, err := db.Table(new(TradePairOpt)).Where("symbol=?", symbol).Get(&item)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, fmt.Errorf("not found trade pair %s", symbol)
	}
	return &item, err
}

func Init(db *xorm.Engine, rdc *redis.Client) {
	db_engine = db

	db_engine.Sync2(
		new(SymbolInfo),
		new(TradePairOpt),
	)
}
