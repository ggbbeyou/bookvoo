package pages

import (
	"net/http"
	"strings"

	"github.com/yzimhao/bookvoo/base/symbols"
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
		symbol := strings.ToLower(c.Param("symbol"))
		tp, err := symbols.GetExchangeBySymbol(symbol)
		if err != nil {
			c.HTML(http.StatusNotFound, "", nil)
			return
		}

		c.HTML(200, "demo.html", gin.H{
			"symbol": symbol,
			"tp":     tp,
		})
	})

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
