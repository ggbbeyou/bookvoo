package orders

import (
	"testing"

	"github.com/yzimhao/bookvoo/base"
	"github.com/yzimhao/bookvoo/common"
	"github.com/yzimhao/bookvoo/user/assets"
	"github.com/yzimhao/utilgo"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	test_symbol       = "ethusd"
	test_user1  int64 = 101
	test_user2  int64 = 102
)

func init() {
	utilgo.ViperInit("../../config.toml")
	db := common.Default_db()
	db.ShowSQL(true)
	Init(db, nil)
	db_engine.ShowSQL(true)

	table := TradeOrder{Symbol: test_symbol}
	table1 := TradeRecord{Symbol: test_symbol}

	db_engine.DropTables(
		new(UnfinishedOrder),
		table,
		table1,
	)

	base.Init(db_engine, nil)
	assets.Init(db_engine, nil)
}

func Test_main(t *testing.T) {
	db := db_engine.NewSession()
	defer db.Close()

	Convey("", t, func() {
		//在结算订单部分测试了这部分下单
	})

	// Convey("限价买单", t, func() {
	// 	order, err := limit_order(test_user1, "ethusd", OrderSideBuy, "1.00", "2")
	// 	So(err, ShouldBeNil)
	// 	So(order.OrderId, ShouldStartWith, "B")
	// })

	// Convey("限价卖单", t, func() {
	// 	order, err := limit_order(test_user1, "ethusd", OrderSideSell, "1.00", "2")
	// 	So(err, ShouldBeNil)
	// 	So(order.OrderId, ShouldStartWith, "A")
	// })

	// Convey("市价按数量", t, func() {
	// 	order, err := market_order_qty(1, "ethusd", OrderSideSell, "1")
	// 	So(err, ShouldBeNil)
	// 	So(order.OrderId, ShouldStartWith, "A")

	// 	order, err = market_order_qty(1, "ethusd", OrderSideBuy, "1")
	// 	So(err, ShouldBeNil)
	// 	So(order.OrderId, ShouldStartWith, "B")

	// 	// assets.UnfreezeAssets(db, true, order.UserId, order.OrderId, "0")
	// })

	// Convey("市价按成交额", t, func() {
	// 	order, err := market_order_amount(1, "ethusd", OrderSideSell, "100.00")
	// 	So(err, ShouldBeNil)
	// 	So(order.OrderId, ShouldStartWith, "A")

	// 	order, err = market_order_amount(1, "ethusd", OrderSideBuy, "100.00")
	// 	So(err, ShouldBeNil)
	// 	So(order.OrderId, ShouldStartWith, "B")
	// })
}
