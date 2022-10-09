package quotation

import (
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/bookvoo/base/symbols"
	"github.com/yzimhao/bookvoo/common"
	"github.com/yzimhao/bookvoo/quotation/models"
)

func GetRouter(r *gin.Engine) *gin.Engine {
	return setupRouter(r)
}

func setupRouter(router *gin.Engine) *gin.Engine {
	apiV1 := router.Group("/api/v1/market")
	{
		apiV1.GET("/ping", ping)
		apiV1.GET("/klines", apiKlines)
		apiV1.GET("/trade_log", nil)
	}
	return router
}

// @Summary 行情k线数据
// @Description 行情k线数据接口
// @Tags 行情接口
// @Accept application/json
// @Produce application/json
// @Param symbol query string true "eg: ethusd"
// @Param period query string true "m1,m3,m5..."
// @Param limit query int false "default: 10"
// @Success 200 {object} common.Response
// @Success 200 {object} models.Kline
// @Router /api/v1/market/klines [get]
func apiKlines(c *gin.Context) {
	symbol := c.Query("symbol")
	period := c.Query("period")

	limit := func() int {
		l := c.GetInt("limit")
		if l <= 0 {
			l = 10
		}
		if l > 100 {
			l = 100
		}
		return l
	}()

	tp, err := symbols.GetExchangeBySymbol(symbol)
	if err != nil {
		common.Fail(c, err.Error())
		return
	}

	intv := models.Period(period)
	kk := models.NewKline(symbol, intv)

	kks := []models.Kline{}
	models.DbEngine().Table(kk.TableName()).OrderBy("open_at desc").Limit(limit).Find(&kks)

	data := []interface{}{}
	for _, item := range kks {

		a := []interface{}{
			item.OpenAt.Unix(),
			tp.FormatAmount(item.Open),
			tp.FormatAmount(item.High),
			tp.FormatAmount(item.Low),
			tp.FormatAmount(item.Close),
			tp.FormatQty(item.Volume),
			tp.FormatAmount(item.Amount),
			item.CloseAt.Unix(),
		}
		data = append(data, a)
	}

	c.JSON(200, data)
}

func ping(c *gin.Context) {
	c.String(200, "pong")
}
