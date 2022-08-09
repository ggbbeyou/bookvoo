package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yzimhao/bookvoo/base/symbols"
	"github.com/yzimhao/bookvoo/common"
)

func symbol_info(c *gin.Context) {
	symbol := strings.ToLower(c.Query("symbol"))
	tp, err := symbols.GetExchangeBySymbol(symbol)
	if err != nil {
		c.HTML(http.StatusNotFound, "", nil)
		return
	}
	common.Success(c, tp)
}
