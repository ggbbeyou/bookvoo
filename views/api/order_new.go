package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/yzimhao/bookvoo/base"
	"github.com/yzimhao/bookvoo/common"
	"github.com/yzimhao/bookvoo/match"
	"github.com/yzimhao/bookvoo/user/orders"
	gowss "github.com/yzimhao/bookvoo/wss"
	te "github.com/yzimhao/trading_engine"
)

type new_order_request struct {
	Symbol    string           `json:"symbol" binding:"required" example:"ethusd"`
	Side      orders.OrderSide `json:"side" binding:"required" example:"sell/buy"`
	OrderType orders.OrderType `json:"order_type" binding:"required" example:"limit/market"`
	Price     string           `json:"price" example:"1.00"`
	Quantity  string           `json:"quantity" example:"12"`
	Amount    string           `json:"amount" example:"100.00"`
}

// @Summary 委托订单
// @Tags 订单相关
// @Description 新订单，支持限价单、市价单
// @Description 不同订单类型的参数要求：
// @Description 限价单: {"symbol": "ethusd", "order_type": "limit", "side": "sell/buy", "price": "1.00", "quantity": "100"}
// @Description 市价-按数量: {"symbol": "ethusd", "order_type": "market", "side": "sell/buy", "quantity": "100"}
// @Description 市价-按金额: {"symbol": "ethusd", "order_type": "market", "side": "sell/buy", "amount": "1000.00"}
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param object body new_order_request true "请求参数"
// @Security ApiKeyAuth
// @Success 200 {object} common.Response
// @Router /api/v1/order/new [post]
func order_new(c *gin.Context) {
	var req new_order_request
	if err := c.BindJSON(&req); err != nil {
		common.Fail(c, err.Error())
		return
	}

	var newOrder *orders.TradeOrder
	var err error
	if req.OrderType == orders.OrderTypeLimit {
		newOrder, err = limit_order(c, req)
		match.Send <- *newOrder
	} else if req.OrderType == orders.OrderTypeMarket {
		if d(req.Amount).Cmp(decimal.Zero) > 0 {
			//按金额操作
			newOrder, err = market_order_by_amount(c, req.Symbol, req.Side, req.Amount)
		} else if d(req.Quantity).Cmp(decimal.Zero) > 0 {
			//按数量操作
			newOrder, err = market_order_by_qty(c, req.Symbol, req.Side, req.Quantity)
		}
	} else {
		common.Fail(c, "invalid side")
		return
	}

	if err != nil {
		common.Fail(c, err.Error())
		return
	}

	base.WssPush(gowss.MsgBody{
		To: fmt.Sprintf("%d", newOrder.UserId),
		Response: gowss.Response{
			Type: "new_order",
			Body: newOrder,
		},
	})

	common.Success(c, gin.H{"order_id": newOrder.OrderId})
}

func limit_order(c *gin.Context, req new_order_request) (*orders.TradeOrder, error) {
	uid := getUserId(c)
	order, err := orders.NewLimitOrder(uid, req.Symbol, req.Side, req.Price, req.Quantity)
	if err != nil {
		return nil, err
	}
	return order, nil
}

//市价按数量操作
func market_order_by_qty(c *gin.Context, symbol string, side orders.OrderSide, qty string) (*orders.TradeOrder, error) {
	uid := getUserId(c)
	order, err := orders.NewMarketOrderByQty(uid, symbol, side, qty)
	if err != nil {
		return nil, err
	}

	t, _ := match.Engine.Get(symbol)
	if side == orders.OrderSideSell {
		t.ChNewOrder <- te.NewAskMarketQtyItem(order.OrderId, d(order.Quantity), order.CreateTime)
	} else if side == orders.OrderSideBuy {
		t.ChNewOrder <- te.NewBidMarketQtyItem(order.OrderId, d(order.Quantity), d(order.FreezeQty), order.CreateTime)
	}

	return order, nil
}

//市价按成交量操作
func market_order_by_amount(c *gin.Context, symbol string, side orders.OrderSide, amount string) (*orders.TradeOrder, error) {
	uid := getUserId(c)
	order, err := orders.NewMarketOrderByAmount(uid, symbol, side, amount)
	if err != nil {
		return nil, err
	}

	t, _ := match.Engine.Get(symbol)
	if side == orders.OrderSideSell {
		t.ChNewOrder <- te.NewAskMarketAmountItem(order.OrderId, d(amount), d(order.FreezeQty), order.CreateTime)
	} else if side == orders.OrderSideBuy {
		t.ChNewOrder <- te.NewBidMarketAmountItem(order.OrderId, d(order.Amount), order.CreateTime)
	}
	return order, nil
}
