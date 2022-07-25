package base

import (
	"github.com/yzimhao/trading_engine"
	te "github.com/yzimhao/trading_engine"
)

var (
	Engine map[string]*te.TradePair
)

func RunMatching() {
	Engine = make(map[string]*te.TradePair)

	db := db_engine.NewSession()
	defer db.Close()

	rows := []TradePairOpt{}
	db.Table(new(TradePairOpt)).Where("status=?", statusEnable).Find(&rows)

	for _, row := range rows {
		Engine[row.Symbol] = trading_engine.NewTradePair(row.Symbol, row.PricePrec, row.QtyPrec)
	}
}
