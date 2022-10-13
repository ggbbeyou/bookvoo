package orders

import (
	"github.com/yzimhao/bookvoo/base/symbols"
	"github.com/yzimhao/bookvoo/user/assets"
)

func NewMarketOrderByQty(user_id int64, trade_symbol string, side OrderSide, qty string) (*TradeOrder, error) {
	return market_order_qty(user_id, trade_symbol, side, qty)
}

func market_order_qty(user_id int64, trade_symbol string, side OrderSide, qty string) (order *TradeOrder, err error) {
	tp, err := symbols.GetExchangeBySymbol(trade_symbol)
	if err != nil {
		return nil, err
	}

	//todo 检查交易对限制

	neworder := TradeOrder{
		OrderId:          order_id_by_side(side),
		Symbol:           trade_symbol,
		PairId:           tp.Id,
		OrderSide:        side,
		OrderType:        OrderTypeMarket,
		UserId:           user_id,
		OriginalPrice:    "0",
		OriginalQuantity: qty,
		OriginalAmount:   "0",
		TradeAvgPrice:    "0",
		TradeQty:         "0",
		TradeAmount:      "0",

		Fee:         "0",
		FreezeAsset: "0",

		FeeRate: string(tp.FeeRate),
		Status:  OrderStatusNew,
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

	//冻结资产
	if neworder.OrderSide == OrderSideSell {
		_, err = assets.FreezeAssets(db, false, user_id, tp.TargetSymbolId, qty, neworder.OrderId, assets.Behavior_Trade)
		if err != nil {
			return nil, err
		}
		neworder.FreezeAsset = qty
	} else if neworder.OrderSide == OrderSideBuy {
		//冻结所有可用
		_, err = assets.FreezeTotalAssets(db, false, user_id, tp.StandardSymbolId, neworder.OrderId, assets.Behavior_Trade)
		if err != nil {
			return nil, err
		}

		freeze, err := assets.QueryFreeze(db, neworder.OrderId)
		if err != nil {
			return nil, err
		}
		neworder.FreezeAsset = freeze.FreezeAmount
	}

	if err = neworder.Save(db); err != nil {
		return nil, err
	}

	return &neworder, nil
}

func NewMarketOrderByAmount(user_id int64, trade_symbol string, side OrderSide, amount string) (order *TradeOrder, err error) {
	return market_order_amount(user_id, trade_symbol, side, amount)
}

func market_order_amount(user_id int64, trade_symbol string, side OrderSide, amount string) (order *TradeOrder, err error) {

	tp, err := symbols.GetExchangeBySymbol(trade_symbol)
	if err != nil {
		return nil, err
	}

	//todo 检查交易对限制

	neworder := TradeOrder{
		OrderId:          order_id_by_side(side),
		Symbol:           trade_symbol,
		PairId:           tp.Id,
		OrderSide:        side,
		OrderType:        OrderTypeMarket,
		UserId:           user_id,
		OriginalPrice:    "0",
		OriginalQuantity: "0",

		TradeAvgPrice: "0",
		TradeQty:      "0",
		TradeAmount:   "0",
		Fee:           "0",
		FreezeAsset:   "0",
		FeeRate:       string(tp.FeeRate),
		Status:        OrderStatusNew,
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

	if neworder.OrderSide == OrderSideSell {
		_, err = assets.FreezeTotalAssets(db, false, user_id, tp.TargetSymbolId, neworder.OrderId, assets.Behavior_Trade)
		if err != nil {
			return nil, err
		}

		freeze, err := assets.QueryFreeze(db, neworder.OrderId)
		if err != nil {
			return nil, err
		}
		neworder.FreezeAsset = freeze.FreezeAmount

	} else if neworder.OrderSide == OrderSideBuy {
		_, err = assets.FreezeAssets(db, false, user_id, tp.StandardSymbolId, amount, neworder.OrderId, assets.Behavior_Trade)
		if err != nil {
			return nil, err
		}
		neworder.FreezeAsset = amount
	}

	if err = neworder.Save(db); err != nil {
		return nil, err
	}

	return &neworder, nil
}
