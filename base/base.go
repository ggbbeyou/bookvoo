package base

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/yzimhao/bookvoo/base/symbols"
	"github.com/yzimhao/gowss"
	"xorm.io/xorm"
)

var (
	Wss *gowss.Hub
)

func Init(db *xorm.Engine, rdc *redis.Client) {
	symbols.Init(db, rdc)
	Wss = gowss.NewHub()
}

func WsHandler(ctx *gin.Context) {
	Wss.ServeWs(ctx.Writer, ctx.Request)
}
