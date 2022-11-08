package clearing

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/goccy/go-json"
	"github.com/sirupsen/logrus"
	"github.com/yzimhao/bookvoo/base"
	"github.com/yzimhao/bookvoo/base/symbols"
	"github.com/yzimhao/bookvoo/common/types"
	"github.com/yzimhao/bookvoo/user/orders"
	gowss "github.com/yzimhao/bookvoo/wss"
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
				logrus.Infof("[clearing] %s ask: %s bid: %s price: %s vol: %s", data.Symbol, data.AskOrderId, data.BidOrderId, data.TradePrice.String(), data.TradeQuantity.String())
				func(res te.TradeResult) {
					err := NewClearing(res)
					if err != nil {
						logrus.Errorf("[clearing] 结算出错 %s %s %s %s", data.Symbol, data.AskOrderId, data.BidOrderId, err)
					}
				}(data)
			}
		}
	}()
}

//结算一条成交记录
func NewClearing(data te.TradeResult) (err error) {
	tradeInfo, err := symbols.GetPairBySymbol(data.Symbol)
	if err != nil {
		return err
	}

	db := db_engine.NewSession()
	defer db.Close()

	err = db.Begin()
	if err != nil {
		return err
	}

	//标记双方订单 防止还未结算完成，就被撤单了
	flag := NewClearingLock(data.AskOrderId, data.BidOrderId)
	flag.Lock()

	defer func() {
		if err != nil {
			logrus.Errorf("[clearing]出现异常 %s %s %s", data.AskOrderId, data.BidOrderId, err)
			db.Rollback()
		} else {
			//正常结算结束，释放掉redis缓存的锁
			flag.UnLock()
			db.Commit()
		}
	}()

	cl := clearing{
		db:     db,
		symbol: data.Symbol,

		symbol_id:          tradeInfo.TargetSymbolId,
		standard_symbol_id: tradeInfo.StandardSymbolId,

		raw: data,

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
		tag := types.SubscribeTradeRecord.Format(map[string]string{"symbol": data.Symbol})
		base.WssPush(gowss.MsgBody{
			To: tag,
			Response: gowss.Response{
				Type: tag,
				Body: map[string]interface{}{
					"price":    tradeInfo.FormatAmount(cl.raw.TradePrice.String()),
					"quantity": tradeInfo.FormatQty(cl.raw.TradeQuantity.String()),
					"amount":   tradeInfo.FormatAmount(cl.raw.TradeAmount.String()),
					"trade_at": data.TradeTime,
				},
			},
		})

		//这份数据传输到k线计算
		ctx := context.Background()
		s, _ := json.Marshal(data)
		rdc.LPush(ctx, types.MarketSubscribe.String(), string(s))
	}
	return nil
}
