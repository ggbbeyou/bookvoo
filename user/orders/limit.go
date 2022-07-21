package orders

import (
	"github.com/yzimhao/bookvoo/core"
	"github.com/yzimhao/bookvoo/core/base"
	"github.com/yzimhao/bookvoo/user/assets"
)

func trade_limit_order(user_id int64, trade_symbol string, side orderSide, price, qty, fee_rate string) (order_id string, err error) {
	tp, err := base.GetTradePairBySymbol(trade_symbol)
	if err != nil {
		return "", err
	}

	//todo 检查交易对限制

	neworder := TradeOrder{
		OrderId:       order_id_by_side(side),
		TradeSymbol:   trade_symbol,
		TradingPair:   tp.Id,
		OrderSide:     side,
		OrderType:     orderTypeLimit,
		UserId:        user_id,
		Price:         price,
		Quantity:      qty,
		UnfinishedQty: qty,
		FeeRate:       fee_rate,
		Status:        orderStatusNew,
	}

	db := db_engine.NewSession()
	defer db.Close()

	//todo 开启事务
	db.Begin()
	defer func() {
		if err != nil {
			db.Rollback()
		} else {
			db.Commit()
		}
	}()

	//冻结相应资产
	if neworder.OrderSide == orderSideAsk {

		_, err = assets.FreezeAssets(db, false, user_id, tp.TradeSymbolId, qty, neworder.OrderId, assets.Behavior_Trade)
		if err != nil {
			return "", err
		}
	} else if neworder.OrderSide == orderSideBid {
		//买单的冻结金额
		amount := core.D(price).Mul(core.D(qty))
		fee := amount.Mul(core.D(fee_rate))
		freeze_amount := amount.Add(fee).String()

		neworder.Fee = fee.String()
		neworder.TradeAmount = amount.String()
		neworder.TotalAmount = freeze_amount
		_, err = assets.FreezeAssets(db, false, user_id, tp.BaseSymbolId, freeze_amount, neworder.OrderId, assets.Behavior_Trade)
		if err != nil {
			return "", err
		}
	}

	//save order
	if err = neworder.Save(db); err != nil {
		return "", err
	}

	_, err = db.Table(new(UnfinishedOrder)).Insert(&neworder)
	if err != nil {
		return "", err
	}
	return neworder.OrderId, nil
}
