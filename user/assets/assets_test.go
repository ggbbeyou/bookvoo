package assets

import (
	"fmt"
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
		new(Assets),
		new(assetsLog),
		new(assetFreezeRecord),
	)

	db_engine.Sync2(
		new(Assets),
		new(assetsLog),
		new(assetFreezeRecord),
	)
}

func Test_main(t *testing.T) {
	db := db_engine.NewSession()
	defer db.Close()

	Convey("从根账户充值100", t, func() {
		f, err := transfer(db, true, ROOTUSERID, 1, 1, "100", "r0001", "recharge")
		So(err, ShouldBeNil)
		So(f, ShouldBeTrue)
	})

	Convey("冻结用户资产", t, func() {
		f, err := freezeAssets(db, true, 1, 1, "10", "a001", "trade")
		So(err, ShouldBeNil)
		So(f, ShouldBeTrue)
	})

	Convey("冻结负数的资产", t, func() {
		f, err := freezeAssets(db, true, 1, 1, "-10", "a002", "trade")
		So(err, ShouldBeError, fmt.Errorf("freeze amount should be gt zero"))
		So(f, ShouldBeFalse)
	})

	Convey("冻结数量0的资产", t, func() {
		f, err := freezeAssets(db, true, 1, 1, "0", "a003", "trade")
		So(err, ShouldBeError, fmt.Errorf("freeze amount should be gt zero"))
		So(f, ShouldBeFalse)
	})

	Convey("解冻业务订单号", t, func() {
		f, err := unfreezeAssets(db, true, 1, 1, "a001", "10")
		So(err, ShouldBeNil)
		So(f, ShouldBeTrue)
	})

	Convey("解冻业务订单部分金额", t, func() {
		freezeAssets(db, true, 1, 1, "10", "a005", "trade")
		f, err := unfreezeAssets(db, true, 1, 1, "a005", "8")
		So(err, ShouldBeNil)
		So(f, ShouldBeTrue)
	})

	Convey("解冻超过业务订单金额的数量", t, func() {
		freezeAssets(db, true, 1, 1, "10", "a006", "trade")
		f, err := unfreezeAssets(db, true, 1, 1, "a006", "20")
		So(err, ShouldBeError, fmt.Errorf("unfreeze amount must lt freeze amount"))
		So(f, ShouldBeFalse)
	})

	Convey("解冻不存在的业务订单号", t, func() {
		f, err := unfreezeAssets(db, true, 1, 1, "a004", "10")
		So(err, ShouldBeError, fmt.Errorf("not found"))
		So(f, ShouldBeFalse)
	})
}
