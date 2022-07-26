package core

import (
	"github.com/yzimhao/bookvoo/core/base"
	"xorm.io/xorm"
)

func SetDbEngine(db *xorm.Engine) {
	base.SetDbEngine(db)
}
