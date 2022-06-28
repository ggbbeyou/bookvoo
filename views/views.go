package views

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/yzimhao/haoex/tradecore"
	"github.com/yzimhao/haoex/wss"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/yzimhao/trading_engine"
)

var sendMsg chan []byte
var web *gin.Engine

var recentTrade []interface{}
var rdc *redis.Client

func Run(config string, r *gin.Engine) {

	trading_engine.Debug = false
	recentTrade = make([]interface{}, 0)

	rdc = redis.NewClient(&redis.Options{
		Addr:     ":6379",
		DB:       0,
		Password: "",
	})
	setupRouter(r)
}

func setupRouter(router *gin.Engine) {

	router.LoadHTMLGlob("./template/default/*.html")

	sendMsg = make(chan []byte, 100)

	go pushDepth()
	go watchTradeLog()

	router.GET("/api/depth", depth)
	router.GET("/api/trade_log", trade_log)
	router.POST("/api/new_order", newOrder)
	router.POST("/api/cancel_order", cancelOrder)
	router.GET("/api/test_rand", testOrder)

	router.GET("/demo", func(c *gin.Context) {
		c.HTML(200, "demo.html", nil)
	})

	//websocket
	{
		wss.HHub = wss.NewHub()
		go wss.HHub.Run()
		go func() {
			for {
				select {
				case data := <-sendMsg:
					wss.HHub.Send(data)
				default:
					time.Sleep(time.Duration(100) * time.Millisecond)
				}
			}
		}()

		router.GET("/ws", wss.ServeWs)
		router.GET("/pong", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
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
	msg := gin.H{
		"tag":  tag,
		"data": data,
	}
	msgByte, _ := json.Marshal(msg)
	sendMsg <- []byte(msgByte)
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
				sendMessage("trade", relog)

				if len(recentTrade) >= 10 {
					recentTrade = recentTrade[1:]
				}
				recentTrade = append(recentTrade, relog)

				//latest price
				sendMessage("latest_price", gin.H{
					"latest_price": tradecore.DemoPair.Price2String(log.TradePrice),
				})

			}
		case cancelOrderId := <-tradecore.DemoPair.ChCancelResult:
			sendMessage("cancel_order", gin.H{
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

		sendMessage("depth", gin.H{
			"ask": ask,
			"bid": bid,
		})
	}
}

func newOrder(c *gin.Context) {
	type args struct {
		OrderId   string `json:"order_id"`
		OrderType string `json:"order_type"`
		PriceType string `json:"price_type"`
		Price     string `json:"price"`
		Quantity  string `json:"quantity"`
		Amount    string `json:"amount"`
	}

	var param args
	c.BindJSON(&param)

	orderId := uuid.NewString()
	param.OrderId = orderId

	amount := string2decimal(param.Amount)
	price := string2decimal(param.Price)
	quantity := string2decimal(param.Quantity)

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

	go sendMessage("new_order", param)

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

	go sendMessage("cancel_order", param)

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
