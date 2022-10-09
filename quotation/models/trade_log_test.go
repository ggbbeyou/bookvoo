package models

import (
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"xorm.io/xorm"
)

var symbol = "eurusdtest"

func init() {
	driver := "mysql"
	dsn := "root:root@tcp(localhost:13306)/test?charset=utf8&loc=Local"
	conn, err := xorm.NewEngine(driver, dsn)
	if err != nil {
		logrus.Panic(err)
	}
	SetDbEngine(conn)
}

func Test_ParseAt(t *testing.T) {
	type testStruct struct {
		name   string
		obj    TradeLog
		expect map[Period]time.Time
	}

	examples := []testStruct{
		testStruct{
			name: "成交记录的k线时间节点",
			obj:  PushTradeLog(symbol, ParseTime("2022-06-18 23:18:03"), "askid", "bidid", "1.00", "200", "200"),
			expect: map[Period]time.Time{
				PERIOD_M1:  ParseTime("2022-06-18 23:18:00"),
				PERIOD_M3:  ParseTime("2022-06-18 23:18:00"),
				PERIOD_M5:  ParseTime("2022-06-18 23:15:00"),
				PERIOD_M15: ParseTime("2022-06-18 23:15:00"),
				PERIOD_M30: ParseTime("2022-06-18 23:00:00"),

				PERIOD_H1:  ParseTime("2022-06-18 23:00:00"),
				PERIOD_H2:  ParseTime("2022-06-18 22:00:00"),
				PERIOD_H4:  ParseTime("2022-06-18 20:00:00"),
				PERIOD_H6:  ParseTime("2022-06-18 18:00:00"),
				PERIOD_H8:  ParseTime("2022-06-18 16:00:00"),
				PERIOD_H12: ParseTime("2022-06-18 12:00:00"),

				PERIOD_D1: ParseTime("2022-06-18 00:00:00"),
				PERIOD_D3: ParseTime("2022-06-18 00:00:00"),

				PERIOD_W1: ParseTime("2022-06-13 00:00:00"),

				PERIOD_MN: ParseTime("2022-06-01 00:00:00"),
			},
		},

		testStruct{
			name: "成交记录的k线时间节点",
			obj:  PushTradeLog(symbol, ParseTime("2022-05-13 11:13:03"), "askid1", "bidid1", "1.00", "200", "200"),
			expect: map[Period]time.Time{
				PERIOD_M1:  ParseTime("2022-05-13 11:13:00"),
				PERIOD_M3:  ParseTime("2022-05-13 11:12:00"),
				PERIOD_M5:  ParseTime("2022-05-13 11:10:00"),
				PERIOD_M15: ParseTime("2022-05-13 11:00:00"),
				PERIOD_M30: ParseTime("2022-05-13 11:00:00"),

				PERIOD_H1:  ParseTime("2022-05-13 11:00:00"),
				PERIOD_H2:  ParseTime("2022-05-13 10:00:00"),
				PERIOD_H4:  ParseTime("2022-05-13 08:00:00"),
				PERIOD_H6:  ParseTime("2022-05-13 06:00:00"),
				PERIOD_H8:  ParseTime("2022-05-13 08:00:00"),
				PERIOD_H12: ParseTime("2022-05-13 00:00:00"),

				PERIOD_D1: ParseTime("2022-05-13 00:00:00"),
				PERIOD_D3: ParseTime("2022-05-12 00:00:00"),

				PERIOD_W1: ParseTime("2022-05-09 00:00:00"),

				PERIOD_MN: ParseTime("2022-05-01 00:00:00"),
			},
		},
	}

	for _, item := range examples {
		Convey(item.name, t, func() {
			for k, v := range item.expect {
				st, _ := item.obj.GetAt(k)
				So(st, ShouldEqual, v)
			}
			item.obj.Clean()
		})
	}

}
