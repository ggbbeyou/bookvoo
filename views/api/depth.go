package api

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yzimhao/bookvoo/common"
	"github.com/yzimhao/bookvoo/match"
)

// @Summary 深度信息
// @Produce application/json
// @Param symbol query string true "eg: ethusd"
// @Param limit  query int false "默认100，最大5000"
// @Success 200 {object} common.Response
// @Router /api/v1/depth [get]
func depth(c *gin.Context) {
	symbol := strings.ToLower(c.Query("symbol"))
	limit := c.Query("limit")
	limitInt, _ := strconv.Atoi(limit)
	if limitInt <= 0 || limitInt > 5000 {
		limitInt = 100
	}

	te, err := match.Engine.Get(symbol)
	if err != nil {
		common.Fail(c, err.Error())
		return
	}

	a := te.GetAskDepth(limitInt)
	b := te.GetBidDepth(limitInt)
	common.Success(c, gin.H{
		"ask": a,
		"bid": b,
	})
}
