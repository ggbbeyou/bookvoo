package orders

import (
	"github.com/yzimhao/bookvoo/core/base"
	"github.com/yzimhao/bookvoo/user/assets"
)

func NewLimitOrder(user_id int64, trade_symbol string, side OrderSide, price, qty string) (order *TradeOrder, err error) {
	return limit_order(user_id, trade_symbol, side, price, qty)
}

func limit_order(user_id int64, trade_symbol string, side OrderSide, price, qty string) (order *TradeOrder, err error) {
	tp, err := base.GetTradePairBySymbol(trade_symbol)
	if err != nil {
		return nil, err
	}

	//todo 检查交易对限制

	neworder := TradeOrder{
		OrderId:     order_id_by_side(side),
		TradeSymbol: trade_symbol,
		TradingPair: tp.Id,
		OrderSide:   side,
		OrderType:   OrderTypeLimit,
		UserId:      user_id,
		Price:       price,
		Quantity:    qty,
		FinishedQty: "0",
		FeeRate:     tp.FeeRate,
		Status:      OrderStatusNew,
	}

	db := db_engine.NewSession()
	defer db.Close()

	err = db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			db.Rollback()
		} else {
			db.Commit()
		}
	}()

	//冻结相应资产
	if neworder.OrderSide == OrderSideSell {
		//卖单部分fee在订单成交后结算的部分收取
		_, err = assets.FreezeAssets(db, false, user_id, tp.SymbolId, qty, neworder.OrderId, assets.Behavior_Trade)
		if err != nil {
			return nil, err
		}
		neworder.Fee = "0"
		neworder.TradeAmount = "0"
		neworder.TotalAmount = "0"

	} else if neworder.OrderSide == OrderSideBuy {
		//买单的冻结金额加上手续费
		amount := d(price).Mul(d(qty))
		fee := amount.Mul(d(neworder.FeeRate))
		freeze_amount := amount.Add(fee).String()

		neworder.Fee = fee.String()
		neworder.TradeAmount = amount.String()
		neworder.TotalAmount = freeze_amount
		_, err = assets.FreezeAssets(db, false, user_id, tp.StandardSymbolId, freeze_amount, neworder.OrderId, assets.Behavior_Trade)
		if err != nil {
			return nil, err
		}
	}

	if err = neworder.Save(db); err != nil {
		return nil, err
	}

	_, err = db.Table(new(UnfinishedOrder)).Insert(&neworder)
	if err != nil {
		return nil, err
	}
	return &neworder, nil
}
