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

	usd, eth := int64(0), int64(0)
	Convey("添加测试用资产symbol信息", t, func() {
		usd_id, err := db.Insert(&SymbolInfo{
			Symbol:       "usd",
			Name:         "美元",
			MinPrecision: 8,
			ShowPrec:     2,
			Standard:     true,
			Status:       statusEnable,
		})
		So(err, ShouldBeNil)
		usd = usd_id

		eth, err = db.Insert(&SymbolInfo{
			Symbol:       "eth",
			Name:         "以太坊",
			MinPrecision: 18,
			ShowPrec:     8,
			Status:       statusEnable,
		})
		So(err, ShouldBeNil)
	})

	Convey("添加测试用交易对信息", t, func() {
		_, err := db.Insert(&TradePairOpt{
			Symbol:         "ethusd",
			Name:           "ETHUSD",
			TradeSymbolId:  int(eth),
			BaseSymbolId:   int(usd),
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
	})
}
