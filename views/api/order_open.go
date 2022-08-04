package api

import "github.com/gin-gonic/gin"

// @Summary 查询当前挂单
// @Description 查询当前还未完全成交的挂单
// @Tags 订单相关
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"

// @Param symbol query string true "eg: ethusd"

// @Security ApiKeyAuth
// @Success 200 {object} common.Response
// @Router /api/v1/order/open [get]
func order_open(c *gin.Context) {

}
