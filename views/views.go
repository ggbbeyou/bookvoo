package views

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/yzimhao/bookvoo/views/api"
	"github.com/yzimhao/bookvoo/views/pages"
	"github.com/yzimhao/gowss"
)

var (
	socket *gowss.Hub
	rdc    *redis.Client
)

func Init(r *redis.Client) {
	rdc = r
}

func Run(r *gin.Engine) {
	setupRouter(r)
	go message()
}

func setupRouter(router *gin.Engine) {
	//pages
	pages.SetupRouter(router)
	//api
	api.SetupRouter(router)
	//websocket
	{
		socket = gowss.NewHub()
		router.GET("/ws", func(ctx *gin.Context) {
			socket.ServeWs(ctx.Writer, ctx.Request)
		})
	}
}

func message() {
	for {
		socket.Broadcast <- gowss.MsgBody{
			To:   tag,
			Body: data,
		}
	}
}

// func pubTradeLog(log trading_engine.TradeResult) {
// 	ctx := context.Background()
// 	raw, _ := json.Marshal(log)
// 	fmt.Println(string(raw))
// 	rdc.LPush(ctx, "list:trade_log", raw)
// }

// func watchTradeLog() {
// 	for {

// 		time.Sleep(time.Second * time.Duration(3))
// 		select {
// 		// case nk, ok := <-kline.ChNewKline:
// 		// 	if ok {
// 		// 		sendMessage(fmt.Sprintf("kline.%s.%s", nk.Period, nk.Symbol), nk)
// 		// 	}
// 		case log, ok := <-base.Engine["ethusd"].ChTradeResult:
// 			if ok {
// 				//
// 				pubTradeLog(log)

// 				relog := gin.H{
// 					"TradePrice":    base.Engine["ethusd"].Price2String(log.TradePrice),
// 					"TradeAmount":   base.Engine["ethusd"].Price2String(log.TradeAmount),
// 					"TradeQuantity": base.Engine["ethusd"].Qty2String(log.TradeQuantity),
// 					"TradeTime":     time.Unix(log.TradeTime/1e9, 0),
// 					"AskOrderId":    log.AskOrderId,
// 					"BidOrderId":    log.BidOrderId,
// 				}
// 				sendMessage("trade.ethusd", relog)

// 				if len(recentTrade) >= 10 {
// 					recentTrade = recentTrade[1:]
// 				}
// 				recentTrade = append(recentTrade, relog)

// 				//latest price
// 				sendMessage("latest_price.ethusd", gin.H{
// 					"latest_price": base.Engine["ethusd"].Price2String(log.TradePrice),
// 				})

// 			}
// 		case cancelOrderId := <-base.Engine["ethusd"].ChCancelResult:
// 			sendMessage("cancel_order.ethusd", gin.H{
// 				"OrderId": cancelOrderId,
// 			})
// 		default:
// 			time.Sleep(time.Duration(100) * time.Millisecond)
// 		}

// 	}
// }

// func pushDepth() {
// 	for {

// 		time.Sleep(time.Duration(150) * time.Millisecond)

// 		ask := base.Engine["ethusd"].GetAskDepth(10)
// 		bid := base.Engine["ethusd"].GetBidDepth(10)

// 		sendMessage("depth.ethusd", gin.H{
// 			"ask": ask,
// 			"bid": bid,
// 		})
// 	}
// }

// func testOrder(c *gin.Context) {
// 	op := strings.ToLower(c.Query("op_type"))
// 	if op != "ask" {
// 		op = "bid"
// 	}

// 	func() {
// 		cnt := 10
// 		for i := 0; i < cnt; i++ {
// 			orderId := uuid.NewString()
// 			if op == "ask" {
// 				orderId = fmt.Sprintf("a-%s", orderId)
// 				item := trading_engine.NewAskLimitItem(orderId, randDecimal(20, 50), randDecimal(20, 100), time.Now().UnixNano())
// 				base.Engine["ethusd"].ChNewOrder <- item
// 			} else {
// 				orderId = fmt.Sprintf("b-%s", orderId)
// 				item := trading_engine.NewBidLimitItem(orderId, randDecimal(1, 20), randDecimal(20, 100), time.Now().UnixNano())
// 				base.Engine["ethusd"].ChNewOrder <- item
// 			}

// 		}
// 	}()

// 	c.JSON(200, gin.H{
// 		"ok": true,
// 		"data": gin.H{
// 			"ask_len": base.Engine["ethusd"].AskLen(),
// 			"bid_len": base.Engine["ethusd"].BidLen(),
// 		},
// 	})
// }

// func cancelOrder(c *gin.Context) {
// 	type args struct {
// 		OrderId string `json:"order_id"`
// 	}

// 	var param args
// 	c.BindJSON(&param)

// 	if param.OrderId == "" {
// 		c.Abort()
// 		return
// 	}
// 	if strings.HasPrefix(param.OrderId, "a-") {
// 		base.Engine["ethusd"].CancelOrder(trading_engine.OrderSideSell, param.OrderId)
// 	} else {
// 		base.Engine["ethusd"].CancelOrder(trading_engine.OrderSideBuy, param.OrderId)
// 	}

// 	go sendMessage("cancel_order.ethusd", param)

// 	c.JSON(200, gin.H{
// 		"ok": true,
// 	})
// }

// func string2decimal(a string) decimal.Decimal {
// 	d, _ := decimal.NewFromString(a)
// 	return d
// }

// func randDecimal(min, max int64) decimal.Decimal {
// 	rand.Seed(time.Now().UnixNano())

// 	d := decimal.New(rand.Int63n(max-min)+min, 0)
// 	return d
// }
