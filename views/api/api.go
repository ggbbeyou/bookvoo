package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/yzimhao/bookvoo/user"
)

var (
	USERID int64 = 101
)

func SetupRouter(router *gin.Engine) {
	user.InitJwt()

	apiV1 := router.Group("/api/v1")
	apiV1.GET("/user/login", user.AuthMiddleware.LoginHandler)
	apiV1.GET("/user/logout", user.AuthMiddleware.LogoutHandler)
	apiV1.GET("/user/refresh", user.AuthMiddleware.RefreshHandler)

	//交易对信息查询
	apiV1.GET("/exchange/info", exchange_info)
	//深度
	apiV1.GET("/depth", depth)
	//成交记录
	apiV1.GET("/trade/record", trade_record)

	//需要验证登录的接口
	apiV1.Use(user.AuthMiddleware.MiddlewareFunc())
	{
		//用户信息查询
		apiV1.GET("/user/query", nil)
		//用户资产查询
		apiV1.GET("/assets/query", assets_query)

		order := apiV1.Group("/order")
		{
			//查询订单
			order.GET("/", order_query)
			//创建订单
			order.POST("/new", order_new)
			//取消订单
			order.POST("/cancel", order_cancel)
			//当前挂单
			order.GET("/open", order_open)
			//查询所有订单 获取所有帐户订单； 有效，已取消或已完成。 带有symbol
			order.GET("/all", order_all)
		}
	}
}

func getUserId(c *gin.Context) int64 {
	uinfo, _ := c.Get("user")
	return uinfo.(*user.User).UserId
}

func d(ss string) decimal.Decimal {
	s, _ := decimal.NewFromString(ss)
	return s
}
