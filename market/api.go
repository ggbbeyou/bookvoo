package market

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/yzimhao/bookvoo/market/models"
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
		apiV1.GET("/latest/price", nil)
	}
	return router
}

func apiKlines(c *gin.Context) {
	symbol := c.Query("symbol")
	period := c.Query("period")

	// startTime := c.Query("start_time")
	// endTime := c.Query("end_time")
	// limit := c.GetInt("limit") //default: 500, max: 1000

	intv := models.Period(period)
	kk := models.NewKline(symbol, intv)

	kks := []models.Kline{}
	models.DbEngine().Table(kk.TableName()).OrderBy("open_at desc").Limit(500).Find(&kks)

	data := []interface{}{}
	for _, item := range kks {
		//todo 小数点位数处理
		a := []interface{}{
			item.OpenAt.Unix(),
			fmtStringDigit(item.Open, 2),
			fmtStringDigit(item.High, 2),
			fmtStringDigit(item.Low, 2),
			fmtStringDigit(item.Close, 2),
			fmtStringDigit(item.Volume, 0),
			fmtStringDigit(item.Amount, 2),
			item.CloseAt.Unix(),
		}
		data = append(data, a)
	}

	c.JSON(200, data)
}

func ping(c *gin.Context) {
	c.String(200, "pong")
}

func fmtStringDigit(a string, digit int) string {
	d, _ := decimal.NewFromString(a)
	return d.StringFixedBank(int32(digit))
}
