package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/yzimhao/bookvoo/base/symbols"
	"github.com/yzimhao/bookvoo/common"
	"github.com/yzimhao/bookvoo/match"
	"github.com/yzimhao/bookvoo/user/orders"
)

// @Summary 最新价格
// @Tags 行情相关
// @Accept application/json
// @Produce application/json
// @Param symbol query string true "symbol"
// @Security ApiKeyAuth
// @Success 200 {object} common.Response
// @Router /api/v1/latest/price [get]
func latest_price(c *gin.Context) {
	symbol := c.Query("symbol")

	tp, err := symbols.GetPairBySymbol(symbol)
	if err != nil {
		common.Fail(c, err.Error())
		return
	}

	te, err := match.Engine.Get(symbol)
	if err != nil {
		common.Fail(c, err.Error())
		return
	}

	price := te.LatestPrice()
	if te.LatestPrice().Equal(decimal.Zero) {
		//from db
		db := orders.Db().NewSession()
		defer db.Close()

		row := orders.TradeRecord{Symbol: symbol}
		db.Table(row.TableName()).OrderBy("create_time desc").Limit(1).Get(&row)
		price, _ = decimal.NewFromString(row.Price)
	}

	common.Success(c, gin.H{
		"price": tp.FormatAmount(price.String()),
	})
}
