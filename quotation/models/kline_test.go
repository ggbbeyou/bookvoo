package models

import (
	"testing"

	"github.com/yzimhao/bookvoo/common"
	"github.com/yzimhao/utilgo"

	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	utilgo.ViperInit("../../config.toml")
	db_engine := common.Default_db()
	// redis_conn := common.Default_redis()
	SetDbEngine(db_engine)
}

func TestNewKline(t *testing.T) {
	Convey("table name", t, func() {
		a := NewKline("eurusdtest", PERIOD_D1)
		So(a.TableName(), ShouldEqual, "kline_eurusdtest_d1")
	})

}
