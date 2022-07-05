package models

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/yzimhao/utilgo"
)

func init() {
	InitDbEngine(utilgo.ViperInit("../config.toml"))
}

func TestNewKline(t *testing.T) {
	Convey("table name", t, func() {
		a := NewKline("eurusdtest", PERIOD_D1)
		So(a.TableName(), ShouldEqual, "kline_eurusdtest_d1")

	})

}
