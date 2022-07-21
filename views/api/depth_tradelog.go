package api

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yzimhao/bookvoo/core/base"
)

func depth(c *gin.Context) {
	symbol := strings.ToLower(c.Param("symbol"))
	limit := c.Query("limit")
	limitInt, _ := strconv.Atoi(limit)
	if limitInt <= 0 || limitInt > 100 {
		limitInt = 10
	}

	//todo 验证是否存在
	a := base.MatchingEngine[symbol].GetAskDepth(limitInt)
	b := base.MatchingEngine[symbol].GetBidDepth(limitInt)

	c.JSON(200, gin.H{
		"ask": a,
		"bid": b,
	})
}

func tradelog(c *gin.Context) {
	symbol := strings.ToLower(c.Param("symbol"))

	c.JSON(200, gin.H{
		"ok": true,
		"data": gin.H{
			"latest_price": base.MatchingEngine[symbol].Price2String(base.MatchingEngine[symbol].LatestPrice()),
			"trade_log":    "", //recentTrade,
		},
	})

}
