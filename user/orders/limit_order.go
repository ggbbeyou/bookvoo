package orders

import (
	"fmt"

	"github.com/yzimhao/bookvoo/base/symbols"
	"github.com/yzimhao/bookvoo/user/assets"
)

//限价订单
func NewLimitOrder(user_id int64, trade_symbol string, side OrderSide, price, qty string) (order *TradeOrder, err error) {
	return limit_order(user_id, trade_symbol, side, price, qty)
}

//检查同用户的反向订单价格和新订单是否能匹配成交，防止买卖双方为同一个用户的成交
// true 没有同用户的买卖订单成交的情况，允许本次下单
// false 存在同用户的买卖订单成交，不允许本次下单
func is_allow_open_reverse_order(user_id int64, pair_id int, cur_side OrderSide, cur_price string) bool {
	db := db_engine.NewSession()
	defer db.Close()

	un_order := new(UnfinishedOrder)
	var has bool
	if cur_side == OrderSideSell {
		has, _ = db.Table(un_order.TableName()).Where("user_id=? and pair_id=? and order_side!=?", user_id, pair_id, cur_side).Where("original_price >= ?", cur_price).Get(un_order)
	} else {
		has, _ = db.Table(un_order.TableName()).Where("user_id=? and pair_id=? and order_side!=?", user_id, pair_id, cur_side).Where("original_price <= ?", cur_price).Get(un_order)
	}
	if has {
		return false
	}
	return true
}

func limit_order(user_id int64, trade_symbol string, side OrderSide, price, qty string) (order *TradeOrder, err error) {
	tp, err := symbols.GetPairBySymbol(trade_symbol)
	if err != nil {
		return nil, err
	}

	//todo 检查交易对限制
	if !is_allow_open_reverse_order(user_id, tp.Id, side, price) {
		return nil, fmt.Errorf("有反向订单未成交价格冲突，不允许下单")
	}

	neworder := TradeOrder{
		OrderId:          order_id_by_side(side),
		Symbol:           trade_symbol,
		PairId:           tp.Id,
		OrderSide:        side,
		OrderType:        OrderTypeLimit,
		UserId:           user_id,
		OriginalPrice:    tp.FormatAmount(price),
		OriginalQuantity: tp.FormatQty(qty),
		OriginalAmount:   "0",
		TradeAvgPrice:    "0",
		TradeQty:         "0",
		TradeAmount:      "0",
		FeeRate:          string(tp.FeeRate),
		Fee:              "0",
		Status:           OrderStatusNew,
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
		_, err = assets.FreezeAssets(db, false, user_id, tp.TargetSymbolId, qty, neworder.OrderId, assets.Behavior_Trade)
		if err != nil {
			return nil, err
		}
		neworder.FreezeAsset = qty
	} else if neworder.OrderSide == OrderSideBuy {
		//买单的冻结金额加上手续费，这里预估全部成交所需要的手续费，
		amount := d(price).Mul(d(qty))
		fee := amount.Mul(d(neworder.FeeRate))
		freeze_amount := amount.Add(fee).String()

		//fee、tradeamount字段在结算程序中修改

		neworder.FreezeAsset = freeze_amount
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
