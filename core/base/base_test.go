package base

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	. "github.com/smartystreets/goconvey/convey"
	"xorm.io/xorm"
)

func init() {
	driver := "mysql"
	dsn := "root:root@tcp(localhost:13306)/test?charset=utf8&loc=Local"
	logrus.Infof("dsn: %s", dsn)
	conn, err := xorm.NewEngine(driver, dsn)
	if err != nil {
		logrus.Panic(err)
	}
	db_engine = conn
	db_engine.ShowSQL(true)

	db_engine.DropTables(
		new(SymbolInfo),
		new(TradePairOpt),
	)
	SetDbEngine(db_engine)
}

func Test_main(t *testing.T) {
	db := db_engine.NewSession()
	defer db.Close()

	usd_id, eth_id := 0, 0
	Convey("添加测试用资产symbol信息", t, func() {
		usdInfo := SymbolInfo{
			Symbol:       "usd",
			Name:         "美元",
			MinPrecision: 8,
			ShowPrec:     2,
			Standard:     true,
			Status:       statusEnable,
		}
		_, err := db.Insert(&usdInfo)
		So(err, ShouldBeNil)
		usd_id = usdInfo.Id

		ethInfo := SymbolInfo{
			Symbol:       "eth",
			Name:         "以太坊",
			MinPrecision: 18,
			ShowPrec:     8,
			Status:       statusEnable,
		}
		_, err = db.Insert(&ethInfo)

		eth_id = ethInfo.Id
		So(err, ShouldBeNil)
	})

	Convey("添加测试用交易对信息", t, func() {
		id, err := db.Insert(&TradePairOpt{
			Symbol:         "ethusd",
			Name:           "ETHUSD",
			TradeSymbolId:  int(eth_id),
			BaseSymbolId:   int(usd_id),
			PricePrec:      2,
			QtyPrec:        4,
			AllowMinQty:    "0.0001",
			AllowMaxQty:    "100",
			AllowMinAmount: "1",
			AllowMaxAmount: "1000000",
			FeeRate:        "0.001",
			Status:         statusEnable,
		})
		So(err, ShouldBeNil)
		logrus.Info(id)
	})
}
