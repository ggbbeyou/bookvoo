package orders

import (
	"fmt"

	"github.com/yzimhao/bookvoo/core"
	"github.com/yzimhao/bookvoo/core/base"
	"github.com/yzimhao/bookvoo/user/assets"
)

func trade_limit_order(user_id int64, trade_pair_id int, side OrderSide, price, qty, fee_rate string) (order_id string, err error) {
	order := TradeOrder{
		OrderId:       order_id_by_side(side),
		TradingPair:   trade_pair_id,
		OrderSide:     side,
		OrderType:     OrderTypeLimit,
		UserId:        user_id,
		Price:         price,
		Quantity:      qty,
		UnfinishedQty: qty,
		FeeRate:       fee_rate,
		Status:        OrderStatusNew,
	}

	//冻结相应资产
	sess := db_engine.NewSession()
	defer sess.Close()

	//todo 开启事务

	tp := base.GetTradePairById(trade_pair_id)
	if tp == nil {
		return "", fmt.Errorf("invalid trade pair id")
	}

	//冻结相应资产
	if order.OrderSide == OrderSideAsk {
		_, err := assets.FreeeBalance(sess, user_id, tp.TradeSymbolId, qty, order.OrderId, "trade order")
		if err != nil {
			return "", err
		}
	} else if order.OrderSide == OrderSideBid {
		//买单的冻结金额
		amount := core.D(price).Mul(core.D(qty))
		fee := amount.Mul(core.D(fee_rate))
		freeze_amount := amount.Add(fee).String()
		_, err := assets.FreeeBalance(sess, user_id, tp.BaseSymbolId, freeze_amount, order.OrderId, "trade order")
		if err != nil {
			return "", err
		}
	}

	//save order
	_, err = sess.Table(new(TradeOrder)).Insert(&order)
	if err != nil {
		return "", err
	}

	_, err = sess.Table(new(UnfinishedOrder)).Insert(&order)
	if err != nil {
		return "", err
	}
	return order.OrderId, nil
}
