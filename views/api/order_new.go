package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/yzimhao/bookvoo/common"
	"github.com/yzimhao/bookvoo/match"
	"github.com/yzimhao/bookvoo/user/orders"
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

	//测试代码，买卖不同账户 避免结算时出问题
	if req.Side == orders.OrderSideBuy {
		c.Set("user_id", int64(101))
	} else {
		c.Set("user_id", int64(102))
	}

	if req.OrderType == orders.OrderTypeLimit {
		limit_order(c, req)
		return
	} else if req.OrderType == orders.OrderTypeMarket {
		//todo
		if d(req.Amount).Cmp(decimal.Zero) > 0 {
			//按金额操作
			market_order_by_amount(c, req.Symbol, req.Side, req.Amount)
		} else if d(req.Quantity).Cmp(decimal.Zero) > 0 {
			//按数量操作
			market_order_by_qty(c, req.Symbol, req.Side, req.Quantity)
		}
		return
	} else {
		common.Fail(c, "invalid side")
		return
	}
}

func limit_order(c *gin.Context, req new_order_request) {
	uid := getUserId(c)
	order, err := orders.NewLimitOrder(uid, req.Symbol, req.Side, req.Price, req.Quantity)
	if err != nil {
		common.Fail(c, err.Error())
		return
	}

	match.Send <- order
	// if req.Side == orders.OrderSideSell {
	// 	match.Engine[req.Symbol].ChNewOrder <- te.NewAskLimitItem(order.OrderId, d(order.Price), d(order.Quantity), order.CreateTime)
	// } else if req.Side == orders.OrderSideBuy {
	// 	match.Engine[req.Symbol].ChNewOrder <- te.NewBidLimitItem(order.OrderId, d(order.Price), d(order.Quantity), order.CreateTime)
	// }
	common.Success(c, gin.H{"order_id": order.OrderId})
}

//市价按数量操作
func market_order_by_qty(c *gin.Context, symbol string, side orders.OrderSide, qty string) {
	uid := getUserId(c)
	order, err := orders.NewMarketOrderByQty(uid, symbol, side, qty)
	if err != nil {
		common.Fail(c, err.Error())
		return
	}

	if side == orders.OrderSideSell {
		match.Engine.Symbols[symbol].ChNewOrder <- te.NewAskMarketQtyItem(order.OrderId, d(order.Quantity), order.CreateTime)
	} else if side == orders.OrderSideBuy {
		match.Engine.Symbols[symbol].ChNewOrder <- te.NewBidMarketQtyItem(order.OrderId, d(order.Quantity), d(order.FreezeQty), order.CreateTime)
	}

	common.Success(c, gin.H{"order_id": order.OrderId})
}

//市价按成交量操作
func market_order_by_amount(c *gin.Context, symbol string, side orders.OrderSide, amount string) {
	uid := getUserId(c)
	order, err := orders.NewMarketOrderByAmount(uid, symbol, side, amount)
	if err != nil {
		common.Fail(c, err.Error())
		return
	}

	if side == orders.OrderSideSell {

		match.Engine.Symbols[symbol].ChNewOrder <- te.NewAskMarketAmountItem(order.OrderId, d(amount), d(order.FreezeQty), order.CreateTime)
	} else if side == orders.OrderSideBuy {
		match.Engine.Symbols[symbol].ChNewOrder <- te.NewBidMarketAmountItem(order.OrderId, d(order.Amount), order.CreateTime)
	}

	common.Success(c, gin.H{"order_id": order.OrderId})
}
