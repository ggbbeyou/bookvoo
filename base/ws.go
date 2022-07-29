package base

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/yzimhao/bookvoo/common/types"
	"github.com/yzimhao/gowss"
)

func WsHandler(ctx *gin.Context) {
	Wss.ServeWs(ctx.Writer, ctx.Request)
}

func TradeResultPush(rdc *redis.Client, msg gowss.MsgBody) {
	ctx := context.Background()
	rdc.LPush(ctx, types.Message.Format(nil), string(msg.GetBody()))
}
