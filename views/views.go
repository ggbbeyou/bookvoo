package views

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/yzimhao/bookvoo/tradecore"
	"github.com/yzimhao/gowss"
	kline "github.com/yzimhao/kline/core"
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
	router.LoadHTMLGlob("./template/default/*.html")
	router.StaticFS("/statics", http.Dir("./template/default/statics"))

	//迁移到别处
	router.GET("/api/depth", depth)
	router.GET("/api/trade_log", trade_log)
	router.GET("/api/test_rand", testOrder)

	api := router.Group("/api/v1/order")

	api.POST("/new", newOrder)
	api.POST("/cancel", cancelOrder)

	//pages
	router.GET("/t/:symbol", func(c *gin.Context) {
		c.HTML(200, "demo.html", gin.H{
			"symbol": c.Param("symbol"),
		})
	})

	//websocket
	{
		socket = gowss.NewHub()
		router.GET("/ws", func(ctx *gin.Context) {
			socket.ServeWs(ctx.Writer, ctx.Request)
		})
	}
}

func depth(c *gin.Context) {
	limit := c.Query("limit")
	limitInt, _ := strconv.Atoi(limit)
	if limitInt <= 0 || limitInt > 100 {
		limitInt = 10
	}
	a := tradecore.DemoPair.GetAskDepth(limitInt)
	b := tradecore.DemoPair.GetBidDepth(limitInt)

	c.JSON(200, gin.H{
		"ask": a,
		"bid": b,
	})
}

func trade_log(c *gin.Context) {
	c.JSON(200, gin.H{
		"ok": true,
		"data": gin.H{
			"latest_price": tradecore.DemoPair.Price2String(tradecore.DemoPair.LatestPrice()),
			"trade_log":    recentTrade,
		},
	})
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
		select {
		case nk, ok := <-kline.ChNewKline:
			if ok {
				sendMessage(fmt.Sprintf("kline.%s.%s", nk.Period, nk.Symbol), nk)
			}
		case log, ok := <-tradecore.DemoPair.ChTradeResult:
			if ok {
				//
				pubTradeLog(log)

				relog := gin.H{
					"TradePrice":    tradecore.DemoPair.Price2String(log.TradePrice),
					"TradeAmount":   tradecore.DemoPair.Price2String(log.TradeAmount),
					"TradeQuantity": tradecore.DemoPair.Qty2String(log.TradeQuantity),
					"TradeTime":     time.Unix(log.TradeTime/1e9, 0),
					"AskOrderId":    log.AskOrderId,
					"BidOrderId":    log.BidOrderId,
				}
				sendMessage("trade.demo", relog)

				if len(recentTrade) >= 10 {
					recentTrade = recentTrade[1:]
				}
				recentTrade = append(recentTrade, relog)

				//latest price
				sendMessage("latest_price.demo", gin.H{
					"latest_price": tradecore.DemoPair.Price2String(log.TradePrice),
				})

			}
		case cancelOrderId := <-tradecore.DemoPair.ChCancelResult:
			sendMessage("cancel_order.demo", gin.H{
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

		ask := tradecore.DemoPair.GetAskDepth(10)
		bid := tradecore.DemoPair.GetBidDepth(10)

		sendMessage("depth.demo", gin.H{
			"ask": ask,
			"bid": bid,
		})
	}
}

func newOrder(c *gin.Context) {
	type args struct {
		OrderId    string    `json:"order_id"`
		OrderType  string    `json:"order_type"`
		PriceType  string    `json:"price_type"`
		Price      string    `json:"price"`
		Quantity   string    `json:"quantity"`
		Amount     string    `json:"amount"`
		CreateTime time.Time `json:"create_time"`
	}

	var param args
	c.BindJSON(&param)

	orderId := uuid.NewString()
	param.OrderId = orderId

	amount := string2decimal(param.Amount)
	price := string2decimal(param.Price)
	quantity := string2decimal(param.Quantity)
	param.CreateTime = time.Now()

	var pt trading_engine.PriceType
	if param.PriceType == "market" {
		param.Price = "0"
		pt = trading_engine.PriceTypeMarket
		if param.Amount != "" {
			pt = trading_engine.PriceTypeMarketAmount
			//市价按成交金额卖出时，默认持有该资产1000个
			param.Quantity = "100"
			if amount.Cmp(decimal.NewFromFloat(100000000)) > 0 || amount.Cmp(decimal.Zero) <= 0 {
				c.JSON(200, gin.H{
					"ok":    false,
					"error": "金额必须大于0，且不能超过 100000000",
				})
				return
			}

		} else if param.Quantity != "" {
			pt = trading_engine.PriceTypeMarketQuantity
			//市价按数量买入资产时，需要用户账户所有可用资产数量，测试默认100块
			param.Amount = "100"
			if quantity.Cmp(decimal.NewFromFloat(100000000)) > 0 || quantity.Cmp(decimal.Zero) <= 0 {
				c.JSON(200, gin.H{
					"ok":    false,
					"error": "数量必须大于0，且不能超过 100000000",
				})
				return
			}
		}
	} else {
		pt = trading_engine.PriceTypeLimit
		param.Amount = "0"
		if price.Cmp(decimal.NewFromFloat(100000000)) > 0 || price.Cmp(decimal.Zero) < 0 {
			c.JSON(200, gin.H{
				"ok":    false,
				"error": "价格必须大于等于0，且不能超过 100000000",
			})
			return
		}
		if quantity.Cmp(decimal.NewFromFloat(100000000)) > 0 || quantity.Cmp(decimal.Zero) <= 0 {
			c.JSON(200, gin.H{
				"ok":    false,
				"error": "数量必须大于0，且不能超过 100000000",
			})
			return
		}
	}

	if strings.ToLower(param.OrderType) == "ask" {
		param.OrderId = fmt.Sprintf("a-%s", orderId)
		item := trading_engine.NewAskItem(pt, param.OrderId, string2decimal(param.Price), string2decimal(param.Quantity), string2decimal(param.Amount), time.Now().UnixNano())
		tradecore.DemoPair.ChNewOrder <- item

	} else {
		param.OrderId = fmt.Sprintf("b-%s", orderId)
		item := trading_engine.NewBidItem(pt, param.OrderId, string2decimal(param.Price), string2decimal(param.Quantity), string2decimal(param.Amount), time.Now().UnixNano())
		tradecore.DemoPair.ChNewOrder <- item
	}

	go sendMessage("new_order.demo", param)

	c.JSON(200, gin.H{
		"ok": true,
		"data": gin.H{
			"ask_len": tradecore.DemoPair.AskLen(),
			"bid_len": tradecore.DemoPair.BidLen(),
		},
	})
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
				tradecore.DemoPair.ChNewOrder <- item
			} else {
				orderId = fmt.Sprintf("b-%s", orderId)
				item := trading_engine.NewBidLimitItem(orderId, randDecimal(1, 20), randDecimal(20, 100), time.Now().UnixNano())
				tradecore.DemoPair.ChNewOrder <- item
			}

		}
	}()

	c.JSON(200, gin.H{
		"ok": true,
		"data": gin.H{
			"ask_len": tradecore.DemoPair.AskLen(),
			"bid_len": tradecore.DemoPair.BidLen(),
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
		tradecore.DemoPair.CancelOrder(trading_engine.OrderSideSell, param.OrderId)
	} else {
		tradecore.DemoPair.CancelOrder(trading_engine.OrderSideBuy, param.OrderId)
	}

	go sendMessage("cancel_order.demo", param)

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
