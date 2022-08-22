package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/bookvoo/base/symbols"
	"github.com/yzimhao/bookvoo/common"
	"github.com/yzimhao/bookvoo/user/orders"
)

// @Summary 查询当前挂单
// @Description 查询当前还未完全成交的挂单
// @Tags 订单相关
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param symbol query string true "eg: ethusd"
// @Security ApiKeyAuth
// @Success 200 {object} common.Response
// @Router /api/v1/order/open [get]
func order_open(c *gin.Context) {
	symbol := c.Query("symbol")
	db := orders.Db().NewSession()
	defer db.Close()

	es, err := symbols.GetExchangeBySymbol(symbol)
	if err != nil {
		common.Fail(c, err.Error())
		return
	}

	rows := []orders.TradeOrder{}
	db.Table(new(orders.UnfinishedOrder)).Where("user_id=? and pair_id=?", getUserId(c), es.Id).Find(&rows)
	for i, item := range rows {
		rows[i].Amount = es.FormatAmount(item.Amount)
		rows[i].Price = es.FormatAmount(item.Price)
		rows[i].AvgPrice = es.FormatAmount(item.AvgPrice)
		rows[i].FinishedQty = es.FormatQty(item.FinishedQty)
		rows[i].Quantity = es.FormatQty(item.Quantity)
		rows[i].TradeAmount = es.FormatAmount(item.TradeAmount)
	}
	common.Success(c, rows)
}
