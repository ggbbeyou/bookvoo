package clearing

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/yzimhao/bookvoo/user/assets"
	"github.com/yzimhao/bookvoo/user/orders"
	"xorm.io/xorm"
)

type clearing struct {
	db                 *xorm.Session
	symbol             string
	symbol_id          int
	standard_symbol_id int
	ask_order_id       string
	bid_order_id       string
	trade_price        decimal.Decimal
	trade_qty          decimal.Decimal
	trade_amount       decimal.Decimal
	record             *orders.TradeRecord
	ask                *orders.TradeOrder
	bid                *orders.TradeOrder
}

func (c *clearing) check() error {
	//

	_, err := c.db.Table(orders.GetOrderTableName(c.symbol)).Where("order_id=?", c.ask_order_id).ForUpdate().Get(c.ask)
	if err != nil {
		return err
	}

	_, err = c.db.Table(orders.GetOrderTableName(c.symbol)).Where("order_id=?", c.bid_order_id).ForUpdate().Get(c.bid)
	if err != nil {
		return err
	}

	if c.ask.Status != orders.OrderStatusNew {
		return fmt.Errorf("%s status error", c.ask_order_id)
	}

	if c.bid.Status != orders.OrderStatusNew {
		return fmt.Errorf("%s status error", c.bid_order_id)
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
	var order *orders.TradeOrder
	if side == orders.OrderSideSell {
		order = c.ask
		order.Fee = d(order.Fee).Add(d(c.record.AskFee)).String()

	} else {
		order = c.bid
		order.Fee = d(order.Fee).Add(d(c.record.BidFee)).String()
	}

	order.Symbol = c.symbol
	order.FinishedQty = d(order.FinishedQty).Add(c.trade_qty).String()
	order.TradeAmount = d(order.TradeAmount).Add(c.trade_amount).String()
	order.AvgPrice = d(order.TradeAmount).Div(d(order.FinishedQty)).String()
	//todo 一些必要的边界值检查

	if order.OrderType == orders.OrderTypeLimit {
		be := d(order.FinishedQty).Cmp(d(order.Quantity))
		if be > 0 {
			return fmt.Errorf("finished quantity must be  <= order.Quantity")
		}
		if be == 0 {
			order.Status = orders.OrderStatusDone
		}

		_, err := c.db.Table(order.TableName()).Where("order_id=?", order.OrderId).AllCols().Update(order)
		if err != nil {
			return err
		}

		if order.Status == orders.OrderStatusNew {
			_, err := c.db.Table(new(orders.UnfinishedOrder)).Where("order_id=?", order.OrderId).AllCols().Update(order)
			if err != nil {
				return err
			}
		} else {
			_, err := c.db.Table(new(orders.UnfinishedOrder)).Where("order_id=?", order.OrderId).Delete()
			if err != nil {
				return err
			}
		}
	} else if order.OrderType == orders.OrderTypeMarket {
		_, err := c.db.Table(order.TableName()).Where("order_id=?", order.OrderId).AllCols().Update(order)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *clearing) tradeRecord() error {

	trade := orders.TradeRecord{
		Symbol:  c.symbol,
		TradeId: trade_id(c.ask_order_id, c.bid_order_id),
		Ask:     c.ask_order_id,
		Bid:     c.bid_order_id,
		TradeBy: func() orders.TradeBy {
			if c.ask.CreateTime > c.bid.CreateTime {
				return orders.TradeBySell
			} else {
				return orders.TradeByBuy
			}
		}(),

		AskUid:   c.ask.UserId,
		BidUid:   c.bid.UserId,
		Price:    c.trade_price.String(),
		Quantity: c.trade_qty.String(),
		Amount:   c.trade_price.Mul(c.trade_qty).String(),

		AskFeeRate: c.ask.FeeRate,
		AskFee:     c.trade_amount.Mul(d(c.ask.FeeRate)).String(),

		BidFeeRate: c.bid.FeeRate,
		BidFee:     c.trade_amount.Mul(d(c.bid.FeeRate)).String(),
	}

	if err := trade.Save(c.db); err != nil {
		logrus.Debugf("%+v, %s", trade, err)
		return err
	}
	c.record = &trade
	return nil
}

func (c *clearing) transfer() error {
	//给买家结算交易物品
	_, err := assets.UnfreezeAssets(c.db, false, c.ask.UserId, c.ask_order_id, c.trade_qty.String())
	if err != nil {
		return err
	}
	_, err = assets.Transfer(c.db, false, c.ask.UserId, c.bid.UserId, c.symbol_id, c.trade_qty.String(), c.record.TradeId, assets.Behavior_Trade)
	if err != nil {
		return err
	}

	//卖家结算本位币
	amount := d(c.record.Amount).Add(d(c.record.BidFee))
	_, err = assets.UnfreezeAssets(c.db, false, c.bid.UserId, c.bid_order_id, amount.String())
	if err != nil {
		return err
	}
	//卖家收到的本位币需要扣除fee
	fee := d(c.record.BidFee).Add(d(c.record.AskFee))
	_, err = assets.Transfer(c.db, false, c.bid.UserId, c.ask.UserId, c.standard_symbol_id, amount.Sub(fee).String(), c.record.TradeId, assets.Behavior_Trade)
	if err != nil {
		return err
	}

	//手续费收入到一个全局的账号里
	_, err = assets.Transfer(c.db, false, c.bid.UserId, assets.SystemFeeUserID, c.standard_symbol_id, fee.String(), c.record.TradeId, assets.Behavior_Trade)
	if err != nil {
		return err
	}
	return nil
}

func d(s string) decimal.Decimal {
	dd, _ := decimal.NewFromString(s)
	return dd
}

func trade_id(ask_id, bid_id string) string {
	times := time.Now().Format("060102")
	hash := hash256(fmt.Sprintf("%s%s", ask_id, bid_id))
	return fmt.Sprintf("T%s%s", times, hash[0:17])
}

func hash256(data interface{}) string {
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%v", data)))
	hashed := fmt.Sprintf("%x", hash.Sum(nil))
	return hashed
}
