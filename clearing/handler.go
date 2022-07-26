package clearing

import "xorm.io/xorm"

var (
	db_engine *xorm.Engine
)

func SetDbEngine(db *xorm.Engine) {
	db_engine = db
}

//结算一条成交记录
func NewClearing(symbol string, ask_id, bid_id string, price, qty string) (err error) {
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

		ask_order_id: ask_id,
		bid_order_id: bid_id,
		trade_price:  d(price),
		trade_qty:    d(qty),
		trade_amount: d(price).Mul(d(qty)),
	}
	//检查双方订单状态
	err = cl.check()
	if err != nil {
		return err
	}

	//修改卖方订单信息
	err = cl.updateAsk()
	if err != nil {
		return err
	}
	//修改买方订单信息
	err = cl.updateBid()
	if err != nil {
		return err
	}
	//写成交日志
	err = cl.tradeRecord()
	if err != nil {
		return err
	}

	//结算三方资产

	return nil
}
