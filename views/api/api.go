// @title 这里写标题
// @version 1.0
// @description 这里写描述信息
// @termsOfService http://swagger.io/terms/
// @contact.name 这里写联系人信息
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host 这里写接口服务的host
// @BasePath 这里写base path
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
