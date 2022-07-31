package types

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_RedisKey(t *testing.T) {
	Convey("key中无参数", t, func() {
		key := WsMessage.Format(nil)
		So(key, ShouldEqual, "message")
	})

	Convey("新订单消息队列 redis key的格式化", t, func() {
		key := NewOrder.Format(map[string]string{"symbol": "ethusd"})
		So(key, ShouldEqual, "order.new.ethusd")
	})

	Convey("成交结果 redis key的格式化", t, func() {
		key := TradeResult.Format(map[string]string{"symbol": "ethusd"})
		So(key, ShouldEqual, "trade.result.ethusd")
	})

}
