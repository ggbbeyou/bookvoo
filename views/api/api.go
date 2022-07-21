package api

import "github.com/gin-gonic/gin"

func SetupRouter(router *gin.Engine) {
	apiV1 := router.Group("/api/v1/")

	apiV1.GET("/depth/:symbol", depth)
	apiV1.GET("/tradelog/:symbol", tradelog)

	//todo 验证登录状态
	apiV1.POST("/order/new", order_new)
	apiV1.POST("/order/cancel", order_cancel)
}
