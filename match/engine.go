package match

import (
	"github.com/sirupsen/logrus"
	"github.com/yzimhao/bookvoo/core/base"
	te "github.com/yzimhao/trading_engine"
	"xorm.io/xorm"
)

var (
	Engine    map[string]*te.TradePair
	db_engine *xorm.Engine
)

func SetDbEngine(db *xorm.Engine) {
	db_engine = db
	base.SetDbEngine(db)
}

func RunMatching() {
	Engine = make(map[string]*te.TradePair)

	db := db_engine.NewSession()
	defer db.Close()

	rows := []base.TradePairOpt{}
	db.Table(new(base.TradePairOpt)).Where("status=?", base.StatusEnable).Find(&rows)

	for _, row := range rows {

		Engine[row.Symbol] = te.NewTradePair(row.Symbol, row.PricePrec, row.QtyPrec)

		go func(item base.TradePairOpt) {
			for {
				select {
				case result := <-Engine[item.Symbol].ChTradeResult:
					logrus.Error(result)
				case cancel := <-Engine[item.Symbol].ChCancelResult:
					logrus.Error(cancel)
				}
			}
		}(row)
	}
}
