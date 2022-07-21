package pages

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine) {
	router.LoadHTMLGlob("./template/default/*.html")
	router.StaticFS("/statics", http.Dir("./template/default/statics"))

	//交易界面
	router.GET("/t/:symbol", func(c *gin.Context) {
		c.HTML(200, "demo.html", gin.H{
			"symbol": c.Param("symbol"),
		})
	})
}
