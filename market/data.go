package market

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yzimhao/bookvoo/common/types"
	"github.com/yzimhao/bookvoo/market/models"
	te "github.com/yzimhao/trading_engine"
)

var (
	kdh        *kdataHandler
	ChNewKline chan Klinetips
)

type Klinetips struct {
	Symbol string `json:"symbol"`
	Period string `json:"period"`

	OpenAt  int64  `json:"open_at"`  //开盘时间
	Open    string `json:"open"`     //开盘价
	High    string `json:"high"`     // 最高价
	Low     string `json:"low"`      //最低价
	Close   string `json:"close"`    //收盘价(当前K线未结束的即为最新价)
	Volume  string `json:"volume"`   //成交量
	CloseAt int64  `json:"close_at"` // 收盘时间
	Amount  string `json:"amount"`   //成交额

}

func Run(router *gin.Engine) {
	logrus.Info("[market] run")
	initConfig()
	setupRouter(router)
	go handleKLDataService()
}

func initConfig() {
	models.SetDbEngine(db_engine)
	ChNewKline = make(chan Klinetips, 1000)
}

func handleKLDataService() {
	kdh = NewKdataHandler(rdc, viper.GetStringSlice("kline.interval"))
	//todo 初始化kline的最新缓存
	kdh.RebuildCache()

	//通过list获取成交记录
	popTradeLog()
}

func popTradeLog() {
	subcribeKey := types.MarketSubscribe
	ctx := context.Background()
	logrus.Infof("subscribe key=%s", subcribeKey)

	for {
		msg := rdc.BRPop(ctx, time.Duration(30)*time.Second, subcribeKey.String()).Val()
		if len(msg) > 1 {
			kdh.WaitGroupAdd(len(kdh.NeedPeriods()) * 1)
			handleData(msg[1])
		}
	}
}

func handleData(msg string) error {
	var tr te.TradeResult
	err := json.Unmarshal([]byte(msg), &tr)
	if err != nil {
		logrus.Errorf("%s, err: %s", msg, err)
		logrus.Errorf("数据解析出错，请参考 %s", "[]")
		return err
	}

	tl := models.TradeLog{
		Symbol:   tr.Symbol,
		At:       time.Unix(tr.TradeTime/1e9, 0),
		Price:    tr.TradePrice.String(),
		Quantity: tr.TradeQuantity.String(),
		Amount:   tr.TradeAmount.String(),
		AskId:    tr.AskOrderId,
		BidId:    tr.BidOrderId,
	}

	kdh.InputTradeLog <- tl
	tl.Save()
	return nil
}
