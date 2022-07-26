package clearing

import (
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/yzimhao/bookvoo/user/orders"
	"xorm.io/xorm"
)

type clearing struct {
	db           *xorm.Session
	symbol       string
	ask_order_id string
	bid_order_id string
	trade_price  decimal.Decimal
	trade_qty    decimal.Decimal
	trade_amount decimal.Decimal
	record       *orders.TradeRecord
	ask          *orders.TradeOrder
	bid          *orders.TradeOrder
}

func (c *clearing) check() error {
	//
	_, err := c.db.Table(new(orders.TradeOrder)).Where("order_id=?", c.ask_order_id).ForUpdate().Get(&c.ask)
	if err != nil {
		return err
	}

	_, err = c.db.Table(new(orders.TradeOrder)).Where("order_id=?", c.bid_order_id).ForUpdate().Get(&c.bid)
	if err != nil {
		return err
	}

	if c.ask.Status != orders.OrderStatusNew {
		return fmt.Errorf("ask status error")
	}

	if c.bid.Status != orders.OrderStatusNew {
		return fmt.Errorf("bid status error")
	}
	return nil
}

func (c *clearing) updateAsk() error {
	return c.updateOrder(orders.OrderSideSell)
}

func (c *clearing) updateBid() error {
	return c.updateOrder(orders.OrderSideBuy)
}

func (c *clearing) updateOrder(side orders.OrderSide) error {
	var order orders.TradeOrder
	if side == orders.OrderSideSell {
		order = *c.ask
		order.Fee = d(order.Fee).Add(d(c.record.AskFee)).String()
	} else {
		order = *c.bid
		order.Fee = d(order.Fee).Add(d(c.record.BidFee)).String()
	}
	order.FinishedQty = d(order.FinishedQty).Add(c.trade_qty).String()
	order.TradeAmount = d(order.TradeAmount).Add(c.trade_amount).String()
	order.TradeSymbol = c.symbol

	//todo 一些必要的边界值检查

	// if d(c.ask.FinishedQty).Cmp(d(c.ask.Quantity)) <= 0 {
	// }
	_, err := c.db.Table(order.TableName()).Where("order_id=?", order.OrderId).Update(order)
	if err != nil {
		return err
	}

	if order.OrderType == orders.OrderTypeLimit {
		if order.Status == orders.OrderStatusNew {
			_, err := c.db.Table(new(orders.UnfinishedOrder)).Where("order_id=?", order.OrderId).Update(order)
			if err != nil {
				return err
			}
		} else {
			_, err := c.db.Table(new(orders.UnfinishedOrder)).Where("order_id=?", order.OrderId).Delete()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *clearing) tradeRecord() error {

	trade := orders.TradeRecord{
		Symbol: c.symbol,
		Ask:    c.ask_order_id,
		Bid:    c.bid_order_id,
		TradeBy: func() orders.TradeBy {
			if c.ask.CreateTime > c.bid.CreateTime {
				return orders.TradeBySell
			} else {
				return orders.TradeByBuy
			}
		}(),

		AskUid:   c.ask.UserId,
		Biduid:   c.bid.UserId,
		Price:    c.trade_price.String(),
		Quantity: c.trade_qty.String(),

		AskFeeRate: c.ask.FeeRate,
		AskFee:     c.trade_amount.Mul(d(c.ask.FeeRate)).String(),

		BidFeeRate: c.bid.FeeRate,
		BidFee:     c.trade_amount.Mul(d(c.bid.FeeRate)).String(),
	}

	if err := trade.Save(c.db); err != nil {
		return err
	}
	c.record = &trade
	return nil
}

func d(s string) decimal.Decimal {
	dd, _ := decimal.NewFromString(s)
	return dd
}
