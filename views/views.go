package views

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/yzimhao/bookvoo/base"
	"github.com/yzimhao/bookvoo/common/types"
	"github.com/yzimhao/bookvoo/match"
	"github.com/yzimhao/bookvoo/views/api"
	"github.com/yzimhao/bookvoo/views/pages"
	"github.com/yzimhao/gowss"
)

var (
	rdc *redis.Client
)

func Init(r *redis.Client) {
	rdc = r
}

func Run(r *gin.Engine) {
	setupRouter(r)
	pushDepth()
	botNewOrder()
}

func setupRouter(router *gin.Engine) {
	//pages
	pages.SetupRouter(router)
	//api
	api.SetupRouter(router)
	//websocket
	{
		router.GET("/ws", func(ctx *gin.Context) {
			base.WsHandler(ctx)
		})
	}
}

func pushDepth() {
	go func() {
		for {
			for symbol, obj := range match.Engine.Symbols {
				ask := obj.GetAskDepth(10)
				bid := obj.GetBidDepth(10)

				base.Wss.Broadcast <- gowss.MsgBody{
					To: types.SubscribeDepth.Format(map[string]string{"symbol": symbol}),
					Body: gin.H{
						"ask": ask,
						"bid": bid,
					},
				}
			}
			time.Sleep(time.Duration(100) * time.Millisecond)
		}
	}()
}

func botNewOrder() {
	go func() {
		for {
			for symbol, obj := range match.Engine.Symbols {
				ask := obj.GetAskDepth(10)
				bid := obj.GetBidDepth(10)
				//demo模式下自动挂单
				autoDemoDepthData(symbol, ask, bid, obj.LatestPrice())

			}
			time.Sleep(time.Duration(500) * time.Second)
		}
	}()
}
