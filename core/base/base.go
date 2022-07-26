package base

import (
	"time"

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

func GetTradePairById(pair_id int) (*TradePairOpt, error) {
	db := db_engine.NewSession()
	defer db.Close()

	item := TradePairOpt{}
	_, err := db.Table(new(TradePairOpt)).Where("id=?", pair_id).Get(&item)
	return &item, err
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
		return nil, err
	}
	return &item, err
}

func SetDbEngine(db *xorm.Engine) {
	db_engine = db

	db_engine.Sync2(
		new(SymbolInfo),
		new(TradePairOpt),
	)
}
