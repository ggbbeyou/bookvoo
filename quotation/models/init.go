package models

import (
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"xorm.io/xorm"
)

var (
	engine           *xorm.Engine
	klineTableMap    sync.Map
	tradeLogTableMap sync.Map
)

func SetDbEngine(db *xorm.Engine) {
	engine = db
}

func DbEngine() *xorm.Engine {
	return engine
}

func ParseTime(tt string) time.Time {
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", tt, time.Local)
	return t
}

func DeleteTableMapCache() {
	klineTableMap.Range(func(key, value any) bool {
		klineTableMap.Delete(key)
		return true
	})

	tradeLogTableMap.Range(func(key, value any) bool {
		tradeLogTableMap.Delete(key)
		return true
	})
}
