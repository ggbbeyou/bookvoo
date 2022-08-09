package clearings

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/goccy/go-json"
	"github.com/sirupsen/logrus"
	"github.com/yzimhao/bookvoo/base"
	"github.com/yzimhao/bookvoo/base/symbols"
	"github.com/yzimhao/bookvoo/common/types"
	"github.com/yzimhao/bookvoo/user/orders"
	"github.com/yzimhao/gowss"
	te "github.com/yzimhao/trading_engine"
	"xorm.io/xorm"
)

var (
	db_engine *xorm.Engine
	Notify    chan te.TradeResult
	rdc       *redis.Client
)

func Init(db *xorm.Engine, r *redis.Client) {
	db_engine = db
	rdc = r
}

func Run() {
	Notify = make(chan te.TradeResult, 1000)
	go func() {
		for {
			if data, ok := <-Notify; ok {
				go func(res te.TradeResult) {
					err := NewClearing(res)
					if err != nil {
						logrus.Errorf("[clearings] %s", err)
					}
				}(data)
			}
		}
	}()
}

//结算一条成交记录
func NewClearing(data te.TradeResult) (err error) {
	logrus.Infof("[clearings] %#v", data)

	tradeInfo, err := symbols.GetExchangeBySymbol(data.Symbol)
	if err != nil {
		return err
	}

	db := db_engine.NewSession()
	defer db.Close()

	err = db.Begin()
	if err != nil {
		return err
	}

	//todo lock 双方订单

	defer func() {
		if err != nil {
			db.Rollback()
		} else {
			db.Commit()
		}
	}()

	cl := clearing{
		db:     db,
		symbol: data.Symbol,

		symbol_id:          tradeInfo.TargetSymbolId,
		standard_symbol_id: tradeInfo.StandardSymbolId,

		ask_order_id: data.AskOrderId,
		bid_order_id: data.BidOrderId,
		trade_price:  data.TradePrice,
		trade_qty:    data.TradeQuantity,
		trade_amount: data.TradePrice.Mul(data.TradeQuantity),

		ask:    new(orders.TradeOrder),
		bid:    new(orders.TradeOrder),
		record: new(orders.TradeRecord),
	}
	//检查双方订单状态
	err = cl.check()
	if err != nil {
		return err
	}

	//写成交日志
	err = cl.tradeRecord()
	if err != nil {
		return err
	}

	//修改买方订单信息
	err = cl.updateBid()
	if err != nil {
		return err
	}

	//修改卖方订单信息
	err = cl.updateAsk()
	if err != nil {
		return err
	}

	//结算三方资产
	err = cl.transfer()
	if err != nil {
		return err
	}

	//成交记录推送到下游
	if rdc != nil {
		base.WssPush(rdc, gowss.MsgBody{
			To: types.SubscribeTradeRecord.Format(map[string]string{"symbol": data.Symbol}),
			Body: map[string]interface{}{
				"price":    tradeInfo.FormatAmount(cl.trade_price.String()),
				"quantity": tradeInfo.FormatQty(cl.trade_qty.String()),
				"amount":   tradeInfo.FormatAmount(cl.trade_amount.String()),
				"trade_at": data.TradeTime,
			},
		})

		//这份数据传输到k线计算
		ctx := context.Background()
		s, _ := json.Marshal(data)
		rdc.LPush(ctx, types.MarketSubscribe.String(), string(s))
	}
	return nil
}
