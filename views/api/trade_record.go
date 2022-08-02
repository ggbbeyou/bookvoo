package api

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yzimhao/bookvoo/common"
)

// @Summary 成交记录
// @Tags 交易相关
// @Produce application/json
// @Param symbol query string true "eg: ethusd"
// @Param limit  query int false "默认10，最大100"
// @Success 200 {object} []orders.TradeRecord
// @Router /api/v1/trade/record [get]
func trade_record(c *gin.Context) {
	symbol := strings.ToLower(c.Param("symbol"))

	common.Success(c, symbol)
}
