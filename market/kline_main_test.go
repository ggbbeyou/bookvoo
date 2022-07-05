package market

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/go-redis/redis/v8"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yzimhao/bookvoo/market/models"
	"github.com/yzimhao/utilgo"
)

var (
	symbol = "eurusdtest"

	ks *kdataHandler
)

func init() {
	rdc = redis.NewClient(&redis.Options{
		Addr:     ":6379",
		DB:       0,
		Password: "",
	})

	models.InitDbEngine(utilgo.ViperInit("../config.toml"))

	rdc.FlushDB(context.Background())
	deleteTestTable()

	all := models.Periods()
	need := []string{}
	for _, v := range all {
		need = append(need, string(v))
	}

	ks = NewKdataHandler(rdc, need)

}

func deleteTestTable() {
	db := models.DbEngine().NewSession()
	defer db.Close()

	type tables struct {
		TableName string
	}

	rows := []tables{}
	db.Table("information_schema.tables").Find(&rows)
	for _, a := range rows {
		if strings.Contains(a.TableName, symbol) {
			db.DropIndexes(a.TableName)
			db.DropTable(a.TableName)
			models.DeleteTableMapCache()
		}
	}
}

func TestNewKLine(t *testing.T) {
	deleteTestTable()

	Convey("1分钟内的成交记录 k线结果", t, func() {

		logs := []models.TradeLog{
			models.PushTradeLog(symbol, models.ParseTime("2022-06-19 22:04:00"), "askid", "bidid", "1.00", "100", "100"),
			models.PushTradeLog(symbol, models.ParseTime("2022-06-19 22:04:13"), "askid1", "bidid1", "3.00", "100", "300"),
			models.PushTradeLog(symbol, models.ParseTime("2022-06-19 22:04:13"), "askid2", "bidid2", "2.00", "100", "200"),
			models.PushTradeLog(symbol, models.ParseTime("2022-06-19 22:04:30"), "askid3", "bidid3", "3.00", "100", "300"),
		}

		ks.WaitGroupAdd(len(ks.needPeriod) * len(logs))
		for _, item := range logs {
			ks.InputTradeLog <- item
		}

		ks.wg.Wait()
		db := models.DbEngine().NewSession()
		defer db.Close()

		for _, period := range models.Periods() {
			obj := models.NewKline(symbol, period)
			table := obj.TableName()

			rows := []models.Kline{}
			db.Table(table).Find(&rows)

			So(len(rows), ShouldBeGreaterThan, 0)

			switch period {
			case models.PERIOD_M1:
				Convey(fmt.Sprintf("%s %s 开盘价", symbol, period), func() {
					So(rows[0].Open, ShouldEqual, "1.00000000000000000000")
				})
				Convey(fmt.Sprintf("%s %s 最高价", symbol, period), func() {
					So(rows[0].High, ShouldEqual, "3.00000000000000000000")
				})
				Convey(fmt.Sprintf("%s %s 最低价", symbol, period), func() {
					So(rows[0].Low, ShouldEqual, "1.00000000000000000000")
				})
				Convey(fmt.Sprintf("%s %s 成交量", symbol, period), func() {
					So(rows[0].Volume, ShouldEqual, "400.00000000000000000000")
				})
				Convey(fmt.Sprintf("%s %s 成交额", symbol, period), func() {
					So(rows[0].Amount, ShouldEqual, "900.00000000000000000000")
				})
				Convey(fmt.Sprintf("%s %s 成交次数", symbol, period), func() {
					So(rows[0].TradeCnt, ShouldEqual, 4)
				})
				Convey(fmt.Sprintf("%s %s 开盘时间", symbol, period), func() {
					So(rows[0].OpenAt, ShouldEqual, models.ParseTime("2022-06-19 22:04:00"))
				})
				Convey(fmt.Sprintf("%s %s 收盘时间", symbol, period), func() {
					So(rows[0].CloseAt, ShouldEqual, models.ParseTime("2022-06-19 22:04:59"))
				})

			case models.PERIOD_D1:
				Convey(fmt.Sprintf("%s %s 开盘价", symbol, period), func() {
					So(rows[0].Open, ShouldEqual, "1.00000000000000000000")
				})
				Convey(fmt.Sprintf("%s %s 最高价", symbol, period), func() {
					So(rows[0].High, ShouldEqual, "3.00000000000000000000")
				})
				Convey(fmt.Sprintf("%s %s 最低价", symbol, period), func() {
					So(rows[0].Low, ShouldEqual, "1.00000000000000000000")
				})
				Convey(fmt.Sprintf("%s %s 成交量", symbol, period), func() {
					So(rows[0].Volume, ShouldEqual, "400.00000000000000000000")
				})
				Convey(fmt.Sprintf("%s %s 成交额", symbol, period), func() {
					So(rows[0].Amount, ShouldEqual, "900.00000000000000000000")
				})
				Convey(fmt.Sprintf("%s %s 成交次数", symbol, period), func() {
					So(rows[0].TradeCnt, ShouldEqual, 4)
				})
				Convey(fmt.Sprintf("%s %s 开盘时间", symbol, period), func() {
					So(rows[0].OpenAt, ShouldEqual, models.ParseTime("2022-06-19 00:00:00"))
				})
				Convey(fmt.Sprintf("%s %s 收盘时间", symbol, period), func() {
					So(rows[0].CloseAt, ShouldEqual, models.ParseTime("2022-06-19 23:59:59"))
				})

			default:
			}
		}

		//测完一个清掉redis
		rdc.FlushDB(context.Background())
	})

}

// func BenchmarkNewKLine(b *testing.B) {

// 	// s1 := rand.NewSource(time.Now().UnixNano())
// 	// r1 := rand.New(s1)

// 	// for i := 0; i < b.N; i++ {
// 	// 	ask := uuid.NewString()
// 	// 	bid := uuid.NewString()

// 	// 	tl := models.PushTradeLog("btcusdt", time.Now(), ask, bid, fmt.Sprintf("%d", r1.Intn(100)), "200")
// 	// 	ChKLine <- tl
// 	// }
// }

// func parseT(tt string) time.Time {
// 	t, _ := time.ParseInLocation("2006-01-02 15:04:05", tt, time.Local)
// 	return t
// }
