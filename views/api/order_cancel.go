package api

import (
	"github.com/gin-gonic/gin"
)

type cancel_order_request struct {
	OrderId string `json:"order_id"`
}

// @Summary 取消一个委托订单
// @Description 取消还未完成的订单
// @Tags 交易相关
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body cancel_order_request true "请求参数"
// @Security ApiKeyAuth
// @Success 200 {object} common.Response
// @Router /api/v1/order/cancel [post]
func order_cancel(c *gin.Context) {
	//todo
}
