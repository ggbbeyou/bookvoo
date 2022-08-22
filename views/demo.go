package views

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yzimhao/bookvoo/match"
	"github.com/yzimhao/bookvoo/user"
	"github.com/yzimhao/bookvoo/user/orders"
)

func autoDemoDepthData(symbol string, ask, bid [][2]string, latest decimal.Decimal) {

	size := 10
	min_price := decimal.NewFromFloat(5.0)
	if viper.GetString("main.mode") == "demo" {
		if latest.Cmp(min_price) == -1 {
			latest, _ = decimal.NewFromString("10")
		}

		if len(ask) <= size {
			//auto new order
			autoOrder(orders.OrderSideSell, symbol, latest, size-len(ask))
		}
		if len(bid) <= size {
			//auto new order
			autoOrder(orders.OrderSideBuy, symbol, latest, size-len(bid))
		}
	}
}

func autoOrder(side orders.OrderSide, symbol string, price decimal.Decimal, n int) {

	for i := 0; i < n; i++ {
		rand.Seed(time.Now().Unix())
		qty := rand.Float64()

		float := decimal.NewFromFloat(qty)
		if side == orders.OrderSideSell {
			price = price.Add(float)
		} else {
			price = price.Sub(float)
		}

		_pirce := price.String()
		_vol := fmt.Sprintf("%.f", qty*100)

		order, err := orders.NewLimitOrder(user.BotUserId, symbol, side, _pirce, _vol)
		if err != nil {
			logrus.Errorf("[autoOrder] %s price: %s  vol: %s - %s", symbol, _pirce, _vol, err)
			return
		}
		match.Send <- *order
	}
}
