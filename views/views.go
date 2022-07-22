package views

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/yzimhao/bookvoo/core/base"
	"github.com/yzimhao/bookvoo/views/api"
	"github.com/yzimhao/bookvoo/views/pages"
	"github.com/yzimhao/gowss"
	"github.com/yzimhao/utilgo"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/yzimhao/trading_engine"
)

var (
	socket *gowss.Hub

	recentTrade []interface{}
	rdc         *redis.Client
	conf        *viper.Viper
)

func Run(config string, r *gin.Engine) {

	conf = utilgo.ViperInit(config)

	trading_engine.Debug = false
	recentTrade = make([]interface{}, 0)

	rdc = redis.NewClient(&redis.Options{
		Addr:     conf.GetString("main.redis.host"),
		DB:       conf.GetInt("main.redis.db"),
		Password: conf.GetString("main.redis.password"),
	})
	setupRouter(r)

	go pushDepth()
	go watchTradeLog()
}

func setupRouter(router *gin.Engine) {
	//pages
	pages.SetupRouter(router)
	api.SetupRouter(router)

	//websocket
	{
		socket = gowss.NewHub()
		router.GET("/ws", func(ctx *gin.Context) {
			socket.ServeWs(ctx.Writer, ctx.Request)
		})
	}
}

func sendMessage(tag string, data interface{}) {

	socket.Broadcast <- gowss.MsgBody{
		To:   tag,
		Body: data,
	}
}

func pubTradeLog(log trading_engine.TradeResult) {
	ctx := context.Background()
	raw, _ := json.Marshal(log)
	fmt.Println(string(raw))
	rdc.LPush(ctx, "list:trade_log", raw)
}

func watchTradeLog() {
	for {

		time.Sleep(time.Second * time.Duration(3))
		select {
		// case nk, ok := <-kline.ChNewKline:
		// 	if ok {
		// 		sendMessage(fmt.Sprintf("kline.%s.%s", nk.Period, nk.Symbol), nk)
		// 	}
		case log, ok := <-base.MatchingEngine["ethusd"].ChTradeResult:
			if ok {
				//
				pubTradeLog(log)

				relog := gin.H{
					"TradePrice":    base.MatchingEngine["ethusd"].Price2String(log.TradePrice),
					"TradeAmount":   base.MatchingEngine["ethusd"].Price2String(log.TradeAmount),
					"TradeQuantity": base.MatchingEngine["ethusd"].Qty2String(log.TradeQuantity),
					"TradeTime":     time.Unix(log.TradeTime/1e9, 0),
					"AskOrderId":    log.AskOrderId,
					"BidOrderId":    log.BidOrderId,
				}
				sendMessage("trade.ethusd", relog)

				if len(recentTrade) >= 10 {
					recentTrade = recentTrade[1:]
				}
				recentTrade = append(recentTrade, relog)

				//latest price
				sendMessage("latest_price.ethusd", gin.H{
					"latest_price": base.MatchingEngine["ethusd"].Price2String(log.TradePrice),
				})

			}
		case cancelOrderId := <-base.MatchingEngine["ethusd"].ChCancelResult:
			sendMessage("cancel_order.ethusd", gin.H{
				"OrderId": cancelOrderId,
			})
		default:
			time.Sleep(time.Duration(100) * time.Millisecond)
		}

	}
}

func pushDepth() {
	for {

		time.Sleep(time.Duration(150) * time.Millisecond)

		ask := base.MatchingEngine["ethusd"].GetAskDepth(10)
		bid := base.MatchingEngine["ethusd"].GetBidDepth(10)

		sendMessage("depth.ethusd", gin.H{
			"ask": ask,
			"bid": bid,
		})
	}
}

