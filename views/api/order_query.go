package api

import "github.com/gin-gonic/gin"

// @Summary 查询挂单
// @Description 查询挂单
// @Tags 订单相关
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"

// @Param symbol query string true "eg: ethusd"
// @Param order_id query string true "eg: A22080117295286970700066"

// @Security ApiKeyAuth
// @Success 200 {object} orders.TradeOrder
// @Router /api/v1/order [get]
func order_query(c *gin.Context) {

}
