package orders

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/yzimhao/bookvoo/base"
	gowss "github.com/yzimhao/bookvoo/wss"
)

var ChCancel chan TradeOrder

func service() {
	ChCancel = make(chan TradeOrder)
	for {
		if data, ok := <-ChCancel; ok {
			detail, err := cancel_order(data.Symbol, data.OrderId)
			if err == nil {
				base.WssPush(gowss.MsgBody{
					To: fmt.Sprintf("%d", detail.UserId),
					Body: map[string]any{
						"type": "cancel_order",
						"data": data.OrderId,
					},
				})
			} else {
				logrus.Errorf("[orders] service cancel %s %s", data.OrderId, err)
			}
		}
	}
}
