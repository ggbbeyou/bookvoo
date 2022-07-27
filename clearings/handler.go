package clearings

import (
	"github.com/yzimhao/bookvoo/core/base"
	"github.com/yzimhao/bookvoo/user/orders"
	"github.com/yzimhao/trading_engine"
	"xorm.io/xorm"
)

var (
	db_engine *xorm.Engine
	Notify    chan trading_engine.TradeResult
)

func SetDbEngine(db *xorm.Engine) {
	db_engine = db
}

//结算一条成交记录
func NewClearing(symbol string, ask_id, bid_id string, price, qty string) (err error) {
	tradeInfo, err := base.GetTradePairBySymbol(symbol)
	if err != nil {
		return err
	}

	db := db_engine.NewSession()
	defer db.Close()

	err = db.Begin()
	if err != nil {
		return err
	}
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
