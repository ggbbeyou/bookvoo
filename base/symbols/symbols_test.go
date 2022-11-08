package symbols

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/yzimhao/bookvoo/common"
	"github.com/yzimhao/utilgo"

	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	utilgo.ViperInit("../../config.toml")
	db := common.Default_db()
	db.DropTables(
		new(Symbol),
		new(Pairs),
	)
	Init(db, nil)
}

func Test_main(t *testing.T) {
	db := db_engine.NewSession()
	defer db.Close()

	usd_id, eth_id := 0, 0
	Convey("添加测试用资产symbol信息", t, func() {
		usdInfo := Symbol{
			Symbol:       "usd",
			Name:         "美元",
			MinPrecision: 8,
			ShowPrec:     2,
			Standard:     true,
			Status:       StatusEnable,
		}
		_, err := db.Insert(&usdInfo)
		So(err, ShouldBeNil)
		usd_id = usdInfo.Id

		ethInfo := Symbol{
			Symbol:       "eth",
			Name:         "以太坊",
			MinPrecision: 18,
			ShowPrec:     4,
			Status:       StatusEnable,
		}
		_, err = db.Insert(&ethInfo)

		eth_id = ethInfo.Id
		So(err, ShouldBeNil)
	})

	Convey("添加测试用交易对信息", t, func() {
		id, err := db.Insert(&Pairs{
			Symbol:           "ethusd",
			Name:             "ETHUSD",
			TargetSymbolId:   int(eth_id),
			StandardSymbolId: int(usd_id),
			PricePrec:        2,
			QtyPrec:          4,
			AllowMinQty:      "0.0001",
			AllowMaxQty:      "100",
			AllowMinAmount:   "1",
			AllowMaxAmount:   "1000000",
			FeeRate:          "0.001",
			Status:           StatusEnable,
		})
		So(err, ShouldBeNil)
		logrus.Info(id)
	})
}
