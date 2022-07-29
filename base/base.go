package base

import (
	"github.com/go-redis/redis/v8"
	"github.com/yzimhao/bookvoo/base/symbols"
	"github.com/yzimhao/gowss"
	"xorm.io/xorm"
)

var (
	Wss *gowss.Hub
)

func Init(db *xorm.Engine, rdc *redis.Client) {
	symbols.Init(db, rdc)
	Wss = gowss.NewHub()
}
