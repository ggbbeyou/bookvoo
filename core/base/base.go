package base

import (
	"time"

	"github.com/yzimhao/trading_engine"
	"xorm.io/xorm"
)

var (
	db_engine *xorm.Engine

	MatchingEngine map[string]*trading_engine.TradePair
)

type status int

const (
	statusDisable status = 0
	statusEnable  status = 1
)

type SymbolInfo struct {
	Id        int    `xorm:"pk autoincr int"`
	Symbol    string `xorm:"varchar(100) notnull unique(symbol)"`
	Name      string `xorm:"varchar(250) notnull"`
	Precision int    `xorm:"default(0)"`

	Status     status    `xorm:"default(0) notnull"`
	CreateTime time.Time `xorm:"timestamp created"`
	UpdateTime time.Time `xorm:"timestamp updated"`
}

type TradePairOpt struct {
	Id     int    `xorm:"pk autoincr int"`
	Symbol string `xorm:"varchar(100) notnull unique(symbol)"`
	Name   string `xorm:"varchar(250) notnull"`

	TradeSymbolId int `xorm:"default(0) unique(symbol_base)"`
	BaseSymbolId  int `xorm:"default(0) unique(symbol_base)"`

	PricePerc      int    `xorm:"default(2)"`
	QtyPerc        int    `xorm:"default(0)"`
	AllowMinQty    string `xorm:"decimal(40,20) notnull"`
	AllowMaxQty    string `xorm:"decimal(40,20) notnull"`
	AllowMinAmount string `xorm:"decimal(40,20) notnull"`
	AllowMaxAmount string `xorm:"decimal(40,20) notnull"`
	FeeRate        string `xorm:"decimal(40,20) notnull default(0)"`

	Status     status    `xorm:"default(0) notnull"`
	CreateTime time.Time `xorm:"timestamp created"`
	UpdateTime time.Time `xorm:"timestamp updated"`
}

func GetTradePairById(pair_id int) *TradePairOpt {
	db := db_engine.NewSession()
	defer db.Close()

	item := TradePairOpt{}
	db.Table(new(TradePairOpt)).Where("id=?", pair_id).Get(&item)
	return &item
}

func GetTradePairBySymbol(symbol string) *TradePairOpt {
	db := db_engine.NewSession()
	defer db.Close()

	item := TradePairOpt{}
	db.Table(new(TradePairOpt)).Where("symbol=?", symbol).Get(&item)
	return &item
}

func RunMatching() {
	MatchingEngine = make(map[string]*trading_engine.TradePair)

	trade_symbol := "demo"
	MatchingEngine[trade_symbol] = trading_engine.NewTradePair(trade_symbol, 2, 0)
}

func SetDbEngine(db *xorm.Engine) {
	db_engine = db

	db_engine.Sync2(
		new(SymbolInfo),
		new(TradePairOpt),
	)
}
