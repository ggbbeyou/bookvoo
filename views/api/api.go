package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type _response struct {
	Ok     int         `json:"ok"`
	Reason string      `json:"reason,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

func SetupRouter(router *gin.Engine) {
	apiV1 := router.Group("/api/v1/")

	apiV1.GET("/depth", depth)
	apiV1.GET("/trade/log", tradelog)

	//todo 验证登录状态
	apiV1.POST("/order/new", order_new)
	apiV1.POST("/order/cancel", order_cancel)
}

func response(c *gin.Context, ok int, reason string, data interface{}) {
	res := _response{
		Ok:     ok,
		Reason: reason,
		Data:   data,
	}
	c.JSON(http.StatusOK, res)
}

func success(c *gin.Context, data interface{}) {
	response(c, 1, "", data)
}

func fail(c *gin.Context, reason string) {
	response(c, 0, reason, nil)
}
