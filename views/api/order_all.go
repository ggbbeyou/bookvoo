package api

import "github.com/gin-gonic/gin"

// @Summary 查询所有订单
// @Description 查询历史所有订单
// @Tags 订单相关
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"

// @Param symbol query string true "eg: ethusd"

// @Security ApiKeyAuth
// @Success 200 {object} common.Response
// @Router /api/v1/order/all [get]
func order_all(c *gin.Context) {

}