func newOrder(c *gin.Context) {
	// symbol := "ethusd"

	// amount := string2decimal(param.Amount)
	// price := string2decimal(param.Price)
	// quantity := string2decimal(param.Quantity)
	// param.CreateTime = time.Now()

	// // var pt trading_engine.PriceType
	// if param.PriceType == "market" {
	// 	param.Price = "0"
	// 	// pt = trading_engine.PriceTypeMarket
	// 	if param.Amount != "" {
	// 		// pt = trading_engine.PriceTypeMarketAmount
	// 		//市价按成交金额卖出时，默认持有该资产1000个
	// 		param.Quantity = "100"
	// 		if amount.Cmp(decimal.NewFromFloat(100000000)) > 0 || amount.Cmp(decimal.Zero) <= 0 {
	// 			c.JSON(200, gin.H{
	// 				"ok":    false,
	// 				"error": "金额必须大于0，且不能超过 100000000",
	// 			})
	// 			return
	// 		}

	// 	} else if param.Quantity != "" {
	// 		// pt = trading_engine.PriceTypeMarketQuantity
	// 		//市价按数量买入资产时，需要用户账户所有可用资产数量，测试默认100块
	// 		param.Amount = "100"
	// 		if quantity.Cmp(decimal.NewFromFloat(100000000)) > 0 || quantity.Cmp(decimal.Zero) <= 0 {
	// 			c.JSON(200, gin.H{
	// 				"ok":    false,
	// 				"error": "数量必须大于0，且不能超过 100000000",
	// 			})
	// 			return
	// 		}
	// 	}
	// } else {
	// 	// pt = trading_engine.PriceTypeLimit
	// 	param.Amount = "0"
	// 	if price.Cmp(decimal.NewFromFloat(100000000)) > 0 || price.Cmp(decimal.Zero) < 0 {
	// 		c.JSON(200, gin.H{
	// 			"ok":    false,
	// 			"error": "价格必须大于等于0，且不能超过 100000000",
	// 		})
	// 		return
	// 	}
	// 	if quantity.Cmp(decimal.NewFromFloat(100000000)) > 0 || quantity.Cmp(decimal.Zero) <= 0 {
	// 		c.JSON(200, gin.H{
	// 			"ok":    false,
	// 			"error": "数量必须大于0，且不能超过 100000000",
	// 		})
	// 		return
	// 	}
	// }

	// if strings.ToLower(param.OrderType) == "ask" {
	// 	order, err := orders.NewLimitOrder(1, symbol, orders.OrderSideAsk, param.Price, param.Quantity)
	// 	if err != nil {
	// 		c.JSON(200, gin.H{
	// 			"ok":    false,
	// 			"error": err.Error(),
	// 		})
	// 		return
	// 	}
	// 	item := trading_engine.NewAskLimitItem(order.OrderId, string2decimal(order.Price), string2decimal(order.Quantity), order.CreateTime)
	// 	base.MatchingEngine[symbol].ChNewOrder <- item

	// } else {
	// 	order, err := orders.NewLimitOrder(1, symbol, orders.OrderSideBid, param.Price, param.Quantity)
	// 	if err != nil {
	// 		c.JSON(200, gin.H{
	// 			"ok":    false,
	// 			"error": err.Error(),
	// 		})
	// 		return
	// 	}
	// 	item := trading_engine.NewAskLimitItem(order.OrderId, string2decimal(order.Price), string2decimal(order.Quantity), order.CreateTime)
	// 	base.MatchingEngine[symbol].ChNewOrder <- item
	// }

	// go sendMessage(fmt.Sprintf("new_order.%s", symbol), param)

	// c.JSON(200, gin.H{
	// 	"ok": true,
	// 	"data": gin.H{
	// 		"ask_len": base.MatchingEngine[symbol].AskLen(),
	// 		"bid_len": base.MatchingEngine[symbol].BidLen(),
	// 	},
	// })
}

func testOrder(c *gin.Context) {
	op := strings.ToLower(c.Query("op_type"))
	if op != "ask" {
		op = "bid"
	}

	func() {
		cnt := 10
		for i := 0; i < cnt; i++ {
			orderId := uuid.NewString()
			if op == "ask" {
				orderId = fmt.Sprintf("a-%s", orderId)
				item := trading_engine.NewAskLimitItem(orderId, randDecimal(20, 50), randDecimal(20, 100), time.Now().UnixNano())
				base.MatchingEngine["ethusd"].ChNewOrder <- item
			} else {
				orderId = fmt.Sprintf("b-%s", orderId)
				item := trading_engine.NewBidLimitItem(orderId, randDecimal(1, 20), randDecimal(20, 100), time.Now().UnixNano())
				base.MatchingEngine["ethusd"].ChNewOrder <- item
			}

		}
	}()

	c.JSON(200, gin.H{
		"ok": true,
		"data": gin.H{
			"ask_len": base.MatchingEngine["ethusd"].AskLen(),
			"bid_len": base.MatchingEngine["ethusd"].BidLen(),
		},
	})
}

func cancelOrder(c *gin.Context) {
	type args struct {
		OrderId string `json:"order_id"`
	}

	var param args
	c.BindJSON(&param)

	if param.OrderId == "" {
		c.Abort()
		return
	}
	if strings.HasPrefix(param.OrderId, "a-") {
		base.MatchingEngine["ethusd"].CancelOrder(trading_engine.OrderSideSell, param.OrderId)
	} else {
		base.MatchingEngine["ethusd"].CancelOrder(trading_engine.OrderSideBuy, param.OrderId)
	}

	go sendMessage("cancel_order.ethusd", param)

	c.JSON(200, gin.H{
		"ok": true,
	})
}

func string2decimal(a string) decimal.Decimal {
	d, _ := decimal.NewFromString(a)
	return d
}

func randDecimal(min, max int64) decimal.Decimal {
	rand.Seed(time.Now().UnixNano())

	d := decimal.New(rand.Int63n(max-min)+min, 0)
	return d
}
