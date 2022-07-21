// Package classification User API.
//
// The purpose of this service is to provide an application
// that is using plain go code to define an API
//
//      Host: localhost
//      Version: 0.0.1
//
// swagger:meta

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type _response struct {
	Ok     bool        `json:"ok"`
	Reason string      `json:"reason"`
	Data   interface{} `json:"data"`
}

func SetupRouter(router *gin.Engine) {
	apiV1 := router.Group("/api/v1/")

	apiV1.GET("/depth/:symbol", depth)
	apiV1.GET("/tradelog/:symbol", tradelog)

	//todo 验证登录状态
	apiV1.POST("/order/new", order_new)
	apiV1.POST("/order/cancel", order_cancel)
}

func response(c *gin.Context, ok bool, reason string, data interface{}) {
	res := _response{
		Ok:     ok,
		Reason: reason,
		Data:   data,
	}
	c.JSON(http.StatusOK, res)
}

func success(c *gin.Context, data interface{}) {
	response(c, true, "", data)
}

func fail(c *gin.Context, reason string) {
	response(c, false, reason, nil)
}
