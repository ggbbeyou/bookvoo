package tradecore

import (
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine"
)

var DemoPair *trading_engine.TradePair

func Run(config string, router *gin.Engine) {
	DemoPair = trading_engine.NewTradePair("demo", 2, 4)
}
