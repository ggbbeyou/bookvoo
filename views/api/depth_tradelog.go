package api

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yzimhao/bookvoo/core/base"
)

// @Summary 深度信息
// @Tags 交易相关
// @Produce application/json
// @Param symbol query string true "eg: ethusd"
// @Param limit  query int false "默认100，最大5000"
// @Success 200 {object} _response
// @Router /api/v1/depth [get]
func depth(c *gin.Context) {
	symbol := strings.ToLower(c.Query("symbol"))
	limit := c.Query("limit")
	limitInt, _ := strconv.Atoi(limit)
	if limitInt <= 0 || limitInt > 5000 {
		limitInt = 100
	}

	if _, ok := base.Engine[symbol]; !ok {
		fail(c, "invalid symbol")
		return
	}

	a := base.Engine[symbol].GetAskDepth(limitInt)
	b := base.Engine[symbol].GetBidDepth(limitInt)
	success(c, gin.H{
		"ask": a,
		"bid": b,
	})
}

func tradelog(c *gin.Context) {
	symbol := strings.ToLower(c.Param("symbol"))

	c.JSON(200, gin.H{
		"ok": true,
		"data": gin.H{
			"latest_price": base.Engine[symbol].Price2String(base.Engine[symbol].LatestPrice()),
			"trade_log":    "", //recentTrade,
		},
	})

}
