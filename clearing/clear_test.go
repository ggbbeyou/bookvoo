package clearing

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"github.com/sirupsen/logrus"
	"github.com/yzimhao/bookvoo/core/base"
	"github.com/yzimhao/bookvoo/user/assets"
	"github.com/yzimhao/bookvoo/user/orders"
	"xorm.io/xorm"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	test_symbol       = "ethusd"
	test_user1  int64 = 101
	test_user2  int64 = 102
)

func init() {
	driver := "mysql"
	dsn := "root:root@tcp(localhost:13306)/test?charset=utf8&loc=Local"

	conn, err := xorm.NewEngine(driver, dsn)
	if err != nil {
		logrus.Panic(err)
	}
	db_engine = conn
	db_engine.ShowSQL(true)

	table1 := orders.TradeRecord{Symbol: test_symbol}
	db_engine.DropTables(
		table1,
	)

	SetDbEngine(db_engine)
	base.SetDbEngine(db_engine)
	assets.SetDbEngine(db_engine)
	orders.SetDbEngine(db_engine)
}

func Test_Main(t *testing.T) {
	Convey("限价单的结算", t, func() {
		buy, err := orders.NewLimitOrder(test_user1, test_symbol, orders.OrderSideBuy, "1.00", "2")
		So(err, ShouldBeNil)
		So(buy.OrderId, ShouldStartWith, "B")

		sell, err := orders.NewLimitOrder(test_user2, test_symbol, orders.OrderSideSell, "1.00", "2")
		So(err, ShouldBeNil)
		So(sell.OrderId, ShouldStartWith, "A")

		// err = NewClearing(test_symbol, sell.OrderId, buy.OrderId, "1", "2")
		// So(err, ShouldBeNil)
	})
}
