package models

import (
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"xorm.io/xorm"
)

var (
	engine           *xorm.Engine
	klineTableMap    sync.Map
	tradeLogTableMap sync.Map
)

func InitDbEngine(opt *viper.Viper) {
	opt.SetDefault("kline.db.driver", "mysql")
	opt.SetDefault("kline.db.dsn", "root:root@tcp(localhost:3306)/test?charset=utf8&loc=Local")
	opt.SetDefault("kline.db.show_sql", false)

	dsn := opt.GetString("kline.db.dsn")
	driver := opt.GetString("kline.db.driver")

	conn, err := xorm.NewEngine(driver, dsn)
	if err != nil {
		logrus.Panic(err)
	}
	engine = conn

	if opt.GetBool("kline.db.show_sql") {
		engine.ShowSQL(true)
	}
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
