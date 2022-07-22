package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/bookvoo/core"
	"github.com/yzimhao/bookvoo/core/base"
	"github.com/yzimhao/bookvoo/user/orders"
	te "github.com/yzimhao/trading_engine"
)

type new_order_request struct {
	Symbol    string           `json:"symbol" binding:"required"`
	Side      orders.OrderSide `json:"side" binding:"required"`
	OrderType orders.OrderType `json:"order_type" binding:"required"`

	Price    string `json:"price"`
	Quantity string `json:"quantity"`
	Amount   string `json:"amount"`
}

// order_new 创建一个新订单
// @Summary 创建一个新订单
// @Description 新订单，支持限价单、市价单
// @Tags 订单相关
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body new_order_request false "请求参数"
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
		market_order(c, req)
		return
	}

}

func limit_order(c *gin.Context, req new_order_request) {
	order, err := orders.NewLimitOrder(1, req.Symbol, req.Side, req.Price, req.Quantity)
	if err != nil {
		fail(c, err.Error())
		return
	}
	if req.Side == orders.OrderSideAsk {
		base.MatchingEngine[req.Symbol].ChNewOrder <- te.NewAskLimitItem(order.OrderId, core.D(order.Price), core.D(order.Quantity), order.CreateTime)
	} else if req.Side == orders.OrderSideBid {
		base.MatchingEngine[req.Symbol].ChNewOrder <- te.NewBidLimitItem(order.OrderId, core.D(order.Price), core.D(order.Quantity), order.CreateTime)
	}
	success(c, gin.H{"order_id": order.OrderId})
}

func market_order(c *gin.Context, req new_order_request) {

}

func order_cancel(c *gin.Context) {

}
