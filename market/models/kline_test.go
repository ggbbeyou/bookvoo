package models

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
	conn, err := xorm.NewEngine(driver, dsn)
	if err != nil {
		logrus.Panic(err)
	}
	SetDbEngine(conn)
}

func TestNewKline(t *testing.T) {
	Convey("table name", t, func() {
		a := NewKline("eurusdtest", PERIOD_D1)
		So(a.TableName(), ShouldEqual, "kline_eurusdtest_d1")
	})

}
