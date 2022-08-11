package api

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yzimhao/bookvoo/base/symbols"
	"github.com/yzimhao/bookvoo/common"
	"github.com/yzimhao/bookvoo/user/assets"
)

// @Summary 用户资产查询
// @Tags
// @Description 获取交易规则、交易对信息
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param symbols query string true "资产symbols用逗号分隔 eg: eth,usd"
// @Success 200 {object} common.Response
// @Router /api/v1/assets/query [get]
func assets_query(c *gin.Context) {
	ss := strings.ToLower(c.Query("symbols"))
	slice := strings.Split(ss, ",")

	data := make(map[string]assets.Assets)

	for _, symbol := range slice {
		info, err := symbols.GetSymbolInfoBySymbol(symbol)
		if err != nil {
			continue
		}
		row := assets.UserAssets(getUserId(c), info.Id)

		row.Available = info.FormatNumber(row.Available)
		row.Total = info.FormatNumber(row.Total)
		row.Freeze = info.FormatNumber(row.Freeze)

		data[symbol] = row
	}
	common.Success(c, data)
}
