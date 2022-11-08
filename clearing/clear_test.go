package clearing

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"github.com/yzimhao/bookvoo/base"
	"github.com/yzimhao/bookvoo/common"
	"github.com/yzimhao/bookvoo/user/assets"
	"github.com/yzimhao/bookvoo/user/orders"
	"github.com/yzimhao/trading_engine"
	"github.com/yzimhao/utilgo"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	test_symbol       = "ethusd"
	test_user1  int64 = 101
	test_user2  int64 = 102
)

func init() {

	utilgo.ViperInit("../config.toml")
	db_engine = common.Default_db()
	redis_conn := common.Default_redis()

	Init(db_engine, redis_conn)
	db_engine.ShowSQL(true)
	deleteTestTable()

	base.Init(db_engine, redis_conn)
	assets.Init(db_engine, redis_conn)
	orders.Init(db_engine, redis_conn)
}

func deleteTestTable() {
	table1 := orders.TradeRecord{Symbol: test_symbol}
	table2 := orders.GetOrderTableName(test_symbol)
	db_engine.DropTables(
		table1.TableName(),
		table2,
		new(orders.UnfinishedOrder),
	)
}

func Test_Main(t *testing.T) {
	defer func() {
		deleteTestTable()
	}()
	Convey("不同账号交易 限价单的结算", t, func() {
		buy, err := orders.NewLimitOrder(test_user1, test_symbol, orders.OrderSideBuy, "1.00", "10")
		So(err, ShouldBeNil)
		So(buy.OrderId, ShouldStartWith, "B")

		sell, err := orders.NewLimitOrder(test_user2, test_symbol, orders.OrderSideSell, "1.00", "10")
		So(err, ShouldBeNil)
		So(sell.OrderId, ShouldStartWith, "A")

		tr := trading_engine.TradeResult{
			Symbol:        test_symbol,
			AskOrderId:    sell.OrderId,
			BidOrderId:    buy.OrderId,
			TradePrice:    decimal.NewFromFloat(1.00),
			TradeQuantity: decimal.NewFromFloat(10),
			TradeTime:     time.Now().Unix(),
		}
		err = NewClearing(tr)
		So(err, ShouldBeNil)
	})

	// Convey("同账号交易 限价单的结算", t, func() {
	// 	//todo 买卖双方为同一个用户时，结算数据会出现脏数据
	// 	buy, err := orders.NewLimitOrder(test_user1, test_symbol, orders.OrderSideBuy, "1.00", "10")
	// 	So(err, ShouldBeNil)
	// 	So(buy.OrderId, ShouldStartWith, "B")

	// 	sell, err := orders.NewLimitOrder(test_user1, test_symbol, orders.OrderSideSell, "1.00", "10")
	// 	So(err, ShouldBeNil)
	// 	So(sell.OrderId, ShouldStartWith, "A")

	// 	tr := trading_engine.TradeResult{
	// 		Symbol:        test_symbol,
	// 		AskOrderId:    sell.OrderId,
	// 		BidOrderId:    buy.OrderId,
	// 		TradePrice:    decimal.NewFromFloat(1.00),
	// 		TradeQuantity: decimal.NewFromFloat(10),
	// 		TradeTime:     time.Now().Unix(),
	// 	}
	// 	err = NewClearing(tr)
	// 	So(err, ShouldBeNil)
	// })
}
