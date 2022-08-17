package base

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	gowss "github.com/yzimhao/bookvoo/wss"
)

func WsHandler(ctx *gin.Context) {
	Wss.ServeWs(ctx.Writer, ctx.Request)
}

func WssPush(rdc *redis.Client, msg gowss.MsgBody) {
	//除了广播到前端外，还需要推送一份到k线计算
	if Wss != nil {
		Wss.Broadcast <- msg
	}
}
