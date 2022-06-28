package market

import (
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/kline/core"
)

func Run(config string, router *gin.Engine) {
	core.RunWithGinRouter(config, router)
}
