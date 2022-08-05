package market

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/yzimhao/bookvoo/base"
	"github.com/yzimhao/bookvoo/common/types"
	"github.com/yzimhao/bookvoo/market/models"
	"github.com/yzimhao/gowss"
	"github.com/yzimhao/utilgo"
)

type kdataHandler struct {
	InputTradeLog chan models.TradeLog

	wg         sync.WaitGroup
	rdc        *redis.Client
	needPeriod []string
	sync.Mutex
}

func NewKdataHandler(r *redis.Client, needPeriod []string) *kdataHandler {
	ks := &kdataHandler{
		rdc:           r,
		needPeriod:    needPeriod,
		InputTradeLog: make(chan models.TradeLog, 15),
	}

	go func() {
		for {
			if tradeLog, ok := <-ks.InputTradeLog; ok {
				for _, period := range models.Periods() {
					if !utilgo.ArrayIn(string(period), needPeriod) {
						continue
					}
					ks.updateKline(tradeLog, period)

				}
			}
		}
	}()
	return ks
}

func (ks *kdataHandler) RebuildCache() {
	ks.Lock()
	defer ks.Unlock()
	logrus.Info("RebuildCache...")
	//todo 重建缓存

}

func (ks *kdataHandler) NeedPeriods() []string {
	return ks.needPeriod
}

func (ks *kdataHandler) WaitGroupAdd(n int) {
	ks.wg.Add(n)
}

func (ks *kdataHandler) updateKline(tl models.TradeLog, period models.Period) {
	ks.Lock()
	defer ks.wg.Done()

	st, et := tl.GetAt(period)

	key := fmt.Sprintf("kline:%s_%s_%s_%s", tl.Symbol, period, time2str(st, "20060102150405"), time2str(et, "20060102150405"))

	ctx := context.Background()
	raw := ks.rdc.Get(ctx, key).Val()

	lastK := models.NewKline(tl.Symbol, period)
	ok := json.Unmarshal([]byte(raw), &lastK)

	newK := models.NewKline(tl.Symbol, period)
	newK.OpenAt = types.Time(st)
	newK.CloseAt = types.Time(et)

	if ok != nil {
		lastK.Open = tl.Price
		lastK.High = tl.Price
		lastK.Low = tl.Price
		lastK.Close = tl.Price
		lastK.Volume = "0"
		lastK.Amount = "0"
		lastK.TradeCnt = 0
	}

	{
		newK.Open = lastK.Open
		newK.High = func() string {
			d1 := str2decimal(tl.Price)
			d2 := str2decimal(lastK.High)
			return decimal.Max(d1, d2).String()
		}()
		newK.Low = func() string {
			d1 := str2decimal(tl.Price)
			d2 := str2decimal(lastK.Low)
			return decimal.Min(d1, d2).String()
		}()
		newK.Close = tl.Price
		newK.Volume = addStrNum(tl.Quantity, lastK.Volume)
		newK.Amount = addStrNum(tl.Amount, lastK.Amount)
		newK.TradeCnt = lastK.TradeCnt + 1
	}

	expire := et.Sub(time.Now())
	//每一个缓存多存储30天再过期
	ks.rdc.Set(ctx, key, newK.ToJson(), expire+time.Hour*24*30)

	ks.Unlock()

	go func() {
		ChNewKline <- Klinetips{
			Symbol: tl.Symbol,
			Period: string(period),

			OpenAt:  time.Time(newK.OpenAt).Unix(),
			Open:    newK.Open,
			High:    newK.High,
			Low:     newK.Low,
			Close:   newK.Close,
			Volume:  newK.Volume,
			CloseAt: time.Time(newK.CloseAt).Unix(),
			Amount:  newK.Amount,
		}
	}()

	base.WssPush(rdc, gowss.MsgBody{
		To: types.SubscribeKline.Format(map[string]string{
			"symbol": tl.Symbol,
			"period": string(period),
		}),
		Body: newK,
	})

	//todo 异步入库
	err := newK.Save()

	if err != nil {
		logrus.Error(err)
		return
	}
}

func time2str(t time.Time, layout string) string {
	if layout == "" {
		return t.Format("2006-01-02 15:04:05")
	}
	return t.Format(layout)
}

func str2decimal(str string) decimal.Decimal {
	d, _ := decimal.NewFromString(str)
	return d
}

func addStrNum(n1, n2 string) string {
	d1 := str2decimal(n1)
	d2 := str2decimal(n2)
	return d1.Add(d2).String()
}
