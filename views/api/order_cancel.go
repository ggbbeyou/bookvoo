package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/bookvoo/base/symbols"
	"github.com/yzimhao/bookvoo/common"
	"github.com/yzimhao/bookvoo/match"
	"github.com/yzimhao/bookvoo/user/orders"
	"github.com/yzimhao/trading_engine"
)

type cancel_order_request struct {
	Symbol  string `json:"symbol"`
	OrderId string `json:"order_id"`
}

// @Summary 取消委托
// @Description 取消还未完成的订单
// @Tags 订单相关
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param object body cancel_order_request true "请求参数"
// @Security ApiKeyAuth
// @Success 200 {object} common.Response
// @Router /api/v1/order/cancel [post]
func order_cancel(c *gin.Context) {
	var req cancel_order_request
	if err := c.BindJSON(&req); err != nil {
		common.Fail(c, err.Error())
		return
	}

	tp, err := symbols.GetPairBySymbol(req.Symbol)
	if err != nil {
		common.Fail(c, err.Error())
		return
	}

	side := orders.OrderIDSide(req.OrderId)
	tside := func() trading_engine.OrderSide {
		if side == orders.OrderSideBuy {
			return trading_engine.OrderSideBuy
		}
		return trading_engine.OrderSideSell
	}()

	t, _ := match.Engine.Get(tp.Symbol)
	t.CancelOrder(tside, req.OrderId)
	common.Success(c, nil)
}
