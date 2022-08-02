package api

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yzimhao/bookvoo/base/symbols"
	"github.com/yzimhao/bookvoo/common"
	"github.com/yzimhao/bookvoo/user/orders"
)

// @Summary 成交记录
// @Tags 交易相关
// @Produce application/json
// @Param symbol query string true "eg: ethusd"
// @Param limit  query int false "默认10，最大100"
// @Success 200 {object} []orders.TradeRecord
// @Router /api/v1/trade/record [get]
func trade_record(c *gin.Context) {
	symbol := strings.ToLower(c.Query("symbol"))
	limit := func() int {
		_limit := c.Query("limit")
		n, _ := strconv.Atoi(_limit)
		if n <= 0 {
			n = 10
		}
		if n > 100 {
			n = 100
		}
		return n
	}()

	tp, err := symbols.GetTradePairBySymbol(symbol)
	if err != nil {
		common.Fail(c, err.Error())
		return
	}

	db := orders.Db().NewSession()
	defer db.Close()

	rows := []orders.TradeRecord{}

	tr := orders.TradeRecord{}
	table := tr.GetTableName(symbol)
	db.Table(table).OrderBy("create_time desc").Limit(limit).Find(&rows)

	for i, row := range rows {
		rows[i].Price = tp.FormatAmount(row.Price)
		rows[i].Quantity = tp.FormatQty(row.Quantity)
		rows[i].Amount = tp.FormatAmount(row.Amount)
		// rows[i].CreateTime = row.CreateTime.Unix()
	}

	common.Success(c, rows)
}
