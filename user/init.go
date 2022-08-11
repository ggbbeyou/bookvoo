package user

import (
	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

var (
	db_engine *xorm.Engine
)

func Run(config string, router *gin.Engine) {

}

func SetDbEngine(db *xorm.Engine) {
	db_engine = db
}
