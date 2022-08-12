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
	if viper.GetString("main.mode") == "demo" {
		if latest.Cmp(decimal.Zero) == 0 {
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

		order, err := orders.NewLimitOrder(user.BotUserId, symbol, side, price.StringFixedBank(4), fmt.Sprintf("%.4f", qty))
		if err != nil {
			logrus.Error(err)
			return
		}
		match.Send <- order
	}
}
