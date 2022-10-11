package match

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/yzimhao/bookvoo/base/symbols"
	"github.com/yzimhao/bookvoo/clearing"
	"github.com/yzimhao/bookvoo/user/orders"
	te "github.com/yzimhao/trading_engine"
	"xorm.io/xorm"
)

var (
	Send chan orders.TradeOrder

	db_engine *xorm.Engine
	Engine    *engine
)

type engine struct {
	symbols sync.Map
	sync.Mutex

	wgRebuild sync.WaitGroup
}

func Init(db *xorm.Engine, rdc *redis.Client) {
	db_engine = db
	Send = make(chan orders.TradeOrder, 10000)
	Engine = new(engine)
}

func Run() {

	Engine.init()
	Engine.wgRebuild.Add(1)
	Engine.service()
	Engine.handler()
	Engine.rebuild()
	logrus.Info("[match] run4")
}

func (e *engine) rebuild() {
	defer e.wgRebuild.Done()

	db := orders.Db().NewSession()
	defer db.Close()

	e.Lock()
	defer e.Unlock()

	e.symbols.Range(func(key, value any) bool {
		symbol := key.(string)
		tp, _ := symbols.GetExchangeBySymbol(symbol)
		rows := []orders.TradeOrder{}
		db.Table(new(orders.UnfinishedOrder)).Where("pair_id=?", tp.Id).OrderBy("create_time asc").Find(&rows)
		for i, row := range rows {
			row.Symbol = symbol
			//rebuild的时候总下单数量减去已经成交的重新加载到撮合
			row.OriginalQuantity = d(row.OriginalQuantity).Sub(d(row.TradeQty)).String()
			logrus.Infof("[match] rebuild (%d) %s", i, row.OrderId)
			Send <- row
		}

		return true
	})

}

func (e *engine) init() {

	db := db_engine.NewSession()
	defer db.Close()

	rows := []symbols.Exchange{}
	db.Table(new(symbols.Exchange)).Where("status=?", symbols.StatusEnable).Find(&rows)
	for _, row := range rows {
		e.symbols.Store(row.Symbol, te.NewTradePair(row.Symbol, row.PricePrec, row.QtyPrec))
	}
}

func (e *engine) service() {
	e.symbols.Range(func(k, v any) bool {
		go func(symbol string, obj *te.TradePair) {
			for {
				select {
				case result := <-obj.ChTradeResult:
					logrus.Infof("[match] %s ask: %s bid: %s price: %s vol: %s", symbol, result.AskOrderId, result.BidOrderId, result.TradePrice.String(), result.TradeQuantity.String())
					clearing.Notify <- result
				case order_id := <-obj.ChCancelResult:
					logrus.Infof("[match] %s 取消订单 %s", symbol, order_id)
					if !clearing.ClearingLockExist(order_id) {
						time.Sleep(time.Duration(1) * time.Second)
						//取消订单之前，需要确认订单释放已经有锁
						// orders.ChCancel <- orders.TradeOrder{
						// 	Symbol:  symbol,
						// 	OrderId: order_id,
						// }
					} else {
						logrus.Warnf("[match] 订单结算锁还未释放 %s", order_id)
					}
				}
			}
		}(k.(string), v.(*te.TradePair))
		return true
	})

}

func (e *engine) handler() {
	go func() {
		for {
			e.wgRebuild.Wait()

			select {
			case data := <-Send:
				logrus.Infof("[match] handler order: %s", data.OrderId)
				go func() {
					t, err := e.Get(data.Symbol)
					if err != nil {
						return
					}
					if data.OrderType == orders.OrderTypeLimit {
						if data.OrderSide == orders.OrderSideSell {
							t.ChNewOrder <- te.NewAskLimitItem(data.OrderId, d(data.OriginalPrice), d(data.OriginalQuantity), data.CreateTime)
						} else if data.OrderSide == orders.OrderSideBuy {
							t.ChNewOrder <- te.NewBidLimitItem(data.OrderId, d(data.OriginalPrice), d(data.OriginalQuantity), data.CreateTime)
						}
					}
					//todo 市价单

				}()
			}

		}
	}()
}

func (e *engine) Get(symbol string) (*te.TradePair, error) {
	v, ok := e.symbols.Load(symbol)
	if !ok {
		return nil, fmt.Errorf("%s tradepair not found", symbol)
	}
	return v.(*te.TradePair), nil
}

func (e *engine) Foreach(a func(k string, v *te.TradePair)) {
	e.symbols.Range(func(key, value any) bool {
		a(key.(string), value.(*te.TradePair))
		return true
	})
}

func d(ss string) decimal.Decimal {
	s, _ := decimal.NewFromString(ss)
	return s
}
