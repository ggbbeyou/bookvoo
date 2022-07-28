package core

import (
	"github.com/go-redis/redis/v8"
	"github.com/yzimhao/bookvoo/core/base"
	"xorm.io/xorm"
)

func Init(db *xorm.Engine, rdc *redis.Client) {
	base.SetDbEngine(db)
}
