package orders

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"xorm.io/xorm"
)

var (
	db_engine *xorm.Engine
)

func SetDbEngine(db *xorm.Engine) {
	db_engine = db

	err := db_engine.Sync2(
		new(UnfinishedOrder),
	)
	if err != nil {
		logrus.Errorf("sync2: %s", err)
	}

}

func order_id_by_side(side orderSide) string {
	if side == orderSideAsk {
		return make_order_id("A")
	} else {
		return make_order_id("B")
	}
}

func make_order_id(pre string) string {
	pre = strings.ToUpper(pre)
	s := time.Now().Format("060102150405")
	ns := time.Now().Nanosecond()
	rand.Seed(time.Now().UnixNano())
	rn := rand.Intn(99)
	return fmt.Sprintf("%s%s%09d%02d", pre, s, ns, rn)
}
