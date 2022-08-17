package views

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/yzimhao/bookvoo/base"
	"github.com/yzimhao/bookvoo/common/types"
	"github.com/yzimhao/bookvoo/match"
	"github.com/yzimhao/bookvoo/views/api"
	"github.com/yzimhao/bookvoo/views/pages"
	gowss "github.com/yzimhao/bookvoo/wss"
	"github.com/yzimhao/trading_engine"
)

var (
	rdc *redis.Client
)

func Init(r *redis.Client) {
	rdc = r
}

func Run(r *gin.Engine) {
	logrus.Info("[views] run")
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
			// for symbol, obj := range match.Engine.Symbols {
			// 	ask := obj.GetAskDepth(6)
			// 	bid := obj.GetBidDepth(6)

			// 	base.Wss.Broadcast <- gowss.MsgBody{
			// 		To: types.SubscribeDepth.Format(map[string]string{"symbol": symbol}),
			// 		Body: gin.H{
			// 			"ask": ask,
			// 			"bid": bid,
			// 		},
			// 	}
			// }

			match.Engine.Foreach(func(symbol string, v *trading_engine.TradePair) {
				ask := v.GetAskDepth(6)
				bid := v.GetBidDepth(6)

				base.Wss.Broadcast <- gowss.MsgBody{
					To: types.SubscribeDepth.Format(map[string]string{"symbol": symbol}),
					Body: gin.H{
						"ask": ask,
						"bid": bid,
					},
				}
			})

			time.Sleep(time.Duration(100) * time.Millisecond)
		}
	}()
}

func botNewOrder() {
	go func() {
		for {

			match.Engine.Foreach(func(symbol string, v *trading_engine.TradePair) {
				ask := v.GetAskDepth(10)
				bid := v.GetBidDepth(10)
				//demo模式下自动挂单
				autoDemoDepthData(symbol, ask, bid, v.LatestPrice())
			})
			time.Sleep(time.Duration(30) * time.Second)
		}
	}()
}
