package assets

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_transfer(t *testing.T) {
	Convey("充值", t, func() {

		db := db_engine.NewSession()
		defer db.Close()
		f, err := transfer(db, ROOTUSERID, 1, 1, "100", "r0001", "recharge")
		So(err, ShouldBeNil)
		So(f, ShouldBeTrue)
	})
}
