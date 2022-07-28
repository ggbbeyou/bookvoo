package clearings

import (
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/yzimhao/bookvoo/base/symbols"
	"github.com/yzimhao/bookvoo/user/orders"
	te "github.com/yzimhao/trading_engine"
	"xorm.io/xorm"
)

var (
	db_engine *xorm.Engine
	Notify    chan te.TradeResult
	rdc       *redis.Client
)

func Init(db *xorm.Engine, r *redis.Client) {
	db_engine = db
	rdc = r
}

func Run() {
	Notify = make(chan te.TradeResult, 1000)
	for {
		if data, ok := <-Notify; ok {
			go NewClearing(data.Symbol, data.AskOrderId, data.BidOrderId, data.TradePrice.String(), data.TradeQuantity.String())
		}
	}
}

//结算一条成交记录
func NewClearing(symbol string, ask_id, bid_id string, price, qty string) (err error) {
	logrus.Infof("[clearings] %s %s %s %s %s", symbol, ask_id, bid_id, price, qty)

	tradeInfo, err := symbols.GetTradePairBySymbol(symbol)
	if err != nil {
		return err
	}

	db := db_engine.NewSession()
	defer db.Close()

	err = db.Begin()
	if err != nil {
		return err
	}

	//todo lock 双方订单

	defer func() {
		if err != nil {
			db.Rollback()
		} else {
			db.Commit()
		}
	}()

	cl := clearing{
		db:     db,
		symbol: symbol,

		symbol_id:          tradeInfo.SymbolId,
		standard_symbol_id: tradeInfo.StandardSymbolId,

		ask_order_id: ask_id,
		bid_order_id: bid_id,
		trade_price:  d(price),
		trade_qty:    d(qty),
		trade_amount: d(price).Mul(d(qty)),

		ask:    new(orders.TradeOrder),
		bid:    new(orders.TradeOrder),
		record: new(orders.TradeRecord),
	}
	//检查双方订单状态
	err = cl.check()
	if err != nil {
		return err
	}

	//写成交日志
	err = cl.tradeRecord()
	if err != nil {
		return err
	}

	//修改买方订单信息
	err = cl.updateBid()
	if err != nil {
		return err
	}

	//修改卖方订单信息
	err = cl.updateAsk()
	if err != nil {
		return err
	}

	//结算三方资产
	err = cl.transfer()
	if err != nil {
		return err
	}
	return nil
}
