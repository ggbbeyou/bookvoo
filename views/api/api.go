package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

var (
	USERID int64 = 101
)

func SetupRouter(router *gin.Engine) {
	apiV1 := router.Group("/api/v1")

	apiV1.GET("/depth", depth)
	apiV1.GET("/trade/log", tradelog)

	//todo 验证登录状态
	order := apiV1.Group("/order")
	{
		order.Use(func(ctx *gin.Context) {
			//todo 登陆中间件
			ctx.Set("user_id", USERID)
		})
		//查询订单
		order.GET("/", nil)
		//创建订单
		order.POST("/new", order_new)
		//取消订单
		order.POST("/cancel", order_cancel)
		//当前挂单
		order.GET("/open", nil)
		//查询所有订单 获取所有帐户订单； 有效，已取消或已完成。 带有symbol
		order.GET("/all", nil)
	}
}

func getUserId(c *gin.Context) int64 {
	val, _ := c.Get("user_id")
	switch val.(type) {
	case int64:
		return val.(int64)
	}
	return -1
}

func d(ss string) decimal.Decimal {
	s, _ := decimal.NewFromString(ss)
	return s
}
