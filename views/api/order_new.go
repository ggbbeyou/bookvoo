package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/yzimhao/bookvoo/core"
	"github.com/yzimhao/bookvoo/core/base"
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

// 新委托订单
// @Summary 创建一个新委托订单
// @Description 新订单，支持限价单、市价单
// @Description 不同订单类型的参数要求：
// @Description 限价单: {"symbol": "ethusd", "order_type": "limit", "side": "sell", "price": "1.00", "quantity": "100"}
// @Description 市价-按数量: {"symbol": "ethusd", "order_type": "market", "side": "sell", "quantity": "100"}
// @Description 市价-按金额: {"symbol": "ethusd", "order_type": "market", "side": "sell", "amount": "1000.00"}
// @Tags 交易相关
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body new_order_request true "请求参数"
// @Security ApiKeyAuth
// @Success 200 {object} _response
// @Router /api/v1/order/new [post]
func order_new(c *gin.Context) {

	var req new_order_request
	if err := c.BindJSON(&req); err != nil {
		fail(c, err.Error())
		return
	}

	if req.OrderType == orders.OrderTypeLimit {
		limit_order(c, req)
		return
	} else if req.OrderType == orders.OrderTypeMarket {
		//todo
		if core.D(req.Amount).Cmp(decimal.Zero) > 0 {
			//按金额操作
			market_order_by_amount(c, req.Symbol, req.Side, req.Amount)
		} else if core.D(req.Quantity).Cmp(decimal.Zero) > 0 {
			//按数量操作
			market_order_by_qty(c, req.Symbol, req.Side, req.Quantity)
		}
		return
	}

}

func limit_order(c *gin.Context, req new_order_request) {
	order, err := orders.NewLimitOrder(1, req.Symbol, req.Side, req.Price, req.Quantity)
	if err != nil {
		fail(c, err.Error())
		return
	}
	if req.Side == orders.OrderSideSell {
		base.MatchingEngine[req.Symbol].ChNewOrder <- te.NewAskLimitItem(order.OrderId, core.D(order.Price), core.D(order.Quantity), order.CreateTime)
	} else if req.Side == orders.OrderSideBuy {
		base.MatchingEngine[req.Symbol].ChNewOrder <- te.NewBidLimitItem(order.OrderId, core.D(order.Price), core.D(order.Quantity), order.CreateTime)
	}
	success(c, gin.H{"order_id": order.OrderId})
}

func market_order_by_amount(c *gin.Context, symbol string, side orders.OrderSide, amount string) {

}

func market_order_by_qty(c *gin.Context, symbol string, side orders.OrderSide, qty string) {

}
