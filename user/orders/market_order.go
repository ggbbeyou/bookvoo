package orders

import (
	"github.com/sirupsen/logrus"
	"github.com/yzimhao/bookvoo/core"
	"github.com/yzimhao/bookvoo/core/base"
	"github.com/yzimhao/bookvoo/user/assets"
)

func NewMarketOrderByQty(user_id int64, trade_symbol string, side OrderSide, qty string) (*TradeOrder, error) {
	return market_order_qty(user_id, trade_symbol, side, qty)
}

func market_order_qty(user_id int64, trade_symbol string, side OrderSide, qty string) (order *TradeOrder, err error) {
	tp, err := base.GetTradePairBySymbol(trade_symbol)
	if err != nil {
		return nil, err
	}

	//todo 检查交易对限制

	neworder := TradeOrder{
		OrderId:       order_id_by_side(side),
		TradeSymbol:   trade_symbol,
		TradingPair:   tp.Id,
		OrderSide:     side,
		OrderType:     OrderTypeMarket,
		UserId:        user_id,
		Price:         "-1",
		Quantity:      qty,
		UnfinishedQty: qty,
		FeeRate:       tp.FeeRate,
		Status:        orderStatusNew,
	}

	db := db_engine.NewSession()
	defer db.Close()

	err = db.Begin()
	if err != nil {
		logrus.Error(err, " 22")
		return nil, err
	}

	defer func() {
		if err != nil {
			db.Rollback()
		} else {
			db.Commit()
		}
	}()

	//冻结资产
	if neworder.OrderSide == OrderSideSell {
		_, err = assets.FreezeAssets(db, false, user_id, tp.TradeSymbolId, qty, neworder.OrderId, assets.Behavior_Trade)
		if err != nil {
			return nil, err
		}
		neworder.Fee = "0"
		neworder.TradeAmount = "0"
		neworder.TotalAmount = "0"
	} else if neworder.OrderSide == OrderSideBuy {
		//冻结所有可用
		_, err = assets.FreezeTotalAssets(db, false, user_id, tp.BaseSymbolId, neworder.OrderId, assets.Behavior_Trade)
		if err != nil {
			return nil, err
		}

		freeze, err := assets.QueryFreeze(db, neworder.OrderId)
		if err != nil {
			return nil, err
		}

		neworder.Fee = "0"
		neworder.TradeAmount = core.D(freeze.FreezeAmount).Mul(core.D("1").Sub(core.D(neworder.FeeRate))).String()
		neworder.TotalAmount = freeze.FreezeAmount
	}

	if err = neworder.Save(db); err != nil {
		logrus.Error(err, " 26")
		return nil, err
	}

	return &neworder, nil
}

func NewMarketOrderByAmount(user_id int64, trade_symbol string, side OrderSide, amount string) (order *TradeOrder, err error) {
	return market_order_amount(user_id, trade_symbol, side, amount)
}

func market_order_amount(user_id int64, trade_symbol string, side OrderSide, amount string) (order *TradeOrder, err error) {

	tp, err := base.GetTradePairBySymbol(trade_symbol)
	if err != nil {
		return nil, err
	}

	//todo 检查交易对限制

	neworder := TradeOrder{
		OrderId:       order_id_by_side(side),
		TradeSymbol:   trade_symbol,
		TradingPair:   tp.Id,
		OrderSide:     side,
		OrderType:     OrderTypeMarket,
		UserId:        user_id,
		Price:         "-1",
		Quantity:      "0",
		UnfinishedQty: "0",
		FeeRate:       tp.FeeRate,
		Status:        orderStatusNew,
	}

	db := db_engine.NewSession()
	defer db.Close()

	err = db.Begin()
	if err != nil {
		logrus.Error(err, " 22")
		return nil, err
	}

	defer func() {
		if err != nil {
			db.Rollback()
		} else {
			db.Commit()
		}
	}()

	if neworder.OrderSide == OrderSideSell {
		_, err = assets.FreezeTotalAssets(db, false, user_id, tp.TradeSymbolId, neworder.OrderId, assets.Behavior_Trade)
		if err != nil {
			return nil, err
		}

		neworder.Fee = "0"
		neworder.TradeAmount = core.D(amount).Mul(core.D("1").Sub(core.D(neworder.FeeRate))).String()
		neworder.TotalAmount = amount
	} else if neworder.OrderSide == OrderSideBuy {
		_, err = assets.FreezeAssets(db, false, user_id, tp.BaseSymbolId, amount, neworder.OrderId, assets.Behavior_Trade)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}
