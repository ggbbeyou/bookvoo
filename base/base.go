package base

import (
	"github.com/go-redis/redis/v8"
	"github.com/yzimhao/bookvoo/base/symbols"
	"xorm.io/xorm"
)

func Init(db *xorm.Engine, rdc *redis.Client) {
	symbols.Init(db, rdc)
}
