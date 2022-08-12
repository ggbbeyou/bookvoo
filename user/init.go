package user

import (
	"github.com/yzimhao/bookvoo/user/assets"
	"xorm.io/xorm"
)

var (
	db_engine *xorm.Engine
)

func Run() {
	assets.InitAssetsForDemo(BotUserId, DemoUsdSymbol, "1000000", "R001")
	assets.InitAssetsForDemo(BotUserId, DemoEthSymbol, "1000000", "R001")
}

func SetDbEngine(db *xorm.Engine) {
	db_engine = db
}
