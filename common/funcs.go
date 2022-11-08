package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"xorm.io/xorm"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type Response struct {
	Ok     int         `json:"ok"`
	Reason string      `json:"reason,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

func ResponseJson(c *gin.Context, ok int, reason string, data interface{}) {
	res := Response{
		Ok:     ok,
		Reason: reason,
		Data:   data,
	}
	c.JSON(http.StatusOK, res)
}

func Success(c *gin.Context, data interface{}) {
	ResponseJson(c, 1, "", data)
}

func Fail(c *gin.Context, reason string) {
	logrus.Debugf("[fail] %s, %s", c.Request.RequestURI, reason)
	ResponseJson(c, 0, reason, nil)
}

func Default_db() *xorm.Engine {
	dsn := viper.GetString("db.dsn")
	driver := viper.GetString("db.driver")
	conn, err := xorm.NewEngine(driver, dsn)
	if err != nil {
		logrus.Panic(err)
	}

	// tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, "ex_")
	// conn.SetTableMapper(tbMapper)

	if viper.GetBool("db.show_sql") {
		conn.ShowSQL(true)
	}
	return conn
}

func Default_redis() *redis.Client {
	rdc := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.host"),
		DB:       viper.GetInt("redis.db"),
		Password: viper.GetString("redis.password"),
	})
	return rdc
}
