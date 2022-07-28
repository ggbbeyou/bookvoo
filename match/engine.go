package match

import (
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/yzimhao/bookvoo/base/symbols"
	"github.com/yzimhao/bookvoo/clearings"
	te "github.com/yzimhao/trading_engine"
	"xorm.io/xorm"
)

var (
	Engine    map[string]*te.TradePair
	db_engine *xorm.Engine
)

func Init(db *xorm.Engine, rdc *redis.Client) {
	db_engine = db
}

func RunMatching() {
	Engine = make(map[string]*te.TradePair)

	db := db_engine.NewSession()
	defer db.Close()

	rows := []symbols.TradePairOpt{}
	db.Table(new(symbols.TradePairOpt)).Where("status=?", symbols.StatusEnable).Find(&rows)

	for _, row := range rows {
		Engine[row.Symbol] = te.NewTradePair(row.Symbol, row.PricePrec, row.QtyPrec)
		go func(item symbols.TradePairOpt) {
			for {
				select {
				case result := <-Engine[item.Symbol].ChTradeResult:
					logrus.Debugf("[tradeResult] %v", result)
					clearings.Notify <- result
				case cancel := <-Engine[item.Symbol].ChCancelResult:
					logrus.Debugf("[cancelOrder] %v", cancel)
				}
			}
		}(row)
	}
}
