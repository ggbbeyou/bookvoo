package user

import (
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/bookvoo/user/assets"
	"github.com/yzimhao/bookvoo/user/orders"
	"xorm.io/xorm"
)

var (
	db_engine *xorm.Engine
)

func Run(config string, router *gin.Engine) {

}

func SetDbEngine(db *xorm.Engine) {
	db_engine = db
	assets.SetDbEngine(db)
	orders.SetDbEngine(db)
}
