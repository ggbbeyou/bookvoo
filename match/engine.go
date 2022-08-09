package match

import (
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/yzimhao/bookvoo/base/symbols"
	"github.com/yzimhao/bookvoo/clearings"
	"github.com/yzimhao/bookvoo/user/orders"
	te "github.com/yzimhao/trading_engine"
	"xorm.io/xorm"
)

var (
	Send chan *orders.TradeOrder

	db_engine *xorm.Engine
	Engine    *engine
)

type engine struct {
	Symbols map[string]*te.TradePair
	sync.Mutex
}

func Init(db *xorm.Engine, rdc *redis.Client) {
	db_engine = db
	Send = make(chan *orders.TradeOrder)
	Engine = new(engine)
}

func Run() {
	Engine.init()
	Engine.service()
	Engine.handler()
	Engine.rebuild()
}

func (e *engine) rebuild() {
	db := orders.Db().NewSession()
	defer db.Close()

	for symbol, _ := range e.Symbols {

		tp, _ := symbols.GetExchangeBySymbol(symbol)
		rows := []orders.TradeOrder{}
		db.Table(new(orders.UnfinishedOrder)).Where("pair_id=?", tp.Id).OrderBy("create_time asc").Find(&rows)
		for _, row := range rows {
			row.Symbol = symbol
			Send <- &row
		}
	}
}

func (e *engine) init() {
	e.Lock()
	defer e.Unlock()

	e.Symbols = make(map[string]*te.TradePair)

	db := db_engine.NewSession()
	defer db.Close()

	rows := []symbols.Exchange{}
	db.Table(new(symbols.Exchange)).Where("status=?", symbols.StatusEnable).Find(&rows)
	for _, row := range rows {
		e.Symbols[row.Symbol] = te.NewTradePair(row.Symbol, row.PricePrec, row.QtyPrec)
	}
}

func (e *engine) service() {
	for symbol, item := range e.Symbols {
		go func(symbol string, obj *te.TradePair) {
			for {
				select {
				case result := <-obj.ChTradeResult:
					logrus.Debugf("[tradeResult] %s %v", symbol, result)
					clearings.Notify <- result
				case cancel := <-obj.ChCancelResult:
					logrus.Debugf("[cancelOrder] %s %v", symbol, cancel)
				}
			}
		}(symbol, item)
	}
}

func (e *engine) handler() {
	go func() {
		for {
			select {
			case data := <-Send:
				func() {
					e.Lock()
					defer e.Unlock()

					if data.OrderType == orders.OrderTypeLimit {
						if data.OrderSide == orders.OrderSideSell {
							e.Symbols[data.Symbol].ChNewOrder <- te.NewAskLimitItem(data.OrderId, d(data.Price), d(data.Quantity), data.CreateTime)
						} else if data.OrderSide == orders.OrderSideBuy {
							e.Symbols[data.Symbol].ChNewOrder <- te.NewBidLimitItem(data.OrderId, d(data.Price), d(data.Quantity), data.CreateTime)
						}
					}

				}()
			}

		}
	}()
}

func (e *engine) Get(symbol string) (*te.TradePair, error) {
	e.Lock()
	defer e.Unlock()

	if _, ok := e.Symbols[symbol]; !ok {
		return nil, fmt.Errorf("invalid symbol")
	}
	return e.Symbols[symbol], nil
}

func d(ss string) decimal.Decimal {
	s, _ := decimal.NewFromString(ss)
	return s
}
