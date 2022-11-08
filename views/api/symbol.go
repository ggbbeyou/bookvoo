package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yzimhao/bookvoo/base/symbols"
	"github.com/yzimhao/bookvoo/common"
)

// @Summary 交易对信息
// @Tags
// @Description 获取交易规则、交易对信息
// @Accept application/json
// @Produce application/json
// @Param symbol query string true "交易对symbol"
// @Success 200 {object} common.Response
// @Router /api/v1/exchange/info [get]
func exchange_info(c *gin.Context) {
	symbol := strings.ToLower(c.Query("symbol"))
	tp, err := symbols.GetPairBySymbol(symbol)
	if err != nil {
		c.HTML(http.StatusNotFound, "", nil)
		return
	}
	common.Success(c, tp)
}
