package quotation

import (
	"github.com/go-redis/redis/v8"
	"xorm.io/xorm"
)

var (
	db_engine *xorm.Engine
	rdc       *redis.Client
)

func Init(db *xorm.Engine, r *redis.Client) {
	db_engine = db
	rdc = r
}
