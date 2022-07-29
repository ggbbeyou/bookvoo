package types

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Subscribe(t *testing.T) {
	Convey("消息订阅 深度的tag", t, func() {
		key := SubscribeDepth.Format(map[string]string{
			"symbol": "ethusd",
		})
		So(key, ShouldEqual, "depth.ethusd")
	})

	Convey("消息订阅 k线数据tag", t, func() {
		key := SubscribeKline.Format(map[string]string{
			"period": "m1",
			"symbol": "ethusd",
		})
		So(key, ShouldEqual, "kline.m1.ethusd")
	})
}
