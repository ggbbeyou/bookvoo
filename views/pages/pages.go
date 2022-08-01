package pages

import (
	"net/http"

	_ "github.com/yzimhao/bookvoo/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func SetupRouter(router *gin.Engine) {
	router.LoadHTMLGlob("./template/default/*.html")
	router.StaticFS("/statics", http.Dir("./template/default/statics"))

	router.Any("/", func(c *gin.Context) {
		default_symbol := "ethusd"
		c.Redirect(http.StatusMovedPermanently, "/t/"+default_symbol)
	})

	//交易界面
	router.GET("/t/:symbol", func(c *gin.Context) {
		c.HTML(200, "demo.html", gin.H{
			"symbol": c.Param("symbol"),
		})
	})

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
