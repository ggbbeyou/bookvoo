package api

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yzimhao/bookvoo/base/symbols"
	"github.com/yzimhao/bookvoo/common"
	"github.com/yzimhao/bookvoo/user/assets"
)

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
