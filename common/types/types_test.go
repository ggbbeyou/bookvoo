package types

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_RedisKey(t *testing.T) {
	Convey("新订单消息队列 redis key的格式化", t, func() {
		key := NewOrder.Symbol("ethusd")
		So(key, ShouldEqual, "new_order_ethusd")
	})

	Convey("成交结果 redis key的格式化", t, func() {
		key := TradeResult.Symbol("ethusd")
		So(key, ShouldEqual, "trade_result_ethusd")
	})

}
