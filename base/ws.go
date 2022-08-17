package base

import (
	"github.com/gin-gonic/gin"
	gowss "github.com/yzimhao/bookvoo/wss"
)

func WsHandler(ctx *gin.Context) {
	Wss.ServeWs(ctx.Writer, ctx.Request)
}

func WssPush(msg gowss.MsgBody) {
	//除了广播到前端外，还需要推送一份到k线计算
	if Wss != nil {
		Wss.Broadcast <- msg
	}
}
